package currencies

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// Currencies is the currency endpoints, bound to one client.
type Currencies struct {
	requester http.Requester
}

// ListCurrencies lists the currencies by ascending id — GET /admin/currencies; search matches name or code, active filters, sort and dir order, limit is 1 to 100.
func (c *Currencies) ListCurrencies(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, CurrenciesPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CurrencyPageShape})
}

// HeadCurrencies asks whether the collection is reachable — HEAD /admin/currencies; no body.
func (c *Currencies) HeadCurrencies(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendApiRequest(ctx, "HEAD", CurrenciesPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CountCurrencies is how many currencies match — GET /admin/currencies/count, a bare {"total": N}; search and active narrow it as they narrow the listing, limit is ignored.
func (c *Currencies) CountCurrencies(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, CurrenciesCountPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CurrencyCountShape})
}

// HeadCurrenciesCount asks whether the count is reachable — HEAD /admin/currencies/count; no body.
func (c *Currencies) HeadCurrenciesCount(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendApiRequest(ctx, "HEAD", CurrenciesCountPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CreateCurrency creates a currency — POST /admin/currencies as a bare record, answered 201; name, code and symbol are required, and name and code must be free. transactionNo is sent as X-Transaction-No and files the currency under that number; empty sends none.
func (c *Currencies) CreateCurrency(ctx context.Context, payload any, transactionNo string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	var headers map[string]string
	if transactionNo != "" {
		headers = map[string]string{CurrenciesTransactionNoHeader: transactionNo}
	}
	return c.requester.SendPostRequest(ctx, CurrenciesPath, &http.RequestOptions{
		JSONBody: payload, Headers: headers, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CurrencyShape})
}

// GetCurrency retrieves one currency, with an ETag — GET /admin/currencies/{currency_id}; field_metadata adds the schema beside it.
func (c *Currencies) GetCurrency(ctx context.Context, currencyID any, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, CurrencyPath, &http.RequestOptions{
		PathParams:  map[string]any{"currency_id": currencyID},
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CurrencyShape})
}

// HeadCurrency asks whether one currency is there — HEAD /admin/currencies/{currency_id}; no body, the ETag its GET carries.
func (c *Currencies) HeadCurrency(ctx context.Context, currencyID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendApiRequest(ctx, "HEAD", CurrencyPath, &http.RequestOptions{
		PathParams:   map[string]any{"currency_id": currencyID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// UpdateCurrency replaces one currency — PUT /admin/currencies/{currency_id}.
func (c *Currencies) UpdateCurrency(ctx context.Context, currencyID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendPutRequest(ctx, CurrencyPath, &http.RequestOptions{
		PathParams: map[string]any{"currency_id": currencyID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CurrencyShape})
}

// PatchCurrency changes only the fields sent — PATCH /admin/currencies/{currency_id}; an empty body is 400, and the base is read-only here as on PUT.
func (c *Currencies) PatchCurrency(ctx context.Context, currencyID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendPatchRequest(ctx, CurrencyPath, &http.RequestOptions{
		PathParams: map[string]any{"currency_id": currencyID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CurrencyShape})
}

// DeleteCurrency deletes a currency — DELETE /admin/currencies/{currency_id}; 204, and the base currency is 409.
func (c *Currencies) DeleteCurrency(ctx context.Context, currencyID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendDeleteRequest(ctx, CurrencyPath, &http.RequestOptions{
		PathParams:   map[string]any{"currency_id": currencyID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// SetBaseCurrency makes an active currency the store's base — PUT /admin/currencies/{currency_id}/base with no body; every rate is divided by the new base's old rate.
func (c *Currencies) SetBaseCurrency(ctx context.Context, currencyID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendPutRequest(ctx, CurrencyBasePath, &http.RequestOptions{
		PathParams:   map[string]any{"currency_id": currencyID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CurrencyBaseChangeShape})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "currencies"

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

// New binds the currency endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *Currencies {
	return &Currencies{requester: requester}
}
