package dashboardreports

import (
	"context"
	"fmt"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// DashboardReports is the dashboard report endpoints, bound to one client. Every report takes the same period query: duration (today, yesterday, last_24_hours, last_7_days, last_30_days, this_month, last_month, last_3_months, last_6_months, last_12_months, last_24_months, last_36_months, last_48_months, last_60_months or custom; last_30_days when left out) or start_date with end_date (YYYY-MM-DD, sent together, not reversed and not in the future). A bad value, a repeated duration and a lone date are 400 with every bad field listed; a name the report does not know is ignored.
type DashboardReports struct {
	requester http.Requester
}

// GetDashboardReportQuick reads the quick report — GET /admin/dashboard/reports/quick: period, kpis, grouped totals, the most viewed pages under view_report and trends. kpis is a comma list of metric keys (an unknown one is 400) that narrows kpis to those, in that order; entity_name is page (the default) or product and picks what view_report ranks.
func (r *DashboardReports) GetDashboardReportQuick(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return r.requester.SendGetRequest(ctx, DashboardReportQuickPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: QuickReportShape})
}

// HeadDashboardReportQuick asks whether the quick report is reachable — HEAD /admin/dashboard/reports/quick; the status of the read, with no body. It runs the read's query, so a bad value is 400.
func (r *DashboardReports) HeadDashboardReportQuick(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return r.requester.SendApiRequest(ctx, "HEAD", DashboardReportQuickPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError, ResponseFormat: format})
}

// GetDashboardReportSales reads the sales report — GET /admin/dashboard/reports/sales: rows by day or month and totals. breakdown (true or false) adds orders, shipping, tax, discount, refunds and the top product to each row and total; sort is date (the default) or sales and order is asc (the default) or desc; any other value is 400.
func (r *DashboardReports) GetDashboardReportSales(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return r.requester.SendGetRequest(ctx, DashboardReportSalesPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: SalesReportShape})
}

// HeadDashboardReportSales asks whether the sales report is reachable — HEAD /admin/dashboard/reports/sales; the status of the read, with no body. It runs the read's query, so a bad value is 400.
func (r *DashboardReports) HeadDashboardReportSales(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return r.requester.SendApiRequest(ctx, "HEAD", DashboardReportSalesPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError, ResponseFormat: format})
}

// GetDashboardReportSummary reads the summary — GET /admin/dashboard/reports/summary: the quick report without kpis, orders by month, sales rows, the most visited and best selling products and the order updates, in one answer. compare_previous (true or false) is checked as a boolean; any other value is 400.
func (r *DashboardReports) GetDashboardReportSummary(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return r.requester.SendGetRequest(ctx, DashboardReportSummaryPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: SummaryReportShape})
}

// HeadDashboardReportSummary asks whether the summary is reachable — HEAD /admin/dashboard/reports/summary; the status of the read, with no body. It runs the read's query, so a bad value is 400.
func (r *DashboardReports) HeadDashboardReportSummary(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return r.requester.SendApiRequest(ctx, "HEAD", DashboardReportSummaryPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "dashboard_reports"

func oneResponseFormat(responseFormat []string) (string, error) {
	switch len(responseFormat) {
	case 0:
		return "", nil
	case 1:
		return responseFormat[0], nil
	default:
		return "", wcerrors.NewUsageError(fmt.Sprintf(
			"responseFormat: pass one format or none, got %d — the trailing argument is optional, "+
				"and omitting it defers to the client", len(responseFormat)))
	}
}

// New binds the dashboard report endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *DashboardReports {
	return &DashboardReports{requester: requester}
}
