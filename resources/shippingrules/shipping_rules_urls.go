package shippingrules

// The path templates; the client fills {api_version}.
const (
	ShippingRulesPath = "/api/{api_version}/admin/shipping_rules"

	ShippingRulesCountPath = "/api/{api_version}/admin/shipping_rules/count"

	ShippingRulePath = "/api/{api_version}/admin/shipping_rules/{shipping_rule_id}"
)

// ShippingRulesKey is the key the listing's rows arrive under.
const ShippingRulesKey = "shipping_rules"

// ShippingRuleKey is the key one row arrives under, and the wrapper every write takes; a bare body is 400.
const ShippingRuleKey = "shipping_rule"

// ShippingRulesCountKey is the one key the count sub-path answers with, bare at the top of the body.
const ShippingRulesCountKey = "total"

// ShippingRulesPageSize is the page size applied when no limit is sent; limit takes 1 to ShippingRulesPageSizeCap, and 0, above the cap or "all" is 400.
const ShippingRulesPageSize = 20

const ShippingRulesPageSizeCap = 100

// ShippingRulesTransactionNoHeader files a new rule under this transaction number and wins over a transaction_no in the body.
const ShippingRulesTransactionNoHeader = "X-Transaction-No"
