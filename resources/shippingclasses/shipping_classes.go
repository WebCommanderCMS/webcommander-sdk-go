package shippingclasses

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// ShippingClasses is the shipping class endpoints, under /api/{api_version}/admin, bound to one client.
type ShippingClasses struct {
	requester http.Requester
}

// ListShippingClasses lists the classes — GET /admin/shipping_classes; limit caps at 20 and sort is 400.
func (s *ShippingClasses) ListShippingClasses(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, ShippingClassesPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ShippingClassPageShape})
}

// HeadShippingClasses asks whether the collection answers — HEAD /admin/shipping_classes.
func (s *ShippingClasses) HeadShippingClasses(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", ShippingClassesPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CountShippingClasses is how many classes the store holds — GET /admin/shipping_classes/count; no parameter narrows it.
func (s *ShippingClasses) CountShippingClasses(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, ShippingClassesCountPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ShippingClassCountShape})
}

// CreateShippingClass creates a class — POST /admin/shipping_classes, wrapped in shipping_class; id and is_default are discarded.
func (s *ShippingClasses) CreateShippingClass(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPostRequest(ctx, ShippingClassesPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ShippingClassShape})
}

// GetShippingClass retrieves one class — GET /admin/shipping_classes/{shipping_class_id}; field_metadata=true adds the schema.
func (s *ShippingClasses) GetShippingClass(ctx context.Context, shippingClassID any, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, ShippingClassPath, &http.RequestOptions{
		PathParams:  map[string]any{"shipping_class_id": shippingClassID},
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ShippingClassShape})
}

// HeadShippingClass is HEAD /admin/shipping_classes/{shipping_class_id}; it is 200 for an unknown id, so it checks the path only.
func (s *ShippingClasses) HeadShippingClass(ctx context.Context, shippingClassID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", ShippingClassPath, &http.RequestOptions{
		PathParams:   map[string]any{"shipping_class_id": shippingClassID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// UpdateShippingClass replaces one class — PUT /admin/shipping_classes/{shipping_class_id}; an omitted description comes back null.
func (s *ShippingClasses) UpdateShippingClass(ctx context.Context, shippingClassID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPutRequest(ctx, ShippingClassPath, &http.RequestOptions{
		PathParams: map[string]any{"shipping_class_id": shippingClassID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ShippingClassShape})
}

// DeleteShippingClass deletes one class — DELETE /admin/shipping_classes/{shipping_class_id}; the default row is 409.
func (s *ShippingClasses) DeleteShippingClass(ctx context.Context, shippingClassID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendDeleteRequest(ctx, ShippingClassPath, &http.RequestOptions{
		PathParams:   map[string]any{"shipping_class_id": shippingClassID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// SetShippingClassDefault makes one class the default — PUT /admin/shipping_classes/{shipping_class_id}/default, no body.
func (s *ShippingClasses) SetShippingClassDefault(ctx context.Context, shippingClassID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPutRequest(ctx, ShippingClassDefaultPath, &http.RequestOptions{
		PathParams:   map[string]any{"shipping_class_id": shippingClassID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ShippingClassShape})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "shipping_classes"

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

// New binds the shipping class endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *ShippingClasses { return &ShippingClasses{requester: requester} }
