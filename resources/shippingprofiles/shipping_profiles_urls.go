package shippingprofiles

// The path templates; the client fills {api_version}.
const (
	ShippingProfilesPath = "/api/{api_version}/admin/shipping_profiles"

	ShippingProfilesCountPath = "/api/{api_version}/admin/shipping_profiles/count"

	ShippingProfilesWizardPath = "/api/{api_version}/admin/shipping_profiles/wizard"

	ShippingProfilePath = "/api/{api_version}/admin/shipping_profiles/{shipping_profile_id}"
)

// ShippingProfilesKey is the key the listing's rows arrive under.
const ShippingProfilesKey = "shipping_profiles"

// ShippingProfileKey is the key one row arrives under, and the wrapper every write takes; a bare body is 400.
const ShippingProfileKey = "shipping_profile"

// ShippingProfilesCountKey is the one key the count sub-path answers with, bare at the top of the body.
const ShippingProfilesCountKey = "total"

// ShippingProfilesPageSize is the page size applied when no limit is sent; limit takes 1 to ShippingProfilesPageSizeCap, and 0, above the cap or "all" is 400.
const ShippingProfilesPageSize = 20

const ShippingProfilesPageSizeCap = 100

// ShippingProfilesTransactionNoHeader files a new profile under this transaction number and wins over a transaction_no in the body.
const ShippingProfilesTransactionNoHeader = "X-Transaction-No"
