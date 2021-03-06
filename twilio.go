package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

const TwilioNumberFrom string = "+15033005651"
const TwilioSid string = "AC785587cdbdd787fd35de9c2440f6ec26"
const TwilioUrl string = "https://api.twilio.com/2010-04-01/Accounts/AC785587cdbdd787fd35de9c2440f6ec26/Messages.json"

func phoneNumberValidator(phoneNumber string) bool {
	valid, err := regexp.MatchString(`^\+1[0-9]{10}$`, phoneNumber)
	if err != nil {
		return false
	}
	return valid
}

// getNumbersToText expects a comma delimited list of phone numbers
// in an environment variable, e.g. NUMBERS=+16666666666,+16666666666
func getNumbersToText() ([]string, error) {
	var numbersToText []string
	var validNumbersToText []string

	numbersString := os.Getenv("PHONE_NUMBERS")
	numbersToText = strings.Split(numbersString, ",")
	log.Infof("numbers to text: %v", numbersToText)
	for _, num := range numbersToText {
		if phoneNumberValidator(num) {
			validNumbersToText = append(validNumbersToText, num)
		}
	}
	log.Infof("valid numbers to text: %v", validNumbersToText)
	if len(validNumbersToText) < 1 {
		log.Error("couldn't find any numbers to text")
		return nil, errors.New("couldn't find any numbers to text")
	}
	return validNumbersToText, nil
}

func getTwilioAuth() (string, error) {
	authKey := os.Getenv("TWILIO_AUTH")
	if authKey != "" {
		return authKey, nil
	}
	return "", errors.New("no twilio auth found in environment vars")
}

func buildTextString(quotes *QuoteObject, dogFriday bool) string {
	if dogFriday {
		return fmt.Sprintf("🐕 It's Dog Friday! 🐕\n\n\"%v\"\n\n-%v", quotes.Quote, quotes.Author)
	}
	return fmt.Sprintf("\"%v\"\n\n-%v", quotes.Quote, quotes.Author)
}

func buildTwilioMsgData(quote *QuoteObject, animalUrl string, dogFriday bool, numberTo string) *strings.Reader {
	msgString := buildTextString(quote, dogFriday)
	msgData := url.Values{}
	msgData.Set("To", numberTo)
	msgData.Set("From", TwilioNumberFrom)
	msgData.Set("Body", msgString)
	msgData.Set("MediaUrl", animalUrl)
	return strings.NewReader(msgData.Encode())
}

func buildTwilioMessage(quote *QuoteObject, animalUrl string, dogFriday bool, numberTo string) http.Request {
	msgDataReader := buildTwilioMsgData(quote, animalUrl, dogFriday, numberTo)
	req, err := http.NewRequest("POST", TwilioUrl, msgDataReader)
	if err != nil {
		log.Fatal(err)
	}
	auth, err := getTwilioAuth()
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(TwilioSid, auth)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return *req
}

func SendMessage(quote *QuoteObject, animalUrl string, dogFriday bool, numberTo string) error {
	msgReq := buildTwilioMessage(quote, animalUrl, dogFriday, numberTo)
	client := &http.Client{}

	log.Printf("Sending Request to %s via Twilio API: %v\n", numberTo, msgReq)
	resp, err := client.Do(&msgReq)
	if err != nil {
		log.Fatal(err)
		return err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Printf("Twilio API Response: %v", string(respBody))
	return nil
}
