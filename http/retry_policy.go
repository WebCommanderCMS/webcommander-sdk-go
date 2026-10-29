package http

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"maps"
	"math"
	"math/rand"
	"net"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	wccfg "github.com/WebCommanderCMS/webcommander-sdk-go/v4/config"
	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
)

type RetryPolicy struct {
	// MaxAttempts is total attempts including the first, so 1 disables retrying.
	MaxAttempts int
	// BackoffSeconds is the base backoff before the first retry.
	BackoffSeconds float64
	// BackoffMaxSeconds is the ceiling on any single backoff.
	BackoffMaxSeconds float64

	retryable       map[int]bool
	unsafeRetryable map[int]bool
}

// safeMethods is RFC 9110 section 9.2.1's safe set, governing the status gates only.
var safeMethods = map[string]bool{"GET": true, "HEAD": true, "OPTIONS": true, "TRACE": true}

// idempotentMethods is the safe set plus PUT and DELETE, governing transport failures only.
var idempotentMethods = map[string]bool{
	"GET": true, "HEAD": true, "OPTIONS": true, "TRACE": true, "PUT": true, "DELETE": true,
}

// RetryableStatusCodes returns a fresh sorted copy of the statuses this policy retries.
func (p *RetryPolicy) RetryableStatusCodes() []int { return sortedCodes(p.retryable) }

// UnsafeRetryableStatusCodes returns a fresh sorted copy of the statuses also retried on an unsafe method; ShouldRetryStatusCode retries one only if it is in RetryableStatusCodes as well.
func (p *RetryPolicy) UnsafeRetryableStatusCodes() []int { return sortedCodes(p.unsafeRetryable) }

func sortedCodes(set map[int]bool) []int {
	return slices.Sorted(maps.Keys(set))
}

// NewRetryPolicyFromConfig builds a retry policy from a config's retry fields.
func NewRetryPolicyFromConfig(config wccfg.SdkConfig) (*RetryPolicy, error) {
	return NewRetryPolicy(config.MaxRetries, config.RetryBackoffSecondsValue(), config.RetryBackoffMaxSeconds,
		config.RetryableHTTPCodes, config.UnsafeRetryableHTTPCodes)
}

// NewRetryPolicy constructs a policy, validating the backoff bounds; a nil code list takes the default.
func NewRetryPolicy(maxAttempts int, backoffSeconds, backoffMaxSeconds float64,
	retryable, unsafeRetryable []int) (*RetryPolicy, error) {
	if maxAttempts < 1 {
		return nil, wcerrors.NewConfigError(fmt.Sprintf("maxAttempts: must be at least 1 (1 disables retry), got %d", maxAttempts))
	}
	if wccfg.IsNaNOrInf(backoffSeconds) {
		return nil, wcerrors.NewConfigError(fmt.Sprintf("backoff_seconds: expected a number of seconds, got %s", wcerrors.FormatSeconds(backoffSeconds)))
	}
	if wccfg.IsNaNOrInf(backoffMaxSeconds) {
		return nil, wcerrors.NewConfigError(fmt.Sprintf("backoff_max_seconds: expected a number of seconds, got %s", wcerrors.FormatSeconds(backoffMaxSeconds)))
	}
	if backoffSeconds < 0 {
		return nil, wcerrors.NewConfigError(fmt.Sprintf("backoff_seconds: must be >= 0, got %s", wcerrors.FormatSeconds(backoffSeconds)))
	}
	if backoffMaxSeconds < backoffSeconds {
		return nil, wcerrors.NewConfigError(fmt.Sprintf("backoff_max_seconds: must be >= backoff_seconds, got %s < %s",
			wcerrors.FormatSeconds(backoffMaxSeconds), wcerrors.FormatSeconds(backoffSeconds)))
	}
	// The same range every SdkConfig field is held to, in the same words.
	if problems := append(wccfg.HTTPCodeSetProblems("retryable_http_codes", retryable),
		wccfg.HTTPCodeSetProblems("unsafe_retryable_http_codes", unsafeRetryable)...); len(problems) > 0 {
		return nil, wcerrors.NewConfigError(strings.Join(problems, "; "))
	}
	if retryable == nil {
		retryable = wccfg.DefaultRetryableCodes()
	}
	if unsafeRetryable == nil {
		unsafeRetryable = wccfg.DefaultUnsafeRetryableCodes()
	}
	return &RetryPolicy{
		MaxAttempts:       maxAttempts,
		BackoffSeconds:    backoffSeconds,
		BackoffMaxSeconds: backoffMaxSeconds,
		retryable:         intSet(retryable),
		unsafeRetryable:   intSet(unsafeRetryable),
	}, nil
}

