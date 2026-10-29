package auth

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	wccfg "github.com/WebCommanderCMS/webcommander-sdk-go/v4/config"
	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/redact"
)

// TokenProvider fetches, caches and renews the access token.
type TokenProvider struct {
	config    wccfg.SdkConfig
	transport *http.HTTPTransport
	store     TokenStore
	tokenPath string
}

const refreshGrantType = "refresh_token"

// TokenPathTemplate is the token endpoint path with its version slot unfilled; BuildTokenPath fills it.
const TokenPathTemplate = "/api/{api_version}/oauth2/token" //#nosec G101 -- the token endpoint path, not a token

// BuildTokenPath builds the token endpoint path for an API version: "v4" gives "/api/v4/oauth2/token".
func BuildTokenPath(apiVersion string) string {
	return strings.Replace(TokenPathTemplate, "{api_version}", apiVersion, 1)
}

// NewTokenProvider builds a provider over a transport and a store.
func NewTokenProvider(config wccfg.SdkConfig, transport *http.HTTPTransport, store TokenStore) *TokenProvider {
	return &TokenProvider{
		config:    config,
		transport: transport,
		store:     store,
		tokenPath: BuildTokenPath(config.APIVersion),
	}
}

// TokenURL is the absolute token endpoint URL.
func (p *TokenProvider) TokenURL() string { return p.config.Auth.BaseURL + p.tokenPath }

// GetAccessToken returns a usable token, fetching one only when the cached one is missing or expiring.
func (p *TokenProvider) GetAccessToken(ctx context.Context) (string, error) {
	stored, err := p.store.ReadTokenRecord()
	if err != nil {
		return "", err
	}
	tokenRecord, ok, invalid := TokenRecordFromMap(stored)
	if invalid != nil {
		return "", invalid
	}
	if ok && !tokenRecord.NeedsRefresh() {
		return tokenRecord.AccessToken, nil
	}
	return p.renew(ctx, "", true)
}

// RefreshAccessToken renews after the API rejected a token, accepting only one different from it.
func (p *TokenProvider) RefreshAccessToken(ctx context.Context, staleAccessToken string) (string, error) {
	return p.renew(ctx, staleAccessToken, false)
}

// acquireRefreshLock takes the store's refresh lock, passing ctx through when the store honours it.
func (p *TokenProvider) acquireRefreshLock(ctx context.Context) (TokenRefreshScope, error) {
	if store, ok := p.store.(ContextTokenStore); ok {
		return store.AcquireTokenRefreshLockContext(ctx)
	}
	return p.store.AcquireTokenRefreshLock()
}

// renew is the double-checked refresh.
func (p *TokenProvider) renew(ctx context.Context, staleAccessToken string, anyCachedTokenWillDo bool) (string, error) {
	scope, err := p.acquireRefreshLock(ctx)
	if err != nil {
		return "", err
	}
	defer scope.Release()

	stored, err := scope.ReadTokenRecord()
	if err != nil {
		return "", err
	}
	tokenRecord, ok, invalid := TokenRecordFromMap(stored)
	if invalid != nil {
		return "", invalid
	}
	if ok && !tokenRecord.NeedsRefresh() {
		// Any fresh token answers GetAccessToken; a refresh replaces a rejected token, so it needs a different one.
		if anyCachedTokenWillDo || (staleAccessToken != "" && tokenRecord.AccessToken != staleAccessToken) {
			return tokenRecord.AccessToken, nil
		}
	}

	var renewed *TokenRecord
	if ok && tokenRecord.RefreshToken != "" {
		renewed = p.requestTokenWithRefreshGrant(ctx, tokenRecord.RefreshToken)
	}
	if renewed == nil {
		fresh, err := p.requestTokenWithCredentials(ctx)
		if err != nil {
			return "", err
		}
		renewed = fresh
	}

	if err := scope.WriteTokenRecord(renewed.ToMap()); err != nil {
		// Reported, not absorbed: a cache that never takes a write turns every call into a token grant.
		return "", err
	}
	return renewed.AccessToken, nil
}

func (p *TokenProvider) requestTokenWithCredentials(ctx context.Context) (*TokenRecord, error) {
	response, err := p.transport.SendHTTPRequest(ctx, http.TransportRequest{
		HTTPMethod: "POST", URL: p.TokenURL(), EndpointTemplate: TokenPathTemplate,
		JSONBody: p.credentialsPayload(), Headers: p.tokenHeaders(),
	})
	if err != nil {
		return nil, err
	}
	if !response.Succeeded() {
		messages := response.GetErrorMessages()
		// The status and count, not the endpoint's server-controlled words, which Error() would render.
		return nil, wcerrors.NewAuthResponseError(
			fmt.Sprintf("token request failed with %d (%d reported)", response.StatusCode, len(messages)),
			response.StatusCode, response.ErrorType, redact.RedactValue(response.Map()),
			redact.RedactMapping(response.Headers), response.RequestInfo, messages)
	}
	return buildTokenRecord(response, "")
}

