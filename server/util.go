package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	// jsoniter "github.com/json-iterator/go"

	"gold-rush/infrastructure"
)

func doRequest(client *http.Client, method, path string, input interface{}) ([]byte, error) {
	reqBody, err := createRequestBody(input)
	if err != nil {
		log.Fatalf("doRequest(%s)/NewRequest: %v\n", path, err)
	}

	req, err := http.NewRequest(method, getURL(path), reqBody)
	if err != nil {
		log.Fatalf("doRequest(%s)/NewRequest: %v\n", path, err)
	}

	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := readBody(resp.Body)
	if err != nil {
		log.Printf("doRequest(%s)/readBody: %v\n", path, err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode >= http.StatusInternalServerError {
			return nil, err
		}

		err = newBusinessError(body)
		log.Printf("doRequest(%s): %v\n", path, err)
		return nil, err
	}

	log.Println("BODY: ", string(body))

	return body, nil
}

func newBusinessError(body []byte) error {
	var businessError infrastructure.BusinessError
	if err := json.Unmarshal(body, &businessError); err != nil {
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
		return bytes.NewReader(nil), nil
	}

	data, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}

func getURL(path string) string {
	return fmt.Sprintf("%s://%s:%s%s", schema, host, port, path)
}

func envOrDefault(name, def string) string {
	env := os.Getenv(name)
	if env == "" {
		env = def
	}

	return env
}
