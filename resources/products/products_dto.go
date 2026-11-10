package products

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// ProductUrls is the product's slug and its absolute storefront address.
type ProductUrls struct {
	URL        *string
	ProductURL *string
}

// ProductNamedReference is an id and a name: an author, a page, a customer or a customer group.
type ProductNamedReference struct {
	ID   *int64
	Name *string
}

// ProductParent is one category the product is filed under.
type ProductParent struct {
	ID              *int64
	Name            *string
	IsInTrash       *bool
	IsParentInTrash *bool
	IsDisposable    *bool
}

// ProductSelectedCustomer is one customer an audience names: ID is the customer's customer_id, InternalID its internal_id.
type ProductSelectedCustomer struct {
	ID         *int64
	InternalID *int64
	Name       *string
}

// ProductSelectedCustomers is who available_for: selected admits, or whom a restriction names; both lists are empty when nobody is named.
type ProductSelectedCustomers struct {
	Customers []*ProductSelectedCustomer
	Groups    []*ProductNamedReference
}

// ProductSelectedCustomersInput is an audience to write: customer and group ids; a group written as {"id": …} is 400. Pass empty slices, not nil, to send an empty audience.
type ProductSelectedCustomersInput struct {
	Customers []int64
	Groups    []int64
}

// ProductWeekday is one day of the weekday schedule.
type ProductWeekday struct {
	Day         *string
	StartHour   *int64
	StartMinute *int64
	EndHour     *int64
	EndMinute   *int64
	IsFullDay   *bool
	Enabled     *bool
}

// ProductAvailableOnWeekdays is the weekday schedule: a switch and the seven days.
type ProductAvailableOnWeekdays struct {
	Enabled  *bool
	Weekdays []*ProductWeekday
}

// ProductProperties is the dimensions as written; the product's top level repeats the four numbers read-only.
type ProductProperties struct {
	Length     *float64
	Width      *float64
	Height     *float64
	Weight     *float64
	LengthUnit *string
	WeightUnit *string
}

// ProductTaxSettings is the per-product tax overrides, applied only when OverwriteDefaultConfig is true.
type ProductTaxSettings struct {
	OverwriteDefaultConfig *bool
	IsTaxExempt            *bool
	ShowPriceWithTax       *bool
	IsPriceWithTax         *bool
	SystemDefaultTax       *bool
	TaxMessage             *string
	BasePriceRounding      *string
	DefaultTaxCode         *string
	CustomTaxProfileID     *string
}

// ProductVariationPriceDisplay is which variant price a listing shows, and its labels; LowestPriceLabel and HighestPriceLabel are second names for the two LabelFor fields.
type ProductVariationPriceDisplay struct {
	Method               *string
	LabelForLowestPrice  *string
	LabelForHighestPrice *string
	LowestPriceLabel     *string
	HighestPriceLabel    *string
}

// ProductWarehouseStock is one location's units in the legacy shape; Warehouse is the location id as text.
type ProductWarehouseStock struct {
	Warehouse *string
	Stock     *int64
}

// ProductStockLocation is one location's holding of the product; exactly one row of a product holds IsDefault.
type ProductStockLocation struct {
	Warehouse     *int64
	WarehouseName *string
	Stock         *int64
	// LowStockLevel is null on a row written without one.
	LowStockLevel *int64
	IsDefault     *bool
}

// ProductPricingAndStock is the pricing and stock block; Locations matches the stock sub-path's rows, and LowLevelStock follows the default row's low_stock_level.
type ProductPricingAndStock struct {
	TrackInventory                       *bool
	AvailableStock                       *int64
	LowLevelStock                        *int64
	Warehouses                           []*ProductWarehouseStock
	MinimumOrderQuantity                 *int64
	MaximumOrderQuantity                 *int64
	TaxProfile                           any
	ShippingProfile                      any
	EnableMultipleOrder                  *bool
	MultipleOrderQuantity                *int64
	OnSale                               *bool
	OnSalePriceType                      *string
	OnSalePrice                          *float64
	ExpectedToPay                        *bool
	ExpectedToPayPrice                   *float64
	EnableCallForPrice                   *bool
	RestrictPriceFor                     *string
	RestrictPriceForExceptCustomers      []any
	RestrictPurchaseFor                  *string
	RestrictPurchaseForSelectedCustomers []any
	HidePrice                            *bool
	EnableCombinedStockSellingLimit      *bool
	CombinedStockSellingLimit            *int64
	EnableAdvancePricing                 *bool
	Suppliers                            []any
	IsMultiLocationEnabled               *bool
	Locations                            []*ProductStockLocation
}

// ProductImage is one product image.
type ProductImage struct {
	ID        *int64
	Filename  *string
	AltText   *string
	Thumbnail *string
	Link      *string
}

// ProductVideo is one video, on the product or on a variant.
type ProductVideo struct {
	ID        *int64
	Filename  *string
	Thumbnail *string
	Link      *string
}

// ProductImageAndVideo is the media block.
type ProductImageAndVideo struct {
	Images []*ProductImage
	Videos []*ProductVideo
	Specs  []any
}