// requestTokenWithRefreshGrant is the cheap grant; any failure falls back to full credentials.
func (p *TokenProvider) requestTokenWithRefreshGrant(ctx context.Context, refreshToken string) *TokenRecord {
	payload := map[string]any{
		"grant_type":    refreshGrantType,
		"refresh_token": refreshToken,
		"redirect_uri":  p.config.Auth.RedirectURI,
		"client_id":     p.config.Auth.ClientID,
		"client_secret": p.config.Auth.ClientSecret,
	}
	response, err := p.transport.SendHTTPRequest(ctx, http.TransportRequest{
		HTTPMethod: "POST", URL: p.TokenURL(), EndpointTemplate: TokenPathTemplate,
		JSONBody: payload, Headers: p.tokenHeaders(),
	})
	if err != nil || !response.Succeeded() {
		return nil
	}
	tokenRecord, buildError := buildTokenRecord(response, refreshToken)
	if buildError != nil {
		return nil
	}
	return tokenRecord
}

func (p *TokenProvider) credentialsPayload() map[string]any {
	return map[string]any{
		"grant_type":    p.config.Auth.GrantType,
		"auth_string":   p.config.Auth.AuthString,
		"redirect_uri":  p.config.Auth.RedirectURI,
		"client_id":     p.config.Auth.ClientID,
		"client_secret": p.config.Auth.ClientSecret,
	}
}

func (p *TokenProvider) tokenHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
		"User-Agent":   p.config.UserAgent,
	}
}

func buildTokenRecord(response *http.ApiResponse, previousRefreshToken string) (*TokenRecord, error) {
	body, isObject := response.Map().(map[string]any)
	if !isObject {
		return nil, grantRefusal(response, fmt.Sprintf(
			"the token response is %s rather than an object", http.JSONKind(response.Map())))
	}
	rawAccess, present := body["access_token"]
	if !present || rawAccess == nil {
		return nil, grantRefusal(response, "the token response carried no access_token")
	}
	access, isString := rawAccess.(string)
	if !isString {
		return nil, grantRefusal(response, fmt.Sprintf(
			"the token response's access_token is %s rather than a string", http.JSONKind(rawAccess)))
	}
	refresh := previousRefreshToken
	// RFC 6749 §6: a response without a refresh token means keep the one already held.
	if rawRefresh, present := body["refresh_token"]; present && rawRefresh != nil {
		text, isText := rawRefresh.(string)
		if !isText {
			return nil, grantRefusal(response, fmt.Sprintf(
				"the token response's refresh_token is %s rather than a string", http.JSONKind(rawRefresh)))
		}
		if text != "" {
			refresh = text
		}
	}
	timeToLiveSeconds := readTTLSeconds(body["expires_in"])
	record, err := NewTokenRecord(access, refresh,
		nowSeconds()+timeToLiveSeconds,
		math.Min(ExpirySkewSeconds, timeToLiveSeconds/2)) // a short lifetime leaves a usable window
	if err != nil {
		return nil, grantRefusal(response, err.Error())
	}
	return &record, nil
}

// grantRefusal is the one refusal for a grant response this SDK will not build a record from.
func grantRefusal(response *http.ApiResponse, problem string) error {
	return wcerrors.NewAuthResponseError(problem,
		response.StatusCode, response.ErrorType, redact.RedactValue(response.Map()),
		nil, response.RequestInfo, nil)
}

// readTTLSeconds reads the lifetime; NaN, infinite, non-positive or unparseable take the default, capped at a day.
func readTTLSeconds(expiresIn any) float64 {
	timeToLiveSeconds := DefaultTokenTTLSeconds
	if n, ok := finiteNumber(expiresIn); ok {
		timeToLiveSeconds = n
	} else if s, ok := expiresIn.(string); ok {
		if parsed, err := strconv.ParseFloat(s, 64); err == nil {
			timeToLiveSeconds = parsed
		}
	}
	if wccfg.IsNaNOrInf(timeToLiveSeconds) || timeToLiveSeconds <= 0 {
		timeToLiveSeconds = DefaultTokenTTLSeconds
	}
	return math.Min(timeToLiveSeconds, MaxTokenTTLSeconds)
}
