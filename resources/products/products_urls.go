package products

// The path templates; the client fills {api_version}.
const (
	ProductsPath = "/api/{api_version}/admin/products"

	ProductsCountPath = "/api/{api_version}/admin/products/count"

	ProductsCheckSkuPath = "/api/{api_version}/admin/products/check-sku"

	ProductsSpecialPath = "/api/{api_version}/admin/products/special"

	ProductsBulkDeletePath = "/api/{api_version}/admin/products/bulk-delete"

	ProductPath = "/api/{api_version}/admin/products/{product_id}"

	ProductReviewsPath = "/api/{api_version}/admin/products/{product_id}/reviews"

	ProductReviewPath = "/api/{api_version}/admin/products/{product_id}/reviews/{review_id}"

	ProductInventoriesPath = "/api/{api_version}/admin/products/{product_id}/inventories"

	ProductImagePath = "/api/{api_version}/admin/products/{product_id}/image"

	ProductMediaPath = "/api/{api_version}/admin/products/{product_id}/media"

	ProductMediaImagePath = "/api/{api_version}/admin/products/{product_id}/media/images/{image_id}"

	ProductStockPath = "/api/{api_version}/admin/products/{product_id}/stock"

	ProductStockAdjustmentsPath = "/api/{api_version}/admin/products/{product_id}/stock/adjustments"

	ProductSeoPath = "/api/{api_version}/admin/products/{product_id}/seo"

	ProductMediaOrderPath = "/api/{api_version}/admin/products/{product_id}/media/order"

	ProductStockLocationPath = "/api/{api_version}/admin/products/{product_id}/stock/locations/{inventory_location_id}"
)

// ProductKey is the key the row read and every row write answer under, and the review import's wrapper.
const ProductKey = "product"

// ProductSeoKey is the key the SEO read answers under.
const ProductSeoKey = "seo"

// ProductReviewKey is the key a single review answers under.
const ProductReviewKey = "review"

// ProductsTransactionNoHeader tags a create, a review import or an inventory top-up with a transaction number; it wins over a body transaction_no it disagrees with.
const ProductsTransactionNoHeader = "X-Transaction-No"

// ProductStockLocationsKey is the key the stock rows arrive under, and the one key the stock write requires.
const ProductStockLocationsKey = "locations"

var productQueryParams = [1]string{"field_metadata"}

// ProductQueryParams returns a fresh copy of the one name the row read accepts; any other is 400 unknown_parameter.
func ProductQueryParams() []string {
	return append(make([]string, 0, len(productQueryParams)), productQueryParams[:]...)
}
