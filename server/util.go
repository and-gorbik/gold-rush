package server

import (
	"bytes"
	"gold-rush/infrastructure"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

func doRequest(client *http.Client, method, url string, input interface{}) ([]byte, error) {
	reqBody, err := createRequestBody(input)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, exploreURL, reqBody)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := readBody(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode >= http.StatusInternalServerError {
			return nil, err
		}

		log.Fatal(newBusinessError(body))
	}

	return body, nil
}

func newBusinessError(body []byte) error {
	var businessError infrastructure.BusinessError
	if err := jsoniter.Unmarshal(body, &businessError); err != nil {
		return err
	}

	return businessError
}

func readBody(r io.ReadCloser) ([]byte, error) {
	defer r.Close()
	return ioutil.ReadAll(r)
}

func createRequestBody(model interface{}) (io.Reader, error) {
	if model == nil {
		return nil, nil
	}

	data, err := jsoniter.Marshal(model)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}
