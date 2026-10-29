package config

import (
	"errors"
	"fmt"
	"maps"
	"math"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"runtime/debug"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/redact"
)

// fallbackVersion is reported when the build carries no module version. Bump it on every release.
const fallbackVersion = "4.0.0"

var sdkVersion = readModuleVersion()

// SDKVersion is the SDK's own version, reported in the default User-Agent header.
func SDKVersion() string { return sdkVersion }

// Version is SDKVersion under a shorter name — the same value.
func Version() string { return sdkVersion }

// DefaultUserAgent is the User-Agent this SDK sends unless SdkConfig.UserAgent replaces it.
func DefaultUserAgent() string { return "webcommander-sdk-go/" + sdkVersion }

func readModuleVersion() string {
	const modulePath = "github.com/WebCommanderCMS/webcommander-sdk-go/v4"
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return fallbackVersion
	}
	candidates := append([]*debug.Module{&info.Main}, info.Deps...)
	for _, m := range candidates {
		if m == nil || m.Path != modulePath {
			continue
		}
		if m.Replace != nil {
			break // a replace directive's version is a placeholder, not a release
		}
		v := strings.TrimPrefix(m.Version, "v")
		if v == "" || v == "(devel)" || strings.HasPrefix(v, "0.0.0") {
			break // source checkout / untagged pseudo-version
		}
		return v
	}
	return fallbackVersion
}

const (
	DefaultAPIVersion             = "v4"
	DefaultTimeoutSeconds         = 30.0
	DefaultConnectTimeoutSeconds  = 10.0
	DefaultMaxRetries             = 1
	DefaultRetryBackoffSeconds    = 0.5
	DefaultRetryBackoffMaxSeconds = 20.0
	DefaultPoolMaxsize            = 10
	DefaultLogLevel               = "WARNING"
	DefaultPrintFormat            = PrintFormatPretty

	// DefaultMaxResponseBytes caps how much of one response the SDK buffers in memory: 64 MiB.
	DefaultMaxResponseBytes = 64 << 20

	// PrintFormatPretty is the multi-line console print format.
	PrintFormatPretty = "pretty"
	// PrintFormatCompact is the single-line console print format.
	PrintFormatCompact = "compact"

	// ResponseFormatMap presents ApiResponse.Data as the decoded body — a map, a list or bytes.
	ResponseFormatMap = "map"
	// ResponseFormatJSON presents ApiResponse.Data as the server's own body text.
	ResponseFormatJSON = "json"
	// ResponseFormatDTO presents ApiResponse.Data as the endpoint's typed view.
	ResponseFormatDTO = "dto"

	DefaultResponseFormat = ResponseFormatMap

	// ResponseFormatEnvVar names the environment variable read when neither the call nor the config sets a response format.
	ResponseFormatEnvVar = "WC_RESPONSE_FORMAT"

	maxHTTPCodes = 500
)

var (
	defaultRetryableHTTPCodes = []int{429, 500, 502, 503, 504}
	// Empty: neither 429 nor 503 proves an unsafe request was refused rather than run.
	defaultUnsafeRetryableHTTPCodes = []int{}

	// Severity order, not alphabetical: the refusal message names them in this order.
	logLevelNames    = []string{"DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL"}
	printFormatNames = []string{PrintFormatPretty, PrintFormatCompact}
	apiVersionRe     = regexp.MustCompile(`(?i)^v?[1-9][0-9]*$`)

	// The three response formats, in the order every refusal names them.
	responseFormatNames = []string{ResponseFormatMap, ResponseFormatJSON, ResponseFormatDTO}
)

