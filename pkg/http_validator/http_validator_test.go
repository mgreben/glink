package httpvalidator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type urlRequest struct {
	URL string `validate:"http_url"`
}

func TestValidateHTTPURL(t *testing.T) {
	t.Parallel()

	validate := validator.New()
	require.NoError(t, validate.RegisterValidation("http_url", ValidateHTTPURL))

	testCases := []struct {
		name  string
		url   string
		valid bool
	}{
		{name: "HTTPS URL", url: "https://example.com/path?key=value#fragment", valid: true},
		{name: "HTTP URL", url: "http://localhost:8080/health", valid: true},
		{name: "whitespace around URL", url: "  https://example.com  ", valid: true},
		{name: "empty URL", url: "", valid: false},
		{name: "host is missing", url: "https:///path", valid: false},
		{name: "scheme is missing", url: "example.com/path", valid: false},
		{name: "unsupported scheme", url: "ftp://example.com/file", valid: false},
		{name: "javascript URL", url: "javascript:alert(1)", valid: false},
		{name: "malformed URL", url: "https://[invalid", valid: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validate.Struct(urlRequest{URL: tc.url})

			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
