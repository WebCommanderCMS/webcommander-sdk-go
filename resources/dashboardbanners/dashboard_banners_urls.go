package dashboardbanners

// The path templates; the client fills {api_version}. {banner_key} takes the banner's key, snake_case.
const (
	DashboardBannersPath = "/api/{api_version}/admin/dashboard/banners"

	DashboardBannerPath = "/api/{api_version}/admin/dashboard/banners/{banner_key}"

	DashboardContentInfoPath = "/api/{api_version}/admin/dashboard/content_info"
)
