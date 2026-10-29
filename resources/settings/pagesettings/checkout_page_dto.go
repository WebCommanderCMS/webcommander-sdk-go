package pagesettings

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// AddressField is one entry of an address section's fields: a form field and how it is drawn.
type AddressField struct {
	Key      *string
	Label    *string
	Active   *bool
	Required *bool
	Order    *int64
	Group    *string
}

// AddressForm is an address section's record: active_fields from the facade, fields from the section.
type AddressForm struct {
	// An element is nil where the reader would not coerce the entry.
	ActiveFields []*string
	Fields       []*AddressField
}

// BillingAddress is the billing_address section, read-only, so it has no ToMap.
type BillingAddress AddressForm

// EnvelopeKey is the section's own name, the key a read arrives under.
func (BillingAddress) EnvelopeKey() string { return "billing_address" }

// ShippingAddress is the shipping_address section, read-only like BillingAddress.
type ShippingAddress AddressForm

// EnvelopeKey is the section's own name, the key a read arrives under.
func (ShippingAddress) EnvelopeKey() string { return "shipping_address" }

// SignUp is the allow_sign_up block: whether checkout may create an account, and whether it demands one.
type SignUp struct {
	Enabled  *bool
	Required *bool
}

// GuestCheckout is the checkout_as_guest block: guest orders, and the one-click sign-up offered.
type GuestCheckout struct {
	Enabled     *bool
	QuickSignup *bool
}

// CustomerCheckout is the customer_checkout section; turning all three ways to check out off is 400.
type CustomerCheckout struct {
	AllowSignIn             *bool
	AllowSignUp             *SignUp
	CheckoutAsGuest         *GuestCheckout
	NewsletterSubscription  *bool
	MultipleShippingMethods *bool
}

// EnvelopeKey is the section's own name, the one wrapper a write takes.
func (CustomerCheckout) EnvelopeKey() string { return "customer_checkout" }

// ToMap renders the request body.
func (c CustomerCheckout) ToMap() map[string]any { return dto.ToMap(c) }

// DeliveryOptions is the delivery_options section; turning all four options off is 400.
type DeliveryOptions struct {
	EnableShipping        *bool
	EnableStorePickup     *bool
	EnableShippingQuote   *bool
	EnableDeliveryService *bool
}

// EnvelopeKey is the section's own name, the one wrapper a write takes.
func (DeliveryOptions) EnvelopeKey() string { return "delivery_options" }

func (d DeliveryOptions) ToMap() map[string]any { return dto.ToMap(d) }

// OrderSummary is the order_summary section: what the checkout's summary panel shows.
type OrderSummary struct {
	ShowProductThumbnail    *bool
	ShowQuantity            *bool
	EnablePaymentExpandView *bool
	OrderConfirmCaptcha     *bool
}

// EnvelopeKey is the section's own name, the one wrapper a write takes.
func (OrderSummary) EnvelopeKey() string { return "order_summary" }

func (o OrderSummary) ToMap() map[string]any { return dto.ToMap(o) }

// OrderComment is the comment block; SendNotification requires Enabled.
type OrderComment struct {
	Enabled          *bool
	SendNotification *bool
}

// TermsReference is the terms_and_condition.ref block; the body carries only the key matching Type.
type TermsReference struct {
	Page any
	Link *string
	Text *string
}

// TermsAndCondition is the terms_and_condition block: whether terms must be accepted, and where they live.
type TermsAndCondition struct {
	Enabled   *bool
	Text      *string
	LinkLabel *string
	Type      *string
	Ref       *TermsReference
}

// Others is the others section: the checkout page's logo, login link, own HTML and two blocks.
type Others struct {
	ShowLogo          *bool
	LoginURL          *string
	CustomHTML        *string
	Comment           *OrderComment
	TermsAndCondition *TermsAndCondition
}

// EnvelopeKey is the section's own name, the one wrapper a write takes.
func (Others) EnvelopeKey() string { return "others" }

func (o Others) ToMap() map[string]any { return dto.ToMap(o) }

// CheckoutPageSectionMeta is one group_meta entry: what a section is, whether writable, where managed.
type CheckoutPageSectionMeta struct {
	Kind     *string
	Writable *bool
	Href     *string
}

// CheckoutPageGroupMeta is the group_meta block, one entry per section; only the facade sends it.
type CheckoutPageGroupMeta struct {
	BillingAddress   *CheckoutPageSectionMeta
	ShippingAddress  *CheckoutPageSectionMeta
	CustomerCheckout *CheckoutPageSectionMeta
	DeliveryOptions  *CheckoutPageSectionMeta
	OrderSummary     *CheckoutPageSectionMeta
	Others           *CheckoutPageSectionMeta
}

// CheckoutPage is the facade's body or one section's; ToMap renders the facade PATCH, not a section's.
type CheckoutPage struct {
	BillingAddress   *BillingAddress
	ShippingAddress  *ShippingAddress
	CustomerCheckout *CustomerCheckout
	DeliveryOptions  *DeliveryOptions
	OrderSummary     *OrderSummary
	Others           *Others
	GroupMeta        *CheckoutPageGroupMeta
	// FieldMetadata is the raw field_metadata block, keyed by section on the facade and by field on a section.
	FieldMetadata any
}

func (CheckoutPage) EnvelopeKey() string { return "checkout_page" }

func (c CheckoutPage) ToMap() map[string]any { return dto.ToMap(c) }

// CheckoutPageShape is the typed view of the facade's body and of one section's.
func CheckoutPageShape(decoded any) any { return dto.Read[CheckoutPage](decoded) }
