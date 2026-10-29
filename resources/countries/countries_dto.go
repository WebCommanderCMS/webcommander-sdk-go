package countries

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// Country is one country: a listing row, and the row path's body.
type Country struct {
	ID         *int64
	Code       *string
	Name       *string
	StateCount *int64
	Active     *bool
	Default    *bool
}

// EnvelopeKey is the key the row path wraps the record in; a listing row carries none.
func (Country) EnvelopeKey() string { return CountryKey }

// CountryReference is the short parent country a state listing and a state read carry beside the state.
type CountryReference struct {
	ID   *int64
	Code *string
	Name *string
}

// CountryState is one state: a listing row, and the state path's body with its parent country beside it.
type CountryState struct {
	ID      *int64
	Code    *string
	Name    *string
	Active  *bool
	Default *bool
	// Country is nil on a listing row; the state read sends it beside the envelope.
	Country *CountryReference
}

// EnvelopeKey is the key the state path wraps the record in.
func (CountryState) EnvelopeKey() string { return CountryStateKey }

// Pagination is the five-key pagination block of both listings; under limit=all Limit reads the rows that arrived.
type Pagination struct {
	Total       *int64
	Limit       *int64
	Offset      *int64
	HasPrevious *bool
	HasNext     *bool
}

// CountryPage is one page of the country listing: the rows and the pagination block, unwrapped.
type CountryPage struct {
	Countries  []*Country
	Pagination *Pagination
}

// CountryStatePage is one page of a country's states, beside the short parent country.
type CountryStatePage struct {
	Country    *CountryReference
	States     []*CountryState
	Pagination *Pagination
}

// CountryShape is the typed view of a body shaped {"country": {…}}.
func CountryShape(decoded any) any { return dto.Read[Country](decoded) }

// CountryPageShape is the typed view of the country listing body.
func CountryPageShape(decoded any) any { return dto.Read[CountryPage](decoded) }

// CountryStateShape is the typed view of a body shaped {"country": {…}, "state": {…}}.
func CountryStateShape(decoded any) any { return dto.Read[CountryState](decoded) }

// CountryStatePageShape is the typed view of a state listing body.
func CountryStatePageShape(decoded any) any { return dto.Read[CountryStatePage](decoded) }
