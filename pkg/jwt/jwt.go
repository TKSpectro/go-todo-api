package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/TKSpectro/go-todo-api/app/models"
	"github.com/TKSpectro/go-todo-api/app/types"
	"github.com/TKSpectro/go-todo-api/config"
	"github.com/TKSpectro/go-todo-api/core"
	"github.com/TKSpectro/go-todo-api/pkg/jwk"
	"github.com/lestrrat-go/jwx/v2/jwa"
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

// Instance of the pub/priv key sets. We keep them in memory so we don't have IO overhead on every request
var JWKS jwk.JWK_DATA

// Init initializes the jwk set from the jwk.json file
func Init() {
	JWKS = jwk.JWK_DATA{
		PRV_SET: jwk.Read(),
	}

	JWKS.PUB_SET = jwk.PublicSetOf(JWKS.PRV_SET)

	fmt.Println("JWKs initialized. Found keys:", JWKS.PRV_SET.Len())
}

// Generate generates both a token and a refresh token
func Generate(account *models.Account) (types.AuthResponseBody, error) {
	token, err := jwt.NewBuilder().
		Issuer(`github.com/TKSpectro/go-todo-api`).
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(config.JWT_TOKEN_EXP)).
		Claim(CLAIM_ACCOUNT_ID, account.ID).
		Claim(CLAIM_TYPE, "auth").
		Claim(CLAIM_SECRET, account.TokenSecret).
		Build()
	if err != nil {
		return types.AuthResponseBody{}, core.RequestErrorFrom(&core.TOKEN_GENERATION_ERROR, err.Error())
	}

	refreshToken, err := jwt.NewBuilder().
		Issuer(`github.com/TKSpectro/go-todo-api`).
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(config.JWT_REFRESH_EXP)).
		Claim(CLAIM_ACCOUNT_ID, account.ID).
		Claim(CLAIM_TYPE, "refresh").
		Claim(CLAIM_SECRET, account.TokenSecret).
		Build()
	if err != nil {
		return types.AuthResponseBody{}, core.RequestErrorFrom(&core.TOKEN_GENERATION_ERROR, err.Error())
	}

	// Get the last key in the set
	jwkKey, ok := JWKS.PRV_SET.Key(JWKS.PRV_SET.Len() - 1)
	if !ok {
		return types.AuthResponseBody{}, core.RequestErrorFrom(&core.TOKEN_GENERATION_ERROR, "failed to get last key in set")
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.ES256K, jwkKey))
	if err != nil {
		return types.AuthResponseBody{}, core.RequestErrorFrom(&core.TOKEN_GENERATION_ERROR, err.Error())
	}

	signedRefresh, err := jwt.Sign(refreshToken, jwt.WithKey(jwa.ES256K, jwkKey))
	if err != nil {
		return types.AuthResponseBody{}, core.RequestErrorFrom(&core.TOKEN_GENERATION_ERROR, err.Error())
	}

	return types.AuthResponseBody{
		Token:        string(signed),
		RefreshToken: string(signedRefresh),
	}, nil
}

// Parse parses the jwt token (Validate against the public keys and return the token)
func Parse(token string) (jwt.Token, error) {
	// When parsing we do it against the public key
	tok, err := jwt.Parse([]byte(token), jwt.WithKeySet(JWKS.PUB_SET))
	if err != nil {
		fmt.Printf("jwt.Parse failed: %s\n", err)
		return nil, err
	}

	return tok, nil
}

// Verify Verifies the token and returns the payload
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

// RotateJWK rotates the jwk set (generates a new key and adds it to the set)
func RotateJWK() {
	jwk.Rotate(&JWKS)
}
