package orderpreviews

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/orders"
)

// PreviewSession is the session block. SessionID starts PRV- and is the preview_session_id of every other call; ExpiresAt is 30 minutes after CreatedAt, and a coupon or gift card call, applying or removing, moves it 30 minutes past that call, while a read does not.
type PreviewSession struct {
	SessionID *string
	Status    *string
	CreatedAt *string
	ExpiresAt *string
}

// PreviewCustomer is the customer as the preview shows it: the billing address names win over the record's.
type PreviewCustomer struct {
	CustomerID *int64
	FirstName  *string
	LastName   *string
	Email      *string
}

// PreviewItem is one line of the preview, priced by the store.
type PreviewItem struct {
	ProductID   *int64
	ProductName *string
	SKU         *string
	Quantity    *int64
	Price       *float64
	Discount    *float64
	Tax         *float64
	TotalAmount *float64
	Taxable     *bool
	Shippable   *bool
}

// PreviewShipping is the shipping block; ShippingTax is not part of the totals block.
type PreviewShipping struct {
	DeliveryType *string
	ShippingType *string
	ShippingCost *float64
	ShippingTax  *float64
}

// PreviewCoupon is the applied coupon.
type PreviewCoupon struct {
	Code   *string
	Status *string
	Amount *float64
}

// PreviewGiftCard is the applied gift card; RemainingBalance is always nil today.
type PreviewGiftCard struct {
	Code             *string
	Status           *string
	AppliedAmount    *float64
	RemainingBalance *float64
}

// PreviewDiscounts holds the coupon and the gift card, each nil when none is applied.
type PreviewDiscounts struct {
	Coupon   *PreviewCoupon
	GiftCard *PreviewGiftCard
}

// PreviewTotals are the priced totals: GrandTotal is SubTotal minus TotalDiscount plus TotalTax, the shipping tax, ShippingCost, HandlingCost and TotalSurcharge, and Due is GrandTotal minus GiftCardApplied.
type PreviewTotals struct {
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

// OrderPreview is a preview session: what the start and the read answer. The coupon and gift card calls answer only Session, Discounts and Totals, so the rest stay nil.
type OrderPreview struct {
	Session      *PreviewSession
	Customer     *PreviewCustomer
	DeliveryType *string
	CurrencyCode *string
	Items        []*PreviewItem
	Shipping     *PreviewShipping
	Discounts    *PreviewDiscounts
	Totals       *PreviewTotals
}

// EnvelopeKey is the key every preview answer wraps the record in.
func (OrderPreview) EnvelopeKey() string { return OrderPreviewKey }

// OrderPreviewInput is the start body, with no envelope. Products is required, each line a product_id and a quantity of at least 1. Without CustomerID a billing address is required, and a shipping address too unless DeliveryType is store_pickup. DeliveryType is shipping (the default), delivery, others_shipping or store_pickup.
type OrderPreviewInput struct {
	CustomerID      *int64
	DeliveryType    string
	BillingAddress  *orders.OrderAddressInput
	ShippingAddress *orders.OrderAddressInput
	Products        []*orders.OrderProductInput
}

// ToMap renders the request body.
func (i OrderPreviewInput) ToMap() map[string]any { return dto.ToMap(i) }

// OrderPreviewCodeInput is the body of the coupon and the gift card call.
type OrderPreviewCodeInput struct {
	Code string
}

// ToMap renders the request body.
func (i OrderPreviewCodeInput) ToMap() map[string]any { return dto.ToMap(i) }

// OrderPreviewPaymentInput is the payment a confirm records; PaymentMethod is required, and Amount may not exceed what is due.
type OrderPreviewPaymentInput struct {
	Amount               *float64
	PaymentStatus        string
	PaymentMethod        string
	TransactionReference string
	TransactionDate      string
}

// OrderPreviewConfirmInput is the confirm body, with no envelope; every field is optional.
type OrderPreviewConfirmInput struct {
	PaymentGateway string
	TransactionNo  string
	SendEmail      *bool
	Payment        *OrderPreviewPaymentInput
}

// ToMap renders the request body.
func (i OrderPreviewConfirmInput) ToMap() map[string]any { return dto.ToMap(i) }

// OrderPreviewShape is the typed view of a body shaped {"order": {…}} in the preview shape.
func OrderPreviewShape(decoded any) any { return dto.Read[OrderPreview](decoded) }
