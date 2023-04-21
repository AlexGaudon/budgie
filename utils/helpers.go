package utils

import (
	"encoding/json"
	"net/http"
)

func DecodeBody(r *http.Request, rv any) error {
	if err := json.NewDecoder(r.Body).Decode(&rv); err != nil {
		return err
	}

	return nil
}
