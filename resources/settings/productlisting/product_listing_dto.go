package productlisting

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// InteractionButtons is a listing tile's add to cart and view details buttons, and whether they wait for hover.
type InteractionButtons struct {
	Enabled         *bool
	AddToCart       *bool
	ShowViewDetails *bool
	ShowOnHover     *bool
}

// ListingPrice is which prices a listing tile shows.
type ListingPrice struct {
	Enabled            *bool
	ExpectToPay        *bool
	ExpectToPayWithTax *bool
}

// Reviews is the rating stars on a listing tile and the count beside them.
type Reviews struct {
	ShowRating      *bool
	ShowReviewCount *bool
}

// VariationCombination is which variation types a listing tile combines; unknown names are 400.
type VariationCombination struct {
	Enabled *bool
	Types   []string
}

// ProductInformation is the product_information section: what each product shows on a listing tile.
type ProductInformation struct {
	InteractionButtons       *InteractionButtons
	ShowSummary              *bool
	ShowSalePercentage       *bool
	ShowStockAvailability    *bool
	ShowParentCategory       *bool
	ShowLoyaltyPoints        *bool
	Price                    *ListingPrice
	Reviews                  *Reviews
	VariationCombination     *VariationCombination
	ShowVariationsAsProducts *bool
}

// Pagination is where the pager sits and how many products a page holds; both are enforced.
type Pagination struct {
	Position         *string
	PerPage          *int64
	AllowUserPerPage *bool
}

// LoadMore is how many products each press of load more adds.
type LoadMore struct {
	PerLoad *int64
}

// Load is the load section: pagination or load more, with Type selecting between the two blocks.
type Load struct {
	Enabled    *bool
	Type       *string
	Pagination *Pagination
	LoadMore   *LoadMore
}

// Sorting is the sorting section: the order a listing arrives in; an unknown SortBy is 400.
type Sorting struct {
	Enabled *bool
	SortBy  *string
}

// SortableOptions is which orders a shopper may pick from.
type SortableOptions struct {
	Alphabetic  *bool
	Featured    *bool
	New         *bool
	TopSelling  *bool
	Price       *bool
	CreatedDate *bool
}

// SortableItems is the sortable_items section: whether a shopper may re-sort a listing, and by what.
type SortableItems struct {
	Enabled *bool
	Options *SortableOptions
}

// SectionMeta is the one group_meta entry: what the area is, whether it is writable and where.
type SectionMeta struct {
	Kind     *string
	Writable *bool
	Href     *string
}

// GroupMeta is the group_meta block — a single entry for the whole area, not one per section.
type GroupMeta struct {
	ProductListing *SectionMeta
}

// ProductListing is the area as the API sends it; ToMap renders the PATCH body, as a bare section is refused.
type ProductListing struct {
	ProductInformation *ProductInformation
	Load               *Load
	Sorting            *Sorting
	SortableItems      *SortableItems
	GroupMeta          *GroupMeta
	// FieldMetadata is the raw field_metadata block, nested as each section is.
	FieldMetadata any
}

func (ProductListing) EnvelopeKey() string { return "product_listing" }

// ToMap renders the request body.
func (p ProductListing) ToMap() map[string]any { return dto.ToMap(p) }

// ProductListingShape is the typed view of the facade's body, the one shape this area sends.
func ProductListingShape(decoded any) any { return dto.Read[ProductListing](decoded) }