// ProductAdvanceInformation is identifiers and origin; mpn, upc, barcode and hs_code are written here, the top level refusing them.
type ProductAdvanceInformation struct {
	GlobalTradeItemNumber    *string
	EnableProductIdentifiers *bool
	GTIN                     *string
	MPN                      *string
	UPC                      *string
	Barcode                  *string
	CountryOfOrigin          *string
	HSCode                   *string
	ProductCondition         *string
	Model                    *string
}

// ProductFlag is a flag and the label shown with it.
type ProductFlag struct {
	Enabled *bool
	Label   *string
}

// ProductFlags is the flags in blocks: is_featured, is_new, sold_out and model again; SoldOut is read-only.
type ProductFlags struct {
	Featured    *ProductFlag
	NewProduct  *bool
	SoldOut     *ProductFlag
	ModelNumber *string
}

// ProductDimension is product_properties again; the two units are read-only here.
type ProductDimension struct {
	LengthUnit *string
	WeightUnit *string
	Width      *float64
	Height     *float64
	Length     *float64
	Weight     *float64
}

// ProductIdentifiers is advance_information again: a write here changes it there, Condition being its product_condition.
type ProductIdentifiers struct {
	Enabled         *bool
	GTIN            *string
	Condition       *string
	Barcode         *string
	HSCode          *string
	CountryOfOrigin *string
	MPN             *string
	UPC             *string
}

// ProductWebtool is webtool_disable_tracking again, read-only.
type ProductWebtool struct {
	DisableTracking *bool
}

// ProductAdvanced is the advanced block.
type ProductAdvanced struct {
	AdvanceInformation     *ProductAdvanceInformation
	WebtoolDisableTracking *bool
	MetaTags               []any
	ProductFlags           *ProductFlags
	ProductDimension       *ProductDimension
	ProductIdentifiers     *ProductIdentifiers
	Webtool                *ProductWebtool
}

// ProductLink is a related, similar or up-sell product; a write takes {"id": …}, and a list of bare ids in UpSellProducts is answered 200 and leaves it empty.
type ProductLink struct {
	ID      *int64
	Name    *string
	URL     *string
	Status  *string
	SKU     *string
	InStock *bool
	Price   *float64
}

// ProductBundle is the bundle block; Product.Bundle is nil on a product that is not a bundle.
type ProductBundle struct {
	PriceFixed       *bool
	QuantityFlexible *bool
	Items            []any
}

// ProductSeoConfig is one SEO setting on the product.
type ProductSeoConfig struct {
	ID        *int64
	Value     *string
	Type      *string
	ConfigKey *string
	MetaTags  []any
}

// ProductTagValue is one value of a product tag.
type ProductTagValue struct {
	ID    *int64
	Title *string
	Value *string
	URL   *string
}

// ProductTag is one tag on the product, with its values.
type ProductTag struct {
	ID                   *int64
	TagID                *int64
	Name                 *string
	Values               []*ProductTagValue
	DisplayInProductPage *bool
	EnableExpandedView   *bool
}

// ProductCustomProperty is one custom property.
type ProductCustomProperty struct {
	ID          *int64
	Label       *string
	Description *string
}

// ProductReviewer is who wrote a review; Name and Email again.
type ProductReviewer struct {
	Name  *string
	Email *string
}

// ProductReview is one review. ReviewedProductID, Comment and CreatedAt repeat productId, review and date; CreatedAt carries a Z. A single review answers under review, which inside a row is its text.
type ProductReview struct {
	ID *int64
	// ProductID travels as productId, the camelCase key in the body.
	ProductID *int64 `wire:"productId"`
	Name      *string
	Email     *string
	Review    *string
	Rating    *float64
	Date      *string
	Status    *string
	Files     []any
	// ReviewedProductID travels as product_id, beside productId.
	ReviewedProductID *int64 `wire:"product_id"`
	Reviewer          *ProductReviewer
	Comment           *string
	CreatedAt         *string
}

// EnvelopeKey is the key a single review answers under.
func (ProductReview) EnvelopeKey() string { return ProductReviewKey }

// ProductVariationTypeOption is one option of a variation type the product uses.
type ProductVariationTypeOption struct {
	ID      *int64
	Label   *string
	Value   *string
	Default *bool
}

// ProductVariationType is one variation type the product uses, with its options.
type ProductVariationType struct {
	ID       *int64
	Name     *string
	Standard *string
	Options  []*ProductVariationTypeOption
}

// ProductCombinationOptionType is the type one option of a combination belongs to.
type ProductCombinationOptionType struct {
	ID           *int64
	Role         *string
	Name         *string
	Standard     *string
	IsDisposable *bool
}

// ProductCombinationOption is one option of a combination; Type is a list of one type.
type ProductCombinationOption struct {
	ID               *int64
	Index            *int64
	Type             []*ProductCombinationOptionType
	Label            *string
	Value            *string
	ColorHasSwatch   *bool
	ColorSwatchValue *string
	ImageBaseURL     *string
}

// ProductCombinationImage is one variant image.
type ProductCombinationImage struct {
	ID        *int64
	Thumbnail *string
	URL       *string
}

// ProductCombinationPriceAndStock is a variant's pricing and stock block.
type ProductCombinationPriceAndStock struct {
	TrackInventory            *bool
	AvailableStock            *int64
	LowLevelStock             *int64
	MinimumOrderQuantity      *int64
	MaximumOrderQuantity      *int64
	SupportedMaxOrderQuantity *int64
	EnableMultipleOrder       *bool
	MultipleOrderQuantity     *int64
	OnSale                    *bool
	SoldOut                   *bool
	SoldOutLabel              *string
	OnSalePriceType           *string
	OnSalePrice               *float64
	ExpectedToPay             *bool
	ExpectedToPayPrice        *float64
	EnableCallForPrice        *bool
	HidePrice                 *bool
}

