package main

import (
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type FunFact struct {
	Fact string `json:"fact"`
}

const factApiUrl string = "https://api.api-ninjas.com/v1/facts?limit=1"
const factApiRetries int = 5

func BuildFactTwilioMessage(fact *FunFact) (string, error) {
	log.Debug("extracting fun fact from response")

	if fact.Fact == "" {
		return "", errors.New("no fun fact found in decoded object")
	}

	return fmt.Sprintf("📣 It's Fun Fact Saturday! 📣\n\nToday's fun fact: %s", fact.Fact), nil
}

func parseFactJsonResponse(responseBody []byte) (*FunFact, error) {
	log.Debug("parsing fun fact res into json")
	var val []FunFact

	err := json.Unmarshal(responseBody, &val)
	if err != nil {
		log.WithError(err).Error("error parsing fun fact into json")
		return nil, err
	}

	if len(val) > 0 {
		return &val[0], nil
	}

	return nil, errors.New("fun fact api response did not contain anything fun")
}

func getFactFromApi() (*FunFact, error) {
	log.Debug("Grabbing fun facts from api")

	body, err := GetApiNinja(factApiUrl, factApiRetries)
	if err != nil {
		return nil, err
	}

	fact, err := parseFactJsonResponse(body)
	if err != nil {
		return nil, err
	}

	return fact, nil
}

func GetFact() (*FunFact, error) {
	// creating wrapper to prepare for hitting DB if the API is down
	// redundancy, yo

	return getFactFromApi()
}
