package inventorylocations

// The path templates; the client fills {api_version}.
const (
	InventoryLocationsPath = "/api/{api_version}/admin/inventory_locations"

	InventoryLocationPath = "/api/{api_version}/admin/inventory_locations/{inventory_location_id}"

	InventoryLocationStockPath = "/api/{api_version}/admin/inventory_locations/{inventory_location_id}/stock"

	InventoryLocationDefaultPath = "/api/{api_version}/admin/inventory_locations/{inventory_location_id}/default"

	InventoryLocationStatusPath = "/api/{api_version}/admin/inventory_locations/{inventory_location_id}/status"
)

// InventoryLocationsKey is the key the listing's rows arrive under — the resource's own plural name.
const InventoryLocationsKey = "inventory_locations"

// InventoryLocationKey is the key one row arrives under; no write takes it, as each takes a bare record.
const InventoryLocationKey = "inventory_location"

// InventoryLocationStockKey is the key the stock sub-path's rows arrive under.
const InventoryLocationStockKey = "items"

var inventoryLocationsQueryParams = [3]string{"q", "limit", "offset"}

// InventoryLocationsQueryParams returns a fresh copy of the three names the listing accepts; any other is 400.
func InventoryLocationsQueryParams() []string {
	return append(make([]string, 0, len(inventoryLocationsQueryParams)), inventoryLocationsQueryParams[:]...)
}

// InventoryLocationsPageSize is the page size applied when no limit is sent; limit takes 1 to 100.
const InventoryLocationsPageSize = 20

// UnroutedMethodBody is the body the store answers at 200, in place of a 405, for an unrouted method.
func UnroutedMethodBody() map[string]any {
	return map[string]any{"isSuccess": false, "message": "Invalid API Request"}
}
