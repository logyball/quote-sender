# Quote Messager

This is a dumb program that uses golang to fetch the quote of the day from `http://quote.rest/` and then uses twilio to send it to my phone.  The real value of it is that it's a simple thing that can be served up via k8s in a cronjob.

### Building

`go build`

### Testing

`go test -v ./...`

### K8s

The real value of this is deploying it to kubernetes as a cronjob and then looking to add prometheus push gateway metrics.

### Secrets

The phone numbers to text as well as the twilio api key are obscured as secrets.  Set them as environment variables for local testing:

```shell
TWILIO_AUTH=api-key PHONE_NUMBERS=+16666666666,+16666666666 go run
```