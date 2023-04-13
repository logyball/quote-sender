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

type TriviaObject struct {
	Category string `json:category`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

const triviaApiUrl string = "https://api.api-ninjas.com/v1/trivia"
const triviaApiRetries int = 5

func parseTriviaJsonResponse(responseBody []byte) (*TriviaObject, error) {
	log.Debug("parsing trivia res into json")
	var val []TriviaObject

	err := json.Unmarshal(responseBody, &val)
	if err != nil {
		log.WithError(err).Error("error parsing trivia into json")
		return nil, err
	}

	if len(val) > 0 {
		return &val[0], nil
	}

	return nil, errors.New("no trivia found")
}

func getTriviaFromApi() (*TriviaObject, error) {
	log.Debug("Grabbing trivia from api")

	apiNinjaKey := os.Getenv("API_NINJA_KEY")
	if apiNinjaKey == "" {
		log.Error("API_NINJA_KEY not found in environment vars")
		return nil, errors.New("API_NINJA_KEY not found in environment vars")
	}

	var resp *http.Response
	var body []byte

	req, err := http.NewRequest("GET", triviaApiUrl, nil)
	if err != nil {
		log.WithError(err).Error("Could not make req object")
		return nil, err
	}
	req.Header.Add("X-Api-Key", apiNinjaKey)
	client := &http.Client{}

	for i := 0; i < triviaApiRetries; i++ {
		resp, err = client.Do(req)
		if err != nil {
			log.WithError(err).Error("failed to HTTP get trivia")
			continue
		}
		defer resp.Body.Close()

		if !(resp.StatusCode < 300) {
			log.Errorf("HTTP Status code check failed: %d", resp.StatusCode)
			continue
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.WithError(err).Error("Error reading trivia api return into struct")
			continue
		}

		triviaObj, err := parseTriviaJsonResponse(body)
		if err != nil {
			continue
		}

		log.Infof("Trivia object returned: %+v", triviaObj)
		return triviaObj, nil
	}

	return nil, fmt.Errorf("failed to find a suitable trivia object after %d tries", triviaApiRetries)
}

func MakeTriviaTwilioMessage(trivia *TriviaObject) string {
	category := trivia.Category
	if category == "" {
		category = "Free for all"
	}
	return fmt.Sprintf("🏆❓ It's Trivia Tuesday! ❓🏆\n\nToday's category is: %s!\n\nQuestion: %s?\n\n...\n...\n...\n...\n\nAnswer: %s", category, trivia.Question, trivia.Answer)
}

func GetTrivia() (*TriviaObject, error) {
	// creating wrapper to prepare for hitting DB if the API is down
	// redundancy, yo

	return getTriviaFromApi()
}
