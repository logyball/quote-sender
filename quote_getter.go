package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type QuoteObject struct {
	Quote  string `json: quote`
	Author string `json: author`
}

type QuoteContents struct {
	Quotes []QuoteObject `json: quotes`
}

type QuoteResponse struct {
	Contents QuoteContents `json: contents`
}

const quoteApiUrl string = "http://quotes.rest/qod?category=inspire&language=en"

func parseQuoteJsonResponse(responseBody []byte) *QuoteObject {
	var val QuoteResponse
	jsonErr := json.Unmarshal(responseBody, &val)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return &val.Contents.Quotes[0]
}

func GetQuoteFromApi() (*QuoteObject, error) {
	log.Info("Grabbing quotes from api")

	resp, err := http.Get(quoteApiUrl)
	defer resp.Body.Close()
	if err != nil {
		log.WithError(err).Error("failed to HTTP get quote")
		return nil, err
	}
	if !(resp.StatusCode < 300) {
		log.Errorf("HTTP Status code check failed: %d", resp.StatusCode)
		return nil, fmt.Errorf("HTTP Status code check failed: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	log.Infof("Quotes API returned: %v", string(body))
	if err != nil {
		log.WithError(err).Error("Error reading quote api return into struct")
		return nil, err
	}

	quote := parseQuoteJsonResponse(body)
	return quote, nil
}
