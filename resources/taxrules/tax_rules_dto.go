package taxrules

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// TaxRuleCode is the tax code a rule points at, as the rule nests it; no write here sets it.
type TaxRuleCode struct {
	ID        *int64
	Name      *string
	IsDefault *bool
}

// TaxRuleZone is one tax zone a rule points at, four keys as the rule nests it; no write here sets it.
type TaxRuleZone struct {
	ID                *int64
	Name              *string
	IsSystemGenerated *bool
	IsDefault         *bool
}

// TaxRule is one tax rule; Code, Zones, DecimalPoint, RoundingType and Default are read-only.
type TaxRule struct {
	ID           *int64
	Name         *string
	Code         *TaxRuleCode
	Description  *string
	Default      *bool
	RoundingType *string
	DecimalPoint *int64
	Zones        []*TaxRuleZone
	CreatedAt    *string
	UpdatedAt    *string
}

// EnvelopeKey is the key a row path and the writes' answers wrap the record in, and both writes' wrapper.
func (TaxRule) EnvelopeKey() string { return TaxRuleKey }

// TaxRuleInput is the writable form, two of the row's ten fields; the update merges.
type TaxRuleInput struct {
	Name        string
	Description string
}

// EnvelopeKey is the wrapper both writes require.
func (TaxRuleInput) EnvelopeKey() string { return TaxRuleKey }

// ToMap renders the request body.
func (t TaxRuleInput) ToMap() map[string]any { return dto.ToMap(t) }

// Pagination is the listing's ten-key pagination block; Limit is the size asked for, capped at 100.
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

// TaxRulePage is one page of the listing: the rows and the pagination block, unwrapped.
type TaxRulePage struct {
	TaxRules []*TaxRule
	// Pagination is nil when the body carried none.
	Pagination *Pagination
}

// TaxRuleShape is the typed view of a body shaped {"tax_rule": {…}}.
func TaxRuleShape(decoded any) any { return dto.Read[TaxRule](decoded) }

// TaxRulePageShape is the typed view of the listing body, {"tax_rules": […], "pagination": {…}}.
func TaxRulePageShape(decoded any) any { return dto.Read[TaxRulePage](decoded) }
