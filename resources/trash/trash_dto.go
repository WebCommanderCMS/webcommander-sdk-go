package trash

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// TrashItem is one trashed row, and the item a restore answers, named by its type and its own resource's id; the trash keeps no id of its own.
type TrashItem struct {
	// Type is the snake_case name the type list answers, whatever case was sent.
	Type *string
	ID   *int64
	Name *string
	// RemovedAt is when the row was trashed, in UTC; a restore and a bulk delete send none.
	RemovedAt *string
	// URL is the record's own path, nil for a blog_post.
	URL *string
}

// EnvelopeKey is the key a restore answers the row under; a listing or bulk row carries none and reads as is.
func (TrashItem) EnvelopeKey() string { return "item" }

// Pagination is the listing's five-key pagination block; there is no next-page link, so step offset.
type Pagination struct {
	Total       *int64
	Limit       *int64
	Offset      *int64
	HasPrevious *bool
	HasNext     *bool
}

// TrashPage is one page of the listing: the rows and where they sit.
type TrashPage struct {
	Items []*TrashItem
	// Pagination is nil when the body carried none.
	Pagination *Pagination
}

// TrashType is one type the trash holds, as a path's {trash_type} takes it.
type TrashType struct {
	Type *string
}

// TrashTypeList is the type list's body: every type, alphabetical.
type TrashTypeList struct {
	Types []*TrashType
}

// TrashBulkError is one item a bulk call could not act on: an error Code such as not_found or not_in_trash, and why in words.
type TrashBulkError struct {
	Type    *string
	ID      *int64
	Code    *string
	Message *string
}

// TrashBulkRestoreResult is the bulk restore's body: the rows restored, and an entry per failed item.
type TrashBulkRestoreResult struct {
	Restored []*TrashItem
	Failed   []*TrashBulkError
}

// TrashBulkDeleteResult is the bulk delete's body: the rows deleted, each Type, ID and Name, and an entry per failed item.
type TrashBulkDeleteResult struct {
	Deleted []*TrashItem
	Failed  []*TrashBulkError
}

// TrashBulkItemInput is one item of a bulk body; Type in any letter case, and an extra key is 400.
type TrashBulkItemInput struct {
	Type string
	ID   *int64
}

// TrashBulkInput is the body both bulk calls take, {"items": [...]} with 1 to 100 items and no envelope.
type TrashBulkInput struct {
	Items []*TrashBulkItemInput
}

// ToMap renders the request body.
func (i TrashBulkInput) ToMap() map[string]any { return dto.ToMap(i) }

// TrashItemShape is the typed view of a restore's answer, the row under item.
func TrashItemShape(decoded any) any { return dto.Read[TrashItem](decoded) }

// TrashPageShape is the typed view of the listing body.
func TrashPageShape(decoded any) any { return dto.Read[TrashPage](decoded) }

// TrashTypeListShape is the typed view of the type list's body.
func TrashTypeListShape(decoded any) any { return dto.Read[TrashTypeList](decoded) }

// TrashBulkRestoreResultShape is the typed view of the bulk restore's body.
func TrashBulkRestoreResultShape(decoded any) any { return dto.Read[TrashBulkRestoreResult](decoded) }

// TrashBulkDeleteResultShape is the typed view of the bulk delete's body.
func TrashBulkDeleteResultShape(decoded any) any { return dto.Read[TrashBulkDeleteResult](decoded) }
