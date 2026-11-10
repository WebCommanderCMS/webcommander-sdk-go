package dashboardsections

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// SectionPeriod is the period a section covers and the one before it, as dates; TimeZone is the store's and Currency the code money figures are in.
type SectionPeriod struct {
	Duration          *string
	StartDate         *string
	EndDate           *string
	PreviousStartDate *string
	PreviousEndDate   *string
	TimeZone          *string
	Currency          *string
}

// SectionEntity is a product, page or post by id and name.
type SectionEntity struct {
	ID   *int64
	Name *string
}

// SectionOrderUpdates is what needs attention now; PendingFulfillment is the orders waiting to ship.
type SectionOrderUpdates struct {
	LowStockItems      *int64
	PendingFulfillment *int64
	PaymentsToCapture  *int64
}

// SectionBestSeller is a product and the units sold, highest first.
type SectionBestSeller struct {
	Product   *SectionEntity
	UnitsSold *int64
}

// SectionContentUpdate is one content item to review; Count is an empty string on the rows measured.
type SectionContentUpdate struct {
	Key   *string
	Label *string
	Count *string
}

// SectionTopPage is a page and its views.
type SectionTopPage struct {
	Page  *SectionEntity
	Views *int64
}

// SectionNeedsAttention is the bookings that need attention: how many, and the latest of them. Latest is empty on the store measured, so the row's shape is not known and it is kept as it arrived.
type SectionNeedsAttention struct {
	Count  *int64
	Latest []any
}

// SectionDraft is a blog draft waiting to be published.
type SectionDraft struct {
	Post      *SectionEntity
	UpdatedAt *string
}

// SectionDrafts is the blog drafts to publish: how many, and the latest.
type SectionDrafts struct {
	Count  *int64
	Latest []*SectionDraft
}

// SectionViewedPost is a blog post and its views.
type SectionViewedPost struct {
	Post  *SectionEntity
	Views *int64
}

// DashboardSectionReport is one tab of the Dashboard. Section echoes the path value, and only the fields of that tab are present. ecommerce: OrderUpdates, BestSellingProducts. content: ContentUpdates, TopPages. bookings: NeedsAttention, Upcoming. courses: NewEnrolments, UnpaidOrders, PopularCourses. blog: CommentsToModerate, DraftsToPublish, MostViewedPosts. Upcoming and PopularCourses are empty on the store measured, so their rows are kept as they arrived.
type DashboardSectionReport struct {
	Section             *string
	Period              *SectionPeriod
	OrderUpdates        *SectionOrderUpdates
	BestSellingProducts []*SectionBestSeller
	ContentUpdates      []*SectionContentUpdate
	TopPages            []*SectionTopPage
	NeedsAttention      *SectionNeedsAttention
	Upcoming            []any
	NewEnrolments       *int64
	UnpaidOrders        *int64
	PopularCourses      []any
	CommentsToModerate  *int64
	DraftsToPublish     *SectionDrafts
	MostViewedPosts     []*SectionViewedPost
}

// DashboardSectionReportShape is the typed view of any of the five tabs.
func DashboardSectionReportShape(decoded any) any { return dto.Read[DashboardSectionReport](decoded) }
