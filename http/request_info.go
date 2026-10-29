package http

import (
	"fmt"
	"maps"
	"slices"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/redact"
)

// RequestInfo is everything about the request that produced a response.
type RequestInfo struct {
	// ClientName is the logical instance name.
	ClientName string
	// APIVersion is the API version segment.
	APIVersion string
	// HTTPMethod is the HTTP method.
	HTTPMethod string
	// URL is the request URL without the normalised query parameters, which are QueryParams below.
	URL string
	// EndpointTemplate is the path template this call was built from, placeholders unresolved.
	EndpointTemplate string
	// QueryParams are the query parameters, secret keys and values redacted.
	QueryParams map[string]any
	// Body is the request body, redacted and deep-copied.
	Body any
	// Headers are the request headers, already redacted.
	Headers map[string]string
	// AttemptCount is which attempt produced the response (1-based).
	AttemptCount int
	// MaxAttempts is the attempt ceiling in force.
	MaxAttempts int
	// ElapsedSeconds is the wall-clock seconds the exchange took.
	ElapsedSeconds float64
}

// RequestSummary lets the errors package read a RequestInfo without importing http.
func (r *RequestInfo) RequestSummary() (queryParams map[string]any, attemptCount, maxAttempts int,
	elapsedSeconds float64, httpMethod, url string) {
	if r == nil { // an error can carry a typed nil; reading it must not panic
		return nil, 0, 0, 0, "", ""
	}
	return redactQueryParams(r.QueryParams), r.AttemptCount, r.MaxAttempts, r.ElapsedSeconds,
		r.HTTPMethod, redact.RedactText(r.URL)
}

// RequestEndpoint reports the endpoint template, or "" when the sender supplied none.
func (r *RequestInfo) RequestEndpoint() string {
	if r == nil {
		return ""
	}
	return r.EndpointTemplate
}

// NewRequestInfo builds a record, redacting as it goes.
func NewRequestInfo(clientName, apiVersion, httpMethod, url string, queryParams map[string]any,
	body any, headers map[string]string, attemptCount, maxAttempts int, elapsedSeconds float64) *RequestInfo {
	if headers == nil {
		headers = map[string]string{}
	}
	return &RequestInfo{
		ClientName:     clientName,
		APIVersion:     apiVersion,
		HTTPMethod:     httpMethod,
		URL:            redact.RedactText(url),
		QueryParams:    redactQueryParams(queryParams),
		Body:           redact.RedactValue(body),
		Headers:        redact.RedactMapping(headers),
		AttemptCount:   attemptCount,
		MaxAttempts:    maxAttempts,
		ElapsedSeconds: elapsedSeconds,
	}
}

// String renders the record in the retention-safe form.
func (r RequestInfo) String() string {
	safe := r.Redacted()
	return fmt.Sprintf("RequestInfo{ClientName:%q APIVersion:%q HTTPMethod:%q Endpoint:%q "+
		"QueryParams:%v Body:%v Headers:%v AttemptCount:%d MaxAttempts:%d ElapsedSeconds:%v}",
		safe.ClientName, safe.APIVersion, safe.HTTPMethod,
		wcerrors.EndpointLabel(safe.EndpointTemplate, safe.URL),
		slices.Sorted(maps.Keys(safe.QueryParams)), safe.Body, slices.Sorted(maps.Keys(safe.Headers)),
		safe.AttemptCount, safe.MaxAttempts, safe.ElapsedSeconds)
}

// Redacted returns a copy with the four secret-bearing fields put through the redactor again.
func (r *RequestInfo) Redacted() *RequestInfo {
	if r == nil {
		return &RequestInfo{}
	}
	out := *r
	out.URL = redact.RedactText(r.URL)
	out.QueryParams = redactQueryParams(r.QueryParams)
	out.Headers = redact.RedactMapping(r.Headers)
	out.Body = redact.RedactValue(r.Body)
	return &out
}

func redactQueryParams(params map[string]any) map[string]any {
	shown := make(map[string]any, len(params))
	for key, value := range params {
		if redact.ShouldRedact(key) {
			shown[key] = redact.Redacted
		} else {
			shown[key] = redact.RedactValue(value)
		}
	}
	return shown
}
