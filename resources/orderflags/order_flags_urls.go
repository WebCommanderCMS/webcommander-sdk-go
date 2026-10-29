package orderflags

// The path templates for the order flag endpoints; the client fills {api_version}.
const (
	OrderFlagsPath = "/api/{api_version}/admin/order_flags"

	OrderFlagPath = "/api/{api_version}/admin/order_flags/{order_flag_id}"
)

// OrderFlagsKey is the key the listing's rows arrive under — the resource's own plural name.
const OrderFlagsKey = "order_flags"

// OrderFlagKey is the key one row arrives under, and the wrapper both writes require.
const OrderFlagKey = "order_flag"

// orderFlagsQueryParams are the four names the listing accepts; any other is 400.
var orderFlagsQueryParams = [4]string{"limit", "offset", "page", "field_metadata"}

// OrderFlagsQueryParams returns a fresh copy of the four names the listing accepts.
func OrderFlagsQueryParams() []string {
	return append(make([]string, 0, len(orderFlagsQueryParams)), orderFlagsQueryParams[:]...)
}
