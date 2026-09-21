package base62

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCode(t *testing.T) {
	code, err := GenerateCode(8)

	require.NoError(t, err)
	assert.Len(t, code, 8)

	for _, char := range code {
		assert.True(t, strings.ContainsRune(alphabet, char))
	}
}

func TestGenerateCodeInvalidLength(t *testing.T) {
	assert.Panics(t, func() {
		_, _ = GenerateCode(0)
	})
}
