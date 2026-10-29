package taxrules

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// TaxRules is the tax rule endpoints, bound to one client; code and zones read as joins but are unwritable.
type TaxRules struct {
	requester http.Requester
}

// ListTaxRules lists the tax rules — GET /admin/tax_rules; limit is capped at 100, so page on pagination.total.
func (t *TaxRules) ListTaxRules(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendGetRequest(ctx, TaxRulesPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: TaxRulePageShape})
}

// HeadTaxRules asks whether the collection is reachable — HEAD /admin/tax_rules.
func (t *TaxRules) HeadTaxRules(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendApiRequest(ctx, "HEAD", TaxRulesPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CreateTaxRule creates a tax rule — POST /admin/tax_rules, wrapped in tax_rule; only name and description are stored.
func (t *TaxRules) CreateTaxRule(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendPostRequest(ctx, TaxRulesPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: TaxRuleShape})
}

// GetTaxRule retrieves one tax rule — GET /admin/tax_rules/{tax_rule_id}.
func (t *TaxRules) GetTaxRule(ctx context.Context, taxRuleID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendGetRequest(ctx, TaxRulePath, &http.RequestOptions{
		PathParams:   map[string]any{"tax_rule_id": taxRuleID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: TaxRuleShape})
}

// HeadTaxRule asks whether the path answers — HEAD /admin/tax_rules/{tax_rule_id}; an unknown id is 200 too.
func (t *TaxRules) HeadTaxRule(ctx context.Context, taxRuleID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendApiRequest(ctx, "HEAD", TaxRulePath, &http.RequestOptions{
		PathParams:   map[string]any{"tax_rule_id": taxRuleID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// UpdateTaxRule merges into one tax rule — PUT /admin/tax_rules/{tax_rule_id}, wrapped in tax_rule; code and zones are fixed.
func (t *TaxRules) UpdateTaxRule(ctx context.Context, taxRuleID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendPutRequest(ctx, TaxRulePath, &http.RequestOptions{
		PathParams: map[string]any{"tax_rule_id": taxRuleID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: TaxRuleShape})
}

// DeleteTaxRule deletes one tax rule — DELETE /admin/tax_rules/{tax_rule_id}, 204 and hard.
func (t *TaxRules) DeleteTaxRule(ctx context.Context, taxRuleID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendDeleteRequest(ctx, TaxRulePath, &http.RequestOptions{
		PathParams:   map[string]any{"tax_rule_id": taxRuleID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "tax_rules"

// New binds the tax rule endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *TaxRules {
	return &TaxRules{requester: requester}
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
