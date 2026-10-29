package paymentgateways

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// PaymentGatewayReference is an {id, name} pair: the gateway's zone, and each audience member it names.
type PaymentGatewayReference struct {
	ID   *int64
	Name *string
}

// PaymentGatewayVisibility is who sees the gateway at checkout; the lists fill only under type selected.
type PaymentGatewayVisibility struct {
	Type           *string
	Customers      []*PaymentGatewayReference
	CustomerGroups []*PaymentGatewayReference
}

// PaymentGatewayAllowedFor is the delivery types the gateway is offered for; nil on gateways without the choice.
type PaymentGatewayAllowedFor struct {
	Shipping    *bool
	Delivery    *bool
	StorePickup *bool
}

// PaymentGatewaySurcharge is the surcharge block, nil on gateways without one, and its write; PercentAmount reads nil when no percentage is set.
type PaymentGatewaySurcharge struct {
	Type          *string
	FlatAmount    *float64
	PercentAmount *float64
	// Ranges holds each range as sent: {id, from, to, amount}.
	Ranges []any
}

// PaymentGatewayMessages is the pending and completed messages; a message the gateway's form lacks reads nil.
type PaymentGatewayMessages struct {
	Pending   *string
	Completed *string
}

// PaymentGatewayPayPalStandard is PayPal's standard credentials.
type PaymentGatewayPayPalStandard struct {
	ClientID     *string
	ClientSecret *string
}

// PaymentGatewayPayPalClassic is PayPal's classic account.
type PaymentGatewayPayPalClassic struct {
	EmailAddress *string
}

// PaymentGatewayToggle is a one-key {"enabled": …} switch.
type PaymentGatewayToggle struct {
	Enabled *bool
}

// PaymentGatewayCardProcessor is the card gateway's processor; the client id and secret read ******** when set.
type PaymentGatewayCardProcessor struct {
	Type          *string
	Mode          *string
	MerchantCode  *string
	ClientID      *string
	ClientSecret  *string
	WalletPayment *PaymentGatewayToggle
}

// PaymentGatewayCardLogo is the card logos shown at checkout.
type PaymentGatewayCardLogo struct {
	Type        *string
	Cards       []string
	CustomImage any
}

// PaymentGatewayGooglePay is Express Checkout's Google Pay block.
type PaymentGatewayGooglePay struct {
	Enabled                     *bool
	PaymentMode                 *string
	MerchantName                *string
	MerchantID                  *string
	Processor                   *string
	SecretKey                   *string
	PublicKey                   *string
	FreeShippingConditionAmount *float64
	ShippingCost                *float64
}

// PaymentGatewayApplePay is Express Checkout's Apple Pay block.
type PaymentGatewayApplePay struct {
	Enabled            *bool
	PaymentMode        *string
	MerchantIdentifier *string
	DomainName         *string
	DisplayName        *string
	Processor          *string
	SecretKey          *string
	PublicKey          *string
}

// PaymentGatewayExpressPayPal is Express Checkout's PayPal block.
type PaymentGatewayExpressPayPal struct {
	ShowInstallmentAmount *bool
}

// PaymentGatewaySettings is the row read's gateway-specific settings, every key any gateway answers; each gateway fills its own few and offline gateways answer {}.
type PaymentGatewaySettings struct {
	Version               *string
	Mode                  *string
	Standard              *PaymentGatewayPayPalStandard
	Classic               *PaymentGatewayPayPalClassic
	Processor             *PaymentGatewayCardProcessor
	CardLogo              *PaymentGatewayCardLogo
	GooglePay             *PaymentGatewayGooglePay
	ApplePay              *PaymentGatewayApplePay
	Paypal                *PaymentGatewayExpressPayPal
	ShowInstallmentAmount *bool
	ShowInProductDetails  *bool
	APIKey                *string
	DefaultTerm           *int64
	AmountPer100Points    *float64
	MerchantID            *string
	MerchantSecretKey     *string
	Image                 any
}

// PaymentGateway is one gateway: a listing row, and the row path's and every write's body, which add Settings.
type PaymentGateway struct {
	ID          *int64
	Code        *string
	Name        *string
	Added       *bool
	Enabled     *bool
	Default     *bool
	Promotional *bool
	Zone        *PaymentGatewayReference
	Information *string
	Visibility  *PaymentGatewayVisibility
	AllowedFor  *PaymentGatewayAllowedFor
	Surcharge   *PaymentGatewaySurcharge
	Messages    *PaymentGatewayMessages
	// Settings is nil on a listing row, which carries none.
	Settings *PaymentGatewaySettings
}

// EnvelopeKey is the key the row read and every write wrap the record in; no write takes it.
func (PaymentGateway) EnvelopeKey() string { return PaymentGatewayKey }

// PaymentGatewayVisibilityInput is the visibility a write takes: customers and customer groups by id.
type PaymentGatewayVisibilityInput struct {
	Type           string
	Customers      []int64
	CustomerGroups []int64
}

// PaymentGatewayInput is the writable form for the PUT and the PATCH, with no envelope; code and promotional are read-only.
type PaymentGatewayInput struct {
	Name string
	// Added, Enabled and Default are pointers, so an unset one is omitted rather than sent as false.
	Added   *bool
	Enabled *bool
	Default *bool
	// Zone takes a zone's id or a PaymentGatewayReference naming its ID.
	Zone        any
	Information string
	Visibility  *PaymentGatewayVisibilityInput
	AllowedFor  *PaymentGatewayAllowedFor
	Surcharge   *PaymentGatewaySurcharge
	Messages    *PaymentGatewayMessages
	Settings    *PaymentGatewaySettings
}

// ToMap renders the request body.
func (i PaymentGatewayInput) ToMap() map[string]any { return dto.ToMap(i) }

// Pagination is the listing's five-key pagination block; Limit is the page size applied, and HasPrevious is true past the first page and false on a page past the end.
type Pagination struct {
	Total       *int64
	Limit       *int64
	Offset      *int64
	HasPrevious *bool
	HasNext     *bool
}

// PaymentGatewayPage is one page of the listing: the rows, without settings, and where they sit.
type PaymentGatewayPage struct {
	PaymentGateways []*PaymentGateway
	// Pagination is nil when the body carried none.
	Pagination *Pagination
}

// PaymentGatewayCount is the count sub-path's body, a bare {"total": N}.
type PaymentGatewayCount struct {
	Total *int64
}

// PaymentGatewayShape is the typed view of a body shaped {"payment_gateway": {…}}.
func PaymentGatewayShape(decoded any) any { return dto.Read[PaymentGateway](decoded) }

// PaymentGatewayPageShape is the typed view of the listing body.
func PaymentGatewayPageShape(decoded any) any { return dto.Read[PaymentGatewayPage](decoded) }

// PaymentGatewayCountShape is the typed view of the count sub-path's body.
func PaymentGatewayCountShape(decoded any) any { return dto.Read[PaymentGatewayCount](decoded) }
