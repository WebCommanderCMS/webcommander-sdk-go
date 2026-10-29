package paymentgateways

// The path templates; the client fills {api_version}.
const (
	PaymentGatewaysPath = "/api/{api_version}/admin/payment_gateways"

	PaymentGatewaysCountPath = "/api/{api_version}/admin/payment_gateways/count"

	PaymentGatewayPath = "/api/{api_version}/admin/payment_gateways/{payment_gateway_id}"
)

// PaymentGatewaysKey is the key the listing's rows arrive under.
const PaymentGatewaysKey = "payment_gateways"

// PaymentGatewayKey is the key one row arrives under; no write takes it as a wrapper.
const PaymentGatewayKey = "payment_gateway"

// PaymentGatewaysCountKey is the one key the count sub-path answers with, bare at the top of the body.
const PaymentGatewaysCountKey = "total"

// PaymentGatewaysPageSize is the page size applied when no limit is sent; a limit above PaymentGatewaysPageSizeCap is 400 out_of_range, not clamped.
const PaymentGatewaysPageSize = 20

const PaymentGatewaysPageSizeCap = 100
