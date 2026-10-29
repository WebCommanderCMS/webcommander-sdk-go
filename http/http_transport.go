package http

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"maps"
	"mime/multipart"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"

	wccfg "github.com/WebCommanderCMS/webcommander-sdk-go/v4/config"
	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/redact"
)

const contentTypeHeader = "Content-Type"

type HTTPTransport struct {
	config     wccfg.SdkConfig
	retry      *RetryPolicy
	client     *http.Client
	onProgress func(string)
}

const unknownParameterError = "unknown_parameter"

var (
	resendableMethods = map[string]bool{"GET": true, "HEAD": true, "OPTIONS": true}
	parameterNameRe   = regexp.MustCompile(`[A-Za-z0-9_-]+`)
)

type TransportOption func(*transportOptions)

type transportOptions struct {
	retry      *RetryPolicy
	client     *http.Client
	onProgress func(string)
}

// WithRetryPolicy replaces the policy the config would derive; nil leaves the config's own.
func WithRetryPolicy(retry *RetryPolicy) TransportOption {
	return func(o *transportOptions) { o.retry = retry }
}

// WithHTTPClient injects the net/http client to send through; only its CheckRedirect is replaced.
func WithHTTPClient(client *http.Client) TransportOption {
	return func(o *transportOptions) { o.client = client }
}

// WithProgress wires the per-attempt progress hook.
func WithProgress(onProgress func(string)) TransportOption {
	return func(o *transportOptions) { o.onProgress = onProgress }
}

// NewHTTPTransport builds a transport for a config, validating it first.
func NewHTTPTransport(config wccfg.SdkConfig, options ...TransportOption) (*HTTPTransport, error) {
	c := config.Normalised()
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return newHTTPTransport(c, options...)
}

func newHTTPTransport(config wccfg.SdkConfig, options ...TransportOption) (*HTTPTransport, error) {
	applied := &transportOptions{}
	for _, option := range options {
		option(applied)
	}
	retry := applied.retry
	if retry == nil {
		derived, err := NewRetryPolicyFromConfig(config)
		if err != nil {
			return nil, err
		}
		retry = derived
	}
	frozenRetry, err := retry.frozen()
	if err != nil {
		return nil, err
	}
	retry = frozenRetry
	client := applied.client
	if client == nil {
		built, err := buildHTTPClient(config)
		if err != nil {
			return nil, err
		}
		client = built
	} else {
		client = withoutRedirects(client)
	}
	onProgress := applied.onProgress
	if onProgress == nil {
		onProgress = func(string) {}
	}
	return &HTTPTransport{config: config, retry: retry, client: client, onProgress: onProgress}, nil
}

// withoutRedirects copies a caller's client so it cannot follow a redirect, sharing its Transport.
func withoutRedirects(client *http.Client) *http.Client {
	copied := *client
	copied.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return &copied
}