type SdkConfig struct {
	Auth AuthConfig

	// APIVersion is the API version segment, e.g. "v4".
	APIVersion string
	// TimeoutSeconds bounds one attempt, headers and body read alike, not the call as a whole.
	TimeoutSeconds        float64
	ConnectTimeoutSeconds float64
	// MaxRetries is total attempts including the first, so 1 disables retrying.
	MaxRetries int
	// RetryableHTTPCodes are the statuses retried for a safe method.
	RetryableHTTPCodes []int
	// UnsafeRetryableHTTPCodes are the statuses also retried for an unsafe method; a code here is retried only if it is in RetryableHTTPCodes as well.
	UnsafeRetryableHTTPCodes []int
	// RetryBackoffSeconds is the base backoff before the first retry; nil means 0.5, so Ptr(0.0) retries at once.
	RetryBackoffSeconds *float64
	// RetryBackoffMaxSeconds is the ceiling on any single backoff.
	RetryBackoffMaxSeconds float64

	// VerifySsl controls TLS certificate verification; nil means true, and Ptr(false) emits a warning.
	VerifySsl *bool
	// CABundle is a PEM file to trust instead of disabling verification; refused beside VerifySsl=false.
	CABundle string
	ProxyURL string
	// PoolMaxsize bounds idle connections per host.
	PoolMaxsize int
	// MaxResponseBytes caps how much of one response body the transport buffers; 0 takes the default.
	MaxResponseBytes int

	// DefaultHeaders are sent with every request; a Content-Type here labels every encoded body.
	DefaultHeaders map[string]string
	// RaiseOnError also returns an *SdkError for a failed response, judged by ApiResponse.Succeeded.
	RaiseOnError bool
	// ResponseFormat is what ApiResponse.Data follows: "map", "json" or "dto".
	ResponseFormat string
	UserAgent      string

	LogLevel  string
	LogBodies bool
	// LogFile also appends the client's log lines to this path; the handle is released by Close.
	LogFile string

	// The six print switches, off unless set, print the resolved URL, query values, body and a curl line.
	PrintApiInfo        bool
	PrintRequestPayload bool
	PrintResponseBody   bool
	PrintHeaders        bool
	PrintCurl           bool
	PrintOnlyOnError    bool
	// PrintFormat is "pretty" or "compact".
	PrintFormat string
	PrintColor  bool
	// PrintMaxChars truncates the rendered exchange block; 0 = unlimited.
	PrintMaxChars int
}

// VerifySslEnabled reads the tri-state VerifySsl with its default: nil means true.
func (c SdkConfig) VerifySslEnabled() bool { return c.VerifySsl == nil || *c.VerifySsl }

// ResponseFormatValue resolves the format: this config's, then WC_RESPONSE_FORMAT, then the default.
func (c SdkConfig) ResponseFormatValue() string {
	if c.ResponseFormat != "" {
		return c.ResponseFormat
	}
	if fromEnvironment, set := environmentResponseFormat(); set {
		return fromEnvironment
	}
	return DefaultResponseFormat
}

// environmentResponseFormat reads WC_RESPONSE_FORMAT trimmed and lower-cased; blank reads as unset.
func environmentResponseFormat() (string, bool) {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(ResponseFormatEnvVar)))
	return value, value != ""
}

// RetryBackoffSecondsValue reads RetryBackoffSeconds with its default: nil means 0.5.
func (c SdkConfig) RetryBackoffSecondsValue() float64 {
	if c.RetryBackoffSeconds == nil {
		return DefaultRetryBackoffSeconds
	}
	return *c.RetryBackoffSeconds
}

