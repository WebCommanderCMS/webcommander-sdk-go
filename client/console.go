package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"net/url"
	"os"
	"slices"
	"sort"
	"strings"

	wccfg "github.com/WebCommanderCMS/webcommander-sdk-go/v4/config"
	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/redact"
)

// RenderExchange renders one exchange as the console would print it.
func RenderExchange(response *http.ApiResponse, body any, files map[string]any,
	config wccfg.SdkConfig, formEncoded bool) string {
	return renderExchange(response, body, files, formEncoded, config.Normalised())
}

// BuildCurlCommand renders the copy-pasteable curl reproduction of one exchange.
func BuildCurlCommand(response *http.ApiResponse, payload any, files map[string]any,
	config wccfg.SdkConfig, formEncoded bool) string {
	// payload is not summarised again: summarising twice renders bogus -F 'files=<multipart: …>'.
	return buildCurl(response.RequestInfo.Redacted(), payload, files, formEncoded, config.Normalised())
}

func renderExchange(response *http.ApiResponse, reportedBody any, files map[string]any,
	formEncoded bool, config wccfg.SdkConfig) string {
	// Redacted, not read as it stands: a host-built &RequestInfo{} never met the redactor.
	info := response.RequestInfo.Redacted()
	payload := summarisePayload(reportedBody, files)
	status := statusLabel(response.StatusCode)

	if config.PrintFormat == wccfg.PrintFormatCompact {
		var b strings.Builder
		fmt.Fprintf(&b, "wc | %s | %s | %s | %s | %d/%d | %ss",
			info.ClientName, info.HTTPMethod, info.URL, status,
			info.AttemptCount, info.MaxAttempts, wcerrors.FormatSeconds(info.ElapsedSeconds))
		if config.PrintRequestPayload {
			b.WriteString(" | payload=")
			b.WriteString(RenderJSONValue(redact.RedactValue(payload), true))
		}
		if config.PrintResponseBody {
			b.WriteString(" | response=")
			b.WriteString(RenderJSONValue(redact.RedactValue(response.Map()), true))
		}
		return truncate(neutralise(b.String()), config.PrintMaxChars)
	}

	colour := newColour(config.PrintColor)

	lines := []string{colour.dim(fmt.Sprintf("WebCommander  %s  %s  %ss",
		info.HTTPMethod, status, wcerrors.FormatSeconds(info.ElapsedSeconds)))}

	if config.PrintApiInfo {
		lines = append(lines,
			infoRow("client", info.ClientName),
			infoRow("api version", info.APIVersion),
			infoRow("method", info.HTTPMethod),
			infoRow("endpoint", info.URL),
			infoRow("params", RenderJSONValue(info.QueryParams, true)),
			fmt.Sprintf("%-12s: %s", "status",
				colour.status(status, response.Succeeded())),
			infoRow("attempt", fmt.Sprintf("%d / %d", info.AttemptCount, info.MaxAttempts)),
			infoRow("elapsed", wcerrors.FormatSeconds(info.ElapsedSeconds)+"s"),
		)
	}
	if config.PrintHeaders {
		lines = append(lines, colour.section("request headers"))
		lines = append(lines, headerRows(info.Headers)...)
		lines = append(lines, colour.section("response headers"))
		lines = append(lines, headerRows(redact.RedactMapping(response.Headers))...)
	}
	if config.PrintRequestPayload && payload != nil {
		lines = append(lines, colour.section("payload"))
		lines = append(lines, rawBlock(RenderJSONValue(redact.RedactValue(payload), false))...)
	}
	if config.PrintResponseBody {
		lines = append(lines, colour.section("response"))
		// Map(), not Data, so the block reads the same whichever presentation Data carries.
		lines = append(lines, rawBlock(RenderJSONValue(redact.RedactValue(response.Map()), false))...)
	}
	if config.PrintCurl {
		lines = append(lines, colour.section("curl (POSIX shell)"))
		lines = append(lines, rawBlock(buildCurl(info, payload, files, formEncoded, config))...)
	}
	return truncate(strings.Join(lines, "\n"), config.PrintMaxChars)
}

