package discountprofiles

// The path templates; the client fills {api_version}.
const (
	DiscountProfilesPath = "/api/{api_version}/admin/discount_profiles"

	DiscountProfilePath = "/api/{api_version}/admin/discount_profiles/{discount_profile_id}"

	DiscountProfileCreateCouponsPath = "/api/{api_version}/admin/discount_profiles/{discount_profile_id}/create-coupons"

	DiscountCouponsAssignCustomerPath = "/api/{api_version}/admin/discount_coupons/assign-customer"
)

// DiscountProfilesKey is the key the listing's rows arrive under — the resource's own plural name.
const DiscountProfilesKey = "discount_profiles"

// DiscountProfileKey is the key one row arrives under; no write takes it, as every write takes a bare body.
const DiscountProfileKey = "discount_profile"

// DiscountProfileCouponCodesKey is the key the generated codes arrive under, beside status and message.
const DiscountProfileCouponCodesKey = "codes"

var discountDetailsKeys = [3]string{"amount_details", "shipping_details", "product_details"}

// DiscountDetailsKeys returns a fresh copy of the three details blocks; a row carries the one its discount_details_type names.
func DiscountDetailsKeys() []string {
	return append(make([]string, 0, len(discountDetailsKeys)), discountDetailsKeys[:]...)
}

var discountProfilesQueryParams = [8]string{
	"limit", "offset", "page", "sort", "dir", "type", "name", "field_metadata",
}

// DiscountProfilesQueryParams returns a fresh copy of the names the listing reads; any other is ignored, not refused.
func DiscountProfilesQueryParams() []string {
	return append(make([]string, 0, len(discountProfilesQueryParams)), discountProfilesQueryParams[:]...)
}

var discountProfilesSortFields = [3]string{"id", "name", "type"}

// DiscountProfilesSortFields returns a fresh copy of the fields sort accepts; any other is 400 invalid_sort_field.
func DiscountProfilesSortFields() []string {
	return append(make([]string, 0, len(discountProfilesSortFields)), discountProfilesSortFields[:]...)
}

var discountProfileTypes = [4]string{
	"customer", "promote_product", "offer_incentive_and_sell_more", "offer_coupon",
}

// DiscountProfileTypes returns a fresh copy of the values type takes; a create and the listing filter refuse any other, an update does not.
func DiscountProfileTypes() []string {
	return append(make([]string, 0, len(discountProfileTypes)), discountProfileTypes[:]...)
}

// DiscountProfilesPageSize is the page size when limit is not sent, and the ceiling on one that is.
const DiscountProfilesPageSize = 20

var discountProfilesPaginationKeys = [10]string{
	"count", "current_page", "has_next", "has_previous", "limit",
	"next_page", "offset", "previous_page", "total", "total_pages",
}

// DiscountProfilesPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func DiscountProfilesPaginationKeys() []string {
	return append(make([]string, 0, len(discountProfilesPaginationKeys)), discountProfilesPaginationKeys[:]...)
}
