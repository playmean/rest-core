package restcore

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ParseValidateBody(body []byte, v interface{}) error {
	validate := validator.New()

	err := json.Unmarshal(body, v)
	if err != nil {
		return NewApiError(&ApiErrorOptions{
			Code:     "PARSE",
			Subcode:  "json",
			Message:  "error while parsing json body",
			Original: err,
		})
	}

	err = validate.Struct(v)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return NewApiError(&ApiErrorOptions{
				Code:     "VALIDATE",
				Message:  "error while validating body",
				Original: err,
			})
		}

		for _, err := range err.(validator.ValidationErrors) {
			return NewApiError(&ApiErrorOptions{
				Code:     "VALIDATE",
				Subcode:  err.Field(),
				Message:  fmt.Sprintf("field '%s' is %s", err.Field(), err.Tag()),
				Original: err,
			})
		}
	}

	return nil
}
