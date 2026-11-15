package carts

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/orders"
)

// AdminCartSession is the session block. SessionID starts PRV- and is the sessionID of every other call; ExpiresAt is 30 minutes after the last write.
type AdminCartSession struct {
	SessionID *string
	Status    *string
	CreatedAt *string
	ExpiresAt *string
}

// AdminCartCustomer is the customer the cart is for; CustomerID is the customer's normal id.
type AdminCartCustomer struct {
	CustomerID *int64
	FirstName  *string
	LastName   *string
	Email      *string
}

// AdminCartItem is one line, priced by the store. CartItemID is what the line's update and delete take; a replace of the cart numbers its lines again.
type AdminCartItem struct {
	CartItemID    *int64
	ProductID     *int64
	ProductName   *string
	VariationID   *int64
	VariationName *string
	SKU           *string
	Price         *float64
	Quantity      *int64
	Discount      *float64
	Tax           *float64
	TotalAmount   *float64
	Taxable       *bool
	Shippable     *bool
}

// AdminCartTotals is the cart's money before shipping.
type AdminCartTotals struct {
	SubTotal      *float64
	TotalDiscount *float64
	TotalTax      *float64
	GrandTotal    *float64
}

// AdminCart is a cart: what every cart call but the deletes answers.
type AdminCart struct {
	Session      *AdminCartSession
	Customer     *AdminCartCustomer
	CurrencyCode *string
	Items        []*AdminCartItem
	Totals       *AdminCartTotals
}

// EnvelopeKey is the key every cart answer wraps the record in.
func (AdminCart) EnvelopeKey() string { return AdminCartKey }

