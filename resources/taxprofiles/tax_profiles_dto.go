package taxprofiles

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// TaxProfileRuleCode is the tax code a nested rule points at.
type TaxProfileRuleCode struct {
	ID        *int64
	Name      *string
	IsDefault *bool
}

// TaxProfileRuleZone is one tax zone a nested rule points at.
type TaxProfileRuleZone struct {
	ID                *int64
	Name              *string
	IsSystemGenerated *bool
	IsDefault         *bool
}

// TaxProfileRule is one tax rule as a profile nests it, the same ten keys a /admin/tax_rules read sends.
type TaxProfileRule struct {
	ID           *int64
	Name         *string
	Code         *TaxProfileRuleCode
	Description  *string
	Default      *bool
	RoundingType *string
	DecimalPoint *int64
	Zones        []*TaxProfileRuleZone
	CreatedAt    *string
	UpdatedAt    *string
}

// TaxProfile is one tax profile; Default is exclusive, and the tax settings' default_tax_profile follows it.
type TaxProfile struct {
	ID          *int64
	Name        *string
	Description *string
	Default     *bool
	// TaxRules is empty on a new profile.
	TaxRules  []*TaxProfileRule
	CreatedAt *string
	UpdatedAt *string
}

// EnvelopeKey is the key a row path and the writes' answers wrap the record in, and both writes' wrapper.
func (TaxProfile) EnvelopeKey() string { return TaxProfileKey }

// TaxProfileInput is the writable form; the update merges, and Default is one-way, so it is a pointer.
type TaxProfileInput struct {
	Name        string
	Description string
	Default     *bool
}

// EnvelopeKey is the wrapper both writes require.
func (TaxProfileInput) EnvelopeKey() string { return TaxProfileKey }

// ToMap renders the request body.
func (t TaxProfileInput) ToMap() map[string]any { return dto.ToMap(t) }

// Pagination is the listing's ten-key pagination block; Limit is the size asked for, up to TaxProfilesPageSize.
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

// TaxProfilePage is one page of the listing: the rows and the pagination block, unwrapped.
type TaxProfilePage struct {
	TaxProfiles []*TaxProfile
	// Pagination is nil when the body carried none.
	Pagination *Pagination
}

// TaxProfileShape is the typed view of a body shaped {"tax_profile": {…}}.
func TaxProfileShape(decoded any) any { return dto.Read[TaxProfile](decoded) }

// TaxProfilePageShape is the typed view of the listing body.
func TaxProfilePageShape(decoded any) any { return dto.Read[TaxProfilePage](decoded) }
