package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"
	"sync"
	"unicode/utf8"

	wccfg "github.com/WebCommanderCMS/webcommander-sdk-go/v4/config"
	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/redact"
)

const (
	maxReportedMessages = 50
	maxMessageDepth     = 4
)

// DTOShape builds the typed view of one response body, or nil where there is no typed form.
type DTOShape func(decoded any) any

type ApiResponse struct {
	StatusCode int
	// Data is the body under the resolved response format, live and never a copy; nil when there is no such form.
	Data        any
	Headers     map[string]string
	ContentType string
	RequestInfo *RequestInfo
	ErrorType   ErrorType

	// Unexported, so an ApiResponse stays buildable as a literal; each view falls back when unset.
	decoded      any      // the decoded body, whatever Data was moved to
	bodyRecorded bool     // whether decoded is the answer — a 204 decodes to nil, legitimately
	rawBody      []byte   // the body as it arrived — not a copy, the slice the transport read
	bodyIsText   bool     // whether the content type made those bytes text
	dtoShape     DTOShape // the typed view this endpoint produces, or nil

	// views guards the five fields below and nothing else.
	views       sync.Mutex
	jsonBuilt   bool
	jsonText    string
	jsonHasText bool // whether the body has a text form at all, which keeps Data nil rather than "" under JSON
	dtoBuilt    bool
	dtoView     any
}

// Map is always the decoded body, whatever ResponseFormat says, and not a copy.
func (r *ApiResponse) Map() any {
	if r.bodyRecorded {
		return r.decoded
	}
	return r.Data
}

// JSON is always the server's own body text, the bytes as they arrived, never re-serialised.
func (r *ApiResponse) JSON() string {
	text, _ := r.jsonView()
	return text
}

// jsonView is JSON() plus whether this body has a text form at all.
func (r *ApiResponse) jsonView() (text string, hasText bool) {
	r.views.Lock()
	defer r.views.Unlock()
	if !r.jsonBuilt {
		r.jsonText, r.jsonHasText = r.buildJSON()
		r.jsonBuilt = true
	}
	return r.jsonText, r.jsonHasText
}

// buildJSON renders the body text once; called under the views lock.
func (r *ApiResponse) buildJSON() (text string, hasText bool) {
	if text, isText := r.Map().(string); isText {
		return text, true
	}
	if !r.bodyIsText || len(r.rawBody) == 0 {
		return "", false
	}
	return string(r.rawBody), true
}

// DTO is always the typed view of this body, built from the shape the endpoint declared.
func (r *ApiResponse) DTO() any {
	r.views.Lock()
	defer r.views.Unlock()
	if !r.dtoBuilt {
		if r.dtoShape != nil {
			r.dtoView = r.dtoShape(r.Map())
		}
		r.dtoBuilt = true
	}
	return r.dtoView
}

// present records the body this response was decoded from and moves Data to the resolved format.
func (r *ApiResponse) present(rawBody []byte, contentType, responseFormat string, shape DTOShape) {
	r.decoded = r.Data
	r.bodyRecorded = true
	r.rawBody = rawBody
	r.bodyIsText = bodyIsTextual(contentType)
	r.dtoShape = shape
	r.jsonBuilt, r.jsonText, r.jsonHasText = false, "", false
	r.dtoBuilt, r.dtoView = false, nil
	switch responseFormat {
	case wccfg.ResponseFormatJSON:
		// nil, not "", for a body with no text form, which "" would misreport as an empty body.
		if text, hasText := r.jsonView(); hasText {
			r.Data = text
		} else {
			r.Data = nil
		}
	case wccfg.ResponseFormatDTO:
		r.Data = r.DTO()
	}
}

// bodyIsTextual applies DecodeResponseBody's two rules, so JSON() and Map() agree on text bodies.
func bodyIsTextual(contentType string) bool {
	lowered := strings.ToLower(contentType)
	return isJSONMediaType(lowered) || strings.HasPrefix(lowered, "text/")
}

const jsonMediaType = "application/json"

// isJSONMediaType reports whether a lower-cased Content-Type names JSON, "+json" suffixes included.
func isJSONMediaType(loweredContentType string) bool {
	mediaType := loweredContentType
	if semicolon := strings.IndexByte(mediaType, ';'); semicolon >= 0 {
		mediaType = mediaType[:semicolon]
	}
	mediaType = strings.TrimSpace(mediaType)
	return mediaType == jsonMediaType || strings.HasSuffix(mediaType, "+json")
}