// summarisePayload describes what was sent without reading a byte of any upload.
func summarisePayload(body any, files map[string]any) any {
	if len(files) == 0 {
		return body
	}
	names := make([]string, 0, len(files))
	for field := range files {
		names = append(names, field+"="+multipartFilename(field, files[field]))
	}
	sort.Strings(names)
	summary := multipartPrefix + strings.Join(names, ", ") + ">"
	if body == nil {
		return summary
	}
	return map[string]any{payloadFilesKey: summary, payloadFormKey: body}
}

// multipartFilename is the name the wire carries for one file entry: the [filename, bytes] form's own first.
func multipartFilename(field string, value any) string {
	if pair, ok := value.([]any); ok && len(pair) > 0 {
		if name, ok := pair[0].(string); ok && name != "" {
			return name
		}
	}
	return field
}

// formPart returns the form half of a summarisePayload result, or the payload itself.
func formPart(payload any) any {
	m, ok := payload.(map[string]any)
	if !ok || len(m) != 2 {
		return payload
	}
	summary, ok := m[payloadFilesKey].(string)
	if !ok || !strings.HasPrefix(summary, multipartPrefix) {
		return payload
	}
	return m[payloadFormKey]
}

// formFields returns the redacted, sorted name/value pairs of a form body, or nil for a non-mapping.
func formFields(body any) [][2]string {
	m, ok := redact.RedactValue(body).(map[string]any)
	if !ok {
		return nil
	}
	fields := make([][2]string, 0, len(m))
	for _, name := range slices.Sorted(maps.Keys(m)) {
		fields = append(fields, [2]string{name, fmt.Sprint(m[name])})
	}
	return fields
}

// buildCurl renders a pasteable POSIX-shell curl of the call, shell-quoted and with headers redacted.
func buildCurl(info *http.RequestInfo, payload any, files map[string]any,
	formEncoded bool, config wccfg.SdkConfig) string {
	var curl strings.Builder
	fmt.Fprintf(&curl, "curl -X %s %s", shellQuote(info.HTTPMethod), shellQuote(withQuery(info)))
	if !config.VerifySslEnabled() {
		curl.WriteString(" \\\n  -k")
	}
	// No --cacert for a configured CABundle: the pasted command must not need the recipient's file.
	for _, name := range slices.Sorted(maps.Keys(info.Headers)) {
		curl.WriteString(" \\\n  -H ")
		curl.WriteString(shellQuote(neutralise(name + ": " + info.Headers[name])))
	}

	body := formPart(payload)
	if len(files) > 0 {
		// One -F per file, neutralised as well as quoted, so a filename cannot erase the lines above it.
		for _, field := range slices.Sorted(maps.Keys(files)) {
			curl.WriteString(" \\\n  -F ")
			curl.WriteString(shellQuote(
				neutralise(field) + "=@" + neutralise(multipartFilename(field, files[field]))))
		}
		// Form fields ride in the same multipart body, so they are -F too.
		for _, field := range formFields(body) {
			curl.WriteString(" \\\n  -F ")
			curl.WriteString(shellQuote(neutralise(field[0]) + "=" + neutralise(field[1])))
		}
		return curl.String()
	}
	if body == nil {
		return curl.String()
	}
	if fields := formFields(body); formEncoded && len(fields) > 0 {
		// curl encodes each --data-urlencode field and sets the form content type itself.
		for _, field := range fields {
			curl.WriteString(" \\\n  --data-urlencode ")
			curl.WriteString(shellQuote(neutralise(field[0]) + "=" + neutralise(field[1])))
		}
		return curl.String()
	}
	if !formEncoded && !hasContentType(info.Headers) {
		// curl defaults --data-raw to form encoding; added only when the caller set none.
		curl.WriteString(" \\\n  -H ")
		curl.WriteString(shellQuote("Content-Type: application/json"))
	}
	// --data-raw carries the JSON body on one line, so a raw newline in it cannot forge a printed line.
	curl.WriteString(" \\\n  --data-raw ")
	curl.WriteString(shellQuote(RenderJSONValue(redact.RedactValue(body), true)))
	return curl.String()
}

