package categories

// The path templates; the client fills {api_version}. The literal sub-paths come before the row path.
const (
	CategoriesPath = "/api/{api_version}/admin/categories"

	CategoriesCountPath = "/api/{api_version}/admin/categories/count"

	CategoriesCheckNamePath = "/api/{api_version}/admin/categories/check-name"

	CategoryPath = "/api/{api_version}/admin/categories/{category_id}"

	CategorySectionPath = "/api/{api_version}/admin/categories/{category_id}/{section}"
)

// The body keys: the listing's rows and the count's record, one row and every write, and the count itself.
const (
	CategoriesKey      = "categories"
	CategoryKey        = "category"
	CategoriesCountKey = "count"
)

// categorySections are the six values {section} takes, each also a key of the row and of group_meta.
var categorySections = [6]string{
	"general", "description", "images_media", "availability_visibility", "tax_shipping", "assign_products",
}

// CategorySections returns a fresh copy of the six values {section} takes.
func CategorySections() []string {
	return append(make([]string, 0, len(categorySections)), categorySections[:]...)
}

// CategoriesPageSize is the page size when limit is not sent, and the ceiling on one that is.
const CategoriesPageSize = 20

var categoriesPaginationKeys = [10]string{
	"count", "current_page", "has_next", "has_previous", "limit",
	"next_page", "offset", "previous_page", "total", "total_pages",
}

// CategoriesPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func CategoriesPaginationKeys() []string {
	return append(make([]string, 0, len(categoriesPaginationKeys)), categoriesPaginationKeys[:]...)
}

// CategoriesTransactionNoHeader files a new category under this transaction number; it beats category.transaction_no.
const CategoriesTransactionNoHeader = "X-Transaction-No"