// String renders every field in declaration order.
func (c SdkConfig) String() string {
	return fmt.Sprintf(
		"SdkConfig(auth=%s, apiVersion=%q, timeoutSeconds=%v, connectTimeoutSeconds=%v, maxRetries=%d, "+
			"retryableHttpCodes=%v, unsafeRetryableHttpCodes=%v, retryBackoffSeconds=%v, "+
			"retryBackoffMaxSeconds=%v, verifySsl=%v, caBundle=%s, proxyUrl=%s, poolMaxsize=%d, "+
			"maxResponseBytes=%d, "+
			"defaultHeaders=%v, raiseOnError=%v, responseFormat=%q, userAgent=%s, logLevel=%q, logBodies=%v, logFile=%s, "+
			"printApiInfo=%v, printRequestPayload=%v, printResponseBody=%v, printHeaders=%v, printCurl=%v, "+
			"printOnlyOnError=%v, printFormat=%q, printColor=%v, printMaxChars=%d)",
		c.Auth, redact.RedactValue(c.APIVersion), c.TimeoutSeconds, c.ConnectTimeoutSeconds,
		c.MaxRetries,
		c.RetryableHTTPCodes, c.UnsafeRetryableHTTPCodes, c.RetryBackoffSecondsValue(),
		c.RetryBackoffMaxSeconds, c.VerifySslEnabled(),
		// Described, not masked: these field names are not secret names, so key masking would miss them.
		redact.DescribePath(c.CABundle), redact.Describe(c.ProxyURL), c.PoolMaxsize,
		c.MaxResponseBytes,
		redact.RedactMapping(c.DefaultHeaders), c.RaiseOnError,
		redact.RedactValue(c.ResponseFormatValue()), redact.Describe(c.UserAgent),
		redact.RedactValue(c.LogLevel), c.LogBodies,
		redact.DescribePath(c.LogFile),
		c.PrintApiInfo, c.PrintRequestPayload, c.PrintResponseBody, c.PrintHeaders, c.PrintCurl,
		c.PrintOnlyOnError, redact.RedactValue(c.PrintFormat), c.PrintColor, c.PrintMaxChars)
}

func (c SdkConfig) Normalised() SdkConfig {
	out := c
	out.Auth = c.Auth.Normalised()
	out.APIVersion = normaliseAPIVersion(c.APIVersion)
	if out.TimeoutSeconds == 0 {
		out.TimeoutSeconds = DefaultTimeoutSeconds
	}
	if out.ConnectTimeoutSeconds == 0 {
		out.ConnectTimeoutSeconds = DefaultConnectTimeoutSeconds
	}
	if out.MaxRetries == 0 {
		out.MaxRetries = DefaultMaxRetries
	}
	// Copied unconditionally, so a caller mutating what they passed cannot reach a built client.
	out.RetryableHTTPCodes = copyCodes(c.RetryableHTTPCodes, defaultRetryableHTTPCodes)
	out.UnsafeRetryableHTTPCodes = copyCodes(c.UnsafeRetryableHTTPCodes, defaultUnsafeRetryableHTTPCodes)
	out.DefaultHeaders = maps.Clone(c.DefaultHeaders)
	if out.RetryBackoffSeconds == nil {
		backoff := DefaultRetryBackoffSeconds
		out.RetryBackoffSeconds = &backoff
	} else {
		backoff := *c.RetryBackoffSeconds
		out.RetryBackoffSeconds = &backoff
	}
	if out.RetryBackoffMaxSeconds == 0 {
		out.RetryBackoffMaxSeconds = DefaultRetryBackoffMaxSeconds
	}
	// Trimmed here so validate and the transport agree: a blank proxyUrl would disable ProxyFromEnvironment.
	out.CABundle = normaliseOptionalPath(c.CABundle)
	out.ProxyURL = normaliseProxyURL(c.ProxyURL)
	if out.VerifySsl == nil {
		verify := true
		out.VerifySsl = &verify
	}
	if out.PoolMaxsize == 0 {
		out.PoolMaxsize = DefaultPoolMaxsize
	}
	if out.MaxResponseBytes == 0 {
		out.MaxResponseBytes = DefaultMaxResponseBytes
	}
	if out.LogLevel == "" {
		out.LogLevel = DefaultLogLevel
	} else {
		out.LogLevel = strings.ToUpper(strings.TrimSpace(c.LogLevel))
	}
	if out.PrintFormat == "" {
		out.PrintFormat = DefaultPrintFormat
	} else {
		out.PrintFormat = strings.ToLower(strings.TrimSpace(c.PrintFormat))
	}
	// Not defaulted, unlike printFormat: blank must survive for the environment and default to decide.
	out.ResponseFormat = strings.ToLower(strings.TrimSpace(c.ResponseFormat))
	if out.UserAgent == "" {
		out.UserAgent = DefaultUserAgent()
	}
	return out
}

