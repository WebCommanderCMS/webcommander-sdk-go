package pagesettings

// The path templates for the category page settings endpoints; the client fills {api_version}.
const (
	CategoryPagePath = "/api/{api_version}/admin/settings/page_settings/category_page"

	CategoryPageSectionPath = "/api/{api_version}/admin/settings/page_settings/category_page/{section}"
)

// categoryPageSections are the two values {section} takes; both are writable.
var categoryPageSections = [2]string{"category_information", "product_information"}

// CategoryPageSections returns a fresh copy of the two values {section} takes.
func CategoryPageSections() []string {
	return append(make([]string, 0, len(categoryPageSections)), categoryPageSections[:]...)
}
