package commonsettings

// The path templates; the client fills {api_version}. Neither sits under /admin.
const (
	CommonSettingsPath = "/api/{api_version}/settings"

	CreateCustomerFieldsPath = "/api/{api_version}/create-customer-fields"
)

// CommonSettingsKey is the key the read and both writes answer under, and the wrapper both writes require.
const CommonSettingsKey = "settings"

// CustomerFieldsKey lists the registration form's active fields; CustomerFieldConfigsKey configures every field.
const (
	CustomerFieldsKey       = "fields"
	CustomerFieldConfigsKey = "fields_configs"
)
