package http

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/redact"
)

var (
	placeholderRe     = regexp.MustCompile(`\{([a-z_][a-z0-9_]*)\}`)
	queryOrFragmentRe = regexp.MustCompile(`[?#]`)
	// The tail is \n?$: a URL copied from console output keeps its newline and must still be refused.
	redactedParameterRe = regexp.MustCompile(`[?#&]([^?#&=]*)=` + regexp.QuoteMeta(redact.Redacted) + `(?:[&#]|\n?$)`)
)

// FindPathParamNames returns the {name} placeholders, in first-seen order, query part excluded.
func FindPathParamNames(pathTemplate string) []string {
	template, _ := splitTemplate(pathTemplate)
	var names []string
	seen := map[string]bool{}
	for _, m := range placeholderRe.FindAllStringSubmatch(template, -1) {
		if !seen[m[1]] {
			seen[m[1]] = true
			names = append(names, m[1])
		}
	}
	return names
}

// SubstitutePathParams fills {name} placeholders; every problem is reported in one error.
func SubstitutePathParams(pathTemplate string, pathParams map[string]any) (string, error) {
	template, queryAndFragment := splitTemplate(pathTemplate)
	expected := FindPathParamNames(pathTemplate)
	if pathParams == nil {
		pathParams = map[string]any{}
	}

	var problems []string

	var missing []string
	for _, n := range expected {
		if _, ok := pathParams[n]; !ok {
			missing = append(missing, n)
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		problems = append(problems, "missing path params: "+strings.Join(missing, ", "))
	}

	expectedSet := map[string]bool{}
	for _, n := range expected {
		expectedSet[n] = true
	}
	var unknown []string
	for k := range pathParams {
		if !expectedSet[k] {
			unknown = append(unknown, k)
		}
	}
	sort.Strings(unknown)
	if len(unknown) > 0 {
		problems = append(problems, "unknown path params: "+strings.Join(unknown, ", "))
	}

	var blank, wrongTyped, dotOnly []string
	for _, name := range expected {
		value, ok := pathParams[name]
		if !ok {
			continue
		}
		if value == nil || isBlankStringValue(value) {
			blank = append(blank, name)
		} else if !usablePathParam(value) {
			wrongTyped = append(wrongTyped, fmt.Sprintf("%s (%T)", name, value))
		} else if isDotOnlySegment(pathParamText(value)) {
			// Judged beside the other argument problems, so it is refused before any encoding or token grant.
			dotOnly = append(dotOnly, name)
		}
	}
	if len(blank) > 0 {
		sort.Strings(blank)
		problems = append(problems, "blank path params: "+strings.Join(blank, ", "))
	}
	if len(wrongTyped) > 0 {
		sort.Strings(wrongTyped)
		problems = append(problems, "wrong-typed path params: "+strings.Join(wrongTyped, ", "))
	}
	if len(dotOnly) > 0 {
		sort.Strings(dotOnly)
		// The name only, never the value, and ASCII only: this text reaches a Windows console.
		problems = append(problems, "dot-only path params: "+strings.Join(dotOnly, ", "))
	}

	if len(problems) > 0 {
		return "", wcerrors.NewUsageError(fmt.Sprintf("%s -> %s", pathTemplate, strings.Join(problems, "; ")))
	}

	substituted := placeholderRe.ReplaceAllStringFunc(template, func(match string) string {
		name := placeholderRe.FindStringSubmatch(match)[1]
		return percentEncode(pathParamText(pathParams[name]))
	})
	if strings.ContainsAny(substituted, "{}") {
		return "", wcerrors.NewUsageError(fmt.Sprintf("%s -> malformed placeholder, unresolved: %s", pathTemplate, substituted))
	}
	return substituted + queryAndFragment, nil
}

// BuildURL joins a base URL and a path.
func BuildURL(baseURL, path string) (string, error) {
	if err := rejectARedactedValue(path); err != nil {
		return "", err
	}
	if err := rejectAFragment("path", path); err != nil {
		return "", err
	}
	lowered := strings.ToLower(path)
	if strings.HasPrefix(lowered, "http://") || strings.HasPrefix(lowered, "https://") {
		if err := requireSameOrigin(baseURL, path); err != nil {
			return "", err
		}
		return path, nil
	}
	trimmedBase := strings.TrimRight(baseURL, "/")
	if strings.HasPrefix(path, "/") {
		return trimmedBase + path, nil
	}
	return trimmedBase + "/" + path, nil
}

// percentEncode applies RFC 3986 percent-encoding with an empty safe set: unreserved only.
func percentEncode(text string) string {
	var b strings.Builder
	for _, by := range []byte(text) {
		unreserved := (by >= 0x41 && by <= 0x5a) ||
			(by >= 0x61 && by <= 0x7a) ||
			(by >= 0x30 && by <= 0x39) ||
			by == 0x2d || by == 0x2e || by == 0x5f || by == 0x7e
		if unreserved {
			b.WriteByte(by)
		} else {
			fmt.Fprintf(&b, "%%%02X", by)
		}
	}
	return b.String()
}

func rejectARedactedValue(path string) error {
	if m := redactedParameterRe.FindStringSubmatch(path); m != nil {
		return wcerrors.NewUsageError(fmt.Sprintf(
			"path: parameter '%s' holds the redaction marker %s", m[1], redact.Redacted))
	}
	return nil
}

// rejectAFragment refuses a '#' anywhere in a path or an absolute URL.
func rejectAFragment(what, path string) error {
	at := strings.IndexByte(path, '#')
	if at < 0 {
		return nil
	}
	return wcerrors.NewUsageError(fmt.Sprintf(
		"%s: carries a '#' at byte %d; pass the query as queryParams", what, at))
}

func requireSameOrigin(baseURL, rawURL string) error {
	parsed, err := parseAbsolute(rawURL)
	if err != nil {
		return err
	}
	if parsed.User != nil {
		// A non-empty username or password, not a non-nil User: some proxies emit an empty userinfo.
		password, _ := parsed.User.Password()
		if parsed.User.Username() != "" || password != "" {
			return wcerrors.NewUsageError(fmt.Sprintf(
				"absolute URL for host %s carries embedded credentials, which would replace the bearer "+
					"token; strip the userinfo before following it", parsed.Hostname()))
		}
	}
	base, err := parseAbsolute(baseURL)
	if err != nil {
		return err
	}
	if OriginOf(parsed) != OriginOf(base) {
		// Neither URL is reproduced; which component differs is what a caller acts on.
		return wcerrors.NewUsageError(fmt.Sprintf(
			"absolute URL: it is not on the configured origin — %s; an absolute URL is requested "+
				"with the access token attached, so only the store's own links may be passed back",
			originDifference(parsed, base)))
	}
	return nil
}

// originDifference names which parts of an origin differ, so a refusal need not echo either URL.
func originDifference(parsed, base *url.URL) string {
	var parts []string
	if !strings.EqualFold(parsed.Scheme, base.Scheme) {
		parts = append(parts, "scheme")
	}
	if !strings.EqualFold(parsed.Hostname(), base.Hostname()) {
		parts = append(parts, "host")
	}
	if originPort(parsed) != originPort(base) {
		parts = append(parts, "port")
	}
	switch len(parts) {
	case 0:
		return "its origin differs"
	case 1:
		return "its " + parts[0] + " differs"
	default:
		return "its " + strings.Join(parts[:len(parts)-1], ", ") + " and " + parts[len(parts)-1] +
			" differ"
	}
}

func parseAbsolute(rawURL string) (*url.URL, error) {
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		// Described, not echoed: an unparseable URL may carry a password that key masking cannot reach.
		refusal := wcerrors.NewUsageError(fmt.Sprintf(
			"%s is not a parseable absolute URL", redact.Describe(rawURL)))
		refusal.Cause = err
		refusal.CauseText = redact.RedactText(rawURL)
		return nil, refusal
	}
	return u, nil
}

