package links

import "errors"

var (
	ErrNotFound   = errors.New("link not found")
	ErrUniqueCode = errors.New("link code already exists")

	ErrCodeGenerationExhausted = errors.New("unique code generation exhausted")
)
