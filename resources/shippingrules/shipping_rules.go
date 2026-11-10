package shippingrules

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// ShippingRules is the shipping rule endpoints, bound to one client.
type ShippingRules struct {
	requester http.Requester
}

// ListShippingRules lists the rules, newest first — GET /admin/shipping_rules; search (part of the name), shipping_class_id, zone_id, method (the method of the rule's rate) and handling_cost (with_handling or without_handling) filter, each rule listed once, sort is id, name, description, created_at, updated_at or transaction_no with order asc or desc, limit is 1 to 100 and defaults to 20, offset pages; an unknown value or an id that names no record is 400.
func (s *ShippingRules) ListShippingRules(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, ShippingRulesPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ShippingRuleListingShape})
}

// HeadShippingRules asks whether the collection is reachable — HEAD /admin/shipping_rules; no body.
func (s *ShippingRules) HeadShippingRules(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", ShippingRulesPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CountShippingRules is how many rules match — GET /admin/shipping_rules/count, a bare {"total": N}; the listing's filters narrow it.
func (s *ShippingRules) CountShippingRules(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, ShippingRulesCountPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ShippingRuleCountShape})
}

// HeadShippingRulesCount asks whether the count is reachable — HEAD /admin/shipping_rules/count; no body.
func (s *ShippingRules) HeadShippingRulesCount(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", ShippingRulesCountPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CreateShippingRule creates a rule — POST /admin/shipping_rules with {"shipping_rule": {…}}, answered 201 with a Location; name is required and unique without regard to case (409), a class, a rate and every zone must exist, zones without a repeat, and profile_id adds the rule to that profile. An unknown field, an id or an unwrapped body is 400. transactionNo is sent as X-Transaction-No and wins over a transaction_no in the body; empty sends none.
func (s *ShippingRules) CreateShippingRule(ctx context.Context, payload any, transactionNo string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	var headers map[string]string
	if transactionNo != "" {
		headers = map[string]string{ShippingRulesTransactionNoHeader: transactionNo}
	}
	return s.requester.SendPostRequest(ctx, ShippingRulesPath, &http.RequestOptions{
		JSONBody: payload, Headers: headers, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ShippingRuleShape})
}

// GetShippingRule retrieves one rule by its id — GET /admin/shipping_rules/{shipping_rule_id}; an id that is not plain digits, or unknown, is 404.
func (s *ShippingRules) GetShippingRule(ctx context.Context, shippingRuleID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, ShippingRulePath, &http.RequestOptions{
		PathParams:   map[string]any{"shipping_rule_id": shippingRuleID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ShippingRuleShape})
}

// HeadShippingRule asks whether a rule exists — HEAD /admin/shipping_rules/{shipping_rule_id}; no body, 404 for an unknown id.
func (s *ShippingRules) HeadShippingRule(ctx context.Context, shippingRuleID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", ShippingRulePath, &http.RequestOptions{
		PathParams:   map[string]any{"shipping_rule_id": shippingRuleID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// UpdateShippingRule replaces one rule — PUT /admin/shipping_rules/{shipping_rule_id}; name, description, shipping_class, shipping_rate and zones are all required, the first three may be null and zones []; the rate and the zones themselves stay.
func (s *ShippingRules) UpdateShippingRule(ctx context.Context, shippingRuleID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPutRequest(ctx, ShippingRulePath, &http.RequestOptions{
		PathParams: map[string]any{"shipping_rule_id": shippingRuleID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ShippingRuleShape})
}

// PatchShippingRule changes only the fields sent — PATCH /admin/shipping_rules/{shipping_rule_id}; zones is the full list, so [] takes every zone off, and a null shipping_rate or shipping_class takes the rate or the class off; the records themselves stay.
func (s *ShippingRules) PatchShippingRule(ctx context.Context, shippingRuleID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, ShippingRulePath, &http.RequestOptions{
		PathParams: map[string]any{"shipping_rule_id": shippingRuleID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ShippingRuleShape})
}

// DeleteShippingRule deletes one rule — DELETE /admin/shipping_rules/{shipping_rule_id}, 204; its rate and zones stay, and a rule on a profile is 409 until it is taken off.
func (s *ShippingRules) DeleteShippingRule(ctx context.Context, shippingRuleID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendDeleteRequest(ctx, ShippingRulePath, &http.RequestOptions{
		PathParams:   map[string]any{"shipping_rule_id": shippingRuleID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "shipping_rules"

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

// New binds the shipping rule endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *ShippingRules {
	return &ShippingRules{requester: requester}
}