func buildHTTPClient(config wccfg.SdkConfig) (*http.Client, error) {
	tlsConfig := &tls.Config{InsecureSkipVerify: !config.VerifySslEnabled()} // #nosec G402 -- opt-in, and off only when the caller sets VerifySsl=false
	if config.CABundle != "" {
		// Neither refusal below renders the path or the OS's own words; both ride CauseText for Explain.
		pem, err := ReadCappedFile(config.CABundle, MaxCABundleBytes)
		if err != nil {
			refusal := wcerrors.NewConfigError(fmt.Sprintf(
				"caBundle: %s could not be read", redact.DescribePath(config.CABundle)))
			refusal.Cause = err
			refusal.CauseText = fmt.Sprintf("%s: %v", config.CABundle, err)
			return nil, refusal
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(pem) {
			refusal := wcerrors.NewConfigError(fmt.Sprintf(
				"caBundle: %s holds no PEM certificates", redact.DescribePath(config.CABundle)))
			refusal.CauseText = config.CABundle
			return nil, refusal
		}
		tlsConfig.RootCAs = pool
	}
	dialer := &net.Dialer{Timeout: durationSeconds(config.ConnectTimeoutSeconds)}
	transport := &http.Transport{
		DialContext:         dialer.DialContext,
		MaxIdleConns:        config.PoolMaxsize * 2,
		MaxIdleConnsPerHost: config.PoolMaxsize,
		ForceAttemptHTTP2:   true,
		TLSClientConfig:     tlsConfig,
	}
	if config.ProxyURL != "" {
		p, _ := url.Parse(config.ProxyURL)
		transport.Proxy = http.ProxyURL(p)
	} else {
		transport.Proxy = http.ProxyFromEnvironment
	}
	return &http.Client{
		Transport:     transport,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}, nil
}

// MaxCABundleBytes caps how much of a CABundle PEM file is read into memory: 4 MiB.
const MaxCABundleBytes = 4 << 20

func ReadCappedFile(path string, limitBytes int) ([]byte, error) {
	handle, err := os.Open(path) //#nosec G304 -- the path is the caller's own configured path
	if err != nil {
		return nil, err
	}
	defer handle.Close() //nolint:errcheck // read-only handle
	content, err := io.ReadAll(io.LimitReader(handle, int64(limitBytes)+1))
	if err != nil {
		return nil, err
	}
	if len(content) > limitBytes {
		return nil, fmt.Errorf("the file is larger than the %d-byte cap this SDK reads it under",
			limitBytes)
	}
	return content, nil
}

type rawResponse struct {
	status      int
	headers     map[string]string
	contentType string
	retryAfter  string
	body        []byte
}

// SendHTTPRequest sends one request, retried per the policy and decoded by content type.
func (t *HTTPTransport) SendHTTPRequest(ctx context.Context, request TransportRequest) (*ApiResponse, error) {
	method, err := NormaliseMethod(request.HTTPMethod)
	if err != nil {
		return nil, err
	}
	// Preparation runs before the attempt timer, so the operation deadline bounds it; read here and after assembly.
	if err := RefuseIfBudgetSpent(ctx, "before the request body was assembled"); err != nil {
		return nil, err
	}

	request, err = request.Materialised()
	if err != nil {
		return nil, err
	}
	if err := CheckRequestHeaders(request.Headers); err != nil {
		return nil, err
	}
	fullURL, err := buildTransportURL(request.URL, request.QueryParams)
	if err != nil {
		return nil, err
	}
	if err := RefuseIfBudgetSpent(ctx, "while the request body was being assembled"); err != nil {
		return nil, err
	}

	started := time.Now()
	attempt := 0
	requestsSent := 0
	staleEchoResent := false
	var lastError error

	for attempt < t.retry.MaxAttempts {
		attempt++
		requestsSent++
		// Named by endpoint template, not resolved URL: a progress hook reaches stdout and the log alike.
		t.onProgress(fmt.Sprintf("%s %s sending (attempt %d/%d, connect %gs, read %gs)",
			method, wcerrors.EndpointLabel(request.EndpointTemplate, fullURL),
			attempt, t.retry.MaxAttempts,
			t.config.ConnectTimeoutSeconds, t.config.TimeoutSeconds))

		raw, sendError := t.doOnce(ctx, method, fullURL, request)
		if sendError != nil {
			lastError = sendError
			if errors.Is(sendError, errResponseTooLarge) {
				// Returned here: the shared exit below would relabel the cap as a network failure.
				info := t.requestInfo(method, request.URL, request, requestsSent,
					t.grantedAttempts(staleEchoResent), started)
				return nil, &wcerrors.SdkError{
					// The refusal's own sentence only; Error() renders the method and endpoint off the record.
					Message:     sendError.Error(),
					ErrorType:   ErrorTypeUnknown,
					RequestInfo: info,
					Cause:       sendError,
				}
			}
			if attempt >= t.retry.MaxAttempts || !t.retry.ShouldRetryAfterException(sendError, method) {
				break
			}
			delay := t.retry.ComputeBackoffSeconds(attempt, "")
			// Reported before the sleep: an unannounced wait reads as a hung process.
			t.onProgress(fmt.Sprintf("%T; retrying in %.1fs (attempt %d/%d)",
				sendError, delay, attempt+1, t.retry.MaxAttempts))
			if delayError := t.delay(ctx, delay); delayError != nil {
				lastError = delayError
				break
			}
			continue
		}

		if attempt < t.retry.MaxAttempts && t.retry.ShouldRetryStatusCode(raw.status, method) {
			delay := t.retry.ComputeBackoffSeconds(attempt, raw.retryAfter)
			t.onProgress(fmt.Sprintf("%d; retrying in %.1fs (attempt %d/%d)",
				raw.status, delay, attempt+1, t.retry.MaxAttempts))
			if delayError := t.delay(ctx, delay); delayError != nil {
				lastError = delayError
				break
			}
			continue
		}

		decoded := DecodeResponseBody(raw.status, raw.body, raw.contentType)
		response, _ := NewApiResponse(raw.status, decoded, raw.headers, raw.contentType,
			t.requestInfo(method, request.URL, request, requestsSent, t.grantedAttempts(staleEchoResent), started),
			ClassifyErrorType(raw.status))

		response.present(raw.body, raw.contentType, request.ResponseFormat, request.DTOShape)

		if !staleEchoResent && resendableMethods[method] && detectStaleParameterEcho(response, request.QueryParams) {
			// Its own slot, not a retry: a caller who disabled retries must still not get another request's response.
			staleEchoResent = true
			attempt--
			// The resend must leave the keep-alive connection behind, or it inherits the shifted slot.
			t.discardPooledConnections()
			continue
		}
		return response, nil
	}

	// Canceled is tested first: it is the caller withdrawing, where DeadlineExceeded is the SDK's own budget.
	verb := "failed"
	errorType := ErrorTypeNetwork
	switch {
	case errors.Is(lastError, context.Canceled):
		verb = "cancelled"
		errorType = ErrorTypeUnknown
	case isTimeoutError(lastError):
		verb = "timed out"
		errorType = ErrorTypeTimeout
	}
	// The failure's type, not its text, which names the host and port; the text rides CauseText.
	failureType := ""
	causeText := ""
	if lastError != nil {
		failureType = fmt.Sprintf("%T", lastError)
		causeText = redact.RedactText(lastError.Error())
	}
	info := t.requestInfo(method, request.URL, request, requestsSent, t.grantedAttempts(staleEchoResent), started)
	return nil, &wcerrors.SdkError{
		// What went wrong, and nothing about which call it was: Error() renders that off the record.
		Message:     fmt.Sprintf("%s after %d attempt(s): %s", verb, requestsSent, failureType),
		ErrorType:   errorType,
		RequestInfo: info,
		Cause:       lastError,
		CauseText:   causeText,
	}
}

// RefuseIfBudgetSpent folds a request's preparation into the operation deadline the caller set.
func RefuseIfBudgetSpent(ctx context.Context, when string) error {
	cause := ctx.Err()
	if cause == nil {
		return nil
	}
	verb, errorType := "cancelled", ErrorTypeUnknown
	if errors.Is(cause, context.DeadlineExceeded) {
		verb, errorType = "timed out", ErrorTypeTimeout
	}
	return &wcerrors.SdkError{
		Message:   fmt.Sprintf("%s %s, so nothing was sent", verb, when),
		ErrorType: errorType,
		Cause:     cause,
	}
}

// Close releases pooled connections.
func (t *HTTPTransport) Close() {
	t.discardPooledConnections()
}

func (t *HTTPTransport) discardPooledConnections() {
	if tr, ok := t.client.Transport.(*http.Transport); ok {
		tr.CloseIdleConnections()
	}
}

func (t *HTTPTransport) doOnce(ctx context.Context, method, fullURL string, request TransportRequest) (*rawResponse, error) {
	requestContext, cancel := context.WithTimeout(ctx, durationSeconds(t.config.TimeoutSeconds))
	defer cancel()

	httpRequest, err := buildRequest(requestContext, method, fullURL, request)
	if err != nil {
		return nil, err
	}
	httpResponse, err := t.client.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	body, err := readCappedBody(httpResponse.Body, httpResponse.ContentLength, t.config.MaxResponseBytes)
	if err != nil {
		return nil, err
	}
	return &rawResponse{
		status:      httpResponse.StatusCode,
		headers:     flattenHeaders(httpResponse.Header),
		contentType: httpResponse.Header.Get("Content-Type"),
		retryAfter:  httpResponse.Header.Get("Retry-After"),
		body:        body,
	}, nil
}

// errResponseTooLarge marks the one doOnce failure that is not a transport failure, so it is not retried.
var errResponseTooLarge = errors.New("webcommander: response body over the buffered-response cap")

// readCappedBody buffers a response body up to limitBytes, refusing a longer one.
func readCappedBody(source io.Reader, declaredBytes int64, limitBytes int) ([]byte, error) {
	if declaredBytes > int64(limitBytes) {
		return nil, refuseOversizedBody(limitBytes, declaredBytes)
	}
	body, err := io.ReadAll(io.LimitReader(source, int64(limitBytes)+1))
	if err != nil {
		return nil, err
	}
	if len(body) > limitBytes {
		return nil, refuseOversizedBody(limitBytes, declaredBytes)
	}
	return body, nil
}

// refuseOversizedBody wraps errResponseTooLarge so the attempt loop leaves without retrying.
func refuseOversizedBody(limitBytes int, declaredBytes int64) error {
	declared := ""
	if declaredBytes >= 0 {
		declared = fmt.Sprintf(" (Content-Length %d)", declaredBytes)
	}
	return fmt.Errorf("%w: the response body is larger than the %d-byte cap this client "+
		"buffers%s; raise SdkConfig.MaxResponseBytes to accept it",
		errResponseTooLarge, limitBytes, declared)
}

func (t *HTTPTransport) delay(ctx context.Context, seconds float64) error {
	if seconds <= 0 {
		return nil
	}
	timer := time.NewTimer(durationSeconds(seconds))
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (t *HTTPTransport) grantedAttempts(staleEchoResent bool) int {
	if staleEchoResent {
		return t.retry.MaxAttempts + 1
	}
	return t.retry.MaxAttempts
}

func (t *HTTPTransport) requestInfo(method, url string, request TransportRequest, attempt, maxAttempts int, started time.Time) *RequestInfo {
	// Rounded where it is measured, so every call site prints the same digits.
	elapsed := wcerrors.RoundSeconds(time.Since(started).Seconds())
	info := NewRequestInfo(t.config.Auth.ClientName, t.config.APIVersion, method, url, request.QueryParams,
		request.RequestInfoBody, request.RequestInfoHeaders, attempt, maxAttempts, elapsed)
	// Set directly: a template holds no values, so the redacting constructor has nothing to do.
	info.EndpointTemplate = request.EndpointTemplate
	return info
}

// NormaliseMethod upper-cases and trims, and refuses anything outside the RFC 9110 token set.
func NormaliseMethod(httpMethod string) (string, error) {
	if strings.TrimSpace(httpMethod) == "" {
		return "", wcerrors.NewUsageError(fmt.Sprintf(
			"http_method: expected a non-blank string like 'GET', got '%s'", httpMethod))
	}
	method := strings.ToUpper(strings.TrimSpace(httpMethod))
	if !wccfg.IsHTTPToken(method) {
		// Neutralised and redacted, not echoed: a refused verb may carry CR/LF or control characters.
		return "", wcerrors.NewUsageError(redact.RedactText(fmt.Sprintf(
			"httpMethod must be a single verb with no spaces or control characters, got '%s'",
			redact.NeutraliseControls(method))))
	}
	return method, nil
}

func buildTransportURL(rawURL string, queryParams map[string]any) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		// Described, not echoed: an unparseable URL may carry userinfo that key masking cannot reach.
		refusal := wcerrors.NewUsageError(fmt.Sprintf(
			"url: must be an absolute http:// or https:// URL, got %s", redact.Describe(rawURL)))
		refusal.CauseText = redact.RedactText(rawURL)
		return "", refusal
	}
	// Ahead of the early return: a fragment is dropped on the wire whether or not there is a query.
	if err := rejectAFragment("url", rawURL); err != nil {
		return "", err
	}
	if len(queryParams) == 0 {
		return rawURL, nil
	}
	var b strings.Builder
	b.WriteString(rawURL)
	first := !strings.Contains(rawURL, "?")
	for _, key := range sortedMapKeys(queryParams) {
		for _, single := range ExpandValue(queryParams[key]) {
			if first {
				b.WriteByte('?')
			} else {
				b.WriteByte('&')
			}
			b.WriteString(url.QueryEscape(key))
			b.WriteByte('=')
			b.WriteString(url.QueryEscape(single))
			first = false
		}
	}
	return b.String(), nil
}

