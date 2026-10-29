package ordercustomfields

// The path templates; the client fills {api_version}.
const (
	OrderCustomFieldsPath = "/api/{api_version}/admin/order_custom_fields"

	OrderCustomFieldsCountPath = "/api/{api_version}/admin/order_custom_fields/count"

	OrderCustomFieldPath = "/api/{api_version}/admin/order_custom_fields/{order_custom_field_id}"

	OrderCustomFieldValidateTimeslotPath = "/api/{api_version}/admin/order_custom_fields/{order_custom_field_id}/validate_timeslot"
)

// OrderCustomFieldsKey wraps the listing's rows, and keys the listing's one group_meta entry.
const OrderCustomFieldsKey = "order_custom_fields"

// OrderCustomFieldKey is the key one row arrives under.
const OrderCustomFieldKey = "order_custom_field"

// OrderCustomFieldsCountKey is the one key the count sub-path answers with, bare at the top of the body.
const OrderCustomFieldsCountKey = "count"

// OrderCustomFieldsPageSize is the page size when limit is not sent, and the ceiling on one that is.
const OrderCustomFieldsPageSize = 20

var orderCustomFieldsPaginationKeys = [10]string{
	"count", "current_page", "has_next", "has_previous", "limit",
	"next_page", "offset", "previous_page", "total", "total_pages",
}

// OrderCustomFieldsPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func OrderCustomFieldsPaginationKeys() []string {
	return append(make([]string, 0, len(orderCustomFieldsPaginationKeys)), orderCustomFieldsPaginationKeys[:]...)
}