// withQuery joins the recorded URL and query parameters into one pasteable URL.
func withQuery(info *http.RequestInfo) string {
	base := neutralise(info.URL)
	if len(info.QueryParams) == 0 {
		return base
	}
	var q strings.Builder
	q.WriteString(base)
	separator := "?"
	if hasQueryComponent(base) {
		separator = "&"
	}
	for _, key := range slices.Sorted(maps.Keys(info.QueryParams)) {
		for _, single := range http.ExpandValue(info.QueryParams[key]) {
			fmt.Fprintf(&q, "%s%s=%s", separator, queryEscapeKeepingStar(key), queryEscapeKeepingStar(single))
			separator = "&"
		}
	}
	return q.String()
}

// hasQueryComponent reports whether the URL already carries a non-empty query, choosing '?' or '&'.
func hasQueryComponent(url string) bool {
	head := url
	if hash := strings.Index(head, "#"); hash >= 0 {
		head = head[:hash]
	}
	question := strings.Index(head, "?")
	return question >= 0 && question+1 < len(head)
}

// queryEscapeKeepingStar is url.QueryEscape with '*' left literal, so a mask reads api_key=***.
func queryEscapeKeepingStar(s string) string {
	return strings.ReplaceAll(url.QueryEscape(s), "%2A", "*")
}

// statusLabel renders the status slot: the code, or "no response" when the exchange produced none.
func statusLabel(statusCode int) string {
	if statusCode == 0 {
		return noResponseStatus
	}
	return fmt.Sprintf("%d", statusCode)
}

func hasContentType(headers map[string]string) bool {
	for n := range headers {
		if strings.ToLower(n) == "content-type" {
			return true
		}
	}
	return false
}

// infoRow renders one "label: value" line, the value neutralised so it cannot forge a row.
func infoRow(label string, value string) string {
	return fmt.Sprintf("%-12s: %s", label, neutralise(value))
}

const (
	ansiEscape = "\x1b"
	ansiReset  = "\x1b[0m"

	multipartPrefix = "<multipart: "
	payloadFilesKey = "files"
	payloadFormKey  = "form"
)

// colourWriter emits ANSI escapes, or nothing at all.
type colourWriter struct{ enabled bool }

func newColour(requested bool) colourWriter {
	return colourWriter{enabled: requested && stdoutSupportsColour()}
}

func stdoutSupportsColour() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	info, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func (c colourWriter) dim(text string) string {
	if !c.enabled {
		return text
	}
	return "\x1b[2m" + text + ansiReset
}

func (c colourWriter) status(text string, succeeded bool) string {
	if !c.enabled {
		return text
	}
	code := "31"
	if succeeded {
		code = "32"
	}
	return "\x1b[" + code + "m" + text + ansiReset
}

func (c colourWriter) section(label string) string { return c.dim(section(label)) }

// headerRows renders one line per header, sorted, name and value both neutralised.
func headerRows(headers map[string]string) []string {
	var rows []string
	for _, name := range slices.Sorted(maps.Keys(headers)) {
		rows = append(rows, fmt.Sprintf("%-14s: %s", neutralise(name), neutralise(headers[name])))
	}
	return rows
}

// section renders a horizontal rule, "── label ────", 62 runes wide rather than 62 bytes.
func section(label string) string {
	rule, _ := ResolveBoxGlyphs()
	prefix := rule + rule + " " + label + " "
	n := max(62-len([]rune(prefix)), 4)
	return prefix + strings.Repeat(rule, n)
}

// rawBlock splits a rendered block into lines, dropping the empty segment after one trailing newline.
func rawBlock(text string) []string {
	if text == "" {
		return nil
	}
	lines := strings.Split(text, "\n")
	if last := len(lines) - 1; lines[last] == "" {
		lines = lines[:last]
	}
	return lines
}