// NewApiResponse builds a response, deriving the classification from the status when none is given.
func NewApiResponse(statusCode int, data any, headers map[string]string, contentType string,
	requestInfo *RequestInfo, errorType ErrorType) (*ApiResponse, error) {
	if !knownErrorType(errorType) {
		available := make([]string, 0, len(errorTypeMembers))
		for _, member := range errorTypeMembers {
			available = append(available, string(member))
		}
		return nil, wcerrors.NewUsageError(fmt.Sprintf(
			"unknown error_type '%s'; expected an ErrorType or one of: %s",
			errorType, strings.Join(available, ", ")))
	}
	if headers == nil {
		headers = map[string]string{}
	}
	denied := bodyDeniesSuccess(data)
	if errorType == ErrorTypeNone && (statusCode >= 300 || denied) {
		errorType = classifyErrorType(statusCode, denied)
	}
	return &ApiResponse{
		StatusCode:  statusCode,
		Data:        data,
		Headers:     headers,
		ContentType: contentType,
		RequestInfo: requestInfo,
		ErrorType:   errorType,
	}, nil
}

// Succeeded reports whether the status is 2xx and the body agrees.
func (r *ApiResponse) Succeeded() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300 && !bodyDeniesSuccess(r.Map())
}

// bodyDeniesSuccess keys on isSuccess alone: this API puts a message in its success envelopes too.
func bodyDeniesSuccess(data any) bool {
	body, isMapping := data.(map[string]any)
	if !isMapping {
		return false
	}
	flag, present := body["isSuccess"]
	return present && !isTruthy(flag)
}

// isTruthy is JSON truthiness spelled out, Go having none.
func isTruthy(value any) bool {
	switch v := value.(type) {
	case nil:
		return false
	case bool:
		return v
	case json.Number:
		// Judged as a number although a string underneath: json.Number("0") is zero.
		if number, err := v.Float64(); err == nil {
			return number != 0
		}
		return strings.TrimSpace(v.String()) != ""
	case float64:
		return v != 0
	case string:
		return len(v) > 0
	case []any:
		return len(v) > 0
	case map[string]any:
		return len(v) > 0
	}
	return true
}

// GetErrorMessages reads the error envelope, most specific first, capped at maxReportedMessages.
func (r *ApiResponse) GetErrorMessages() []string {
	body := r.Map()
	if _, ok := body.(map[string]any); !ok {
		body = parseJSONBody(body)
	}
	m, ok := body.(map[string]any)
	if !ok {
		return nil
	}
	messages := readErrorEntries(m["errors"])
	if len(messages) == 0 {
		messages = readMessageValue(m["message"], 0)
	}
	if len(messages) == 0 {
		messages = readMessageValue(m["error"], 0)
		if len(messages) == 0 {
			messages = readMessageValue(m["errorCode"], 0)
		}
	}
	if len(messages) == 0 {
		// RFC 7807, reached only when an intermediary answers the error rather than the store.
		messages = readMessageValue(m["detail"], 0)
		if len(messages) == 0 {
			messages = readMessageValue(m["title"], 0)
		}
	}
	return capMessages(messages)
}

// RaiseForStatus returns an *SdkError when the call failed, with body and headers redacted, else nil.
func (r *ApiResponse) RaiseForStatus() error {
	if r.Succeeded() {
		return nil
	}
	messages := r.GetErrorMessages()
	return wcerrors.NewResponseError(classificationSummary(r.ErrorType, len(messages)),
		r.StatusCode, r.ErrorType,
		redact.RedactValue(r.Map()), redact.RedactMapping(r.Headers), r.RequestInfo, messages)
}

// classificationSummary is the classification and the message count, none of the messages' text.
func classificationSummary(errorType ErrorType, reported int) string {
	name := string(errorType)
	if name == "" || errorType == ErrorTypeNone {
		name = "error"
	}
	if reported == 0 {
		return name + " response"
	}
	return fmt.Sprintf("%s response (%d reported)", name, reported)
}

func parseJSONBody(body any) any {
	var text string
	switch b := body.(type) {
	case []byte:
		text = string(b)
	case string:
		text = b
	default:
		return nil
	}
	// Strip, unprefix, strip again: "\uFEFF  {...}" would otherwise fail the '{' guard.
	text = strings.TrimSpace(text)
	text = strings.TrimLeft(text, "\uFEFF")
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, "{") && !strings.HasPrefix(text, "[") {
		return nil
	}
	if v, ok := DecodeJSON([]byte(text)); ok {
		if m, ok := v.(map[string]any); ok {
			return m
		}
	}
	return nil
}

func readErrorEntries(entries any) []string {
	var messages []string
	switch e := entries.(type) {
	case map[string]any:
		// Sorted: a Go map ranges randomly and SdkError.Message must be stable across runs.
		for _, field := range slices.Sorted(maps.Keys(e)) {
			prefixField(field, readMessageValue(e[field], 0), &messages)
		}
		return messages
	case []any:
		for _, entry := range e {
			if m, ok := entry.(map[string]any); ok {
				text := m["message"]
				if text == nil {
					text = m["code"]
				}
				prefixField(fieldName(m["field"]), readMessageValue(text, 0), &messages)
			} else {
				messages = append(messages, readMessageValue(entry, 0)...)
			}
		}
		return messages
	case string:
		return readMessageValue(e, 0)
	}
	return messages
}

