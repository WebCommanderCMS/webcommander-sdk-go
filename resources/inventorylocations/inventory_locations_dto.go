package inventorylocations

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// InventoryLocation is one inventory location, the same seven fields in every form the API sends.
type InventoryLocation struct {
	ID        *int64
	Name      *string
	Address   *string
	Active    *bool
	Default   *bool
	CreatedAt *string
	UpdatedAt *string
}

// EnvelopeKey is the key a row path, a create and the sub-paths wrap the record in; no write takes it.
func (InventoryLocation) EnvelopeKey() string { return InventoryLocationKey }

// InventoryLocationInput is the writable form for the create, the update and the patch, with no envelope.
type InventoryLocationInput struct {
	Name    string
	Address string
	// Active and Default are pointers, so an unset flag is omitted rather than sent false.
	Active  *bool
	Default *bool
}

// ToMap renders the request body.
func (i InventoryLocationInput) ToMap() map[string]any { return dto.ToMap(i) }

// InventoryLocationStatusInput is the status sub-path's body; Active is required, hence a pointer.
type InventoryLocationStatusInput struct {
	Active *bool
}

func (i InventoryLocationStatusInput) ToMap() map[string]any { return dto.ToMap(i) }

// Pagination is the five-key pagination block of the listing and the stock sub-path; Limit is the page size applied.
type Pagination struct {
	Total       *int64
	Limit       *int64
	Offset      *int64
	HasPrevious *bool
	HasNext     *bool
}

// InventoryLocationPage is one page of the listing: the rows and the pagination block, unwrapped.
type InventoryLocationPage struct {
	InventoryLocations []*InventoryLocation
	// Pagination is nil when the body carried none.
	Pagination *Pagination
}

// InventoryLocationStockProduct is the product one stock item belongs to.
type InventoryLocationStockProduct struct {
	ID   *int64
	Name *string
	SKU  *string
}

// InventoryLocationStockItem is one product's holding at one location; DefaultForProduct is the product's own flag.
type InventoryLocationStockItem struct {
	Product           *InventoryLocationStockProduct
	Variation         any
	AvailableStock    *int64
	LowStockLevel     *int64
	DefaultForProduct *bool
}

// InventoryLocationStock is the stock sub-path's body: the location, its items and their pagination.
type InventoryLocationStock struct {
	InventoryLocation *InventoryLocation
	Items             []*InventoryLocationStockItem
	Pagination        *Pagination
}

// InventoryLocationShape is the typed view of a body shaped {"inventory_location": {…}}.
func InventoryLocationShape(decoded any) any { return dto.Read[InventoryLocation](decoded) }

// InventoryLocationPageShape is the typed view of the listing body.
func InventoryLocationPageShape(decoded any) any { return dto.Read[InventoryLocationPage](decoded) }

// InventoryLocationStockShape is the typed view of the stock sub-path's body.
func InventoryLocationStockShape(decoded any) any { return dto.Read[InventoryLocationStock](decoded) }
