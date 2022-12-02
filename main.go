package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"golang.org/x/sync/errgroup"

	log "github.com/sirupsen/logrus"
)

const pushGatewayUri string = "http://prometheus-pushgateway.monitoring:9091/"
const promJobName string = "quote-messenger"
const notificationTopicBase string = "https://ntfy.sh/%s"

var LastSuccess = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "quote_messenger_last_successful_run",
	Help: "The Last time the quote messanger job completed successfully.",
})

func errHandling(err error, explaination string) {
	topic := os.Getenv("ERR_NOTIFICATION_TOPIC")
	if topic == "" {
		log.Fatal("Error in the err handling, ironic")
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

func main() {
	var quote *QuoteObject
	var err error
	var animalUrl string
	var numbersToText []string

	isItDogFridayBabeee := time.Now().Weekday().String() == "Friday"
	eg := new(errgroup.Group)

	eg.Go(func() error {
		quote, err = GetQuoteFromApi()
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
		err := SendMessage(quote, animalUrl, isItDogFridayBabeee, phoneNumber)
		if err != nil {
			errHandling(err, "Error sending messages w/ twilio")
			log.Fatal(err)
		}
	}

	LastSuccess.SetToCurrentTime()
	err = push.New(pushGatewayUri, promJobName).Collector(LastSuccess).Push()
	if err != nil {
		log.Info("Sending success metric failed.  the irony")
		errHandling(err, "Error sending metric to push gateway")
		log.Fatal(err)
	}
}