// frozen is this policy's private copy, taken at construction and re-validated on the way through.
func (p *RetryPolicy) frozen() (*RetryPolicy, error) {
	return NewRetryPolicy(p.MaxAttempts, p.BackoffSeconds, p.BackoffMaxSeconds,
		p.RetryableStatusCodes(), p.UnsafeRetryableStatusCodes())
}

func (p *RetryPolicy) ShouldRetryStatusCode(statusCode int, method string) bool {
	if !p.retryable[statusCode] {
		return false
	}
	if IsUnsafeMethod(method) {
		return p.unsafeRetryable[statusCode]
	}
	return true
}

// ShouldRetryAfterException retries a never-delivered failure on any method, an unknown one if idempotent.
func (p *RetryPolicy) ShouldRetryAfterException(err error, method string) bool {
	switch classifyDelivery(err) {
	case proofNeverDelivered:
		return true
	case proofDeliveryUnknown:
		return IsIdempotentMethod(method)
	default:
		return false
	}
}

// maxBackoffDoublings clamps the exponent: 2^1024 is +Inf, and a 0 backoff would then make a NaN window.
const maxBackoffDoublings = 1023

// ComputeBackoffSeconds is exponential with full jitter; a numeric Retry-After wins, clamped the same.
func (p *RetryPolicy) ComputeBackoffSeconds(attempt int, retryAfter string) float64 {
	if serverDelay, ok := parseRetryAfter(retryAfter); ok {
		return math.Min(serverDelay, p.BackoffMaxSeconds)
	}
	doublings := math.Min(math.Max(0, float64(attempt-1)), maxBackoffDoublings)
	window := math.Min(p.BackoffSeconds*math.Pow(2, doublings), p.BackoffMaxSeconds)
	if window <= 0 {
		return 0
	}
	return rand.Float64() * window //#nosec G404 -- retry jitter spreads load; it is not a secret
}

// IsUnsafeMethod negates the safe allow-list, so an extension method or the empty string is unsafe.
func IsUnsafeMethod(method string) bool {
	return !safeMethods[strings.ToUpper(strings.TrimSpace(method))]
}

// IsIdempotentMethod reports whether a method may be replayed after a failure leaving delivery unknown.
func IsIdempotentMethod(method string) bool {
	return idempotentMethods[strings.ToUpper(strings.TrimSpace(method))]
}

// deliveryProof is what a transport failure proves about delivery: the SDK's one retry-safety rule.
type deliveryProof int

const (
	// proofNoRetry is a failure retry cannot address: a cancellation, or no transport failure at all.
	proofNoRetry deliveryProof = iota
	// proofNeverDelivered is a failure to establish the connection, so nothing can have arrived.
	proofNeverDelivered
	// proofDeliveryUnknown is a failure after the request was fully sent, so the store may have acted.
	proofDeliveryUnknown
)

func classifyDelivery(err error) deliveryProof {
	if !isTransportError(err) {
		return proofNoRetry
	}
	if isConnectionEstablishmentError(err) {
		return proofNeverDelivered
	}
	return proofDeliveryUnknown
}

// isConnectionEstablishmentError reports whether a failure happened before the request was written.
func isConnectionEstablishmentError(err error) bool {
	var dnsError *net.DNSError
	if errors.As(err, &dnsError) {
		return true
	}
	var operationError *net.OpError
	if errors.As(err, &operationError) && operationError.Op == "dial" {
		return true
	}
	return isTLSError(err)
}

// The date grammars a Retry-After may arrive in, applied after the day-of-week is removed.
var httpDateLayouts = []string{
	// RFC 5322 / IMF-fixdate, with and without seconds, zone named or numeric or absent.
	"2 Jan 2006 15:04:05 MST",
	"2 Jan 2006 15:04:05 -0700",
	"2 Jan 2006 15:04:05",
	"2 Jan 2006 15:04 MST",
	"2 Jan 2006 15:04 -0700",
	"2 Jan 2006 15:04",
	// RFC 850's dashed date, two- and four-digit years.
	"2-Jan-06 15:04:05 MST",
	"2-Jan-06 15:04:05 -0700",
	"2-Jan-06 15:04:05",
	"2-Jan-2006 15:04:05 MST",
	"2-Jan-2006 15:04:05 -0700",
	"2-Jan-2006 15:04:05",
	// asctime, space-padded day or not.
	"Jan _2 15:04:05 2006",
	"Jan 2 15:04:05 2006",
}