// OriginOf returns scheme://host:port with the scheme's implied port made explicit.
func OriginOf(u *url.URL) string {
	return fmt.Sprintf("%s://%s:%s",
		strings.ToLower(u.Scheme), strings.ToLower(u.Hostname()), originPort(u))
}

// originPort makes the scheme's implied port explicit, so origins compare alike either way.
func originPort(u *url.URL) string {
	if port := u.Port(); port != "" {
		return port
	}
	switch strings.ToLower(u.Scheme) {
	case "https":
		return "443"
	case "http":
		return "80"
	}
	return "-1"
}

func splitTemplate(pathTemplate string) (template, queryAndFragment string) {
	location := queryOrFragmentRe.FindStringIndex(pathTemplate)
	if location == nil {
		return pathTemplate, ""
	}
	return pathTemplate[:location[0]], pathTemplate[location[0]:]
}

// usablePathParam accepts a string, an integer or an integral json.Number, as responses decode numbers.
func usablePathParam(value any) bool {
	switch v := value.(type) {
	case string, int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return true
	case json.Number:
		_, err := v.Int64()
		return err == nil
	}
	return false
}

func isBlankStringValue(value any) bool {
	s, ok := value.(string)
	return ok && strings.TrimSpace(s) == ""
}

// isDotOnlySegment reports whether a path parameter is nothing but dots.
func isDotOnlySegment(text string) bool {
	return text != "" && strings.Trim(text, ".") == ""
}

func pathParamText(value any) string {
	if s, ok := value.(string); ok {
		return strings.TrimSpace(s)
	}
	return fmt.Sprint(value)
}
