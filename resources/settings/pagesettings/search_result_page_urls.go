package pagesettings

// The path templates for the search result page settings endpoints; the client fills {api_version}.
const (
	SearchResultPagePath = "/api/{api_version}/admin/settings/page_settings/search_result_page"

	SearchResultPageSectionPath = "/api/{api_version}/admin/settings/page_settings/search_result_page/{section}"

	SearchResultPageProductPromotionsPath = "/api/{api_version}/admin/settings/page_settings/search_result_page/product_promotion"

	SearchResultPageProductPromotionPath = "/api/{api_version}/admin/settings/page_settings/search_result_page/product_promotion/{promotion_id}"

	SearchResultPageProductPromotionActivePath = "/api/{api_version}/admin/settings/page_settings/search_result_page/product_promotion/{promotion_id}/active"
)

// sections are the three values {section} takes; product_promotion is a resource with methods of its own.
var sections = [3]string{"product_information", "seo", "reinitialise"}

// SearchResultPageSections returns a fresh copy of the three values {section} takes.
func SearchResultPageSections() []string {
	return append(make([]string, 0, len(sections)), sections[:]...)
}
