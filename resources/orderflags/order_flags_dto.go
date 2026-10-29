package orderflags

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// OrderFlag is one order flag, the same five fields on a listing row and a row path; Type is not writable.
type OrderFlag struct {
	ID          *int64
	Name        *string
	Color       *string
	Description *string
	Type        *string
}

// EnvelopeKey is the key a row path, a create and an update wrap the record in; a listing row has none.
func (OrderFlag) EnvelopeKey() string { return OrderFlagKey }

// OrderFlagInput is the writable form; both writes require name and a hex color, and the update replaces.
type OrderFlagInput struct {
	Name        string
	Color       string
	Description string
}

func (OrderFlagInput) EnvelopeKey() string { return OrderFlagKey }

// ToMap renders the request body.
func (o OrderFlagInput) ToMap() map[string]any { return dto.ToMap(o) }

// Pagination is the listing's five-key pagination block; Limit is the size applied.
type Pagination struct {
	Total       *int64
	Limit       *int64
	Offset      *int64
	HasNext     *bool
	HasPrevious *bool
}

// OrderFlagPage is one page of the listing: the rows and the pagination block, unwrapped.
type OrderFlagPage struct {
	OrderFlags []*OrderFlag
	// Pagination is nil when the body carried none.
	Pagination *Pagination
}

// OrderFlagShape is the typed view of a body shaped {"order_flag": {…}}.
func OrderFlagShape(decoded any) any { return dto.Read[OrderFlag](decoded) }

// OrderFlagPageShape is the typed view of the listing body, {"order_flags": […], "pagination": {…}}.
func OrderFlagPageShape(decoded any) any { return dto.Read[OrderFlagPage](decoded) }
