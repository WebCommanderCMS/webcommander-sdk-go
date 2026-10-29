package shippingclasses

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// ShippingClass is one shipping class, the same four fields on a listing row and on a row path.
type ShippingClass struct {
	ID            *int64
	Name          *string
	Description   *string
	IsDefault     *bool
	FieldMetadata any
}

// EnvelopeKey is the key a row path and every write's answer wrap the record in.
func (ShippingClass) EnvelopeKey() string { return ShippingClassKey }

// ShippingClassInput is the body both writes take; Name is required and unique.
type ShippingClassInput struct {
	Name        string
	Description string
}

// EnvelopeKey is the wrapper both writes require.
func (ShippingClassInput) EnvelopeKey() string { return ShippingClassKey }

// ToMap renders the request body.
func (s ShippingClassInput) ToMap() map[string]any { return dto.ToMap(s) }

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

// ShippingClassPage is one page of the listing, with the field_metadata block when it was asked for.
type ShippingClassPage struct {
	ShippingClasses []*ShippingClass
	Pagination      *Pagination
	FieldMetadata   any
}

// ShippingClassCount is the count sub-path's {"shipping_classes": {"count": N}}.
type ShippingClassCount struct {
	Count *int64
}

// EnvelopeKey is the plural key the count record arrives under.
func (ShippingClassCount) EnvelopeKey() string { return ShippingClassesKey }

// ShippingClassShape is the typed view of a body shaped {"shipping_class": {…}}.
func ShippingClassShape(decoded any) any { return dto.Read[ShippingClass](decoded) }

// ShippingClassPageShape is the typed view of the listing body.
func ShippingClassPageShape(decoded any) any { return dto.Read[ShippingClassPage](decoded) }

// ShippingClassCountShape is the typed view of the count sub-path's body.
func ShippingClassCountShape(decoded any) any { return dto.Read[ShippingClassCount](decoded) }
