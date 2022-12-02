package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildTwilioMessage(t *testing.T) {
	_ = os.Setenv("TWILIO_AUTH", "asdf")
	testQuoteObj := QuoteObject{
		Quote:  "quote",
		Author: "author",
	}
	animalUrl := "www.google.com"
	dogFriday := false
	numberTo := "12345"
	req, err := buildTwilioMessage(&testQuoteObj, animalUrl, dogFriday, numberTo)

	assert.Equal(t, "https://api.twilio.com/2010-04-01/Accounts/AC785587cdbdd787fd35de9c2440f6ec26/Messages.json", req.URL.String())
	assert.Contains(t, req.Header.Get("Accept"), "application/json")
	assert.Contains(t, req.Header.Get("Content-Type"), "application/x-www-form-urlencoded")
	assert.Equal(t, "POST", req.Method)
	assert.NotEmpty(t, req.Body)
	assert.Nil(t, err)
	_ = os.Unsetenv("TWILIO_AUTH")
}

func TestPhoneNumberValidator(t *testing.T) {
	goodNumber := "+16666666666"
	badNumbers := []string{"+1666666666", "+26666666666", "+1a666666666", "16666666666", "+166666666660"}
	assert.Equal(t, true, phoneNumberValidator(goodNumber))
	for _, num := range badNumbers {
		assert.Equal(t, false, phoneNumberValidator(num))
	}
}

func TestGetPhoneNumbersFromEnv(t *testing.T) {
	_ = os.Setenv("PHONE_NUMBERS", "+16666666666,+166666666667")
	numbersToText, err := getNumbersToText()
	assert.Empty(t, err)
	assert.NotEmpty(t, numbersToText)
	assert.Contains(t, numbersToText, "+16666666666")
	assert.NotContains(t, numbersToText, "+166666666667")
	_ = os.Unsetenv("PHONE_NUMBERS")
}

func TestGetTwilioAuthErrorWhenEnvNotPresent(t *testing.T) {
	_ = os.Unsetenv("TWILIO_AUTH")
	_, err := getTwilioAuth()
	assert.Error(t, err)
}

func TestGetTwilioAuthFromEnv(t *testing.T) {
	os.Clearenv()
	err := os.Setenv("TWILIO_AUTH", "Not_Null")
	if err != nil {
		t.Fatal(err)
	}

	auth, err := getTwilioAuth()

	assert.Nil(t, err)
	assert.Equal(t, "Not_Null", auth)
	_ = os.Unsetenv("TWILIO_AUTH")
}
