package orders

// The path templates; the client fills {api_version}.
const (
	OrdersPath = "/api/{api_version}/admin/orders"

	OrdersCountPath = "/api/{api_version}/admin/orders/count"

	OrdersWithPaymentPath = "/api/{api_version}/admin/orders/with-payment"

	OrderPath = "/api/{api_version}/admin/orders/{order_id}"

	OrderAddressesPath = "/api/{api_version}/admin/orders/{order_id}/addresses"

	OrderCancelPath = "/api/{api_version}/admin/orders/{order_id}/cancel"

	OrderCommentsPath = "/api/{api_version}/admin/orders/{order_id}/comments"

	OrderPaymentsPath = "/api/{api_version}/admin/orders/{order_id}/payments"

	OrderRefundsPath = "/api/{api_version}/admin/orders/{order_id}/refunds"

	OrderShipmentStatusPath = "/api/{api_version}/admin/orders/{order_id}/shipment-status"

	OrderShipmentsPath = "/api/{api_version}/admin/orders/{order_id}/shipments"

	OrderShipmentPath = "/api/{api_version}/admin/orders/{order_id}/shipments/{shipment_id}"

	OrderUpdateStatusPath = "/api/{api_version}/admin/orders/{order_id}/update-status"

	// These two resolve the slot against the order's internal_id; an order_id there names another order, or none.
	OrderCustomFieldsPath = "/api/{api_version}/admin/orders/{order_internal_id}/custom_fields"
	OrderSendInvoicePath  = "/api/{api_version}/admin/orders/{order_internal_id}/send-invoice"

	PaymentRefundsPath = "/api/{api_version}/admin/payment/{payment_id}/refunds"

	OrderStatusAddPath = "/api/{api_version}/admin/order/add-order-status"

	CustomerOrdersPath = "/api/{api_version}/admin/customers/{customer_id}/orders"
)

// The body keys of the order endpoints.
const (
	OrdersKey         = "orders"
	OrderKey          = "order"
	OrderCommentsKey  = "comments"
	OrderCommentKey   = "comment"
	OrderShipmentsKey = "shipments"
	OrderShipmentKey  = "shipment"
	OrderPaymentKey   = "payment"
)

// OrdersPageSize is the listing's page size when limit is absent, and its ceiling when limit is larger.
const OrdersPageSize = 20

var ordersSortFields = [10]string{
	"orderNo", "customerName", "created", "orderStatus", "paymentStatus",
	"shippingStatus", "deliveryType", "scheduledDate", "total", "items",
}

// OrdersSortFields returns a fresh copy of the ten values the listing's sort refusal names; id is also accepted and changes nothing.
func OrdersSortFields() []string {
	return append(make([]string, 0, len(ordersSortFields)), ordersSortFields[:]...)
}

// OrdersIdempotencyKeyHeader is the header a payment refund is refused without; a repeat under one key replays the first refund.
const OrdersIdempotencyKeyHeader = "Idempotency-Key"

// OrdersTransactionNoHeader files a created order under this transaction number; it beats a transaction_no body field.
const OrdersTransactionNoHeader = "X-Transaction-No"
