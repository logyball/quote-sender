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
	log.Debug("reading in phone numbers to text from environment")

	var validNumbersToText []string

	numbersString := os.Getenv("PHONE_NUMBERS")
	numbersToText := strings.Split(numbersString, ",")
	for _, num := range numbersToText {
		if phoneNumberValidator(num) {
			validNumbersToText = append(validNumbersToText, num)
		}
	}
	if len(validNumbersToText) < 1 {
		log.Error("couldn't find any numbers to text")
		return nil, errors.New("couldn't find any numbers to text")
	}
	log.Infof("valid numbers to text: %+v", validNumbersToText)

	return validNumbersToText, nil
}

func getTwilioAuth() (string, error) {
	log.Debug("getting twilio auth")

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

func buildTwilioMsgData(msgString string, animalUrl string, dogFriday bool, numberTo string) *strings.Reader {
	msgData := url.Values{}
	msgData.Set("To", numberTo)
	msgData.Set("From", TwilioNumberFrom)
	msgData.Set("Body", msgString)
	msgData.Set("MediaUrl", animalUrl)
	return strings.NewReader(msgData.Encode())
}

func buildTwilioPayload(msgString string, animalUrl string, dogFriday bool, numberTo string) (*http.Request, error) {
	log.WithFields(log.Fields{
		"msg":       msgString,
		"url":       animalUrl,
		"dogFriday": dogFriday,
		"numberTo":  numberTo,
	}).Debug("building twilio message")

	msgDataReader := buildTwilioMsgData(msgString, animalUrl, dogFriday, numberTo)
	req, err := http.NewRequest(http.MethodPost, TwilioUrl, msgDataReader)
	if err != nil {
		log.WithError(err).Error("error building post req for twilio")
		return nil, err
	}

	auth, err := getTwilioAuth()
	if err != nil {
		log.WithError(err).Error("error building post req for twilio")
		return nil, err
	}

	req.SetBasicAuth(TwilioSid, auth)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func SendMessage(msgString string, animalUrl string, dogFriday bool, numberTo string) error {
	log.Debug("Sending message to twilio")

	msgReq, err := buildTwilioPayload(msgString, animalUrl, dogFriday, numberTo)
	if err != nil {
		return err
	}

	client := &http.Client{}

	log.Infof("Sending Request to %s via Twilio API: %v\n", numberTo, msgReq)
	resp, err := client.Do(msgReq)
	if err != nil {
		log.WithError(err).Error("error sending post request")
		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("error reading response")
		return err
	}

	log.Infof("Twilio API Response: %v", string(respBody))
	return nil
}
