package pagesettings

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// CartPage is the cart_page section: what the cart page itself offers, three booleans.
type CartPage struct {
	ShowProductThumbnail *bool
	EnableClearCart      *bool
	SaveCart             *bool
}

// EnvelopeKey is the section's own name, the wrapper a write must carry.
func (CartPage) EnvelopeKey() string { return "cart_page" }

// ToMap renders the request body.
func (c CartPage) ToMap() map[string]any { return dto.ToMap(c) }

// FlyingCart is the flying_cart section: the slide-out cart; its other two switches require Enabled.
type FlyingCart struct {
	Enabled         *bool
	EnableClearCart *bool
	SaveCart        *bool
}

// EnvelopeKey is the section's own name, the wrapper a write must carry.
func (FlyingCart) EnvelopeKey() string { return "flying_cart" }

func (f FlyingCart) ToMap() map[string]any { return dto.ToMap(f) }

// ShoppingCartSectionMeta is one group_meta entry: what a section is, whether writable, where managed.
type ShoppingCartSectionMeta struct {
	Kind     *string
	Writable *bool
	Href     *string
}

// ShoppingCartGroupMeta is the group_meta block, one entry per section; only the facade sends it.
type ShoppingCartGroupMeta struct {
	CartPage   *ShoppingCartSectionMeta
	FlyingCart *ShoppingCartSectionMeta
}

// ShoppingCart is the facade's body or one section's, under the section's own key.
type ShoppingCart struct {
	CartPage   *CartPage
	FlyingCart *FlyingCart
	GroupMeta  *ShoppingCartGroupMeta
	// FieldMetadata is the raw field_metadata block, keyed by section on the facade and by field on a section.
	FieldMetadata any
}

func (ShoppingCart) EnvelopeKey() string { return "shopping_cart" }

func (s ShoppingCart) ToMap() map[string]any { return dto.ToMap(s) }

// ShoppingCartShape is the typed view of the facade's body and of one section's.
func ShoppingCartShape(decoded any) any { return dto.Read[ShoppingCart](decoded) }
