package dashboardsections

// DashboardSectionPath is the path template; the client fills {api_version}. {section} takes one of DashboardSectionTabs.
const DashboardSectionPath = "/api/{api_version}/admin/dashboard/sections/{section}"

// dashboardSectionTabs are the five values {section} takes, the tabs of the Dashboard's section_tabs section.
var dashboardSectionTabs = [5]string{"ecommerce", "content", "bookings", "courses", "blog"}

// DashboardSectionTabs returns a fresh copy of the five values {section} takes.
func DashboardSectionTabs() []string {
	return append(make([]string, 0, len(dashboardSectionTabs)), dashboardSectionTabs[:]...)
}
