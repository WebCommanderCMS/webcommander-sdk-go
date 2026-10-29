package customergroups

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// CustomerGroupShippingProfile is the shipping profile a group names; a write sends the bare id, not this record.
type CustomerGroupShippingProfile struct {
	ID   *int64
	Name *string
}

// CustomerGroupBackgroundImage is the stored image; a write sends a base64 or data: URI string, and null removes it.
type CustomerGroupBackgroundImage struct {
	URL    *string
	Width  *int64
	Height *int64
}

// CustomerGroupAmountSetting is an enabled switch and its amount: minimum_purchase, max_purchase and credit_limit all take it.
type CustomerGroupAmountSetting struct {
	Enabled *bool
	Amount  *float64
}

// CustomerGroupPaymentTerms holds terms written NET_30, NET_60 or NET_90 and read back in lower case.
type CustomerGroupPaymentTerms struct {
	Enabled *bool
	Terms   *string
}

type CustomerGroupOnAccountPurchase struct {
	Enabled      *bool
	MaxPurchase  *CustomerGroupAmountSetting
	CreditLimit  *CustomerGroupAmountSetting
	PaymentTerms *CustomerGroupPaymentTerms
}

type CustomerGroupPurchasePreference struct {
	MinimumPurchase   *CustomerGroupAmountSetting
	OnAccountPurchase *CustomerGroupOnAccountPurchase
}

// CustomerGroupTaxConfiguration is read and written alike; BasePriceRounding is nearest, up or down.
type CustomerGroupTaxConfiguration struct {
	OverrideStoreTax             *bool
	TaxExempt                    *bool
	TaxMessage                   *string
	ShowPricesWithTaxIncluded    *bool
	PricesEnteredWithTaxIncluded *bool
	BasePriceRounding            *string
}

// CustomerGroupProductDisplay is read and written alike; PriceDisplayMethod is lowest or highest.
type CustomerGroupProductDisplay struct {
	OverwriteDefault             *bool
	VariationDisplayType         *string
	DisplayParentImage           *bool
	ActivateColourVariationImage *bool
	ShowSelectedVariationValue   *bool
	StrikeThroughOutOfStock      *bool
	SingleClickVariationLoading  *bool
	PriceDisplayMethod           *string
}

// CustomerGroupMember is one member by its customer_id: an entry of Customers and the member read.
type CustomerGroupMember struct {
	ID        *int64
	FirstName *string
	LastName  *string
	Email     *string
}

// EnvelopeKey is the key the member read wraps one member in.
func (CustomerGroupMember) EnvelopeKey() string { return CustomerGroupCustomerKey }

// CustomerGroup is one group; a listing row and a create, patch or assign answer stop at UpdatedAt, the row read and the member write add the last two.
type CustomerGroup struct {
	ID                           *int64
	Name                         *string
	Description                  *string
	Status                       *string
	DefaultTaxCode               *string
	AddAsNewCustomerDefaultGroup *bool
	ShippingProfile              *CustomerGroupShippingProfile
	BackgroundImage              *CustomerGroupBackgroundImage
	PurchasePreference           *CustomerGroupPurchasePreference
	TaxConfiguration             *CustomerGroupTaxConfiguration
	ProductDisplay               *CustomerGroupProductDisplay
	CreatedAt                    *string
	UpdatedAt                    *string
	CustomerCount                *int64
	Customers                    []*CustomerGroupMember
	// FieldMetadata is the schema a field_metadata=true row read sends beside the envelope.
	FieldMetadata any
}

// EnvelopeKey is the key one group arrives under and every group write takes.
func (CustomerGroup) EnvelopeKey() string { return CustomerGroupKey }

// Pagination is the listing's five-key pagination block; Limit is the page size applied.
type Pagination struct {
	Total       *int64
	Limit       *int64
	Offset      *int64
	HasPrevious *bool
	HasNext     *bool
}

// CustomerGroupPage is one page of the listing, with the field_metadata block when it was asked for.
type CustomerGroupPage struct {
	CustomerGroups []*CustomerGroup
	Pagination     *Pagination
	FieldMetadata  any
}

// CustomerGroupCount is the count sub-path's bare {"total": N}.
type CustomerGroupCount struct {
	Total *int64
}

// CustomerGroupAssignment is the assign answer; Added is false when the customer was already a member.
type CustomerGroupAssignment struct {
	Status     *string
	CustomerID *int64
	Email      *string
	Group      *CustomerGroup
	Added      *bool
}

// CustomerGroupInput is the writable form for the create and the patch; Status takes A, I, active or inactive.
type CustomerGroupInput struct {
	Name           string
	Description    string
	Status         string
	DefaultTaxCode string
	// AddAsNewCustomerDefaultGroup and ShippingProfile are pointers, so an unset one is omitted rather than sent false or 0.
	AddAsNewCustomerDefaultGroup *bool
	ShippingProfile              *int64
	BackgroundImage              string
	PurchasePreference           *CustomerGroupPurchasePreference
	TaxConfiguration             *CustomerGroupTaxConfiguration
	ProductDisplay               *CustomerGroupProductDisplay
}

// EnvelopeKey is the wrapper both writes require.
func (CustomerGroupInput) EnvelopeKey() string { return CustomerGroupKey }

// ToMap renders the request body.
func (c CustomerGroupInput) ToMap() map[string]any { return dto.ToMap(c) }

// CustomerGroupCustomersInput is the member write's body, replacing the whole list; an empty non-nil slice empties it.
type CustomerGroupCustomersInput struct {
	Customers []any
}

func (c CustomerGroupCustomersInput) ToMap() map[string]any { return dto.ToMap(c) }

// CustomerGroupAssignInput is the assign body, both ids bare: Customer is a customer_id, as a number or a string.
type CustomerGroupAssignInput struct {
	Customer any
	Group    any
}

func (c CustomerGroupAssignInput) ToMap() map[string]any { return dto.ToMap(c) }

// CustomerGroupShape is the typed view of a body shaped {"customer_group": {…}}.
func CustomerGroupShape(decoded any) any { return dto.Read[CustomerGroup](decoded) }

// CustomerGroupPageShape is the typed view of the listing body.
func CustomerGroupPageShape(decoded any) any { return dto.Read[CustomerGroupPage](decoded) }

// CustomerGroupCountShape is the typed view of the count body.
func CustomerGroupCountShape(decoded any) any { return dto.Read[CustomerGroupCount](decoded) }

// CustomerGroupMemberShape is the typed view of the member read.
func CustomerGroupMemberShape(decoded any) any { return dto.Read[CustomerGroupMember](decoded) }

// CustomerGroupAssignmentShape is the typed view of the assign answer.
func CustomerGroupAssignmentShape(decoded any) any { return dto.Read[CustomerGroupAssignment](decoded) }
