package discountprofiles

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// AmountDetails is the block an amount profile carries.
type AmountDetails struct {
	Type             *string
	ApplyTo          *string
	SingleAmount     *float64
	SingleAmountType *string
	MinimumAmountOn  *string
	// Tiers reads as tier ids.
	Tiers []*int64
}

// ShippingDetails is the block a shipping profile carries.
type ShippingDetails struct {
	Type             *string
	ApplyTo          *string
	MinimumQtyOn     *string
	AmountType       *string
	SingleAmountType *string
	MinimumQty       *int64
	MinimumAmount    *float64
	SingleAmount     any
	CapAmount        any
	MaximumTime      any
	Zone             any
	ShippingClass    any
	Tiers            []any
}

// ProductDetails is the block a product profile carries.
type ProductDetails struct {
	Type                 *string
	FreeProductMinQtyOn  *string
	FreeProductMaxQty    *int64
	FreeProductMinAmount *float64
	MinimumQtyOn         *string
	AmountType           *string
	SingleAmountType     *string
	SingleAmount         any
	CapPrice             any
	CapPriceMaxQty       any
	ProductIds           []any
	CategoryIds          []any
	Tiers                []any
}

// CustomerCoupon is one generated coupon; CustomerEmail is "" until a customer is assigned.
type CustomerCoupon struct {
	ID            *int64
	Code          *string
	CustomerEmail *string
}

// DiscountProfile is one profile, a listing row or a row path's {"discount_profile": {…}}, with the one details block its discount_details_type names.
type DiscountProfile struct {
	ID                                     *int64
	UUID                                   *string
	Name                                   *string
	Type                                   *string
	DiscountDetailsType                    *string
	AmountDetails                          *AmountDetails
	ShippingDetails                        *ShippingDetails
	ProductDetails                         *ProductDetails
	IsActive                               *bool
	StartDate                              *string
	IsSpecifyEndDate                       *bool
	EndDate                                *string
	IsApplyCouponCode                      *bool
	DefaultCouponCode                      *string
	IsCouponCodeAutoGenerate               *bool
	IsCreateUniqueCouponEachCustomer       *bool
	IsImportedCoupon                       *bool
	IsMaximumUseTotal                      *bool
	MaximumUseCount                        *int64
	IsMaximumUseCustomer                   *bool
	MaximumUseCustomerCount                *int64
	IsMaximumDiscountAllowed               *bool
	MaximumDiscountAllowedAmount           *float64
	IsExcludeProductsOnSale                *bool
	IsDiscountUsedWithOtherDiscount        *bool
	IsDisplayDiscountInformationProdDetail *bool
	IsDisplayTextCoupon                    *bool
	DisplayTextCoupon                      *string
	IsDisplayTextCart                      *bool
	DisplayTextCart                        *string
	IsDisplayTextPartialDiscountCondition  *bool
	DisplayTextPartialDiscountCondition    *string
	InvoiceNote                            *string
	ExcludeProducts                        []*int64
	CustomerCoupons                        []*CustomerCoupon
	Usage                                  []*int64
	CreatedOn                              *string
	LastUpdatedOn                          *string
	// FieldMetadata is the schema a field_metadata=true row read sends beside the envelope.
	FieldMetadata any
}

// EnvelopeKey is the key a row path and every write's answer wrap the record in; no write takes it.
func (DiscountProfile) EnvelopeKey() string { return DiscountProfileKey }

// DiscountProfileInput is the writable form, with no envelope; the update replaces, resetting every flag, date and limit it omits.
type DiscountProfileInput struct {
	Name                string
	Type                string
	DiscountDetailsType string
	AmountDetails       *AmountDetails
	ShippingDetails     *ShippingDetails
	ProductDetails      *ProductDetails
	// The flags and the limits are pointers, so an unset one is omitted rather than sent false or 0.
	IsActive                               *bool
	StartDate                              string
	IsSpecifyEndDate                       *bool
	EndDate                                string
	IsApplyCouponCode                      *bool
	DefaultCouponCode                      string
	IsCouponCodeAutoGenerate               *bool
	IsCreateUniqueCouponEachCustomer       *bool
	IsImportedCoupon                       *bool
	IsMaximumUseTotal                      *bool
	MaximumUseCount                        *int64
	IsMaximumUseCustomer                   *bool
	MaximumUseCustomerCount                *int64
	IsMaximumDiscountAllowed               *bool
	MaximumDiscountAllowedAmount           *float64
	IsExcludeProductsOnSale                *bool
	IsDiscountUsedWithOtherDiscount        *bool
	IsDisplayDiscountInformationProdDetail *bool
	IsDisplayTextCoupon                    *bool
	DisplayTextCoupon                      string
	IsDisplayTextCart                      *bool
	DisplayTextCart                        string
	IsDisplayTextPartialDiscountCondition  *bool
	DisplayTextPartialDiscountCondition    string
	InvoiceNote                            string
	ExcludeProducts                        []int64
}

// ToMap renders the request body.
func (d DiscountProfileInput) ToMap() map[string]any { return dto.ToMap(d) }

// Pagination is the listing's ten-key pagination block; Limit is the size applied, capped at 20.
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

// ResourceMeta is the one group_meta entry: what the collection is, whether it is writable and where.
type ResourceMeta struct {
	Kind     *string
	Writable *bool
	Href     *string
}

// GroupMeta is the listing's group_meta block.
type GroupMeta struct {
	DiscountProfiles *ResourceMeta
}

// DiscountProfilePage is one page of the listing: the rows, where they sit, the group_meta and, when asked for, the schema.
type DiscountProfilePage struct {
	DiscountProfiles []*DiscountProfile
	Pagination       *Pagination
	GroupMeta        *GroupMeta
	FieldMetadata    any
}

// DiscountProfileCoupons is the create-coupons answer: the codes generated, beside status and message.
type DiscountProfileCoupons struct {
	Status  *string
	Message *string
	Codes   []*string
}

// DiscountProfileCouponsInput is the create-coupons body; RequestedCodes is required and at least 1.
type DiscountProfileCouponsInput struct {
	RequestedCodes int64
}

func (d DiscountProfileCouponsInput) ToMap() map[string]any { return dto.ToMap(d) }

// DiscountCouponCustomerInput is the assign-customer body: a generated code and the email of an existing customer.
type DiscountCouponCustomerInput struct {
	Code  string
	Email string
}

func (d DiscountCouponCustomerInput) ToMap() map[string]any { return dto.ToMap(d) }

// DiscountProfileShape is the typed view of a body shaped {"discount_profile": {…}}.
func DiscountProfileShape(decoded any) any { return dto.Read[DiscountProfile](decoded) }

// DiscountProfilePageShape is the typed view of the listing body.
func DiscountProfilePageShape(decoded any) any { return dto.Read[DiscountProfilePage](decoded) }

// DiscountProfileCouponsShape is the typed view of the create-coupons body.
func DiscountProfileCouponsShape(decoded any) any { return dto.Read[DiscountProfileCoupons](decoded) }
