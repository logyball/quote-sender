package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCatDecode(t *testing.T) {
	resBytes := []byte("[{\"breeds\":[],\"id\":\"6XxCfLGUC\",\"url\":\"https://cdn2.thecatapi.com/images/6XxCfLGUC.png\",\"width\":418,\"height\":700}]")

	url, err := parseCatJsonResponse(resBytes)

	assert.Equal(t, "https://cdn2.thecatapi.com/images/6XxCfLGUC.png", url)
	assert.Nil(t, err, "error was not nil that should have been")
}

func TestDogDecode(t *testing.T) {
	resBytes := []byte(`{
		"message":"https://images.dog.ceo/breeds/ovcharka-caucasian/IMG_20190602_204319.jpg",
		"status":"success"
	 }`)

	url, err := parseDogJsonResponse(resBytes)

	assert.Equal(t, "https://images.dog.ceo/breeds/ovcharka-caucasian/IMG_20190602_204319.jpg", url)
	assert.Nil(t, err, "error was not nil that should have been")

}
