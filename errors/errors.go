package errors

import (
	"errors"
	"fmt"
	"maps"
	"math"
	"slices"
	"strconv"
	"strings"
)

// errorKind classifies an SdkError for errors.Is, via the sentinels below.
type errorKind int

const (
	kindResponse  errorKind = iota // a non-2xx HTTP response
	kindConfig                     // configuration is missing or invalid
	kindUsage                      // a caller mistake caught before any I/O
	kindAttribute                  // a name that does not exist on a client or resource
	kindAuth                       // authentication failed
)

// Sentinels for errors.Is. Example: errors.Is(err, ErrUsage).
var (
	// ErrConfig matches configuration errors raised at construction.
	ErrConfig = errors.New("webcommander: configuration error")
	// ErrUsage matches caller mistakes caught before any network I/O.
	ErrUsage = errors.New("webcommander: usage error")
	// ErrAttribute matches lookups of a name that does not exist on a client or a bound resource.
	ErrAttribute = errors.New("webcommander: attribute error")
	ErrAuth      = errors.New("webcommander: authentication error")
)

// ErrorType classifies a failed call, from the HTTP status plus an isSuccess denial under a 2xx.
type ErrorType string

// Callers compare against these values and serialised exchanges store them, so do not respell them.
const (
	ErrorTypeNone             ErrorType = "none"               // 2xx — no error
	ErrorTypeRedirect         ErrorType = "redirect"           // 3xx — an unfollowed redirect
	ErrorTypeValidation       ErrorType = "validation"         // 400/422 — malformed or invalid
	ErrorTypeAuthentication   ErrorType = "authentication"     // 401 — missing/rejected/expired
	ErrorTypeAuthorization    ErrorType = "authorization"      // 403 — not permitted
	ErrorTypeNotFound         ErrorType = "not_found"          // 404 — no such record
	ErrorTypeMethodNotAllowed ErrorType = "method_not_allowed" // 405 — wrong method
	ErrorTypeConflict         ErrorType = "conflict"           // 409 — conflict with existing data
	ErrorTypeRateLimit        ErrorType = "rate_limit"         // 429 — rate limited
	ErrorTypeServer           ErrorType = "server"             // 5xx — the server failed
	ErrorTypeNetwork          ErrorType = "network"            // transport failure, no HTTP status
	ErrorTypeTimeout          ErrorType = "timeout"            // the request timed out
	ErrorTypeUnknown          ErrorType = "unknown"            // any other non-2xx status
)

// RequestContext is the request record an SdkError carries.
type RequestContext interface {
	RequestSummary() (queryParams map[string]any, attemptCount, maxAttempts int,
		elapsedSeconds float64, httpMethod, url string)
}

// EndpointNamer is the optional half of RequestContext: it names the endpoint template, not the path.
type EndpointNamer interface {
	// RequestEndpoint is the endpoint template, or "" when the sender supplied none.
	RequestEndpoint() string
}

// SdkError is the one error type the SDK returns.
type SdkError struct {
	Message string
	// StatusCode is the HTTP status when the error came from a response; 0 otherwise.
	StatusCode int
	// ErrorType is the failure classification; the zero value "" means unclassified.
	ErrorType ErrorType
	// RawResponse is the decoded response body, already redacted.
	RawResponse any
	// ResponseHeaders are the response headers, already redacted; empty if no response.
	ResponseHeaders map[string]string
	// RequestInfo describes the call that failed; nil when the failure preceded a built request.
	RequestInfo RequestContext
	Cause       error
	// CauseText is the underlying failure's own words, already redacted.
	CauseText string

	kind        errorKind
	fieldErrors []string
}

