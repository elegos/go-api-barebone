package bbAuth_test

import (
	auth "api-barebone/src/bbAuth"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSecrets(t *testing.T) {
	handler := &auth.PlainTextAuthHandler{
		AuthFile:    "./test_auth.csv",
		SecretsFile: "./test_secrets.json",
	}

	_, err := handler.GetSecrets()
	assert.Nil(t, err, "The handler should be able to retrieve the secrets")

	secrets, err := handler.GetSecrets()
	assert.Nil(t, err, "The handler should be able to retrieve the secrets")
	assert.NotNil(t, secrets, "Secrets should be returned by the function")

	assert.Len(t, secrets, 2)
	assert.Equal(t, "p4ssw0rd", secrets["us3rn4m3"].Password, "us3rn4m3's password should be p4ssw0rd")
	assert.Equal(t, float64(0), secrets["us3rn4m3"].Token.ExpireEpochTime, "The plain password token should never expire")
	assert.NotZero(t, secrets["us3rn4m3"].Token.RefreshToken)
	assert.NotZero(t, secrets["us3rn4m3"].Token.Token)
	assert.NotEqual(t, secrets["us3rn4m3"].Token.Token, secrets["us3rn4m3"].Token.RefreshToken)
}

func TestUsernamePasswordLogin(t *testing.T) {
	handler := auth.PlainTextAuthHandler{
		AuthFile:    "./test_auth.csv",
		SecretsFile: "./test_secrets.json",
	}

	token, err := handler.UsernamePasswordLogin("us3rn4m3", "p4ssw0rd")
	assert.Nil(t, err, "The username should be correctly identified")
	assert.NotZero(t, token.Token, "The username's token should be correctly returned")

	token, err = handler.UsernamePasswordLogin("wrong", "password")
	assert.NotNil(t, err, "The username should not be identified")
	assert.Zero(t, token.Token, "A non-existing user should have no token")
}

func TestRefreshToken(t *testing.T) {
	handler := auth.PlainTextAuthHandler{
		AuthFile:    "./test_auth.csv",
		SecretsFile: "./test_secrets.json",
	}

	secrets, err := handler.GetSecrets()
	assert.NoError(t, err)

	token := secrets["us3rn4m3"].Token.Token
	refreshToken := secrets["us3rn4m3"].Token.RefreshToken

	newToken, err := handler.RefreshToken(refreshToken)
	assert.NoError(t, err)
	assert.Equal(t, token, newToken, "The tokens should always be the same (no expiration)")

	newToken, err = handler.RefreshToken("invalid refresh token")
	assert.Error(t, err)
	assert.Equal(t, "", newToken, "It should return no token")
}

func TestLogout(t *testing.T) {
	handler := auth.PlainTextAuthHandler{
		AuthFile:    "./test_auth.csv",
		SecretsFile: "./test_secrets.json",
	}

	secrets, err := handler.GetSecrets()
	assert.NoError(t, err)

	integrity := secrets["us3rn4m3"].IntegrityCheck
	token := secrets["us3rn4m3"].Token.Token
	refreshToken := secrets["us3rn4m3"].Token.RefreshToken

	handler.Logout(token)
	secrets, _ = handler.GetSecrets()

	assert.Equal(t, integrity, secrets["us3rn4m3"].IntegrityCheck)
	assert.NotEqual(t, token, secrets["us3rn4m3"].Token.Token)
	assert.NotEqual(t, refreshToken, secrets["us3rn4m3"].Token.RefreshToken)
}
