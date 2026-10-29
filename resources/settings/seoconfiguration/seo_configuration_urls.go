package seoconfiguration

// The path templates for the SEO configuration settings endpoints; the client fills {api_version}.
const (
	SeoConfigurationPath = "/api/{api_version}/admin/settings/seo_configuration"

	SeoConfigurationSectionPath = "/api/{api_version}/admin/settings/seo_configuration/{section}"

	SeoConfigurationDocumentRootPath = "/api/{api_version}/admin/settings/seo_configuration/document_root"

	SeoConfigurationDocumentRootFilePath = "/api/{api_version}/admin/settings/seo_configuration/document_root/{config_key}"

	SeoConfigurationRedirectsPath = "/api/{api_version}/admin/settings/seo_configuration/redirects"

	SeoConfigurationRedirectPath = "/api/{api_version}/admin/settings/seo_configuration/redirects/{redirect_id}"
)

// sections are the three values {section} takes; document_root and redirects are collections of their own.
var sections = [3]string{"global_store_settings", "custom_code", "url_settings"}

// SeoConfigurationSections returns a fresh copy of the three values {section} takes.
func SeoConfigurationSections() []string {
	return append(make([]string, 0, len(sections)), sections[:]...)
}
