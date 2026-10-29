package products

// The path templates; the client fills {api_version}.
const (
	ProductsPath = "/api/{api_version}/admin/products"

	ProductPath = "/api/{api_version}/admin/products/{product_id}"

	ProductReviewsPath = "/api/{api_version}/admin/products/{product_id}/reviews"

	ProductStockPath = "/api/{api_version}/admin/products/{product_id}/stock"

	ProductStockLocationPath = "/api/{api_version}/admin/products/{product_id}/stock/locations/{inventory_location_id}"
)

// ProductKey is the key the row read and every row write answer under, and the review import's wrapper.
const ProductKey = "product"

// ProductsTransactionNoHeader tags a create or a review import with a transaction number; it wins over a body transaction_no it disagrees with.
const ProductsTransactionNoHeader = "X-Transaction-No"

// ProductStockLocationsKey is the key the stock rows arrive under, and the one key the stock write requires.
const ProductStockLocationsKey = "locations"

var productQueryParams = [1]string{"field_metadata"}

// ProductQueryParams returns a fresh copy of the one name the row read accepts; any other is 400 unknown_parameter.
func ProductQueryParams() []string {
	return append(make([]string, 0, len(productQueryParams)), productQueryParams[:]...)
}
