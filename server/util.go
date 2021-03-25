package server

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	jsoniter "github.com/json-iterator/go"

	"gold-rush/infrastructure"
)

func doRequest(client *http.Client, method, path string, input interface{}) ([]byte, error) {
	reqBody, err := createRequestBody(input)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, getURL(path), reqBody)
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

func getURL(path string) string {
	return fmt.Sprintf("%s://%s%s:%s", schema, host, path, port)
}

func envOrDefault(name, def string) string {
	env := os.Getenv(name)
	if env == "" {
		env = def
	}

	return env
}
