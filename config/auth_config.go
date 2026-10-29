package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strings"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/redact"
)

// TokenCacheMode selects where access and refresh tokens are cached between calls.
type TokenCacheMode int

// The cache modes. The zero value is TokenCachePlatformDefault.
const (
	// TokenCachePlatformDefault uses the platform's per-user cache directory.
	TokenCachePlatformDefault TokenCacheMode = iota
	// TokenCacheMemoryOnly keeps tokens in memory for the life of the process.
	TokenCacheMemoryOnly
	// TokenCacheExplicitPath uses the file named by AuthConfig.TokenCachePath.
	TokenCacheExplicitPath
)

// UseDefaultTokenCache is an alias for TokenCachePlatformDefault.
const UseDefaultTokenCache = TokenCachePlatformDefault

// String renders the mode as its stable name, which logged config shows and log filters match.
func (m TokenCacheMode) String() string {
	switch m {
	case TokenCacheMemoryOnly:
		return "memoryOnly"
	case TokenCacheExplicitPath:
		return "explicitPath"
	default:
		return "platformDefault"
	}
}

// DefaultClientName is the implicit name of an unnamed credentials block; Normalised does not apply it.
const DefaultClientName = "default"

// DefaultGrantType is the OAuth grant type Normalised applies when GrantType is blank.
const DefaultGrantType = "client_credentials"

// AuthConfig is who you are and which instance you are talking to.
type AuthConfig struct {
	// BaseURL is the instance origin (trailing slashes stripped on normalisation).
	BaseURL string `wc:"required"`
	// ClientID is the OAuth client id.
	ClientID string `wc:"required"`
	// ClientSecret is the OAuth client secret; redacted by String.
	ClientSecret string `wc:"required"`
	// RedirectURI is the registered redirect URI.
	RedirectURI string `wc:"required"`
	// AuthString is the base64 admin credential; redacted by String.
	AuthString string `wc:"required"`
	// ClientName labels this instance in logs and selects it by name.
	ClientName string `wc:"required"`
	// GrantType is the OAuth grant type ("client_credentials").
	GrantType string
	// TokenCacheMode selects where tokens are cached; zero = platform default.
	TokenCacheMode TokenCacheMode
	// TokenCachePath is the explicit cache file; naming it promotes the mode to explicit path.
	TokenCachePath string
}

// A known grant type is a protocol constant and prints as itself; anything else is only described.
var knownGrantTypes = map[string]bool{"client_credentials": true, "refresh_token": true}

// String names which values are set and no part of any, so it is safe in any diagnostic.
func (c AuthConfig) String() string {
	grantType := redact.Describe(c.GrantType)
	if knownGrantTypes[c.GrantType] {
		grantType = fmt.Sprintf("%q", c.GrantType)
	}
	return fmt.Sprintf(
		"AuthConfig(baseUrl=%s, clientId=%s, clientSecret=%q, redirectUri=%s, authString=%q, "+
			"clientName=%q, grantType=%s, tokenCacheMode=%s, tokenCachePath=%s)",
		redact.Describe(c.BaseURL), redact.Describe(c.ClientID), redact.Redacted,
		redact.Describe(c.RedirectURI), redact.Redacted, redact.RedactValue(c.ClientName),
		grantType, c.TokenCacheMode, redact.DescribePath(c.TokenCachePath))
}

// Normalised returns a copy with defaults applied and the base URL and cache mode resolved.
func (c AuthConfig) Normalised() AuthConfig {
	out := c
	out.BaseURL = strings.TrimRight(strings.TrimSpace(c.BaseURL), "/")
	// ClientName is left blank so Validate reports it, rather than colliding with "default".
	if out.GrantType == "" {
		out.GrantType = DefaultGrantType
	}
	if out.TokenCacheMode == TokenCachePlatformDefault && out.TokenCachePath != "" {
		out.TokenCacheMode = TokenCacheExplicitPath
	}
	if out.TokenCacheMode != TokenCacheExplicitPath {
		out.TokenCachePath = ""
	}
	return out
}