// ProductCombinationDetails is one variant.
type ProductCombinationDetails struct {
	VariationID             *int64
	Available               *bool
	ProductName             *string
	SKU                     *string
	CustomSKU               *string
	Urls                    *ProductUrls
	MetaTags                []any
	IsInventoryEnabled      *bool
	AvailableStock          *int64
	LowStockLevel           *int64
	Spec                    *string
	IsMultipleOrderQuantity *bool
	IsOnSale                *bool
	IsExpectToPay           *bool
	IsCallForPriceEnabled   *bool
	IsNew                   *bool
	IsFeatured              *bool
	OnSaleAmountType        *string
	Title                   *string
	Heading                 *string
	Summary                 *string
	Description             *string
	Model                   *string
	ProductCondition        *string
	FeaturedLabel           *string
	BasePrice               *float64
	CostPrice               *float64
	CompareAtPrice          *float64
	SalePrice               *float64
	ExpectToPayPrice        *float64
	// Weight, Height, Length and Width arrive as a number on some products and as text on others.
	Weight                    any
	Height                    any
	Length                    any
	Width                     any
	DisplayPrice              *float64
	PreviousPrice             *float64
	PriceAndStock             *ProductCombinationPriceAndStock
	Images                    []*ProductCombinationImage
	Videos                    []*ProductVideo
	Specs                     []any
	SupportedMaxOrderQuantity *int64
	IsBase                    *bool
	Barcode                   *string
	CountryOfOrigin           *string
	HSCode                    *string
	Tax                       *float64
	TaxMessage                *string
}

// ProductCombination is one option combination and the variant it makes.
type ProductCombination struct {
	Options []*ProductCombinationOption
	Details *ProductCombinationDetails
}

// ProductVariation is the variation block; {} on a product with no variations, so every field reads nil.
type ProductVariation struct {
	Types                 []*ProductVariationType
	AvailableCombinations []*ProductCombination
}

// ProductCategory is one category the product is filed under, as Categories lists it; Parents again.
type ProductCategory struct {
	ID            *int64
	Name          *string
	InTrash       *bool
	ParentInTrash *bool
	Disposable    *bool
}

// ProductMediaImage is one image in the media block.
type ProductMediaImage struct {
	ID           *int64
	Filename     *string
	Position     *int64
	Primary      *bool
	AltText      *string
	URL          *string
	ThumbnailURL *string
}

// ProductMediaVideo is one video in the media block.
type ProductMediaVideo struct {
	ID           *int64
	Filename     *string
	Position     *int64
	URL          *string
	ThumbnailURL *string
}

// ProductMedia is the media block. A write names files the product already holds, by id; another id is 400. SpecificationFile was null on every product read, so it is untyped.
type ProductMedia struct {
	Images            []*ProductMediaImage
	Videos            []*ProductMediaVideo
	SpecificationFile any
}

// ProductWeekdayWindow is one day of the schedule; times are HH:mm.
type ProductWeekdayWindow struct {
	Enabled   *bool
	FullDay   *bool
	StartTime *string
	EndTime   *string
}

// ProductScheduleWeekdays is the seven days of the schedule, each always present.
type ProductScheduleWeekdays struct {
	Monday    *ProductWeekdayWindow
	Tuesday   *ProductWeekdayWindow
	Wednesday *ProductWeekdayWindow
	Thursday  *ProductWeekdayWindow
	Friday    *ProductWeekdayWindow
	Saturday  *ProductWeekdayWindow
	Sunday    *ProductWeekdayWindow
}

// ProductSchedule is when the product is on sale. Type is schedule_type (date_range or custom_schedule), Preset is custom_schedule_option (all_days, weekdays, business_hour, weekends or custom).
type ProductSchedule struct {
	Enabled   *bool
	Type      *string
	StartDate *string
	EndDate   *string
	Preset    *string
	Weekdays  *ProductScheduleWeekdays
}

// ProductPasswordProtection is password_protected and password again; a set password reads ***.
type ProductPasswordProtection struct {
	Enabled  *bool
	Password *string
}

// ProductAvailableForPurchase is availability in one block: Enabled is available, Active is administrative_status, WhoCanView is available_for. RestrictForSelectedCustomers is read-only.
type ProductAvailableForPurchase struct {
	Enabled                      *bool
	Active                       *bool
	WhoCanView                   *string
	SelectedCustomers            *ProductSelectedCustomers
	RestrictForSelectedCustomers *ProductSelectedCustomers
	Schedule                     *ProductSchedule
	PasswordProtection           *ProductPasswordProtection
}

// ProductAvailableForPurchaseInput is the availability block to write. While Enabled is false the store answers 200 and keeps none of WhoCanView, Schedule or PasswordProtection.
type ProductAvailableForPurchaseInput struct {
	Enabled            *bool
	Active             *bool
	WhoCanView         string
	SelectedCustomers  *ProductSelectedCustomersInput
	Schedule           *ProductSchedule
	PasswordProtection *ProductPasswordProtection
}

