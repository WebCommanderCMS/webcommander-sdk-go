package dashboardkpis

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// DashboardKPIs is the dashboard KPI endpoints, bound to one client.
type DashboardKPIs struct {
	requester http.Requester
}

// GetDashboardKPIOptions lists the metrics a tile can show — GET /admin/dashboard/kpi_options, grouped, with min 2 and max 5 tiles and the keys now in use under selected. It reads the same tiles as kpi_selection.
func (k *DashboardKPIs) GetDashboardKPIOptions(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return k.requester.SendGetRequest(ctx, DashboardKPIOptionsPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: KPIOptionsShape})
}

// GetDashboardKPISelection reads the tiles in use with their figures — GET /admin/dashboard/kpi_selection, the tiles of the default layout under selected in position order. duration (last_30_days when left out) or start_date with end_date (YYYY-MM-DD, sent together) choose the period; a bad value is 400.
func (k *DashboardKPIs) GetDashboardKPISelection(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return k.requester.SendGetRequest(ctx, DashboardKPISelectionPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: KPISelectionShape})
}

// UpdateDashboardKPISelection saves the tiles — PUT /admin/dashboard/kpi_selection with {"metrics": [...]}, answered 200 like the read. The list replaces the default layout's tiles, which is store-wide state, and is read back at once by that layout and by kpi_options. Fewer than 2 or more than 5, a repeat, an unknown key, a missing metrics and a non-list are 400. duration or the dates in queryParams set the period of the figures in the answer.
func (k *DashboardKPIs) UpdateDashboardKPISelection(ctx context.Context, payload any, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return k.requester.SendPutRequest(ctx, DashboardKPISelectionPath, &http.RequestOptions{
		JSONBody: payload, QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: KPISelectionShape})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "dashboard_kpis"

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

func requirePayload(payload any) error {
	if payload == nil || isNilValue(payload) {
		return wcerrors.NewUsageError("payload is required: this endpoint takes a body, and nil would " +
			"send the request with no body at all")
	}
	return nil
}

func isNilValue(value any) bool {
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Map, reflect.Slice, reflect.Pointer, reflect.Interface, reflect.Func:
		return reflected.IsNil()
	}
	return false
}

// New binds the dashboard KPI endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *DashboardKPIs {
	return &DashboardKPIs{requester: requester}
}