// Removed, not parsed: a server sending the wrong weekday still gets its wait honoured.
var dayNames = map[string]bool{
	"mon": true, "tue": true, "wed": true, "thu": true, "fri": true, "sat": true, "sun": true,
	"monday": true, "tuesday": true, "wednesday": true, "thursday": true, "friday": true,
	"saturday": true, "sunday": true,
}

// stripDayName leaves a token that is not a day name in place, so "21 Oct 2099 ..." keeps its day.
func stripDayName(value string) string {
	head, rest, found := strings.Cut(value, " ")
	if !found {
		return value
	}
	if !dayNames[strings.ToLower(strings.TrimSuffix(head, ","))] {
		return value
	}
	return strings.TrimSpace(rest)
}

// parseRetryAfter reads delay-seconds then an HTTP-date; false leaves the caller on its computed jitter.
func parseRetryAfter(retryAfter string) (float64, bool) {
	retryAfter = strings.TrimSpace(retryAfter)
	if retryAfter == "" {
		return 0, false
	}
	if isNumericCandidate(retryAfter) {
		if seconds, err := strconv.ParseFloat(retryAfter, 64); err == nil {
			return clampSeconds(seconds), true
		}
	}
	dated := stripDayName(retryAfter)
	for _, layout := range httpDateLayouts {
		if when, err := time.Parse(layout, dated); err == nil {
			// A zone-less spelling parses as UTC; a date already past clamps to 0.
			return clampSeconds(time.Until(resolveNamedZone(when)).Seconds()), true
		}
	}
	return 0, false
}

// The North-American zone abbreviations RFC 5322 still admits, in seconds east of UTC.
var namedZoneOffsets = map[string]int{
	"EST": -5 * 3600, "EDT": -4 * 3600,
	"CST": -6 * 3600, "CDT": -5 * 3600,
	"MST": -7 * 3600, "MDT": -6 * 3600,
	"PST": -8 * 3600, "PDT": -7 * 3600,
	"AST": -4 * 3600, "ADT": -3 * 3600,
}

// resolveNamedZone corrects a zone time.Parse left at a fabricated offset 0, so EST is not read as UTC.
func resolveNamedZone(when time.Time) time.Time {
	name, offset := when.Zone()
	if offset != 0 {
		return when // the runtime already resolved the zone
	}
	tableOffset, ok := namedZoneOffsets[name]
	if !ok {
		return when
	}
	// The wall clock was read as UTC; the real instant is that clock minus the zone's offset.
	return when.Add(-time.Duration(tableOffset) * time.Second)
}

// isNumericCandidate excludes the underscore and hexadecimal spellings delta-seconds disallows.
func isNumericCandidate(value string) bool {
	if strings.Contains(value, "_") {
		return false
	}
	return !strings.Contains(strings.ToLower(value), "0x")
}

// clampSeconds floors a parsed delay at zero, answering 0 for NaN rather than propagating it.
func clampSeconds(seconds float64) float64 {
	if math.IsNaN(seconds) {
		return 0
	}
	return math.Max(0, seconds)
}

// isTLSError reports a certificate or handshake problem.
func isTLSError(err error) bool {
	if err == nil {
		return false
	}
	var unknownAuthority x509.UnknownAuthorityError
	var invalidCert x509.CertificateInvalidError
	var hostname x509.HostnameError
	var recordHeader tls.RecordHeaderError
	var verification *tls.CertificateVerificationError
	return errors.As(err, &unknownAuthority) || errors.As(err, &invalidCert) ||
		errors.As(err, &hostname) || errors.As(err, &recordHeader) ||
		errors.As(err, &verification)
}

func isTransportError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) {
		// The caller withdrew the request; another attempt is not what cancelling asked for.
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	// io.ErrUnexpectedEOF carries neither net.Error nor *url.Error, so it is named here to reach the table.
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return true
	}
	var networkError net.Error
	if errors.As(err, &networkError) {
		return true
	}
	var urlError *url.Error
	return errors.As(err, &urlError)
}

func isTimeoutError(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var networkError net.Error
	return errors.As(err, &networkError) && networkError.Timeout()
}

func intSet(codes []int) map[int]bool {
	set := make(map[int]bool, len(codes))
	for _, c := range codes {
		set[c] = true
	}
	return set
}
