package inventorylocations

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// InventoryLocations is the inventory location endpoints, bound to one client.
type InventoryLocations struct {
	requester http.Requester
}

// ListInventoryLocations lists the locations — GET /admin/inventory_locations; q, limit (1 to 100) and offset only, any other name is 400.
func (i *InventoryLocations) ListInventoryLocations(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return i.requester.SendGetRequest(ctx, InventoryLocationsPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: InventoryLocationPageShape})
}

// HeadInventoryLocations asks whether the collection is reachable — HEAD /admin/inventory_locations.
func (i *InventoryLocations) HeadInventoryLocations(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return i.requester.SendApiRequest(ctx, "HEAD", InventoryLocationsPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CreateInventoryLocation creates a location — POST /admin/inventory_locations as a bare record; default true moves the flag.
func (i *InventoryLocations) CreateInventoryLocation(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return i.requester.SendPostRequest(ctx, InventoryLocationsPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: InventoryLocationShape})
}

// GetInventoryLocation retrieves one location — GET /admin/inventory_locations/{inventory_location_id}.
func (i *InventoryLocations) GetInventoryLocation(ctx context.Context, inventoryLocationID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return i.requester.SendGetRequest(ctx, InventoryLocationPath, &http.RequestOptions{
		PathParams:   map[string]any{"inventory_location_id": inventoryLocationID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: InventoryLocationShape})
}

// HeadInventoryLocation asks whether one location is there — HEAD /admin/inventory_locations/{inventory_location_id}.
func (i *InventoryLocations) HeadInventoryLocation(ctx context.Context, inventoryLocationID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return i.requester.SendApiRequest(ctx, "HEAD", InventoryLocationPath, &http.RequestOptions{
		PathParams:   map[string]any{"inventory_location_id": inventoryLocationID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// UpdateInventoryLocation replaces one location — PUT /admin/inventory_locations/{inventory_location_id}; an omitted address is cleared, an omitted active reads true.
func (i *InventoryLocations) UpdateInventoryLocation(ctx context.Context, inventoryLocationID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return i.requester.SendPutRequest(ctx, InventoryLocationPath, &http.RequestOptions{
		PathParams: map[string]any{"inventory_location_id": inventoryLocationID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: InventoryLocationShape})
}

// PatchInventoryLocation merges fields into one location — PATCH /admin/inventory_locations/{inventory_location_id} as a bare record of at least one field.
func (i *InventoryLocations) PatchInventoryLocation(ctx context.Context, inventoryLocationID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return i.requester.SendPatchRequest(ctx, InventoryLocationPath, &http.RequestOptions{
		PathParams: map[string]any{"inventory_location_id": inventoryLocationID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: InventoryLocationShape})
}

// DeleteInventoryLocation deletes a location — DELETE /admin/inventory_locations/{inventory_location_id}; the default row is 409.
func (i *InventoryLocations) DeleteInventoryLocation(ctx context.Context, inventoryLocationID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return i.requester.SendDeleteRequest(ctx, InventoryLocationPath, &http.RequestOptions{
		PathParams:   map[string]any{"inventory_location_id": inventoryLocationID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// GetInventoryLocationStock reads what one location holds — GET /admin/inventory_locations/{inventory_location_id}/stock; paged by limit and offset.
func (i *InventoryLocations) GetInventoryLocationStock(ctx context.Context, inventoryLocationID any, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return i.requester.SendGetRequest(ctx, InventoryLocationStockPath, &http.RequestOptions{
		PathParams:  map[string]any{"inventory_location_id": inventoryLocationID},
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: InventoryLocationStockShape})
}

// SetInventoryLocationDefault moves the default flag to one location — PUT /admin/inventory_locations/{inventory_location_id}/default.
func (i *InventoryLocations) SetInventoryLocationDefault(ctx context.Context, inventoryLocationID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return i.requester.SendPutRequest(ctx, InventoryLocationDefaultPath, &http.RequestOptions{
		PathParams:   map[string]any{"inventory_location_id": inventoryLocationID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: InventoryLocationShape})
}

// UpdateInventoryLocationStatus sets active on one location — PUT /admin/inventory_locations/{inventory_location_id}/status.
func (i *InventoryLocations) UpdateInventoryLocationStatus(ctx context.Context, inventoryLocationID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return i.requester.SendPutRequest(ctx, InventoryLocationStatusPath, &http.RequestOptions{
		PathParams: map[string]any{"inventory_location_id": inventoryLocationID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: InventoryLocationShape})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "inventory_locations"

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

// New binds the inventory location endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *InventoryLocations {
	return &InventoryLocations{requester: requester}
}
