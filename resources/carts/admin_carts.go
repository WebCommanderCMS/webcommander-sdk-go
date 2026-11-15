package carts

import (
	"context"

	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/orders"
)

// AdminCarts is the admin cart and checkout endpoints, bound to one client. Storefront carts are Carts.
type AdminCarts struct {
	requester http.Requester
}

// CreateAdminCart starts a cart for a customer — POST /admin/orders/carts, 201. customer_id is the customer's customer_id and products (product_id and quantity) is required. Keep cart.session.session_id: every other call takes it as sessionID. The session lasts 30 minutes from the last write.
func (a *AdminCarts) CreateAdminCart(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendPostRequest(ctx, AdminCartsPath, &http.RequestOptions{
		JSONBody:     payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: AdminCartShape})
}

// GetAdminCart reads the cart — GET /admin/orders/carts/{session_id}, where the slot takes cart.session.session_id. A deleted cart is 410 session_deleted and one an order was placed from is 410 session_confirmed.
func (a *AdminCarts) GetAdminCart(ctx context.Context, sessionID string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendGetRequest(ctx, AdminCartPath, &http.RequestOptions{
		PathParams:   map[string]any{"session_id": sessionID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: AdminCartShape})
}

// ReplaceAdminCart replaces the cart's customer and lines — PUT /admin/orders/carts/{session_id}, 200. customer_id and products are both sent. The lines are put on again, so each gets a new cart_item_id.
func (a *AdminCarts) ReplaceAdminCart(ctx context.Context, sessionID string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendPutRequest(ctx, AdminCartPath, &http.RequestOptions{
		PathParams:   map[string]any{"session_id": sessionID},
		JSONBody:     payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: AdminCartShape})
}

// UpdateAdminCart changes the cart's customer — PATCH /admin/orders/carts/{session_id}, 200. Only customer_id is taken, and it is required: products is 400 unknown_field.
func (a *AdminCarts) UpdateAdminCart(ctx context.Context, sessionID string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendPatchRequest(ctx, AdminCartPath, &http.RequestOptions{
		PathParams:   map[string]any{"session_id": sessionID},
		JSONBody:     payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: AdminCartShape})
}

// DeleteAdminCart deletes the cart — DELETE /admin/orders/carts/{session_id}; 204 with no body, after which every call on the session is 410.
func (a *AdminCarts) DeleteAdminCart(ctx context.Context, sessionID string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendDeleteRequest(ctx, AdminCartPath, &http.RequestOptions{
		PathParams:   map[string]any{"session_id": sessionID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// AddAdminCartItems puts products on the cart — POST /admin/orders/carts/{session_id}/cart_items, 201 with the whole cart. products is required; one already on the cart raises that line's quantity.
func (a *AdminCarts) AddAdminCartItems(ctx context.Context, sessionID string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendPostRequest(ctx, AdminCartItemsPath, &http.RequestOptions{
		PathParams:   map[string]any{"session_id": sessionID},
		JSONBody:     payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: AdminCartShape})
}

// ClearAdminCartItems takes every line off the cart — DELETE /admin/orders/carts/{session_id}/cart_items, 200 with the empty cart. The session stays.
func (a *AdminCarts) ClearAdminCartItems(ctx context.Context, sessionID string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendDeleteRequest(ctx, AdminCartItemsPath, &http.RequestOptions{
		PathParams:   map[string]any{"session_id": sessionID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: AdminCartShape})
}

// UpdateAdminCartItem sets one line's quantity — PATCH /admin/orders/carts/{session_id}/cart_items/{cart_item_id}, 200 with the cart. cartItemID is a line's cart_item_id and quantity is required. A line the cart does not hold is 404.
func (a *AdminCarts) UpdateAdminCartItem(ctx context.Context, sessionID string, cartItemID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendPatchRequest(ctx, AdminCartItemPath, &http.RequestOptions{
		PathParams:   map[string]any{"session_id": sessionID, "cart_item_id": cartItemID},
		JSONBody:     payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: AdminCartShape})
}

// DeleteAdminCartItem takes one line off the cart — DELETE /admin/orders/carts/{session_id}/cart_items/{cart_item_id}, 200 with the cart. cartItemID is a line's cart_item_id; one the cart does not hold is 404.
func (a *AdminCarts) DeleteAdminCartItem(ctx context.Context, sessionID string, cartItemID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendDeleteRequest(ctx, AdminCartItemPath, &http.RequestOptions{
		PathParams:   map[string]any{"session_id": sessionID, "cart_item_id": cartItemID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: AdminCartShape})
}

// GetCheckout reads the checkout of a cart — GET /admin/orders/checkout/{session_id}. issues lists what still blocks a confirm, delivery_type required until one is set; the addresses start as the customer's.
func (a *AdminCarts) GetCheckout(ctx context.Context, sessionID string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendGetRequest(ctx, AdminCheckoutPath, &http.RequestOptions{
		PathParams:   map[string]any{"session_id": sessionID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: AdminCheckoutShape})
}

// UpdateCheckout sets the delivery and addresses — PATCH /admin/orders/checkout/{session_id}, 200 with the priced checkout. Only delivery_type, billing_address and shipping_address are taken. The shipping and handling cost follow the shipping address, so set it with delivery_type.
func (a *AdminCarts) UpdateCheckout(ctx context.Context, sessionID string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendPatchRequest(ctx, AdminCheckoutPath, &http.RequestOptions{
		PathParams:   map[string]any{"session_id": sessionID},
		JSONBody:     payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: AdminCheckoutShape})
}

// ConfirmCheckout places the order — POST /admin/orders/checkout/{session_id}/confirm; 201 with the order, and the body optional (nil sends none). This creates a real order and closes the session, so every later call on it is 410 session_confirmed. A checkout with issues is 422 checkout_incomplete. Without a payment the order is unpaid. transactionNo is sent as X-Transaction-No and files the order under that number; empty sends none.
func (a *AdminCarts) ConfirmCheckout(ctx context.Context, sessionID string, payload any, transactionNo string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	if payload != nil && isNilValue(payload) {
		payload = nil
	}
	var headers map[string]string
	if transactionNo != "" {
		headers = map[string]string{AdminCartsTransactionNoHeader: transactionNo}
	}
	return a.requester.SendPostRequest(ctx, AdminCheckoutConfirmPath, &http.RequestOptions{
		PathParams: map[string]any{"session_id": sessionID}, JSONBody: payload, Headers: headers,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: orders.FlatOrderShape})
}

// ApplyCheckoutCoupon applies a coupon — POST /admin/orders/checkout/{session_id}/coupon, 200. code is required. One coupon at a time: a second is 409, and a code the store does not accept is 422 invalid_coupon.
func (a *AdminCarts) ApplyCheckoutCoupon(ctx context.Context, sessionID string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendPostRequest(ctx, AdminCheckoutCouponPath, &http.RequestOptions{
		PathParams:   map[string]any{"session_id": sessionID},
		JSONBody:     payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: AdminCheckoutShape})
}

// RemoveCheckoutCoupon removes the coupon — DELETE /admin/orders/checkout/{session_id}/coupon. With none applied it is 404 coupon_not_applied.
func (a *AdminCarts) RemoveCheckoutCoupon(ctx context.Context, sessionID string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendDeleteRequest(ctx, AdminCheckoutCouponPath, &http.RequestOptions{
		PathParams:   map[string]any{"session_id": sessionID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: AdminCheckoutShape})
}

// ApplyCheckoutGiftCard applies a gift card — POST /admin/orders/checkout/{session_id}/gift_card, 200. code is required. The card must be usable by the cart's customer, else 422 invalid_gift_card; a second card is 409. It lowers due and spends nothing until an order is placed.
func (a *AdminCarts) ApplyCheckoutGiftCard(ctx context.Context, sessionID string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendPostRequest(ctx, AdminCheckoutGiftCardPath, &http.RequestOptions{
		PathParams:   map[string]any{"session_id": sessionID},
		JSONBody:     payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: AdminCheckoutShape})
}

// RemoveCheckoutGiftCard removes the gift card — DELETE /admin/orders/checkout/{session_id}/gift_card. With none applied it is 404 gift_card_not_applied.
func (a *AdminCarts) RemoveCheckoutGiftCard(ctx context.Context, sessionID string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendDeleteRequest(ctx, AdminCheckoutGiftCardPath, &http.RequestOptions{
		PathParams:   map[string]any{"session_id": sessionID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: AdminCheckoutShape})
}

// GetCheckoutDeliverySchedules reads the chosen delivery schedule — GET /admin/orders/checkout/{session_id}/delivery_schedules, with session and delivery_schedule, which is nil until one is set.
func (a *AdminCarts) GetCheckoutDeliverySchedules(ctx context.Context, sessionID string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendGetRequest(ctx, AdminCheckoutDeliverySchedulesPath, &http.RequestOptions{
		PathParams:   map[string]any{"session_id": sessionID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: AdminCheckoutShape})
}

// SetCheckoutDeliverySchedule chooses a delivery schedule and date — PUT /admin/orders/checkout/{session_id}/delivery_schedules. id and date are both required. Only a delivery_type of delivery takes one (else 422), an unknown id is 404 and one that does not apply to the order is 422 schedule_not_applicable.
func (a *AdminCarts) SetCheckoutDeliverySchedule(ctx context.Context, sessionID string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendPutRequest(ctx, AdminCheckoutDeliverySchedulesPath, &http.RequestOptions{
		PathParams:   map[string]any{"session_id": sessionID},
		JSONBody:     payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: AdminCheckoutShape})
}

// ClearCheckoutDeliverySchedule clears the delivery schedule — DELETE /admin/orders/checkout/{session_id}/delivery_schedules; 204 with no body, and also when none was set.
func (a *AdminCarts) ClearCheckoutDeliverySchedule(ctx context.Context, sessionID string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendDeleteRequest(ctx, AdminCheckoutDeliverySchedulesPath, &http.RequestOptions{
		PathParams:   map[string]any{"session_id": sessionID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// AdminCartsResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const AdminCartsResourceName = "admin_carts"

// NewAdminCarts binds the admin cart and checkout endpoints to any requester (the WebCommanderClient).
func NewAdminCarts(requester http.Requester) *AdminCarts {
	return &AdminCarts{requester: requester}
}