// ExpandValue renders one query value as the strings the wire carries, a slice as one per element.
func ExpandValue(value any) []string {
	switch v := value.(type) {
	case string:
		return []string{v}
	case []any:
		out := make([]string, 0, len(v))
		for _, item := range v {
			out = append(out, stringifyValue(item))
		}
		return out
	default:
		return []string{stringifyValue(value)}
	}
}

func stringifyValue(value any) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

// CheckRequestHeaders refuses, before any I/O, a header the wire cannot carry.
func CheckRequestHeaders(headers map[string]string) error {
	for _, name := range slices.Sorted(maps.Keys(headers)) {
		problem := wccfg.HeaderNameProblem(name)
		if problem == "" {
			problem = wccfg.HeaderValueProblem(headers[name])
		}
		if problem != "" {
			return wcerrors.NewUsageError(fmt.Sprintf(
				"header '%s': not something the transport can send — %s. The value is not "+
					"shown: it may be a credential", redact.NeutraliseControls(name), problem))
		}
	}
	return nil
}

func buildRequest(ctx context.Context, method, fullURL string, request TransportRequest) (*http.Request, error) {
	assembled := request.wire
	var body io.Reader
	if assembled != nil {
		// A fresh reader over the same bytes, starting at their start on every attempt.
		body = bytes.NewReader(assembled.body)
	}
	httpRequest, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		// net/http's refusal names the whole URL it was handed, so the text rides CauseText.
		refusal := wcerrors.NewUsageError(
			"the request could not be built for this method and URL")
		refusal.Cause = err
		refusal.CauseText = redact.RedactText(err.Error())
		return nil, refusal
	}
	applyHeaders(httpRequest, request.Headers)
	switch {
	case assembled == nil:
		// No body, so nothing to label.
	case assembled.kind == wireMultipart:
		// Unconditional: the boundary is this SDK's, so a caller's Content-Type is overwritten, not rejected.
		httpRequest.Header.Set(contentTypeHeader, assembled.contentType)
	case httpRequest.Header.Get(contentTypeHeader) == "":
		// Set only when the caller has not, so a vendor type survives unchecked.
		httpRequest.Header.Set(contentTypeHeader, assembled.contentType)
	}
	return httpRequest, nil
}

