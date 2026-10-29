package taxsettings

// TaxSettingsPath is the tax settings facade's one path; it takes no sections and serves no HEAD.
const TaxSettingsPath = "/api/{api_version}/admin/tax/settings"

// TaxSettingsKey is the key the read answers under; the PUT is wrapped in TaxSettingsWriteKey instead.
const TaxSettingsKey = "tax_settings"

// TaxSettingsWriteKey is the wrapper the PUT requires, which is not the key it answers under.
const TaxSettingsWriteKey = "settings"

// TaxSettingsUnknownPathMessage is the tax router's 404 text for a path or method it does not route.
const TaxSettingsUnknownPathMessage = "Unknown tax path"
