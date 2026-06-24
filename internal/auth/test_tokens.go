package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRefreshToken(t *testing.T) {
	token1 := MakeRefreshToken()
	token2 := MakeRefreshToken()

	assert.NotEmpty(t, token1, "Refresh token should not be empty")
	assert.NotEmpty(t, token2, "Refresh token should not be empty")
	assert.NotEqual(t, token1, token2, "Two generated refresh tokens should not be equal")
}

func TestRefreshTokenrevocation(t *testing.T) {
	// Create a refresh token
	refreshToken := MakeRefreshToken()

	// Simulate storing the refresh token in a database (in-memory for this test)
	storedTokens := map[string]bool{
		refreshToken: true,
	}

	// Function to revoke the refresh token
	revokeRefreshToken := func(token string) {
		if _, exists := storedTokens[token]; exists {
			delete(storedTokens, token)
		}
	}

	// Revoke the refresh token
	revokeRefreshToken(refreshToken)

	// Check if the refresh token is revoked
	_, exists := storedTokens[refreshToken]
	assert.False(t, exists, "Refresh token should be revoked and not exist in the stored tokens")
}