// ProductOnSale is the sale price; AmountType is on_sale_price_type.
type ProductOnSale struct {
	Enabled         *bool
	DiscountedPrice *float64
	AmountType      *string
}

// ProductQuantityMultiple is order in multiples of MultipleOf.
type ProductQuantityMultiple struct {
	Enabled    *bool
	MultipleOf *int64
}

// ProductExpectToPay is the expect-to-pay price.
type ProductExpectToPay struct {
	Enabled *bool
	Price   *float64
}

// ProductAdvancedPricing is the advanced pricing switches.
type ProductAdvancedPricing struct {
	Enabled          *bool
	QuantityMultiple *ProductQuantityMultiple
	ExpectToPay      *ProductExpectToPay
	CallForPrice     *bool
}

// ProductPricing is prices in one block; RrpPrice is compare_at_price. On a create BasePrice here stands in for the top-level one.
type ProductPricing struct {
	BasePrice       *float64
	CostPrice       *float64
	RrpPrice        *float64
	OnSale          *ProductOnSale
	AdvancedPricing *ProductAdvancedPricing
}

// ProductInventoryTracking is stock tracking; Suppliers echoes the ids written.
type ProductInventoryTracking struct {
	Enabled        *bool
	AvailableStock *int64
	LowStockLevel  *int64
	Suppliers      []any
}

// ProductCombinedStockLimit is the combined stock selling limit.
type ProductCombinedStockLimit struct {
	Enabled *bool
	Units   *int64
}

// ProductInventory is stock in one block; Preorder is preorder_enabled.
type ProductInventory struct {
	Tracking           *ProductInventoryTracking
	CombinedStockLimit *ProductCombinedStockLimit
	Preorder           *bool
}

// ProductOrdering is the order quantity limits; a MaxQuantity of 0 is no maximum.
type ProductOrdering struct {
	MinQuantity *int64
	MaxQuantity *int64
}

// ProductRestriction is a price or purchase restriction. Only Mode is written; Enabled is read-only and Customers is answered 200 but kept empty.
type ProductRestriction struct {
	Enabled   *bool
	Mode      *string
	Customers []any
}

// ProductRestrictions is restrict_price_for and restrict_purchase_for in one block; Enabled and DisplayOnly are read-only.
type ProductRestrictions struct {
	Enabled     *bool
	Price       *ProductRestriction
	Purchase    *ProductRestriction
	DisplayOnly *bool
}

// ProductPriceIncludesTax is whether base_price includes tax, and how it is rounded.
type ProductPriceIncludesTax struct {
	Enabled           *bool
	BasePriceRounding *string
}

// ProductTax is tax_settings in one block; Enabled is overwrite_default_config. A write takes CustomTaxProfile as an id or {"id": …} of a custom tax profile.
type ProductTax struct {
	Enabled          *bool
	Exempt           *bool
	CustomTaxProfile *ProductNamedReference
	Message          *string
	ShowPriceWithTax *bool
	PriceIncludesTax *ProductPriceIncludesTax
}

// ProductShipping is the shipping profile; Enabled is read-only and true while a profile is set. A create drops Profile; an update takes it as an id or {"id": …}.
type ProductShipping struct {
	Enabled *bool
	Profile *ProductNamedReference
}

// ProductLinks is API paths of the product and its sub-paths, read-only.
type ProductLinks struct {
	Self       *string
	Variations *string
	Stock      *string
	Reviews    *string
}

// Product is one product: the row read's body and every row write's answer.
type Product struct {
	ID                           *int64
	Name                         *string
	SKU                          *string
	Status                       *string
	Urls                         *ProductUrls
	Heading                      *string
	Parents                      []*ProductParent
	CreatedBy                    *ProductNamedReference
	CreatedAt                    *string
	UpdatedAt                    *string
	TransactionNo                *string
	ProductType                  *string
	CustomType                   *string
	IsVirtual                    *bool
	InvoiceNote                  *string
	Available                    *bool
	ProductPage                  *ProductNamedReference
	ProductLayout                any
	CustomClass                  *string
	AvailableOnDateRange         *bool
	AvailableFromDate            *string
	AvailableToDate              *string
	AvailableFor                 *string
	SelectedCustomers            *ProductSelectedCustomers
	RestrictForSelectedCustomers *ProductSelectedCustomers
	PasswordProtected            *bool
	// Password reads *** when one is set and "" otherwise.
	Password               *string
	AvailableOnWeekdays    *ProductAvailableOnWeekdays
	BasePrice              *float64
	CostPrice              *float64
	CompareAtPrice         *float64
	IsOnSale               *bool
	SalePrice              *float64
	SoldOut                *bool
	Length                 *float64
	Width                  *float64
	Height                 *float64
	Weight                 *float64
	SoldOutLabel           *string
	IsFeatured             *bool
	FeaturedLabel          *string
	IsNew                  *bool
	Featured               *bool
	New                    *bool
	Model                  *string
	ProductProperties      *ProductProperties
	Version                *int64
	CreatedOn              *string
	LastUpdatedOn          *string
	LastUpdatedBy          *ProductNamedReference
	ScheduleType           *string
	CustomScheduleOption   *string
	TaxSettings            *ProductTaxSettings
	VariationPriceDisplay  *ProductVariationPriceDisplay
	AdministrativeStatus   *bool
	IsDraft                *bool
	ProductSummary         *string
	ProductDescription     *string
	ProductPricingAndStock *ProductPricingAndStock
	ImageAndVideo          *ProductImageAndVideo
	Advanced               *ProductAdvanced
	RelatedProducts        []*ProductLink
	SimilarProducts        []*ProductLink
	Bundle                 *ProductBundle
	SeoConfigs             []*ProductSeoConfig
	Tags                   []*ProductTag
	CustomProperties       []*ProductCustomProperty
	Reviews                []*ProductReview
	ProductVariation       *ProductVariation
	PreorderEnabled        *bool
	HasVariations          *bool
	// The block-shaped copy of the fields above; a write to either changes both.
	ProductURL           *string
	UpdatedBy            *ProductNamedReference
	Virtual              *bool
	Summary              *string
	Description          *string
	Categories           []*ProductCategory
	UpSellProducts       []*ProductLink
	Media                *ProductMedia
	AvailableForPurchase *ProductAvailableForPurchase
	Pricing              *ProductPricing
	Inventory            *ProductInventory
	Ordering             *ProductOrdering
	Restrictions         *ProductRestrictions
	Tax                  *ProductTax
	Shipping             *ProductShipping
	Links                *ProductLinks
	// FieldMetadata is nil unless the call asked for it; it sits beside the envelope.
	FieldMetadata any
}

