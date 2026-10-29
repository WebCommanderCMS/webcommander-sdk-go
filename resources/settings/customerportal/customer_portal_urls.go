package customerportal

// The path templates for the customer portal endpoints; the client fills {api_version}.
const (
	CustomerPortalPath = "/api/{api_version}/admin/settings/customer_portal"

	CustomerPortalSectionPath = "/api/{api_version}/admin/settings/customer_portal/{section}"
)

// customerPortalSections are the six values {section} takes, in facade order.
var customerPortalSections = [6]string{
	"portal",
	"registration",
	"login",
	"sso",
	"profile_page",
	"registration_fields",
}

// customerPortalWritableSections are the four sections group_meta marks writable; the other two are read-only.
var customerPortalWritableSections = [4]string{
	"portal",
	"registration",
	"login",
	"sso",
}

// CustomerPortalSections returns a fresh copy of the six values {section} takes.
func CustomerPortalSections() []string {
	return append(make([]string, 0, len(customerPortalSections)), customerPortalSections[:]...)
}

// CustomerPortalWritableSections returns a fresh copy of the four a PATCH may name.
func CustomerPortalWritableSections() []string {
	return append(make([]string, 0, len(customerPortalWritableSections)), customerPortalWritableSections[:]...)
}
