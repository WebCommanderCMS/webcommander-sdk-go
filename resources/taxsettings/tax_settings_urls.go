package taxsettings

// TaxSettingsPath is the tax settings facade's one path; it takes no sections.
const TaxSettingsPath = "/api/{api_version}/admin/settings/tax"

// TaxSettingsKey is the key the read answers under, and the key both writes are wrapped in.
const TaxSettingsKey = "tax_settings"

// TaxSettingsWriteKey is the wrapper PUT and PATCH require; the store also accepts tax.
const TaxSettingsWriteKey = "tax_settings"

// TaxSettingsUnknownPathMessage is the tax router's 404 text for a path or method it does not route; the old /admin/tax/settings path answers it.
const TaxSettingsUnknownPathMessage = "Unknown tax path"
