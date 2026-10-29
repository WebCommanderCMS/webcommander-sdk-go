package ecommercesettings

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/inventorylocations"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/orderflags"
)

// ContinueShopping is the continue_shopping section; its target is a page slug, else 400 invalid_reference.
type ContinueShopping struct {
	EnableContinueShopping          *bool
	ContinueShoppingTarget          *string
	ContinueShoppingSpecifiedTarget *string
}

// EnvelopeKey is the section's own name, also one of the two wrappers a write takes.
func (ContinueShopping) EnvelopeKey() string { return "continue_shopping" }

// ToMap renders the request body.
func (c ContinueShopping) ToMap() map[string]any { return dto.ToMap(c) }

// ProductInventory is the product_inventory section: out-of-stock behaviour and when a sale draws stock.
type ProductInventory struct {
	ShowOutOfStockProducts            *bool
	LowStockNotification              *bool
	OutOfStockNotification            *bool
	HideOnSaleTagForOutOfStockProduct *bool
	OrderQuantityOverStock            *string
	UpdateStock                       *string
}

// EnvelopeKey is the section's own name, also one of the two wrappers a write takes.
func (ProductInventory) EnvelopeKey() string { return "product_inventory" }

func (p ProductInventory) ToMap() map[string]any { return dto.ToMap(p) }

// Pricing is the pricing section; each depends_on switch is enforced, so an amount of 0 while on is 400.
type Pricing struct {
	UpdatePriceRange                  *bool
	EnableMinimumPurchaseAmount       *bool
	MinimumPurchaseAmount             *float64
	EnableMinimumPurchaseAmountInCart *bool
	MinimumPurchaseAmountInCart       *float64
	EnableOnAccountPurchase           *bool
	AssignCreditLimitPerOrder         *bool
	CreditLimitPerOrder               *float64
	EnablePaymentTerms                *bool
	PaymentTerms                      *string
	EnableCustomerPriceRestriction    *bool
}

// EnvelopeKey is the section's own name, also one of the two wrappers a write takes.
func (Pricing) EnvelopeKey() string { return "pricing" }

func (p Pricing) ToMap() map[string]any { return dto.ToMap(p) }

// ProductRestriction is the product_restriction section: one switch.
type ProductRestriction struct {
	EnableProductRestrictionBasedOnParentCategory *bool
}

// EnvelopeKey is the section's own name, also one of the two wrappers a write takes.
func (ProductRestriction) EnvelopeKey() string { return "product_restriction" }

func (p ProductRestriction) ToMap() map[string]any { return dto.ToMap(p) }

// OrderFilters is the order_filters section: the status each admin order list opens on, null when unset.
type OrderFilters struct {
	DefaultOrderStatus    *string
	DefaultPaymentStatus  *string
	DefaultShipmentStatus *string
}

// EnvelopeKey is the section's own name, also one of the two wrappers a write takes.
func (OrderFilters) EnvelopeKey() string { return "order_filters" }

func (o OrderFilters) ToMap() map[string]any { return dto.ToMap(o) }

// InventoryLocation is one row of location_settings.locations, the inventory location collection's type.
type InventoryLocation = inventorylocations.InventoryLocation

// LocationSettings is the location_settings section; EnableMultiLocation is the only writable field.
type LocationSettings struct {
	EnableMultiLocation *bool
	Locations           []*InventoryLocation
}

// EnvelopeKey is the section's own name, also one of the two wrappers a write takes.
func (LocationSettings) EnvelopeKey() string { return "location_settings" }

func (l LocationSettings) ToMap() map[string]any { return dto.ToMap(l) }

// OrderFlag is one row of order_flags, the order flag collection's own type.
type OrderFlag = orderflags.OrderFlag

// SectionMeta is one group_meta entry; location_settings' adds LocationsHref.
type SectionMeta struct {
	Kind          *string
	Writable      *bool
	Href          *string
	LocationsHref *string
}

// GroupMeta is the group_meta block, one entry per area key; only the facade sends it.
type GroupMeta struct {
	ContinueShopping   *SectionMeta
	ProductInventory   *SectionMeta
	Pricing            *SectionMeta
	ProductRestriction *SectionMeta
	OrderFilters       *SectionMeta
	LocationSettings   *SectionMeta
	OrderFlags         *SectionMeta
}

// EcommerceSettings is the facade's body or one section's; ToMap renders the facade PATCH, not a section's.
type EcommerceSettings struct {
	ContinueShopping   *ContinueShopping
	ProductInventory   *ProductInventory
	Pricing            *Pricing
	ProductRestriction *ProductRestriction
	OrderFilters       *OrderFilters
	LocationSettings   *LocationSettings
	OrderFlags         []*OrderFlag
	GroupMeta          *GroupMeta
	// FieldMetadata is the raw field_metadata block.
	FieldMetadata any
}

func (EcommerceSettings) EnvelopeKey() string { return "ecommerce_settings" }

func (e EcommerceSettings) ToMap() map[string]any { return dto.ToMap(e) }

// EcommerceSettingsShape is the typed view of the facade's body and of one section's.
func EcommerceSettingsShape(decoded any) any { return dto.Read[EcommerceSettings](decoded) }
