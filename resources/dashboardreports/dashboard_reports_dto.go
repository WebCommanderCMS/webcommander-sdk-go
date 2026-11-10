package dashboardreports

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// ReportPeriod is the period a report covers and the one before it, as dates; TimeZone is the store's and Currency the code money figures are in. Duration is custom for a date range.
type ReportPeriod struct {
	Duration          *string
	StartDate         *string
	EndDate           *string
	PreviousStartDate *string
	PreviousEndDate   *string
	TimeZone          *string
	Currency          *string
}

// ReportKPI is one figure with the previous period's. Format is money, count or percent and Direction is up, down or flat; ChangePercent is nil when the previous value was zero. In trends there is no Group.
type ReportKPI struct {
	Key           *string
	Code          *string
	Label         *string
	Group         *string
	Format        *string
	Value         *float64
	PreviousValue *float64
	ChangePercent *float64
	Direction     *string
}

// ReportOrders is the order counts of the period.
type ReportOrders struct {
	Count     *int64
	Pending   *int64
	Cancelled *int64
}

// ReportSales is the sales figures of the period, rounded to 2 decimals.
type ReportSales struct {
	Subtotal          *float64
	Total             *float64
	ItemsSold         *int64
	AverageOrderValue *float64
}

// ReportCustomers is the customers registered in the period.
type ReportCustomers struct {
	New *int64
}

// ReportProducts is the products added in the period, and those out of stock now.
type ReportProducts struct {
	New        *int64
	OutOfStock *int64
}

// ReportPageViews is the page views of the period.
type ReportPageViews struct {
	Total  *int64
	Unique *int64
}

// ReportEntity is a page, product or post by id and name.
type ReportEntity struct {
	ID   *int64
	Name *string
}

// ReportViewItem is one row of the most viewed list: a Page or, with entity_name product, a Product.
type ReportViewItem struct {
	Page    *ReportEntity
	Product *ReportEntity
	Views   *int64
}

// ReportViews is the most viewed pages, or products when the query's entity_name is product; Entity says which.
type ReportViews struct {
	Entity *string
	Items  []*ReportViewItem
}

// QuickReport is the quick report. KPIs holds the 27 metrics, or only the ones the query's kpis names, in that order. Inside a summary only the grouped totals and ViewReport are present.
type QuickReport struct {
	Period     *ReportPeriod
	KPIs       []*ReportKPI `wire:"kpis"`
	Orders     *ReportOrders
	Sales      *ReportSales
	Customers  *ReportCustomers
	Products   *ReportProducts
	PageViews  *ReportPageViews
	ViewReport *ReportViews
	Trends     []*ReportKPI
}

// SalesRow is one row of the sales report: a day for a period up to a month, else a month. Without breakdown only the date fields, Label, MonthName, Sales and TotalSales are present; Day is present on day rows only.
type SalesRow struct {
	Date              *string
	Year              *int64
	Month             *int64
	Day               *int64
	Label             *string
	MonthName         *string
	Sales             *float64
	TotalSales        *float64
	Orders            *int64
	OrdersValue       *float64
	GrossSales        *float64
	Shipping          *float64
	Handling          *float64
	Tax               *float64
	Discount          *float64
	Refunds           *float64
	AverageOrderValue *float64
	TopProduct        *ReportEntity
}

// SalesTotals is the rows added up; Sales always, the rest with breakdown.
type SalesTotals struct {
	Sales             *float64
	Orders            *int64
	OrdersValue       *float64
	GrossSales        *float64
	Shipping          *float64
	Handling          *float64
	Tax               *float64
	Discount          *float64
	Refunds           *float64
	AverageOrderValue *float64
}

// SalesReport is the sales report. GroupBy is day for a period up to a month and month for a longer one; Sort (date or sales) and Order (asc or desc) echo what was applied. A period with no sales has no rows.
type SalesReport struct {
	Period  *ReportPeriod
	GroupBy *string
	Sort    *string
	Order   *string
	Rows    []*SalesRow
	Totals  *SalesTotals
}

// OrderCount is the orders placed in one month.
type OrderCount struct {
	Date  *string
	Year  *int64
	Month *int64
	Count *int64
}

// MostVisitedProduct is a product and its views.
type MostVisitedProduct struct {
	Product *ReportEntity
	Views   *int64
}

// BestSellingProduct is a product and the units sold, highest first.
type BestSellingProduct struct {
	Product   *ReportEntity
	UnitsSold *int64
}

// OrderUpdates is what needs attention now; PendingFulfillment is the orders waiting to ship.
type OrderUpdates struct {
	LowStockItems      *int64
	PendingFulfillment *int64
	PaymentsToCapture  *int64
}

// SummaryReport is every section of the Dashboard in one answer. QuickReport has no KPIs.
type SummaryReport struct {
	Period              *ReportPeriod
	QuickReport         *QuickReport
	OrderByMonth        []*OrderCount
	OrdersOverTime      []*OrderCount
	SalesReport         []*SalesRow
	MostVisitedProducts []*MostVisitedProduct
	BestSellingProducts []*BestSellingProduct
	OrderUpdates        *OrderUpdates
}

// QuickReportShape is the typed view of the quick report.
func QuickReportShape(decoded any) any { return dto.Read[QuickReport](decoded) }

// SalesReportShape is the typed view of the sales report.
func SalesReportShape(decoded any) any { return dto.Read[SalesReport](decoded) }

// SummaryReportShape is the typed view of the summary.
func SummaryReportShape(decoded any) any { return dto.Read[SummaryReport](decoded) }