// Validate reports every problem at once, judging the config as NewClient would see it.
func (c SdkConfig) Validate() error {
	return c.Normalised().validate()
}

// validate is the body of Validate; it requires an already-Normalised receiver.
func (c SdkConfig) validate() error {
	var problems []string
	// details carries the described values on CauseText, which Explain renders and Error does not.
	var details []string

	problems = append(problems, headerProblems(c.DefaultHeaders)...)

	if err := c.Auth.Validate(); err != nil {
		problems = append(problems, err.Error())
		var authError *wcerrors.SdkError
		if errors.As(err, &authError) && authError.CauseText != "" {
			details = append(details, authError.CauseText)
		}
	}
	if c.APIVersion == "" || !apiVersionRe.MatchString(c.APIVersion) {
		problems = append(problems, fmt.Sprintf("apiVersion: expected something like 'v4' or '4', got '%s'", c.APIVersion))
	}
	if c.MaxRetries < 1 {
		problems = append(problems, fmt.Sprintf("maxRetries: must be at least 1 (1 disables retry), got %d", c.MaxRetries))
	}
	if !finitePositive(c.TimeoutSeconds) {
		problems = append(problems, "timeoutSeconds: must be a finite number > 0")
	}
	if !finitePositive(c.ConnectTimeoutSeconds) {
		problems = append(problems, "connectTimeoutSeconds: must be a finite number > 0")
	}
	if IsNaNOrInf(c.RetryBackoffSecondsValue()) || c.RetryBackoffSecondsValue() < 0 {
		problems = append(problems, "retryBackoffSeconds: must be a finite number >= 0")
	}
	if IsNaNOrInf(c.RetryBackoffMaxSeconds) || c.RetryBackoffMaxSeconds < c.RetryBackoffSecondsValue() {
		problems = append(problems, "retryBackoffMaxSeconds: must be finite and >= retryBackoffSeconds")
	}
	if c.PoolMaxsize < 1 {
		problems = append(problems, "poolMaxsize: must be at least 1")
	}
	if c.MaxResponseBytes < 1 {
		// Only a negative value reaches here; 0 already became the default.
		problems = append(problems, fmt.Sprintf(
			"maxResponseBytes: must be at least 1 byte (0 takes the %d-byte default), got %d",
			DefaultMaxResponseBytes, c.MaxResponseBytes))
	}
	if c.ProxyURL != "" {
		// Refused at construction: a proxy that fails to parse and silently goes direct is worse.
		if p, err := url.Parse(c.ProxyURL); err != nil || p.Scheme == "" || p.Host == "" {
			problems = append(problems, "proxyUrl: is not a URL — expected something like "+
				"http://proxy.internal:8080 (scheme and host required)")
		}
	}
	problems = append(problems, HTTPCodeSetProblems("retryableHttpCodes", c.RetryableHTTPCodes)...)
	problems = append(problems, HTTPCodeSetProblems("unsafeRetryableHttpCodes", c.UnsafeRetryableHTTPCodes)...)
	if !slices.Contains(logLevelNames, c.LogLevel) {
		problems = append(problems, fmt.Sprintf("logLevel: expected one of %s, got '%s'",
			strings.Join(logLevelNames, ", "), c.LogLevel))
	}
	if !slices.Contains(printFormatNames, c.PrintFormat) {
		problems = append(problems, fmt.Sprintf("printFormat: expected one of %s, got '%s'",
			strings.Join(printFormatNames, ", "), c.PrintFormat))
	}
	validateResponseFormat(&problems, c.ResponseFormat)
	if c.PrintMaxChars < 0 {
		problems = append(problems, "printMaxChars: must be >= 0 (0 = unlimited)")
	}
	if strings.TrimSpace(c.UserAgent) == "" {
		problems = append(problems, "userAgent: must not be blank")
	}
	if c.CABundle != "" {
		if !c.VerifySslEnabled() {
			problems = append(problems, "caBundle and verifySsl=false are contradictory: one supplies a "+
				"trust store, the other disables verification")
		}
		if info, err := os.Stat(c.CABundle); err != nil || info.IsDir() {
			problems = append(problems, fmt.Sprintf(
				"caBundle: %s does not name a readable file", redact.DescribePath(c.CABundle)))
			details = append(details, "caBundle: "+c.CABundle)
		}
	}

	if len(problems) > 0 {
		refusal := wcerrors.NewConfigError(strings.Join(problems, "; "))
		refusal.CauseText = strings.Join(details, "; ")
		return refusal
	}
	return nil
}

