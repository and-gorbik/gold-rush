package server

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"gold-rush/models"
)

func (s *GoldRushServer) doRequest(method, url string, input interface{}) ([]byte, error) {
	reqBody, err := createRequestBody(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, exploreURL, reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := readBody(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newBusinessError(body)
	}

	return body, nil
}

func newBusinessError(body []byte) error {
	var businessError models.BusinessError
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