// EnvelopeKey is the key the row read and every row write answer under; no row write takes it.
func (Product) EnvelopeKey() string { return ProductKey }

// Pagination is the listing's five-key pagination block; Limit is the page size applied.
type Pagination struct {
	Total       *int64
	Limit       *int64
	Offset      *int64
	HasPrevious *bool
	HasNext     *bool
}

// ProductPage is one page of the listing: full product records, where they sit and any schema.
type ProductPage struct {
	Products      []*Product
	Pagination    *Pagination
	FieldMetadata any
}

// ProductInput is the create's, the PUT's and the PATCH's bare record; an update merges it at every level, so only what is set is sent and changed, and a PUT also needs the ten fields UpdateProduct lists. product_url, links and restrict_for_selected_customers are read-only, and updated_by is answered 200 and ignored.
type ProductInput struct {
	Name string
	SKU  string
	// TransactionNo is stored on a create; an update ignores it.
	TransactionNo          string
	Heading                string
	ProductType            string
	CustomType             string
	IsVirtual              *bool
	InvoiceNote            string
	Available              *bool
	AvailableOnDateRange   *bool
	AvailableFromDate      string
	AvailableToDate        string
	AvailableFor           string
	SelectedCustomers      *ProductSelectedCustomersInput
	PasswordProtected      *bool
	Password               string
	AvailableOnWeekdays    *ProductAvailableOnWeekdays
	BasePrice              *float64
	CostPrice              *float64
	CompareAtPrice         *float64
	IsOnSale               *bool
	SalePrice              *float64
	IsFeatured             *bool
	FeaturedLabel          string
	IsNew                  *bool
	ProductProperties      *ProductProperties
	ScheduleType           string
	CustomScheduleOption   string
	TaxSettings            *ProductTaxSettings
	VariationPriceDisplay  *ProductVariationPriceDisplay
	AdministrativeStatus   *bool
	ProductSummary         string
	ProductDescription     string
	ProductPricingAndStock *ProductPricingAndStock
	ImageAndVideo          *ProductImageAndVideo
	Advanced               *ProductAdvanced
	// Parents takes category ids.
	Parents         []*int64
	RelatedProducts []any
	SimilarProducts []any
	Tags            []any
	PreorderEnabled *bool
	Virtual         *bool
	Summary         string
	Description     string
	// Categories takes category ids; [] is answered 200 and changes nothing.
	Categories           []int64
	UpSellProducts       []*ProductLink
	Media                *ProductMedia
	AvailableForPurchase *ProductAvailableForPurchaseInput
	Pricing              *ProductPricing
	Inventory            *ProductInventory
	Ordering             *ProductOrdering
	Restrictions         *ProductRestrictions
	Tax                  *ProductTax
	Shipping             *ProductShipping
}

// ToMap renders the request body.
func (p ProductInput) ToMap() map[string]any { return dto.ToMap(p) }

// ProductReviews is the review import's answer: every review the product holds, the ones just added last.
type ProductReviews struct {
	Reviews []*ProductReview
}

// EnvelopeKey is the key the review import answers under.
func (ProductReviews) EnvelopeKey() string { return ProductKey }

// ProductReviewInput is one review to add; Name, Email and a Rating of 1 to 5 are required, and a second review under an Email the product has overwrites the first, keeping its id.
type ProductReviewInput struct {
	Name   string
	Email  string
	Rating *float64
	Review string
	// Comment carries the same text as Review, under the comment key.
	Comment string
	Status  string
	// DateCreated is yyyy-MM-ddTHH:mm:ss; the read's date key is refused on the write.
	DateCreated string `wire:"dateCreated"`
}

// ProductReviewsImportInput is the review import's body, {"product": {"reviews": [...]}}; Reviews must not be empty.
type ProductReviewsImportInput struct {
	Reviews       []*ProductReviewInput
	TransactionNo string
}

func (ProductReviewsImportInput) EnvelopeKey() string { return ProductKey }

