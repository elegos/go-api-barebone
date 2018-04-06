package bbAuth

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"strings"
)

// This authentication handler is just an example.
// A more advanced (database-driven?) implementation should be coded.

type secretsFile struct {
	Salt string `json:"salt"`
}

// PlainTextAuthHandler a plain-text authentication handler
type PlainTextAuthHandler struct {
	AuthFile    string // path to the auth file
	SecretsFile string // path to the secrets file

	secrets map[string]PlainPasswordToken // username / Token hash map
	seed    string                        // the tokens seed
}

// PlainPasswordToken the token containing the token itself and the password
type PlainPasswordToken struct {
	Password       string
	IntegrityCheck string
	Token          Token
}

func randomToken(salt string) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	runesLen := len(letterRunes)

	result := ""

	for i := int64(0); i < 512; i++ {
		result += string(letterRunes[rand.Int31n(int32(runesLen-1))])
	}

	md5Hasher := md5.New()
	md5Hasher.Write([]byte(result + ":" + salt))

	return hex.EncodeToString(md5Hasher.Sum(nil))
}

func containsString(slice []string, search string) bool {
	for _, element := range slice {
		if element == search {
			return true
		}
	}
	return false
}

func uniqueRandomString(existingStrings *[]string, salt string) string {
	str := ""
	nonUnique := true

	for nonUnique {
		str = randomToken(salt)
		nonUnique = containsString(*existingStrings, str)
	}

	*existingStrings = append(*existingStrings, str)

	return str
}

// GetSecrets get the username/password hash map
func (handler *PlainTextAuthHandler) GetSecrets() (map[string]PlainPasswordToken, error) {
	if handler.AuthFile == "" {
		return nil, errors.New("You need to configure the plain text authentication handler first! Missing AuthFile")
	}

	if handler.SecretsFile == "" {
		return nil, errors.New("You need to configure the plain text authentication handler first! Missing SecretsFile")
	}

	data, err := ioutil.ReadFile(handler.AuthFile)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(strings.NewReader(string(data)))

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadFile(handler.SecretsFile)
	if err != nil {
		return nil, err
	}

	secrets := secretsFile{}
	json.Unmarshal(data, &secrets)

	authTokens := []string{}
	refreshTokens := []string{}
	if handler.secrets == nil {
		handler.secrets = map[string]PlainPasswordToken{}
	} else {
		for _, secret := range handler.secrets {
			authTokens = append(authTokens, secret.Token.Token)
			refreshTokens = append(refreshTokens, secret.Token.RefreshToken)
		}
	}

	md5Hasher := md5.New()
	for _, record := range records {
		// md5(username,password,seed)
		md5Hasher.Write([]byte(secrets.Salt + record[0] + record[1]))
		integrity := hex.EncodeToString(md5Hasher.Sum(nil))

		if secret, ok := handler.secrets[record[0]]; ok && secret.IntegrityCheck == integrity {
			continue
		} else {
			// Ensure tokens are unique
			authToken := uniqueRandomString(&authTokens, secrets.Salt)
			refreshToken := uniqueRandomString(&refreshTokens, secrets.Salt)

			handler.secrets[record[0]] = PlainPasswordToken{
				IntegrityCheck: integrity,
				Password:       record[1],
				Token: Token{
					Token:           authToken,
					RefreshToken:    refreshToken,
					ExpireEpochTime: 0,
				},
			}
		}
	}

	return handler.secrets, nil
}

// UsernamePasswordLogin login via username / password
func (handler *PlainTextAuthHandler) UsernamePasswordLogin(username string, password string) (token Token, err error) {
	secrets, err := handler.GetSecrets()

	if err != nil {
		return Token{}, err
	}

	if secret, ok := secrets[username]; ok && secret.Password == password {
		return secret.Token, nil
	}

	return Token{}, errors.New("Wrong username and / or password")
}

// RefreshToken manage the token refresh
func (handler *PlainTextAuthHandler) RefreshToken(refreshToken string) (newToken string, err error) {
	secrets, err := handler.GetSecrets()

	if err != nil {
		return "", err
	}

	for _, secret := range secrets {
		if secret.Token.RefreshToken == refreshToken {
			return secret.Token.Token, nil
		}
	}

	return "", errors.New("Invalid refresh token")
}

// Logout manage the logout process
func (handler *PlainTextAuthHandler) Logout(token string) (loggedOut bool, err error) {
	for key, secret := range handler.secrets {
		if secret.Token.Token == token {
			// Remove token and refresh token to avoid further interaction
			secret.Token.Token = ""
			secret.Token.RefreshToken = ""
			// Remove the integrity check to force a new token tuple the next time
			secret.IntegrityCheck = ""

			// assign the modified secret to the handler's memory
			handler.secrets[key] = secret

			return true, nil
		}
	}

	return false, errors.New("Invalid token")
}
