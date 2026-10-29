package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/url"
	"strings"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
)

// TransportRequest is one request's content in one value.
type TransportRequest struct {
	// HTTPMethod is the HTTP method.
	HTTPMethod string
	// URL is the absolute request URL.
	URL string
	// EndpointTemplate is the path template this call was built from, placeholders unresolved.
	EndpointTemplate string
	// QueryParams are parameters to append, already normalised by the caller.
	QueryParams map[string]any
	// JSONBody is a JSON-encodable value, or a raw JSON string whose bytes are sent as they stand.
	JSONBody any
	// FormData is form-encoded data to send.
	FormData map[string]any
	// Files are multipart files to send: each value is []byte or a string, or [filename, contents[, contentType]].
	Files map[string]any
	// Headers are headers for this request.
	Headers map[string]string
	// RequestInfoHeaders are already-redacted headers to record; nil records nothing.
	RequestInfoHeaders map[string]string
	// RequestInfoBody is an already-redacted body to record on the request info.
	RequestInfoBody any
	// ResponseFormat is the already-resolved format for this call's Data; "" is the built-in one.
	ResponseFormat string
	// DTOShape is the typed view this call produces for DTO(); nil where there is no typed form.
	DTOShape DTOShape

	// wire is the assembled body, built once by Materialised and reused unchanged by every attempt.
	wire *wireBody
}

// wireBodyKind matters once: a multipart Content-Type overrides the caller's, the other two defer to it.
type wireBodyKind int

const (
	wireJSON wireBodyKind = iota
	wireForm
	wireMultipart
)

// wireBody is the exact bytes and their Content-Type; the final-wire ceiling is enforced on its length.
type wireBody struct {
	kind        wireBodyKind
	contentType string
	body        []byte
}

// Materialised resolves a request for sending: reference fields copied, JSON rendered, body assembled once.
func (r TransportRequest) Materialised() (TransportRequest, error) {
	out := r
	out.Headers = snapshotStringMap(r.Headers)
	out.RequestInfoHeaders = snapshotStringMap(r.RequestInfoHeaders)
	out.RequestInfoBody = snapshotValue(r.RequestInfoBody)
	out.Files = snapshotAnyMap(r.Files)
	// The body's containers are copied first, or a String method on a leaf could rewrite the mapping.
	out.JSONBody = snapshotValue(r.JSONBody)
	querySource := snapshotAnyMap(r.QueryParams)
	formSource := snapshotAnyMap(r.FormData)

	out.QueryParams = normaliseWireLeaves(querySource)
	out.FormData = normaliseWireLeaves(formSource)
	var encodedJSON []byte
	if out.JSONBody != nil {
		encoded, err := encodeRequestJSONBody(out.JSONBody)
		if err != nil {
			return TransportRequest{}, RefuseUnencodableJSONBody(err)
		}
		out.JSONBody = string(encoded)
		// Kept as bytes too, so assembly skips a second 64 MiB copy back from text.
		encodedJSON = encoded
	}
	// Assembled once, and skipped when already assembled; both wire ceilings are enforced here.
	if out.wire == nil {
		assembled, err := assembleWireBody(out, encodedJSON)
		if err != nil {
			return TransportRequest{}, err
		}
		out.wire = assembled
	}
	return out, nil
}

// MaxRequestBodyBytes caps the final wire body of one request, whichever kind it is: 64 MiB.
const MaxRequestBodyBytes = 64 << 20

// assembleWireBody renders a materialised body into the exact bytes to send; nil, nil for no body.
func assembleWireBody(request TransportRequest, encodedJSON []byte) (*wireBody, error) {
	switch {
	case len(request.Files) > 0:
		return assembleMultipartBody(request)
	case request.JSONBody != nil:
		if err := refuseOversizedWireBody("jsonBody", "the JSON document", len(encodedJSON)); err != nil {
			return nil, err
		}
		return &wireBody{kind: wireJSON, contentType: "application/json; charset=utf-8", body: encodedJSON}, nil
	case request.FormData != nil:
		encoded, err := encodeFormBody(request.FormData)
		if err != nil {
			return nil, err
		}
		return &wireBody{kind: wireForm, contentType: "application/x-www-form-urlencoded", body: encoded}, nil
	}
	return nil, nil
}

// encodeFormBody produces the same bytes url.Values.Encode would, fields sorted by name.
func encodeFormBody(formData map[string]any) ([]byte, error) {
	var b strings.Builder
	for _, name := range sortedMapKeys(formData) {
		chunk := url.QueryEscape(name) + "=" + url.QueryEscape(stringifyValue(formData[name]))
		if b.Len() > 0 {
			chunk = "&" + chunk
		}
		if err := refuseOversizedWireBody("formData", "the URL-encoded form body",
			b.Len()+len(chunk)); err != nil {
			return nil, err
		}
		b.WriteString(chunk)
	}
	return []byte(b.String()), nil
}

