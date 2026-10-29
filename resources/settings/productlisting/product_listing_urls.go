package productlisting

// ProductListingPath is the area's one path template; /product_listing/{section} does not exist.
const ProductListingPath = "/api/{api_version}/admin/settings/product_listing"

// productListingSections are the four sections, in facade order; they are body keys, not path values.
var productListingSections = [4]string{
	"product_information",
	"load",
	"sorting",
	"sortable_items",
}

// ProductListingSections returns a fresh copy of the four sections a PATCH body may name.
func ProductListingSections() []string {
	return append(make([]string, 0, len(productListingSections)), productListingSections[:]...)
}
