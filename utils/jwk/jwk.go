package jwk

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type JWK_DATA struct {
	PRV_SET jwk.Set
	PUB_SET jwk.Set
}

func Read() jwk.Set {
	if _, err := os.Stat("./jwk.json"); os.IsNotExist(err) {
		// create file if not exists
		file, err := os.Create("./jwk.json")
		if err != nil {
			fmt.Printf("failed to create jwk.json: %s\n", err)
		}

		file.Write([]byte(`{}`))
		defer file.Close()
	}

	set, err := jwk.ReadFile("./jwk.json")
	if err != nil {
		if err.Error() == "failed to unmarshal JWK set: failed to parse sole key in key set" {
			fmt.Println("jwk.json is empty, creating new set")

			// Create a new set if the file is empty
			set = jwk.NewSet()
			key, err := Generate()
			if err != nil {
				fmt.Printf("failed to generate new JWK: %s\n", err)
			}

			set.AddKey(key)
		} else {
			fmt.Printf("failed to read jwk.json: %s\n", err)
			panic(err)
		}
	}

	return set
}

func PublicSetOf(set jwk.Set) jwk.Set {
	pubSet, err := jwk.PublicSetOf(set)
	if err != nil {
		fmt.Printf("jwk.PublicSetOf failed: %s\n", err)
		panic(err)
	}

	return pubSet
}

func Generate() (jwk.Key, error) {
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

func Rotate(jwks *JWK_DATA) {
	key, err := Generate()
	if err != nil {
		fmt.Printf("failed to generate new JWK: %s\n", err)
	}

	jwks.PRV_SET.AddKey(key)

	pubKey, err := key.PublicKey()
	if err != nil {
		fmt.Printf("failed to get public key from JWK: %s\n", err)
	}

	jwks.PUB_SET.AddKey(pubKey)

	enc, err := json.MarshalIndent(jwks.PRV_SET, "", "    ")
	if err != nil {
		fmt.Printf("failed to marshal JWK set: %s\n", err)
	}

	os.WriteFile("./jwk.json", []byte(enc), 0644)
}