// refuseOversizedWireBody names the cap and the field that raises it, for all three body kinds.
func refuseOversizedWireBody(field, what string, size int) error {
	if size <= MaxRequestBodyBytes {
		return nil
	}
	return wcerrors.NewUsageError(fmt.Sprintf(
		"%s: %s is %d bytes, over the %d-byte cap this SDK assembles in memory; send a body this "+
			"size with a request of your own, streaming it from disk rather than buffering it whole",
		field, what, size, MaxRequestBodyBytes))
}

// normaliseWireLeaves renders every leaf as the text the wire will carry.
func normaliseWireLeaves(source map[string]any) map[string]any {
	if source == nil {
		return nil
	}
	out := make(map[string]any, len(source))
	for _, name := range sortedMapKeys(source) {
		out[name] = normaliseWireValue(source[name])
	}
	return out
}

// normaliseWireValue mirrors ExpandValue's three arms, so expanding the result gives the same values.
func normaliseWireValue(value any) any {
	switch held := value.(type) {
	case string:
		return held
	case []any:
		if held == nil {
			return held
		}
		out := make([]any, len(held))
		for index, item := range held {
			out[index] = stringifyValue(item)
		}
		return out
	default:
		return stringifyValue(value)
	}
}

// snapshotAnyMap deep-copies the JSON-shaped values inside. nil in, nil out.
func snapshotAnyMap(source map[string]any) map[string]any {
	if source == nil {
		return nil
	}
	out := make(map[string]any, len(source))
	for name, value := range source {
		out[name] = snapshotValue(value)
	}
	return out
}

func snapshotStringMap(source map[string]string) map[string]string {
	return maps.Clone(source)
}

// SnapshotJSONValue deep-copies one JSON-shaped value at any depth; anything else is returned as is.
func SnapshotJSONValue(value any) any { return snapshotValue(value) }

// snapshotValue rebuilds mappings, lists and byte slices and returns anything else as it stands.
func snapshotValue(value any) any {
	switch held := value.(type) {
	case map[string]any:
		return snapshotAnyMap(held)
	case map[string]string:
		return snapshotStringMap(held)
	case []any:
		if held == nil {
			return held
		}
		out := make([]any, len(held))
		for index, entry := range held {
			out[index] = snapshotValue(entry)
		}
		return out
	case []byte:
		if held == nil {
			return held
		}
		out := make([]byte, len(held))
		copy(out, held)
		return out
	case []string:
		if held == nil {
			return held
		}
		out := make([]string, len(held))
		copy(out, held)
		return out
	default:
		return value
	}
}

// DTOPayload is a typed request body: anything that can hand over the map the SDK encodes.
type DTOPayload interface {
	// ToMap renders this payload as the body the SDK encodes.
	ToMap() map[string]any
}

// Requester is the subset of the client that resource modules call.
type Requester interface {
	SendApiRequest(ctx context.Context, httpMethod, path string, options *RequestOptions) (*ApiResponse, error)
	SendGetRequest(ctx context.Context, path string, options *RequestOptions) (*ApiResponse, error)
	SendPostRequest(ctx context.Context, path string, options *RequestOptions) (*ApiResponse, error)
	SendPutRequest(ctx context.Context, path string, options *RequestOptions) (*ApiResponse, error)
	SendPatchRequest(ctx context.Context, path string, options *RequestOptions) (*ApiResponse, error)
	SendDeleteRequest(ctx context.Context, path string, options *RequestOptions) (*ApiResponse, error)
}

// RefuseUnencodableJSONBody is the one refusal for a JSON body encoding/json will not render.
func RefuseUnencodableJSONBody(err error) error {
	refusal := wcerrors.NewUsageError("jsonBody: could not be encoded to JSON — " + describeJSONFailure(err))
	refusal.Cause = err
	refusal.CauseText = err.Error()
	return refusal
}

// describeJSONFailure names why encoding/json refused a value, structurally.
func describeJSONFailure(err error) string {
	var unsupportedType *json.UnsupportedTypeError
	if errors.As(err, &unsupportedType) {
		return fmt.Sprintf("it holds a %s, which has no JSON form", unsupportedType.Type)
	}
	var unsupportedValue *json.UnsupportedValueError
	if errors.As(err, &unsupportedValue) {
		// Named together because the distinguishing text is the value's own.
		return "it holds a value with no JSON form — a NaN, an infinity, or a cycle"
	}
	var marshaler *json.MarshalerError
	if errors.As(err, &marshaler) {
		return fmt.Sprintf("the MarshalJSON method on %s refused it; its own words are on Explain()",
			marshaler.Type)
	}
	return "encoding/json refused it; its own words are on Explain()"
}

// EncodeJSONBody encodes without HTML escaping or json.Encoder's trailing newline: the wire's exact bytes.
func EncodeJSONBody(v any) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buffer.Bytes(), "\n"), nil
}

// encodeRequestJSONBody passes a string body through untouched: it is already JSON, verified upstream.
func encodeRequestJSONBody(jsonBody any) ([]byte, error) {
	if raw, ok := jsonBody.(string); ok {
		return []byte(raw), nil
	}
	return EncodeJSONBody(jsonBody)
}