// Error renders one structural line: the status, the call by its endpoint template, and the summary.
func (sdkError *SdkError) Error() string {
	var parts []string
	if sdkError.StatusCode != 0 {
		parts = append(parts, fmt.Sprintf("%d", sdkError.StatusCode))
	}
	if info := sdkError.RequestInfo; info != nil {
		_, _, _, _, method, url := info.RequestSummary()
		endpointTemplate := ""
		if named, isNamed := info.(EndpointNamer); isNamed {
			endpointTemplate = named.RequestEndpoint()
		}
		if label := EndpointLabel(endpointTemplate, url); method != "" && label != "" {
			parts = append(parts, flattenControl(method)+" "+flattenControl(label))
		}
	}
	reason := flattenControl(sdkError.Message)
	if len(parts) == 0 {
		return reason
	}
	return strings.Join(parts, " ") + " - " + reason
}

// EndpointLabel names a call by its endpoint template, else by the resolved path's segment count.
func EndpointLabel(endpointTemplate, resolvedURL string) string {
	if endpointTemplate != "" {
		return endpointTemplate
	}
	return pathShape(resolvedURL)
}

// pathShape describes a resolved path by how many segments it has — "<4-segment path>".
func pathShape(resolvedURL string) string {
	path := shortenURL(resolvedURL)
	if path == "" {
		return ""
	}
	// The query and any fragment are values, not shape.
	if cut := strings.IndexAny(path, "?#"); cut >= 0 {
		path = path[:cut]
	}
	segments := 0
	for _, segment := range strings.Split(path, "/") {
		if segment != "" {
			segments++
		}
	}
	return fmt.Sprintf("<%d-segment path>", segments)
}

// shortenURL drops the scheme and host, keeping the path that identifies the call.
func shortenURL(url string) string {
	_, remainder, found := strings.Cut(url, "://")
	if !found {
		return url
	}
	path := strings.Index(remainder, "/")
	if path < 0 {
		return url
	}
	return remainder[path:]
}

func (sdkError *SdkError) Unwrap() error { return sdkError.Cause }

func (sdkError *SdkError) Is(target error) bool {
	switch sdkError.kind {
	case kindConfig:
		return target == ErrConfig
	case kindUsage:
		return target == ErrUsage
	case kindAttribute:
		// An attribute error is also a usage error.
		return target == ErrAttribute || target == ErrUsage
	case kindAuth:
		return target == ErrAuth
	}
	return false
}

func (sdkError *SdkError) GetErrorMessages() []string {
	out := make([]string, len(sdkError.fieldErrors))
	copy(out, sdkError.fieldErrors)
	return out
}

// GetHint returns a suggestion for what to do about this failure, or "".
func (sdkError *SdkError) GetHint() string {
	return hintsByStatus[sdkError.StatusCode]
}

// Explain returns the full diagnostic, the only render of what Error withholds: path, query, cause, body.
func (sdkError *SdkError) Explain() string {
	var b strings.Builder
	b.WriteString(sdkError.Error())

	// Each row is gated on its own value, so a record carrying only parameters still prints.
	if info := sdkError.RequestInfo; info != nil {
		params, attempt, maxAttempts, elapsed, _, url := info.RequestSummary()
		if url != "" {
			fmt.Fprintf(&b, "\n  path    : %s", flattenControl(shortenURL(url)))
		}
		if len(params) > 0 {
			fmt.Fprintf(&b, "\n  params  : %s", flattenControl(params))
		}
		if attempt > 0 && maxAttempts > 0 {
			fmt.Fprintf(&b, "\n  attempt : %d/%d", attempt, maxAttempts)
		}
		fmt.Fprintf(&b, "\n  elapsed : %ss", FormatSeconds(elapsed))
	}

	// Read from CauseText, not Cause, so redaction stays with the layer that knows the failure.
	if sdkError.CauseText != "" {
		fmt.Fprintf(&b, "\n  cause   : %s", flattenControl(sdkError.CauseText))
	}

	if location := redirectTarget(sdkError.ResponseHeaders); location != "" {
		fmt.Fprintf(&b, "\n  location: %s", flattenControl(location))
	}

	reported := sdkError.GetErrorMessages()
	if len(reported) > 1 || (len(reported) == 1 && !strings.Contains(sdkError.Error(), reported[0])) {
		b.WriteString("\n  reported:")
		for _, reportedMessage := range reported {
			fmt.Fprintf(&b, "\n    - %s", flattenControl(reportedMessage))
		}
	}

	if hint := sdkError.GetHint(); hint != "" {
		fmt.Fprintf(&b, "\n  hint    : %s", hint)
	}

	if sdkError.RawResponse != nil {
		body := flattenControl(sdkError.RawResponse)
		// Cut in runes: a byte-wise cut would split a rune and emit invalid UTF-8.
		if runes := []rune(body); len(runes) > 500 {
			body = fmt.Sprintf("%s... +%d chars", string(runes[:500]), len(runes)-500)
		}
		fmt.Fprintf(&b, "\n  body    : %s", body)
	}
	return b.String()
}

