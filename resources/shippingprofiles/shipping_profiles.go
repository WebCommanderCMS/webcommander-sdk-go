package shippingprofiles

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// ShippingProfiles is the shipping profile endpoints, bound to one client.
type ShippingProfiles struct {
	requester http.Requester
}

// ListShippingProfiles lists the profiles, newest first — GET /admin/shipping_profiles; search (part of the name) and default (true or false) filter, sort is id, name, description, rule_preference, created_at, updated_at or transaction_no with order asc or desc, limit is 1 to 100 and defaults to 20, offset pages; an unknown value is 400 and an unknown name is ignored.
func (s *ShippingProfiles) ListShippingProfiles(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, ShippingProfilesPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ShippingProfileListingShape})
}

// HeadShippingProfiles asks whether the collection is reachable — HEAD /admin/shipping_profiles; no body.
func (s *ShippingProfiles) HeadShippingProfiles(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", ShippingProfilesPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CountShippingProfiles is how many profiles match — GET /admin/shipping_profiles/count, a bare {"total": N}; search and default narrow it.
func (s *ShippingProfiles) CountShippingProfiles(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, ShippingProfilesCountPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ShippingProfileCountShape})
}

// HeadShippingProfilesCount asks whether the count is reachable — HEAD /admin/shipping_profiles/count; no body.
func (s *ShippingProfiles) HeadShippingProfilesCount(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", ShippingProfilesCountPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CreateShippingProfile creates a profile — POST /admin/shipping_profiles with {"shipping_profile": {…}}, answered 201 with a Location; name is required and unique without regard to case (409), shipping_rules lists rule ids once each, and an unknown field, an id or an unwrapped body is 400. transactionNo is sent as X-Transaction-No and wins over a transaction_no in the body; empty sends none.
func (s *ShippingProfiles) CreateShippingProfile(ctx context.Context, payload any, transactionNo string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	var headers map[string]string
	if transactionNo != "" {
		headers = map[string]string{ShippingProfilesTransactionNoHeader: transactionNo}
	}
	return s.requester.SendPostRequest(ctx, ShippingProfilesPath, &http.RequestOptions{
		JSONBody: payload, Headers: headers, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ShippingProfileShape})
}

// CreateShippingProfileWizard creates a profile with its rules and their rates in one call — POST /admin/shipping_profiles/wizard with {"shipping_profile": {…}}; a new profile name is 201, and a name that exists adds the rules to that profile and is 200. Every rule's rate, class and zone is checked first and one bad value saves nothing; a zone given as {"name": …} is created. The wizard takes no transaction number.
func (s *ShippingProfiles) CreateShippingProfileWizard(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPostRequest(ctx, ShippingProfilesWizardPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ShippingProfileShape})
}

// GetShippingProfile retrieves one profile by its id — GET /admin/shipping_profiles/{shipping_profile_id}; an id that is not plain digits, or unknown, is 404.
func (s *ShippingProfiles) GetShippingProfile(ctx context.Context, shippingProfileID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, ShippingProfilePath, &http.RequestOptions{
		PathParams:   map[string]any{"shipping_profile_id": shippingProfileID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ShippingProfileShape})
}

// HeadShippingProfile asks whether a profile exists — HEAD /admin/shipping_profiles/{shipping_profile_id}; no body, 404 for an unknown id.
func (s *ShippingProfiles) HeadShippingProfile(ctx context.Context, shippingProfileID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", ShippingProfilePath, &http.RequestOptions{
		PathParams:   map[string]any{"shipping_profile_id": shippingProfileID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// UpdateShippingProfile replaces one profile — PUT /admin/shipping_profiles/{shipping_profile_id}; name, description, rule_preference and shipping_rules are all required (description may be null), and the rules not listed are taken off the profile and kept.
func (s *ShippingProfiles) UpdateShippingProfile(ctx context.Context, shippingProfileID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPutRequest(ctx, ShippingProfilePath, &http.RequestOptions{
		PathParams: map[string]any{"shipping_profile_id": shippingProfileID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ShippingProfileShape})
}

// PatchShippingProfile changes only the fields sent — PATCH /admin/shipping_profiles/{shipping_profile_id}; shipping_rules is the full list, so [] takes every rule off, and default true moves the store default to this profile; the default cannot be switched off directly (409).
func (s *ShippingProfiles) PatchShippingProfile(ctx context.Context, shippingProfileID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, ShippingProfilePath, &http.RequestOptions{
		PathParams: map[string]any{"shipping_profile_id": shippingProfileID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ShippingProfileShape})
}

// DeleteShippingProfile deletes one profile — DELETE /admin/shipping_profiles/{shipping_profile_id}, 204; its rules and zones stay, and the store default is 409 until another profile holds the flag.
func (s *ShippingProfiles) DeleteShippingProfile(ctx context.Context, shippingProfileID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendDeleteRequest(ctx, ShippingProfilePath, &http.RequestOptions{
		PathParams:   map[string]any{"shipping_profile_id": shippingProfileID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "shipping_profiles"

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

// New binds the shipping profile endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *ShippingProfiles {
	return &ShippingProfiles{requester: requester}
}
