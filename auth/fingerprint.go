package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"unicode/utf8"
)

// BuildCredentialsNamespace is the fixed 16-hex namespace these five credential values key the token cache under.
func BuildCredentialsNamespace(baseURL, clientID, clientSecret, authString, grantType string) string {
	var b strings.Builder
	for _, value := range []string{"v3", baseURL, clientID, clientSecret, authString, grantType} {
		b.WriteString(strconv.Itoa(utf8.RuneCountInString(value)))
		b.WriteByte(':')
		b.WriteString(value)
	}
	sum := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(sum[:])[:16]
}
