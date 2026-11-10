package shippingprofiles

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// ShippingReference is an {id, name} pair: a rule's rate, or a zone the wizard names; a write sends {"id": …} for an existing zone and {"name": …} to create one.
type ShippingReference struct {
	ID   *int64
	Name *string
}

// ShippingProfileRule is a rule on a profile; ShippingRate is nil on a rule that has no rate.
type ShippingProfileRule struct {
	ID           *int64
	Name         *string
	ShippingRate *ShippingReference
}

// ShippingProfile is one profile: a listing row, the row read and every write's answer, the same ten keys; RulePreference reads product_assigned_rule, lowest_shipping_cost_rule or highest_shipping_cost_rule, and exactly one profile has Default true.
type ShippingProfile struct {
	ID             *int64
	Name           *string
	Description    *string
	RulePreference *string
	Default        *bool
	RulesCount     *int64
	ShippingRules  []*ShippingProfileRule
	TransactionNo  *string
	CreatedAt      *string
	UpdatedAt      *string
}

// EnvelopeKey is the key a row read and every write's answer wrap the record in.
func (ShippingProfile) EnvelopeKey() string { return ShippingProfileKey }

// ShippingProfileInput is the writable form for the create, the PUT replace and the PATCH, wrapped in shipping_profile; an unset field is not sent, so to send null pass a map. ShippingRules is the whole list of rule ids, each once, and an empty list takes every rule off; Default takes only true or false.
type ShippingProfileInput struct {
	Name           string
	Description    string
	RulePreference string
	Default        *bool
	ShippingRules  []int64
	TransactionNo  string
	SuggestName    *bool
}

// EnvelopeKey is the wrapper every write requires; a bare record is 400.
func (ShippingProfileInput) EnvelopeKey() string { return ShippingProfileKey }

// ToMap renders the request body.
func (i ShippingProfileInput) ToMap() map[string]any { return dto.ToMap(i) }

// ShippingProfileWizardRate is the rate a wizard rule creates: the fields of a rate write, Name and Method required.
type ShippingProfileWizardRate struct {
	Name                    string
	Method                  string
	CalculateForEntireOrder *bool
	PriceEnteredWithTax     *bool
	ShippingCost            any
	HandlingCost            any
	Conditions              []any
	AdditionalRule          any
	TransactionNo           string
}

// ShippingProfileWizardRule is a rule the wizard creates beside its rate: Name is required, ShippingClass is a class id, and Zones name existing zones by ID or new ones by Name.
type ShippingProfileWizardRule struct {
	Name          string
	Description   string
	ShippingClass *int64
	ShippingRate  *ShippingProfileWizardRate
	Zones         []ShippingReference
}

// ShippingProfileWizardInput is the wizard's body, wrapped in shipping_profile: the profile fields and the rules, each with its rate, to create under it.
type ShippingProfileWizardInput struct {
	Name           string
	Description    string
	RulePreference string
	ShippingRules  []ShippingProfileWizardRule
}

// EnvelopeKey is the wrapper the wizard requires.
func (ShippingProfileWizardInput) EnvelopeKey() string { return ShippingProfileKey }

// ToMap renders the request body.
func (i ShippingProfileWizardInput) ToMap() map[string]any { return dto.ToMap(i) }

// Pagination is the listing's five-key pagination block; Limit is the page size applied.
type Pagination struct {
	Total       *int64
	Limit       *int64
	Offset      *int64
	HasPrevious *bool
	HasNext     *bool
}

// ShippingProfileListing is one page of the listing, newest first.
type ShippingProfileListing struct {
	ShippingProfiles []*ShippingProfile
	// Pagination is nil when the body carried none.
	Pagination *Pagination
}

// ShippingProfileCount is the count sub-path's body, a bare {"total": N}.
type ShippingProfileCount struct {
	Total *int64
}

// ShippingProfileShape is the typed view of a body shaped {"shipping_profile": {…}}.
func ShippingProfileShape(decoded any) any { return dto.Read[ShippingProfile](decoded) }

// ShippingProfileListingShape is the typed view of the listing body.
func ShippingProfileListingShape(decoded any) any { return dto.Read[ShippingProfileListing](decoded) }

// ShippingProfileCountShape is the typed view of the count sub-path's body.
func ShippingProfileCountShape(decoded any) any { return dto.Read[ShippingProfileCount](decoded) }
