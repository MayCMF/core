package jwtauth

import (
	"encoding/json"
)

// tokenInfo - Token information
type tokenInfo struct {
	AccessToken string `json:"access_token"` // Access token
	TokenType   string `json:"token_type"`   // Token type
	ExpiresAt   int64  `json:"expires_at"`   // Token expiration time
}

func (t *tokenInfo) GetAccessToken() string {
	return t.AccessToken
}

func (t *tokenInfo) GetTokenType() string {
	return t.TokenType
}

func (t *tokenInfo) GetExpiresAt() int64 {
	return t.ExpiresAt
}

func (t *tokenInfo) EncodeToJSON() ([]byte, error) {
	return json.Marshal(t)
}
