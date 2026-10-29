package pages

// The path templates; the client fills {api_version}.
const (
	PagesPath = "/api/{api_version}/admin/pages"

	PagesCountPath = "/api/{api_version}/admin/pages/count"

	PagePath = "/api/{api_version}/admin/pages/{page_id}"

	PageURLPath = "/api/{api_version}/admin/pages/{page_id}/url"

	PageLandingPath = "/api/{api_version}/admin/pages/{page_id}/landing"

	PageDeleteBlockersPath = "/api/{api_version}/admin/pages/{page_id}/delete-blockers"

	PageDuplicatePath = "/api/{api_version}/admin/pages/{page_id}/duplicate"
)

// PagesKey is the key the listing's rows arrive under.
const PagesKey = "pages"

// PageKey is the key one row arrives under, and the wrapper every write but the URL write takes.
const PageKey = "page"

// PagesCountKey is the one key the count sub-path answers with, bare at the top of the body.
const PagesCountKey = "total"

// PagesPageSize is the page size applied when no limit is sent; limit takes 1 to PagesPageSizeCap or "all", and 0 or above the cap is 400 out_of_range.
const PagesPageSize = 20

const PagesPageSizeCap = 100

// PagesIfMatchHeader carries the ETag a read answered; a stale one makes the write 412.
const PagesIfMatchHeader = "If-Match"
