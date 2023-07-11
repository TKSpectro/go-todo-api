package jwt

import (
	"fmt"
	"time"

	"github.com/TKSpectro/go-todo-api/config"

	"github.com/golang-jwt/jwt/v5"
)

// TokenPayload defines the payload for the token
type TokenPayload struct {
	AccountID uint
	Type      string
	Secret    string
}

const (
	CLAIM_EXPIRE = "exp"

	CLAIM_ACCOUNT_ID = "accountID"
	CLAIM_TYPE       = "type"
	CLAIM_SECRET     = "tokenSecret"
)

// Generate generates the jwt token based on payload
func Generate(payload *TokenPayload) string {
	v, err := time.ParseDuration(config.JWT_TOKEN_EXP)
	if err != nil {
		panic("Invalid time duration. Should be time.ParseDuration string")
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		CLAIM_EXPIRE:     time.Now().Add(v).Unix(),
		CLAIM_ACCOUNT_ID: payload.AccountID,
		CLAIM_TYPE:       payload.Type,
		CLAIM_SECRET:     payload.Secret,
	})

	token, err := t.SignedString([]byte(config.JWT_TOKEN_SECRET))
	if err != nil {
		panic(err)
	}

	return token
}

func Parse(token string) (*jwt.Token, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.JWT_TOKEN_SECRET), nil
	})
}

// Verify verifies the jwt token against the secret
func Verify(token string) (*TokenPayload, error) {
	parsed, err := Parse(token)
	if err != nil {
		return nil, err
	}

	// Parsing token claims
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	return &TokenPayload{
		AccountID: uint(claims[CLAIM_ACCOUNT_ID].(float64)),
		Type:      claims[CLAIM_TYPE].(string),
		Secret:    claims[CLAIM_SECRET].(string),
	}, nil
}
