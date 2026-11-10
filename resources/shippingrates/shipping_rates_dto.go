package shippingrates

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// ShippingReference is an {id, name} pair: a rule that uses the rate.
type ShippingReference struct {
	ID   *int64
	Name *string
}

// ShippingRateCost is an amount and how it is charged; Type is fixed or percent.
type ShippingRateCost struct {
	Amount *float64
	Type   *string
}

// ShippingRateAdditionalRule is the extra charge for each further unit; enabled needs ForEachAdditional of at least 1 and AdditionalCost of at least 0.
type ShippingRateAdditionalRule struct {
	Enabled           *bool
	ForEachAdditional *int64
	AdditionalCost    *float64
}

// ShippingRate is one rate: a listing row, the row read and every write's answer, the same 14 keys. Which costs a rate carries follows its Method: a flat_rate has ShippingCost and HandlingCost at the rate level, free_shipping has no ShippingCost and no AdditionalRule, shipping_variant has none of the costs, and the ship_by_ methods carry their costs per row in Conditions, so the rate-level costs read nil. WeightUnit is "pounds" on ship_by_weight and nil elsewhere.
type ShippingRate struct {
	ID                      *int64
	Name                    *string
	Method                  *string
	CalculateForEntireOrder *bool
	PriceEnteredWithTax     *bool
	ShippingCost            *ShippingRateCost
	HandlingCost            *ShippingRateCost
	WeightUnit              *string
	// Conditions are the per-row price, quantity, weight or variation ranges; a row has a from key, so the rows stay untyped.
	Conditions     []any
	AdditionalRule *ShippingRateAdditionalRule
	ShippingRules  []*ShippingReference
	TransactionNo  *string
	CreatedAt      *string
	UpdatedAt      *string
}

// EnvelopeKey is the key a row read and every write's answer wrap the record in.
func (ShippingRate) EnvelopeKey() string { return ShippingRateKey }

// ShippingRateInput is the writable form for the create, the PUT replace and the PATCH, wrapped in shipping_rate; an unset field is not sent, so to send null pass a map. Conditions rows are maps: price and quantity rows carry from, to, shipping_cost and handling_cost, a weight row adds packet_weight, and a variation row carries variation_combination ([{"id": n}]) in place of the range. ShippingRuleID attaches a created rate to a rule that has none.
type ShippingRateInput struct {
	Name                    string
	Method                  string
	CalculateForEntireOrder *bool
	PriceEnteredWithTax     *bool
	ShippingCost            *ShippingRateCost
	HandlingCost            *ShippingRateCost
	Conditions              []any
	AdditionalRule          *ShippingRateAdditionalRule
	ShippingRuleID          *int64
	TransactionNo           string
	SuggestName             *bool
}

// EnvelopeKey is the wrapper every write requires; a bare record is 400.
func (ShippingRateInput) EnvelopeKey() string { return ShippingRateKey }

// ToMap renders the request body.
func (i ShippingRateInput) ToMap() map[string]any { return dto.ToMap(i) }

// Pagination is the listing's five-key pagination block; Limit is the page size applied.
type Pagination struct {
	Total       *int64
	Limit       *int64
	Offset      *int64
	HasPrevious *bool
	HasNext     *bool
}

// ShippingRateListing is one page of the listing, newest first.
type ShippingRateListing struct {
	ShippingRates []*ShippingRate
	// Pagination is nil when the body carried none.
	Pagination *Pagination
}

// ShippingRateCount is the count sub-path's body, a bare {"total": N}.
type ShippingRateCount struct {
	Total *int64
}

// ShippingRateShape is the typed view of a body shaped {"shipping_rate": {…}}.
func ShippingRateShape(decoded any) any { return dto.Read[ShippingRate](decoded) }

// ShippingRateListingShape is the typed view of the listing body.
func ShippingRateListingShape(decoded any) any { return dto.Read[ShippingRateListing](decoded) }

// ShippingRateCountShape is the typed view of the count sub-path's body.
func ShippingRateCountShape(decoded any) any { return dto.Read[ShippingRateCount](decoded) }
