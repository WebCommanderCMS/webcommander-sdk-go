package seoconfiguration

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// SeoTemplates is one content type's block of global_store_settings; allowed tokens differ per block.
type SeoTemplates struct {
	DisableIndexing *bool
	// An element is nil where the reader would not coerce the entry.
	TitleTemplate       []*string
	DescriptionTemplate []*string
}

// GlobalStoreSettings is the global_store_settings section: the store name, title separator and templates.
type GlobalStoreSettings struct {
	StoreName      *string
	TitleSeparator *string
	Pages          *SeoTemplates
	Products       *SeoTemplates
	Categories     *SeoTemplates
}

// EnvelopeKey is the section's own name, which is the wrapper a write to this section must carry.
func (GlobalStoreSettings) EnvelopeKey() string { return "global_store_settings" }

// ToMap renders the request body.
func (g GlobalStoreSettings) ToMap() map[string]any { return dto.ToMap(g) }

// CustomCodeBlock is the header or footer block of custom_code: raw HTML or JavaScript on every page.
type CustomCodeBlock struct {
	Enabled *bool
	Code    *string
}

// CustomCode is the custom_code section — the two blocks whose Code reaches every storefront page.
type CustomCode struct {
	Header *CustomCodeBlock
	Footer *CustomCodeBlock
}

// EnvelopeKey is the section's own name, the wrapper a write to this section must carry.
func (CustomCode) EnvelopeKey() string { return "custom_code" }

func (c CustomCode) ToMap() map[string]any { return dto.ToMap(c) }

// Redirect is one redirect row: the store's id, and the pair of URLs.
type Redirect struct {
	ID             *int64
	SourceURL      *string
	DestinationURL *string
}

// EnvelopeKey is redirect: how one row arrives, and a body the create also accepts.
func (Redirect) EnvelopeKey() string { return "redirect" }

func (r Redirect) ToMap() map[string]any { return dto.ToMap(r) }

// RedirectsPreview is url_settings' inlined redirects block, unpaginated, so its count is not the total.
type RedirectsPreview struct {
	Items []*Redirect
}

// URLSettings is the url_settings section: URL rewrites on rename, redirects left behind, and inlined rows.
type URLSettings struct {
	RegenerateURLOnItemUpdate *bool
	AutoRedirectOnURLChange   *bool
	Redirects                 *RedirectsPreview
}

// EnvelopeKey is the section's own name, the wrapper a write to this section must carry.
func (URLSettings) EnvelopeKey() string { return "url_settings" }

func (u URLSettings) ToMap() map[string]any { return dto.ToMap(u) }

// DocumentRootFile is one file served from the document root; ConfigKey, not the name, is the delete's id.
type DocumentRootFile struct {
	ConfigKey *string
	FileName  *string
	FilePath  *string
}

// EnvelopeKey is document_root_file, the upload answer's wrapper; the upload is multipart, so no ToMap.
func (DocumentRootFile) EnvelopeKey() string { return "document_root_file" }

// DocumentRootPreview is the facade's inlined document_root block, unpaginated like RedirectsPreview.
type DocumentRootPreview struct {
	Items []*DocumentRootFile
}

// SectionMeta is one group_meta entry; url_settings' nests a fourth, Redirects, of this same type.
type SectionMeta struct {
	Kind      *string
	Writable  *bool
	Href      *string
	Redirects *SectionMeta
}

// GroupMeta is the group_meta block, one SectionMeta per thing the area declares; only the facade sends it.
type GroupMeta struct {
	GlobalStoreSettings *SectionMeta
	CustomCode          *SectionMeta
	DocumentRoot        *SectionMeta
	URLSettings         *SectionMeta
}

// SeoConfiguration is the facade's body or one section's, under the section's own key.
type SeoConfiguration struct {
	GlobalStoreSettings *GlobalStoreSettings
	CustomCode          *CustomCode
	DocumentRoot        *DocumentRootPreview
	URLSettings         *URLSettings
	GroupMeta           *GroupMeta
	// FieldMetadata is the raw field_metadata block, keyed by section on the facade and by field on a section.
	FieldMetadata any
}

func (SeoConfiguration) EnvelopeKey() string { return "seo_configuration" }

func (s SeoConfiguration) ToMap() map[string]any { return dto.ToMap(s) }

// SeoConfigurationShape is the typed view of the facade's body and of one section's.
func SeoConfigurationShape(decoded any) any { return dto.Read[SeoConfiguration](decoded) }

// RedirectShape is the typed view of a body shaped {"redirect": {…}}.
func RedirectShape(decoded any) any { return dto.Read[Redirect](decoded) }

// DocumentRootFileShape is the typed view of a body shaped {"document_root_file": {…}}.
func DocumentRootFileShape(decoded any) any { return dto.Read[DocumentRootFile](decoded) }
