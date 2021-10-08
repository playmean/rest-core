package restcore

import (
	"log"
	"net/http"
)

type Sender interface {
	JSON(interface{}) error
	SendStatus(int) error
}

func SendError(s Sender, err error, status int) error {
	if apiError, ok := err.(*ApiError); ok {
		original := apiError.Original()

		if original != nil {
			log.Print(original.Error())
		}

		s.JSON(map[string]interface{}{
			"errorCode":    apiError.Code(),
			"errorMessage": apiError.Message(),
		})

		return s.SendStatus(status)
	}

	return s.SendStatus(http.StatusInternalServerError)
}
