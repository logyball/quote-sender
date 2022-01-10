package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCatDecode(t *testing.T) {
	responseJsonStr := "[{\"breeds\":[],\"id\":\"6XxCfLGUC\",\"url\":\"https://cdn2.thecatapi.com/images/6XxCfLGUC.png\",\"width\":418,\"height\":700}]"
	resByteArr := []byte(responseJsonStr)

	url := parseCatJsonResponse(resByteArr)

	assert.Equal(t, "https://cdn2.thecatapi.com/images/6XxCfLGUC.png", url)
}