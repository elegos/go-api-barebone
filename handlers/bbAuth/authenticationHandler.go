package bbAuth

import (
	"api-barebone/src/bbAuth"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// AuthHandler the authentication handler to use
var AuthHandler bbAuth.AuthenticationHandler

type authError struct {
	Code         int    `json:"code"`
	ErrorMessage string `json:"errorMessage"`
}

type authPassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const errInternalServerErrorCode = 500
const errInternalServerErrorMessage = "Internal server error."

const errUnhandledBodyCode = 3000
const errUnhandledBodyMessage = "Unable to parse the request's body"

const errInvalidJSONPayloadCode = 3001
const errInvalidJSONPayloadMessage = "Invalid request JSON payload"

const errInvalidCredentialsCode = 3002
const errInvalidCredentialsMessage = "Invalid credentials"

func writeJSON(w http.ResponseWriter, data interface{}) {
	bytes, _ := json.Marshal(data)
	w.Write(bytes)
}

// AuthorizationGrantTypePassword handles the authentication via grant type password
func AuthorizationGrantTypePassword(w http.ResponseWriter, r *http.Request) {
	if AuthHandler == nil {
		bytes, _ := json.Marshal(authError{
			Code:         errInternalServerErrorCode,
			ErrorMessage: errInternalServerErrorMessage,
		})
		w.Write(bytes)

		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		writeJSON(w, authError{
			Code:         errUnhandledBodyCode,
			ErrorMessage: errUnhandledBodyMessage,
		})

		return
	}

	authCoords := authPassword{}
	err = json.Unmarshal(body, &authCoords)

	if err != nil {
		writeJSON(w, authError{
			Code:         errInvalidJSONPayloadCode,
			ErrorMessage: errInvalidJSONPayloadMessage,
		})

		return
	}

	token, err := AuthHandler.UsernamePasswordLogin(authCoords.Username, authCoords.Password)

	if err != nil {
		writeJSON(w, authError{
			Code:         errInvalidCredentialsCode,
			ErrorMessage: errInvalidCredentialsMessage,
		})

		return
	}

	token.Token = base64.StdEncoding.EncodeToString([]byte(token.Token))

	writeJSON(w, token)
}

// Logout the logout handler
func Logout(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Authentication")

	if strings.Index(header, "Bearer ") == -1 {
		return // no token provided, silent success (no action)
	}

	token := string([]rune(header)[7:])
	tokenBytes, err := base64.StdEncoding.DecodeString(token)
	token = string(tokenBytes)

	if err != nil {
		logrus.WithField("Authentication", header).Error("Error decoding bearer token")

		return
	}

	_, err = AuthHandler.Logout(token)

	if err != nil {
		address := r.RemoteAddr
		if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
			address = forwardedFor
		}

		logrus.WithFields(logrus.Fields{
			"address": address,
			"token":   token,
		}).Warn("User logout attempt failed")
	}
}
