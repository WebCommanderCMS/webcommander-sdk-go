package dashboardsections

import (
	"context"
	"fmt"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// DashboardSections is the dashboard section endpoints, bound to one client. The Dashboard shows five tabs; a tab is a path value.
type DashboardSections struct {
	requester http.Requester
}

// GetDashboardSection reads one tab — GET /admin/dashboard/sections/{section}, where section is one of DashboardSectionTabs (ecommerce, content, bookings, courses, blog); the answer is bare, section echoing it, with period and that tab's own fields. Any other value is 404 not_found on key. duration or start_date with end_date choose the period, as in the reports.
func (s *DashboardSections) GetDashboardSection(ctx context.Context, section string, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, DashboardSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, QueryParams: queryParams,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: DashboardSectionReportShape})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "dashboard_sections"

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

// New binds the dashboard section endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *DashboardSections {
	return &DashboardSections{requester: requester}
}
