package taxzones

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// TaxCountry is a country as a tax zone nests it; a zone write sends this record, as a bare id is 500.
type TaxCountry struct {
	ID   *int64
	Name *string
	Code *string
}

// ToMap lets a country stand on its own as a reference in a zone body.
func (t TaxCountry) ToMap() map[string]any { return dto.ToMap(t) }

// TaxState is a state as a tax zone nests it; a zone write sends this record too.
type TaxState struct {
	ID   *int64
	Name *string
	Code *string
}

// ToMap lets a state stand on its own as a reference in a zone body.
func (t TaxState) ToMap() map[string]any { return dto.ToMap(t) }

// TaxZone is one tax zone; PostCode reads as the flat list of strings a write sends, and SystemGenerated rows are the store's.
type TaxZone struct {
	ID              *int64
	Name            *string
	SystemGenerated *bool
	Default         *bool
	UseRadius       *bool
	Radius          *float64
	StartLocation   *string
	Countries       []*TaxCountry
	States          []*TaxState
	// Pointers, since the reader coerces into a declared pointer and leaves a bare string at its zero.
	PostCode  []*string
	CreatedAt *string
	UpdatedAt *string
}

// EnvelopeKey is the key a row path and the writes' answers wrap the record in, and both writes' wrapper.
func (TaxZone) EnvelopeKey() string { return TaxZoneKey }

// TaxZoneInput is the writable form; Name is required by both writes, and the update replaces the lists.
type TaxZoneInput struct {
	Name      string
	Countries []*TaxCountry
	States    []*TaxState
	PostCode  []string
	Radius    *float64
	// StartLocation is text, so an empty one is omitted; the store reads it only beside UseRadius.
	StartLocation string
	UseRadius     *bool
}

// EnvelopeKey is the wrapper both writes require.
func (TaxZoneInput) EnvelopeKey() string { return TaxZoneKey }

// ToMap renders the request body.
func (t TaxZoneInput) ToMap() map[string]any { return dto.ToMap(t) }

// Pagination is the listing's ten-key pagination block; Limit is the size asked for, up to TaxZonesPageSize.
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

// TaxZonePage is one page of the listing: the rows and the pagination block, unwrapped.
type TaxZonePage struct {
	TaxZones []*TaxZone
	// Pagination is nil when the body carried none.
	Pagination *Pagination
}

// TaxZoneShape is the typed view of a body shaped {"tax_zone": {…}}.
func TaxZoneShape(decoded any) any { return dto.Read[TaxZone](decoded) }

// TaxZonePageShape is the typed view of the listing body, {"tax_zones": […], "pagination": {…}}.
func TaxZonePageShape(decoded any) any { return dto.Read[TaxZonePage](decoded) }
