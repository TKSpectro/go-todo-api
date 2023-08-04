package jwk

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"

	"github.com/TKSpectro/go-todo-api/config"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type JWK_DATA struct {
	PRV_SET jwk.Set
	PUB_SET jwk.Set
}

func Read() jwk.Set {
	path := "./jwk.json"
	if config.ROOT_PATH != "" {
		path = config.ROOT_PATH + "/jwk.json"
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		// create file if not exists
		file, err := os.Create(path)
		if err != nil {
			fmt.Printf("failed to create jwk.json: %s\n", err)
		}

		file.Write([]byte(`{}`))
		defer file.Close()
	}

	set, err := jwk.ReadFile(path)
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
	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Printf("failed to generate new ECDSA private key: %s\n", err)
		return nil, err
	}

	key, err := jwk.FromRaw(ecdsaKey)
	if err != nil {
		fmt.Printf("failed to create symmetric key: %s\n", err)
		return nil, err
	}

	key.Set(jwk.KeyIDKey, uuid.New().String())
	key.Set(jwk.AlgorithmKey, jwa.ES256K)

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

	path := "./jwk.json"
	if config.ROOT_PATH != "" {
		path = config.ROOT_PATH + "/jwk.json"
	}

	os.WriteFile(path, []byte(enc), 0644)
}
