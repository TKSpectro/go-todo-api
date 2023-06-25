package utils

import "crypto/rand"

const ALLOWED_SECRET_TOKEN_CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func RandomString(length int, allowedCharsOverwrite string) string {
	var allowedChars string
	if len(allowedCharsOverwrite) > 0 {
		allowedChars = allowedCharsOverwrite
	} else {
		allowedChars = ALLOWED_SECRET_TOKEN_CHARS
	}

	ll := len(allowedChars)
	// 8 comes from db max length of secretToken
	result := make([]byte, length)

	rand.Read(result)
	for i := 0; i < length; i++ {
		result[i] = allowedChars[int(result[i])%ll]
	}

	return string(result)
}
