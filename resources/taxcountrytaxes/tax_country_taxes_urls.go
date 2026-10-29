package taxcountrytaxes

// The path templates, hyphenated: the underscore path is 404 Unknown tax path. The client fills {api_version}.
const (
	TaxCountryTaxesPath = "/api/{api_version}/admin/tax/country-taxes"

	TaxCountryTaxPath = "/api/{api_version}/admin/tax/country-taxes/{country_tax_id}"
)

// TaxCountryTaxesKey is the key the listing's rows arrive under — the resource's own plural name.
const TaxCountryTaxesKey = "country_taxes"

// TaxCountryTaxKey is the key one row arrives under, and the wrapper both writes require.
const TaxCountryTaxKey = "country_tax"

var taxCountryTaxesQueryParams = [4]string{"limit", "offset", "page", "sort"}

// TaxCountryTaxesQueryParams returns a fresh copy of the four names the listing acts on; sort takes id or name.
func TaxCountryTaxesQueryParams() []string {
	return append(make([]string, 0, len(taxCountryTaxesQueryParams)), taxCountryTaxesQueryParams[:]...)
}

// TaxCountryTaxesPageSize is the page size when limit is not sent; limit takes 1 to 100, a larger one clamped to 100.
const TaxCountryTaxesPageSize = 20

var taxCountryTaxesPaginationKeys = [10]string{
	"count", "current_page", "has_next", "has_previous", "limit",
	"next_page", "offset", "previous_page", "total", "total_pages",
}

// TaxCountryTaxesPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func TaxCountryTaxesPaginationKeys() []string {
	return append(make([]string, 0, len(taxCountryTaxesPaginationKeys)), taxCountryTaxesPaginationKeys[:]...)
}

// UnroutedMethodBody is the 200 body the legacy underscore row path returns for an unrouted method; the hyphenated paths return 404.
func UnroutedMethodBody() map[string]any {
	return map[string]any{"isSuccess": false, "message": "Invalid API Request"}
}

// TaxCountryTaxesUnknownPathMessage opens the tax router's 404 text for a method the listing path does not route.
const TaxCountryTaxesUnknownPathMessage = "Unknown tax path"

// TaxCountryTaxesTransactionNoHeader files a new country tax under this transaction number.
const TaxCountryTaxesTransactionNoHeader = "X-Transaction-No"
