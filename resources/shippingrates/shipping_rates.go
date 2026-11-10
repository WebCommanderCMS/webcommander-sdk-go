package shippingrates

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// ShippingRates is the shipping rate endpoints, bound to one client.
type ShippingRates struct {
	requester http.Requester
}

// ListShippingRates lists the rates, newest first — GET /admin/shipping_rates; search (part of the name), method (such as flat_rate or ship_by_weight) and handling_cost (with_handling or without_handling) filter, sort is id, name, method, created_at, updated_at or transaction_no with order asc or desc, limit is 1 to 100 and defaults to 20, offset pages; an unknown value is 400.
func (s *ShippingRates) ListShippingRates(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, ShippingRatesPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ShippingRateListingShape})
}

// HeadShippingRates asks whether the collection is reachable — HEAD /admin/shipping_rates; no body.
func (s *ShippingRates) HeadShippingRates(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", ShippingRatesPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CountShippingRates is how many rates match — GET /admin/shipping_rates/count, a bare {"total": N}; the listing's filters narrow it.
func (s *ShippingRates) CountShippingRates(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, ShippingRatesCountPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ShippingRateCountShape})
}

// HeadShippingRatesCount asks whether the count is reachable — HEAD /admin/shipping_rates/count; no body.
func (s *ShippingRates) HeadShippingRatesCount(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", ShippingRatesCountPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CreateShippingRate creates a rate — POST /admin/shipping_rates with {"shipping_rate": {…}}, answered 201 with a Location; name is required and unique without regard to case (409), and so is method. A flat_rate needs shipping_cost; ship_by_price, ship_by_quantity, ship_by_weight and ship_by_variation take their costs per row in conditions and refuse a rate-level shipping_cost; free_shipping always saves calculate_for_entire_order false. shipping_rule_id attaches the rate to a rule that has none (409 otherwise). An unknown field, an id or an unwrapped body is 400. transactionNo is sent as X-Transaction-No and wins over a transaction_no in the body; empty sends none.
func (s *ShippingRates) CreateShippingRate(ctx context.Context, payload any, transactionNo string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	var headers map[string]string
	if transactionNo != "" {
		headers = map[string]string{ShippingRatesTransactionNoHeader: transactionNo}
	}
	return s.requester.SendPostRequest(ctx, ShippingRatesPath, &http.RequestOptions{
		JSONBody: payload, Headers: headers, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ShippingRateShape})
}

// GetShippingRate retrieves one rate by its id — GET /admin/shipping_rates/{shipping_rate_id}; an id that is not plain digits, or unknown, is 404.
func (s *ShippingRates) GetShippingRate(ctx context.Context, shippingRateID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, ShippingRatePath, &http.RequestOptions{
		PathParams:   map[string]any{"shipping_rate_id": shippingRateID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ShippingRateShape})
}

// HeadShippingRate asks whether a rate exists — HEAD /admin/shipping_rates/{shipping_rate_id}; no body, 404 for an unknown id.
func (s *ShippingRates) HeadShippingRate(ctx context.Context, shippingRateID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", ShippingRatePath, &http.RequestOptions{
		PathParams:   map[string]any{"shipping_rate_id": shippingRateID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// UpdateShippingRate replaces one rate — PUT /admin/shipping_rates/{shipping_rate_id}; name, method, calculate_for_entire_order, price_entered_with_tax, the costs the method uses, handling_cost and additional_rule are all required, and one the rate does not use is sent as null; a field left out is 400 naming it.
func (s *ShippingRates) UpdateShippingRate(ctx context.Context, shippingRateID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPutRequest(ctx, ShippingRatePath, &http.RequestOptions{
		PathParams: map[string]any{"shipping_rate_id": shippingRateID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ShippingRateShape})
}

// PatchShippingRate changes only the fields sent — PATCH /admin/shipping_rates/{shipping_rate_id}; the costs, rows and boxes not sent are kept, and an empty shipping_rate changes nothing and is 200.
func (s *ShippingRates) PatchShippingRate(ctx context.Context, shippingRateID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, ShippingRatePath, &http.RequestOptions{
		PathParams: map[string]any{"shipping_rate_id": shippingRateID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ShippingRateShape})
}

// DeleteShippingRate deletes one rate — DELETE /admin/shipping_rates/{shipping_rate_id}, 204; a rate a rule uses is 409 until the rule lets it go.
func (s *ShippingRates) DeleteShippingRate(ctx context.Context, shippingRateID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendDeleteRequest(ctx, ShippingRatePath, &http.RequestOptions{
		PathParams:   map[string]any{"shipping_rate_id": shippingRateID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "shipping_rates"

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

// New binds the shipping rate endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *ShippingRates {
	return &ShippingRates{requester: requester}
}
