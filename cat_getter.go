package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const catApiUrl string = "https://api.thecatapi.com/v1/images/search?api_key=61c67453-a15e-4a0e-8254-ade03fb0ec05&mime_types=png"

type CatObject struct {
	Url string `json: url`
}

func downloadFile(filepath string, url string) error {
	log.Infof("downloading %v", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

// isCatImageSmallEnough makes sure that the returned cat URL is <5mb
func isCatImageSmallEnough(url string) bool {
	err := downloadFile("./tmpCat", url)
	if err != nil {
		log.Fatal(err)
	}
	fileInfo, err := os.Stat("./tmpCat")
	err = os.Remove("./tmpCat")
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Cat %v was %v bytes", url, fileInfo.Size())
	return fileInfo.Size() <= 5000000
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
	for {
		log.Info("getting cat from api")
		resp, err := http.Get(catApiUrl)
		defer resp.Body.Close()
		if err != nil{
			log.Fatal(err)
		}
			if !(resp.StatusCode < 300){
			log.Fatalf("Status code was %v", resp.StatusCode)
		}
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil{
			log.Fatal(err)
		}
		url := parseCatJsonResponse(respBody)
		if isCatImageSmallEnough(url) {
			log.Infof("cat %v was under the file limit, returning", url)
			return url
		}
		log.Infof("cat %v was too large :(, trying again", url)
	}
}