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
			expected: "9caf06bb4436cdbfa20af9121a626bc1093c4f54b31c0fa937957856135345b6",
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
