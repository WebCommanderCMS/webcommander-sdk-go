package taxcountrytaxes

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// TaxCountryTaxes is the country tax endpoints, bound to one client: a full collection with no HEAD.
type TaxCountryTaxes struct {
	requester http.Requester
}

// ListTaxCountryTaxes lists the country taxes, newest first — GET /admin/tax/country-taxes; limit is 1 to 100, sort takes id or name.
func (t *TaxCountryTaxes) ListTaxCountryTaxes(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendGetRequest(ctx, TaxCountryTaxesPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CountryTaxPageShape})
}

// CreateTaxCountryTax creates a country tax — POST /admin/tax/country-taxes, wrapped in country_tax, answered 201; country_code is required. transactionNo is sent as X-Transaction-No and files the row under that number; empty sends none.
func (t *TaxCountryTaxes) CreateTaxCountryTax(ctx context.Context, payload any, transactionNo string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	var headers map[string]string
	if transactionNo != "" {
		headers = map[string]string{TaxCountryTaxesTransactionNoHeader: transactionNo}
	}
	return t.requester.SendPostRequest(ctx, TaxCountryTaxesPath, &http.RequestOptions{
		JSONBody: payload, Headers: headers, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CountryTaxShape})
}

// GetTaxCountryTax retrieves one country tax, with an ETag and state_rules — GET /admin/tax/country-taxes/{country_tax_id}.
func (t *TaxCountryTaxes) GetTaxCountryTax(ctx context.Context, countryTaxID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendGetRequest(ctx, TaxCountryTaxPath, &http.RequestOptions{
		PathParams:   map[string]any{"country_tax_id": countryTaxID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CountryTaxShape})
}

// UpdateTaxCountryTax merges into one country tax — PUT /admin/tax/country-taxes/{country_tax_id}, wrapped in country_tax; 200 with an ETag.
func (t *TaxCountryTaxes) UpdateTaxCountryTax(ctx context.Context, countryTaxID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendPutRequest(ctx, TaxCountryTaxPath, &http.RequestOptions{
		PathParams: map[string]any{"country_tax_id": countryTaxID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CountryTaxShape})
}

// DeleteTaxCountryTax deletes one country tax — DELETE /admin/tax/country-taxes/{country_tax_id}, 204 and hard.
func (t *TaxCountryTaxes) DeleteTaxCountryTax(ctx context.Context, countryTaxID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendDeleteRequest(ctx, TaxCountryTaxPath, &http.RequestOptions{
		PathParams:   map[string]any{"country_tax_id": countryTaxID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "tax_country_taxes"

// New binds the country tax endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *TaxCountryTaxes {
	return &TaxCountryTaxes{requester: requester}
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
