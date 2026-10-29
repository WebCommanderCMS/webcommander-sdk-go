package taxcodes

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// TaxCode is one tax code, the same eight keys on a listing row and a row path; Default is not exclusive.
type TaxCode struct {
	ID           *int64
	Name         *string
	Label        *string
	Description  *string
	Rate         *float64
	ResolverType *string
	Priority     *int64
	Default      *bool
}

// EnvelopeKey is the key a row path and the writes' answers wrap the record in, and both writes' wrapper.
func (TaxCode) EnvelopeKey() string { return TaxCodeKey }

// TaxCodeInput is the writable form; an update that leaves Priority or ResolverType unset resets it.
type TaxCodeInput struct {
	Name        string
	Label       string
	Description string
	// Pointers, so an unset number is omitted rather than stored as 0.
	Rate         *float64
	ResolverType string
	Priority     *int64
}

// EnvelopeKey is the wrapper both writes require.
func (TaxCodeInput) EnvelopeKey() string { return TaxCodeKey }

// ToMap renders the request body.
func (t TaxCodeInput) ToMap() map[string]any { return dto.ToMap(t) }

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

// TaxCodePage is one page of the listing: the rows and the pagination block, unwrapped.
type TaxCodePage struct {
	TaxCodes []*TaxCode
	// Pagination is nil when the body carried none.
	Pagination *Pagination
}

// TaxCodeShape is the typed view of a body shaped {"tax_code": {…}}.
func TaxCodeShape(decoded any) any { return dto.Read[TaxCode](decoded) }

// TaxCodePageShape is the typed view of the listing body, {"tax_codes": […], "pagination": {…}}.
func TaxCodePageShape(decoded any) any { return dto.Read[TaxCodePage](decoded) }
