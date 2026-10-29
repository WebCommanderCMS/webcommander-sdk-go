package orders

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// Orders is the order endpoints, bound to one client.
type Orders struct {
	requester http.Requester
}

// ListOrders lists orders — GET /admin/orders, newest first; cancelled only under order_status. limit caps at 20, customerId takes an internal_id, an unknown name is 400.
func (o *Orders) ListOrders(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendGetRequest(ctx, OrdersPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: OrderPageShape})
}

// CountOrders counts orders — GET /admin/orders/count, {"orders": {"count": N}}; the listing's filters narrow it, its paging, sort and field_metadata change nothing.
func (o *Orders) CountOrders(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendGetRequest(ctx, OrdersCountPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: OrderCountShape})
}

// CreateOrder creates an order — POST /admin/orders. transactionNo is sent as X-Transaction-No and files the order under that number; empty sends none.
func (o *Orders) CreateOrder(ctx context.Context, payload any, transactionNo string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	var headers map[string]string
	if transactionNo != "" {
		headers = map[string]string{OrdersTransactionNoHeader: transactionNo}
	}
	return o.requester.SendPostRequest(ctx, OrdersPath, &http.RequestOptions{
		JSONBody: payload, Headers: headers, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: FlatOrderShape})
}

// CreateOrderWithPayment creates an order and records a payment on it — POST /admin/orders/with-payment.
func (o *Orders) CreateOrderWithPayment(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendPostRequest(ctx, OrdersWithPaymentPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: FlatOrderShape})
}

// GetOrder retrieves one order — GET /admin/orders/{order_id}. field_metadata adds the listing's schema block; any other parameter is ignored.
func (o *Orders) GetOrder(ctx context.Context, orderID any, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendGetRequest(ctx, OrderPath, &http.RequestOptions{
		PathParams:  map[string]any{"order_id": orderID},
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: OrderShape})
}

