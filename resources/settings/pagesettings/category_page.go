package pagesettings

import (
	"context"

	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// CategoryPageSettings is the category page settings endpoints, bound to one client: a facade over two sections.
type CategoryPageSettings struct {
	requester http.Requester
}

// GetCategoryPage reads the whole area — GET /admin/settings/page_settings/category_page.
func (s *CategoryPageSettings) GetCategoryPage(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, CategoryPagePath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CategoryPageShape})
}

// PatchCategoryPage writes one or both sections — PATCH /admin/settings/page_settings/category_page, wrapped in category_page.
func (s *CategoryPageSettings) PatchCategoryPage(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, CategoryPagePath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CategoryPageShape})
}

// HeadCategoryPage asks whether the area is reachable — HEAD /admin/settings/page_settings/category_page.
func (s *CategoryPageSettings) HeadCategoryPage(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", CategoryPagePath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// GetCategoryPageSection reads one section — GET /admin/settings/page_settings/category_page/{section}, under its own key.
func (s *CategoryPageSettings) GetCategoryPageSection(ctx context.Context, section string, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, CategoryPageSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, QueryParams: queryParams,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CategoryPageShape})
}

// PatchCategoryPageSection writes one section — PATCH /admin/settings/page_settings/category_page/{section}, under its own key.
func (s *CategoryPageSettings) PatchCategoryPageSection(ctx context.Context, section string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, CategoryPageSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CategoryPageShape})
}

// CategoryPageResourceName is this resource's registry key.
const CategoryPageResourceName = "category_page"

// NewCategoryPageSettings binds the category page settings endpoints to any requester.
func NewCategoryPageSettings(requester http.Requester) *CategoryPageSettings {
	return &CategoryPageSettings{requester: requester}
}