// TopLevel names the field that sits beside that envelope rather than inside it.
func (ProductReviewsImportInput) TopLevel() []string { return []string{"transaction_no"} }

// ToMap renders the request body.
func (p ProductReviewsImportInput) ToMap() map[string]any { return dto.ToMap(p) }

// ProductStockVariation is one variant's stock on a product with variations.
type ProductStockVariation struct {
	VariationID    *int64
	ProductID      *int64
	Name           *string
	SKU            *string
	AvailableStock *int64
	Locations      []*ProductStockLocation
}

// ProductStock is the stock sub-path's body, read and written; AvailableStock is the sum of the rows' stock, and Variations is there only on a product with variations.
type ProductStock struct {
	Status               *string
	ProductID            *int64
	MultiLocationEnabled *bool
	TrackInventory       *bool
	AvailableStock       *int64
	Locations            []*ProductStockLocation
	Variations           []*ProductStockVariation
}

// ProductStockLocationInput is one row of the stock write; with no row marked is_default, the first row sent becomes the default.
type ProductStockLocationInput struct {
	Warehouse *int64
	Stock     *int64
	// LowStockLevel left nil is stored null.
	LowStockLevel *int64
	IsDefault     *bool
}

// ProductStockInput is the stock write's body: Locations is required and replaces the product's whole set.
type ProductStockInput struct {
	Note      string
	Locations []*ProductStockLocationInput
}

func (p ProductStockInput) ToMap() map[string]any { return dto.ToMap(p) }

// ProductCount is the count's answer; Count repeats Total.
type ProductCount struct {
	Total *int64
	Count *int64
}

// ProductSkuHolder is the product that holds a SKU; Status is trashed for one in the trash.
type ProductSkuHolder struct {
	ID     *int64
	Name   *string
	Status *string
}

// ProductSkuCheck is a SKU's verdict: Available and IsExist are opposites, Product and ID name the holder and are nil when the SKU is free. A variant's SKU names the variant's own product.
type ProductSkuCheck struct {
	Available *bool
	IsExist   *bool
	Product   *ProductSkuHolder
	ID        *int64
}

// ProductSpecial is one row of a special listing: a product's storefront prices.
type ProductSpecial struct {
	ID                      *int64
	SKU                     *string
	Name                    *string
	Image                   *string
	BasePrice               *float64
	IsOnSale                *bool
	OnSaleAmountType        *string
	EffectivePrice          *float64
	ActualBasePrice         *float64
	OnSaleAmount            *float64
	ExpectToPayPrice        *float64
	ExpectToPayPriceWithTax *float64
	TaxExclusiveBasePrice   *float64
	TaxExclusiveSalePrice   *float64
	Tax                     *float64
	PreviousPriceToDisplay  *float64
	ActualPriceToDisplay    *float64
	SalePrice               *float64
	Updated                 *string
	URL                     *string
	Summary                 *string
	Description             *string
}

// ProductSpecialPage is a special listing. Max is the page size applied; Status is success on a per-customer answer with no rows, and absent otherwise.
type ProductSpecialPage struct {
	Status     *string
	Total      *int64
	Max        *int64
	Offset     *int64
	Products   []*ProductSpecial
	Pagination *Pagination
}

// ProductsBulkDeleteInput is the bulk delete's body, no envelope: IDs is required and must not be empty.
type ProductsBulkDeleteInput struct {
	IDs []int64 `wire:"ids"`
}

// ToMap renders the request body.
func (p ProductsBulkDeleteInput) ToMap() map[string]any { return dto.ToMap(p) }

// ProductBulkDeleteError is one id the bulk delete did not trash; Reason is not_found for an unknown id.
type ProductBulkDeleteError struct {
	ID     *int64
	Reason *string
}

// ProductsBulkDeleteResult is the bulk delete's answer, 200 even when an id failed: counts, and the failed ids.
type ProductsBulkDeleteResult struct {
	Deleted *int64
	Failed  *int64
	Errors  []*ProductBulkDeleteError
}

// ProductReviewList is one product's reviews; the listing has no paging block.
type ProductReviewList struct {
	Reviews []*ProductReview
}

// ProductReviewPatchInput is the review moderation body, no envelope; Status is pending or approved.
type ProductReviewPatchInput struct {
	Status  string
	Rating  *float64
	Comment string
}

// ToMap renders the request body.
func (p ProductReviewPatchInput) ToMap() map[string]any { return dto.ToMap(p) }

// ProductImageFileInput is one image to store: Filename and its Base64 content.
type ProductImageFileInput struct {
	Filename string
	Base64   string `wire:"base64"`
}

// ProductImageAndVideoInput is the image_and_video block of an image upload.
type ProductImageAndVideoInput struct {
	Images []*ProductImageFileInput
}

// ProductImageUploadInput is the image upload's body, {"product": {"image_and_video": {"images": [...]}}}.
type ProductImageUploadInput struct {
	ImageAndVideo *ProductImageAndVideoInput
}

// EnvelopeKey is the key the upload's body is wrapped in.
func (ProductImageUploadInput) EnvelopeKey() string { return ProductKey }

// ToMap renders the request body.
func (p ProductImageUploadInput) ToMap() map[string]any { return dto.ToMap(p) }

