package shippingmethods

// The path templates; the client fills {api_version}.
const (
	ShippingMethodsPath = "/api/{api_version}/admin/shipping_methods"

	ShippingMethodsCountPath = "/api/{api_version}/admin/shipping_methods/count"

	ShippingMethodPath = "/api/{api_version}/admin/shipping_methods/{shipping_method_id}"

	ShippingMethodDefaultPath = "/api/{api_version}/admin/shipping_methods/{shipping_method_id}/default"
)

// ShippingMethodsKey wraps the listing's rows, and the count sub-path's record.
const ShippingMethodsKey = "shipping_methods"

// ShippingMethodKey is the key one row arrives under, and the wrapper both writes require.
const ShippingMethodKey = "shipping_method"

// ShippingMethodsCountKey is the one key the count record carries.
const ShippingMethodsCountKey = "count"

var shippingMethodsQueryParams = [4]string{"limit", "offset", "page", "field_metadata"}

// ShippingMethodsQueryParams returns a fresh copy of the four names the listing accepts; others are 400.
func ShippingMethodsQueryParams() []string {
	return append(make([]string, 0, len(shippingMethodsQueryParams)), shippingMethodsQueryParams[:]...)
}

// ShippingMethodsPageSize is the page size when limit is not sent, and the ceiling on one that is.
const ShippingMethodsPageSize = 20

var shippingMethodsPaginationKeys = [10]string{
	"count", "current_page", "has_next", "has_previous", "limit",
	"next_page", "offset", "previous_page", "total", "total_pages",
}

// ShippingMethodsPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func ShippingMethodsPaginationKeys() []string {
	return append(make([]string, 0, len(shippingMethodsPaginationKeys)), shippingMethodsPaginationKeys[:]...)
}

// UnroutedMethodBody is the body the store answers at 200, in place of a 405, for an unrouted method.
func UnroutedMethodBody() map[string]any {
	return map[string]any{"isSuccess": false, "message": "Invalid API Request"}
}
