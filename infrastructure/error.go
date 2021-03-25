package infrastructure

type BusinessError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err BusinessError) Error() string {
	return err.Message
}
