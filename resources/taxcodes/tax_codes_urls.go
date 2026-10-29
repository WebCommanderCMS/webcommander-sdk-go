package taxcodes

// The path templates; the client fills {api_version}.
const (
	TaxCodesPath = "/api/{api_version}/admin/tax_codes"

	TaxCodePath = "/api/{api_version}/admin/tax_codes/{tax_code_id}"
)

// TaxCodesKey is the key the listing's rows arrive under — the resource's own plural name.
const TaxCodesKey = "tax_codes"

// TaxCodeKey is the key one row arrives under, and the wrapper both writes require.
const TaxCodeKey = "tax_code"

var taxCodesQueryParams = [3]string{"limit", "offset", "page"}

// TaxCodesQueryParams returns a fresh copy of the three names the listing acts on.
func TaxCodesQueryParams() []string {
	return append(make([]string, 0, len(taxCodesQueryParams)), taxCodesQueryParams[:]...)
}

// TaxCodesPageSize is the page size when limit is not sent; a limit above 100 is held to 100.
const TaxCodesPageSize = 20

var taxCodesPaginationKeys = [10]string{
	"count", "current_page", "has_next", "has_previous", "limit",
	"next_page", "offset", "previous_page", "total", "total_pages",
}

// TaxCodesPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func TaxCodesPaginationKeys() []string {
	return append(make([]string, 0, len(taxCodesPaginationKeys)), taxCodesPaginationKeys[:]...)
}

// UnroutedMethodBody is the body the store answers at 200, in place of a 405, for an unrouted method.
func UnroutedMethodBody() map[string]any {
	return map[string]any{"isSuccess": false, "message": "Invalid API Request"}
}
