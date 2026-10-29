package sitesettings

// The path templates; the client fills {api_version}.
const (
	SiteSettingsPath = "/api/{api_version}/admin/settings/site_settings"

	SiteSettingsSectionPath = "/api/{api_version}/admin/settings/site_settings/{section}"
)

var siteSettingsSections = [8]string{
	"administration",
	"brand_kit",
	"date_time",
	"domain_settings",
	"email_settings",
	"favicon",
	"store_details",
	"store_logo",
}

// siteSettingsMediaSections are the two sections of kind media.
var siteSettingsMediaSections = [2]string{
	"favicon",
	"store_logo",
}

var siteSettingsWritableSections = [6]string{
	"administration",
	"brand_kit",
	"date_time",
	"domain_settings",
	"email_settings",
	"store_details",
}

// SiteSettingsSections returns a fresh copy of the eight values {section} takes.
func SiteSettingsSections() []string {
	return append(make([]string, 0, len(siteSettingsSections)), siteSettingsSections[:]...)
}

// SiteSettingsMediaSections returns a fresh copy of the two this SDK reads and never writes.
func SiteSettingsMediaSections() []string {
	return append(make([]string, 0, len(siteSettingsMediaSections)), siteSettingsMediaSections[:]...)
}

// SiteSettingsWritableSections returns a fresh copy of the six a PATCH may name here.
func SiteSettingsWritableSections() []string {
	return append(make([]string, 0, len(siteSettingsWritableSections)), siteSettingsWritableSections[:]...)
}
