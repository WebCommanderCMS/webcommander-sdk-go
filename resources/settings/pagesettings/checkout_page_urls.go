package pagesettings

// The path templates for the checkout page settings endpoints; the client fills {api_version}.
const (
	CheckoutPagePath = "/api/{api_version}/admin/settings/page_settings/checkout_page"

	CheckoutPageSectionPath = "/api/{api_version}/admin/settings/page_settings/checkout_page/{section}"
)

// checkoutPageSections are the six values {section} takes, in facade order.
var checkoutPageSections = [6]string{
	"billing_address",
	"shipping_address",
	"customer_checkout",
	"delivery_options",
	"order_summary",
	"others",
}

// checkoutPageWritableSections are the four sections a PATCH may name; the two address sections are read-only.
var checkoutPageWritableSections = [4]string{
	"customer_checkout",
	"delivery_options",
	"order_summary",
	"others",
}

// CheckoutPageSections returns a fresh copy of the six values {section} takes.
func CheckoutPageSections() []string {
	return append(make([]string, 0, len(checkoutPageSections)), checkoutPageSections[:]...)
}

// CheckoutPageWritableSections returns a fresh copy of the four of those six a PATCH may name.
func CheckoutPageWritableSections() []string {
	return append(make([]string, 0, len(checkoutPageWritableSections)), checkoutPageWritableSections[:]...)
}