func readMessageValue(value any, depth int) []string {
	if value == nil || depth > maxMessageDepth {
		return nil
	}
	switch v := value.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return nil
		}
		return []string{v}
	case map[string]any:
		text := v["message"]
		if text == nil {
			text = v["code"]
		}
		return readMessageValue(text, depth+1)
	case []any:
		var messages []string
		for _, item := range v {
			messages = append(messages, readMessageValue(item, depth+1)...)
		}
		return messages
	default:
		text := fmt.Sprint(v)
		if strings.TrimSpace(text) == "" {
			return nil
		}
		return []string{text}
	}
}

// prefixField tests the field for presence, not truthiness.
func prefixField(field string, fieldMessages []string, out *[]string) {
	// Blankness is judged trimmed, but the prefix keeps the store's own spelling.
	blank := strings.TrimSpace(field) == ""
	for _, m := range fieldMessages {
		if blank {
			*out = append(*out, m)
		} else {
			*out = append(*out, field+": "+m)
		}
	}
}

func fieldName(v any) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

func capMessages(messages []string) []string {
	if len(messages) <= maxReportedMessages {
		return messages
	}
	capped := append([]string{}, messages[:maxReportedMessages]...)
	capped = append(capped, fmt.Sprintf("+%d more problems", len(messages)-maxReportedMessages))
	return capped
}

// DecodeResponseBody decodes JSON to maps and lists, text/* to a string, else bytes; 204 or empty is nil.
func DecodeResponseBody(status int, body []byte, contentType string) any {
	if status == 204 || len(body) == 0 {
		return nil
	}
	lowered := strings.ToLower(contentType)
	if isJSONMediaType(lowered) {
		if v, ok := DecodeJSON(body); ok {
			return v
		}
		return string(body)
	}
	if strings.HasPrefix(lowered, "text/") {
		return decodeTextBody(body, lowered)
	}
	return body
}

// Tested as a substring, so a malformed parameter list cannot make body decoding fail.
const charsetParameter = "charset="

// decodeTextBody tries the declared charset, then UTF-8, then ISO-8859-1, which cannot fail.
func decodeTextBody(body []byte, loweredContentType string) string {
	if strings.Contains(loweredContentType, charsetParameter) {
		if decode := resolveCodec(declaredCharset(loweredContentType)); decode != nil {
			return decode(body)
		}
	}
	if utf8.Valid(body) {
		return string(body)
	}
	return decodeLatin1(body)
}

// declaredCharset reads the value after "charset="; lower-cased already by the caller.
func declaredCharset(loweredContentType string) string {
	value := loweredContentType[strings.Index(loweredContentType, charsetParameter)+len(charsetParameter):]
	if semicolon := strings.IndexByte(value, ';'); semicolon >= 0 {
		value = value[:semicolon]
	}
	return strings.Trim(strings.TrimSpace(value), `"'`)
}

// resolveCodec covers the charsets this package decodes without a dependency; anything else is nil.
func resolveCodec(name string) func([]byte) string {
	switch name {
	case "utf-8", "utf8":
		return decodeUTF8Replacing
	case "us-ascii", "ascii":
		return decodeASCIIReplacing
	case "iso-8859-1", "iso8859-1", "iso_8859-1", "latin-1", "latin1", "latin", "l1", "cp819":
		return decodeLatin1
	}
	return nil
}

func decodeUTF8Replacing(body []byte) string {
	if utf8.Valid(body) {
		return string(body)
	}
	return strings.ToValidUTF8(string(body), string(utf8.RuneError))
}

func decodeASCIIReplacing(body []byte) string {
	runes := make([]rune, len(body))
	for i, b := range body {
		if b < 0x80 {
			runes[i] = rune(b)
		} else {
			runes[i] = utf8.RuneError
		}
	}
	return string(runes)
}

// decodeLatin1 maps each byte to the code point of the same value; it cannot fail.
func decodeLatin1(body []byte) string {
	runes := make([]rune, len(body))
	for i, b := range body {
		runes[i] = rune(b)
	}
	return string(runes)
}

// DecodeJSON decodes exactly one JSON document, numbers kept as json.Number, and reports whether it parsed.
func DecodeJSON(body []byte) (any, bool) {
	dec := json.NewDecoder(bytes.NewReader(stripBOM(body)))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, false
	}
	if _, err := dec.Token(); !errors.Is(err, io.EOF) {
		return nil, false
	}
	return v, true
}

// JSONKind names a value's shape in JSON's vocabulary, never the value itself.
func JSONKind(value any) string {
	switch value.(type) {
	case nil:
		return "null"
	case bool:
		return "a boolean"
	case float64, float32, int, int32, int64, json.Number:
		return "a number"
	case string:
		return "a string"
	case map[string]any:
		return "an object"
	case []any:
		return "an array"
	}
	return "not a JSON value"
}

func stripBOM(b []byte) []byte {
	if len(b) >= 3 && b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF {
		return b[3:]
	}
	return b
}
