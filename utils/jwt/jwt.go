package jwt

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/TKSpectro/go-todo-api/config"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
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

	token, err := jwt.NewBuilder().
		Issuer(`github.com/TKSpectro/go-todo-api`).
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(v)).
		Claim(CLAIM_ACCOUNT_ID, payload.AccountID).
		Claim(CLAIM_TYPE, payload.Type).
		Claim(CLAIM_SECRET, payload.Secret).
		Build()
	if err != nil {
		fmt.Printf("failed to build token: %s\n", err)
		panic(err)
	}

	keySet, err := jwk.ReadFile("./jwk.json")
	if err != nil {
		fmt.Printf("failed to read jwk.json: %s\n", err)
		panic(err)
	}

	// Get the last key in the set
	jwkKey, ok := keySet.Key(keySet.Len() - 1)
	if !ok {
		fmt.Printf("failed to get last key in set: %s\n", err)
		panic(err)
	}

	// Sign a JWT!
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, jwkKey))
	if err != nil {
		fmt.Printf("failed to sign token: %s\n", err)
		panic(err)
	}

	return string(signed)
}

func Parse(token string) (jwt.Token, error) {
	raw, err := os.ReadFile("./jwk.json")
	if err != nil {
		fmt.Printf("failed to read jwk.json: %s\n", err)
		return nil, err
	}

	privSet, err := jwk.Parse(raw)
	if err != nil {
		fmt.Printf("jwk.ParseKey failed: %s\n", err)
		return nil, err
	}

	pubSet, err := jwk.PublicSetOf(privSet)
	if err != nil {
		fmt.Printf("jwk.PublicSetOf failed: %s\n", err)
		return nil, err
	}

	// When parsing we do it against the public key
	tok, err := jwt.Parse([]byte(token), jwt.WithKeySet(pubSet))
	if err != nil {
		fmt.Printf("jwt.Parse failed: %s\n", err)
		return nil, err
	}

	return tok, nil
}

// Verify verifies the jwt token against the secret
func Verify(token string) (*TokenPayload, error) {
	tok, err := Parse(token)
	if err != nil {
		return nil, err
	}

	claims, err := tok.AsMap(context.Background())
	if err != nil {
		return nil, err
	}

	return &TokenPayload{
		AccountID: uint(claims[CLAIM_ACCOUNT_ID].(float64)),
		Type:      claims[CLAIM_TYPE].(string),
		Secret:    claims[CLAIM_SECRET].(string),
	}, nil
}

func GenerateNewJWK() (jwk.Key, error) {
	raw, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("failed to generate new RSA private key: %s\n", err)
		return nil, err
	}

	key, err := jwk.FromRaw(raw)
	if err != nil {
		fmt.Printf("failed to create symmetric key: %s\n", err)
		return nil, err
	}

	key.Set(jwk.KeyIDKey, uuid.New().String())
	key.Set(jwk.AlgorithmKey, jwa.RS256)

	return key, nil
}