// RenderJSONValue renders one value the way the console prints it: indented, or one line when compact.
func RenderJSONValue(value any, compact bool) string {
	if value == nil {
		return "<empty>"
	}
	if b, ok := value.([]byte); ok {
		return fmt.Sprintf("<%d bytes>", len(b))
	}
	if s, ok := value.(string); ok {
		// A text body is shown as text, not as a JSON scalar.
		if !compact {
			return wcerrors.EscapeControls(s, true)
		}
		// Raw JSON text is rendered as the document it is: encoding it would make the curl line send a string.
		if document, isJSON := compactJSONDocument(s); isJSON {
			return neutralise(document)
		}
		encoded, _ := http.EncodeJSONBody(s)
		quoted := string(encoded)
		if len(quoted) >= 2 && quoted[0] == '"' && quoted[len(quoted)-1] == '"' {
			quoted = quoted[1 : len(quoted)-1]
		}
		return neutralise(quoted)
	}
	// Nested binary is summarised first: encoding/json would base64 a nested []byte into the output.
	value = summariseNestedBinary(value, 0)
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	if !compact {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(value); err != nil {
		return wcerrors.EscapeControls(fmt.Sprint(value), !compact)
	}
	return wcerrors.EscapeControls(strings.TrimRight(buffer.String(), "\n"), !compact)
}

// summariseNestedBinary replaces every nested []byte with the "<N bytes>" summary, depth-capped.
func summariseNestedBinary(value any, depth int) any {
	if depth >= redact.MaxRedactionDepth {
		return value
	}
	switch v := value.(type) {
	case []byte:
		return fmt.Sprintf("<%d bytes>", len(v))
	case map[string]any:
		out := make(map[string]any, len(v))
		for key, entry := range v {
			out[key] = summariseNestedBinary(entry, depth+1)
		}
		return out
	case []any:
		out := make([]any, 0, len(v))
		for _, entry := range v {
			out = append(out, summariseNestedBinary(entry, depth+1))
		}
		return out
	}
	return value
}

// compactJSONDocument reports the one-line form of a string that is a JSON object or array.
func compactJSONDocument(text string) (string, bool) {
	trimmed := strings.TrimSpace(text)
	if !strings.HasPrefix(trimmed, "{") && !strings.HasPrefix(trimmed, "[") {
		return "", false
	}
	var out bytes.Buffer
	if err := json.Compact(&out, []byte(trimmed)); err != nil {
		return "", false
	}
	return out.String(), true
}

// shellQuote quotes one curl fragment for a POSIX shell only where it needs it.
func shellQuote(fragment string) string {
	if fragment != "" && !strings.ContainsFunc(fragment, shellUnsafe) {
		return fragment
	}
	return "'" + strings.ReplaceAll(fragment, "'", `'"'"'`) + "'"
}

// shellUnsafe reports whether a rune forces the fragment to be quoted.
func shellUnsafe(r rune) bool {
	switch {
	case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
		return false
	case r == '_', r == '@', r == '%', r == '+', r == '=', r == ':', r == ',',
		r == '.', r == '/', r == '-':
		return false
	}
	return true
}

// neutralise escapes every control character, newlines included, so a one-row value cannot forge a row.
func neutralise(text string) string { return wcerrors.EscapeControls(text, false) }

// truncate cuts a rendered block to maxChars in characters, notice included, closing a cut colour run.
func truncate(text string, maxChars int) string {
	if maxChars <= 0 {
		return text
	}
	runes := []rune(text)
	if len(runes) <= maxChars {
		return text
	}
	_, ellipsis := ResolveBoxGlyphs()
	reserved := len([]rune(fmt.Sprintf("%s +%d chars", ellipsis, len(runes))))
	if strings.Contains(text, ansiEscape) {
		reserved += len(ansiReset)
	}
	cut := max(maxChars-reserved, 0)
	return fmt.Sprintf("%s%s +%d chars", closeColourRun(string(runes[:cut])), ellipsis, len(runes)-cut)
}

// closeColourRun repairs an ANSI sequence a cut landed inside: a half escape is dropped, a run closed.
func closeColourRun(text string) string {
	if !strings.Contains(text, ansiEscape) {
		return text
	}
	if marker := strings.LastIndex(text, ansiEscape); !strings.Contains(text[marker:], "m") {
		text = text[:marker]
	}
	last := strings.LastIndex(text, "\x1b[")
	if last != -1 && !strings.HasPrefix(text[last:], ansiReset) {
		return text + ansiReset
	}
	return text
}
