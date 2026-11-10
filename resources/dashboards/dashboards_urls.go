package dashboards

// The path templates; the client fills {api_version}. {dashboard_id} takes the layout's id, a UUID.
const (
	DashboardsPath = "/api/{api_version}/admin/dashboards"

	DashboardPath = "/api/{api_version}/admin/dashboards/{dashboard_id}"
)
