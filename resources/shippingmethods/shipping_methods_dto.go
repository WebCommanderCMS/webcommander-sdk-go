package shippingmethods

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// ShippingMethod is one shipping method, the same three fields on a listing row and a row path.
type ShippingMethod struct {
	ID            *int64
	Name          *string
	IsDefault     *bool
	FieldMetadata any
}

// EnvelopeKey is the key a row path and every write's answer wrap the record in; a listing row has none.
func (ShippingMethod) EnvelopeKey() string { return ShippingMethodKey }

// ShippingMethodInput is the writable form, shared by both writes: Name, required, unique, at most 100 characters.
type ShippingMethodInput struct {
	Name string
}

// EnvelopeKey is the wrapper both writes require.
func (ShippingMethodInput) EnvelopeKey() string { return ShippingMethodKey }

// ToMap renders the request body.
func (s ShippingMethodInput) ToMap() map[string]any { return dto.ToMap(s) }

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

// ShippingMethodPage is one page of the listing: rows, pagination, and field_metadata when asked for.
type ShippingMethodPage struct {
	ShippingMethods []*ShippingMethod
	// Pagination is nil when the body carried none.
	Pagination *Pagination
	// FieldMetadata is nil unless the call asked for it.
	FieldMetadata any
}

// ShippingMethodCount is the count sub-path's body, {"shipping_methods": {"count": N}}.
type ShippingMethodCount struct {
	Count *int64
}

// EnvelopeKey is the plural key the count record arrives under.
func (ShippingMethodCount) EnvelopeKey() string { return ShippingMethodsKey }

// ShippingMethodShape is the typed view of a body shaped {"shipping_method": {…}}.
func ShippingMethodShape(decoded any) any { return dto.Read[ShippingMethod](decoded) }

// ShippingMethodPageShape is the typed view of the listing body.
func ShippingMethodPageShape(decoded any) any { return dto.Read[ShippingMethodPage](decoded) }

// ShippingMethodCountShape is the typed view of the count sub-path's body.
func ShippingMethodCountShape(decoded any) any { return dto.Read[ShippingMethodCount](decoded) }
