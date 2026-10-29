package aliases

// The path templates; the client fills {api_version}.
const (
	AliasesPath = "/api/{api_version}/admin/aliases"

	AliasPath = "/api/{api_version}/admin/aliases/{alias_id}"
)

// The body keys: the listing's rows, and the create body's one field, which also keys the new row.
const (
	AliasesKey = "aliases"
	AliasKey   = "alias"
)

var aliasesQueryParams = [4]string{"limit", "offset", "page", "field_metadata"}

// AliasesQueryParams returns a fresh copy of the names the listing accepts; others are 400.
func AliasesQueryParams() []string {
	return append(make([]string, 0, len(aliasesQueryParams)), aliasesQueryParams[:]...)
}

// AliasesPageSize is the page size when limit is not sent, and the ceiling on one that is.
const AliasesPageSize = 20

var aliasesPaginationKeys = [5]string{"has_next", "has_previous", "limit", "offset", "total"}

// AliasesPaginationKeys returns a fresh copy of the five keys the pagination block carries.
func AliasesPaginationKeys() []string {
	return append(make([]string, 0, len(aliasesPaginationKeys)), aliasesPaginationKeys[:]...)
}
