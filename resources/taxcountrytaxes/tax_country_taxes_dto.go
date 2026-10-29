package taxcountrytaxes

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// TaxCountry is a country as a country tax nests it, read-only; nil on the store's rest-of-the-world row.
type TaxCountry struct {
	ID   *int64
	Name *string
	Code *string
}

// ToMap lets a country stand on its own as a reference in a body that takes one.
func (t TaxCountry) ToMap() map[string]any { return dto.ToMap(t) }

// TaxState is a state as a state rule nests it.
type TaxState struct {
	ID   *int64
	Name *string
	Code *string
}

// CountryTaxStateRule is one entry of StateRules, which only the row path and the update send.
type CountryTaxStateRule struct {
	ID               *int64
	Name             *string
	State            *TaxState
	AddressLine1     *string
	AddressLine2     *string
	City             *string
	PostCode         *string
	TaxRate          *float64
	TaxRule          *string
	EffectiveTaxRate *float64
	IsDefault        *bool
	CreatedOn        *string
	LastUpdatedOn    *string
	LastUpdatedBy    any
}

// CountryTax is one country tax; the row path and the update add StateRules to the twelve keys a listing row and the create carry.
type CountryTax struct {
	ID      *int64
	Name    *string
	Country *TaxCountry
	TaxCode *string
	TaxRate *float64
	// TaxMethod reads lower case: fixed_rate, api or no_tax.
	TaxMethod       *string
	IsDefault       *bool
	IsRestOfWorld   *bool
	StateRulesCount *int64
	CreatedOn       *string
	LastUpdatedOn   *string
	LastUpdatedBy   any
	StateRules      []*CountryTaxStateRule
}

// EnvelopeKey is the key a row path and the writes' answers wrap the record in, and both writes' wrapper.
func (CountryTax) EnvelopeKey() string { return TaxCountryTaxKey }

// CountryTaxInput is the writable form; the create requires CountryCode, and the update merges.
type CountryTaxInput struct {
	CountryCode string
	Name        string
	TaxCode     string
	TaxRate     *float64
	TaxMethod   string
	IsDefault   *bool
	// TransactionNo files the create under that number; the update re-files the row under it.
	TransactionNo string
}

// EnvelopeKey is the wrapper both writes require.
func (CountryTaxInput) EnvelopeKey() string { return TaxCountryTaxKey }

// TopLevel names the fields that sit beside that envelope; inside it the store files nothing.
func (CountryTaxInput) TopLevel() []string { return []string{"transaction_no"} }

// ToMap renders the request body.
func (c CountryTaxInput) ToMap() map[string]any { return dto.ToMap(c) }

// ResourceMeta is one group_meta entry: what the path is, whether it is writable and where.
type ResourceMeta struct {
	Kind     *string
	Writable *bool
	Href     *string
}

// GroupMeta is the group_meta block only the listing sends.
type GroupMeta struct {
	CountryTax *ResourceMeta
	StateRules *ResourceMeta
}

// Pagination is the listing's ten-key pagination block; Limit is the size applied, capped at 100.
type Pagination struct {
	Total        *int64
	Limit        *int64
	Offset       *int64
	Count        *int64
	CurrentPage  *int64
	TotalPages   *int64
	HasNext      *bool
	HasPrevious  *bool
	PreviousPage *string
	NextPage     *string
}

// CountryTaxPage is one page of the listing: the rows, group_meta and the pagination block, unwrapped.
type CountryTaxPage struct {
	CountryTaxes []*CountryTax
	GroupMeta    *GroupMeta
	// Pagination is nil when the body carried none.
	Pagination *Pagination
}

// CountryTaxShape is the typed view of a body shaped {"country_tax": {…}}.
func CountryTaxShape(decoded any) any { return dto.Read[CountryTax](decoded) }

// CountryTaxPageShape is the typed view of the listing body.
func CountryTaxPageShape(decoded any) any { return dto.Read[CountryTaxPage](decoded) }
