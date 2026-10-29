// Package redact decides what counts as a secret, for every path that renders SDK output.
package redact

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	Redacted = "***"
	// TooDeep replaces a structure nested deeper than MaxRedactionDepth.
	TooDeep = "<...>"
	// Circular replaces a container that contains itself.
	Circular = "<circular>"
	// Malformed opens the placeholder for text that opens as JSON but does not parse: its size, no bytes.
	Malformed = "<malformed json of "
	// MaxRedactionDepth is how deep RedactValue walks before writing TooDeep.
	MaxRedactionDepth = 64
)

var markers = []string{
	"secret", "password", "passphrase", "passwd", "pwd", "token", "bearer", "jwt",
	"authorization", "authstring", "apikey", "assertion", "credential", "verifier",
	"cookie", "session", "signature", "privatekey",
}

var nonSecretTail = map[string]bool{
	"amount": true, "at": true, "consent": true, "count": true, "duration": true,
	"level": true, "path": true, "policy": true, "remaining": true, "status": true,
	"type":    true,
	"expired": true, // customer_password_expired describes a password's state, never the password
}

var nonSecretLead = map[string]bool{"has": true, "is": true, "was": true}

var (
	camelBoundaryRe = regexp.MustCompile(`([a-z0-9])([A-Z])`)
	nonAlnumRe      = regexp.MustCompile(`[^A-Za-z0-9]+`)
	urlRunRe        = regexp.MustCompile(`[A-Za-z][A-Za-z0-9+.\-]{0,31}://[^\s<>"']*`)
	pathQueryRunRe  = regexp.MustCompile(`/[^\s<>"'?#]*[?#][^\s<>"']+`)
	// A scheme-relative authority carrying a password, "//user:pass@host".
	authorityRunRe = regexp.MustCompile(`//[^\s<>"'/?#@]*:[^\s<>"'/?#@:]*@[^\s<>"'/?#]*`)
	// The same without leading slashes, "user:pass@host:port".
	bareAuthorityRunRe = regexp.MustCompile(`[^\s<>"'/?#@:]+:[^\s<>"'/?#@:]*@[^\s<>"'/?#@]*[.:][^\s<>"'/?#@]*`)
)

// SensitiveKeyMarkers returns a fresh copy of the substrings that make a key name secret.
func SensitiveKeyMarkers() []string {
	return append(make([]string, 0, len(markers)), markers...)
}

func ShouldRedact(key string) bool {
	words := keyWords(key)
	if len(words) > 0 && (nonSecretLead[words[0]] || nonSecretTail[words[len(words)-1]]) {
		return false
	}
	normalised := strings.Join(words, "")
	for _, marker := range markers {
		if strings.Contains(normalised, marker) {
			return true
		}
	}
	return false
}

func keyWords(key string) []string {
	spaced := camelBoundaryRe.ReplaceAllString(key, "$1 $2")
	var words []string
	for _, w := range nonAlnumRe.Split(spaced, -1) {
		if w != "" {
			words = append(words, strings.ToLower(w))
		}
	}
	return words
}

// RedactValue deep-copies a value, masking secret-named entries; an over-deep or cyclic one degrades to a marker.
func RedactValue(value any) any {
	return redact(value, 0, map[uintptr]bool{})
}

func redact(value any, depth int, seen map[uintptr]bool) any {
	if s, ok := value.(string); ok {
		// In the order a value nests them: JSON that arrived as text, then any URL in the result.
		return RedactText(redactJSONText(s, depth, seen))
	}
	if b, ok := value.([]byte); ok {
		// Copied though it has no keys: sharing the array would leave the record's bytes writable by the caller.
		return append([]byte(nil), b...)
	}
	reflectValue := reflect.ValueOf(value)
	if !reflectValue.IsValid() {
		return value
	}
	switch reflectValue.Kind() {
	case reflect.Map:
		if depth >= MaxRedactionDepth {
			return TooDeep
		}
		pointer := reflectValue.Pointer()
		if seen[pointer] {
			return Circular
		}
		seen[pointer] = true
		defer delete(seen, pointer)
		out := make(map[string]any, reflectValue.Len())
		for _, k := range reflectValue.MapKeys() {
			key := mapKeyName(k)
			if ShouldRedact(key) {
				out[key] = Redacted
			} else {
				out[key] = redact(reflectValue.MapIndex(k).Interface(), depth+1, seen)
			}
		}
		return out
	case reflect.Slice, reflect.Array:
		if depth >= MaxRedactionDepth {
			return TooDeep
		}
		if reflectValue.Kind() == reflect.Slice {
			pointer := reflectValue.Pointer()
			if seen[pointer] {
				return Circular
			}
			seen[pointer] = true
			defer delete(seen, pointer)
		}
		out := make([]any, 0, reflectValue.Len())
		for i := 0; i < reflectValue.Len(); i++ {
			out = append(out, redactEntry(reflectValue.Index(i).Interface(), depth+1, seen))
		}
		return out
	default:
		return keepOrSummarise(value, reflectValue)
	}
}

