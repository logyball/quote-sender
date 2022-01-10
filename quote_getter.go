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

func GetQuoteFromApi(dogFriday bool) *QuoteObject {
	log.Info("Grabbing quotes from api")
	resp, err := http.Get(quoteApiUrl)
	if err != nil {
		defer resp.Body.Close()
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if !(resp.StatusCode < 300) {
		log.Fatalf("Status Code was %v", resp.StatusCode)
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	log.Infof("Quotes API returned: %v", string(body))
	if readErr != nil {
		log.Fatal(readErr)
	}
	quote := parseQuoteJsonResponse(body)
	if dogFriday {
		quote.Quote = fmt.Sprintf("It's Dog Friday!\n\n%s", quote.Quote)
	}
	return quote
}
