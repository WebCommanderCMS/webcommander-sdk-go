package shippingsettings

// ShippingSettingsPath is the area's one path template; the old /shipping_settings path answers 404 Unknown settings page.
const ShippingSettingsPath = "/api/{api_version}/admin/settings/shipping_settings/general"

// shippingSettingsFields are the seven body keys, in facade order; they are fields, not sections.
var shippingSettingsFields = [7]string{
	"default_shipping_profile",
	"apply_store_tax_settings_to_shipping",
	"default_delivery_profile",
	"apply_store_tax_settings_to_delivery",
	"shipping_tax_profile",
	"enable_shipping_class",
	"enable_class_in_flat_view",
}

// ShippingSettingsFields returns a fresh copy of the seven names a PATCH body may carry.
func ShippingSettingsFields() []string {
	return append(make([]string, 0, len(shippingSettingsFields)), shippingSettingsFields[:]...)
}

// shippingSettingsReferenceFields are the three fields written by id that read back as {id, name} or null.
var shippingSettingsReferenceFields = [3]string{
	"default_shipping_profile",
	"default_delivery_profile",
	"shipping_tax_profile",
}

// ShippingSettingsReferenceFields returns a fresh copy of the three fields that hold a record, not a scalar.
func ShippingSettingsReferenceFields() []string {
	return append(make([]string, 0, len(shippingSettingsReferenceFields)),
		shippingSettingsReferenceFields[:]...)
}
