package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization") // api key in format "Authorization: Apikey <API_KEY>"
	if apiKey == "" {
		return "", errors.New("Authorization header is missing")
	}
	apiKey = strings.TrimSpace(apiKey)
	const apiKeyPrefix = "ApiKey "
	if len(apiKey) <= len(apiKeyPrefix) || apiKey[:len(apiKeyPrefix)] != apiKeyPrefix {
		return "", errors.New("Invalid Authorization header format")
	}

	return apiKey[len(apiKeyPrefix):], nil

}
