package shippingrules

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// ShippingReference is an {id, name} pair: a rule's class, rate, zone or profile; a write names a class, a rate or a zone by ID.
type ShippingReference struct {
	ID   *int64
	Name *string
}

// ShippingRule is one rule: a listing row, the row read and every write's answer, the same ten keys; ShippingClass and ShippingRate are nil when the rule has none, and Zones and Profiles are empty lists.
type ShippingRule struct {
	ID            *int64
	Name          *string
	Description   *string
	ShippingClass *ShippingReference
	ShippingRate  *ShippingReference
	Zones         []*ShippingReference
	Profiles      []*ShippingReference
	TransactionNo *string
	CreatedAt     *string
	UpdatedAt     *string
}

// EnvelopeKey is the key a row read and every write's answer wrap the record in.
func (ShippingRule) EnvelopeKey() string { return ShippingRuleKey }

// ShippingRuleInput is the writable form for the create, the PUT replace and the PATCH, wrapped in shipping_rule; an unset field is not sent, so to send null for the class or the rate pass a map. ShippingClass and ShippingRate take an id or a ShippingReference, and Zones a list of ids or ShippingReferences, each zone once; ProfileID adds a created rule to that profile.
type ShippingRuleInput struct {
	Name          string
	Description   string
	ShippingClass any
	ShippingRate  any
	Zones         []any
	ProfileID     *int64
	TransactionNo string
	SuggestName   *bool
}

// EnvelopeKey is the wrapper every write requires; a bare record is 400.
func (ShippingRuleInput) EnvelopeKey() string { return ShippingRuleKey }

// ToMap renders the request body.
func (i ShippingRuleInput) ToMap() map[string]any { return dto.ToMap(i) }

// Pagination is the listing's five-key pagination block; Limit is the page size applied.
type Pagination struct {
	Total       *int64
	Limit       *int64
	Offset      *int64
	HasPrevious *bool
	HasNext     *bool
}

// ShippingRuleListing is one page of the listing, newest first.
type ShippingRuleListing struct {
	ShippingRules []*ShippingRule
	// Pagination is nil when the body carried none.
	Pagination *Pagination
}

// ShippingRuleCount is the count sub-path's body, a bare {"total": N}.
type ShippingRuleCount struct {
	Total *int64
}

// ShippingRuleShape is the typed view of a body shaped {"shipping_rule": {…}}.
func ShippingRuleShape(decoded any) any { return dto.Read[ShippingRule](decoded) }

// ShippingRuleListingShape is the typed view of the listing body.
func ShippingRuleListingShape(decoded any) any { return dto.Read[ShippingRuleListing](decoded) }

// ShippingRuleCountShape is the typed view of the count sub-path's body.
func ShippingRuleCountShape(decoded any) any { return dto.Read[ShippingRuleCount](decoded) }
