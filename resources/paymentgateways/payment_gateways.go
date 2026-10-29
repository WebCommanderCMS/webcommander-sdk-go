package paymentgateways

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// PaymentGateways is the payment gateway endpoints, bound to one client.
type PaymentGateways struct {
	requester http.Requester
}

// ListPaymentGateways lists the gateways by ascending id — GET /admin/payment_gateways; limit 1 to 100 (20 by default), offset or page, search, exclude, promotional and added; other names are ignored.
func (c *PaymentGateways) ListPaymentGateways(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, PaymentGatewaysPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: PaymentGatewayPageShape})
}

// HeadPaymentGateways asks whether the collection is reachable — HEAD /admin/payment_gateways; no body.
func (c *PaymentGateways) HeadPaymentGateways(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendApiRequest(ctx, "HEAD", PaymentGatewaysPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CountPaymentGateways counts the gateways the listing's filters match — GET /admin/payment_gateways/count, a bare {"total": N}.
func (c *PaymentGateways) CountPaymentGateways(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, PaymentGatewaysCountPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: PaymentGatewayCountShape})
}

// HeadPaymentGatewaysCount asks whether the count is reachable — HEAD /admin/payment_gateways/count; no body.
func (c *PaymentGateways) HeadPaymentGatewaysCount(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendApiRequest(ctx, "HEAD", PaymentGatewaysCountPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// GetPaymentGateway retrieves one gateway with its settings — GET /admin/payment_gateways/{payment_gateway_id}; a stored secret reads ********.
func (c *PaymentGateways) GetPaymentGateway(ctx context.Context, paymentGatewayID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, PaymentGatewayPath, &http.RequestOptions{
		PathParams:   map[string]any{"payment_gateway_id": paymentGatewayID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: PaymentGatewayShape})
}

// HeadPaymentGateway asks whether one gateway is there — HEAD /admin/payment_gateways/{payment_gateway_id}; no body.
func (c *PaymentGateways) HeadPaymentGateway(ctx context.Context, paymentGatewayID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendApiRequest(ctx, "HEAD", PaymentGatewayPath, &http.RequestOptions{
		PathParams:   map[string]any{"payment_gateway_id": paymentGatewayID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// UpdatePaymentGateway writes one gateway — PUT /admin/payment_gateways/{payment_gateway_id}; a bare record merged like the PATCH, so the fields left out keep their value.
func (c *PaymentGateways) UpdatePaymentGateway(ctx context.Context, paymentGatewayID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendPutRequest(ctx, PaymentGatewayPath, &http.RequestOptions{
		PathParams: map[string]any{"payment_gateway_id": paymentGatewayID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: PaymentGatewayShape})
}

// PatchPaymentGateway changes only the fields sent — PATCH /admin/payment_gateways/{payment_gateway_id}; a bare record, and an empty body is 400.
func (c *PaymentGateways) PatchPaymentGateway(ctx context.Context, paymentGatewayID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendPatchRequest(ctx, PaymentGatewayPath, &http.RequestOptions{
		PathParams: map[string]any{"payment_gateway_id": paymentGatewayID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: PaymentGatewayShape})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "payment_gateways"

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

// New binds the payment gateway endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *PaymentGateways {
	return &PaymentGateways{requester: requester}
}
