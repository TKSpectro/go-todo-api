package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/TKSpectro/go-todo-api/config"
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

// your JWK
const jwkStr = `{
		"kty": "RSA",
		"n": "mmO0OvOPQ53HRxV4eHOkTTxLVfk6zcq8KAD86gbnydYBNO_Si4Q1twyvefd58-BaO4N4NCEA97QrYm57ThKCe8agLGwWPHhxgbu_SAuYQehXxkf4sWy7Q17kGFG5k5AfQGZBqTY-YaawQqLlF6ILVbWab_AoEF4yB7pI3AnNnXs",
		"e": "AQAB",
		"d": "RzsrI2vONJcuIyjPzVslehEQfRkhPWOFTjuudNc8yA25vs_LZ11XXx42M-KvXIqtdvngUsTLan2w6pgowcuecX3t_2wUx0GJJgARfkN7gsWIS3CyXZBEEMjLGVU4vHt5zNE3GJKo3hb1TwEiulpL_Ix6hfcTSJpEaBWrBxjxV-E",
		"p": "5EA0bi6ui1H1wsG85oc7i9O7UH58WPIK_ytzBWXFIwcaSFFBqqNYNnZaHFsMe4cbHSBgShWHO3UueGVgOKmB8Q",
		"q": "rSi7CosQZmj_RFIYW10ef7XTZsdpIdOXV9-1dThAJUvkslKiTfdU7T0IYYsJ2K58ekJqdpcoKAVLB2SZVvdqKw",
		"dp": "S9yjEHPng1qsShzGQgB0ZBbtTOWdQpq_2OuCAStACFJWA-8t2h8MNJ3FeWMxlOTkuBuIpVbeaX6bAV0ATBTaoQ",
		"dq": "ZssMJhkh1jm0d-FoVix0Y4oUAiqUzaDnciH6faiz47AnBnkporEV-HPH2ugII1qJyKZOvzHCg-eIf84HfWoI2w",
		"qi": "lyVz1HI2b1IjzOMENkmUTaVEO6DM6usZi3c3_MobUUM05yyBhnHtPjWzqWn1uJ_Gt5bkJDdcpfvmkPAhKWEU9Q"
	}`

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

	//convert jwk in bytes and return a new key
	jwkKey, err := jwk.ParseKey([]byte(jwkStr))
	if err != nil {
		fmt.Printf("jwk.ParseKey failed: %s\n", err)
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
	key, err := jwk.ParseKey([]byte(jwkStr))
	if err != nil {
		fmt.Printf("jwk.ParseKey failed: %s\n", err)
		panic(err)
	}

	// TODO: Gotta fix this
	tok, err := jwt.Parse([]byte(token), jwt.WithKey(jwa.RS256, key))
	if err != nil {
		fmt.Printf("jwt.Parse failed: %s\n", err)
		panic(err)
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
