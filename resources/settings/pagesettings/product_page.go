package pagesettings

import (
	"context"

	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// ProductPageSettings is the product page settings endpoints, bound to one client: a facade over seven sections.
type ProductPageSettings struct {
	requester http.Requester
}

// GetProductPage reads the whole area — GET /admin/settings/page_settings/product_page.
func (s *ProductPageSettings) GetProductPage(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, ProductPagePath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ProductPageShape})
}

// PatchProductPage writes one or more sections — PATCH /admin/settings/page_settings/product_page, wrapped or bare.
func (s *ProductPageSettings) PatchProductPage(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, ProductPagePath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ProductPageShape})
}

// HeadProductPage asks whether the area is reachable — HEAD /admin/settings/page_settings/product_page.
func (s *ProductPageSettings) HeadProductPage(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", ProductPagePath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// GetProductPageSection reads one section — GET /admin/settings/page_settings/product_page/{section}, under the facade's key.
func (s *ProductPageSettings) GetProductPageSection(ctx context.Context, section string, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, ProductPageSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, QueryParams: queryParams,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ProductPageShape})
}

// PatchProductPageSection writes one section — PATCH /admin/settings/page_settings/product_page/{section}, under either key.
func (s *ProductPageSettings) PatchProductPageSection(ctx context.Context, section string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, ProductPageSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ProductPageShape})
}

// ProductPageResourceName is this resource's registry key.
const ProductPageResourceName = "product_page"

// NewProductPageSettings binds the product page settings endpoints to any requester.
func NewProductPageSettings(requester http.Requester) *ProductPageSettings {
	return &ProductPageSettings{requester: requester}
}
