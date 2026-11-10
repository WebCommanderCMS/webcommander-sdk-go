package dashboardkpis

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// DashboardRef is the default layout the KPI tiles belong to; ID is its dashboard_id.
type DashboardRef struct {
	ID   *string
	Name *string
}

// KPIPeriod is the period a figure covers and the one before it, as dates; TimeZone is the store's and Currency the code money figures are in.
type KPIPeriod struct {
	Duration          *string
	StartDate         *string
	EndDate           *string
	PreviousStartDate *string
	PreviousEndDate   *string
	TimeZone          *string
	Currency          *string
}

// KPIMetric is one metric a tile can show. Key is what a selection takes; Code is the report's own name for it and differs for some, such as aov and average_order_value.
type KPIMetric struct {
	Key    *string
	Code   *string
	Label  *string
	Format *string
	InUse  *bool
}

// KPIGroup is a group of metrics: ecommerce, content, blog, bookings or courses.
type KPIGroup struct {
	Key     *string
	Label   *string
	Metrics []*KPIMetric
}

// KPIOptions is every metric a tile can show. Min and Max bound the number of tiles, and Selected lists the keys in use, in tile order.
type KPIOptions struct {
	Dashboard *DashboardRef
	Min       *int64
	Max       *int64
	Selected  []*string
	Groups    []*KPIGroup
}

// KPITile is one selected tile with its figure. Format is money, count or percent and Direction is up, down or flat; ChangePercent is nil when the previous period was zero. Every count is a float on the wire.
type KPITile struct {
	Position      *int64
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

// KPISelection is the tiles in use with their figures for the period; the read and the save answer alike.
type KPISelection struct {
	Dashboard   *DashboardRef
	OnlineStore *bool
	Period      *KPIPeriod
	Selected    []*KPITile
}

// KPISelectionInput is the save body, bare: Metrics is the full list of 2 to 5 distinct metric keys, in tile order. It replaces the default layout's tiles.
type KPISelectionInput struct {
	Metrics []string
}

// ToMap renders the request body.
func (i KPISelectionInput) ToMap() map[string]any { return dto.ToMap(i) }

// KPIOptionsShape is the typed view of the metric options.
func KPIOptionsShape(decoded any) any { return dto.Read[KPIOptions](decoded) }

// KPISelectionShape is the typed view of the tiles in use, as the read and the save answer.
func KPISelectionShape(decoded any) any { return dto.Read[KPISelection](decoded) }
