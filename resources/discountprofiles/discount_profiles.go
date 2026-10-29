package discountprofiles

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// DiscountProfiles is the discount profile endpoints, bound to one client.
type DiscountProfiles struct {
	requester http.Requester
}

// ListDiscountProfiles lists the profiles, newest id first — GET /admin/discount_profiles.
func (d *DiscountProfiles) ListDiscountProfiles(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return d.requester.SendGetRequest(ctx, DiscountProfilesPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: DiscountProfilePageShape})
}

// CreateDiscountProfile creates a profile — POST /admin/discount_profiles as a bare record; name, discount_details_type and the details block it names are required.
func (d *DiscountProfiles) CreateDiscountProfile(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return d.requester.SendPostRequest(ctx, DiscountProfilesPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: DiscountProfileShape})
}

// GetDiscountProfile retrieves one profile — GET /admin/discount_profiles/{discount_profile_id}; field_metadata adds the field schema beside the row.
func (d *DiscountProfiles) GetDiscountProfile(ctx context.Context, discountProfileID any, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return d.requester.SendGetRequest(ctx, DiscountProfilePath, &http.RequestOptions{
		PathParams:  map[string]any{"discount_profile_id": discountProfileID},
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: DiscountProfileShape})
}

// UpdateDiscountProfile replaces one profile — PUT /admin/discount_profiles/{discount_profile_id}.
func (d *DiscountProfiles) UpdateDiscountProfile(ctx context.Context, discountProfileID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return d.requester.SendPutRequest(ctx, DiscountProfilePath, &http.RequestOptions{
		PathParams: map[string]any{"discount_profile_id": discountProfileID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: DiscountProfileShape})
}

// DeleteDiscountProfile deletes one profile and its coupons — DELETE /admin/discount_profiles/{discount_profile_id}; 204, no body.
func (d *DiscountProfiles) DeleteDiscountProfile(ctx context.Context, discountProfileID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return d.requester.SendDeleteRequest(ctx, DiscountProfilePath, &http.RequestOptions{
		PathParams:   map[string]any{"discount_profile_id": discountProfileID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CreateDiscountProfileCoupons generates codes on a profile that has is_apply_coupon_code — POST /admin/discount_profiles/{discount_profile_id}/create-coupons.
func (d *DiscountProfiles) CreateDiscountProfileCoupons(ctx context.Context, discountProfileID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return d.requester.SendPostRequest(ctx, DiscountProfileCreateCouponsPath, &http.RequestOptions{
		PathParams: map[string]any{"discount_profile_id": discountProfileID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: DiscountProfileCouponsShape})
}

// AssignDiscountCouponCustomer assigns a generated code to a customer — POST /admin/discount_coupons/assign-customer, {"code", "email"}; 201 and a verdict, so DTO() is nil.
func (d *DiscountProfiles) AssignDiscountCouponCustomer(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return d.requester.SendPostRequest(ctx, DiscountCouponsAssignCustomerPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "discount_profiles"

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

// New binds the discount profile endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *DiscountProfiles {
	return &DiscountProfiles{requester: requester}
}