func validateResponseFormat(problems *[]string, configured string) {
	if configured != "" {
		if !IsResponseFormat(configured) {
			*problems = append(*problems, fmt.Sprintf("responseFormat: expected one of %s, got '%s'",
				strings.Join(responseFormatNames, ", "), configured))
		}
		return
	}
	fromEnvironment, set := environmentResponseFormat()
	if set && !IsResponseFormat(fromEnvironment) {
		*problems = append(*problems, fmt.Sprintf("%s: expected one of %s, got '%s'",
			ResponseFormatEnvVar, strings.Join(responseFormatNames, ", "), fromEnvironment))
	}
}

func HTTPCodeSetProblems(name string, codes []int) []string {
	if len(codes) > maxHTTPCodes {
		return []string{fmt.Sprintf(
			"%s: at most %d status codes (100-599 is every code there is), got more than that",
			name, maxHTTPCodes)}
	}
	var problems []string
	for _, code := range codes {
		if code < 100 || code > 599 {
			problems = append(problems, fmt.Sprintf("%s: %d is not an HTTP status code", name, code))
		}
	}
	return problems
}

// headerProblems refuses a DefaultHeaders entry the transport could not send.
func headerProblems(headers map[string]string) []string {
	var problems []string
	for name, value := range headers {
		problem := HeaderNameProblem(name)
		if problem == "" {
			problem = HeaderValueProblem(value)
		}
		if problem != "" {
			problems = append(problems, fmt.Sprintf(
				"defaultHeaders: '%s' cannot be sent as a header — %s. The value is not shown: "+
					"it may be a credential", redact.NeutraliseControls(name), problem))
		}
	}
	sort.Strings(problems)
	return problems
}

// HeaderNameProblem names what makes a header name unsendable, or "" when it is sendable.
func HeaderNameProblem(name string) string {
	if name == "" {
		return "a blank name"
	}
	if !IsHTTPToken(name) {
		return "a name outside the RFC 9110 token set, which is letters, digits and " +
			"!#$%&'*+-.^_`|~ — a space, a colon and a non-ASCII byte are all refused there"
	}
	return ""
}

// HeaderValueProblem names what makes a header value unsendable, or "" when it is sendable.
func HeaderValueProblem(value string) string {
	if value == "" {
		return ""
	}
	if value[0] == ' ' || value[0] == '\t' {
		return "a value opening with a space or a tab, which several HTTP clients trim, so the " +
			"value on the wire would not be the one written"
	}
	for index := 0; index < len(value); index++ {
		byteValue := value[index]
		switch {
		case byteValue == '\t':
			continue // legal inside a field value; only a leading or trailing one is refused
		case byteValue == '\r':
			return fmt.Sprintf("a carriage return at byte %d, which would end the header line "+
				"and let whatever follows be read as a header of its own", index)
		case byteValue == '\n':
			return fmt.Sprintf("a newline at byte %d, which would end the header line and let "+
				"whatever follows be read as a header of its own", index)
		case byteValue == 0x00:
			return fmt.Sprintf("a NUL at byte %d", index)
		case byteValue < 0x20 || byteValue == 0x7f:
			return fmt.Sprintf("a control character at byte %d", index)
		case byteValue >= 0x80:
			return fmt.Sprintf("a non-ASCII byte at byte %d; a header value has no charset on "+
				"the wire, so encode it yourself — base64, or RFC 8187 percent-encoding", index)
		}
	}
	// Last, so a trailing CR or NUL is named as a smuggling primitive rather than as whitespace.
	if last := value[len(value)-1]; last == ' ' || last == '\t' {
		return "a value ending in a space or a tab, which several HTTP clients trim, so the " +
			"value on the wire would not be the one written"
	}
	return ""
}

