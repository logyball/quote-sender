package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"golang.org/x/sync/errgroup"

	log "github.com/sirupsen/logrus"
)

const promJobName string = "quote-messenger"
const notificationTopicBase string = "https://ntfy.sh/%s"

var LastSuccess = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "quote_messenger_last_successful_run",
	Help: "The Last time the quote messanger job completed successfully.",
})

func errHandling(err error, explaination string) {
	errReportingEnabled := os.Getenv("ERROR_REPORT_ENABLED")
	if errReportingEnabled == "" || errReportingEnabled == "false" {
		return
	}

	topic := os.Getenv("ERR_NOTIFICATION_TOPIC")
	if topic == "" {
		log.Fatal("Failed to retrieve error notification topic from environment vars")
	}

	log.WithField("err", err).WithField("explaination", explaination).Info("Error handling")

	fullUrl := fmt.Sprintf(notificationTopicBase, topic)
	req, err := http.NewRequest(
		http.MethodPost,
		fullUrl,
		bytes.NewBuffer([]byte(fmt.Sprintf("COTD ERROR\nErr: %v\nExplaination: %s", err, explaination))),
	)
	if err != nil {
		log.WithError(err).Fatal("Error making HTTP req body, ironic")
	}

	resp, err := http.Post(fullUrl, "application/json", req.Body)
	if err != nil {
		log.WithError(err).Fatal("Error making HTTP POST request, ironic")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Fatalf("HTTP status code was non-successful, ironic.  Code: %d", resp.StatusCode)
	}
}

func setLogLevel() {
	log.SetReportCaller(true)
	env := os.Getenv("ENVIRONMENT")
	if env == "prod" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.TraceLevel)
	}
}

func pushToPrometheus() {
	promGatewayEnabled := os.Getenv("PUSH_TO_PROMETHEUS")
	if promGatewayEnabled == "" || promGatewayEnabled == "false" {
		return
	}

	pushGatewayUri := os.Getenv("PROMETHEUS_GATEWAY_URI")
	if pushGatewayUri == "" {
		errHandling(errors.New("prometheus push gateway URI not defined"), "Error sending metric to push gateway")
		log.Fatal("Sending success metric failed.  the irony")
	}

	LastSuccess.SetToCurrentTime()
	err := push.New(pushGatewayUri, promJobName).Collector(LastSuccess).Push()
	if err != nil {
		log.Error("Sending success metric failed.  the irony")
		errHandling(err, "Error sending metric to push gateway")
	}
}

func main() {
	var quote *QuoteObject
	var err error
	var animalUrl string
	var numbersToText []string

	setLogLevel()
	isItDogFridayBabeee := time.Now().Weekday().String() == "Friday"

	eg := new(errgroup.Group)

	eg.Go(func() error {
		quote, err = GetQuote()
		if err != nil {
			errHandling(err, "Error getting quotes from the quote of the day API")
		}
		return err
	})
	eg.Go(func() error {
		animalUrl, err = GetAnimalFromApi(isItDogFridayBabeee)
		if err != nil {
			errHandling(err, "Error getting animals from the cat/dog of the day API")
		}
		return err
	})
	eg.Go(func() error {
		numbersToText, err = getNumbersToText()
		if err != nil {
			errHandling(err, "Error getting which numbers to dial")
		}
		return err
	})

	err = eg.Wait()
	if err != nil {
		log.Fatal(err)
	}

	for _, phoneNumber := range numbersToText {
		p := phoneNumber

		eg.Go(func() error {
			err := SendMessage(quote, animalUrl, isItDogFridayBabeee, p)
			if err != nil {
				errHandling(err, "Error sending messages w/ twilio")
			}
			return err
		})
	}

	err = eg.Wait()
	if err != nil {
		log.Fatal(err)
	}

	pushToPrometheus()
}
