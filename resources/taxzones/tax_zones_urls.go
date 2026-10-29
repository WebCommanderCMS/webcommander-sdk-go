package taxzones

// The path templates; the client fills {api_version}.
const (
	TaxZonesPath = "/api/{api_version}/admin/tax_zones"

	TaxZonePath = "/api/{api_version}/admin/tax_zones/{tax_zone_id}"
)

// TaxZonesKey is the key the listing's rows arrive under — the resource's own plural name.
const TaxZonesKey = "tax_zones"

// TaxZoneKey is the key one row arrives under, and the wrapper both writes require.
const TaxZoneKey = "tax_zone"

var taxZonesQueryParams = [3]string{"limit", "offset", "page"}

// TaxZonesQueryParams returns a fresh copy of the three names the listing acts on.
func TaxZonesQueryParams() []string {
	return append(make([]string, 0, len(taxZonesQueryParams)), taxZonesQueryParams[:]...)
}

// TaxZonesPageSize is the page size when limit is not sent, and the ceiling on one that is.
const TaxZonesPageSize = 20

var taxZonesPaginationKeys = [10]string{
	"count", "current_page", "has_next", "has_previous", "limit",
	"next_page", "offset", "previous_page", "total", "total_pages",
}

// TaxZonesPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func TaxZonesPaginationKeys() []string {
	return append(make([]string, 0, len(taxZonesPaginationKeys)), taxZonesPaginationKeys[:]...)
}

// UnroutedMethodBody is the body the store answers at 200, in place of a 405, for an unrouted method.
func UnroutedMethodBody() map[string]any {
	return map[string]any{"isSuccess": false, "message": "Invalid API Request"}
}
