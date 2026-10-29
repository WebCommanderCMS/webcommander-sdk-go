package auth

import (
	"encoding/json"
	"fmt"
	"time"

	wccfg "github.com/WebCommanderCMS/webcommander-sdk-go/v4/config"
	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/redact"
)

// The token lifetime bounds.
const (
	// DefaultTokenTTLSeconds is the lifetime assumed when the endpoint gives none.
	DefaultTokenTTLSeconds = 3600.0
	// ExpirySkewSeconds is how early a token is treated as expired.
	ExpirySkewSeconds = 60.0
	// MaxTokenTTLSeconds is the furthest ahead a cached expiry is believed (one day).
	MaxTokenTTLSeconds = 24 * 60 * 60.0
)

// TokenRecord is an access token and when it stops being usable.
type TokenRecord struct {
	// AccessToken is the bearer token.
	AccessToken string
	// RefreshToken is the refresh token, or "".
	RefreshToken string
	// ExpiresAt is the epoch seconds when the token expires.
	ExpiresAt float64
	// SkewSeconds is how early to treat it as expired — carried per record.
	SkewSeconds float64
}

// NeedsRefresh reports whether the token should be treated as expired.
func (r TokenRecord) NeedsRefresh() bool {
	now := nowSeconds()
	if r.ExpiresAt-now > MaxTokenTTLSeconds {
		return true // a wild clock's stamp is not believed
	}
	return now >= r.ExpiresAt-r.SkewSeconds
}

// ToMap is the plain mapping a token store persists — the shared cache format.
func (r TokenRecord) ToMap() map[string]any {
	var refresh any
	if r.RefreshToken != "" {
		refresh = r.RefreshToken
	}
	return map[string]any{
		"access_token":  r.AccessToken,
		"refresh_token": refresh,
		"expires_at":    r.ExpiresAt,
		"skew_seconds":  r.SkewSeconds,
	}
}

// String redacts both tokens and reports the record's state rather than its exact expiry.
func (r TokenRecord) String() string {
	refresh := "null"
	if r.RefreshToken != "" {
		refresh = redact.Redacted
	}
	state := "usable"
	if r.NeedsRefresh() {
		state = "due for refresh"
	}
	return fmt.Sprintf("TokenRecord(accessToken=%q, refreshToken=%s, state=%s)",
		redact.Redacted, refresh, state)
}

// NewTokenRecord builds a record, holding each token to the grammar of the header or response it travels in.
func NewTokenRecord(accessToken, refreshToken string, expiresAt, skewSeconds float64) (TokenRecord, error) {
	if problem := bearerCredentialProblem(accessToken); problem != "" {
		return TokenRecord{}, wcerrors.NewAuthError(
			"accessToken: not a value this SDK can send in an Authorization: Bearer header — " +
				problem + ". The value is not shown: it is a credential")
	}
	if refreshToken != "" {
		if problem := refreshCredentialProblem(refreshToken); problem != "" {
			return TokenRecord{}, wcerrors.NewAuthError(
				"refreshToken: not a value a token response may carry — " + problem +
					". The value is not shown: it is a credential")
		}
	}
	if wccfg.IsNaNOrInf(expiresAt) {
		return TokenRecord{}, wcerrors.NewAuthError(fmt.Sprintf(
			"expiresAt: expected epoch seconds, got %s", wcerrors.FormatSeconds(expiresAt)))
	}
	if wccfg.IsNaNOrInf(skewSeconds) || skewSeconds < 0 {
		return TokenRecord{}, wcerrors.NewAuthError(fmt.Sprintf(
			"skewSeconds: must be a number of seconds and at least 0, got %s",
			wcerrors.FormatSeconds(skewSeconds)))
	}
	return TokenRecord{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		SkewSeconds:  min(skewSeconds, MaxTokenTTLSeconds),
	}, nil
}

