package orderflags

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// OrderFlags is the order flag endpoints, bound to one client.
type OrderFlags struct {
	requester http.Requester
}

// ListOrderFlags lists the flags — GET /admin/order_flags; an unknown query name is 400, and limit caps at 20.
func (o *OrderFlags) ListOrderFlags(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendGetRequest(ctx, OrderFlagsPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: OrderFlagPageShape})
}

// HeadOrderFlags asks whether the collection is reachable — HEAD /admin/order_flags.
func (o *OrderFlags) HeadOrderFlags(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendApiRequest(ctx, "HEAD", OrderFlagsPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CreateOrderFlag creates a flag — POST /admin/order_flags, wrapped in order_flag; name and a hex color are required.
func (o *OrderFlags) CreateOrderFlag(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendPostRequest(ctx, OrderFlagsPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: OrderFlagShape})
}

// GetOrderFlag retrieves one flag — GET /admin/order_flags/{order_flag_id}.
func (o *OrderFlags) GetOrderFlag(ctx context.Context, orderFlagID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendGetRequest(ctx, OrderFlagPath, &http.RequestOptions{
		PathParams: map[string]any{"order_flag_id": orderFlagID}, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: OrderFlagShape})
}

// HeadOrderFlag asks whether one flag is there — HEAD /admin/order_flags/{order_flag_id}.
func (o *OrderFlags) HeadOrderFlag(ctx context.Context, orderFlagID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendApiRequest(ctx, "HEAD", OrderFlagPath, &http.RequestOptions{
		PathParams: map[string]any{"order_flag_id": orderFlagID}, RaiseOnError: raiseOnError,
		ResponseFormat: format})
}

// UpdateOrderFlag replaces one flag — PUT /admin/order_flags/{order_flag_id}; send the full record, as PATCH is 405.
func (o *OrderFlags) UpdateOrderFlag(ctx context.Context, orderFlagID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendPutRequest(ctx, OrderFlagPath, &http.RequestOptions{
		PathParams: map[string]any{"order_flag_id": orderFlagID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: OrderFlagShape})
}

// DeleteOrderFlag deletes a flag — DELETE /admin/order_flags/{order_flag_id}, 204 and hard.
func (o *OrderFlags) DeleteOrderFlag(ctx context.Context, orderFlagID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendDeleteRequest(ctx, OrderFlagPath, &http.RequestOptions{
		PathParams: map[string]any{"order_flag_id": orderFlagID}, RaiseOnError: raiseOnError,
		ResponseFormat: format})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "order_flags"

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

// New binds the order flag endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *OrderFlags { return &OrderFlags{requester: requester} }
