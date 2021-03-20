package infrastructure

type BusinessError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err BusinessError) Error() string {
	return err.Message
}

func ReadError(err error) (string, bool) {
	if e, ok := err.(*BusinessError); ok {
		return e.Message, true
	}

	return err.Error(), false
}
