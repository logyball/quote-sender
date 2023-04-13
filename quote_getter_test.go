package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildQuoteMessage(t *testing.T) {
	testQuoteObj := QuoteObject{
		Quote:  "quote",
		Author: "author",
	}

	tests := map[string]struct {
		input     QuoteObject
		dogFriday bool
		want      string
	}{
		"simple":     {input: testQuoteObj, dogFriday: false, want: "\"quote\"\n\n-author"},
		"dog friday": {input: testQuoteObj, dogFriday: true, want: "🐕 It's Dog Friday! 🐕\n\n\"quote\"\n\n-author"},
	}

	for name, tst := range tests {
		t.Run(name, func(t *testing.T) {
			res := BuildTwilioMessage(&tst.input, tst.dogFriday)
			assert.Equal(t, tst.want, res, "not right message")
		})
	}
}

func TestQuoteDecode(t *testing.T) {
	resByteArr := []byte("{\n  \"success\": {\n    \"total\": 1\n  },\n  \"contents\": {\n    \"quotes\": [\n      {\n        \"quote\": \"Many a false step was made by standing still.\",\n        \"length\": \"45\",\n        \"author\": \"Fortune Cookie\",\n        \"tags\": [\n          \"inspire\",\n          \"standing-still\"\n        ],\n        \"category\": \"inspire\",\n        \"language\": \"en\",\n        \"date\": \"2021-04-14\",\n        \"permalink\": \"https://theysaidso.com/quote/fortune-cookie-many-a-false-step-was-made-by-standing-still\",\n        \"id\": \"N0OhL98JryRfjljJqVtwGgeF\",\n        \"background\": \"https://theysaidso.com/img/qod/qod-inspire.jpg\",\n        \"title\": \"Inspiring Quote of the day\"\n      }\n    ]\n  },\n  \"baseurl\": \"https://theysaidso.com\",\n  \"copyright\": {\n    \"year\": 2023,\n    \"url\": \"https://theysaidso.com\"\n  }\n}")
	res, err := parseQuoteJsonResponse(resByteArr)

	assert.Nil(t, err)
	assert.Equal(t, "Many a false step was made by standing still.", res.Quote)
	assert.Equal(t, "Fortune Cookie", res.Author)
}

func TestBuildQuoteMsgStr(t *testing.T) {
	quote := QuoteObject{
		Quote:  "This is a Quote",
		Author: "Author",
	}
	assert.Equal(t, "\"This is a Quote\"\n\n-Author", buildTextString(&quote, false))
}
