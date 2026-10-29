package taxcustomprofiles

// The path templates, hyphenated: the underscore path is 404 Unknown tax path. The client fills {api_version}.
const (
	TaxCustomProfilesPath = "/api/{api_version}/admin/tax/custom-profiles"

	TaxCustomProfilePath = "/api/{api_version}/admin/tax/custom-profiles/{custom_profile_id}"
)

// TaxCustomProfilesKey is the key the listing's rows arrive under — the resource's own plural name.
const TaxCustomProfilesKey = "custom_profiles"

// TaxCustomProfileKey is the key one row arrives under, and the wrapper both writes require.
const TaxCustomProfileKey = "custom_profile"

var taxCustomProfilesQueryParams = [4]string{"limit", "offset", "page", "sort"}

// TaxCustomProfilesQueryParams returns a fresh copy of the four names the listing acts on; sort takes id or name.
func TaxCustomProfilesQueryParams() []string {
	return append(make([]string, 0, len(taxCustomProfilesQueryParams)), taxCustomProfilesQueryParams[:]...)
}

// TaxCustomProfilesPageSize is the page size when limit is not sent; limit takes 1 to 100, a larger one clamped to 100.
const TaxCustomProfilesPageSize = 20

var taxCustomProfilesPaginationKeys = [10]string{
	"count", "current_page", "has_next", "has_previous", "limit",
	"next_page", "offset", "previous_page", "total", "total_pages",
}

// TaxCustomProfilesPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func TaxCustomProfilesPaginationKeys() []string {
	return append(make([]string, 0, len(taxCustomProfilesPaginationKeys)), taxCustomProfilesPaginationKeys[:]...)
}

// UnroutedMethodBody is the 200 body the legacy underscore row path returns for an unrouted method; the hyphenated paths return 404.
func UnroutedMethodBody() map[string]any {
	return map[string]any{"isSuccess": false, "message": "Invalid API Request"}
}

// TaxCustomProfilesUnknownPathMessage opens the tax router's 404 text for a method the listing path does not route.
const TaxCustomProfilesUnknownPathMessage = "Unknown tax path"

// TaxCustomProfilesTransactionNoHeader files a new custom profile under this transaction number.
const TaxCustomProfilesTransactionNoHeader = "X-Transaction-No"
