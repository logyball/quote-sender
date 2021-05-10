package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

const catApiUrl string = "https://api.thecatapi.com/v1/images/search?api_key=61c67453-a15e-4a0e-8254-ade03fb0ec05&mime_types=png"

type CatObject struct {
	Url string `json: url`
}


func parseCatJsonResponse(responseBody []byte) string {
	var val []CatObject
	err := json.Unmarshal(responseBody, &val)
	if err != nil {
		log.Fatal(err)
	}
	return val[0].Url
}


// GetCatFromApi returns a URL with a random cat pic
func GetCatFromApi() string {
	log.Info("getting cat from api")
	resp, err := http.Get(catApiUrl)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if !(resp.StatusCode < 300) {
		log.Fatalf("Status code was %v", resp.StatusCode)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return parseCatJsonResponse(respBody)
}