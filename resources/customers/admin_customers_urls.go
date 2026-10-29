package customers

// The path templates for the admin customer endpoints; the client fills {api_version} from the config.
const (
	AdminCustomersPath = "/api/{api_version}/admin/customers"

	AdminCustomerPath = "/api/{api_version}/admin/customers/{customer_id}"

	AdminCustomersCountPath = "/api/{api_version}/admin/customers/count"

	AdminCustomersCheckEmailPath = "/api/{api_version}/admin/customers/check-email"

	AdminCustomersExportPath = "/api/{api_version}/admin/customers/export"

	AdminCustomersImportPath = "/api/{api_version}/admin/customers/import"

	AdminCustomersBulkResetPasswordPath = "/api/{api_version}/admin/customers/bulk-reset-password" //#nosec G101 -- a URL path template, not a credential

	AdminCustomerAddressesPath = "/api/{api_version}/admin/customers/{customer_id}/addresses"

	AdminCustomerResetPasswordPath = "/api/{api_version}/admin/customers/{customer_id}/reset-password" //#nosec G101 -- a URL path template, not a credential

	AdminCustomerChangePasswordPath = "/api/{api_version}/admin/customers/{customer_id}/change-password" //#nosec G101 -- a URL path template, not a credential

	AdminCustomerStoreCreditAdjustmentsPath = "/api/{api_version}/admin/customers/{customer_id}/store-credit/adjustments" //#nosec G101 -- a URL path template, not a credential
)

// AdminCustomersTransactionNoHeader files the new customer under this transaction number; it beats a transaction_no body field.
const AdminCustomersTransactionNoHeader = "X-Transaction-No"
