package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashRefreshToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:     "non-empty string",
			input:    "test-refresh-token",
			expected: "0a9b110d5e553bd98e9965c70a601c15c36805016ba60d54f20f5830c39edcde",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hashRefreshToken(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	plainRefreshToken, hashedRefreshToken, err := generateRefreshToken()
	assert.NoError(t, err)
	assert.NotEmpty(t, plainRefreshToken)
	assert.NotEmpty(t, hashedRefreshToken)
	assert.Equal(t, hashRefreshToken(plainRefreshToken), hashedRefreshToken)
}