// IsHTTPToken reports whether a value is an RFC 9110 token, the grammar of header names and methods.
func IsHTTPToken(value string) bool {
	for index := 0; index < len(value); index++ {
		c := value[index]
		allowed := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') ||
			strings.IndexByte("!#$%&'*+-.^_`|~", c) >= 0
		if !allowed {
			return false
		}
	}
	return true
}

func normaliseAPIVersion(raw string) string {
	if raw == "" {
		return DefaultAPIVersion
	}
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" {
		return ""
	}
	if !strings.HasPrefix(value, "v") {
		value = "v" + value
	}
	return value
}

// normaliseOptionalPath keeps a whitespace-only path verbatim, so it fails the file check.
func normaliseOptionalPath(raw string) string {
	if trimmed := strings.TrimSpace(raw); trimmed != "" || raw == "" {
		return trimmed
	}
	return raw
}

// DeepCopy returns a copy that shares no mutable state with the receiver.
func (c SdkConfig) DeepCopy() SdkConfig {
	out := c
	out.RetryableHTTPCodes = slices.Clone(c.RetryableHTTPCodes)
	out.UnsafeRetryableHTTPCodes = slices.Clone(c.UnsafeRetryableHTTPCodes)
	out.DefaultHeaders = maps.Clone(c.DefaultHeaders)
	if c.RetryBackoffSeconds != nil {
		backoff := *c.RetryBackoffSeconds
		out.RetryBackoffSeconds = &backoff
	}
	if c.VerifySsl != nil {
		verify := *c.VerifySsl
		out.VerifySsl = &verify
	}
	return out
}

// WithOverrides derives a new config, changing the snake_case-named fields in changes.
func (c SdkConfig) WithOverrides(changes map[string]any) (SdkConfig, error) {
	return applyOverrides(c, changes)
}

func SdkConfigFromCredentials(credentials []any, instance any, overrides map[string]any) (SdkConfig, error) {
	blocks := make([]any, len(credentials))
	copy(blocks, credentials)
	chosen, err := SelectCredentials(blocks, instance)
	if err != nil {
		return SdkConfig{}, err
	}

	switch typed := chosen.(type) {
	case SdkConfig:
		return applyOverrides(typed, overrides)
	case *SdkConfig:
		return applyOverrides(*typed, overrides)
	case AuthConfig:
		return applyOverrides(SdkConfig{Auth: typed}, overrides)
	case *AuthConfig:
		return applyOverrides(SdkConfig{Auth: *typed}, overrides)
	case map[string]any:
		return fromMapping(typed, instance, overrides)
	default:
		return SdkConfig{}, wcerrors.NewConfigError(fmt.Sprintf(
			"credentials[%s]: expected an SdkConfig, an AuthConfig or a mapping of config "+
				"fields, got %T", selectorLabel(instance), chosen))
	}
}

// selectorLabel names the selector for a message; an absent one reads as 0, the entry it chose.
func selectorLabel(instance any) string {
	if instance == nil {
		return "0"
	}
	if name, ok := instance.(string); ok {
		return fmt.Sprintf("%q", name)
	}
	return fmt.Sprintf("%v", instance)
}

