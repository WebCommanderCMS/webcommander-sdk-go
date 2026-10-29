package aliases

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// Alias is one alias, the same four fields on a listing row and in a create's answer.
type Alias struct {
	ID        *int64
	Alias     *string
	CreatedAt *string
	UpdatedAt *string
}

// EnvelopeKey is the key a create's answer wraps the record in.
func (Alias) EnvelopeKey() string { return AliasKey }

// AliasInput is the create body, {"alias": "<hostname>"}; it declares no envelope, as a wrapped one is 400.
type AliasInput struct {
	Alias string
}

// ToMap renders the request body.
func (a AliasInput) ToMap() map[string]any { return dto.ToMap(a) }

// Pagination is the listing's five-key pagination block; Limit is the size the store applied.
type Pagination struct {
	Total       *int64
	Limit       *int64
	Offset      *int64
	HasNext     *bool
	HasPrevious *bool
}

// AliasPage is one page of the listing; field_metadata adds no key here.
type AliasPage struct {
	Aliases    []*Alias
	Pagination *Pagination
}

// AliasShape is the typed view of a body shaped {"alias": {…}}.
func AliasShape(decoded any) any { return dto.Read[Alias](decoded) }

// AliasPageShape is the typed view of the listing body.
func AliasPageShape(decoded any) any { return dto.Read[AliasPage](decoded) }
