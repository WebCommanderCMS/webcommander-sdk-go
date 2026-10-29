package variations

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// VariationOptionImage is an image option's stored picture; nil on text and color options.
type VariationOptionImage struct {
	URL      *string
	FileName *string
}

// VariationOption is one option of a type; Value is nil on an image option.
type VariationOption struct {
	ID      *int64
	Label   *string
	Value   *string
	Order   *int64
	Default *bool
	Image   *VariationOptionImage
}

// Variation is one variation type: a listing row, a row read and every type write's answer.
type Variation struct {
	ID         *int64
	Name       *string
	Standard   *string
	Disposable *bool
	Options    []*VariationOption
	// FieldMetadata is nil unless the call asked for it.
	FieldMetadata any
}

// EnvelopeKey is the key a row read and every type write's answer wrap the record in; a listing row has none.
func (Variation) EnvelopeKey() string { return VariationKey }

// Pagination is the listing's five-key pagination block; Limit is the page size applied.
type Pagination struct {
	Total       *int64
	Limit       *int64
	Offset      *int64
	HasPrevious *bool
	HasNext     *bool
}

// VariationPage is one page of the listing, and the schema when field_metadata is asked for.
type VariationPage struct {
	Variations    []*Variation
	FieldMetadata any
	// Pagination is nil when the body carried none.
	Pagination *Pagination
}

// VariationOptionImageInput is an image option's upload; written as filename, read back as file_name.
type VariationOptionImageInput struct {
	Base64   string `wire:"base64"`
	Filename string
}

// VariationOptionInput is one option in a type write; an ID names an option of this type to change.
type VariationOptionInput struct {
	ID    *int64
	Label string
	Value string
	// Order and Default are pointers, so an unset one is omitted rather than sent as 0 or false.
	Order   *int64
	Default *bool
	Image   *VariationOptionImageInput
}

// VariationInput is the type writes' body, {"variation": {...}}; Standard cannot change after the create.
type VariationInput struct {
	Name     string
	Standard string
	Options  []*VariationOptionInput
}

// EnvelopeKey is the wrapper the create, the replace and the merge require.
func (VariationInput) EnvelopeKey() string { return VariationKey }

// ToMap renders the request body.
func (v VariationInput) ToMap() map[string]any { return dto.ToMap(v) }

// ProductVariantExpectToPay is a variant's expected-price display.
type ProductVariantExpectToPay struct {
	IsEnabled *bool
	Price     *float64
}

// ProductVariantPricing is a variant's prices; the bulk write's full answer carries the first four fields.
type ProductVariantPricing struct {
	BasePrice      *float64
	CostPrice      *float64
	SalePrice      *float64
	IsOnSale       *bool
	SaleAmountType *string
	HidePrice      *bool
	IsCallForPrice *bool
	ExpectToPay    *ProductVariantExpectToPay
	TieredPricing  any
}

// ProductVariantInventory is a variant's stock; the bulk write's full answer carries the first three fields.
type ProductVariantInventory struct {
	IsEnabled         *bool
	Quantity          *int64
	LowStockThreshold *int64
	IsSoldOut         *bool
	SoldOutLabel      *string
}

// ProductVariantOrdering is a variant's order quantities; MaxQuantity 0 is no maximum.
type ProductVariantOrdering struct {
	MinQuantity                *int64
	MaxQuantity                *int64
	MultipleOf                 *int64
	IsMultipleQuantityRequired *bool
}

// ProductVariantDimensions is a variant's size and weight.
type ProductVariantDimensions struct {
	Length *float64
	Width  *float64
	Height *float64
	Weight *float64
}

// ProductVariantBadges is a variant's two badges.
type ProductVariantBadges struct {
	IsFeatured   *bool
	IsNewArrival *bool
}

// ProductVariant is one variant of a product; Key joins its option labels with /.
type ProductVariant struct {
	VariationID     *int64
	Available       *bool
	IsBase          *bool
	Key             *string
	Name            *string
	SKU             *string
	CustomSKU       *string
	Title           *string
	Summary         *string
	Description     *string
	Barcode         *string
	CountryOfOrigin *string
	HSCode          *string
	Pricing         *ProductVariantPricing
	Inventory       *ProductVariantInventory
	Ordering        *ProductVariantOrdering
	Dimensions      *ProductVariantDimensions
	Badges          *ProductVariantBadges
	Images          []any
}