// ProductMediaUploaded is what an image upload answers: Product holds the files just stored, TransactionNo is the product's own.
type ProductMediaUploaded struct {
	Product       *ProductImageAndVideo
	TransactionNo *string
}

// ProductMediaUploadInput is the media upload's body, no envelope: Type (only images is accepted), Filename and Content, base64 or a data: URI. Position is in the collection and refused unknown_field by the store.
type ProductMediaUploadInput struct {
	Type     string
	Filename string
	Content  string
	Position *int64
}

// ToMap renders the request body.
func (p ProductMediaUploadInput) ToMap() map[string]any { return dto.ToMap(p) }

// ProductInventoryStockInput is the stock figures a top-up writes.
type ProductInventoryStockInput struct {
	TrackInventory bool
	AvailableStock *int64
	LowLevelStock  *int64
}

// ProductInventoryInput is the top-up's body, {"product": {"product_pricing_and_stock": {...}}}.
type ProductInventoryInput struct {
	ProductPricingAndStock *ProductInventoryStockInput
}

// EnvelopeKey is the key the top-up's body is wrapped in.
func (ProductInventoryInput) EnvelopeKey() string { return ProductKey }

// ToMap renders the request body.
func (p ProductInventoryInput) ToMap() map[string]any { return dto.ToMap(p) }

// ProductInventoryResult is the top-up's answer. TransactionNo is the product's own from before the write, nil on its first top-up, though the header is stored on the product by then.
type ProductInventoryResult struct {
	Status        *string
	Message       *string
	TransactionNo *string
}

// ProductPageShape is the typed view of the listing body, {"products": […], "pagination": {…}}.
func ProductPageShape(decoded any) any { return dto.Read[ProductPage](decoded) }

// ProductShape is the typed view of a body shaped {"product": {…}}, field_metadata beside it.
func ProductShape(decoded any) any { return dto.Read[Product](decoded) }

// ProductReviewsShape is the typed view of the review import's answer, {"product": {"reviews": [...]}}.
func ProductReviewsShape(decoded any) any { return dto.Read[ProductReviews](decoded) }

// ProductStockShape is the typed view of the stock sub-path's body.
func ProductStockShape(decoded any) any { return dto.Read[ProductStock](decoded) }

// ProductCountShape is the typed view of the count's answer.
func ProductCountShape(decoded any) any { return dto.Read[ProductCount](decoded) }

// ProductSkuCheckShape is the typed view of the SKU check's answer.
func ProductSkuCheckShape(decoded any) any { return dto.Read[ProductSkuCheck](decoded) }

// ProductSpecialPageShape is the typed view of a special listing.
func ProductSpecialPageShape(decoded any) any { return dto.Read[ProductSpecialPage](decoded) }

// ProductsBulkDeleteResultShape is the typed view of the bulk delete's answer.
func ProductsBulkDeleteResultShape(decoded any) any {
	return dto.Read[ProductsBulkDeleteResult](decoded)
}

// ProductReviewListShape is the typed view of one product's review listing.
func ProductReviewListShape(decoded any) any { return dto.Read[ProductReviewList](decoded) }

// ProductReviewShape is the typed view of a body shaped {"review": {…}}.
func ProductReviewShape(decoded any) any { return dto.Read[ProductReview](decoded) }

// ProductMediaImageShape is the typed view of the media upload's answer, one flat image row.
func ProductMediaImageShape(decoded any) any { return dto.Read[ProductMediaImage](decoded) }

// ProductMediaUploadedShape is the typed view of an image or media upload's answer.
func ProductMediaUploadedShape(decoded any) any { return dto.Read[ProductMediaUploaded](decoded) }

// ProductInventoryResultShape is the typed view of the top-up's answer.
func ProductInventoryResultShape(decoded any) any { return dto.Read[ProductInventoryResult](decoded) }

// ProductMediaOrderInput is the media order body, no envelope: Type images and every image's id once, in the order wanted.
type ProductMediaOrderInput struct {
	Type string
	Ids  []int64
}

// ToMap renders the request body.
func (p ProductMediaOrderInput) ToMap() map[string]any { return dto.ToMap(p) }

// ProductMediaOrder is the media order's answer: the product's images in their new order.
type ProductMediaOrder struct {
	Images []*ProductMediaImage
}

// ProductMediaImagePatchInput is the image move body, no envelope: the 1-based Position, the only field it takes.
type ProductMediaImagePatchInput struct {
	Position *int64
}

// ToMap renders the request body.
func (p ProductMediaImagePatchInput) ToMap() map[string]any { return dto.ToMap(p) }

// ProductSeoAdvancedSettings is the robots and canonical URL settings.
type ProductSeoAdvancedSettings struct {
	Enabled      *bool
	RobotsMeta   *string
	CanonicalURL *string
}

// ProductSeoGoogleAnalyticsEvents is which Google Analytics events the product page sends.
type ProductSeoGoogleAnalyticsEvents struct {
	Enabled         *bool
	ViewItem        *bool
	ViewItemList    *bool
	ViewCart        *bool
	AddToWishList   *bool
	RemoveFromCart  *bool
	BeginCheckout   *bool
	AddShippingInfo *bool
	AddPaymentInfo  *bool
	Purchase        *bool
}

