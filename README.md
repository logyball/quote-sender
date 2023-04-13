# Quote/Trivia and Cat Messenger

This is a dumb program that uses golang to fetch the quote of the day from `http://quote.rest/` as well as some trivia from [API Ninjas](https://api-ninjas.com/api/trivia) and then uses twilio to send it to my phone.  The real value of it is that it's a simple thing that can be served up via k8s in a cronjob.

Starting in v0.10, I added the [CatAPI](https://thecatapi.com/) as bonus picture inside the message. Starting in v1.0, I added the [Dog API](https://dog.ceo/dog-api/) to show on Fridays, as well as added support for way more cats via filetypes supported natively by Twilio. As of v2.0, I've integrated into the awesome service [ntfy.sh](https://ntfy.sh/) for quick notifications on my phone if there are any errors. As of v2.0.10, Trivia Tuesday is implemented.

<img src="./imgs/example_message.jpeg" width="360" height="640" alt="Example Message"/>

## Requirements

1. Go v1.19+
2. Working Kubernetes Cluster w/ the [SealedSecrets CRD](https://github.com/bitnami-labs/sealed-secrets) installed
3. [Twilio Subscription and API Key](https://www.twilio.com/docs/usage/api)
4. [Quote API Key](https://quotes.rest/)
5. [Cat API Key](https://thecatapi.com/signup)
6. [API Ninja Key](https://api-ninjas.com/api)
7. (optional) [ntfy.sh](https://ntfy.sh) for error reporting
8. (optional) a prometheus pushgateway instance

## Configuration

Configure the following environment variables in `ops/*/cronJob.yml` to suit your needs:

1. TZ - your time zone
2. (if error reporting) ERR_NOTIFICATION_TOPIC - your ntfy topic
3. (if not error reporting) ERROR_REPORT_ENABLED - set to false
4. (if not pushing to prometheus) PUSH_TO_PROMETHEUS - set to false
5. (if pushing to prometheus) PROMETHEUS_GATEWAY_URI - where your prometheus push gateway listens

### Building

`go build`

### Testing

`go test -v ./...`

### K8s

The real value of this is deploying it to kubernetes as a cronjob and then looking to add prometheus push gateway metrics.

#### Kustomization

This was my first foray into kustomizing.  I created a base, overlay, and two separate secrets such that deployment of all resources can be done together as well as having separate **sealed** secrets (in this case, phone numbers) for dev and prod.

#### Gotchas with SealedSecrets and Kustomization

Because SealedSecrets use the name and namespace and all other data about the secret to create the hash, you have to make a complete replica of the secret using the "-dev" suffix, but then set the kustomization selector to the original label of the secrets.

### Applying with Kustomize

dev:
`kubectl apply -k ops/overlays/dev`

prod:
`kubectl apply -k ops/base`

### Secrets

The phone numbers to text, the cat api key (get yours [here](https://thecatapi.com/signup)) as well as the twilio api key are obscured as secrets.  Set them as environment variables for local testing:

```shell
TWILIO_AUTH=api-key CAT_API_KEY=api-key PHONE_NUMBERS=+16666666666,+16666666666 go run .
```

If you want to test locally with error reporting and prometheus metrics, it's a bit more involved:

```shell
TWILIO_AUTH=api-key \
  PHONE_NUMBERS=+16666666666 \
  CAT_API_KEY=api-key \
  QUOTE_API_KEY=api-key \
  API_NINJA_KEY=api-key \
  ERROR_REPORT_ENABLED=true \
  ERR_NOTIFICATION_TOPIC=your_topic \
  PUSH_TO_PROMETHEUS=true \
  PROMETHEUS_GATEWAY_URI=your_gateway_uri \
  go run .
```

### Adding New Phone Numbers

Using `kubeseal` and SealedSecrets, you can add new phone numbers:

```bash
$ kubectl create secret generic phone-numbers -n quotes --from-literal=numbers=+15555555555,+16666666666 --dry-run=client -o yaml | kubeseal --format=yaml -
apiVersion: bitnami.com/v1alpha1
kind: SealedSecret
metadata:
  creationTimestamp: null
  name: phone-numbers
  namespace: quotes
spec:
  encryptedData:
    numbers: 
    <gibberish>
  template:
    metadata:
      creationTimestamp: null
      name: phone-numbers
      namespace: quotes
```

Then apply that to your cluster.

## Vanilla Deployment

I had previously deployed this [to my kubernetes cluster on a bunch of raspberry pis](https://loganballard.com/index.php/2021/05/26/running-kubernetes-on-raspberry-pi/).  However, they've recently stopped working and are impossible to buy new ones.  So I bit the bullet and used Google Cloud Free Tier to host a non-containerized, old-school binary that can be easily `scp`'d into the VM.

In order to replicate that, deploy an Ubunutu VM to GCP (or your cloud provider of choice), and give yourself the [ability to SSH into it via an SSH Keypair](https://cloud.google.com/compute/docs/connect/add-ssh-keys).  Once you've got that, and the external IP of the VM, update your `.env` file with:

- twilio API key
- cat API key
- Quote API Key
- API Ninja Key
- phone numbers to text
- (OPTIONALLY) error reporting stuff
- (OPTIONALLY) prometheus things
- `REMOTE_MACHINE_IP` = your VMs IP
- USER_DIR = your users home directory

then run `make deploy`

```sh
$ make deploy
mkdir -p ./vanilla_deploy/tmp
cat ./vanilla_deploy/env.template \
		... a bunch of juicy stuff ...
scp ./vanilla_deploy/startup.sh REMOTE_IP:/home/Ubuntu/startup.sh
scp ./vanilla_deploy/tmp/.env REMOTE_IP:/home/Ubuntu/.env
ssh REMOTE_IP /home/Ubuntu/startup.sh
Hit:1 http://us-west1.gce.archive.ubuntu.com/ubuntu bionic InRelease
Get:2 http://us-west1.gce.archive.ubuntu.com/ubuntu bionic-updates InRelease [88.7 kB]
Get:3 http://us-west1.gce.archive.ubuntu.com/ubuntu bionic-backports InRelease [83.3 kB]
Hit:4 http://security.ubuntu.com/ubuntu bionic-security InRelease
Fetched 172 kB in 1s (232 kB/s)
Reading package lists...
Building dependency tree...
Reading state information...
5 packages can be upgraded. Run 'apt list --upgradable' to see them.
Reading package lists...
Building dependency tree...
Reading state information...
Package 'golang-go' is not installed, so not removed
The following packages were automatically installed and are no longer required:
  golang-1.10-go golang-1.10-race-detector-runtime golang-1.10-src
  golang-race-detector-runtime golang-src libnuma1 pkg-config
Use 'sudo apt autoremove' to remove them.
0 upgraded, 0 newly installed, 0 to remove and 5 not upgraded.
go mod download
mkdir ./dist
go build -o ./dist/quoteCats
rm -rf ./vanilla_deploy/tmp

```

This command will:

- install `go` on the remote VM
- clone this repo
- build the binary
- install a cron that runs this binary every day at 14:00 UTC