// Validate reports every problem at once (call on a Normalised copy).
func (c AuthConfig) Validate() error {
	var problems []string
	var details []string

	var missing []string
	for _, required := range []struct{ name, value string }{
		{"baseUrl", c.BaseURL},
		{"clientId", c.ClientID},
		{"clientSecret", c.ClientSecret},
		{"redirectUri", c.RedirectURI},
		{"authString", c.AuthString},
	} {
		if isBlank(required.value) {
			missing = append(missing, required.name)
		}
	}
	if len(missing) > 0 {
		problems = append(problems, "missing required values: "+strings.Join(missing, ", "))
	}
	if !isBlank(c.BaseURL) {
		if found := baseURLProblems(c.BaseURL); len(found) > 0 {
			problems = append(problems, found...)
			details = append(details, "baseUrl: "+redact.RedactText(c.BaseURL))
		}
	}
	if isBlank(c.ClientName) {
		problems = append(problems, "clientName: must not be blank")
	}
	if isBlank(c.GrantType) {
		problems = append(problems, "grantType: must not be blank")
	}
	if c.TokenCacheMode == TokenCacheExplicitPath {
		if isBlank(c.TokenCachePath) {
			problems = append(problems, "tokenCachePath: must not be blank — use TokenCacheMemoryOnly to "+
				"keep the token in memory, or a path to cache it")
		} else if info, err := os.Stat(c.TokenCachePath); err == nil && info.IsDir() {
			// A directory here makes every write fail, silently costing one authentication per call.
			problems = append(problems, fmt.Sprintf(
				"tokenCachePath: %s names an existing directory, so nothing could be written to it — "+
					"pass the file to write, or use TokenCacheMemoryOnly to keep the token in memory",
				redact.DescribePath(c.TokenCachePath)))
			details = append(details, "tokenCachePath: "+c.TokenCachePath)
		}
	}

	if len(problems) > 0 {
		// Trimmed, not just tested for blankness: a padded name would label the report `instance prod :`.
		label := strings.TrimSpace(c.ClientName)
		if label == "" {
			label = "<unnamed>"
		}
		refusal := wcerrors.NewConfigError(fmt.Sprintf("instance %s: %s", label, strings.Join(problems, "; ")))
		refusal.CauseText = strings.Join(details, "; ")
		return refusal
	}
	return nil
}

// WithOverrides derives a new AuthConfig, changing the snake_case-named fields in changes.
func (c AuthConfig) WithOverrides(changes map[string]any) (AuthConfig, error) {
	out := c
	if len(changes) > 0 {
		if err := assignFields(reflect.ValueOf(&out).Elem(), changes, "overrides"); err != nil {
			return AuthConfig{}, err
		}
	}
	out = out.Normalised()
	if err := out.Validate(); err != nil {
		return AuthConfig{}, err
	}
	return out, nil
}

func baseURLProblems(baseURL string) []string {
	u, err := url.Parse(baseURL)
	if err != nil {
		// Unwrapped: *url.Error prepends the whole input, and an unparseable base URL may carry userinfo.
		var parseErr *url.Error
		if errors.As(err, &parseErr) {
			err = parseErr.Err
		}
		return []string{fmt.Sprintf("baseUrl: is not a URL (%v)", err)}
	}
	if u.Scheme == "" {
		return []string{"baseUrl: has no scheme — must start with http:// or https://"}
	}
	// The refusals below name which part is wrong without echoing the value, which may carry userinfo.
	scheme := strings.ToLower(u.Scheme)
	if (scheme != "http" && scheme != "https") || u.Hostname() == "" {
		return []string{fmt.Sprintf(
			"baseUrl: must start with http:// or https:// and include a host, got %s",
			redact.Describe(baseURL))}
	}
	var problems []string
	if u.User != nil {
		problems = append(problems,
			"baseUrl: the host carries embedded credentials, which would replace the bearer token; "+
				"strip the userinfo")
	}
	// A ';' path-parameter run is not refused: it is a path, not a pasted query or fragment.
	if u.RawQuery != "" || u.Fragment != "" {
		which := "a query string"
		switch {
		case u.RawQuery == "":
			which = "a fragment"
		case u.Fragment != "":
			which = "a query string and a fragment"
		}
		problems = append(problems, fmt.Sprintf(
			"baseUrl: must be the instance origin only, and this one carries %s", which))
	}
	return problems
}