// ProductSeoGoogleAnalytics is Google Analytics on the product page.
type ProductSeoGoogleAnalytics struct {
	Enabled       *bool
	MeasurementID *string
	Events        *ProductSeoGoogleAnalyticsEvents
}

// ProductSeoFacebookPixelEvents is which Facebook Pixel events the product page sends.
type ProductSeoFacebookPixelEvents struct {
	Enabled          *bool
	Search           *bool
	AddToCart        *bool
	AddToWishList    *bool
	InitiateCheckout *bool
	AddPaymentInfo   *bool
	Purchase         *bool
}

// ProductSeoFacebookPixel is the Facebook Pixel on the product page.
type ProductSeoFacebookPixel struct {
	Enabled *bool
	PixelID *string
	Events  *ProductSeoFacebookPixelEvents
}

// ProductSeoGoogleAds is the Google Ads conversion the product page reports.
type ProductSeoGoogleAds struct {
	ConversionID    *string
	ConversionLabel *string
}

// ProductSeoAnalytics is the product page's analytics settings.
type ProductSeoAnalytics struct {
	Enabled         *bool
	Platform        *string
	GoogleAnalytics *ProductSeoGoogleAnalytics
	FacebookPixel   *ProductSeoFacebookPixel
	GoogleAds       *ProductSeoGoogleAds
}

// ProductSeoCode is one custom code block and whether it is on.
type ProductSeoCode struct {
	Enabled *bool
	Code    *string
}

// ProductSeoCustomCode is the custom code added to the product page's header and footer.
type ProductSeoCustomCode struct {
	Header *ProductSeoCode
	Footer *ProductSeoCode
}

// ProductSeoSocialCard is the title, description and image one network shows for the product.
type ProductSeoSocialCard struct {
	Title       *string
	Description *string
	ImageURL    *string
}

// ProductSeoSocial is the product's social card, and one per network.
type ProductSeoSocial struct {
	Title       *string
	Description *string
	ImageURL    *string
	Facebook    *ProductSeoSocialCard
	Twitter     *ProductSeoSocialCard
	Linkedin    *ProductSeoSocialCard
}

// ProductSeoAiVisibility is what the product offers AI search; each field is nil until set.
type ProductSeoAiVisibility struct {
	StructuredAttributes any
	ProductFaq           any
	CitabilityScore      any
	JSONLd               any
}

// ProductSeo is a product's SEO sub-resource, under seo.
type ProductSeo struct {
	URLHandle        *string
	Title            *string
	Description      *string
	AdvancedSettings *ProductSeoAdvancedSettings
	Analytics        *ProductSeoAnalytics
	CustomCode       *ProductSeoCustomCode
	Social           *ProductSeoSocial
	AiVisibility     *ProductSeoAiVisibility
}

// EnvelopeKey is the key the SEO read answers under.
func (ProductSeo) EnvelopeKey() string { return ProductSeoKey }

// ProductStockTracking is a product's stock figure and low-stock level, with whether tracking is on.
type ProductStockTracking struct {
	Enabled        *bool
	AvailableStock *int64
	LowStockLevel  *int64
}

// ProductStockSummary is the stock block of the stock patch's answer.
type ProductStockSummary struct {
	Tracking             *ProductStockTracking
	MultiLocationEnabled *bool
	Locations            []*ProductStockLocation
	Variations           []*ProductStockVariation
}

// ProductStockAdjustmentUser is the admin who made a stock change.
type ProductStockAdjustmentUser struct {
	ID   *int64
	Name *string
}

// ProductStockAdjustment is one stock change; Delta is the new figure minus the old.
type ProductStockAdjustment struct {
	ID            *int64
	Delta         *int64
	Note          *string
	TransactionNo *string
	CreatedBy     *ProductStockAdjustmentUser
	CreatedAt     *string
}

// ProductStockPatchInput is the stock patch's body, no envelope. Tracking merges; Locations sets the stock per location; Note and TransactionNo go on the adjustment.
type ProductStockPatchInput struct {
	Tracking      *ProductStockTracking
	Locations     []*ProductStockLocationInput
	Note          string
	TransactionNo string
}

// ToMap renders the request body.
func (p ProductStockPatchInput) ToMap() map[string]any { return dto.ToMap(p) }

// ProductStockPatchResult is the stock patch's answer: the stock now, and the adjustment it recorded.
type ProductStockPatchResult struct {
	Stock      *ProductStockSummary
	Adjustment *ProductStockAdjustment
}

// ProductStockAdjustmentPage is one page of a product's stock history.
type ProductStockAdjustmentPage struct {
	Adjustments []*ProductStockAdjustment
	Pagination  *Pagination
}

// ProductMediaOrderShape is the typed view of the media order's answer.
func ProductMediaOrderShape(decoded any) any { return dto.Read[ProductMediaOrder](decoded) }

// ProductSeoShape is the typed view of a body shaped {"seo": {…}}.
func ProductSeoShape(decoded any) any { return dto.Read[ProductSeo](decoded) }

// ProductStockPatchResultShape is the typed view of the stock patch's answer.
func ProductStockPatchResultShape(decoded any) any { return dto.Read[ProductStockPatchResult](decoded) }

// ProductStockAdjustmentPageShape is the typed view of a page of stock history.
func ProductStockAdjustmentPageShape(decoded any) any {
	return dto.Read[ProductStockAdjustmentPage](decoded)
}
