package currencies

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// CurrencyCountry is the country a currency is filed under, as every row answers it.
type CurrencyCountry struct {
	ID   *int64
	Code *string
	Name *string
}

// Currency is one currency: a listing row, the row path's body and every row write's answer; the base write answers only ID, Name, Code and Base of it.
type Currency struct {
	ID     *int64
	Name   *string
	Code   *string
	Symbol *string
	// Country is nil on the base write's answer, which carries none.
	Country *CurrencyCountry
	// ConversionRate is a float even when sent as an integer.
	ConversionRate *float64
	Rounding       *string
	Precision      *int64
	Active         *bool
	Base           *bool
	UpdatedAt      *string
	// FieldMetadata is the raw field_metadata block, present when field_metadata=true is asked for.
	FieldMetadata any
}

// EnvelopeKey is the key a row read and every row write wrap the record in; no write takes it.
func (Currency) EnvelopeKey() string { return CurrencyKey }

// CurrencyInput is the writable form for the create, the PUT replace and the PATCH, with no envelope.
type CurrencyInput struct {
	Name   string
	Code   string
	Symbol string
	// Country takes a country's id, code, name or a CurrencyCountry naming its Code.
	Country any
	// ConversionRate, Precision and Active are pointers, so an unset one is omitted rather than sent as zero.
	ConversionRate *float64
	Rounding       string
	Precision      *int64
	Active         *bool
}

// ToMap renders the request body.
func (i CurrencyInput) ToMap() map[string]any { return dto.ToMap(i) }

// Pagination is the listing's five-key pagination block; Limit is the page size applied.
type Pagination struct {
	Total       *int64
	Limit       *int64
	Offset      *int64
	HasPrevious *bool
	HasNext     *bool
}

// CurrencyPage is one page of the listing: the rows, where they sit and any schema.
type CurrencyPage struct {
	Currencies []*Currency
	// Pagination is nil when the body carried none.
	Pagination    *Pagination
	FieldMetadata any
}

// CurrencyCount is the count sub-path's body, a bare {"total": N}.
type CurrencyCount struct {
	Total *int64
}

// CurrencyReference is the row that held the base flag before a base write.
type CurrencyReference struct {
	ID   *int64
	Code *string
	Name *string
}

// CurrencyBaseChange is the base write's body: the new base, and the row that held the flag, nil when none moved.
type CurrencyBaseChange struct {
	Currency     *Currency
	PreviousBase *CurrencyReference
}

// CurrencyShape is the typed view of a body shaped {"currency": {…}}.
func CurrencyShape(decoded any) any { return dto.Read[Currency](decoded) }

// CurrencyPageShape is the typed view of the listing body.
func CurrencyPageShape(decoded any) any { return dto.Read[CurrencyPage](decoded) }

// CurrencyCountShape is the typed view of the count sub-path's body.
func CurrencyCountShape(decoded any) any { return dto.Read[CurrencyCount](decoded) }

// CurrencyBaseChangeShape is the typed view of the base write's body.
func CurrencyBaseChangeShape(decoded any) any { return dto.Read[CurrencyBaseChange](decoded) }
