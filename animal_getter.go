package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	catApiUrlBase string = "https://api.thecatapi.com/v1/images/search?api_key=%s"
	dogApiUrl     string = "https://dog.ceo/api/breeds/image/random"
	urlRetries    int    = 5
	maxFileSize   int    = 5000000 // 5 MB
)

var (
	bannedCatList    = [...]string{"MzbkKPaBt"}
	allowedFileTypes = map[string]bool{
		"jpeg": true,
		"jpg":  true,
		"gif":  true,
		"png":  true,
	}
)

type CatObject struct {
	Url string `json:"url"`
}

type DogObject struct {
	Message string
	Status  string
}

// downloadFile rips a file straight from the internet onto the local filesystem
func downloadFile(filepath string, url string) error {
	log.WithFields(log.Fields{"filepath": filepath, "url": url}).Debug("downloading file")

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
func isImageSmallEnough(url string) (bool, error) {
	log.WithField("url", url).Debugf("checking if file is under the limit of %d bytes", maxFileSize)

	err := downloadFile("./tmpAnimal", url)
	if err != nil {
		log.WithError(err).Error("failed to download file")
		return false, err
	}

	fileInfo, err := os.Stat("./tmpAnimal")
	if err != nil {
		log.WithError(err).Error("failed to get file stats")
		return false, err
	}

	err = os.Remove("./tmpAnimal")
	if err != nil {
		log.WithError(err).Error("failed to remove file")
		return false, err
	}

	return fileInfo.Size() <= int64(maxFileSize), nil
}

// isCatImageBanned checks a blacklist for known bad cat images
func isCatImageBanned(url string) bool {
	log.WithField("url", url).Debug("checking URL for banned cat images")

	for _, bannedImage := range bannedCatList {
		if strings.Contains(url, bannedImage) {
			log.Infof("cat %s was on the blacklist", url)
			return true
		}
	}

	return false
}

func isFileTypeAllowed(url string) (bool, error) {
	log.WithField("url", url).Debug("checking URL contains a file on twilios allowed file type list")

	filepath := path.Base(url)
	filepathSplit := strings.Split(filepath, ".")
	if len(filepathSplit) < 1 {
		log.Errorf("Could not parse filepath string of %s into filetype", filepath)
		return false, fmt.Errorf("could not parse filepath string of %s into filetype", filepath)
	}

	fileExtension := filepathSplit[1]
	if _, ok := allowedFileTypes[fileExtension]; ok {
		return true, nil
	}

	log.Errorf("Filetype was not allowed for %s", url)
	return false, fmt.Errorf("filetype was not allowed for %s", url)
}

func parseCatJsonResponse(responseBody []byte) (string, error) {
	log.Debug("parsing response from the cat API")
	var val []CatObject

	err := json.Unmarshal(responseBody, &val)
	if err != nil {
		log.WithError(err).Error("failed to parse cat API response")
		return "", err
	}

	return val[0].Url, nil
}

func parseDogJsonResponse(responseBody []byte) (string, error) {
	log.Debug("parsing response from the dog API")
	var val DogObject

	err := json.Unmarshal(responseBody, &val)
	if err != nil {
		log.WithError(err).Error("failed to parse dog API response")
		return "", err
	}
	if val.Status != "success" {
		log.Error("dog API returned a non-successful status")
		return "", fmt.Errorf("dog API returned a non-successful status: %+v", val)
	}

	return val.Message, nil
}

func getCatApiResponse() (*http.Response, error) {
	log.Debug("getting info from the cat API")

	catApiKey := os.Getenv("CAT_API_KEY")
	if catApiKey == "" {
		log.Error("CAT_API_KEY not found in environment vars")
		return nil, errors.New("CAT_API_KEY not found in environment vars")
	}

	return http.Get(fmt.Sprintf(catApiUrlBase, catApiKey))
}

func getAnimalApiResponse(dogFriday bool) (*http.Response, error) {
	if dogFriday {
		log.Debug("getting info from the dog API")
		return http.Get(dogApiUrl)
	}
	return getCatApiResponse()
}

func parseAnimalResponse(dogFriday bool, responseBody []byte) (string, error) {
	if dogFriday {
		return parseDogJsonResponse(responseBody)
	}
	url, err := parseCatJsonResponse(responseBody)
	if isCatImageBanned(url) {
		return "", fmt.Errorf("cat %s is on the blacklist", url)
	}

	return url, err
}

// GetAnimalFromApi returns a URL with a random cat pic
// It must be <5mb, and will retry 5 times to get a cat
func GetAnimalFromApi(dogFriday bool) (string, error) {
	log.WithField("isDogFriday", dogFriday).Debug("Getting animal image from dog/cat API")

	var resp *http.Response
	var err error
	var url string

	for i := 0; i < urlRetries; i++ {
		resp, err = getAnimalApiResponse(dogFriday)
		if err != nil {
			log.WithError(err).Error("failed to HTTP get animal URL")
			return "", err
		}
		defer resp.Body.Close()

		if !(resp.StatusCode < 300) {
			log.Errorf("Status code was %v", resp.StatusCode)
			return "", fmt.Errorf("status code was %v", resp.StatusCode)
		}

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.WithError(err).Error("Error reading response body")
			return "", err
		}

		url, err = parseAnimalResponse(dogFriday, respBody)
		if err != nil {
			continue
		}

		allowed, err := isFileTypeAllowed(url)
		if err != nil || !allowed {
			continue
		}

		smallEnough, err := isImageSmallEnough(url)
		if err != nil {
			continue
		}

		if smallEnough {
			log.Infof("url fits all criteria, returning %s", url)
			return url, nil
		}
	}

	return "", fmt.Errorf("couldn't find a suitable animal in %d tries :(, quitting with error", urlRetries)
}
