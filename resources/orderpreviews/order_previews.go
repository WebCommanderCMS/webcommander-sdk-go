package orderpreviews

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/orders"
)

// OrderPreviews is the order preview endpoints, bound to one client.
type OrderPreviews struct {
	requester http.Requester
}

// CreateOrderPreview starts an order preview — POST /admin/orders/preview, 201, a priced session that places no order. products is required; read order.session.session_id, which every other call takes as previewSessionID. The session lasts 30 minutes.
func (p *OrderPreviews) CreateOrderPreview(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendPostRequest(ctx, OrderPreviewsPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: OrderPreviewShape})
}

// GetOrderPreview reads an order preview — GET /admin/orders/preview/{preview_session_id}, where the slot takes the session_id the start answered. A read does not extend the session. An unknown id is 404, and a session already confirmed is 410 session_confirmed.
func (p *OrderPreviews) GetOrderPreview(ctx context.Context, previewSessionID string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendGetRequest(ctx, OrderPreviewPath, &http.RequestOptions{
		PathParams:   map[string]any{"preview_session_id": previewSessionID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: OrderPreviewShape})
}

// ConfirmOrderPreview places the order — POST /admin/orders/preview/{preview_session_id}/confirm; 201 with the created order, and the body optional (nil sends none). This creates a real order and closes the session, so a second confirm, or any other call on it, is 410 session_confirmed. Without a payment the order is unpaid, and so it is after a payment_gateway alone.
func (p *OrderPreviews) ConfirmOrderPreview(ctx context.Context, previewSessionID string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	if payload != nil && isNilValue(payload) {
		payload = nil
	}
	return p.requester.SendPostRequest(ctx, OrderPreviewConfirmPath, &http.RequestOptions{
		PathParams: map[string]any{"preview_session_id": previewSessionID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: orders.FlatOrderShape})
}

// ApplyOrderPreviewCoupon applies a coupon — POST /admin/orders/preview/{preview_session_id}/coupon, 200 with session, discounts and totals. One coupon at a time: a second is 409, and a code the store does not accept is 422 invalid_coupon.
func (p *OrderPreviews) ApplyOrderPreviewCoupon(ctx context.Context, previewSessionID string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendPostRequest(ctx, OrderPreviewCouponPath, &http.RequestOptions{
		PathParams: map[string]any{"preview_session_id": previewSessionID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: OrderPreviewShape})
}

// RemoveOrderPreviewCoupon removes the coupon — DELETE /admin/orders/preview/{preview_session_id}/coupon, 204 with no body. With none applied it is 404 coupon_not_applied.
func (p *OrderPreviews) RemoveOrderPreviewCoupon(ctx context.Context, previewSessionID string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendDeleteRequest(ctx, OrderPreviewCouponPath, &http.RequestOptions{
		PathParams:   map[string]any{"preview_session_id": previewSessionID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ApplyOrderPreviewGiftCard applies a gift card — POST /admin/orders/preview/{preview_session_id}/gift-card, 200 with session, discounts and totals. The card must be the previewed customer's, active, paid, in date and with a balance, else 422 invalid_gift_card; a second card is 409. It lowers due and spends nothing until an order is placed with it applied.
func (p *OrderPreviews) ApplyOrderPreviewGiftCard(ctx context.Context, previewSessionID string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendPostRequest(ctx, OrderPreviewGiftCardPath, &http.RequestOptions{
		PathParams: map[string]any{"preview_session_id": previewSessionID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: OrderPreviewShape})
}

// RemoveOrderPreviewGiftCard removes the gift card — DELETE /admin/orders/preview/{preview_session_id}/gift-card, 204 with no body. With none applied it is 404 gift_card_not_applied.
func (p *OrderPreviews) RemoveOrderPreviewGiftCard(ctx context.Context, previewSessionID string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendDeleteRequest(ctx, OrderPreviewGiftCardPath, &http.RequestOptions{
		PathParams:   map[string]any{"preview_session_id": previewSessionID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "order_previews"

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

// New binds the order preview endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *OrderPreviews {
	return &OrderPreviews{requester: requester}
}
