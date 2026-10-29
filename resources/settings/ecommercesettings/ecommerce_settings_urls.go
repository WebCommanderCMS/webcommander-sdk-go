package ecommercesettings

// The path templates; the client fills {api_version}.
const (
	EcommerceSettingsPath = "/api/{api_version}/admin/settings/ecommerce_settings"

	EcommerceSettingsSectionPath = "/api/{api_version}/admin/settings/ecommerce_settings/{section}"
)

var ecommerceSettingsAreaKeys = [7]string{
	"continue_shopping",
	"product_inventory",
	"pricing",
	"product_restriction",
	"order_filters",
	"location_settings",
	"order_flags",
}

var ecommerceSettingsSections = [6]string{
	"continue_shopping",
	"product_inventory",
	"pricing",
	"product_restriction",
	"order_filters",
	"location_settings",
}

// EcommerceSettingsCollectionKey is the one area key that is no section: order_flags, a collection of its own.
const EcommerceSettingsCollectionKey = "order_flags"

// EcommerceSettingsAreaKeys returns a fresh copy of the seven keys the facade carries, in store order.
func EcommerceSettingsAreaKeys() []string {
	return append(make([]string, 0, len(ecommerceSettingsAreaKeys)), ecommerceSettingsAreaKeys[:]...)
}

// EcommerceSettingsSections returns a fresh copy of the six values {section} takes, all writable.
func EcommerceSettingsSections() []string {
	return append(make([]string, 0, len(ecommerceSettingsSections)), ecommerceSettingsSections[:]...)
}
