package client

import (
	"context"
	"fmt"
	"iter"
	"maps"
	neturl "net/url"
	"slices"
	"sort"
	"strings"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

type PaginateOptions struct {
	// QueryParams apply to the first request only.
	QueryParams map[string]any
	// MaxPages caps the walk; 0 means unlimited.
	MaxPages int
	// RaiseOnError: nil or true stops with an error on a failed page; false yields that page last.
	RaiseOnError *bool
}

// IterPages walks a paginated listing, one page per yield, following pagination.next_page.
func IterPages(ctx context.Context, client http.Requester, path string, options PaginateOptions) iter.Seq2[*http.ApiResponse, error] {
	return func(yield func(*http.ApiResponse, error) bool) {
		if options.MaxPages < 0 {
			yield(nil, wcerrors.NewUsageError(fmt.Sprintf("maxPages: expected a positive integer or 0 for unlimited, got %d", options.MaxPages)))
			return
		}
		raise := options.RaiseOnError == nil || *options.RaiseOnError
		firstQuery, err := http.BuildQueryParams(options.QueryParams)
		if err != nil {
			yield(nil, err)
			return
		}

		visited := map[string]bool{}
		yielded := 0
		target, targetQuery := path, firstQuery

		for {
			key := canonicalPageKey(client, target, targetQuery)
			if visited[key] {
				yield(nil, &wcerrors.SdkError{Message: fmt.Sprintf(
					"pagination loop: a page was served twice after %d page(s) of %s",
					yielded, wcerrors.EndpointLabel("", path))})
				return
			}
			visited[key] = true

			response, err := client.SendGetRequest(ctx, target,
				&http.RequestOptions{QueryParams: targetQuery, RaiseOnError: &raise})
			if err != nil {
				yield(nil, err)
				return
			}
			if !yield(response, nil) {
				return
			}
			yielded++

			if !response.Succeeded() {
				return // only reachable with RaiseOnError=false
			}
			if options.MaxPages > 0 && yielded >= options.MaxPages {
				return
			}
			// Map(), not Data: next_page is a fact about the body, whatever Data presents.
			next := readNextPage(response.Map())
			if next == "" {
				return
			}
			// A next link carries its own query; the caller's params rode the first request only.
			target, targetQuery = next, nil
		}
	}
}

// requestURLResolver is the optional half of http.Requester, asserted for below.
type requestURLResolver interface {
	resolvedRequestURL(path string) (string, error)
}

// resolvedRequestURL assembles a path as a request would, so a self-link is recognisable as one.
func (c *WebCommanderClient) resolvedRequestURL(path string) (string, error) {
	return http.BuildURL(c.config.Auth.BaseURL, path)
}

func canonicalPageKey(client http.Requester, target string, extraQuery map[string]any) string {
	resolved := target
	if resolver, ok := client.(requestURLResolver); ok {
		if absolute, err := resolver.resolvedRequestURL(target); err == nil {
			resolved = absolute
		}
	}
	parsed, err := neturl.Parse(resolved)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "unresolved\x00" + resolved + "\x00" + canonicalQuery(nil, extraQuery)
	}
	var b strings.Builder
	b.WriteString(http.OriginOf(parsed))
	b.WriteString(canonicalPath(parsed.EscapedPath()))
	if query := canonicalQuery(parsed, extraQuery); query != "" {
		b.WriteByte('?')
		b.WriteString(query)
	}
	// The fragment is client-side only (RFC 3986 3.5), so links differing in nothing else are one request.
	return b.String()
}

