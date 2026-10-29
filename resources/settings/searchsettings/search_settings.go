package searchsettings

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// SearchSettingsResource is the search settings endpoints: a facade over three sections, plus the read-only reindex_status.
type SearchSettingsResource struct {
	requester http.Requester
}

// GetSearchSettings reads the whole area — GET /admin/settings/search_settings.
func (s *SearchSettingsResource) GetSearchSettings(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, SearchSettingsPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: SearchSettingsShape})
}

// PatchSearchSettings writes one or more sections — PATCH /admin/settings/search_settings; never name synonyms.
func (s *SearchSettingsResource) PatchSearchSettings(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, SearchSettingsPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: SearchSettingsShape})
}

// PutSearchSettings replaces the whole area — PUT /admin/settings/search_settings; omitted fields reset, so prefer the PATCH.
func (s *SearchSettingsResource) PutSearchSettings(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPutRequest(ctx, SearchSettingsPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: SearchSettingsShape})
}

// HeadSearchSettings asks whether the area is reachable — HEAD /admin/settings/search_settings.
func (s *SearchSettingsResource) HeadSearchSettings(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", SearchSettingsPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// GetSearchSettingsSection reads one section — GET /admin/settings/search_settings/{section}, under the section's key.
func (s *SearchSettingsResource) GetSearchSettingsSection(ctx context.Context, section string, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, SearchSettingsSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, QueryParams: queryParams,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: SearchSettingsShape})
}

// PatchSearchSettingsSection writes one section — PATCH /admin/settings/search_settings/{section}.
func (s *SearchSettingsResource) PatchSearchSettingsSection(ctx context.Context, section string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, SearchSettingsSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: SearchSettingsShape})
}

// PutSearchSettingsSection replaces one section — PUT /admin/settings/search_settings/{section}; send it whole, never synonyms.
func (s *SearchSettingsResource) PutSearchSettingsSection(ctx context.Context, section string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPutRequest(ctx, SearchSettingsSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: SearchSettingsShape})
}

// GetSearchSettingsReindexStatus reads the reindex state — GET /admin/settings/search_settings/reindex_status, unwrapped.
func (s *SearchSettingsResource) GetSearchSettingsReindexStatus(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, SearchSettingsReindexStatusPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ReindexStatusShape})
}

// HeadSearchSettingsReindexStatus asks whether the reindex status is reachable — HEAD /admin/settings/search_settings/reindex_status; no ETag.
func (s *SearchSettingsResource) HeadSearchSettingsReindexStatus(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", SearchSettingsReindexStatusPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "search_settings"

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

// New binds the search settings endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *SearchSettingsResource {
	return &SearchSettingsResource{requester: requester}
}
