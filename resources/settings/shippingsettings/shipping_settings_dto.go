package shippingsettings

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// ProfileReference is what a reference field holds: an id and the name it resolves to, or null.
type ProfileReference struct {
	ID   *int64
	Name *string
}

// AreaMeta is one group_meta entry: what the area is, whether it can be written, and its path.
type AreaMeta struct {
	Kind     *string
	Writable *bool
	Href     *string
}

// GroupMeta is the group_meta block: the settings area and the two collections beside it.
type GroupMeta struct {
	ShippingSettings *AreaMeta
	ShippingMethods  *AreaMeta
	ShippingClasses  *AreaMeta
}

// ShippingSettings is the area's seven flat fields under shipping_settings; ToMap renders the wrapped body, and a reference field is written by id, or cleared to null, through a mapping body instead.
type ShippingSettings struct {
	DefaultShippingProfile          *ProfileReference
	ApplyStoreTaxSettingsToShipping *bool
	DefaultDeliveryProfile          *ProfileReference
	ApplyStoreTaxSettingsToDelivery *bool
	ShippingTaxProfile              *ProfileReference
	EnableShippingClass             *bool
	EnableClassInFlatView           *bool
	GroupMeta                       *GroupMeta
	// FieldMetadata is the raw field_metadata block.
	FieldMetadata any
}

func (ShippingSettings) EnvelopeKey() string { return "shipping_settings" }

// ToMap renders the request body.
func (s ShippingSettings) ToMap() map[string]any { return dto.ToMap(s) }

// ShippingSettingsShape is the typed view of the facade's body, the one shape this area sends.
func ShippingSettingsShape(decoded any) any { return dto.Read[ShippingSettings](decoded) }
