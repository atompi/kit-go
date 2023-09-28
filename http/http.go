package http

import (
	"bytes"
	"io"
	"net/http"
)

func HttpGet(url string) (data []byte, err error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	request.Header.Set("Accept", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	data, err = io.ReadAll(response.Body)
	return
}

func HttpPost(url string, postData []byte) (data []byte, err error) {
	body := bytes.NewBuffer(postData)
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	data, err = io.ReadAll(response.Body)
	return
}
