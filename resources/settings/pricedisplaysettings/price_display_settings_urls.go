package pricedisplaysettings

// The path templates for the price display settings endpoints; the client fills {api_version}.
const (
	PriceDisplaySettingsPath = "/api/{api_version}/admin/settings/price_display_settings"

	PriceDisplaySettingsSectionPath = "/api/{api_version}/admin/settings/price_display_settings/{section}"
)

// priceDisplaySettingsSections are the two values {section} takes, in facade order; both are writable.
var priceDisplaySettingsSections = [2]string{
	"order_item",
	"order_summary",
}

// PriceDisplaySettingsSections returns a fresh copy of the two values {section} takes.
func PriceDisplaySettingsSections() []string {
	return append(make([]string, 0, len(priceDisplaySettingsSections)), priceDisplaySettingsSections[:]...)
}