func fromMapping(block map[string]any, instance any, overrides map[string]any) (SdkConfig, error) {
	remaining := make(map[string]any, len(block))
	maps.Copy(remaining, block)

	// "auth" is a field of the config but not of AuthConfig, so it comes out first.
	suppliedAuth, hasAuth := remaining["auth"]
	delete(remaining, "auth")

	authNames := authFieldNames()
	authValues := map[string]any{}
	for key := range remaining {
		if authNames[key] {
			authValues[key] = remaining[key]
			delete(remaining, key)
		}
	}

	settings := make(map[string]any, len(remaining)+len(overrides))
	maps.Copy(settings, remaining)
	maps.Copy(settings, overrides)

	if hasAuth && len(authValues) > 0 {
		named := slices.Sorted(maps.Keys(authValues))
		return SdkConfig{}, wcerrors.NewConfigError(
			"auth: an AuthConfig was supplied alongside the credential field(s) " +
				strings.Join(named, ", ") + "; pass one or the other")
	}

	var config SdkConfig
	if !hasAuth {
		// Checked here, since a missing required key would otherwise fail deeper, naming no field.
		if absent := absentRequiredAuthFields(authValues); len(absent) > 0 {
			return SdkConfig{}, wcerrors.NewConfigError(
				"credentials: missing required field(s): " + strings.Join(absent, ", "))
		}
		auth, err := authFromMapping(authValues)
		if err != nil {
			return SdkConfig{}, err
		}
		config.Auth = auth
	} else {
		switch typed := suppliedAuth.(type) {
		case AuthConfig:
			config.Auth = typed
		case *AuthConfig:
			config.Auth = *typed
		default:
			return SdkConfig{}, wcerrors.NewConfigError(fmt.Sprintf(
				"credentials[%s]: auth must be an AuthConfig, got %T",
				selectorLabel(instance), suppliedAuth))
		}
	}
	return applyOverrides(config, settings)
}

func authFieldNames() map[string]bool {
	names := map[string]bool{}
	authType := reflect.TypeFor[AuthConfig]()
	for index := 0; index < authType.NumField(); index++ {
		names[dto.WireName(authType.Field(index).Name)] = true
	}
	return names
}

func absentRequiredAuthFields(supplied map[string]any) []string {
	authType := reflect.TypeFor[AuthConfig]()
	var absent []string
	for index := 0; index < authType.NumField(); index++ {
		field := authType.Field(index)
		if field.Tag.Get("wc") != "required" {
			continue
		}
		name := dto.WireName(field.Name)
		if _, found := supplied[name]; !found {
			absent = append(absent, name)
		}
	}
	return absent
}

func authFromMapping(values map[string]any) (AuthConfig, error) {
	var auth AuthConfig
	if err := assignFields(reflect.ValueOf(&auth).Elem(), values, "credentials"); err != nil {
		return AuthConfig{}, err
	}
	return auth, nil
}

func applyOverrides(config SdkConfig, overrides map[string]any) (SdkConfig, error) {
	if len(overrides) > 0 {
		if err := assignFields(reflect.ValueOf(&config).Elem(), overrides, "overrides"); err != nil {
			return SdkConfig{}, err
		}
	}
	if err := config.Validate(); err != nil {
		return SdkConfig{}, err
	}
	return config.Normalised(), nil
}

func assignFields(target reflect.Value, values map[string]any, where string) error {
	targetType := target.Type()
	byName := map[string]int{}
	for index := 0; index < targetType.NumField(); index++ {
		byName[dto.WireName(targetType.Field(index).Name)] = index
	}
	unknown := make([]string, 0)
	for _, key := range slices.Sorted(maps.Keys(values)) {
		index, found := byName[key]
		if !found {
			unknown = append(unknown, key)
			continue
		}
		if err := assignOne(target.Field(index), values[key], key); err != nil {
			return err
		}
	}
	if len(unknown) > 0 {
		known := slices.Sorted(maps.Keys(byName))
		return wcerrors.NewConfigError("unknown config field(s): " + strings.Join(unknown, ", ") +
			" (in " + where + "); known: " + strings.Join(known, ", "))
	}
	return nil
}