func NewConfigError(message string) *SdkError {
	return &SdkError{Message: message, kind: kindConfig}
}

func NewUsageError(message string) *SdkError {
	return &SdkError{Message: message, kind: kindUsage}
}

// NewAttributeError builds the error for an unknown name; it matches ErrAttribute and ErrUsage.
func NewAttributeError(message string) *SdkError {
	return &SdkError{Message: message, kind: kindAttribute}
}

// InsecureTransportWarning says certificate verification is off, or the base URL is plain http://.
type InsecureTransportWarning struct {
	Message string
}

func NewInsecureTransportWarning(message string) *InsecureTransportWarning {
	return &InsecureTransportWarning{Message: message}
}

// Error renders the warning's message with control characters escaped.
func (w *InsecureTransportWarning) Error() string { return flattenControl(w.Message) }

func redirectTarget(headers map[string]string) string {
	for k, v := range headers {
		if strings.ToLower(k) == "location" {
			return v
		}
	}
	return ""
}

// flattenControl renders one server-controlled value as a single escaped line, so a message cannot forge extra log lines.
func flattenControl(value any) string {
	var text string
	switch {
	case value == nil:
		text = "null"
	case isLiteralContainer(value):
		text = renderLiteral(value)
	default:
		text = fmt.Sprint(value)
	}
	return EscapeControls(text, false)
}

