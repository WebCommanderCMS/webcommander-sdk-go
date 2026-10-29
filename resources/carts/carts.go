package carts

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// Carts is the cart endpoints, bound to one client.
type Carts struct {
	requester http.Requester
}

// AddCartItem adds a product to the cart session_id names, creating it — POST /carts, 201; re-adding raises the line's quantity. customer_id is the customer's internal_id.
func (c *Carts) AddCartItem(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendPostRequest(ctx, CartsPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CartShape})
}

// GetCart retrieves the cart — GET /carts/{session_id}, with group_meta beside it. A cart with no line is 400 Your shopping cart is empty.
func (c *Carts) GetCart(ctx context.Context, sessionID string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, CartPath, &http.RequestOptions{
		PathParams:   map[string]any{"session_id": sessionID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CartShape})
}

// ResolveCartSession asks which session id to use — PUT /carts/{session_id}.
func (c *Carts) ResolveCartSession(ctx context.Context, sessionID string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendPutRequest(ctx, CartPath, &http.RequestOptions{
		PathParams: map[string]any{"session_id": sessionID}, JSONBody: payload,
		// The store ignores the body when Content-Type carries a charset.
		Headers:      map[string]string{"Content-Type": "application/json"},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CartSessionShape})
}

// UpdateCartItem changes one line's quantity — PATCH /carts/{session_id}. cart_item_id and quantity are required; the addresses are optional.
func (c *Carts) UpdateCartItem(ctx context.Context, sessionID string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendPatchRequest(ctx, CartPath, &http.RequestOptions{
		PathParams: map[string]any{"session_id": sessionID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CartShape})
}

// DeleteCart empties the cart — DELETE /carts/{session_id}; 204, after which a read is 400.
func (c *Carts) DeleteCart(ctx context.Context, sessionID string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendDeleteRequest(ctx, CartPath, &http.RequestOptions{
		PathParams:   map[string]any{"session_id": sessionID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// DeleteCartItem removes one line — DELETE /carts/{session_id}/cart_items/{cart_item_id}; 204. cartItemID is a line's id; removing the last line empties the cart.
func (c *Carts) DeleteCartItem(ctx context.Context, sessionID string, cartItemID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendDeleteRequest(ctx, CartItemPath, &http.RequestOptions{
		PathParams:   map[string]any{"session_id": sessionID, "cart_item_id": cartItemID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CreateOrderFromCart turns the cart into an order — POST /admin/orders/carts/{session_id}, 201, and the cart is gone. customer_id is the customer's customer_id, not its internal_id. transactionNo is sent as X-Transaction-No and files the order under that number; empty sends none.
func (c *Carts) CreateOrderFromCart(ctx context.Context, sessionID string, payload any, transactionNo string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	var headers map[string]string
	if transactionNo != "" {
		headers = map[string]string{CartsTransactionNoHeader: transactionNo}
	}
	return c.requester.SendPostRequest(ctx, CartOrderPath, &http.RequestOptions{
		PathParams: map[string]any{"session_id": sessionID}, JSONBody: payload, Headers: headers,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CartOrderShape})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "carts"

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

// New binds the cart endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *Carts {
	return &Carts{requester: requester}
}
