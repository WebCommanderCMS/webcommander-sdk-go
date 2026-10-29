package carts

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// CartTotals is the cart's totals; Subtotal sums each line's unit_price times its quantity.
type CartTotals struct {
	Subtotal         *float64
	Discount         *float64
	ShippingDiscount *float64
	Tax              *float64
	Total            *float64
	GrandTotal       *float64
	Paid             *float64
}

// CartShipping is the cart's shipping block; its two taxes follow the shipping address's country.
type CartShipping struct {
	Cost        *float64
	HandlingFee *float64
	HandlingTax *float64
	ShippingTax *float64
	TotalTax    *float64
}

// CartItemPricing is one line's pricing.
type CartItemPricing struct {
	UnitPrice        *float64
	DisplayUnitPrice *float64
	Total            *float64
	Discount         *float64
	DiscountTax      *float64
	Tax              *float64
}

// CartItemProduct is the product one line holds; AvailableStock is the product's, not the line's.
type CartItemProduct struct {
	ID             *int64
	SKU            *string
	BasePrice      *float64
	AvailableStock *int64
	OnSale         *bool
	CostPrice      *float64
	SalePrice      *float64
	URL            *string
	Name           *string
	ProductType    *string
	ImageLink      *string
}

// CartItem is one line of the cart; ID is the cart_item_id the line's writes take.
type CartItem struct {
	ID         *int64
	Quantity   *int64
	Variations any
	Taxable    *bool
	Shippable  *bool
	Pricing    *CartItemPricing
	Product    *CartItemProduct
}

// CartDiscount is the cart's discount block.
type CartDiscount struct {
	Discount              *float64
	OrderDiscount         *float64
	ShippingDiscount      *float64
	DiscountOnTax         *float64
	DiscountOnShippingTax *float64
}

// CartDiscountMessage is the discount entry of messages.
type CartDiscountMessage struct {
	Message *string
	Status  any
}

// CartMessages is the cart's messages block.
type CartMessages struct {
	Discount *CartDiscountMessage
}

// SectionMeta is the one entry of group_meta.
type SectionMeta struct {
	Kind     *string
	Writable *bool
	Href     *string
}

// GroupMeta is the group_meta block, which only a read sends.
type GroupMeta struct {
	Cart *SectionMeta
}

// Cart is one cart: the answer of the add, the read and the line update.
type Cart struct {
	CustomerID     *int64
	SessionID      *string
	CurrencyCode   *string
	CurrencySymbol *string
	DeliveryType   *string
	Totals         *CartTotals
	Shipping       *CartShipping
	Items          []*CartItem
	Discount       *CartDiscount
	Messages       *CartMessages
	// GroupMeta is read from beside the envelope.
	GroupMeta *GroupMeta
}

// EnvelopeKey is the key a cart arrives under.
func (Cart) EnvelopeKey() string { return CartKey }

// CartAddressInput is a shipping_address or billing_address; the cart answers neither back.
type CartAddressInput struct {
	FirstName    string
	LastName     string
	AddressLine1 string `wire:"address_line1"`
	AddressLine2 string `wire:"address_line2"`
	City         string
	PostCode     string
	Phone        string
	Mobile       string
	Fax          string
	Email        string
	CountryCode  string
	StateCode    string
}

// ToMap renders the request body.
func (a CartAddressInput) ToMap() map[string]any { return dto.ToMap(a) }

// CartItemInput is the add's body, no envelope; CustomerID takes the customer's internal_id.
type CartItemInput struct {
	// Pointers, so an unset number is omitted rather than sent 0.
	ProductID       *int64
	Quantity        *int64
	CustomerID      *int64
	SessionID       string
	ShippingAddress *CartAddressInput
}

func (i CartItemInput) ToMap() map[string]any { return dto.ToMap(i) }

// CartItemUpdateInput is the line update's body, no envelope; CartItemID and Quantity are required.
type CartItemUpdateInput struct {
	CartItemID      *int64
	Quantity        *int64
	ShippingAddress *CartAddressInput
	BillingAddress  *CartAddressInput
}

func (i CartItemUpdateInput) ToMap() map[string]any { return dto.ToMap(i) }

// CartSessionInput is the session write's body; an unset SessionID sends {}.
type CartSessionInput struct {
	SessionID string `wire:"sessionId"`
}

func (i CartSessionInput) ToMap() map[string]any { return dto.ToMap(i) }

