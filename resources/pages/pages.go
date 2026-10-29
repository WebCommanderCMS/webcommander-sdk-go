package pages

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// Pages is the page endpoints, bound to one client.
type Pages struct {
	requester http.Requester
}

// ListPages lists the pages outside the trash, newest change first — GET /admin/pages; search, active, has_layout, visibility, created_from/created_to, updated_from/updated_to (YYYY-MM-DD), author, store and not_in_store filter, sort and order order, limit is 1 to 100 or "all"; an unknown name is 400.
func (p *Pages) ListPages(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendGetRequest(ctx, PagesPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: PageListingShape})
}

// HeadPages asks whether the collection is reachable — HEAD /admin/pages; no body.
func (p *Pages) HeadPages(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendApiRequest(ctx, "HEAD", PagesPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CountPages is how many pages match — GET /admin/pages/count, a bare {"total": N}; the listing's filters narrow it, and sort, order, limit and offset are ignored.
func (p *Pages) CountPages(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendGetRequest(ctx, PagesCountPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: PageCountShape})
}

// HeadPagesCount asks whether the count is reachable — HEAD /admin/pages/count; no body.
func (p *Pages) HeadPagesCount(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendApiRequest(ctx, "HEAD", PagesCountPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CreatePage creates a page — POST /admin/pages with {"page": {…}}, answered 201; name and visibility are required, and the URL is made from the name.
func (p *Pages) CreatePage(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendPostRequest(ctx, PagesPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: PageShape})
}

// GetPage retrieves one page, with an ETag — GET /admin/pages/{page_id}; a page in the trash still reads, in_trash true.
func (p *Pages) GetPage(ctx context.Context, pageID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendGetRequest(ctx, PagePath, &http.RequestOptions{
		PathParams:   map[string]any{"page_id": pageID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: PageShape})
}

// HeadPage asks whether one page is there — HEAD /admin/pages/{page_id}; no body, the ETag its GET carries.
func (p *Pages) HeadPage(ctx context.Context, pageID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendApiRequest(ctx, "HEAD", PagePath, &http.RequestOptions{
		PathParams:   map[string]any{"page_id": pageID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// UpdatePage replaces one page — PUT /admin/pages/{page_id}; name, title, subtitle, url, layout, visibility, body_class, active, disable_tracking and password_protection are required, a partial body is 400. ifMatch is sent as If-Match, a stale ETag is 412 and is checked before the body; empty sends none.
func (p *Pages) UpdatePage(ctx context.Context, pageID any, payload any, ifMatch string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	var headers map[string]string
	if ifMatch != "" {
		headers = map[string]string{PagesIfMatchHeader: ifMatch}
	}
	return p.requester.SendPutRequest(ctx, PagePath, &http.RequestOptions{
		PathParams: map[string]any{"page_id": pageID}, JSONBody: payload, Headers: headers,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: PageShape})
}

// PatchPage changes only the fields sent — PATCH /admin/pages/{page_id}; an empty page is 400 and a page in the trash 409. ifMatch is sent as If-Match; empty sends none.
func (p *Pages) PatchPage(ctx context.Context, pageID any, payload any, ifMatch string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	var headers map[string]string
	if ifMatch != "" {
		headers = map[string]string{PagesIfMatchHeader: ifMatch}
	}
	return p.requester.SendPatchRequest(ctx, PagePath, &http.RequestOptions{
		PathParams: map[string]any{"page_id": pageID}, JSONBody: payload, Headers: headers,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: PageShape})
}

// DeletePage moves a page to the trash — DELETE /admin/pages/{page_id}; 204. force=true purges it, from the trash or live; a second plain delete is 400 and the landing page 409. ifMatch is sent as If-Match; empty sends none.
func (p *Pages) DeletePage(ctx context.Context, pageID any, queryParams map[string]any, ifMatch string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	var headers map[string]string
	if ifMatch != "" {
		headers = map[string]string{PagesIfMatchHeader: ifMatch}
	}
	return p.requester.SendDeleteRequest(ctx, PagePath, &http.RequestOptions{
		PathParams:  map[string]any{"page_id": pageID},
		QueryParams: queryParams, Headers: headers,
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// UpdatePageURL moves a page to a new URL — PUT /admin/pages/{page_id}/url with a bare {"url": …}; answers the page, and a URL another page holds is 409.
func (p *Pages) UpdatePageURL(ctx context.Context, pageID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendPutRequest(ctx, PageURLPath, &http.RequestOptions{
		PathParams: map[string]any{"page_id": pageID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: PageShape})
}

// SetLandingPage makes a page the store's landing page — PATCH /admin/pages/{page_id}/landing with no body; the previous landing page stops being it.
func (p *Pages) SetLandingPage(ctx context.Context, pageID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendPatchRequest(ctx, PageLandingPath, &http.RequestOptions{
		PathParams:   map[string]any{"page_id": pageID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: PageShape})
}

// GetPageDeleteBlockers is whether a delete would pass and what blocks it, such as landing_page — GET /admin/pages/{page_id}/delete-blockers.
func (p *Pages) GetPageDeleteBlockers(ctx context.Context, pageID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendGetRequest(ctx, PageDeleteBlockersPath, &http.RequestOptions{
		PathParams:   map[string]any{"page_id": pageID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: PageDeleteBlockersShape})
}

// HeadPageDeleteBlockers asks whether the delete check is reachable — HEAD /admin/pages/{page_id}/delete-blockers; no body.
func (p *Pages) HeadPageDeleteBlockers(ctx context.Context, pageID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendApiRequest(ctx, "HEAD", PageDeleteBlockersPath, &http.RequestOptions{
		PathParams:   map[string]any{"page_id": pageID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// DuplicatePage copies a page — POST /admin/pages/{page_id}/duplicate with no body, answered 201; a body field is 400. The copy is "<name> - Copy" at <url>-copy, then "- Copy 1" at -copy-1.
func (p *Pages) DuplicatePage(ctx context.Context, pageID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendPostRequest(ctx, PageDuplicatePath, &http.RequestOptions{
		PathParams:   map[string]any{"page_id": pageID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: PageShape})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "pages"

func oneResponseFormat(responseFormat []string) (string, error) {
	switch len(responseFormat) {
	case 0:
		return "", nil
	case 1:
		return responseFormat[0], nil
	default:
		return "", wcerrors.NewUsageError(fmt.Sprintf(
			"responseFormat: pass one format or none, got %d — the trailing argument is optional, "+
				"and omitting it defers to the client", len(responseFormat)))
	}
}

func requirePayload(payload any) error {
	if payload == nil || isNilValue(payload) {
		return wcerrors.NewUsageError("payload is required: this endpoint takes a body, and nil would " +
			"send the request with no body at all")
	}
	return nil
}

func isNilValue(value any) bool {
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Map, reflect.Slice, reflect.Pointer, reflect.Interface, reflect.Func:
		return reflected.IsNil()
	}
	return false
}

// New binds the page endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *Pages {
	return &Pages{requester: requester}
}
