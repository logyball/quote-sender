package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"

	log "github.com/sirupsen/logrus"
)

const pushGatewayUri string = "http://prometheus-pushgateway.monitoring:9091/"
const promJobName string = "quote-messenger"

var LastSuccess = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "quote_messenger_last_successful_run",
	Help: "The Last time the quote messanger job completed successfully.",
})

func main() {
	quote := GetQuoteFromApi()
	numbersToText, err := getNumbersToText()
	if err != nil {
		log.Fatal(err)
	}
	for _, phoneNumber := range numbersToText {
		err := SendQuoteMessage(quote, phoneNumber)
		if err != nil {
			log.Fatal(err)
		}
	}
	LastSuccess.SetToCurrentTime()
	err = push.New(pushGatewayUri, promJobName).Collector(LastSuccess).Push()
	if err != nil {
		log.Print("Sending success metric failed.  the irony")
		log.Fatal(err)
	}
}