func canonicalPath(escapedPath string) string {
	if escapedPath == "" {
		return "/"
	}
	segments := strings.Split(escapedPath, "/")
	out := make([]string, 0, len(segments))
	for index, segment := range segments {
		switch segment {
		case ".":
			// A dot segment names the directory it is in.
			if index == len(segments)-1 {
				out = append(out, "")
			}
			continue
		case "..":
			// Pop the last named segment, keeping the leading empty one that makes the path absolute.
			for len(out) > 1 && out[len(out)-1] == "" {
				out = out[:len(out)-1]
			}
			if len(out) > 1 {
				out = out[:len(out)-1]
			}
			if index == len(segments)-1 {
				out = append(out, "")
			}
			continue
		}
		out = append(out, recodeComponent(segment))
	}
	joined := strings.Join(out, "/")
	if !strings.HasPrefix(joined, "/") {
		return "/" + joined
	}
	return joined
}

func canonicalQuery(parsed *neturl.URL, extraQuery map[string]any) string {
	var pairs []string
	if parsed != nil {
		for _, chunk := range strings.Split(parsed.RawQuery, "&") {
			if chunk == "" {
				continue
			}
			name, value, found := strings.Cut(chunk, "=")
			if !found {
				pairs = append(pairs, recodeComponent(name))
				continue
			}
			pairs = append(pairs, recodeComponent(name)+"="+recodeComponent(value))
		}
	}
	for _, name := range slices.Sorted(maps.Keys(extraQuery)) {
		for _, single := range http.ExpandValue(extraQuery[name]) {
			pairs = append(pairs, neturl.QueryEscape(name)+"="+neturl.QueryEscape(single))
		}
	}
	sort.Strings(pairs)
	return strings.Join(pairs, "&")
}

// recodeComponent re-encodes one escaped segment under one rule: "%7E" and "~" match, "%2F" and "/" do not.
func recodeComponent(component string) string {
	decoded, err := neturl.QueryUnescape(component)
	if err != nil {
		return component
	}
	return neturl.QueryEscape(decoded)
}

// IterRecords walks every record of a listing across every page; a failed page always errors and stops.
func IterRecords(ctx context.Context, client http.Requester, path, recordKey string, options PaginateOptions) iter.Seq2[any, error] {
	return func(yield func(any, error) bool) {
		if strings.TrimSpace(recordKey) == "" {
			yield(nil, wcerrors.NewUsageError(fmt.Sprintf(
				"recordKey: expected the envelope key holding the records, like 'customers', got '%s'", recordKey)))
			return
		}
		alwaysRaise := true
		options.RaiseOnError = &alwaysRaise
		for page, err := range IterPages(ctx, client, path, options) {
			if err != nil {
				yield(nil, err)
				return
			}
			records, keyMissing := recordsOf(page, recordKey)
			if keyMissing {
				// A page with no list under the key is malformed, not the end: later pages come only via next_page.
				yield(nil, missingRecordKeyError(page, recordKey))
				return
			}
			for _, record := range records {
				if !yield(record, nil) {
					return
				}
			}
		}
	}
}

// readNextPage returns the next page's URL, or "" at the end of the listing.
func readNextPage(data any) string {
	m, ok := data.(map[string]any)
	if !ok {
		return ""
	}
	pagination, ok := m["pagination"].(map[string]any)
	if !ok {
		return ""
	}
	next, ok := pagination["next_page"].(string)
	if !ok || strings.TrimSpace(next) == "" {
		return ""
	}
	return next
}

func recordsOf(page *http.ApiResponse, recordKey string) (records []any, keyMissing bool) {
	if data, ok := page.Map().(map[string]any); ok {
		if list, ok := data[recordKey].([]any); ok {
			return list, false
		}
	}
	return nil, true
}

// missingRecordKeyError names the keys the page does carry, so a wrong key diagnoses itself.
func missingRecordKeyError(page *http.ApiResponse, recordKey string) error {
	// The envelope's own type, so a list, a string and a number each diagnose differently.
	available := fmt.Sprintf("a %T body", page.Map())
	if data, ok := page.Map().(map[string]any); ok {
		available = strings.Join(slices.Sorted(maps.Keys(data)), ", ")
	}
	return wcerrors.NewUsageError(fmt.Sprintf(
		"recordKey '%s': this page has no list under that key; it carries %s", recordKey, available))
}
