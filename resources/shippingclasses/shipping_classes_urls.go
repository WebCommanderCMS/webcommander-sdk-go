package shippingclasses

// The path templates; the client fills {api_version}.
const (
	ShippingClassesPath = "/api/{api_version}/admin/shipping_classes"

	ShippingClassesCountPath = "/api/{api_version}/admin/shipping_classes/count"

	ShippingClassPath = "/api/{api_version}/admin/shipping_classes/{shipping_class_id}"

	ShippingClassDefaultPath = "/api/{api_version}/admin/shipping_classes/{shipping_class_id}/default"
)

// The body keys: the listing's rows and the count's record, one row and both writes, and the count itself.
const (
	ShippingClassesKey      = "shipping_classes"
	ShippingClassKey        = "shipping_class"
	ShippingClassesCountKey = "count"
)

var shippingClassesQueryParams = [4]string{"limit", "offset", "page", "field_metadata"}

// ShippingClassesQueryParams returns a fresh copy of the names the listing accepts; others are 400.
func ShippingClassesQueryParams() []string {
	return append(make([]string, 0, len(shippingClassesQueryParams)), shippingClassesQueryParams[:]...)
}

// ShippingClassesPageSize is the page size when limit is not sent, and the ceiling on one that is.
const ShippingClassesPageSize = 20

var shippingClassesPaginationKeys = [10]string{
	"count", "current_page", "has_next", "has_previous", "limit",
	"next_page", "offset", "previous_page", "total", "total_pages",
}

// ShippingClassesPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func ShippingClassesPaginationKeys() []string {
	return append(make([]string, 0, len(shippingClassesPaginationKeys)), shippingClassesPaginationKeys[:]...)
}

// UnroutedMethodBody is the 200 body a method with no handler answers, such as PATCH; nothing is stored.
func UnroutedMethodBody() map[string]any {
	return map[string]any{"isSuccess": false, "message": "Invalid API Request"}
}
