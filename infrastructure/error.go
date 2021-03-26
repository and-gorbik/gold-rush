package infrastructure

import "fmt"

type ProviderError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	StatusCode int
}

func (err ProviderError) Error() string {
	return fmt.Sprintf("[%d.%d] %s", err.StatusCode, err.Code, err.Message)
}
