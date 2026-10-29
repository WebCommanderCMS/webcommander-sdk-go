package customergroups

// The path templates; the client fills {api_version}. The literal sub-paths come before the row path.
const (
	CustomerGroupsPath = "/api/{api_version}/admin/customer_groups"

	CustomerGroupsCountPath = "/api/{api_version}/admin/customer_groups/count"

	CustomerGroupsAssignPath = "/api/{api_version}/admin/customer_groups/assign"

	CustomerGroupPath = "/api/{api_version}/admin/customer_groups/{customer_group_id}"

	CustomerGroupCustomersPath = "/api/{api_version}/admin/customer_groups/{customer_group_id}/customers"

	CustomerGroupCustomerPath = "/api/{api_version}/admin/customer_groups/{customer_group_id}/customers/{customer_id}"
)

// CustomerGroupsKey is the key the listing's rows arrive under — the resource's own plural name.
const CustomerGroupsKey = "customer_groups"

// CustomerGroupKey is the wrapper every group write takes and the key one group arrives under.
const CustomerGroupKey = "customer_group"

// CustomerGroupCustomersKey is the member write's body key, a list of customer_id values, and the row read's member list.
const CustomerGroupCustomersKey = "customers"

// CustomerGroupCustomerKey is the key one member arrives under on the member read.
const CustomerGroupCustomerKey = "customer"

// CustomerGroupsCountKey is the one key the count answers with, bare at the top of the body.
const CustomerGroupsCountKey = "total"

// CustomerGroupsPageSize is the page size applied when no limit is sent.
const CustomerGroupsPageSize = 20

// CustomerGroupsPageSizeCap is the largest limit that is not 400.
const CustomerGroupsPageSizeCap = 100

var customerGroupsSortFields = [6]string{"id", "name", "status", "description", "created", "updated"}

// CustomerGroupsSortFields returns a fresh copy of the names sort takes; created_at and updated_at, the names a row carries, are 400.
func CustomerGroupsSortFields() []string {
	return append(make([]string, 0, len(customerGroupsSortFields)), customerGroupsSortFields[:]...)
}
