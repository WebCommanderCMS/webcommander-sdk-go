package giftcards

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// GiftCard is one gift card. Code is what an order preview's gift-card call takes; a card can be applied only while IsActive and IsPaid, with a Balance above zero, before its AvailableTo date and with Email equal to the previewed customer's email.
type GiftCard struct {
	ID                *int64
	Code              *string
	Balance           *float64
	InitialAmount     *float64
	Email             *string
	FirstName         *string
	LastName          *string
	Message           *string
	IsActive          *bool
	IsPaid            *bool
	AvailableTo       *string
	ExternalReference *string
	ProductName       *string
	CreatedOn         *string
}

// Pagination is the listing's ten-key pagination block; Limit is the size the store applied, at most 20.
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

// GiftCardPage is one page of the listing, newest first.
type GiftCardPage struct {
	GiftCards []*GiftCard
	// Pagination is nil when the body carried none.
	Pagination *Pagination
}

// GiftCardPageShape is the typed view of the listing body.
func GiftCardPageShape(decoded any) any { return dto.Read[GiftCardPage](decoded) }
