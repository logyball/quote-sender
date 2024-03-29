package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type QuoteObject struct {
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

type QuoteContents struct {
	Quotes []QuoteObject `json:"quotes"`
}

type QuoteResponse struct {
	Contents QuoteContents `json:"contents"`
}

const quoteApiUrl string = "https://quotes.rest/qod?category=inspire&language=en"
const quoteApiRetries int = 5

func parseQuoteJsonResponse(responseBody []byte) (*QuoteObject, error) {
	log.Debug("parsing quote res into json")
	var val QuoteResponse

	err := json.Unmarshal(responseBody, &val)
	if err != nil {
		log.WithError(err).Error("error parsing quote into json")
		return nil, err
	}

	return &val.Contents.Quotes[0], nil
}

func getQuoteFromApi() (*QuoteObject, error) {
	log.Debug("Grabbing quotes from api")

	quoteApiKey := os.Getenv("QUOTE_API_KEY")
	if quoteApiKey == "" {
		log.Error("QUOTE_API_KEY not found in environment vars")
		return nil, errors.New("QUOTE_API_KEY not found in environment vars")
	}

	var bearer = "Bearer " + quoteApiKey
	var resp *http.Response
	var body []byte

	req, err := http.NewRequest("GET", quoteApiUrl, nil)
	if err != nil {
		log.WithError(err).Error("Could not make req object")
		return nil, err
	}
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}

	for i := 0; i < quoteApiRetries; i++ {
		resp, err = client.Do(req)
		if err != nil {
			log.WithError(err).Error("failed to HTTP get quote")
			continue
		}
		defer resp.Body.Close()

		if !(resp.StatusCode < 300) {
			log.Errorf("HTTP Status code check failed: %d", resp.StatusCode)
			continue
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.WithError(err).Error("Error reading quote api return into struct")
			continue
		}

		quoteObj, err := parseQuoteJsonResponse(body)
		if err != nil {
			continue
		}

		log.Infof("Quote object returned: %+v", quoteObj)
		return quoteObj, nil
	}

	return nil, fmt.Errorf("failed to find a suitable quote after %d tries", quoteApiRetries)
}

func GetQuote() (*QuoteObject, error) {
	// creating wrapper to prepare for hitting DB if the API is down
	// redundancy, yo

	return getQuoteFromApi()
}
