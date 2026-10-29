package trash

// The path templates; the client fills {api_version}.
const (
	TrashPath = "/api/{api_version}/admin/trash"

	TrashTypesPath = "/api/{api_version}/admin/trash/types"

	TrashItemPath = "/api/{api_version}/admin/trash/{trash_type}/{trash_item_id}"

	TrashItemRestorePath = "/api/{api_version}/admin/trash/{trash_type}/{trash_item_id}/restore"

	TrashBulkRestorePath = "/api/{api_version}/admin/trash/bulk-restore"

	TrashBulkDeletePath = "/api/{api_version}/admin/trash/bulk-delete"
)

// TrashItemsKey is the key the listing's rows arrive under, and the key a bulk body lists its items under.
const TrashItemsKey = "items"

// TrashTypesKey is the key the type list arrives under.
const TrashTypesKey = "types"

// TrashPageSize is the page size applied when no limit is sent; a limit above TrashPageSizeCap is 400 out_of_range.
const TrashPageSize = 20

const TrashPageSizeCap = 100
