package app

import (
	"gold-rush/models"
)

func retry(attemps int) {

}

func readError(err error) (string, bool) {
	if e, ok := err.(*models.BusinessError); ok {
		return e.Message, true
	}

	return err.Error(), false
}
