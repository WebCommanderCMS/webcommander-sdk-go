package pagesettings

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// ProductInteractionButtons is the product_interaction_buttons block: which buttons each result carries.
type ProductInteractionButtons struct {
	Enabled                    *bool
	DisplayAddToCartButton     *bool
	DisplayAddToCompareButton  *bool
	DisplayAddToWishlistButton *bool
	DisplayViewDetailsButton   *bool
	DisplayOnHover             *bool
	DisplayQuickViewButton     *bool
}

// NumberOfSearchesPerResult is the number_of_searches_per_result block: pagination and page size.
type NumberOfSearchesPerResult struct {
	Enabled      *bool
	ItemsPerPage *int64
}

// SetProductDisplayPreferences is the set_product_display_preferences block: loading and pager position.
type SetProductDisplayPreferences struct {
	Enabled                            *bool
	DisplayMode                        *string
	SelectPaginationPosition           *string
	NumberOfProductsPerPage            *int64
	AllowUsersToSelectNumberOfProducts *bool
	NumberOfProductsPerLoad            *int64
}

// SetProductSorting is the set_product_sorting block: whether a default sort applies, and which.
type SetProductSorting struct {
	Enabled     *bool
	DefaultSort *string
}

// ConfigureSortingOptions is the configure_sorting_options block; once enabled, one option is required.
type ConfigureSortingOptions struct {
	Enabled      *bool
	Alphabetical *bool
	Featured     *bool
	New          *bool
	TopSelling   *bool
	Price        *bool
	CreatedDate  *bool
}

// DisplayPrice is the display_price block: whether a result shows a price, and which price.
type DisplayPrice struct {
	Enabled                        *bool
	DisplayExpectToPayPrice        *bool
	ShowTaxWithTheExpectToPayPrice *bool
}

// ProductInformation is the product_information section: what each search result shows.
type ProductInformation struct {
	UseDefaultProductListingConfigurations *bool
	ProductInteractionButtons              *ProductInteractionButtons
	DisplayContentPage                     *bool
	NumberOfSearchesPerResult              *NumberOfSearchesPerResult
	DisplayShortDescription                *bool
	DisplayVariationsCombination           *bool
	SetProductDisplayPreferences           *SetProductDisplayPreferences
	SetProductSorting                      *SetProductSorting
	ConfigureSortingOptions                *ConfigureSortingOptions
	DisplayLoyaltyPoints                   *bool
	DisplayPrice                           *DisplayPrice
	LabelForPrice                          *string
}

// EnvelopeKey is the section's own name, the wrapper a write must carry.
func (ProductInformation) EnvelopeKey() string { return "product_information" }

// ToMap renders the request body.
func (p ProductInformation) ToMap() map[string]any { return dto.ToMap(p) }

// Seo is the seo section: two read-only markers, so it has no ToMap.
type Seo struct {
	ReadOnly *bool
	PageKind *string
}

// EnvelopeKey is the section's own name, the wrapper a write to this section would carry.
func (Seo) EnvelopeKey() string { return "seo" }

// Reinitialise is the reinitialise section, a read-only action marker with no ToMap.
type Reinitialise struct {
	Action *string
}

// EnvelopeKey is the section's own name, the wrapper a write to this section would carry.
func (Reinitialise) EnvelopeKey() string { return "reinitialise" }

// ProductPromotionPreview is the facade's product_promotion marker; the rows have their own path.
type ProductPromotionPreview struct {
	ReadOnly *bool
}

// ProductPromotion is one promotion row: store id and stamps, three required fields and six optional.
type ProductPromotion struct {
	ID                *int64
	Name              *string
	Rule              *string
	Position          *string
	SpecificKeywords  *string
	IsActive          *bool
	IsAlwaysAvailable *bool
	AvailableFromDate *string
	AvailableToDate   *string
	// An element is nil where the reader would not coerce the entry.
	ProductIds []*int64
	Created    *string
	Updated    *string
}

// EnvelopeKey is product_promotion: one row's key, and the wrapper every write to the collection needs.
func (ProductPromotion) EnvelopeKey() string { return "product_promotion" }

func (p ProductPromotion) ToMap() map[string]any { return dto.ToMap(p) }

// SectionMeta is one group_meta entry; its Href may name another resource entirely.
type SectionMeta struct {
	Kind     *string
	Writable *bool
	UILabel  *string
	Href     *string
}

// GroupMeta is the group_meta block, one SectionMeta per thing the area declares; only the facade sends it.
type GroupMeta struct {
	ProductInformation *SectionMeta
	Seo                *SectionMeta
	ProductPromotion   *SectionMeta
	Reinitialise       *SectionMeta
}

// SearchResultPage is the facade's body or one section's, both under the facade's key.
type SearchResultPage struct {
	ProductInformation *ProductInformation
	Seo                *Seo
	ProductPromotion   *ProductPromotionPreview
	Reinitialise       *Reinitialise
	GroupMeta          *GroupMeta
	// FieldMetadata is the raw field_metadata block, keyed by section on the facade and by field on a section.
	FieldMetadata any
}

func (SearchResultPage) EnvelopeKey() string { return "search_result_page" }

func (s SearchResultPage) ToMap() map[string]any { return dto.ToMap(s) }

// SearchResultPageShape is the typed view of the facade's body and of one section's.
func SearchResultPageShape(decoded any) any { return dto.Read[SearchResultPage](decoded) }

// ProductPromotionShape is the typed view of a body shaped {"product_promotion": {…}}.
func ProductPromotionShape(decoded any) any { return dto.Read[ProductPromotion](decoded) }
