package taxcodes

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// TaxCodes is the tax code endpoints bound to one client.
type TaxCodes struct {
	requester http.Requester
}

// ListTaxCodes lists the tax codes — GET /admin/tax_codes; limit is capped at 100, so page on pagination.total.
func (t *TaxCodes) ListTaxCodes(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendGetRequest(ctx, TaxCodesPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: TaxCodePageShape})
}

// HeadTaxCodes asks whether the collection is reachable — HEAD /admin/tax_codes.
func (t *TaxCodes) HeadTaxCodes(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendApiRequest(ctx, "HEAD", TaxCodesPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CreateTaxCode creates a tax code — POST /admin/tax_codes, wrapped in tax_code; name is required.
func (t *TaxCodes) CreateTaxCode(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendPostRequest(ctx, TaxCodesPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: TaxCodeShape})
}

// GetTaxCode retrieves one tax code — GET /admin/tax_codes/{tax_code_id}.
func (t *TaxCodes) GetTaxCode(ctx context.Context, taxCodeID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendGetRequest(ctx, TaxCodePath, &http.RequestOptions{
		PathParams:   map[string]any{"tax_code_id": taxCodeID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: TaxCodeShape})
}

// HeadTaxCode answers 200 even for an unknown id — HEAD /admin/tax_codes/{tax_code_id}; use GetTaxCode to check a row.
func (t *TaxCodes) HeadTaxCode(ctx context.Context, taxCodeID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendApiRequest(ctx, "HEAD", TaxCodePath, &http.RequestOptions{
		PathParams:   map[string]any{"tax_code_id": taxCodeID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// UpdateTaxCode updates one tax code — PUT /admin/tax_codes/{tax_code_id}; send priority and resolver_type, as they reset.
func (t *TaxCodes) UpdateTaxCode(ctx context.Context, taxCodeID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendPutRequest(ctx, TaxCodePath, &http.RequestOptions{
		PathParams: map[string]any{"tax_code_id": taxCodeID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: TaxCodeShape})
}

// DeleteTaxCode deletes one tax code — DELETE /admin/tax_codes/{tax_code_id}, 204 and hard.
func (t *TaxCodes) DeleteTaxCode(ctx context.Context, taxCodeID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendDeleteRequest(ctx, TaxCodePath, &http.RequestOptions{
		PathParams:   map[string]any{"tax_code_id": taxCodeID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "tax_codes"

// New binds the tax code endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *TaxCodes {
	return &TaxCodes{requester: requester}
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
