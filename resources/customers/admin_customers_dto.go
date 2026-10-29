package customers

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// Reference is a lookup value sent as an object, not a bare id — {"id": n, "name": "Australia", "code": "AU"}.
type Reference struct {
	ID   *int64
	Name *string
	Code *string
}

type Address struct {
	ID           *int64
	Type         *string
	FirstName    *string
	LastName     *string
	Email        *string
	CompanyName  *string
	AddressLine1 *string
	AddressLine2 *string
	Country      *Reference
	State        *Reference
	City         *string
	// PostCode is text rather than a number: a postal code can carry a leading zero.
	PostCode *string
	Mobile   *string
	Phone    *string
	Fax      *string
}

// CustomProperty is one entry of a customer's custom_properties — a label/description pair with an id of its own.
type CustomProperty struct {
	ID          *int64
	Label       *string
	Description *string
}

// Customer is one customer: the fields a listing sends, then the ones only a detail fetch adds.
type Customer struct {
	CustomerID      *int64
	InternalID      *int64
	UUID            *string
	FirstName       *string
	LastName        *string
	Name            *string
	Email           *string
	Status          *string
	Gender          *string
	CompanyName     *string
	StoreCredit     *float64
	CustomerGroups  []*Reference
	ShippingProfile any
	Addresses       []*Address
	CreatedAt       *string
	UpdatedAt       *string

	// Detail-only from here down; nil on a listing record.
	ABN         *string
	ABNBranch   *string
	AccountHold *bool
	// ActiveBillingAddress and ActiveShippingAddress are address IDs, not addresses.
	ActiveBillingAddress  *int64
	ActiveShippingAddress *int64
	AllowCreditLimit      *bool
	BackgroundImage       *string
	BaseURL               *string
	CountReferralCodeUsed *int64
	CreditLimit           *float64
	// An element is nil where the reader would not coerce the entry.
	CustomProperties        []*CustomProperty
	CustomerPasswordExpired *bool
	DefaultTaxCode          *string
	EmailVerified           *bool
	HidePrice               *bool
	HowDoYouKnow            *string
	// IsInTrash reports a deleted record: DELETE is soft, and a trashed customer still reads back 200.
	IsInTrash         *bool
	Newsletters       []any
	Note              *string
	PaymentTerm       any
	ProfilePicture    any
	ReferralCode      any
	SMSSubscribed     *bool
	Source            *string
	Timezone          *string
	TransactionNumber *string
	UsedReferralCode  any
	Version           *int64
}

// EnvelopeKey is the key the detail, create, update and addresses bodies wrap the record in; a list entry has none.
func (Customer) EnvelopeKey() string { return "customer" }

// CustomerInput is the writable form of Customer, sent through the customer envelope.
type CustomerInput struct {
	FirstName     string
	LastName      string
	Email         string
	EmailVerified *bool
	Password      string
	OldPassword   string
	Status        string
	Company       *bool
	CompanyName   string
	AddressLine1  string
	AddressLine2  string
	City          string
	State         any
	Country       any
	PostCode      string
	Addresses     []*Address
	// ID belongs to an entry's Label/Description pair, so it is honoured while that text is unchanged.
	CustomProperties []*CustomProperty
	TransactionNo    string
}

func (CustomerInput) EnvelopeKey() string { return "customer" }

// TopLevel names the fields that sit beside that envelope rather than inside it.
func (CustomerInput) TopLevel() []string { return []string{"transaction_no"} }

// ToMap renders the request body.
func (c CustomerInput) ToMap() map[string]any { return dto.ToMap(c) }

// StoreCreditAdjustment is one movement of a customer's store credit: the fields a history entry sends, then a write's.
type StoreCreditAdjustment struct {
	ID         *int64
	CustomerID *int64
	InternalID *int64
	Action     *string
	Amount     *float64
	Note       *string
	CreatedAt  *string
	CreatedBy  any

	// Write-only from here down; nil on a history entry.
	BalanceBefore *float64
	BalanceAfter  *float64
	StoreCredit   *float64
}

// EnvelopeKey is the key a write wraps the record in; a history entry carries none.
func (StoreCreditAdjustment) EnvelopeKey() string { return "store_credit_adjustment" }

// StoreCreditAdjustmentInput is the writable form of StoreCreditAdjustment, and the only body here with no envelope.
type StoreCreditAdjustmentInput struct {
	Action string
	Amount float64
	Note   string
}

func (s StoreCreditAdjustmentInput) ToMap() map[string]any { return dto.ToMap(s) }

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

// Page is one page of a listing: the records, and the block that says where they sit.
type Page[T any] struct {
	Items []*T
	// Pagination is nil when the body carried none.
	Pagination *Pagination
}

// CustomerShape is the typed view of a body shaped {"customer": {…}}.
func CustomerShape(decoded any) any { return dto.Read[Customer](decoded) }

// StoreCreditAdjustmentShape is the typed view of a body shaped {"store_credit_adjustment": {…}, "store_credit": 20.0}.
func StoreCreditAdjustmentShape(decoded any) any { return dto.Read[StoreCreditAdjustment](decoded) }

// PageOf reads a listing envelope as a Page of T.
func PageOf[T any](decoded any) any {
	body, isObject := decoded.(map[string]any)
	if !isObject {
		return nil
	}
	page := &Page[T]{}
	if records, sole := soleList(body); sole {
		page.Items = make([]*T, len(records))
		for at, record := range records {
			if item, isRecord := dto.Read[T](record).(*T); isRecord {
				page.Items[at] = item
			}
		}
	}
	if block, isBlock := dto.Read[Pagination](body["pagination"]).(*Pagination); isBlock {
		page.Pagination = block
	}
	return page
}

// soleList is the body's one list-valued entry, and whether there was exactly one of them.
func soleList(body map[string]any) ([]any, bool) {
	var records []any
	found := 0
	for _, value := range body {
		if list, isList := value.([]any); isList {
			records = list
			found++
		}
	}
	return records, found == 1
}