// assignOne converts one loader-shaped value onto one field, refusing by name rather than panicking.
func assignOne(field reflect.Value, value any, name string) error {
	if value == nil {
		return nil // absent and explicitly null read alike: leave the zero value
	}
	incoming := reflect.ValueOf(value)
	if incoming.Type().AssignableTo(field.Type()) {
		field.Set(incoming)
		return nil
	}
	switch field.Kind() {
	case reflect.Float64:
		if number, ok := asFloat(value); ok {
			field.SetFloat(number)
			return nil
		}
	case reflect.Int:
		if number, ok := asFloat(value); ok && number == float64(int(number)) {
			field.SetInt(int64(number))
			return nil
		}
	case reflect.String:
		if text, ok := value.(string); ok {
			field.SetString(text)
			return nil
		}
		// A non-string scalar renders as text, so {"default_headers": {"X-Retry-Count": 3}} works.
		if text, ok := scalarText(value); ok {
			field.SetString(text)
			return nil
		}
	case reflect.Bool:
		if flag, ok := value.(bool); ok {
			field.SetBool(flag)
			return nil
		}
	case reflect.Pointer:
		// A loader supplies the bare value for a *bool / *float64 field; this wraps it.
		holder := reflect.New(field.Type().Elem())
		if err := assignOne(holder.Elem(), value, name); err != nil {
			return err
		}
		field.Set(holder)
		return nil
	case reflect.Slice:
		if items, ok := value.([]any); ok {
			out := reflect.MakeSlice(field.Type(), 0, len(items))
			for _, item := range items {
				element := reflect.New(field.Type().Elem()).Elem()
				if err := assignOne(element, item, name); err != nil {
					return err
				}
				out = reflect.Append(out, element)
			}
			field.Set(out)
			return nil
		}
	case reflect.Map:
		if entries, ok := value.(map[string]any); ok {
			out := reflect.MakeMapWithSize(field.Type(), len(entries))
			for key, entry := range entries {
				element := reflect.New(field.Type().Elem()).Elem()
				if err := assignOne(element, entry, name); err != nil {
					return err
				}
				out.SetMapIndex(reflect.ValueOf(key), element)
			}
			field.Set(out)
			return nil
		}
	}
	return wcerrors.NewConfigError(fmt.Sprintf(
		"%s: expected %s, got %T", name, field.Type(), value))
}

func scalarText(value any) (string, bool) {
	switch v := value.(type) {
	case bool:
		return strconv.FormatBool(v), true
	case int:
		return strconv.Itoa(v), true
	case int64:
		return strconv.FormatInt(v, 10), true
	case float64:
		return wcerrors.FormatSeconds(v), true
	case float32:
		return wcerrors.FormatSeconds(float64(v)), true
	}
	return "", false
}

func asFloat(value any) (float64, bool) {
	switch number := value.(type) {
	case float64:
		return number, true
	case float32:
		return float64(number), true
	case int:
		return float64(number), true
	case int64:
		return float64(number), true
	}
	return 0, false
}

// LogLevelRank orders the log levels by severity; an unknown level ranks below all of them.
func LogLevelRank(level string) int { return slices.Index(logLevelNames, level) }

func IsResponseFormat(value string) bool { return slices.Contains(responseFormatNames, value) }

// ResponseFormatNames returns a fresh copy of the three legal response formats, in refusal order.
func ResponseFormatNames() []string {
	return append(make([]string, 0, len(responseFormatNames)), responseFormatNames...)
}

func DefaultRetryableCodes() []int { return copyCodes(defaultRetryableHTTPCodes, nil) }

func DefaultUnsafeRetryableCodes() []int { return copyCodes(defaultUnsafeRetryableHTTPCodes, nil) }

// copyCodes substitutes fallback only for nil: an empty slice means retry nothing.
func copyCodes(codes, fallback []int) []int {
	if codes == nil {
		codes = fallback
	}
	return append(make([]int, 0, len(codes)), codes...)
}

// normaliseProxyURL gives a schemeless value an http:// scheme, so "proxy.internal:8080" works.
func normaliseProxyURL(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || strings.Contains(trimmed, "://") {
		return trimmed
	}
	if p, err := url.Parse(trimmed); err == nil && p.Host != "" {
		return trimmed
	}
	return "http://" + trimmed
}

func isBlank(value string) bool { return strings.TrimSpace(value) == "" }

func finitePositive(value float64) bool { return !IsNaNOrInf(value) && value > 0 }

func IsNaNOrInf(value float64) bool { return math.IsNaN(value) || math.IsInf(value, 0) }