// MaxMultipartBytes caps the multipart body this SDK assembles in memory: 64 MiB, judged before any copy.
const MaxMultipartBytes = 64 << 20

// The boundary's shape is fixed at 55 bytes, so the framing arithmetic knows its length in advance.
const (
	multipartBoundaryPrefix = "webcommanderSdkBoundary"
	multipartBoundaryDigits = 32
	multipartBoundaryLength = len(multipartBoundaryPrefix) + multipartBoundaryDigits
)

// multipartPart is one resolved part of a multipart body, ready to write and already judged.
type multipartPart struct {
	// what names the part in a refusal — "files part 'file'", "formData field 'name'".
	what    string
	header  textproto.MIMEHeader
	content []byte
}

// assembleMultipartBody renders Files and FormData into the exact bytes of a multipart body, once.
func assembleMultipartBody(request TransportRequest) (*wireBody, error) {
	parts, err := resolveMultipartParts(request)
	if err != nil {
		return nil, err
	}
	predicted, err := refuseOversizedMultipart(parts)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	// One allocation of exactly the document's size, where bytes.Buffer would double its way there.
	buffer.Grow(predicted)
	writer := multipart.NewWriter(&buffer)
	boundary, err := multipartBoundary(parts)
	if err != nil {
		return nil, err
	}
	_ = writer.SetBoundary(boundary)
	for _, part := range parts {
		written, err := writer.CreatePart(part.header)
		if err != nil {
			return nil, wcerrors.NewUsageError(part.what + ": could not create part: " + err.Error())
		}
		if _, err := written.Write(part.content); err != nil {
			return nil, wcerrors.NewUsageError(part.what + ": could not be written: " + err.Error())
		}
	}
	if err := writer.Close(); err != nil {
		return nil, wcerrors.NewUsageError("files: could not finalise multipart body: " + err.Error())
	}
	// The produced document, not the predicted one, so the ceiling holds for the bytes on the wire.
	if err := refuseOversizedWireBody("files", "the multipart body", buffer.Len()); err != nil {
		return nil, err
	}
	return &wireBody{kind: wireMultipart, contentType: writer.FormDataContentType(),
		body: buffer.Bytes()}, nil
}

