package pagesettings

import (
	"context"

	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// CheckoutPageSettings is the checkout page settings endpoints, bound to one client: a facade over six sections.
type CheckoutPageSettings struct {
	requester http.Requester
}

// GetCheckoutPage reads the whole area — GET /admin/settings/page_settings/checkout_page.
func (s *CheckoutPageSettings) GetCheckoutPage(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, CheckoutPagePath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CheckoutPageShape})
}

// PatchCheckoutPage writes one or more sections — PATCH /admin/settings/page_settings/checkout_page, wrapped in checkout_page.
func (s *CheckoutPageSettings) PatchCheckoutPage(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, CheckoutPagePath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CheckoutPageShape})
}

// HeadCheckoutPage asks whether the area is reachable — HEAD /admin/settings/page_settings/checkout_page.
func (s *CheckoutPageSettings) HeadCheckoutPage(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", CheckoutPagePath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// GetCheckoutPageSection reads one section — GET /admin/settings/page_settings/checkout_page/{section}, under its own key.
func (s *CheckoutPageSettings) GetCheckoutPageSection(ctx context.Context, section string, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, CheckoutPageSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, QueryParams: queryParams,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CheckoutPageShape})
}

// PatchCheckoutPageSection writes one writable section — PATCH /admin/settings/page_settings/checkout_page/{section}.
func (s *CheckoutPageSettings) PatchCheckoutPageSection(ctx context.Context, section string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, CheckoutPageSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CheckoutPageShape})
}

// CheckoutPageResourceName is this resource's registry key.
const CheckoutPageResourceName = "checkout_page"

// NewCheckoutPageSettings binds the checkout page settings endpoints to any requester.
func NewCheckoutPageSettings(requester http.Requester) *CheckoutPageSettings {
	return &CheckoutPageSettings{requester: requester}
}
