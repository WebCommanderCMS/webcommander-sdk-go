package customergroups

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// CustomerGroups is the customer group endpoints, bound to one client.
type CustomerGroups struct {
	requester http.Requester
}

// ListCustomerGroups lists the groups by name — GET /admin/customer_groups; limit 1 to 100, search_text or name match a name, status takes A, I or their words.
func (c *CustomerGroups) ListCustomerGroups(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, CustomerGroupsPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CustomerGroupPageShape})
}

// CountCustomerGroups counts the groups — GET /admin/customer_groups/count; {"total": N}, narrowed by search_text, name and status and not by paging.
func (c *CustomerGroups) CountCustomerGroups(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, CustomerGroupsCountPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CustomerGroupCountShape})
}

// CreateCustomerGroup creates a group — POST /admin/customer_groups, wrapped in customer_group; 201, name required and unique in any letter case.
func (c *CustomerGroups) CreateCustomerGroup(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendPostRequest(ctx, CustomerGroupsPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CustomerGroupShape})
}

// GetCustomerGroup retrieves one group with customer_count and its customers — GET /admin/customer_groups/{customer_group_id}; field_metadata adds the schema.
func (c *CustomerGroups) GetCustomerGroup(ctx context.Context, customerGroupID any, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, CustomerGroupPath, &http.RequestOptions{
		PathParams:  map[string]any{"customer_group_id": customerGroupID},
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CustomerGroupShape})
}

// PatchCustomerGroup changes only the fields sent, nested ones included — PATCH /admin/customer_groups/{customer_group_id}; PUT is 405.
func (c *CustomerGroups) PatchCustomerGroup(ctx context.Context, customerGroupID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendPatchRequest(ctx, CustomerGroupPath, &http.RequestOptions{
		PathParams: map[string]any{"customer_group_id": customerGroupID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CustomerGroupShape})
}

// DeleteCustomerGroup deletes one group and its memberships — DELETE /admin/customer_groups/{customer_group_id}; 204, the members stay as customers.
func (c *CustomerGroups) DeleteCustomerGroup(ctx context.Context, customerGroupID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendDeleteRequest(ctx, CustomerGroupPath, &http.RequestOptions{
		PathParams:   map[string]any{"customer_group_id": customerGroupID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// UpdateCustomerGroupCustomers replaces the member list with {"customers": [customer_id, ...]} — PATCH /admin/customer_groups/{customer_group_id}/customers; [] empties it.
func (c *CustomerGroups) UpdateCustomerGroupCustomers(ctx context.Context, customerGroupID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendPatchRequest(ctx, CustomerGroupCustomersPath, &http.RequestOptions{
		PathParams: map[string]any{"customer_group_id": customerGroupID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CustomerGroupShape})
}

// GetCustomerGroupCustomer reads one member by its customer_id — GET /admin/customer_groups/{customer_group_id}/customers/{customer_id}; a non-member is 404.
func (c *CustomerGroups) GetCustomerGroupCustomer(ctx context.Context, customerGroupID any, customerID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, CustomerGroupCustomerPath, &http.RequestOptions{
		PathParams:   map[string]any{"customer_group_id": customerGroupID, "customer_id": customerID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CustomerGroupMemberShape})
}

// AssignCustomerGroup adds one customer to one group, keeping the other members — POST /admin/customer_groups/assign; 201, added false for an existing member.
func (c *CustomerGroups) AssignCustomerGroup(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendPostRequest(ctx, CustomerGroupsAssignPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CustomerGroupAssignmentShape})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "customer_groups"

func oneResponseFormat(responseFormat []string) (string, error) {
	switch len(responseFormat) {
	case 0:
		return "", nil
	case 1:
		return responseFormat[0], nil
	default:
		return "", wcerrors.NewUsageError(fmt.Sprintf(
			"responseFormat: pass one format or none, got %d — the trailing argument is optional, "+
				"and omitting it defers to the client", len(responseFormat)))
	}
}

func requirePayload(payload any) error {
	if payload == nil || isNilValue(payload) {
		return wcerrors.NewUsageError("payload is required: this endpoint takes a body, and nil would " +
			"send the request with no body at all")
	}
	return nil
}

func isNilValue(value any) bool {
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Map, reflect.Slice, reflect.Pointer, reflect.Interface, reflect.Func:
		return reflected.IsNil()
	}
	return false
}

// New binds the customer group endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *CustomerGroups {
	return &CustomerGroups{requester: requester}
}
