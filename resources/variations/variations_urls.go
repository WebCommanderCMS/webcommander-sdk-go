package variations

// The path templates; the client fills {api_version}.
const (
	VariationsPath = "/api/{api_version}/admin/variations"

	VariationPath = "/api/{api_version}/admin/variations/{variation_id}"

	VariationOptionPath = "/api/{api_version}/admin/variations/{variation_id}/options/{variation_option_id}"

	ProductVariationsPath = "/api/{api_version}/admin/products/{product_id}/variations"

	ProductVariationsBulkPath = "/api/{api_version}/admin/products/{product_id}/variations/bulk"
)

// VariationsTransactionNoHeader tags a variant import with a transaction number; it wins over a body transaction_no it disagrees with.
const VariationsTransactionNoHeader = "X-Transaction-No"

// VariationsKey is the key the listing's rows arrive under, and the product read's groups too.
const VariationsKey = "variations"

// VariationKey is the key one type arrives under, and the wrapper every type write requires.
const VariationKey = "variation"

// ProductVariantsKey is the key a product variant group's variants arrive under.
const ProductVariantsKey = "variants"

// ProductVariationsUpdatedKey and ProductVariationsFailedKey list the variants the bulk write wrote and refused.
const (
	ProductVariationsUpdatedKey = "updated"
	ProductVariationsFailedKey  = "failed"
)

// VariationsPageSize is the page size applied when no limit is sent; VariationsPageSizeCap the ceiling on one that is.
const (
	VariationsPageSize    = 20
	VariationsPageSizeCap = 100
)

var variationsQueryParams = [6]string{"search", "standard", "limit", "offset", "page", "field_metadata"}

// VariationsQueryParams returns a fresh copy of the names the listing accepts; any other, q, sort and dir included, is 400.
func VariationsQueryParams() []string {
	return append(make([]string, 0, len(variationsQueryParams)), variationsQueryParams[:]...)
}

var variationStandards = [3]string{"text", "color", "image"}

// VariationStandards returns a fresh copy of the three standards a type takes.
func VariationStandards() []string {
	return append(make([]string, 0, len(variationStandards)), variationStandards[:]...)
}

var productVariationsBulkFields = [6]string{
	"available", "base_price", "cost_price", "available_stock", "min_order_quantity", "max_order_quantity",
}

// ProductVariationsBulkFields returns a fresh copy of the six variant fields the bulk write sets; at least one is required.
func ProductVariationsBulkFields() []string {
	return append(make([]string, 0, len(productVariationsBulkFields)), productVariationsBulkFields[:]...)
}

var productVariationsResponseModes = [2]string{"slim", "full"}

// ProductVariationsResponseModes returns a fresh copy of the bulk write's response values; slim when none is sent.
func ProductVariationsResponseModes() []string {
	return append(make([]string, 0, len(productVariationsResponseModes)), productVariationsResponseModes[:]...)
}
