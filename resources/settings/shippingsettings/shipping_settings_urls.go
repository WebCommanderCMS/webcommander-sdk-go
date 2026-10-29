package shippingsettings

// ShippingSettingsPath is the area's one path template; /shipping_settings/{name} does not exist.
const ShippingSettingsPath = "/api/{api_version}/admin/settings/shipping_settings"

// shippingSettingsFields are the seven body keys, in facade order; they are fields, not sections.
var shippingSettingsFields = [7]string{
	"default_shipping_profile",
	"enable_shipping_default_tax",
	"default_delivery_profile",
	"enable_delivery_default_tax",
	"shipping_tax_profile",
	"enable_shipping_class",
	"enable_flat_view_class",
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
