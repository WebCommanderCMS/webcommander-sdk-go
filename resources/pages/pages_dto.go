package pages

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// PageReference is an {id, name} pair: the layout, an access-list customer or group, the navigation and the author; a write names one by ID, Name is read-only.
type PageReference struct {
	ID   *int64
	Name *string
}

// PagePasswordProtection is the password gate; a listing row masks the password as ***, the row read and the writes answer it as saved.
type PagePasswordProtection struct {
	Enabled  *bool
	Password *string
}

// PageMenuPlacement is whether the page is linked from a navigation; Navigation is nil when it is not, and required when a write sends Enabled true.
type PageMenuPlacement struct {
	Enabled    *bool
	Navigation *PageReference
}

// PagePublishing is the draft and approval state; the draft and published ids read nil, never 0.
type PagePublishing struct {
	Draft                *bool
	DraftPageID          *int64
	PublishedPageID      *int64
	PendingApproval      *bool
	PendingApprovalCount *int64
	PendingApprovalIds   []*int64
	Disposable           *bool
}

// Page is one page: a listing row, the row read and every write's answer, the same 23 keys; Visibility is open, hidden or restricted, VisibleTo all or selected on a restricted page and nil otherwise.
type Page struct {
	ID                 *int64
	Name               *string
	Title              *string
	Subtitle           *string
	URL                *string
	FullURL            *string
	Layout             *PageReference
	Visibility         *string
	VisibleTo          *string
	Customers          []*PageReference
	CustomerGroups     []*PageReference
	BodyClass          *string
	PasswordProtection *PagePasswordProtection
	MenuPlacement      *PageMenuPlacement
	Active             *bool
	DisableTracking    *bool
	Landing            *bool
	InTrash            *bool
	Publishing         *PagePublishing
	// Store is null on every page the store holds, so no store-bound page could be read to type it.
	Store     any
	Author    *PageReference
	CreatedAt *string
	UpdatedAt *string
}

// EnvelopeKey is the key a row read and every write's answer wrap the record in.
func (Page) EnvelopeKey() string { return PageKey }

// PageInput is the writable form for the create, the PUT replace and the PATCH, wrapped in page; an unset field is not sent, so to send null for Subtitle, BodyClass or Layout pass a map.
type PageInput struct {
	Name     string
	Title    string
	Subtitle string
	URL      string
	Layout   *PageReference
	// Visibility is open, hidden or restricted; VisibleTo is all or selected, and only with restricted.
	Visibility string
	VisibleTo  string
	// Customers take public customer_id values and CustomerGroups group ids, each bare or as a PageReference.
	Customers          []any
	CustomerGroups     []any
	BodyClass          string
	PasswordProtection *PagePasswordProtection
	MenuPlacement      *PageMenuPlacement
	Active             *bool
	DisableTracking    *bool
}

// EnvelopeKey is the wrapper the create, the PUT and the PATCH require; a bare record is 400.
func (PageInput) EnvelopeKey() string { return PageKey }

// ToMap renders the request body.
func (i PageInput) ToMap() map[string]any { return dto.ToMap(i) }

// PageURLInput is the URL write's bare body, {"url": …}; a page wrapper is 400.
type PageURLInput struct {
	URL string
}

// ToMap renders the request body.
func (i PageURLInput) ToMap() map[string]any { return dto.ToMap(i) }

// Pagination is the listing's five-key pagination block; Limit is the page size applied, the total under limit=all.
type Pagination struct {
	Total       *int64
	Limit       *int64
	Offset      *int64
	HasPrevious *bool
	HasNext     *bool
}

// PageListing is one page of the listing: the rows outside the trash, and where they sit.
type PageListing struct {
	Pages []*Page
	// Pagination is nil when the body carried none.
	Pagination *Pagination
}

// PageCount is the count sub-path's body, a bare {"total": N}.
type PageCount struct {
	Total *int64
}

// PageDeleteBlocker is one reason the page cannot be deleted, such as landing_page.
type PageDeleteBlocker struct {
	Code    *string
	Message *string
}

// PageDeleteBlockers is the delete check's body: whether a delete would pass, what blocks it and what links to the page.
type PageDeleteBlockers struct {
	CanDelete *bool
	Blockers  []*PageDeleteBlocker
	// Linked was empty on every page read, so its rows stay untyped.
	Linked []any
}

// PageShape is the typed view of a body shaped {"page": {…}}.
func PageShape(decoded any) any { return dto.Read[Page](decoded) }

// PageListingShape is the typed view of the listing body.
func PageListingShape(decoded any) any { return dto.Read[PageListing](decoded) }

// PageCountShape is the typed view of the count sub-path's body.
func PageCountShape(decoded any) any { return dto.Read[PageCount](decoded) }

// PageDeleteBlockersShape is the typed view of the delete check's body.
func PageDeleteBlockersShape(decoded any) any { return dto.Read[PageDeleteBlockers](decoded) }
