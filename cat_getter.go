package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	catApiUrl     string = "https://api.thecatapi.com/v1/images/search?api_key=61c67453-a15e-4a0e-8254-ade03fb0ec05&mime_types=png"
	catUrlRetries int    = 5
)

var (
	bannedCatList = [...]string{"MzbkKPaBt"}
)

type CatObject struct {
	Url string `json: url`
}

// downloadFile rips a file straight from the internet onto the local filesystem
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
	if err != nil {
		log.Fatal(err)
	}
	err = os.Remove("./tmpCat")
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Cat %v was %v bytes", url, fileInfo.Size())
	return fileInfo.Size() <= 5000000
}

// isCatImageBanned checks a blacklist for known bad cat images
func isCatImageBanned(url string) bool {
	for _, bannedImage := range bannedCatList {
		if strings.Contains(url, bannedImage) {
			return true
		}
	}
	return false
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
// It must be <5mb, and will retry 5 times to get a cat
func GetCatFromApi() string {
	for i := 0; i < catUrlRetries; i++ {
		log.Info("getting cat from api")
		resp, err := http.Get(catApiUrl)
		if err != nil {
			defer resp.Body.Close()
			log.Fatal(err)
		}
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
		url := parseCatJsonResponse(respBody)
		if isCatImageBanned(url) {
			log.Infof("cat %v was on the blacklist, trying again", url)
			continue
		}
		if isCatImageSmallEnough(url) {
			log.Infof("cat %v was under the file limit, returning", url)
			return url
		}
		log.Infof("cat %v was too large :(, trying again", url)
	}
	log.Fatalf("Couldn't find a suitable cat in 5 tries :(, quitting with error")
	return ""
}