// bearerCredentialProblem reports why a value cannot be sent as a Bearer credential, or "" when it can.
func bearerCredentialProblem(value string) string {
	if value == "" {
		return "it is empty"
	}
	padding := 0
	for index := len(value) - 1; index >= 0 && value[index] == '='; index-- {
		padding++
	}
	body := value[:len(value)-padding]
	if body == "" {
		return "it is '=' padding and nothing else"
	}
	for index := 0; index < len(body); index++ {
		switch c := body[index]; {
		case c >= 'A' && c <= 'Z', c >= 'a' && c <= 'z', c >= '0' && c <= '9',
			c == '-', c == '.', c == '_', c == '~', c == '+', c == '/':
		case c == '=':
			return fmt.Sprintf("byte %d is '=' padding, which RFC 6750 allows only at the end", index)
		default:
			return fmt.Sprintf("byte %d is not one of RFC 6750's b64token characters "+
				"(letters, digits, '-', '.', '_', '~', '+', '/', then '=' padding)", index)
		}
	}
	return ""
}

// refreshCredentialProblem reports why a value cannot be a refresh credential, or "" when it can.
func refreshCredentialProblem(value string) string {
	if value == "" {
		return "it is empty"
	}
	for index := 0; index < len(value); index++ {
		if c := value[index]; c < 0x20 || c > 0x7e {
			return fmt.Sprintf("byte %d is outside RFC 6749's visible token-response characters "+
				"(%%x20-7E)", index)
		}
	}
	return ""
}

// TokenRecordFromMap rebuilds a record: (record, true, nil), none (zero, false, nil), or unusable (err).
func TokenRecordFromMap(record map[string]any) (TokenRecord, bool, error) {
	if len(record) == 0 {
		return TokenRecord{}, false, nil
	}
	rawAccess, hasAccess := record["access_token"]
	if !hasAccess || rawAccess == nil {
		return TokenRecord{}, false, nil
	}
	accessToken, isString := rawAccess.(string)
	if !isString {
		return TokenRecord{}, false, corruptRecordError(fmt.Sprintf(
			"its access_token is %s rather than a string", http.JSONKind(rawAccess)))
	}
	if problem := bearerCredentialProblem(accessToken); problem != "" {
		return TokenRecord{}, false, corruptRecordError(
			"its access_token is not a value this SDK can send in an Authorization: Bearer " +
				"header — " + problem)
	}
	expiresAt, ok := finiteNumber(record["expires_at"])
	if !ok {
		return TokenRecord{}, false, corruptRecordError(fmt.Sprintf(
			"its expires_at is %s rather than a finite number of epoch seconds",
			http.JSONKind(record["expires_at"])))
	}
	refresh := ""
	if rawRefresh, present := record["refresh_token"]; present && rawRefresh != nil {
		text, isText := rawRefresh.(string)
		if !isText {
			return TokenRecord{}, false, corruptRecordError(fmt.Sprintf(
				"its refresh_token is %s rather than a string", http.JSONKind(rawRefresh)))
		}
		// An empty refresh token means none, as ToMap writes null and a hand-edited file may write "".
		if text != "" {
			if problem := refreshCredentialProblem(text); problem != "" {
				return TokenRecord{}, false, corruptRecordError(
					"its refresh_token is not a value a token response may carry — " + problem)
			}
			refresh = text
		}
	}
	skew, ok := finiteNumber(record["skew_seconds"])
	if !ok || skew < 0 {
		skew = ExpirySkewSeconds
	}
	skew = min(skew, MaxTokenTTLSeconds)
	return TokenRecord{AccessToken: accessToken, RefreshToken: refresh,
		ExpiresAt: expiresAt, SkewSeconds: skew}, true, nil
}

// corruptRecordError is the one refusal for a cache entry this SDK will not use.
func corruptRecordError(problem string) error {
	return wcerrors.NewAuthError(fmt.Sprintf(
		"the cached token entry cannot be used because %s; clear the token to rebuild it", problem))
}

// finiteNumber reads a finite float from a JSON-ish value.
func finiteNumber(value any) (float64, bool) {
	switch v := value.(type) {
	case json.Number:
		f, err := v.Float64()
		if err != nil || wccfg.IsNaNOrInf(f) {
			return 0, false
		}
		return f, true
	case float64:
		if wccfg.IsNaNOrInf(v) {
			return 0, false
		}
		return v, true
	case float32:
		return finiteNumber(float64(v))
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case int32:
		return float64(v), true
	default:
		return 0, false
	}
}

func nowSeconds() float64 {
	return float64(time.Now().UnixNano()) / 1e9
}
