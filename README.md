# Quote and Cat Messenger

This is a dumb program that uses golang to fetch the quote of the day from `http://quote.rest/` and then uses twilio to send it to my phone.  The real value of it is that it's a simple thing that can be served up via k8s in a cronjob.

Starting in v0.10, I added the [CatAPI](https://thecatapi.com/) as bonus picture inside the message. Starting in v1.0, I added the [Dog API](https://dog.ceo/dog-api/) to show on Fridays, as well as added support for way more cats via filetypes supported natively by Twilio. As of v2.0, I've integrated into the awesome service [ntfy.sh](https://ntfy.sh/) for quick notifications on my phone if there are any errors.

<img src="./imgs/example_message.jpeg" width="360" height="640" alt="Example Message"/>

## Requirements

1. Go v1.19+
2. Working Kubernetes Cluster w/ the [SealedSecrets CRD](https://github.com/bitnami-labs/sealed-secrets) installed
3. Twilio Subscription and API Key
4. (optional) [ntfy.sh](https://ntfy.sh) for error reporting
5. (optional) a prometheus pushgateway instance

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
TWILIO_AUTH=api-key CAT_API_KEY=api-key PHONE_NUMBERS=+16666666666,+16666666666 ERR_NOTIFICATION_TOPIC=your_err_notification_topic go run .
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
