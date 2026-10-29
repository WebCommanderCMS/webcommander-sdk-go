package pagesettings

// The path templates for the product page settings endpoints; the client fills {api_version}.
const (
	ProductPagePath = "/api/{api_version}/admin/settings/page_settings/product_page"

	ProductPageSectionPath = "/api/{api_version}/admin/settings/page_settings/product_page/{section}"
)

// productPageSections are the seven values {section} takes, in facade order; all are writable.
var productPageSections = [7]string{
	"configuration",
	"variations",
	"image_and_video",
	"up_sell",
	"related_products",
	"review_and_ratings",
	"quick_view",
}

// ProductPageSections returns a fresh copy of the seven values {section} takes.
func ProductPageSections() []string {
	return append(make([]string, 0, len(productPageSections)), productPageSections[:]...)
}
