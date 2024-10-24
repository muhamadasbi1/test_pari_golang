package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func DecodeAndValidate(r *http.Request, req interface{}, validate *validator.Validate) error {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	return validate.Struct(req)
}