// resolveMultipartParts judges every part in write order: files, then formData, each sorted by name.
func resolveMultipartParts(request TransportRequest) ([]multipartPart, error) {
	parts := make([]multipartPart, 0, len(request.Files)+len(request.FormData))
	for _, field := range sortedMapKeys(request.Files) {
		filename, content, contentType, err := resolveFilePart(field, request.Files[field])
		if err != nil {
			return nil, err
		}
		// CreatePart rather than CreateFormFile, which hardcodes application/octet-stream.
		header := textproto.MIMEHeader{}
		header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeMultipartParam(field), escapeMultipartParam(filename)))
		if contentType != "" {
			header.Set(contentTypeHeader, contentType)
		}
		parts = append(parts, multipartPart{
			what:    fmt.Sprintf("files: '%s'", field),
			header:  header,
			content: content,
		})
	}
	// Not writer.WriteField: its quote escaper writes CR and LF in a field name verbatim.
	for _, name := range sortedMapKeys(request.FormData) {
		if err := refuseMultipartHeaderBytes(formDataLabel(name), name); err != nil {
			return nil, err
		}
		header := textproto.MIMEHeader{}
		header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"`,
			escapeMultipartParam(name)))
		parts = append(parts, multipartPart{
			what:    "formData",
			header:  header,
			content: []byte(stringifyValue(request.FormData[name])),
		})
	}
	return parts, nil
}