// mapKeyName names a map key without calling its own String or Error, which could disguise a secret.
func mapKeyName(key reflect.Value) string {
	if !key.IsValid() {
		return typeSummary(key)
	}
	value := key.Interface()
	if _, isJSONNumber := value.(json.Number); isJSONNumber {
		return key.String()
	}
	if rendersItself(value) {
		return typeSummary(key)
	}
	switch key.Kind() {
	case reflect.String:
		return key.String()
	case reflect.Bool:
		return strconv.FormatBool(key.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(key.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(key.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return canonicalFloatRepr(key.Float())
	}
	return typeSummary(key)
}

// Describe names one configured value by its kind and length and by no part of the value itself.
func Describe(value string) string {
	if value == "" {
		return "<unset>"
	}
	if scheme, _, found := strings.Cut(value, "://"); found && isURLScheme(scheme) {
		return fmt.Sprintf("<an %s URL of %d characters>", strings.ToLower(scheme),
			len([]rune(value)))
	}
	return fmt.Sprintf("<a value of %d characters>", len([]rune(value)))
}

// DescribePath names a configured path by its segment count and last extension, never a segment.
func DescribePath(path string) string {
	if path == "" {
		return "<unset>"
	}
	cleaned := strings.ReplaceAll(path, "\\", "/")
	segments := 0
	last := ""
	for _, segment := range strings.Split(cleaned, "/") {
		if segment != "" {
			segments++
			last = segment
		}
	}
	if dot := strings.LastIndexByte(last, '.'); dot > 0 && dot < len(last)-1 {
		return fmt.Sprintf("<a %d-segment path ending %s>", segments, last[dot:])
	}
	return fmt.Sprintf("<a %d-segment path>", segments)
}

// isURLScheme lets Describe tell "https://host" from a value that merely holds "://".
func isURLScheme(text string) bool {
	if text == "" {
		return false
	}
	for index := 0; index < len(text); index++ {
		c := text[index]
		letter := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
		if index == 0 && !letter {
			return false
		}
		if !letter && !(c >= '0' && c <= '9') && c != '+' && c != '-' && c != '.' {
			return false
		}
	}
	return true
}

// keepOrSummarise keeps a JSON scalar and describes anything else by type and length, calling nothing on it.
func keepOrSummarise(value any, reflectValue reflect.Value) any {
	switch reflectValue.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.String:
		if _, isJSONNumber := value.(json.Number); isJSONNumber || !rendersItself(value) {
			return value
		}
	}
	return typeSummary(reflectValue)
}

// rendersItself reports whether printing a value would call a method the value supplied.
func rendersItself(value any) bool {
	switch value.(type) {
	case fmt.Stringer, error:
		return true
	}
	return false
}

// typeSummary names a value by its type, and by its length where it has one, showing none of it.
func typeSummary(reflectValue reflect.Value) string {
	switch reflectValue.Kind() {
	case reflect.Chan, reflect.String:
		return fmt.Sprintf("<%s of %d>", reflectValue.Type(), reflectValue.Len())
	}
	return fmt.Sprintf("<%s>", reflectValue.Type())
}

// redactEntry honours a [name, value] pair — header lists and form bodies carry secrets as pairs.
func redactEntry(entry any, depth int, seen map[uintptr]bool) any {
	if name, isPair := pairName(entry); isPair && ShouldRedact(name) {
		return []any{name, Redacted}
	}
	return redact(entry, depth, seen)
}

// pairName reports the head of a two-element pair whatever its concrete type.
func pairName(entry any) (string, bool) {
	if list, ok := entry.([]any); ok { // the common decoded-JSON shape, no reflection
		if len(list) != 2 {
			return "", false
		}
		name, ok := list[0].(string)
		return name, ok
	}
	value := reflect.ValueOf(entry)
	if (value.Kind() != reflect.Slice && value.Kind() != reflect.Array) || value.Len() != 2 {
		return "", false
	}
	name, ok := value.Index(0).Interface().(string)
	return name, ok
}

func redactJSONText(text string, depth int, seen map[uintptr]bool) string {
	trimmed := strings.TrimSpace(text)
	if !strings.HasPrefix(trimmed, "{") && !strings.HasPrefix(trimmed, "[") {
		return text
	}
	// NaN and Infinity are folded to same-length placeholders first: encoding/json refuses them.
	source, nonFinites := foldNonFiniteLiterals(text)
	decoder := &orderedDecoder{dec: json.NewDecoder(strings.NewReader(source)), nonFinite: nonFinites}
	decoder.dec.UseNumber()
	parsed, err := decoder.value()
	if err != nil {
		return withheldStructuredText(text)
	}
	// Trailing content is a second document whose keys were never judged.
	if _, err := decoder.dec.Token(); !errors.Is(err, io.EOF) {
		return withheldStructuredText(text)
	}
	redacted := redactOrdered(parsed, depth+1, seen)
	if reflect.DeepEqual(redacted, parsed) {
		return text
	}
	var out strings.Builder
	writeCanonical(&out, redacted)
	return out.String()
}

// withheldStructuredText stands in for text that opens as JSON but will not parse: kind and size only.
func withheldStructuredText(text string) string {
	return fmt.Sprintf("%s%d bytes>", Malformed, len(text))
}

// jsonMember is one key/value entry of a decoded JSON object, order preserved.
type jsonMember struct {
	key   string
	value any
}

type jsonObject []jsonMember

// nonFinite is a decoded NaN or Infinity leaf, typed apart so the re-emit is the source literal.
type nonFinite string

// Each placeholder is a JSON number of the literal's exact length, so every byte offset still matches.
var nonFinitePlaceholders = []struct{ literal, placeholder string }{
	{"-Infinity", "-0.000000"},
	{"Infinity", "0.000000"},
	{"NaN", "0.0"},
}

// foldNonFiniteLiterals rewrites each non-finite literal in a value position, skipping string contents.
func foldNonFiniteLiterals(text string) (string, map[int64]nonFinite) {
	if !strings.Contains(text, "NaN") && !strings.Contains(text, "Infinity") {
		return text, nil
	}
	var folded []byte
	var starts map[int64]nonFinite
	inString, escaped := false, false
	for i := 0; i < len(text); {
		c := text[i]
		switch {
		case inString && escaped:
			escaped = false
		case inString && c == '\\':
			escaped = true
		case inString && c == '"':
			inString = false
		case !inString && c == '"':
			inString = true
		case !inString:
			if literal, placeholder, ok := nonFiniteAt(text, i); ok {
				if folded == nil {
					folded = []byte(text)
					starts = map[int64]nonFinite{}
				}
				copy(folded[i:], placeholder)
				starts[int64(i)] = nonFinite(literal)
				i += len(literal)
				continue
			}
		}
		i++
	}
	if folded == nil {
		return text, nil
	}
	return string(folded), starts
}

// nonFiniteAt reports the literal at i, if one starts there and a JSON value may begin there.
func nonFiniteAt(text string, i int) (literal, placeholder string, ok bool) {
	if i > 0 && strings.IndexByte(" \t\r\n[,:", text[i-1]) < 0 {
		return "", "", false
	}
	for _, candidate := range nonFinitePlaceholders {
		if !strings.HasPrefix(text[i:], candidate.literal) {
			continue
		}
		if next := i + len(candidate.literal); next < len(text) && isWordByte(text[next]) {
			return "", "", false // 1Infinityx and friends: not this token
		}
		return candidate.literal, candidate.placeholder, true
	}
	return "", "", false
}

func isWordByte(c byte) bool {
	return c == '_' || c == '.' || ('0' <= c && c <= '9') ||
		('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

// orderedDecoder carries the folded positions, so a placeholder number goes back as its literal.
type orderedDecoder struct {
	dec       *json.Decoder
	nonFinite map[int64]nonFinite
}

// value decodes into an order-preserving tree: objects become jsonObject, arrays []any.
func (d *orderedDecoder) value() (any, error) {
	token, err := d.dec.Token()
	if err != nil {
		return nil, err
	}
	// Read immediately: InputOffset moves with every further token.
	if number, ok := token.(json.Number); ok && len(d.nonFinite) > 0 {
		start := d.dec.InputOffset() - int64(len(number.String()))
		if literal, ok := d.nonFinite[start]; ok {
			return literal, nil
		}
	}
	return d.valueFrom(token)
}

func (d *orderedDecoder) valueFrom(token json.Token) (any, error) {
	delim, ok := token.(json.Delim)
	if !ok {
		return token, nil
	}
	switch delim {
	case '{':
		object := jsonObject{}
		// Duplicate keys collapse as a decode into a map would: first position, last value.
		index := map[string]int{}
		for d.dec.More() {
			keyTok, err := d.dec.Token()
			if err != nil {
				return nil, err
			}
			key, ok := keyTok.(string)
			if !ok {
				return nil, fmt.Errorf("redact: object key %v is not a string", keyTok)
			}
			value, err := d.value()
			if err != nil {
				return nil, err
			}
			if at, ok := index[key]; ok {
				object[at].value = value
			} else {
				index[key] = len(object)
				object = append(object, jsonMember{key: key, value: value})
			}
		}
		if _, err := d.dec.Token(); err != nil { // consume '}'
			return nil, err
		}
		return object, nil
	case '[':
		array := []any{}
		for d.dec.More() {
			value, err := d.value()
			if err != nil {
				return nil, err
			}
			array = append(array, value)
		}
		if _, err := d.dec.Token(); err != nil { // consume ']'
			return nil, err
		}
		return array, nil
	default:
		return nil, fmt.Errorf("redact: unexpected delimiter %v", delim)
	}
}

// redactOrdered mirrors redact over the ordered tree.
func redactOrdered(value any, depth int, seen map[uintptr]bool) any {
	switch v := value.(type) {
	case jsonObject:
		if depth >= MaxRedactionDepth {
			return TooDeep
		}
		out := make(jsonObject, 0, len(v))
		for _, m := range v {
			if ShouldRedact(m.key) {
				out = append(out, jsonMember{key: m.key, value: Redacted})
			} else {
				out = append(out, jsonMember{key: m.key, value: redactOrdered(m.value, depth+1, seen)})
			}
		}
		return out
	case []any:
		if depth >= MaxRedactionDepth {
			return TooDeep
		}
		out := make([]any, 0, len(v))
		for _, entry := range v {
			out = append(out, redactOrderedEntry(entry, depth+1, seen))
		}
		return out
	case string:
		return RedactText(redactJSONText(v, depth, seen))
	default:
		return v
	}
}

func redactOrderedEntry(entry any, depth int, seen map[uintptr]bool) any {
	if name, isPair := pairName(entry); isPair && ShouldRedact(name) {
		return []any{name, Redacted}
	}
	return redactOrdered(entry, depth, seen)
}

func writeCanonical(out *strings.Builder, value any) {
	switch v := value.(type) {
	case jsonObject:
		out.WriteByte('{')
		for i, m := range v {
			if i > 0 {
				out.WriteString(", ")
			}
			writeCanonicalString(out, m.key)
			out.WriteString(": ")
			writeCanonical(out, m.value)
		}
		out.WriteByte('}')
	case []any:
		out.WriteByte('[')
		for i, entry := range v {
			if i > 0 {
				out.WriteString(", ")
			}
			writeCanonical(out, entry)
		}
		out.WriteByte(']')
	case string:
		writeCanonicalString(out, v)
	case json.Number:
		out.WriteString(canonicalNumberRepr(v))
	case nonFinite:
		// Bare, the spelling canonicalNumberRepr uses, so writer and reader agree.
		out.WriteString(string(v))
	case bool:
		if v {
			out.WriteString("true")
		} else {
			out.WriteString("false")
		}
	case nil:
		out.WriteString("null")
	}
}

// canonicalNumberRepr re-renders numbers canonically, so equal numbers rewrite to the same bytes.
func canonicalNumberRepr(n json.Number) string {
	s := n.String()
	if !strings.ContainsAny(s, ".eE") {
		if s == "-0" {
			return "0" // the canonical form has no negative zero integer
		}
		return s
	}
	f, _ := strconv.ParseFloat(s, 64) // range errors still yield ±Inf / 0, handled below
	switch {
	case math.IsInf(f, 1):
		return "Infinity" // an overflowing literal such as 1e999
	case math.IsInf(f, -1):
		return "-Infinity"
	}
	return canonicalFloatRepr(f)
}

// canonicalFloatRepr renders the shortest round-trip digits, in fixed notation iff -4 < decpt <= 16.
func canonicalFloatRepr(f float64) string {
	formatted := strconv.FormatFloat(f, 'e', -1, 64) // [-]d[.ddd]e±dd, shortest digits
	sign := ""
	if formatted[0] == '-' {
		sign = "-"
		formatted = formatted[1:]
	}
	mantissa, expText, _ := strings.Cut(formatted, "e")
	digits := strings.Replace(mantissa, ".", "", 1)
	exp, _ := strconv.Atoi(expText)
	decpt := exp + 1
	if -4 < decpt && decpt <= 16 {
		switch {
		case decpt <= 0:
			return sign + "0." + strings.Repeat("0", -decpt) + digits
		case decpt >= len(digits):
			return sign + digits + strings.Repeat("0", decpt-len(digits)) + ".0"
		default:
			return sign + digits[:decpt] + "." + digits[decpt:]
		}
	}
	head := digits[:1]
	if len(digits) > 1 {
		head += "." + digits[1:]
	}
	return fmt.Sprintf("%s%se%+03d", sign, head, exp)
}

// writeCanonicalString escapes quote, backslash and control characters; non-ASCII stays raw.
func writeCanonicalString(out *strings.Builder, s string) {
	out.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"':
			out.WriteString(`\"`)
		case '\\':
			out.WriteString(`\\`)
		case '\b':
			out.WriteString(`\b`)
		case '\f':
			out.WriteString(`\f`)
		case '\n':
			out.WriteString(`\n`)
		case '\r':
			out.WriteString(`\r`)
		case '\t':
			out.WriteString(`\t`)
		default:
			if r < 0x20 {
				fmt.Fprintf(out, `\u%04x`, r)
			} else {
				out.WriteRune(r)
			}
		}
	}
	out.WriteByte('"')
}

func RedactMapping(headers map[string]string) map[string]string {
	shown := make(map[string]string, len(headers))
	for k, v := range headers {
		if ShouldRedact(k) {
			shown[k] = Redacted
		} else {
			shown[k] = RedactText(v)
		}
	}
	return shown
}

// NeutraliseControls renders every control character and DEL as \xNN, so a value cannot end a log line.
func NeutraliseControls(value string) string {
	var escaped strings.Builder
	for _, r := range value {
		if r < 0x20 || r == 0x7f {
			fmt.Fprintf(&escaped, "\\x%02X", r)
			continue
		}
		escaped.WriteRune(r)
	}
	return escaped.String()
}

// RedactText redacts every URL, schemeless path?query and scheme-relative authority in free text.
func RedactText(value string) string {
	cleaned := urlRunRe.ReplaceAllStringFunc(value, redactOneURL)
	cleaned = pathQueryRunRe.ReplaceAllStringFunc(cleaned, redactOneURL)
	// Gated on an '@' so text with no userinfo never pays for it.
	if strings.Contains(cleaned, "@") {
		cleaned = redactAuthorityRuns(cleaned, authorityRunRe, ":", len("//"))
		cleaned = redactAuthorityRuns(cleaned, bareAuthorityRunRe, "/:", 0)
	}
	return cleaned
}

// redactAuthorityRuns masks the password in each match past prefix bytes, skipping one after a skipAfter byte.
func redactAuthorityRuns(text string, pattern *regexp.Regexp, skipAfter string, prefix int) string {
	var out strings.Builder
	cursor, position := 0, 0
	for position < len(text) {
		at := pattern.FindStringIndex(text[position:])
		if at == nil {
			break
		}
		start, end := position+at[0], position+at[1]
		if start > 0 && strings.IndexByte(skipAfter, text[start-1]) >= 0 {
			position = start + 1
			continue
		}
		out.WriteString(text[cursor : start+prefix])
		out.WriteString(redactUserinfo(text[start+prefix : end]))
		cursor, position = end, end
	}
	if cursor == 0 {
		return text // nothing matched: hand back the original bytes, no allocation
	}
	out.WriteString(text[cursor:])
	return out.String()
}

func redactOneURL(u string) string {
	// Split by hand so a malformed URL degrades to "unchanged", never panics.
	rest := u
	prefix := ""
	authority, newAuthority := "", ""
	if schemeEnd := strings.Index(rest, "://"); schemeEnd >= 0 {
		authorityEnd := indexOfAny(rest, schemeEnd+3, "/?#")
		authority = rest[schemeEnd+3 : authorityEnd]
		newAuthority = redactUserinfo(authority)
		prefix = rest[:schemeEnd+3] + newAuthority
		rest = rest[authorityEnd:]
	} else if strings.HasPrefix(rest, "//") {
		// Only pathQueryRunRe's runs arrive here schemeless, so this arm never sees a bare //user:pw@h.
		authorityEnd := indexOfAny(rest, 2, "/?#")
		authority = rest[2:authorityEnd]
		newAuthority = redactUserinfo(authority)
		prefix = "//" + newAuthority
		rest = rest[authorityEnd:]
	}
	query := strings.Index(rest, "?")
	fragment := strings.Index(rest, "#")
	var cut int
	switch {
	case query >= 0 && (fragment < 0 || query < fragment):
		cut = query
	case fragment >= 0:
		cut = fragment
	default:
		cut = -1
	}
	if cut < 0 {
		if newAuthority == authority {
			return u
		}
		return prefix + rest
	}
	// Split before either is judged: as one query string, "token=1#session=2" is one parameter.
	path := rest[:cut]
	tail := rest[cut+1:]
	queryPart, fragmentPart := "", ""
	if rest[cut] == '?' {
		queryPart = tail
		if before, after, found := strings.Cut(tail, "#"); found {
			queryPart, fragmentPart = before, after
		}
	} else {
		// A fragment alone is judged as a query string: "#session=2" carries implicit-flow tokens.
		fragmentPart = tail
	}
	newQuery := redactQueryString(queryPart)
	newFragment := redactQueryString(fragmentPart)
	// Rebuilt only on a match: a rebuild drops a trailing empty '?' or '#'.
	if newAuthority == authority && newQuery == queryPart && newFragment == fragmentPart {
		return u
	}
	out := prefix + path
	if newQuery != "" {
		out += "?" + newQuery
	}
	if newFragment != "" {
		out += "#" + newFragment
	}
	return out
}

func indexOfAny(text string, from int, targets string) int {
	for i := from; i < len(text); i++ {
		if strings.IndexByte(targets, text[i]) >= 0 {
			return i
		}
	}
	return len(text)
}

func redactUserinfo(authority string) string {
	at := strings.LastIndex(authority, "@")
	if at < 0 {
		return authority
	}
	userinfo := authority[:at]
	user, _, found := strings.Cut(userinfo, ":")
	if !found {
		return authority
	}
	return user + ":" + Redacted + authority[at:]
}

// redactQueryString leaves every parameter but the secret-named ones byte-exact.
func redactQueryString(query string) string {
	if query == "" {
		return query
	}
	chunks := strings.Split(query, "&")
	matched := false
	for i, chunk := range chunks {
		rawName, _, found := strings.Cut(chunk, "=")
		if !found {
			continue // a parameter with no '=' hides nothing
		}
		if ShouldRedact(safeDecode(rawName)) {
			chunks[i] = rawName + "=" + Redacted
			matched = true
		}
	}
	if matched {
		return strings.Join(chunks, "&")
	}
	return query
}

func safeDecode(s string) string {
	if decoded, err := url.QueryUnescape(s); err == nil {
		return decoded
	}
	return s
}
