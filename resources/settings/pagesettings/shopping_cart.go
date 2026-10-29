package pagesettings

import (
	"context"

	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// ShoppingCartSettings is the shopping cart settings endpoints, bound to one client: a facade over two sections.
type ShoppingCartSettings struct {
	requester http.Requester
}

// GetShoppingCart reads the whole area — GET /admin/settings/page_settings/shopping_cart.
func (s *ShoppingCartSettings) GetShoppingCart(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, ShoppingCartPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ShoppingCartShape})
}

// PatchShoppingCart writes one or both sections — PATCH /admin/settings/page_settings/shopping_cart, wrapped in shopping_cart.
func (s *ShoppingCartSettings) PatchShoppingCart(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, ShoppingCartPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ShoppingCartShape})
}

// HeadShoppingCart asks whether the area is reachable — HEAD /admin/settings/page_settings/shopping_cart.
func (s *ShoppingCartSettings) HeadShoppingCart(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", ShoppingCartPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// GetShoppingCartSection reads one section — GET /admin/settings/page_settings/shopping_cart/{section}, under its own key.
func (s *ShoppingCartSettings) GetShoppingCartSection(ctx context.Context, section string, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, ShoppingCartSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, QueryParams: queryParams,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ShoppingCartShape})
}

// PatchShoppingCartSection writes one section — PATCH /admin/settings/page_settings/shopping_cart/{section}, under its own key.
func (s *ShoppingCartSettings) PatchShoppingCartSection(ctx context.Context, section string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, ShoppingCartSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ShoppingCartShape})
}

// ShoppingCartResourceName is this resource's registry key.
const ShoppingCartResourceName = "shopping_cart"

// NewShoppingCartSettings binds the shopping cart settings endpoints to any requester.
func NewShoppingCartSettings(requester http.Requester) *ShoppingCartSettings {
	return &ShoppingCartSettings{requester: requester}
}