// sortedMapKeys sorts, so a body is one document rather than one per range over the mapping.
func sortedMapKeys(source map[string]any) []string {
	return slices.Sorted(maps.Keys(source))
}

// multipartBoundary derives the boundary from the parts, re-salting until it appears in none of them.
func multipartBoundary(parts []multipartPart) (string, error) {
	sum := multipartPartsDigest(parts)
	for salt := 0; salt < maxMultipartBoundarySalts; salt++ {
		candidate := multipartBoundaryCandidate(sum, salt)
		if !multipartPartsContain(parts, candidate) {
			return candidate, nil
		}
	}
	// Unreachable short of a preimage break: each candidate is 128 bits derived from the whole body.
	return "", wcerrors.NewUsageError(fmt.Sprintf(
		"files: could not derive a multipart boundary absent from the body after %d attempts",
		maxMultipartBoundarySalts))
}

const maxMultipartBoundarySalts = 64

// multipartPartsDigest mixes lengths in, so two different bodies cannot stream the same bytes in.
func multipartPartsDigest(parts []multipartPart) []byte {
	digest := sha256.New()
	for _, part := range parts {
		for _, name := range slices.Sorted(maps.Keys(part.header)) {
			fmt.Fprintf(digest, "%s\x00%s\x00", name, strings.Join(part.header[name], "\x00"))
		}
		fmt.Fprintf(digest, "%d\x00", len(part.content))
		digest.Write(part.content)
	}
	return digest.Sum(nil)
}

// multipartBoundaryCandidate is the salt'th boundary for parts digesting to sum; salt 0 is the digest.
func multipartBoundaryCandidate(sum []byte, salt int) string {
	digits := hex.EncodeToString(sum)
	if salt > 0 {
		resalted := sha256.New()
		resalted.Write(sum)
		fmt.Fprintf(resalted, "\x00salt\x00%d", salt)
		digits = hex.EncodeToString(resalted.Sum(nil))
	}
	return multipartBoundaryPrefix + digits[:multipartBoundaryDigits]
}

