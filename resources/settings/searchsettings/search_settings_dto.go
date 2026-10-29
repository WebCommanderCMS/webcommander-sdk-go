package searchsettings

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// ExcludedRecord is one entry of an excluded list as a read sends it; a write takes ids or these records.
type ExcludedRecord struct {
	ID   *int64
	Name *string
}

// SearchableContentGroup is one content type the site search covers; DisplayHeading caps at 50 characters.
type SearchableContentGroup struct {
	Enabled        *bool
	DisplayHeading *string
	Excluded       []*ExcludedRecord
}

// ProductCategoriesContent is the product_categories block, the one content group with a fourth field.
type ProductCategoriesContent struct {
	Enabled            *bool
	DisplayHeading     *string
	Excluded           []*ExcludedRecord
	HideParentCategory *bool
}

// SearchableContent is the searchable_content section: the search box text and the content it reaches.
type SearchableContent struct {
	PlaceholderText   *string
	Product           *SearchableContentGroup
	ProductCategories *ProductCategoriesContent
	Blog              *SearchableContentGroup
	OtherPages        *SearchableContentGroup
}

// EnvelopeKey is the section's own name, the only wrapper a write takes.
func (SearchableContent) EnvelopeKey() string { return "searchable_content" }

// ToMap renders the request body.
func (s SearchableContent) ToMap() map[string]any { return dto.ToMap(s) }

// SearchFieldPriority is one row of search_fields.fields; Priority is bounded 1 to 10.
type SearchFieldPriority struct {
	Field    *string
	Priority *int64
}

// SearchFields is the search_fields section; each field and each priority unique, title at priority 1.
type SearchFields struct {
	MatchKeywordsExactly *bool
	MatchFullSentences   *bool
	Fields               []*SearchFieldPriority
}

// EnvelopeKey is the section's own name, the only wrapper a write takes.
func (SearchFields) EnvelopeKey() string { return "search_fields" }

func (s SearchFields) ToMap() map[string]any { return dto.ToMap(s) }

// SynonymGroup is one group of synonyms.groups, at least two keywords that a query treats as each other.
type SynonymGroup struct {
	Keywords []*string
}

// Synonyms is the synonyms section; a write to it triggers an irreversible reindex.
type Synonyms struct {
	Groups []*SynonymGroup
}

// EnvelopeKey is the section's own name, the only wrapper a write would take.
func (Synonyms) EnvelopeKey() string { return "synonyms" }

func (s Synonyms) ToMap() map[string]any { return dto.ToMap(s) }

// SectionMeta is one group_meta entry: what a section is, whether it is writable, where it is managed.
type SectionMeta struct {
	Kind     *string
	Writable *bool
	Href     *string
}

// GroupMeta is the group_meta block, one SectionMeta per section; only the facade sends it.
type GroupMeta struct {
	SearchableContent *SectionMeta
	SearchFields      *SectionMeta
	Synonyms          *SectionMeta
}

// ReindexStatus is the reindex_status sub-resource's bare body: where the Elasticsearch reindex stands.
type ReindexStatus struct {
	State      *string
	StartedAt  any
	FinishedAt any
	Error      any
}

// ReindexStatusShape is the typed view of that sub-resource's body, the one shape it sends.
func ReindexStatusShape(decoded any) any { return dto.Read[ReindexStatus](decoded) }

// SearchSettings is the facade's body or one section's; ToMap renders the facade PATCH, not a section's.
type SearchSettings struct {
	SearchableContent *SearchableContent
	SearchFields      *SearchFields
	Synonyms          *Synonyms
	GroupMeta         *GroupMeta
	// FieldMetadata is the raw field_metadata block.
	FieldMetadata any
}

func (SearchSettings) EnvelopeKey() string { return "search_settings" }

func (s SearchSettings) ToMap() map[string]any { return dto.ToMap(s) }

// SearchSettingsShape is the typed view of the facade's body and of one section's.
func SearchSettingsShape(decoded any) any { return dto.Read[SearchSettings](decoded) }
