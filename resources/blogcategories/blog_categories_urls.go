package blogcategories

// The path templates; the client fills {api_version}.
const (
	BlogCategoriesPath = "/api/{api_version}/admin/blog_categories"

	BlogCategoryPath = "/api/{api_version}/admin/blog_categories/{blog_category_id}"
)

// The body keys: the listing's rows, and one row, which is also the wrapper both writes require.
const (
	BlogCategoriesKey = "blog_categories"
	BlogCategoryKey   = "blog_category"
)

// BlogCategoriesPageSize is the page size when limit is not sent, and the ceiling a larger one is clamped to.
const BlogCategoriesPageSize = 20

var blogCategoriesPaginationKeys = [10]string{
	"count", "current_page", "has_next", "has_previous", "limit",
	"next_page", "offset", "previous_page", "total", "total_pages",
}

// BlogCategoriesPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func BlogCategoriesPaginationKeys() []string {
	return append(make([]string, 0, len(blogCategoriesPaginationKeys)), blogCategoriesPaginationKeys[:]...)
}
