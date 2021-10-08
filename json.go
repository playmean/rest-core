package restcore

import (
	"encoding/json"
)

func EncodeJson(input interface{}) ([]byte, error) {
	encoded, err := json.Marshal(input)
	if err != nil {
		return nil, NewApiError(&ApiErrorOptions{
			Code:     "ENCODE",
			Subcode:  "json",
			Message:  "error while encoding json",
			Original: err,
		})
	}

	return encoded, nil
}

func DecodeJson(input []byte, output interface{}) error {
	err := json.Unmarshal(input, output)
	if err != nil {
		return NewApiError(&ApiErrorOptions{
			Code:     "DECODE",
			Subcode:  "json",
			Message:  "error while decoding json",
			Original: err,
		})
	}

	return nil
}
