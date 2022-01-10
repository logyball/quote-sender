package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	catApiUrl  string = "https://api.thecatapi.com/v1/images/search?api_key=61c67453-a15e-4a0e-8254-ade03fb0ec05&mime_types=png"
	dogApiUrl  string = "https://dog.ceo/api/breeds/image/random"
	urlRetries int    = 5
)

var (
	bannedCatList    = [...]string{"MzbkKPaBt"}
	allowedFileTypes = [...]string{"jpeg", "jpg", "gif", "png"}
)

type CatObject struct {
	Url string `json: url`
}

type DogObject struct {
	Message string
	Status  string
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

// isImageSmallEnough makes sure that the returned cat URL is <5mb
func isImageSmallEnough(url string) bool {
	err := downloadFile("./tmpAnimal", url)
	if err != nil {
		log.Fatal(err)
	}
	fileInfo, err := os.Stat("./tmpAnimal")
	if err != nil {
		log.Fatal(err)
	}
	err = os.Remove("./tmpAnimal")
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Animal image %v was %v bytes", url, fileInfo.Size())
	return fileInfo.Size() <= 5000000
}

// isCatImageBanned checks a blacklist for known bad cat images
func isCatImageBanned(url string) bool {
	for _, bannedImage := range bannedCatList {
		if strings.Contains(url, bannedImage) {
			log.Infof("cat %v was on the blacklist, trying again", url)
			return true
		}
	}
	return false
}

func isFileTypeAllowed(url string) bool {
	filepath := path.Base(url)
	filepathSplit := strings.Split(filepath, ".")
	if len(filepathSplit) < 1 {
		log.Infof("Could not parse filepath string of %v into filetype", filepath)
		return false
	}
	fileExtension := filepathSplit[1]
	for _, filetype := range allowedFileTypes {
		if fileExtension == filetype {
			log.Infof("Filetype was allowed for %v", url)
			return true
		}
	}
	log.Infof("Filetype was not allowed for %v", url)
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

func parseDogJsonResponse(responseBody []byte) string {
	var val DogObject
	err := json.Unmarshal(responseBody, &val)
	if err != nil || val.Status != "success" {
		log.Fatal(err)
	}
	return val.Message
}

// GetAnimalFromApi returns a URL with a random cat pic
// It must be <5mb, and will retry 5 times to get a cat
func GetAnimalFromApi(dogFriday bool) string {
	for i := 0; i < urlRetries; i++ {
		log.Info("getting animal from api")
		var resp *http.Response
		var err error
		var url string

		if !dogFriday {
			resp, err = http.Get(catApiUrl)
		} else {
			resp, err = http.Get(dogApiUrl)
		}
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
		if !dogFriday {
			url = parseCatJsonResponse(respBody)
			if isCatImageBanned(url) {
				continue
			}
		} else {
			url = parseDogJsonResponse(respBody)
			if !isFileTypeAllowed(url) {
				continue
			}
		}
		if isImageSmallEnough(url) {
			log.Infof("animal %v was under the file limit, returning", url)
			return url
		}
		log.Infof("animal %v was too large :(, trying again", url)
	}
	log.Fatalf("Couldn't find a suitable animal in 5 tries :(, quitting with error")
	return ""
}