// multipartPartsContain searches for the bare boundary, which a lenient receiver would also match.
func multipartPartsContain(parts []multipartPart, candidate string) bool {
	needle := []byte(candidate)
	for _, part := range parts {
		for name, values := range part.header {
			if strings.Contains(name, candidate) {
				return true
			}
			for _, value := range values {
				if strings.Contains(value, candidate) {
					return true
				}
			}
		}
		if bytes.Contains(part.content, needle) {
			return true
		}
	}
	return false
}

// multipartWireSize is the exact byte length assembleMultipartBody will write, per mime/multipart's framing.
func multipartWireSize(parts []multipartPart, limit int) (int, bool) {
	total := multipartBoundaryLength + 8 // the closing delimiter line
	for index, part := range parts {
		if index == 0 {
			total += multipartBoundaryLength + 4
		} else {
			total += multipartBoundaryLength + 6
		}
		for name, values := range part.header {
			for _, value := range values {
				total += len(name) + 2 + len(value) + 2 // "name: value" CRLF
			}
		}
		total += 2 // the blank line that ends the header block
		total += len(part.content)
		if total > limit {
			return 0, false
		}
	}
	return total, true
}

// refuseOversizedMultipart refuses a body over MaxMultipartBytes before a byte is copied.
func refuseOversizedMultipart(parts []multipartPart) (int, error) {
	size, within := multipartWireSize(parts, MaxMultipartBytes)
	if !within {
		return 0, wcerrors.NewUsageError(fmt.Sprintf(
			"files: multipart body over the %d-byte cap", MaxMultipartBytes))
	}
	return size, nil
}

// CheckFiles refuses, before any I/O and the token grant, files entries that cannot become parts.
func CheckFiles(files map[string]any) error {
	for _, field := range sortedMapKeys(files) {
		if _, _, _, err := resolveFilePart(field, files[field]); err != nil {
			return err
		}
	}
	return nil
}

// CheckFormDataNames refuses, before the token grant, a formData name that cannot label a part.
func CheckFormDataNames(formData map[string]any) error {
	for _, name := range sortedMapKeys(formData) {
		if err := refuseMultipartHeaderBytes(formDataLabel(name), name); err != nil {
			return err
		}
	}
	return nil
}

// formDataLabel names a field by its length, not its text, which would carry the forged header.
func formDataLabel(name string) string {
	return fmt.Sprintf("a formData name of %d byte(s)", len(name))
}

func resolveFilePart(field string, value any) (filename string, content []byte, contentType string, err error) {
	resolved := false
	switch v := value.(type) {
	case []byte, string:
		if payload, ok := partContent(v); ok {
			filename, content, resolved = field, payload, true
		}
	case []any:
		// The triple lets the caller set the part's content type, as an .xlsx import needs.
		if len(v) > 3 {
			return "", nil, "", wcerrors.NewUsageError(fmt.Sprintf(
				"files: '%s' has %d elements; the tuple form is [filename, contents] or "+
					"[filename, contents, contentType]", field, len(v)))
		}
		if len(v) >= 2 {
			if name, ok := v[0].(string); ok {
				if bytesValue, contentOK := partContent(v[1]); contentOK {
					filename, content, resolved = name, bytesValue, true
				}
			}
		}
		if len(v) == 3 && resolved {
			ct, ok := v[2].(string)
			if !ok {
				return "", nil, "", wcerrors.NewUsageError(fmt.Sprintf(
					"files: '%s' third element must be the part's content type as a string, got %T",
					field, v[2]))
			}
			contentType = ct
		}
	}
	if !resolved {
		return "", nil, "", wcerrors.NewUsageError(fmt.Sprintf(
			"files: '%s' must be []byte, a string, or [filename, contents[, contentType]]", field))
	}
	// All three reach a part header, so all three are held to the same rule.
	if err := refuseMultipartHeaderBytes(fmt.Sprintf("a files part name of %d byte(s)", len(field)), field); err != nil {
		return "", nil, "", err
	}
	for _, held := range []struct{ what, value string }{
		{fmt.Sprintf("the filename on part '%s'", field), filename},
		{fmt.Sprintf("the content type on part '%s'", field), contentType},
	} {
		if err := refuseMultipartHeaderBytes(held.what, held.value); err != nil {
			return "", nil, "", err
		}
	}
	return filename, content, contentType, nil
}

