package taxzones

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// TaxZones is the tax zone endpoints, bound to one client; system_generated rows are the store's own.
type TaxZones struct {
	requester http.Requester
}

// ListTaxZones lists the tax zones — GET /admin/tax_zones; limit is clamped, so page on pagination.total.
func (t *TaxZones) ListTaxZones(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendGetRequest(ctx, TaxZonesPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: TaxZonePageShape})
}

// HeadTaxZones asks whether the collection is reachable — HEAD /admin/tax_zones.
func (t *TaxZones) HeadTaxZones(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendApiRequest(ctx, "HEAD", TaxZonesPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CreateTaxZone creates a tax zone — POST /admin/tax_zones, wrapped in tax_zone; countries and states take [{"id": 12}] records.
func (t *TaxZones) CreateTaxZone(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendPostRequest(ctx, TaxZonesPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: TaxZoneShape})
}

// GetTaxZone retrieves one tax zone — GET /admin/tax_zones/{tax_zone_id}.
func (t *TaxZones) GetTaxZone(ctx context.Context, taxZoneID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendGetRequest(ctx, TaxZonePath, &http.RequestOptions{
		PathParams:   map[string]any{"tax_zone_id": taxZoneID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: TaxZoneShape})
}

// HeadTaxZone answers 200 even for an unknown id — HEAD /admin/tax_zones/{tax_zone_id}, so it is no existence check.
func (t *TaxZones) HeadTaxZone(ctx context.Context, taxZoneID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendApiRequest(ctx, "HEAD", TaxZonePath, &http.RequestOptions{
		PathParams:   map[string]any{"tax_zone_id": taxZoneID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// UpdateTaxZone updates one tax zone — PUT /admin/tax_zones/{tax_zone_id}; name is required and the lists are replaced.
func (t *TaxZones) UpdateTaxZone(ctx context.Context, taxZoneID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendPutRequest(ctx, TaxZonePath, &http.RequestOptions{
		PathParams: map[string]any{"tax_zone_id": taxZoneID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: TaxZoneShape})
}

// DeleteTaxZone deletes one tax zone — DELETE /admin/tax_zones/{tax_zone_id}, 204 and hard.
func (t *TaxZones) DeleteTaxZone(ctx context.Context, taxZoneID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendDeleteRequest(ctx, TaxZonePath, &http.RequestOptions{
		PathParams:   map[string]any{"tax_zone_id": taxZoneID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "tax_zones"

// New binds the tax zone endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *TaxZones {
	return &TaxZones{requester: requester}
}

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
