package httpvalidator

import (
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidateHTTPURL(fl validator.FieldLevel) bool {
	raw := strings.TrimSpace(fl.Field().String())

	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	return u.Host != ""
}