// CartSession is the session write's answer.
type CartSession struct {
	SessionID *string `wire:"sessionId"`
}

// CartOrderInput is the order write's body, no envelope; CustomerID takes the customer's customer_id.
type CartOrderInput struct {
	CustomerID     *int64
	PaymentGateway string
	DeliveryType   string
	TransactionNo  string
}

func (i CartOrderInput) ToMap() map[string]any { return dto.ToMap(i) }

// Reference is a country or a state, as an object; a state arrives {} when the address has none.
type Reference struct {
	ID   *int64
	Name *string
	Code *string
}

// CartOrderCustomer is the order's customer_summary, carrying both of the customer's ids.
type CartOrderCustomer struct {
	FirstName     *string
	LastName      *string
	CustomerGroup *string
	Gender        *string
	CustomerID    *int64
	InternalID    *int64
	Email         *string
	AddressLine1  *string
	AddressLine2  *string
	City          *string
	Country       *Reference
	State         *Reference
	PostCode      *string
	Phone         *string
	Mobile        *string
	Fax           *string
	CompanyName   *string
}

// CartOrderLine is one entry of the order's order_line_details.
type CartOrderLine struct {
	ItemID            *int64
	ProductName       *string
	ProductID         *int64
	SKU               *string
	Quantity          *int64
	Price             *float64
	TotalAmount       *float64
	Tax               *float64
	Discount          *float64
	TaxDiscount       *float64
	IsTaxable         *bool
	IsShippable       *bool
	ImageURL          *string
	Variations        []any
	SerialBatchNumber *string
}

// CartOrderPayment is one entry of the order's payments.
type CartOrderPayment struct {
	ID                   *int64
	Status               *string
	Amount               *float64
	PaymentMethod        *string
	TrackInfo            *string
	GatewayName          *string
	GatewayCode          *string
	GatewayResponse      *string
	PayerInfo            *string
	PayingDate           *string
	TransactionDate      *string
	TransactionReference *string
}

// CartOrderAddress is the order's billing_address or shipping_address.
type CartOrderAddress struct {
	ID           *int64
	FirstName    *string
	LastName     *string
	Email        *string
	CompanyName  *string
	AddressLine1 *string
	AddressLine2 *string
	Country      *Reference
	State        *Reference
	City         *string
	PostCode     *string
	Mobile       *string
	Phone        *string
	Fax          *string
}

// CartOrder is the order a cart became, Status read from beside the envelope; TransactionNo is always nil here.
type CartOrder struct {
	Status                   *string
	OrderID                  *int64
	InternalID               *int64
	OrderNo                  *int64
	TransactionNo            *string
	CurrencyCode             *string
	OrderStatus              *string
	PaymentStatus            *string
	ShippingStatus           *string
	SubTotal                 *float64
	ShippingCost             *float64
	ShippingTax              *float64
	HandlingCost             *float64
	TotalSurcharge           *float64
	TotalDiscount            *float64
	TotalTax                 *float64
	GrandTotal               *float64
	Paid                     *float64
	Due                      *float64
	ItemsTotal               *int64
	IPAddress                *string
	OrderChannel             *string
	CreatedAt                *string
	UpdatedAt                *string
	CustomerSummary          *CartOrderCustomer
	OrderLineDetails         []*CartOrderLine
	Payments                 []*CartOrderPayment
	ExternalReferenceOrderID any
	ExternalReferenceStoreID any
	ExternalReferenceOrderNo any
	ExternalReferenceStatus  any
	ExternalReferenceURL     any
	CustomerPONumber         *string
	PaymentTerm              *string
	FulfillmentMode          *string
	FulfillmentTerm          *string
	InvoiceMode              *string
	BillingAddress           *CartOrderAddress
	ShippingAddress          *CartOrderAddress
	CustomFields             []any
}

// EnvelopeKey is the key the order arrives under.
func (CartOrder) EnvelopeKey() string { return CartOrderKey }

// CartShape is the typed view of a body shaped {"cart": {…}}.
func CartShape(decoded any) any { return dto.Read[Cart](decoded) }

// CartSessionShape is the typed view of the session write's answer.
func CartSessionShape(decoded any) any { return dto.Read[CartSession](decoded) }

// CartOrderShape is the typed view of a body shaped {"status": …, "order": {…}}.
func CartOrderShape(decoded any) any { return dto.Read[CartOrder](decoded) }
