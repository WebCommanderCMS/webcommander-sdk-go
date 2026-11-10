package shippingrates

// The path templates; the client fills {api_version}.
const (
	ShippingRatesPath = "/api/{api_version}/admin/shipping_rates"

	ShippingRatesCountPath = "/api/{api_version}/admin/shipping_rates/count"

	ShippingRatePath = "/api/{api_version}/admin/shipping_rates/{shipping_rate_id}"
)

// ShippingRatesKey is the key the listing's rows arrive under.
const ShippingRatesKey = "shipping_rates"

// ShippingRateKey is the key one row arrives under, and the wrapper every write takes; a bare body is 400.
const ShippingRateKey = "shipping_rate"

// ShippingRatesCountKey is the one key the count sub-path answers with, bare at the top of the body.
const ShippingRatesCountKey = "total"

// ShippingRatesPageSize is the page size applied when no limit is sent; limit takes 1 to ShippingRatesPageSizeCap, and 0, above the cap or "all" is 400.
const ShippingRatesPageSize = 20

const ShippingRatesPageSizeCap = 100

// ShippingRatesTransactionNoHeader files a new rate under this transaction number and wins over a transaction_no in the body.
const ShippingRatesTransactionNoHeader = "X-Transaction-No"
