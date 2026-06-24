package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "my-secret-key"
	// Test with a valid token
	token, err := MakeJWT(userID, tokenSecret, time.Minute*5)
	assert.NoError(t, err)

	validatedUserID, err := ValidateJWT(token, tokenSecret)
	assert.NoError(t, err)
	assert.Equal(t, userID, validatedUserID)

	// Test with an invalid token
	invalidToken := "invalid.token.string"
	_, err = ValidateJWT(invalidToken, tokenSecret)
	assert.Error(t, err)

	// Test with an expired token
	expiredToken, err := MakeJWT(userID, tokenSecret, -time.Minute*5) // Token expired 5 minutes ago
	assert.NoError(t, err)

	_, err = ValidateJWT(expiredToken, tokenSecret)
	assert.Error(t, err)
}

func TestValidateJWTWithInvalidIssuer(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "my-secret-key"
	// Create a token with an invalid issuer
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "invalid-issuer",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		Subject:   userID.String(),
	})
	tokenString, err := token.SignedString([]byte(tokenSecret))
	assert.NoError(t, err)

	_, err = ValidateJWT(tokenString, tokenSecret)
	assert.Error(t, err)
}

func TestWrongSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "my-secret-key"
	// Create a valid token
	token, err := MakeJWT(userID, tokenSecret, time.Minute*5)
	assert.NoError(t, err)

	// Validate the token with a wrong secret
	wrongSecret := "wrong-secret-key"
	_, err = ValidateJWT(token, wrongSecret)
	assert.Error(t, err)
}
