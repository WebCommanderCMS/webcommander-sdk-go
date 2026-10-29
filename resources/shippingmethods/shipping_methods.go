package shippingmethods

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// ShippingMethods is the shipping method endpoints, bound to one client.
type ShippingMethods struct {
	requester http.Requester
}

// ListShippingMethods lists the methods — GET /admin/shipping_methods; limit is clamped and an unknown query name is 400.
func (s *ShippingMethods) ListShippingMethods(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, ShippingMethodsPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ShippingMethodPageShape})
}

// HeadShippingMethods asks whether the collection is reachable — HEAD /admin/shipping_methods.
func (s *ShippingMethods) HeadShippingMethods(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", ShippingMethodsPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CountShippingMethods is how many methods the store holds — GET /admin/shipping_methods/count.
func (s *ShippingMethods) CountShippingMethods(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, ShippingMethodsCountPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ShippingMethodCountShape})
}

// CreateShippingMethod creates a method — POST /admin/shipping_methods, wrapped in shipping_method; name is required and unique.
func (s *ShippingMethods) CreateShippingMethod(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPostRequest(ctx, ShippingMethodsPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ShippingMethodShape})
}

// GetShippingMethod retrieves one method — GET /admin/shipping_methods/{shipping_method_id}.
func (s *ShippingMethods) GetShippingMethod(ctx context.Context, shippingMethodID any, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, ShippingMethodPath, &http.RequestOptions{
		PathParams:  map[string]any{"shipping_method_id": shippingMethodID},
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ShippingMethodShape})
}

// HeadShippingMethod answers 200 even for an unknown id — HEAD /admin/shipping_methods/{shipping_method_id}.
func (s *ShippingMethods) HeadShippingMethod(ctx context.Context, shippingMethodID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", ShippingMethodPath, &http.RequestOptions{
		PathParams:   map[string]any{"shipping_method_id": shippingMethodID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// UpdateShippingMethod renames one method — PUT /admin/shipping_methods/{shipping_method_id}; name is the one field read.
func (s *ShippingMethods) UpdateShippingMethod(ctx context.Context, shippingMethodID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPutRequest(ctx, ShippingMethodPath, &http.RequestOptions{
		PathParams: map[string]any{"shipping_method_id": shippingMethodID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ShippingMethodShape})
}

// DeleteShippingMethod deletes a method — DELETE /admin/shipping_methods/{shipping_method_id}; the default row is 409.
func (s *ShippingMethods) DeleteShippingMethod(ctx context.Context, shippingMethodID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendDeleteRequest(ctx, ShippingMethodPath, &http.RequestOptions{
		PathParams:   map[string]any{"shipping_method_id": shippingMethodID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// SetShippingMethodDefault moves the default flag to one method — PUT /admin/shipping_methods/{shipping_method_id}/default.
func (s *ShippingMethods) SetShippingMethodDefault(ctx context.Context, shippingMethodID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPutRequest(ctx, ShippingMethodDefaultPath, &http.RequestOptions{
		PathParams:   map[string]any{"shipping_method_id": shippingMethodID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ShippingMethodShape})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "shipping_methods"

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

// New binds the shipping method endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *ShippingMethods { return &ShippingMethods{requester: requester} }
