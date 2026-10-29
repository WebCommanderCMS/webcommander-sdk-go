package transactions

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// TransactionListEntry is one listing row: a number, how many records carry it, and its read path.
type TransactionListEntry struct {
	TransactionNo *string
	TotalRecords  *int64
	URL           *string
}

// TransactionTypeCount is one by_type row of a number's read; the store form-encodes the number in ItemsURL, a space as +.
type TransactionTypeCount struct {
	Type *string
	// Count is Created plus Updated; Trashed counts those of them now in the trash.
	Count   *int64
	Created *int64
	Updated *int64
	Trashed *int64
	// ItemsURL is sent on every row, a Count of 0 included.
	ItemsURL *string
}

// Transaction is one number's read: a bare body, its 15 by_type rows always in the same order.
type Transaction struct {
	TransactionNo *string
	TotalRecords  *int64
	ByType        []*TransactionTypeCount
}

// TransactionItem is one detail=slim record under a transaction number: Type, ID, Stamp, Trashed plus the fields its type sends, the rest nil; read a detail=full row from Map.
type TransactionItem struct {
	Type *string
	ID   *int64
	// Stamp is "created" or "updated".
	Stamp       *string
	Trashed     *bool
	Name        *string
	SKU         *string
	Rating      *float64
	Email       *string
	FirstName   *string
	LastName    *string
	OrderNo     *int64
	GrandTotal  *float64
	Status      *string
	Code        *string
	Symbol      *string
	CountryCode *string
	StateName   *string
	Postcode    *string
	Rate        *float64
	Priority    *int64
	Description *string
	PolicyType  *string
	URL         *string
}

// Pagination is the five-key pagination block of both listings; there is no next-page link, so step offset.
type Pagination struct {
	Total       *int64
	Limit       *int64
	Offset      *int64
	HasPrevious *bool
	HasNext     *bool
}

// TransactionPage is one page of the numbers, unwrapped; a transaction_no filter answers Pagination too.
type TransactionPage struct {
	Transactions []*TransactionListEntry
	Pagination   *Pagination
}

// TransactionItemPage is one page of the records carrying a number, beside the number; each item names its own Type, and a number or a type with no records answers Items empty and Pagination.Total 0.
type TransactionItemPage struct {
	TransactionNo *string
	Items         []*TransactionItem
	Pagination    *Pagination
}

// TransactionShape is the typed view of a number's bare read body.
func TransactionShape(decoded any) any { return dto.Read[Transaction](decoded) }

// TransactionPageShape is the typed view of the number listing body.
func TransactionPageShape(decoded any) any { return dto.Read[TransactionPage](decoded) }

// TransactionItemPageShape is the typed view of the items listing body.
func TransactionItemPageShape(decoded any) any { return dto.Read[TransactionItemPage](decoded) }