// CheckoutAddress is a billing or shipping address as the checkout reads it back; State is an empty Reference when the address has none.
type CheckoutAddress struct {
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

// CheckoutShipping is the delivery the checkout was priced for; empty until a delivery_type is set.
type CheckoutShipping struct {
	DeliveryType *string
	ShippingType *string
	ShippingCost *float64
	ShippingTax  *float64
}

// CheckoutCoupon is the coupon applied: AppliedAmount is the discount it takes off.
type CheckoutCoupon struct {
	Code          *string
	Status        *string
	AppliedAmount *float64
}

// CheckoutGiftCard is the gift card applied: AppliedAmount is the lesser of its balance and the grand total.
type CheckoutGiftCard struct {
	Code             *string
	Status           *string
	AppliedAmount    *float64
	RemainingBalance *float64
}

// CheckoutDiscounts holds at most one coupon and one gift card, each nil when none is applied.
type CheckoutDiscounts struct {
	Coupon   *CheckoutCoupon
	GiftCard *CheckoutGiftCard
}

// CheckoutTotals is the money. A gift card lowers Due, not GrandTotal.
type CheckoutTotals struct {
	SubTotal        *float64
	TotalDiscount   *float64
	TotalTax        *float64
	ShippingCost    *float64
	HandlingCost    *float64
	TotalSurcharge  *float64
	GiftCardApplied *float64
	GrandTotal      *float64
	Due             *float64
}

// CheckoutIssue is one thing still blocking a confirm, such as delivery_type required.
type CheckoutIssue struct {
	Field   *string
	Code    *string
	Message *string
}

// AdminCheckout is the checkout. The read and the update answer every field; the coupon, gift card and delivery schedule calls answer the parts they change. DeliverySchedule is nil until one is set.
type AdminCheckout struct {
	Session          *AdminCartSession
	Customer         *AdminCartCustomer
	DeliveryType     *string
	CurrencyCode     *string
	Items            []*AdminCartItem
	BillingAddress   *CheckoutAddress
	ShippingAddress  *CheckoutAddress
	Shipping         *CheckoutShipping
	Discounts        *CheckoutDiscounts
	DeliverySchedule any
	Totals           *CheckoutTotals
	Issues           []*CheckoutIssue
}

// EnvelopeKey is the key every checkout answer wraps the record in.
func (AdminCheckout) EnvelopeKey() string { return AdminCheckoutKey }

// AdminCartProductInput is one product line to put on a cart. ProductID and Quantity are required; VariationID picks the variation of a variable product.
type AdminCartProductInput struct {
	ProductID   *int64
	VariationID *int64
	Quantity    *int64
}

// AdminCartInput is the create and replace body, bare. Products is required, and so is CustomerID on a replace; SessionID on a create names the session yourself, in the PRV- form. An unknown field is 400 unknown_field.
type AdminCartInput struct {
	CustomerID *int64
	SessionID  string
	Products   []*AdminCartProductInput
}

// ToMap renders the request body.
func (i AdminCartInput) ToMap() map[string]any { return dto.ToMap(i) }

// AdminCartCustomerInput is the cart update body: the customer to move the cart to. CustomerID is required.
type AdminCartCustomerInput struct {
	CustomerID *int64
}

// ToMap renders the request body.
func (i AdminCartCustomerInput) ToMap() map[string]any { return dto.ToMap(i) }

// AdminCartItemsInput is the add-lines body. Adding a product already on the cart raises that line's quantity.
type AdminCartItemsInput struct {
	Products []*AdminCartProductInput
}

// ToMap renders the request body.
func (i AdminCartItemsInput) ToMap() map[string]any { return dto.ToMap(i) }

// AdminCartItemUpdateInput is the line update body: the new Quantity, required.
type AdminCartItemUpdateInput struct {
	Quantity *int64
}

// ToMap renders the request body.
func (i AdminCartItemUpdateInput) ToMap() map[string]any { return dto.ToMap(i) }

// AdminCheckoutInput is the checkout update body. DeliveryType is shipping, delivery, others_shipping or store_pickup; an address replaces the whole one.
type AdminCheckoutInput struct {
	DeliveryType    string
	BillingAddress  *orders.OrderAddressInput
	ShippingAddress *orders.OrderAddressInput
}

// ToMap renders the request body.
func (i AdminCheckoutInput) ToMap() map[string]any { return dto.ToMap(i) }

// AdminCheckoutCodeInput is the body of a coupon or gift card call: Code, required.
type AdminCheckoutCodeInput struct {
	Code string
}

// ToMap renders the request body.
func (i AdminCheckoutCodeInput) ToMap() map[string]any { return dto.ToMap(i) }

// AdminCheckoutPaymentInput is the payment a confirm records. Amount and PaymentMethod are required; an Amount above what is due is refused.
type AdminCheckoutPaymentInput struct {
	Amount               *float64
	PaymentStatus        string
	PaymentMethod        string
	TransactionReference string
	TransactionDate      string
}

// AdminCheckoutConfirmInput is the confirm body, bare, and optional. PaymentGateway takes one of the store's codes, such as PIS; an unknown field is 400.
type AdminCheckoutConfirmInput struct {
	PaymentGateway string
	TransactionNo  string
	SendEmail      *bool
	Payment        *AdminCheckoutPaymentInput
}

// ToMap renders the request body.
func (i AdminCheckoutConfirmInput) ToMap() map[string]any { return dto.ToMap(i) }

// AdminCheckoutDeliveryScheduleInput is the delivery schedule to choose: its ID and the Date wanted, both required.
type AdminCheckoutDeliveryScheduleInput struct {
	ID   *int64
	Date string
}

// ToMap renders the request body.
func (i AdminCheckoutDeliveryScheduleInput) ToMap() map[string]any { return dto.ToMap(i) }

// AdminCartShape is the typed view of a body shaped {"cart": {…}}.
func AdminCartShape(decoded any) any { return dto.Read[AdminCart](decoded) }

// AdminCheckoutShape is the typed view of a body shaped {"checkout": {…}}.
func AdminCheckoutShape(decoded any) any { return dto.Read[AdminCheckout](decoded) }
