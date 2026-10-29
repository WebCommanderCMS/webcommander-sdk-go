package pricedisplaysettings

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// PriceDisplaySettingsResource is the price display settings endpoints, bound to one client: a facade over two sections.
type PriceDisplaySettingsResource struct {
	requester http.Requester
}

// GetPriceDisplaySettings reads the whole area — GET /admin/settings/price_display_settings.
func (s *PriceDisplaySettingsResource) GetPriceDisplaySettings(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, PriceDisplaySettingsPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: PriceDisplaySettingsShape})
}

// PatchPriceDisplaySettings writes fields in one section or both — PATCH /admin/settings/price_display_settings, wrapped in price_display_settings only.
func (s *PriceDisplaySettingsResource) PatchPriceDisplaySettings(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, PriceDisplaySettingsPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: PriceDisplaySettingsShape})
}

// HeadPriceDisplaySettings asks whether the area is reachable — HEAD /admin/settings/price_display_settings.
func (s *PriceDisplaySettingsResource) HeadPriceDisplaySettings(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", PriceDisplaySettingsPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// GetPriceDisplaySettingsSection reads one section — GET /admin/settings/price_display_settings/{section}, under its own key.
func (s *PriceDisplaySettingsResource) GetPriceDisplaySettingsSection(ctx context.Context, section string, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, PriceDisplaySettingsSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, QueryParams: queryParams,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: PriceDisplaySettingsShape})
}

// PatchPriceDisplaySettingsSection writes one section — PATCH /admin/settings/price_display_settings/{section}; the answer is computed.
func (s *PriceDisplaySettingsResource) PatchPriceDisplaySettingsSection(ctx context.Context, section string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, PriceDisplaySettingsSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: PriceDisplaySettingsShape})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "price_display_settings"

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

// New binds the price display settings endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *PriceDisplaySettingsResource {
	return &PriceDisplaySettingsResource{requester: requester}
}