// UpdateOrder updates an order — PUT /admin/orders/{order_id}: the customer and the addresses, each address sent replacing the stored one. Sections left out are kept; items is 400.
func (o *Orders) UpdateOrder(ctx context.Context, orderID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendPutRequest(ctx, OrderPath, &http.RequestOptions{
		PathParams: map[string]any{"order_id": orderID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: OrderShape})
}

// PatchOrder changes named fields — PATCH /admin/orders/{order_id}.
func (o *Orders) PatchOrder(ctx context.Context, orderID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendPatchRequest(ctx, OrderPath, &http.RequestOptions{
		PathParams: map[string]any{"order_id": orderID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: OrderShape})
}

// DeleteOrder deletes an order — DELETE /admin/orders/{order_id}; hard. Read status: deleted, or not_found for an id no order holds, both at 200.
func (o *Orders) DeleteOrder(ctx context.Context, orderID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendDeleteRequest(ctx, OrderPath, &http.RequestOptions{
		PathParams:   map[string]any{"order_id": orderID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// UpdateOrderAddresses replaces one or both addresses — PATCH /admin/orders/{order_id}/addresses; an address left out is untouched, and country_code or a country object both resolve.
func (o *Orders) UpdateOrderAddresses(ctx context.Context, orderID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendPatchRequest(ctx, OrderAddressesPath, &http.RequestOptions{
		PathParams: map[string]any{"order_id": orderID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: OrderAddressesShape})
}

// CancelOrder cancels an order — POST /admin/orders/{order_id}/cancel; a nil payload sends no body.
func (o *Orders) CancelOrder(ctx context.Context, orderID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	if payload != nil && isNilValue(payload) {
		payload = nil
	}
	return o.requester.SendPostRequest(ctx, OrderCancelPath, &http.RequestOptions{
		PathParams: map[string]any{"order_id": orderID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: FlatOrderShape})
}

// ListOrderComments reads an order's comments, oldest first and unpaged — GET /admin/orders/{order_id}/comments.
func (o *Orders) ListOrderComments(ctx context.Context, orderID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendGetRequest(ctx, OrderCommentsPath, &http.RequestOptions{
		PathParams:   map[string]any{"order_id": orderID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: OrderCommentListShape})
}

// CreateOrderComment adds a comment — POST /admin/orders/{order_id}/comments; 201. A bare record: a comment wrapper hides content and is 400.
func (o *Orders) CreateOrderComment(ctx context.Context, orderID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendPostRequest(ctx, OrderCommentsPath, &http.RequestOptions{
		PathParams: map[string]any{"order_id": orderID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: OrderCommentShape})
}

// GetOrderCustomFields reads an order's custom field values — GET /admin/orders/{order_internal_id}/custom_fields.
func (o *Orders) GetOrderCustomFields(ctx context.Context, orderInternalID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendGetRequest(ctx, OrderCustomFieldsPath, &http.RequestOptions{
		PathParams:   map[string]any{"order_internal_id": orderInternalID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: OrderCustomFieldsShape})
}

// CreateOrderPayment records a payment — POST /admin/orders/{order_id}/payments; 201. The payment wrapper is required, and the receipt's id is what a payment refund takes.
func (o *Orders) CreateOrderPayment(ctx context.Context, orderID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendPostRequest(ctx, OrderPaymentsPath, &http.RequestOptions{
		PathParams: map[string]any{"order_id": orderID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: PaymentReceiptShape})
}

// ListOrderRefunds lists an order's refunds, newest first — GET /admin/orders/{order_id}/refunds, by the order's order_id. limit is 1 to 20 (above 20 is capped, 0 is 400), with offset or a 1-based page; any other name is ignored.
func (o *Orders) ListOrderRefunds(ctx context.Context, orderID any, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendGetRequest(ctx, OrderRefundsPath, &http.RequestOptions{
		PathParams:  map[string]any{"order_id": orderID},
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: OrderRefundPageShape})
}

// CreateOrderRefund refunds against an order — POST /admin/orders/{order_id}/refunds; 201. It raises refunded, leaves paid alone, and cannot pass what was paid less earlier refunds.
func (o *Orders) CreateOrderRefund(ctx context.Context, orderID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendPostRequest(ctx, OrderRefundsPath, &http.RequestOptions{
		PathParams: map[string]any{"order_id": orderID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: OrderRefundShape})
}

// CreatePaymentRefund refunds part of a payment — POST /admin/payment/{payment_id}/refunds; 200. Refused without idempotencyKey; a repeat under the same key replays the first refund.
func (o *Orders) CreatePaymentRefund(ctx context.Context, paymentID any, payload any, idempotencyKey string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendPostRequest(ctx, PaymentRefundsPath, &http.RequestOptions{
		PathParams: map[string]any{"payment_id": paymentID}, JSONBody: payload,
		Headers:      map[string]string{OrdersIdempotencyKeyHeader: idempotencyKey},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: PaymentRefundShape})
}

// SendOrderInvoice emails the invoice — POST /admin/orders/{order_internal_id}/send-invoice; no body. Pass internal_id: an order_id there names the order holding it as its internal_id.
func (o *Orders) SendOrderInvoice(ctx context.Context, orderInternalID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendPostRequest(ctx, OrderSendInvoicePath, &http.RequestOptions{
		PathParams:   map[string]any{"order_internal_id": orderInternalID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: InvoiceDispatchShape})
}

// GetOrderShipmentStatus reads how far an order has shipped — GET /admin/orders/{order_id}/shipment-status.
func (o *Orders) GetOrderShipmentStatus(ctx context.Context, orderID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendGetRequest(ctx, OrderShipmentStatusPath, &http.RequestOptions{
		PathParams:   map[string]any{"order_id": orderID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ShipmentStatusShape})
}

// ListOrderShipments lists an order's shipments — GET /admin/orders/{order_id}/shipments; unpaged.
func (o *Orders) ListOrderShipments(ctx context.Context, orderID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendGetRequest(ctx, OrderShipmentsPath, &http.RequestOptions{
		PathParams:   map[string]any{"order_id": orderID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: OrderShipmentListShape})
}

// CreateOrderShipment ships order lines — POST /admin/orders/{order_id}/shipments; 201. method, shipping_date and shipment_items are required; tracking_info may be left out.
func (o *Orders) CreateOrderShipment(ctx context.Context, orderID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendPostRequest(ctx, OrderShipmentsPath, &http.RequestOptions{
		PathParams: map[string]any{"order_id": orderID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: OrderShipmentShape})
}

// UpdateOrderShipment changes a shipment — PATCH /admin/orders/{order_id}/shipments/{shipment_id}, by its shipping_id; a field left out is kept, and shipment_items needs change_note.
func (o *Orders) UpdateOrderShipment(ctx context.Context, orderID any, shipmentID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendPatchRequest(ctx, OrderShipmentPath, &http.RequestOptions{
		PathParams: map[string]any{"order_id": orderID, "shipment_id": shipmentID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: OrderShipmentShape})
}

// UpdateOrderStatus moves an order's status — PATCH /admin/orders/{order_id}/update-status; answered with the flat order shape. A completed or cancelled order can be moved again.
func (o *Orders) UpdateOrderStatus(ctx context.Context, orderID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendPatchRequest(ctx, OrderUpdateStatusPath, &http.RequestOptions{
		PathParams: map[string]any{"order_id": orderID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: FlatOrderShape})
}

// CreateOrderStatus adds a custom order status — POST /admin/order/add-order-status; 201 {"success": "true"}. No endpoint lists, applies or deletes one, and the order keeps its status.
func (o *Orders) CreateOrderStatus(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendPostRequest(ctx, OrderStatusAddPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ListCustomerOrders reads one customer's orders — GET /admin/customers/{customer_id}/orders, by the public customer_id; unpaged, no parameter read, and cancelled orders left out.
func (o *Orders) ListCustomerOrders(ctx context.Context, customerID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return o.requester.SendGetRequest(ctx, CustomerOrdersPath, &http.RequestOptions{
		PathParams:   map[string]any{"customer_id": customerID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CustomerOrdersShape})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "orders"

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

// New binds the order endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *Orders {
	return &Orders{requester: requester}
}
