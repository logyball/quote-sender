package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type QuoteObject struct {
	Quote string `json: quote`
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

func GetQuoteFromApi() *QuoteObject {
	log.Info("Grabbing quotes from api")
	resp, err := http.Get(quoteApiUrl)
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
	return parseQuoteJsonResponse(body)
}
