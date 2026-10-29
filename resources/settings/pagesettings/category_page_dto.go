package pagesettings

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// CategoryHeading is the heading block: whether the category name is shown, and at which heading level.
type CategoryHeading struct {
	Enabled *bool
	Style   *string
}

// Subcategory is the subcategory block: whether subcategories are listed, and how that list pages.
type Subcategory struct {
	Enabled            *bool
	Display            *bool
	LoadType           *string
	PaginationPosition *string
	PerPage            *int64
	AllowUserPerPage   *bool
	PerLoad            *int64
	DisplayProducts    *bool
	ShowSummary        *bool
}

// CategoryInformation is the category_information section: the category itself and its subcategories.
type CategoryInformation struct {
	Heading         *CategoryHeading
	ShowDescription *bool
	ShowBannerImage *bool
	ShowImage       *bool
	Subcategory     *Subcategory
}

// EnvelopeKey is the section's own name, the wrapper a write must carry.
func (CategoryInformation) EnvelopeKey() string { return "category_information" }

// ToMap renders the request body.
func (c CategoryInformation) ToMap() map[string]any { return dto.ToMap(c) }

// InteractionButtons is the interaction_buttons block: which buttons each product card carries.
type InteractionButtons struct {
	Enabled       *bool
	AddToCart     *bool
	ViewDetails   *bool
	ShowOnHover   *bool
	AddToWishList *bool
	AddToCompare  *bool
	QuickView     *bool
}

// Price is the price block: whether a card shows a price, and which price.
type Price struct {
	Enabled               *bool
	ExpectToPay           *bool
	ExpectToPayWithTax    *bool
	StrikeThroughPrevious *bool
}

// VariationType is one entry of variation_combination.types; take ids from allowed_values.
type VariationType struct {
	ID   *int64
	Name *string
}

// VariationCombination is the variation_combination block: whether combinations show, and which types.
type VariationCombination struct {
	Enabled *bool
	// An element is nil where the reader would not coerce the entry.
	Types []*VariationType
}

// Rating is the rating block: whether the rating shows, and the review count beside it.
type Rating struct {
	Enabled         *bool
	ShowReviewCount *bool
}

// ViewSwitcher is the view_switcher block: whether a shopper may change the display mode, and to which.
type ViewSwitcher struct {
	Enabled *bool
	Grid    *bool
	List    *bool
	Tabular *bool
}

// ColumnsPerRow is default_display's columns_per_row block; each mode has its own max.
type ColumnsPerRow struct {
	Grid *int64
	List *int64
}

// DefaultDisplay is the default_display block: the mode a listing opens in and how wide its rows are.
type DefaultDisplay struct {
	Enabled                            *bool
	Mode                               *string
	ColumnsPerRow                      *ColumnsPerRow
	ShowVariationOptionsLabelInTabular *bool
}

// ProductLoad is the load block: how the product listing pages and where its control sits.
type ProductLoad struct {
	Enabled            *bool
	Type               *string
	PaginationPosition *string
	PerPage            *int64
	AllowUserPerPage   *bool
	PerLoad            *int64
}

// DefaultSort is the default_sort block: whether a default order is applied, and which.
type DefaultSort struct {
	Enabled *bool
	By      *string
}

// SortOptions is the sort_options block: which orders a shopper may switch to.
type SortOptions struct {
	Enabled    *bool
	Alphabetic *bool
	Featured   *bool
	New        *bool
	TopSelling *bool
	Price      *bool
	Date       *bool
}

// CategoryPageProductInformation is the product_information section: each product card and the listing.
type CategoryPageProductInformation struct {
	UseGlobalDefaults     *bool
	InteractionButtons    *InteractionButtons
	ShowSummary           *bool
	ShowSalePercentage    *bool
	ShowStockAvailability *bool
	ShowParentCategory    *bool
	Price                 *Price
	VariationCombination  *VariationCombination
	Rating                *Rating
	ViewSwitcher          *ViewSwitcher
	DefaultDisplay        *DefaultDisplay
	Load                  *ProductLoad
	DefaultSort           *DefaultSort
	SortOptions           *SortOptions
}

// EnvelopeKey is the section's own name, the wrapper a write must carry.
func (CategoryPageProductInformation) EnvelopeKey() string { return "product_information" }

func (p CategoryPageProductInformation) ToMap() map[string]any { return dto.ToMap(p) }

// CategoryPageSectionMeta is one group_meta entry: what a section is, whether writable, where managed.
type CategoryPageSectionMeta struct {
	Kind     *string
	Writable *bool
	Href     *string
}

// CategoryPageGroupMeta is the group_meta block, one entry per section; only the facade sends it.
type CategoryPageGroupMeta struct {
	CategoryInformation *CategoryPageSectionMeta
	ProductInformation  *CategoryPageSectionMeta
}

// CategoryPage is the facade's body or one section's, under the section's own key.
type CategoryPage struct {
	CategoryInformation *CategoryInformation
	ProductInformation  *CategoryPageProductInformation
	GroupMeta           *CategoryPageGroupMeta
	// FieldMetadata is the raw field_metadata block, keyed by section on the facade and by field on a section.
	FieldMetadata any
}

func (CategoryPage) EnvelopeKey() string { return "category_page" }

func (c CategoryPage) ToMap() map[string]any { return dto.ToMap(c) }

// CategoryPageShape is the typed view of the facade's body and of one section's.
func CategoryPageShape(decoded any) any { return dto.Read[CategoryPage](decoded) }
