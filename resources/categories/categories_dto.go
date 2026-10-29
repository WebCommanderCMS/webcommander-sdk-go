package categories

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// CategoryReference is another record a category points at, read as {id, name}; a write needs only ID.
type CategoryReference struct {
	ID   *int64
	Name *string
}

// CategoryGeneral is the general section as read.
type CategoryGeneral struct {
	Name           *string
	Title          *string
	SKU            *string
	URL            *string
	Heading        *string
	ParentCategory *CategoryReference
	FilterProfile  any
	Layout         *CategoryReference
	Summary        *string
}

// CategoryDescription is the description section, the same one field read and written.
type CategoryDescription struct {
	Content *string
}

// CategoryImagesMedia is the images_media section as read: the two stored images' URLs and the image's alt text, each "" when unset.
type CategoryImagesMedia struct {
	ImageURL           *string
	ImageAlt           *string
	BackgroundImageURL *string
}

// CategorySelectedCustomers is who available_for "selected" admits, each customer and group read as {id, name}.
type CategorySelectedCustomers struct {
	Customers []*CategoryReference
	Groups    []*CategoryReference
}

// CategoryAvailabilityVisibility is the availability_visibility section as read; password is never echoed.
type CategoryAvailabilityVisibility struct {
	Available            *bool
	Visible              *bool
	AvailableFor         *string
	SelectedCustomers    *CategorySelectedCustomers
	AvailableOnDateRange *bool
	AvailableFromDate    *string
	AvailableToDate      *string
	PasswordProtected    *bool
	CustomSorting        *bool
	DefaultSorting       *string
	DisableTracking      *bool
}

// CategoryTaxShipping is the tax_shipping section, read and written; enable_shipping false clears ShippingProfile.
type CategoryTaxShipping struct {
	TaxExempt            *bool
	TaxProfile           *CategoryReference
	TaxMessage           *string
	ShowPriceWithTax     *bool
	PricesEnteredWithTax *bool
	BasePriceRounding    *string
	EnableShipping       *bool
	ShippingProfile      *CategoryReference
}

// CategoryProduct is one assigned product, as assign_products reads it.
type CategoryProduct struct {
	ID   *int64
	Name *string
	SKU  *string
}

// CategoryAssignProducts is the assign_products section as read.
type CategoryAssignProducts struct {
	Products []*CategoryProduct
}

// SectionMeta is one group_meta entry; Href is the section's path with the row's id left a placeholder.
type SectionMeta struct {
	Kind     *string
	Writable *bool
	Href     *string
}

// GroupMeta is the block the listing, a row read and every whole-row write send beside the row.
type GroupMeta struct {
	General                *SectionMeta
	Description            *SectionMeta
	ImagesMedia            *SectionMeta
	AvailabilityVisibility *SectionMeta
	TaxShipping            *SectionMeta
	AssignProducts         *SectionMeta
}

// Category is one category: a listing row, a row read, a whole-row write's answer, or a section's answer, which sets that section alone.
type Category struct {
	General                *CategoryGeneral
	Description            *CategoryDescription
	ImagesMedia            *CategoryImagesMedia
	AvailabilityVisibility *CategoryAvailabilityVisibility
	TaxShipping            *CategoryTaxShipping
	AssignProducts         *CategoryAssignProducts
	ID                     *int64
	TransactionNo          *string
	Disposable             *bool
	InTrash                *bool
	ParentInTrash          *bool
	CreatedBy              *string
	CreatedOn              *string
	LastUpdatedBy          *string
	LastUpdatedOn          *string
	GroupMeta              *GroupMeta
	FieldMetadata          any
}

// EnvelopeKey is the key one row, a section's answer and every write wrap the record in.
func (Category) EnvelopeKey() string { return CategoryKey }

// Pagination is the listing's ten-key pagination block; Limit is the size the store applied.
type Pagination struct {
	Total        *int64
	Limit        *int64
	Offset       *int64
	Count        *int64
	CurrentPage  *int64
	TotalPages   *int64
	HasNext      *bool
	HasPrevious  *bool
	PreviousPage *string
	NextPage     *string
}

// CategoryPage is one page of the listing, with the field_metadata block when it was asked for.
type CategoryPage struct {
	Categories    []*Category
	GroupMeta     *GroupMeta
	Pagination    *Pagination
	FieldMetadata any
}

// CategoryCount is the count sub-path's {"categories": {"count": N}}.
type CategoryCount struct {
	Count *int64
}

// EnvelopeKey is the plural key the count record arrives under.
func (CategoryCount) EnvelopeKey() string { return CategoriesKey }

// CategoryGeneralInput is the writable general fields; Name is required on a create and kept by a PUT that omits it.
type CategoryGeneralInput struct {
	Name           string
	Title          string
	SKU            string
	URL            string
	Heading        string
	ParentCategory *CategoryReference
	Layout         *CategoryReference
	Summary        string
}

// CategoryImage is one image upload; a data: URI string in its place is 500.
type CategoryImage struct {
	Filename string
	Base64   string `wire:"base64"`
}

// CategoryImagesMediaInput is the writable images_media fields; neither null nor "" removes a stored image.
type CategoryImagesMediaInput struct {
	Image           *CategoryImage
	ImageAlt        string
	BackgroundImage *CategoryImage
}

// CategorySelectedCustomersInput is the ids available_for "selected" admits; an id nothing holds is dropped unrefused.
type CategorySelectedCustomersInput struct {
	Customers []int64
	Groups    []int64
}

// CategoryAvailabilityVisibilityInput is the read fields plus the write-only Password, SelectedCustomers taking ids.
type CategoryAvailabilityVisibilityInput struct {
	Available            *bool
	Visible              *bool
	AvailableFor         string
	SelectedCustomers    *CategorySelectedCustomersInput
	AvailableOnDateRange *bool
	AvailableFromDate    string
	AvailableToDate      string
	PasswordProtected    *bool
	Password             string
	CustomSorting        *bool
	DefaultSorting       string
	DisableTracking      *bool
}

// CategoryAssignProductsInput is the product ids to assign; a write replaces the whole assignment.
type CategoryAssignProductsInput struct {
	Products []int64
}

// CategoryInput is the body every write takes, {"category": {<section>: {…}}}; a section PATCH reads only its own section.
type CategoryInput struct {
	General                *CategoryGeneralInput
	Description            *CategoryDescription
	ImagesMedia            *CategoryImagesMediaInput
	AvailabilityVisibility *CategoryAvailabilityVisibilityInput
	TaxShipping            *CategoryTaxShipping
	AssignProducts         *CategoryAssignProductsInput
	// TransactionNo is read inside the category wrapper only; beside it, it is ignored and nothing is filed.
	TransactionNo string
}

// EnvelopeKey is the wrapper every write requires.
func (CategoryInput) EnvelopeKey() string { return CategoryKey }

// ToMap renders the request body.
func (c CategoryInput) ToMap() map[string]any { return dto.ToMap(c) }

// CategoryShape is the typed view of a body shaped {"category": {…}}.
func CategoryShape(decoded any) any { return dto.Read[Category](decoded) }

// CategoryPageShape is the typed view of the listing body.
func CategoryPageShape(decoded any) any { return dto.Read[CategoryPage](decoded) }

// CategoryCountShape is the typed view of the count sub-path's body.
func CategoryCountShape(decoded any) any { return dto.Read[CategoryCount](decoded) }
