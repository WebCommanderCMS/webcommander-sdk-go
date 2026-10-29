package pagesettings

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/customers"
)

// SearchResultPageSettings is the search result page settings endpoints: a facade over three sections and one collection.
type SearchResultPageSettings struct {
	requester http.Requester
}

// GetSearchResultPage reads the whole area — GET /admin/settings/page_settings/search_result_page.
func (s *SearchResultPageSettings) GetSearchResultPage(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, SearchResultPagePath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: SearchResultPageShape})
}

// PatchSearchResultPage writes the writable section — PATCH /admin/settings/page_settings/search_result_page.
func (s *SearchResultPageSettings) PatchSearchResultPage(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, SearchResultPagePath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: SearchResultPageShape})
}

// HeadSearchResultPage asks whether the area is reachable — HEAD /admin/settings/page_settings/search_result_page.
func (s *SearchResultPageSettings) HeadSearchResultPage(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", SearchResultPagePath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// GetSearchResultPageSection reads one section — GET /admin/settings/page_settings/search_result_page/{section}, under the facade's key.
func (s *SearchResultPageSettings) GetSearchResultPageSection(ctx context.Context, section string, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, SearchResultPageSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, QueryParams: queryParams,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: SearchResultPageShape})
}

// PatchSearchResultPageSection writes product_information — PATCH /admin/settings/page_settings/search_result_page/{section}.
func (s *SearchResultPageSettings) PatchSearchResultPageSection(ctx context.Context, section string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, SearchResultPageSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: SearchResultPageShape})
}

// GetSearchResultPageProductPromotions lists the promotions — GET /admin/settings/page_settings/search_result_page/product_promotion; limit is not applied.
func (s *SearchResultPageSettings) GetSearchResultPageProductPromotions(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, SearchResultPageProductPromotionsPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: customers.PageOf[ProductPromotion]})
}

// HeadSearchResultPageProductPromotions asks whether the promotions collection is reachable — HEAD /admin/settings/page_settings/search_result_page/product_promotion; no ETag.
func (s *SearchResultPageSettings) HeadSearchResultPageProductPromotions(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", SearchResultPageProductPromotionsPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CreateSearchResultPageProductPromotion adds a promotion — POST /admin/settings/page_settings/search_result_page/product_promotion; name, rule and position required.
func (s *SearchResultPageSettings) CreateSearchResultPageProductPromotion(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPostRequest(ctx, SearchResultPageProductPromotionsPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ProductPromotionShape})
}

// GetSearchResultPageProductPromotion reads one promotion — GET /admin/settings/page_settings/search_result_page/product_promotion/{promotion_id}.
func (s *SearchResultPageSettings) GetSearchResultPageProductPromotion(ctx context.Context, promotionID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, SearchResultPageProductPromotionPath, &http.RequestOptions{
		PathParams: map[string]any{"promotion_id": promotionID}, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ProductPromotionShape})
}

// HeadSearchResultPageProductPromotion asks whether one promotion is there — HEAD /admin/settings/page_settings/search_result_page/product_promotion/{promotion_id}; an ETag, and 404 for an unknown id.
func (s *SearchResultPageSettings) HeadSearchResultPageProductPromotion(ctx context.Context, promotionID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", SearchResultPageProductPromotionPath, &http.RequestOptions{
		PathParams: map[string]any{"promotion_id": promotionID}, RaiseOnError: raiseOnError,
		ResponseFormat: format})
}

// PatchSearchResultPageProductPromotion changes one promotion — PATCH /admin/settings/page_settings/search_result_page/product_promotion/{promotion_id}; partial.
func (s *SearchResultPageSettings) PatchSearchResultPageProductPromotion(ctx context.Context, promotionID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, SearchResultPageProductPromotionPath, &http.RequestOptions{
		PathParams: map[string]any{"promotion_id": promotionID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ProductPromotionShape})
}

// PatchSearchResultPageProductPromotionActive turns one promotion on or off — PATCH /admin/settings/page_settings/search_result_page/product_promotion/{promotion_id}/active; only is_active is read.
func (s *SearchResultPageSettings) PatchSearchResultPageProductPromotionActive(ctx context.Context, promotionID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, SearchResultPageProductPromotionActivePath, &http.RequestOptions{
		PathParams: map[string]any{"promotion_id": promotionID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ProductPromotionShape})
}

// DeleteSearchResultPageProductPromotion removes one promotion — DELETE /admin/settings/page_settings/search_result_page/product_promotion/{promotion_id}.
func (s *SearchResultPageSettings) DeleteSearchResultPageProductPromotion(ctx context.Context, promotionID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendDeleteRequest(ctx, SearchResultPageProductPromotionPath, &http.RequestOptions{
		PathParams: map[string]any{"promotion_id": promotionID}, RaiseOnError: raiseOnError,
		ResponseFormat: format})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "search_result_page"

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

// New binds the search result page settings endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *SearchResultPageSettings {
	return &SearchResultPageSettings{requester: requester}
}
