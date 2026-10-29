package taxrules

// The path templates; the client fills {api_version}.
const (
	TaxRulesPath = "/api/{api_version}/admin/tax_rules"

	TaxRulePath = "/api/{api_version}/admin/tax_rules/{tax_rule_id}"
)

// TaxRulesKey is the key the listing's rows arrive under — the resource's own plural name.
const TaxRulesKey = "tax_rules"

// TaxRuleKey is the key one row arrives under, and the wrapper both writes require.
const TaxRuleKey = "tax_rule"

var taxRulesQueryParams = [3]string{"limit", "offset", "page"}

// TaxRulesQueryParams returns a fresh copy of the three names the listing acts on.
func TaxRulesQueryParams() []string {
	return append(make([]string, 0, len(taxRulesQueryParams)), taxRulesQueryParams[:]...)
}

// TaxRulesPageSize is the page size when limit is not sent; a limit above 100 is held to 100.
const TaxRulesPageSize = 20

var taxRulesPaginationKeys = [10]string{
	"count", "current_page", "has_next", "has_previous", "limit",
	"next_page", "offset", "previous_page", "total", "total_pages",
}

// TaxRulesPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func TaxRulesPaginationKeys() []string {
	return append(make([]string, 0, len(taxRulesPaginationKeys)), taxRulesPaginationKeys[:]...)
}

// UnroutedMethodBody is the body the store answers at 200, in place of a 405, for an unrouted method.
func UnroutedMethodBody() map[string]any {
	return map[string]any{"isSuccess": false, "message": "Invalid API Request"}
}
