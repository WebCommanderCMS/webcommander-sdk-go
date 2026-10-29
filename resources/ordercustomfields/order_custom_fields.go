package ordercustomfields

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// OrderCustomFields is the order custom field endpoints, bound to one client.
type OrderCustomFields struct {
	requester http.Requester
}

// ListOrderCustomFields lists the fields — GET /admin/order_custom_fields; limit caps at 20, and sort takes id, label, name, type, orderNo, isActive, created or updated.
func (o *OrderCustomFields) ListOrderCustomFields(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendGetRequest(ctx, OrderCustomFieldsPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: OrderCustomFieldPageShape})
}

// CountOrderCustomFields is how many fields the store holds — GET /admin/order_custom_fields/count; Map() is {"count": N}, which filterBy narrows as it narrows the listing.
func (o *OrderCustomFields) CountOrderCustomFields(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendGetRequest(ctx, OrderCustomFieldsCountPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError, ResponseFormat: format})
}

// GetOrderCustomField retrieves one field, with nine fields its listing row leaves out — GET /admin/order_custom_fields/{order_custom_field_id}.
func (o *OrderCustomFields) GetOrderCustomField(ctx context.Context, orderCustomFieldID any, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendGetRequest(ctx, OrderCustomFieldPath, &http.RequestOptions{
		PathParams:  map[string]any{"order_custom_field_id": orderCustomFieldID},
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: OrderCustomFieldShape})
}

// ValidateOrderCustomFieldTimeslot checks a timeslot against an order limit — POST /admin/order_custom_fields/{order_custom_field_id}/validate_timeslot.
func (o *OrderCustomFields) ValidateOrderCustomFieldTimeslot(ctx context.Context, orderCustomFieldID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendPostRequest(ctx, OrderCustomFieldValidateTimeslotPath, &http.RequestOptions{
		PathParams: map[string]any{"order_custom_field_id": orderCustomFieldID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "order_custom_fields"

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

// New binds the order custom field endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *OrderCustomFields {
	return &OrderCustomFields{requester: requester}
}
