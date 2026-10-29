package taxcustomprofiles

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// CustomProfileReference is one profile a rule is filed under, on the rule's Profiles.
type CustomProfileReference struct {
	ID   *int64
	Name *string
}

// CustomProfileRule is one entry of Rules, which only the row path sends.
type CustomProfileRule struct {
	ID            *int64
	Name          *string
	Description   *string
	ZonesCount    *int64
	Profiles      []*CustomProfileReference
	CreatedOn     *string
	LastUpdatedOn *string
}

// CustomProfile is one custom tax profile; the row path adds Rules to the eight keys a listing row and both writes' answers carry.
type CustomProfile struct {
	ID            *int64
	Name          *string
	TaxCode       *string
	Description   *string
	IsDefault     *bool
	RulesCount    *int64
	CreatedOn     *string
	LastUpdatedOn *string
	Rules         []*CustomProfileRule
}

// EnvelopeKey is the key a row path and the writes' answers wrap the record in, and both writes' wrapper.
func (CustomProfile) EnvelopeKey() string { return TaxCustomProfileKey }

// CustomProfileInput is the writable form; the create requires Name, and the update replaces but keeps an omitted Name.
type CustomProfileInput struct {
	Name        string
	TaxCode     string
	Description string
	IsDefault   *bool
	// TransactionNo files the create under that number; the update files it only on a profile no number holds yet.
	TransactionNo string
}

// EnvelopeKey is the wrapper both writes require.
func (CustomProfileInput) EnvelopeKey() string { return TaxCustomProfileKey }

// TopLevel names the fields that sit beside that envelope; inside it the store files nothing.
func (CustomProfileInput) TopLevel() []string { return []string{"transaction_no"} }

// ToMap renders the request body.
func (c CustomProfileInput) ToMap() map[string]any { return dto.ToMap(c) }

// ResourceMeta is one group_meta entry: what the path is, whether it is writable and where.
type ResourceMeta struct {
	Kind     *string
	Writable *bool
	Href     *string
}

// GroupMeta is the group_meta block only the listing sends.
type GroupMeta struct {
	CustomProfile *ResourceMeta
	Rules         *ResourceMeta
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

// CustomProfilePage is one page of the listing: the rows, group_meta and the pagination block, unwrapped.
type CustomProfilePage struct {
	CustomProfiles []*CustomProfile
	GroupMeta      *GroupMeta
	// Pagination is nil when the body carried none.
	Pagination *Pagination
}

// CustomProfileShape is the typed view of a body shaped {"custom_profile": {…}}.
func CustomProfileShape(decoded any) any { return dto.Read[CustomProfile](decoded) }

// CustomProfilePageShape is the typed view of the listing body.
func CustomProfilePageShape(decoded any) any { return dto.Read[CustomProfilePage](decoded) }
