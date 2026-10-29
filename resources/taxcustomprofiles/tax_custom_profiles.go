package taxcustomprofiles

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// TaxCustomProfiles is the custom tax profile endpoints, bound to one client: a full collection with no HEAD.
type TaxCustomProfiles struct {
	requester http.Requester
}

// ListTaxCustomProfiles lists the custom profiles, newest first — GET /admin/tax/custom-profiles; limit is 1 to 100, sort takes id or name.
func (t *TaxCustomProfiles) ListTaxCustomProfiles(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendGetRequest(ctx, TaxCustomProfilesPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CustomProfilePageShape})
}

// CreateTaxCustomProfile creates a custom profile — POST /admin/tax/custom-profiles, wrapped in custom_profile, answered 201; name is required. transactionNo is sent as X-Transaction-No and files the profile under that number; empty sends none.
func (t *TaxCustomProfiles) CreateTaxCustomProfile(ctx context.Context, payload any, transactionNo string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	var headers map[string]string
	if transactionNo != "" {
		headers = map[string]string{TaxCustomProfilesTransactionNoHeader: transactionNo}
	}
	return t.requester.SendPostRequest(ctx, TaxCustomProfilesPath, &http.RequestOptions{
		JSONBody: payload, Headers: headers, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CustomProfileShape})
}

// GetTaxCustomProfile retrieves one custom profile, with an ETag and rules — GET /admin/tax/custom-profiles/{custom_profile_id}.
func (t *TaxCustomProfiles) GetTaxCustomProfile(ctx context.Context, customProfileID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendGetRequest(ctx, TaxCustomProfilePath, &http.RequestOptions{
		PathParams:   map[string]any{"custom_profile_id": customProfileID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CustomProfileShape})
}

// UpdateTaxCustomProfile replaces one custom profile — PUT /admin/tax/custom-profiles/{custom_profile_id}; an omitted name is kept, an omitted description or tax_code comes back null.
func (t *TaxCustomProfiles) UpdateTaxCustomProfile(ctx context.Context, customProfileID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendPutRequest(ctx, TaxCustomProfilePath, &http.RequestOptions{
		PathParams: map[string]any{"custom_profile_id": customProfileID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CustomProfileShape})
}

// DeleteTaxCustomProfile deletes one custom profile — DELETE /admin/tax/custom-profiles/{custom_profile_id}, 204 and hard.
func (t *TaxCustomProfiles) DeleteTaxCustomProfile(ctx context.Context, customProfileID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendDeleteRequest(ctx, TaxCustomProfilePath, &http.RequestOptions{
		PathParams:   map[string]any{"custom_profile_id": customProfileID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "tax_custom_profiles"

// New binds the custom profile endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *TaxCustomProfiles {
	return &TaxCustomProfiles{requester: requester}
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
