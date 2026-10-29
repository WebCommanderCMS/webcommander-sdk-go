package ecommercesettings

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// EcommerceSettingsResource is the ecommerce settings endpoints: a facade over six sections beside the order_flags key.
type EcommerceSettingsResource struct {
	requester http.Requester
}

// GetEcommerceSettings reads the whole area — GET /admin/settings/ecommerce_settings.
func (e *EcommerceSettingsResource) GetEcommerceSettings(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return e.requester.SendGetRequest(ctx, EcommerceSettingsPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: EcommerceSettingsShape})
}

// PatchEcommerceSettings writes one or more sections — PATCH /admin/settings/ecommerce_settings, wrapped in ecommerce_settings.
func (e *EcommerceSettingsResource) PatchEcommerceSettings(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return e.requester.SendPatchRequest(ctx, EcommerceSettingsPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: EcommerceSettingsShape})
}

// HeadEcommerceSettings asks whether the area is reachable — HEAD /admin/settings/ecommerce_settings.
func (e *EcommerceSettingsResource) HeadEcommerceSettings(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return e.requester.SendApiRequest(ctx, "HEAD", EcommerceSettingsPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// GetEcommerceSettingsSection reads one section — GET /admin/settings/ecommerce_settings/{section}, under the section's key.
func (e *EcommerceSettingsResource) GetEcommerceSettingsSection(ctx context.Context, section string, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return e.requester.SendGetRequest(ctx, EcommerceSettingsSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, QueryParams: queryParams,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: EcommerceSettingsShape})
}

// PatchEcommerceSettingsSection writes one section — PATCH /admin/settings/ecommerce_settings/{section}; unknown fields are 400.
func (e *EcommerceSettingsResource) PatchEcommerceSettingsSection(ctx context.Context, section string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return e.requester.SendPatchRequest(ctx, EcommerceSettingsSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: EcommerceSettingsShape})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "ecommerce_settings"

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

// New binds the ecommerce settings endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *EcommerceSettingsResource {
	return &EcommerceSettingsResource{requester: requester}
}
