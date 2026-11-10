package dashboards

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// DashboardPart is one part of the seo section. Read back, the title and CanHide come with it; a write needs only Key and Visible. The kpis part cannot be hidden.
type DashboardPart struct {
	Key     *string
	Title   *string
	Visible *bool
	CanHide *bool
}

// DashboardTab is one tab of the section_tabs section: ecommerce, content, bookings, courses or blog. A write needs Key and sets Visible or Position, counted from 1 among the tabs.
type DashboardTab struct {
	Key         *string
	Title       *string
	Description *string
	Visible     *bool
	CanHide     *bool
	Position    *int64
}

// DashboardSection is a section of a layout, read and written in one type. The standard ones are announcement, performance, activity and contact; they only show or hide, and contact cannot hide. The customisable ones are setup_guide, section_tabs (with Tabs) and seo (with Parts); they also take a Position. Read back, the title, description and CanHide come with it; a write needs only Key and what changes.
type DashboardSection struct {
	Key         *string
	Title       *string
	Description *string
	Visible     *bool
	CanHide     *bool
	Position    *int64
	Tabs        []*DashboardTab
	Parts       []*DashboardPart
}

// DashboardKPI is one KPI tile of a layout, in tile order. Key is a kpi_options metric key.
type DashboardKPI struct {
	Position *int64
	Key      *string
	Label    *string
	Group    *string
}

// Dashboard is one layout. ID is the dashboard_id of the other calls, a UUID. Exactly one layout is the default.
type Dashboard struct {
	ID                   *string
	Name                 *string
	Default              *bool
	StandardSections     []*DashboardSection
	CustomisableSections []*DashboardSection
	KPIs                 []*DashboardKPI `wire:"kpis"`
	CreatedAt            *string
	UpdatedAt            *string
}

// DashboardSummary is a layout as the listing shows it, without its sections and tiles.
type DashboardSummary struct {
	ID        *string
	Name      *string
	Default   *bool
	CreatedAt *string
	UpdatedAt *string
}

// DashboardListing is the listing: every layout of the admin, with no paging.
type DashboardListing struct {
	Dashboards []*DashboardSummary
}

// DashboardInput is the write body, bare. A create needs Name, unique among the admin's layouts (409); it starts as a copy of the CopyFrom layout, the default one when that is left out. A patch sends only what changes; the sections named are merged into the layout, the rest is kept. KPIs is the full list of 2 to 5 metric keys in tile order. A standard section takes no Position (400). Default true makes this layout the store's default layout; false on the default layout is 409, as the default moves only by making another layout the default.
type DashboardInput struct {
	Name                 string
	Default              *bool
	CopyFrom             string
	StandardSections     []*DashboardSection
	CustomisableSections []*DashboardSection
	KPIs                 []string `wire:"kpis"`
}

// ToMap renders the request body.
func (i DashboardInput) ToMap() map[string]any { return dto.ToMap(i) }

// DashboardShape is the typed view of one layout, a bare body.
func DashboardShape(decoded any) any { return dto.Read[Dashboard](decoded) }

// DashboardListingShape is the typed view of the listing body.
func DashboardListingShape(decoded any) any { return dto.Read[DashboardListing](decoded) }
