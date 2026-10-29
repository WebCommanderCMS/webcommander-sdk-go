package currencies

// The path templates; the client fills {api_version}.
const (
	CurrenciesPath = "/api/{api_version}/admin/currencies"

	CurrenciesCountPath = "/api/{api_version}/admin/currencies/count"

	CurrencyPath = "/api/{api_version}/admin/currencies/{currency_id}"

	CurrencyBasePath = "/api/{api_version}/admin/currencies/{currency_id}/base"
)

// CurrenciesKey is the key the listing's rows arrive under — the resource's own plural name.
const CurrenciesKey = "currencies"

// CurrencyKey is the key one row arrives under; no write takes it as a wrapper.
const CurrencyKey = "currency"

// CurrencyPreviousBaseKey is the key the base write names the row that held the flag under; null when none moved.
const CurrencyPreviousBaseKey = "previous_base"

// CurrenciesCountKey is the one key the count sub-path answers with, bare at the top of the body.
const CurrenciesCountKey = "total"

// CurrenciesPageSize is the page size applied when no limit is sent; limit takes 1 to CurrenciesPageSizeCap.
const CurrenciesPageSize = 20

const CurrenciesPageSizeCap = 100

var currenciesSortFields = [10]string{
	"id", "name", "code", "symbol", "precision", "rounding", "conversion_rate", "active", "base", "updated_at",
}

// CurrenciesSortFields returns a fresh copy of the ten names sort accepts; any other is 400 invalid_enum.
func CurrenciesSortFields() []string {
	return append(make([]string, 0, len(currenciesSortFields)), currenciesSortFields[:]...)
}

var currencyRoundingValues = [3]string{"nearest", "up", "down"}

// CurrencyRoundingValues returns a fresh copy of the values rounding takes, as the store's field_metadata lists them.
func CurrencyRoundingValues() []string {
	return append(make([]string, 0, len(currencyRoundingValues)), currencyRoundingValues[:]...)
}

// CurrenciesTransactionNoHeader files a new currency under this transaction number, the only way: a body field is 400.
const CurrenciesTransactionNoHeader = "X-Transaction-No"
