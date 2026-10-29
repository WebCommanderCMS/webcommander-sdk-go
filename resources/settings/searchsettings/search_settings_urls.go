package searchsettings

// The path templates for the search settings endpoints; reindex_status is a path of its own, not a section.
const (
	SearchSettingsPath = "/api/{api_version}/admin/settings/search_settings"

	SearchSettingsSectionPath = "/api/{api_version}/admin/settings/search_settings/{section}"

	SearchSettingsReindexStatusPath = "/api/{api_version}/admin/settings/search_settings/reindex_status"
)

// searchSettingsSections are the three values {section} takes, in facade order.
var searchSettingsSections = [3]string{
	"searchable_content",
	"search_fields",
	"synonyms",
}

// SearchSettingsReindexingSection is the one section to leave alone: writing it triggers an irreversible reindex.
const SearchSettingsReindexingSection = "synonyms"

// SearchSettingsSections returns a fresh copy of the three values {section} takes.
func SearchSettingsSections() []string {
	return append(make([]string, 0, len(searchSettingsSections)), searchSettingsSections[:]...)
}
