package dashboards

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// Dashboards is the dashboard layout endpoints, bound to one client.
type Dashboards struct {
	requester http.Requester
}

// ListDashboards lists the admin's layouts — GET /admin/dashboards, {"dashboards": [...]} with no paging and any query name ignored. Each row is id, name, default and the two timestamps, and exactly one row is the default.
func (d *Dashboards) ListDashboards(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return d.requester.SendGetRequest(ctx, DashboardsPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: DashboardListingShape})
}

// CreateDashboard creates a layout — POST /admin/dashboards, answered 201 with the whole layout, bare and with no Location. The body is bare and needs name, unique among the admin's layouts (409 name_taken); an empty body is 400. It starts as a copy of the layout named in copy_from or of the default one, then the sections and kpis sent are applied. Unknown fields are ignored.
func (d *Dashboards) CreateDashboard(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return d.requester.SendPostRequest(ctx, DashboardsPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: DashboardShape})
}

// GetDashboard reads one layout — GET /admin/dashboards/{dashboard_id}, where the slot takes the layout's id, a UUID; the layout is bare. An unknown id is 404 not_found on dashboard_id.
func (d *Dashboards) GetDashboard(ctx context.Context, dashboardID string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return d.requester.SendGetRequest(ctx, DashboardPath, &http.RequestOptions{
		PathParams:   map[string]any{"dashboard_id": dashboardID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: DashboardShape})
}

// PatchDashboard changes only the fields sent — PATCH /admin/dashboards/{dashboard_id}, answered 200 with the whole layout. An empty body is 400; unknown fields are ignored. A blank name or one over 60 characters is 400, a section key the layout does not have is 400, a standard section with a position is 400 and hiding contact is 400. default true moves the store's default layout to this one.
func (d *Dashboards) PatchDashboard(ctx context.Context, dashboardID string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return d.requester.SendPatchRequest(ctx, DashboardPath, &http.RequestOptions{
		PathParams: map[string]any{"dashboard_id": dashboardID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: DashboardShape})
}

// DeleteDashboard deletes one layout — DELETE /admin/dashboards/{dashboard_id}, 204 with no body; the layout then reads 404.
func (d *Dashboards) DeleteDashboard(ctx context.Context, dashboardID string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return d.requester.SendDeleteRequest(ctx, DashboardPath, &http.RequestOptions{
		PathParams:   map[string]any{"dashboard_id": dashboardID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "dashboards"

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

// New binds the dashboard layout endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *Dashboards {
	return &Dashboards{requester: requester}
}