// ProductVariationGroup is a product's variants, under ProductKey: its type names joined with /.
type ProductVariationGroup struct {
	ProductKey *string
	Variants   []*ProductVariant
}

// ProductVariations is the product read's body; Variations is empty for a product without variants.
type ProductVariations struct {
	Variations []*ProductVariationGroup
}

// ProductVariantOptionReference is one option a variant carries, in the bulk write's full answer.
type ProductVariantOptionReference struct {
	ID    *int64
	Type  *string
	Label *string
}

// ProductVariantUpdate is one written variant: ID and Name, and the rest with response=full.
type ProductVariantUpdate struct {
	ID              *int64
	Name            *string
	SKU             *string
	CustomSKU       *string
	Barcode         *string
	CountryOfOrigin *string
	HSCode          *string
	Available       *bool
	Options         []*ProductVariantOptionReference
	Pricing         *ProductVariantPricing
	Inventory       *ProductVariantInventory
	Ordering        *ProductVariantOrdering
	TransactionNo   *string
	URL             *string
}

// ProductVariationsBulkResult is the bulk write's answer; Failed stays empty, a refused id refuses the whole write.
type ProductVariationsBulkResult struct {
	Status  *string
	Message *string
	Updated []*ProductVariantUpdate
	Failed  []any
}

// ProductVariationsBulkInput is the bulk write's body, no envelope: the variant ids and at least one of the six fields.
type ProductVariationsBulkInput struct {
	VariationIds []int64
	// The six fields are pointers, so an unset one is omitted rather than sent as 0 or false.
	Available        *bool
	BasePrice        *float64
	CostPrice        *float64
	AvailableStock   *int64
	MinOrderQuantity *int64
	MaxOrderQuantity *int64
}

// ToMap renders the request body.
func (p ProductVariationsBulkInput) ToMap() map[string]any { return dto.ToMap(p) }

// ProductVariantInput is one variant to import; Key joins one option label per type with / and is required, and a key the product has is that variant, changed.
type ProductVariantInput struct {
	Key             string
	Name            string
	SKU             string
	CustomSKU       string
	Title           string
	Summary         string
	Description     string
	Barcode         string
	CountryOfOrigin string
	HSCode          string
	Pricing         *ProductVariantPricing
	Inventory       *ProductVariantInventory
	Ordering        *ProductVariantOrdering
	Dimensions      *ProductVariantDimensions
	Badges          *ProductVariantBadges
	Images          []any
}

// ProductVariationGroupInput is one group to import: ProductKey names the types, joined with /, and Variants must not be empty.
type ProductVariationGroupInput struct {
	ProductKey string
	Variants   []*ProductVariantInput
}

// ProductVariationsImportInput is the import's body, no envelope; Variations must not be empty.
type ProductVariationsImportInput struct {
	TransactionNo string
	Variations    []*ProductVariationGroupInput
}

// ToMap renders the request body.
func (p ProductVariationsImportInput) ToMap() map[string]any { return dto.ToMap(p) }

// ProductVariationsImportResult is the import's answer: Status, a Message counting the variants, and TotalVariants.
type ProductVariationsImportResult struct {
	Status        *string
	Message       *string
	TotalVariants *int64
}

// VariationShape is the typed view of a body shaped {"variation": {…}}.
func VariationShape(decoded any) any { return dto.Read[Variation](decoded) }

// VariationPageShape is the typed view of the listing body.
func VariationPageShape(decoded any) any { return dto.Read[VariationPage](decoded) }

// ProductVariationsShape is the typed view of the product read's body.
func ProductVariationsShape(decoded any) any { return dto.Read[ProductVariations](decoded) }

// ProductVariationsBulkResultShape is the typed view of the bulk write's answer.
func ProductVariationsBulkResultShape(decoded any) any {
	return dto.Read[ProductVariationsBulkResult](decoded)
}

// ProductVariationsImportResultShape is the typed view of the import's answer.
func ProductVariationsImportResultShape(decoded any) any {
	return dto.Read[ProductVariationsImportResult](decoded)
}
