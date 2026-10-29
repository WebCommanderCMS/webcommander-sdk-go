package taxprofiles

// The path templates; the client fills {api_version}.
const (
	TaxProfilesPath = "/api/{api_version}/admin/tax_profiles"

	TaxProfilePath = "/api/{api_version}/admin/tax_profiles/{tax_profile_id}"
)

// TaxProfilesKey is the key the listing's rows arrive under — the resource's own plural name.
const TaxProfilesKey = "tax_profiles"

// TaxProfileKey is the key one row arrives under, and the wrapper both writes require.
const TaxProfileKey = "tax_profile"

var taxProfilesQueryParams = [3]string{"limit", "offset", "page"}

// TaxProfilesQueryParams returns a fresh copy of the three names the listing acts on.
func TaxProfilesQueryParams() []string {
	return append(make([]string, 0, len(taxProfilesQueryParams)), taxProfilesQueryParams[:]...)
}

// TaxProfilesPageSize is the page size when limit is not sent, and the ceiling on one that is.
const TaxProfilesPageSize = 20

var taxProfilesPaginationKeys = [10]string{
	"count", "current_page", "has_next", "has_previous", "limit",
	"next_page", "offset", "previous_page", "total", "total_pages",
}

// TaxProfilesPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func TaxProfilesPaginationKeys() []string {
	return append(make([]string, 0, len(taxProfilesPaginationKeys)), taxProfilesPaginationKeys[:]...)
}

// UnroutedMethodBody is the body the store answers at 200, in place of a 405, for an unrouted method.
func UnroutedMethodBody() map[string]any {
	return map[string]any{"isSuccess": false, "message": "Invalid API Request"}
}
