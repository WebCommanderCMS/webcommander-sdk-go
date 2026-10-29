package ordercustomfields

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// OrderCustomField is one field: a listing row carries the twelve fields up to MaxLength, the row path all twenty-one.
type OrderCustomField struct {
	ID                      *int64
	FieldType               *string
	Label                   *string
	Name                    *string
	IsActive                *bool
	OrderNo                 *int64
	CreatedOn               *string
	LastUpdatedOn           *string
	IsRequired              *bool
	InputFormat             *string
	HasMaxLength            *bool
	MaxLength               *int64
	HoverText               *string
	Placeholder             *string
	DefaultValue            *string
	HTMLClass               *string
	Options                 []*string
	IsAllowedForShipping    *bool
	IsAllowedForDelivery    *bool
	IsAllowedForStorePickup *bool
	IsAllowedForPreOrder    *bool
	// FieldMetadata is the raw field_metadata block.
	FieldMetadata any
}

// EnvelopeKey is the key the row path wraps the record in; a listing row carries none.
func (OrderCustomField) EnvelopeKey() string { return OrderCustomFieldKey }

// Pagination is the listing's ten-key pagination block; Limit is the size the store applied, capped at 20.
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

// ResourceMeta is the listing's one group_meta entry: Kind "resource", Writable false.
type ResourceMeta struct {
	Kind     *string
	Writable *bool
	Href     *string
}

// GroupMeta is the listing's group_meta block, keyed by the collection's own name; no other path sends it.
type GroupMeta struct {
	OrderCustomFields *ResourceMeta
}

// OrderCustomFieldPage is one page of the listing: the rows, where they sit, group_meta and any schema.
type OrderCustomFieldPage struct {
	OrderCustomFields []*OrderCustomField
	Pagination        *Pagination
	GroupMeta         *GroupMeta
	// FieldMetadata is nil unless the call asked for it.
	FieldMetadata any
}

// OrderCustomFieldTimeslotInput is the validate_timeslot body; SelectedDate is required, as dd/MM/yyyy.
type OrderCustomFieldTimeslotInput struct {
	SelectedDate string
	FieldName    string
	// Pointers, so an unset number is omitted rather than sent 0: an order_limit of 0 is a fully booked slot.
	WeekdayID  *int64
	StartHour  *int64
	EndHour    *int64
	OrderLimit *int64
}

// ToMap renders the request body.
func (o OrderCustomFieldTimeslotInput) ToMap() map[string]any { return dto.ToMap(o) }

// OrderCustomFieldShape is the typed view of a body shaped {"order_custom_field": {…}}.
func OrderCustomFieldShape(decoded any) any { return dto.Read[OrderCustomField](decoded) }

// OrderCustomFieldPageShape is the typed view of the listing body.
func OrderCustomFieldPageShape(decoded any) any { return dto.Read[OrderCustomFieldPage](decoded) }