// EscapeControls rewrites every control character as one escape spelling: \n by name, else \xNN or \uNNNN.
func EscapeControls(text string, keepNewlines bool) string {
	var b strings.Builder
	for _, r := range text {
		if keepNewlines && r == '\n' {
			b.WriteRune(r)
			continue
		}
		switch {
		case r == '\n':
			b.WriteString(`\n`)
		case r == '\r':
			b.WriteString(`\r`)
		case r == '\t':
			b.WriteString(`\t`)
		case r <= 0x1f, r >= 0x7f && r <= 0x9f:
			fmt.Fprintf(&b, `\x%02x`, r)
		case r == 0x2028, r == 0x2029:
			fmt.Fprintf(&b, `\u%04x`, r)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func isLiteralContainer(value any) bool {
	switch value.(type) {
	case map[string]any, []any, map[string]string:
		return true
	}
	return false
}

// renderLiteral renders a decoded-JSON value in literal notation.
func renderLiteral(value any) string {
	switch v := value.(type) {
	case nil:
		return "None"
	case bool:
		if v {
			return "True"
		}
		return "False"
	case string:
		return renderLiteralString(v)
	case float64:
		return FormatSeconds(v)
	case map[string]any:
		return renderLiteralMap(v)
	case map[string]string:
		return renderLiteralMap(v)
	case []any:
		var b strings.Builder
		b.WriteByte('[')
		for i, entry := range v {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(renderLiteral(entry))
		}
		b.WriteByte(']')
		return b.String()
	default:
		return fmt.Sprint(v) // json.Number and the integer kinds print their literal
	}
}

func renderLiteralMap[V any](v map[string]V) string {
	var b strings.Builder
	b.WriteByte('{')
	for i, key := range slices.Sorted(maps.Keys(v)) {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(renderLiteralString(key))
		b.WriteString(": ")
		b.WriteString(renderLiteral(v[key]))
	}
	b.WriteByte('}')
	return b.String()
}

// renderLiteralString single-quotes, or double-quotes text carrying a single quote and no double one.
func renderLiteralString(s string) string {
	quote := "'"
	if strings.Contains(s, "'") && !strings.Contains(s, `"`) {
		quote = `"`
	}
	escaped := strings.ReplaceAll(s, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, quote, `\`+quote)
	return quote + escaped + quote
}

const redirectHint = "the SDK does not follow redirects; read the Location header"

var hintsByStatus = map[int]string{
	301: redirectHint,
	302: redirectHint,
	303: redirectHint,
	307: redirectHint,
	308: redirectHint,
	400: "the API rejected the request itself - check the payload shape and the " +
		"query parameter names against the endpoint documentation",
	401: "credentials were rejected, or the token expired and could not be renewed " +
		"- check clientId, clientSecret and authString",
	403: "authenticated, but this client is not allowed to do that",
	404: "no such record - check the id you passed actually exists",
	405: "wrong HTTP method for this path",
	409: "conflicts with existing data, e.g. an email already in use",
	429: "rate limited - slow down, or raise maxRetries so the SDK waits and retries for you",
	500: "the server failed. Retrying may help on a read; the SDK will not retry a " +
		"write on a 500, because it cannot tell whether the server committed",
}

// NewResponseError builds an error for a non-2xx response; a 401/403 is flagged as an auth error.
func NewResponseError(message string, statusCode int, errorType ErrorType, rawResponse any,
	responseHeaders map[string]string, requestInfo RequestContext, messages []string) *SdkError {
	kind := kindResponse
	if statusCode == 401 || statusCode == 403 {
		kind = kindAuth
	}
	sdkError := &SdkError{Message: message, StatusCode: statusCode, ErrorType: errorType, RawResponse: rawResponse,
		ResponseHeaders: copyHeaders(responseHeaders), RequestInfo: requestInfo, kind: kind}
	sdkError.fieldErrors = append(sdkError.fieldErrors, messages...)
	return sdkError
}

// copyHeaders snapshots the caller's map, so reusing it cannot change an earlier error's Explain.
func copyHeaders(responseHeaders map[string]string) map[string]string {
	copied := make(map[string]string, len(responseHeaders))
	maps.Copy(copied, responseHeaders)
	return copied
}

func NewAuthError(message string) *SdkError {
	return &SdkError{Message: message, kind: kindAuth}
}

func NewAuthResponseError(message string, statusCode int, errorType ErrorType, rawResponse any,
	responseHeaders map[string]string, requestInfo RequestContext, messages []string) *SdkError {
	sdkError := NewResponseError(message, statusCode, errorType, rawResponse, responseHeaders, requestInfo, messages)
	sdkError.kind = kindAuth
	return sdkError
}

// RoundSeconds rounds a duration to four decimal places where it is measured, so every reader agrees.
func RoundSeconds(value float64) float64 {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return value
	}
	return math.Round(value*1e4) / 1e4
}

// FormatSeconds renders a rounded duration as the shortest round-tripping decimal, always fractional.
func FormatSeconds(value float64) string {
	if math.IsNaN(value) {
		return "nan"
	}
	if math.IsInf(value, 1) {
		return "inf"
	}
	if math.IsInf(value, -1) {
		return "-inf"
	}

	decimalPoint := decimalExponent(value) + 1
	if decimalPoint < -3 || decimalPoint > 16 {
		return strconv.FormatFloat(value, 'e', -1, 64)
	}
	return withFraction(strconv.FormatFloat(value, 'f', -1, 64))
}

// decimalExponent is the power of ten in value's shortest round-tripping form: the "17" of 1e+17.
func decimalExponent(value float64) int {
	_, exponentText, found := strings.Cut(strconv.FormatFloat(value, 'e', -1, 64), "e")
	if !found {
		return 0
	}
	parsed, err := strconv.Atoi(exponentText)
	if err != nil {
		return 0
	}
	return parsed
}

// withFraction forces a fractional part on, so a two-second call prints "2.0s".
func withFraction(text string) string {
	if strings.ContainsRune(text, '.') {
		return text
	}
	return text + ".0"
}
