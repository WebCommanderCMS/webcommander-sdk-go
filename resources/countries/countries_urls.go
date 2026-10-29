package countries

// The path templates; the client fills {api_version}.
const (
	CountriesPath = "/api/{api_version}/admin/countries"

	CountryPath = "/api/{api_version}/admin/countries/{country_id}"

	CountryStatesPath = "/api/{api_version}/admin/countries/{country_id}/states"

	CountryStatePath = "/api/{api_version}/admin/countries/{country_id}/states/{state_id}"
)

// CountriesKey is the key the listing's rows arrive under.
const CountriesKey = "countries"

// CountryKey is the key one country arrives under, and the short parent beside a state listing or read.
const CountryKey = "country"

// CountryStatesKey is the key a state listing's rows arrive under.
const CountryStatesKey = "states"

// CountryStateKey is the key one state arrives under.
const CountryStateKey = "state"

// CountriesPageSize is the page size applied when no limit is sent; a sent one runs 1 to 100, or all.
const CountriesPageSize = 20

// CountriesMaxLimit is the largest numeric limit either listing takes.
const CountriesMaxLimit = 100

// CountriesLimitAll is the limit value that answers every row from offset on in one page.
const CountriesLimitAll = "all"

var countriesSortFields = [3]string{"id", "name", "code"}

// CountriesSortFields returns a fresh copy of the names sort takes on both listings; order takes asc or desc.
func CountriesSortFields() []string {
	return append(make([]string, 0, len(countriesSortFields)), countriesSortFields[:]...)
}
