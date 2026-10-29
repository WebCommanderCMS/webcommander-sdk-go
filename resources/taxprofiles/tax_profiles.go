package taxprofiles

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// TaxProfiles is the tax profile endpoints, bound to one client; default is exclusive and one-way.
type TaxProfiles struct {
	requester http.Requester
}

// ListTaxProfiles lists the tax profiles — GET /admin/tax_profiles; limit is clamped, so page on pagination.total.
func (t *TaxProfiles) ListTaxProfiles(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendGetRequest(ctx, TaxProfilesPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: TaxProfilePageShape})
}

// HeadTaxProfiles asks whether the collection is reachable — HEAD /admin/tax_profiles.
func (t *TaxProfiles) HeadTaxProfiles(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendApiRequest(ctx, "HEAD", TaxProfilesPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CreateTaxProfile creates a tax profile — POST /admin/tax_profiles, wrapped in tax_profile; an id in the body is 500.
func (t *TaxProfiles) CreateTaxProfile(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendPostRequest(ctx, TaxProfilesPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: TaxProfileShape})
}

// GetTaxProfile retrieves one tax profile — GET /admin/tax_profiles/{tax_profile_id}.
func (t *TaxProfiles) GetTaxProfile(ctx context.Context, taxProfileID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendGetRequest(ctx, TaxProfilePath, &http.RequestOptions{
		PathParams:   map[string]any{"tax_profile_id": taxProfileID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: TaxProfileShape})
}

// HeadTaxProfile asks whether the path answers — HEAD /admin/tax_profiles/{tax_profile_id}; an unknown id is 200 too.
func (t *TaxProfiles) HeadTaxProfile(ctx context.Context, taxProfileID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendApiRequest(ctx, "HEAD", TaxProfilePath, &http.RequestOptions{
		PathParams:   map[string]any{"tax_profile_id": taxProfileID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// UpdateTaxProfile merges into one tax profile — PUT /admin/tax_profiles/{tax_profile_id}; default true moves the flag.
func (t *TaxProfiles) UpdateTaxProfile(ctx context.Context, taxProfileID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendPutRequest(ctx, TaxProfilePath, &http.RequestOptions{
		PathParams: map[string]any{"tax_profile_id": taxProfileID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: TaxProfileShape})
}

// DeleteTaxProfile deletes one tax profile — DELETE /admin/tax_profiles/{tax_profile_id}; the default row deletes too.
func (t *TaxProfiles) DeleteTaxProfile(ctx context.Context, taxProfileID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendDeleteRequest(ctx, TaxProfilePath, &http.RequestOptions{
		PathParams:   map[string]any{"tax_profile_id": taxProfileID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "tax_profiles"

// New binds the tax profile endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *TaxProfiles {
	return &TaxProfiles{requester: requester}
}

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