func partContent(value any) ([]byte, bool) {
	switch v := value.(type) {
	case []byte:
		return v, v != nil
	case string:
		return []byte(v), true
	}
	return nil, false
}

// multipartParamEscaper percent-escapes the two bytes that can end a quoted parameter early.
var multipartParamEscaper = strings.NewReplacer(`"`, "%22", `\`, "%5C")

func escapeMultipartParam(value string) string {
	return multipartParamEscaper.Replace(value)
}

// refuseMultipartHeaderBytes refuses CR, LF or NUL in a part name, upload filename or content type.
func refuseMultipartHeaderBytes(what, value string) error {
	for index := 0; index < len(value); index++ {
		name := ""
		switch value[index] {
		case '\r':
			name = "a carriage return"
		case '\n':
			name = "a newline"
		case 0x00:
			name = "a NUL"
		}
		if name != "" {
			return wcerrors.NewUsageError(fmt.Sprintf("%s carries %s at byte %d", what, name, index))
		}
	}
	return nil
}

func applyHeaders(httpRequest *http.Request, headers map[string]string) {
	for k, v := range headers {
		httpRequest.Header.Set(k, v)
	}
}

func flattenHeaders(header http.Header) map[string]string {
	out := make(map[string]string, len(header))
	for k, v := range header {
		if len(v) > 0 {
			// Joined, never dropped, so a second Set-Cookie survives.
			out[k] = strings.Join(v, ", ")
		}
	}
	return out
}

// detectStaleParameterEcho detects a 400 that belongs to another request on the same keep-alive connection.
func detectStaleParameterEcho(response *ApiResponse, sentParams map[string]any) bool {
	if response.StatusCode != 400 {
		return false
	}
	body, ok := response.Map().(map[string]any)
	if !ok {
		return false
	}
	namedFields := extractErrorFields(body)
	if len(namedFields) == 0 && stringifyValue(body["error"]) == unknownParameterError {
		namedFields = parseUnknownParameterMessage(body)
	}
	if len(namedFields) == 0 {
		return false
	}
	sent := map[string]bool{}
	for name := range sentParams {
		sent[strings.ToLower(name)] = true
	}
	if len(sent) == 0 {
		return false // with nothing sent, "not ours" is vacuous
	}
	for field := range namedFields {
		if sent[field] {
			return false // an intersection means it is genuinely ours
		}
	}
	return true
}

func extractErrorFields(body map[string]any) map[string]bool {
	fields := map[string]bool{}
	switch entries := body["errors"].(type) {
	case map[string]any:
		// Keyed by field name, so the keys are the objected-to fields.
		for fieldName := range entries {
			if name := strings.ToLower(strings.TrimSpace(fieldName)); name != "" {
				fields[name] = true
			}
		}
	case []any:
		for _, entry := range entries {
			if m, ok := entry.(map[string]any); ok && m["field"] != nil {
				fields[strings.ToLower(strings.TrimSpace(fmt.Sprint(m["field"])))] = true
			}
		}
	}
	return fields
}

func parseUnknownParameterMessage(body map[string]any) map[string]bool {
	fields := map[string]bool{}
	message := stringifyValue(body["message"])
	const marker = "unknown query parameter(s):"
	_, listing, found := strings.Cut(message, marker)
	if !found {
		return fields
	}
	for listed := range strings.SplitSeq(listing, ",") {
		if named := parameterNameRe.FindString(listed); named != "" {
			fields[strings.ToLower(named)] = true
		}
	}
	return fields
}

func durationSeconds(seconds float64) time.Duration {
	return time.Duration(seconds * float64(time.Second))
}
