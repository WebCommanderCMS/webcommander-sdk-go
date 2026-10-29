package pricedisplaysettings

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// DisplayScopes is the surfaces one setting applies to; a scope outside the block's allowed_scopes is 400.
type DisplayScopes struct {
	ShoppingCart *bool
	Checkout     *bool
	OrderSuccess *bool
}

// ScopedToggle is a master switch with optional per-surface control; with Customize off, Scopes follow Enabled.
type ScopedToggle struct {
	Enabled   *bool
	Customize *bool
	Scopes    *DisplayScopes
}

// PlainToggle is a setting with nothing but a master switch.
type PlainToggle struct {
	Enabled *bool
}

// Subtotal is the subtotal block: where the line shows, and whether it carries tax and discount.
type Subtotal struct {
	Show            *DisplayScopes
	IncludeTax      *ScopedToggle
	IncludeDiscount *ScopedToggle
}

// Discount is the discount block: where the line shows, whether it hides at zero, and its tax.
type Discount struct {
	Show       *DisplayScopes
	HideIfZero *ScopedToggle
	IncludeTax *ScopedToggle
}

// ShippingCost is the shipping_cost block: where the line shows, and its tax and discount.
type ShippingCost struct {
	Show            *DisplayScopes
	IncludeTax      *ScopedToggle
	IncludeDiscount *ScopedToggle
}

// ShippingDiscount is the shipping_discount block: where the line shows, and the shipping tax switches.
type ShippingDiscount struct {
	Show        *DisplayScopes
	IncludeTax  *ScopedToggle
	ShippingTax *PlainToggle
}

// Handling is the handling block: whether handling shows apart from shipping, and its two tax settings.
type Handling struct {
	ShowSeparately *DisplayScopes
	IncludeTax     *ScopedToggle
	HandlingTax    *ScopedToggle
}

// PricingInfo is the pricing_info block: a note after the total, and hiding the due amount on success.
type PricingInfo struct {
	Enabled            *bool
	Text               *string
	Customize          *bool
	Scopes             *DisplayScopes
	DoNotShowDueAmount *bool
}

// OrderItem is the order_item section: how a line item's own price shows, in five ScopedToggle blocks.
type OrderItem struct {
	PriceIncludingTax      *ScopedToggle
	PriceIncludingDiscount *ScopedToggle
	TaxInLine              *ScopedToggle
	DiscountInLineItem     *ScopedToggle
	ItemTotalIncludingTax  *ScopedToggle
}

// EnvelopeKey is the section's own name, one of a write's two forms; the other is no wrapper at all.
func (OrderItem) EnvelopeKey() string { return "order_item" }

// ToMap renders the request body.
func (o OrderItem) ToMap() map[string]any { return dto.ToMap(o) }

// OrderSummary is the order_summary section: which lines the summary panel shows, and what each includes.
type OrderSummary struct {
	OrderTotalExcludingTax      *ScopedToggle
	OrderTotalIncludingDiscount *ScopedToggle
	Subtotal                    *Subtotal
	Discount                    *Discount
	SubtotalTax                 *PlainToggle
	TotalTax                    *PlainToggle
	OrderTotal                  *PlainToggle
	ShippingCost                *ShippingCost
	ShippingDiscount            *ShippingDiscount
	Handling                    *Handling
	PaymentSurcharge            *ScopedToggle
	PricingInfo                 *PricingInfo
}

// EnvelopeKey is the section's own name, one of a write's two forms; the other is no wrapper at all.
func (OrderSummary) EnvelopeKey() string { return "order_summary" }

func (o OrderSummary) ToMap() map[string]any { return dto.ToMap(o) }

// SectionMeta is one group_meta entry: what a section is, whether it is writable, where it is managed.
type SectionMeta struct {
	Kind     *string
	Writable *bool
	Href     *string
}

// GroupMeta is the group_meta block, one SectionMeta per section; only the facade sends it.
type GroupMeta struct {
	OrderItem    *SectionMeta
	OrderSummary *SectionMeta
}

// PriceDisplaySettings is the facade's body or one section's; ToMap renders the facade PATCH, not a section's.
type PriceDisplaySettings struct {
	OrderItem    *OrderItem
	OrderSummary *OrderSummary
	GroupMeta    *GroupMeta
	// FieldMetadata is the raw field_metadata block, keyed by section on the facade and on a section alike.
	FieldMetadata any
}

func (PriceDisplaySettings) EnvelopeKey() string { return "price_display_settings" }

func (p PriceDisplaySettings) ToMap() map[string]any { return dto.ToMap(p) }

// PriceDisplaySettingsShape is the typed view of the facade's body and of one section's.
func PriceDisplaySettingsShape(decoded any) any { return dto.Read[PriceDisplaySettings](decoded) }
