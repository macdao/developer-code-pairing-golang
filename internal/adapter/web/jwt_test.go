package web

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	userID := uint64(1001)

	token, err := GenerateToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken_Success(t *testing.T) {
	userID := uint64(1001)

	token, err := GenerateToken(userID)
	assert.NoError(t, err)

	claims, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	_, err := ValidateToken("invalid-token")
	assert.Error(t, err)
}

func TestExtractUserID_Success(t *testing.T) {
	userID := uint64(1001)

	token, err := GenerateToken(userID)
	assert.NoError(t, err)

	extractedUserID, err := ExtractUserID(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, extractedUserID)
}

func TestExtractUserID_InvalidToken(t *testing.T) {
	_, err := ExtractUserID("invalid-token")
	assert.Error(t, err)
}
