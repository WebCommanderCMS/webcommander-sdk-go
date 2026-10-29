package customers

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// AdminCustomers is the admin customer endpoints — /api/{api_version}/admin/customers* — bound to one client.
type AdminCustomers struct {
	requester http.Requester
}

// GetCustomers lists customers — GET /admin/customers; limit caps at 20 and parameter names are case-sensitive.
func (a *AdminCustomers) GetCustomers(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendGetRequest(ctx, AdminCustomersPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: PageOf[Customer]})
}

// GetCustomer fetches one customer by id — GET /admin/customers/{customer_id}; the slot takes customer_id, not internal_id.
func (a *AdminCustomers) GetCustomer(ctx context.Context, customerID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendGetRequest(ctx, AdminCustomerPath, &http.RequestOptions{
		PathParams: map[string]any{"customer_id": customerID}, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CustomerShape})
}

// GetCustomerCount counts matching customers — GET /admin/customers/count; Data is the one key total.
func (a *AdminCustomers) GetCustomerCount(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendGetRequest(ctx, AdminCustomersCountPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format})
}

// GetCustomerEmailExists checks an email — GET /admin/customers/check-email; a trashed customer still exists.
func (a *AdminCustomers) GetCustomerEmailExists(ctx context.Context, email string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if strings.TrimSpace(email) == "" {
		return nil, wcerrors.NewUsageError("email is required and identifies the lookup, so a blank one is " +
			"refused here rather than sent as a different question")
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendGetRequest(ctx, AdminCustomersCheckEmailPath, &http.RequestOptions{
		QueryParams: map[string]any{"email": email}, RaiseOnError: raiseOnError,
		ResponseFormat: format})
}

// GetCustomerExport downloads all customers as XLSX bytes — GET /admin/customers/export; queryParams toggle columns, not rows.
func (a *AdminCustomers) GetCustomerExport(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendGetRequest(ctx, AdminCustomersExportPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CreateCustomer creates a customer — POST /admin/customers; a duplicate email, trashed ones included, is 409. transactionNo is sent as X-Transaction-No and files the customer under that number; empty sends none.
func (a *AdminCustomers) CreateCustomer(ctx context.Context, payload any, transactionNo string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	var headers map[string]string
	if transactionNo != "" {
		headers = map[string]string{AdminCustomersTransactionNoHeader: transactionNo}
	}
	return a.requester.SendPostRequest(ctx, AdminCustomersPath, &http.RequestOptions{
		JSONBody: payload, Headers: headers, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CustomerShape})
}

// UpdateCustomer updates a customer — PUT /admin/customers/{customer_id}.
func (a *AdminCustomers) UpdateCustomer(ctx context.Context, customerID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendPutRequest(ctx, AdminCustomerPath, &http.RequestOptions{
		PathParams: map[string]any{"customer_id": customerID}, JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CustomerShape})
}

// DeleteCustomer deletes a customer — DELETE /admin/customers/{customer_id}; 204 with no body, so DTO() is nil.
func (a *AdminCustomers) DeleteCustomer(ctx context.Context, customerID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendDeleteRequest(ctx, AdminCustomerPath, &http.RequestOptions{
		PathParams: map[string]any{"customer_id": customerID}, RaiseOnError: raiseOnError,
		ResponseFormat: format})
}

// CreateCustomerAddresses adds addresses — POST /admin/customers/{customer_id}/addresses; type is PRIMARY, SHIPPING or BILLING.
func (a *AdminCustomers) CreateCustomerAddresses(ctx context.Context, customerID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendPostRequest(ctx, AdminCustomerAddressesPath, &http.RequestOptions{
		PathParams: map[string]any{"customer_id": customerID}, JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CustomerShape})
}

// ResetCustomerPassword sets a password without the current one — POST /admin/customers/{customer_id}/reset-password.
func (a *AdminCustomers) ResetCustomerPassword(ctx context.Context, customerID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendPostRequest(ctx, AdminCustomerResetPasswordPath, &http.RequestOptions{
		PathParams: map[string]any{"customer_id": customerID}, JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format})
}

// ChangeCustomerPassword changes a password — POST /admin/customers/{customer_id}/change-password, email and old_password included.
func (a *AdminCustomers) ChangeCustomerPassword(ctx context.Context, customerID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendPostRequest(ctx, AdminCustomerChangePasswordPath, &http.RequestOptions{
		PathParams: map[string]any{"customer_id": customerID}, JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format})
}

// BulkResetCustomerPasswords queues resets — POST /admin/customers/bulk-reset-password, body {"customer_ids": […]}.
func (a *AdminCustomers) BulkResetCustomerPasswords(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendPostRequest(ctx, AdminCustomersBulkResetPasswordPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format})
}

// ImportCustomers uploads a spreadsheet — POST /admin/customers/import; the part is named file and its extension picks the parser.
func (a *AdminCustomers) ImportCustomers(ctx context.Context, files any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	normalised, err := normaliseFiles(files)
	if err != nil {
		return nil, err
	}
	if err := requireFiles(files, normalised); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendPostRequest(ctx, AdminCustomersImportPath, &http.RequestOptions{
		Files: normalised, RaiseOnError: raiseOnError,
		ResponseFormat: format})
}

// GetCustomerStoreCreditAdjustments lists store-credit movements, newest first — GET /admin/customers/{customer_id}/store-credit/adjustments.
func (a *AdminCustomers) GetCustomerStoreCreditAdjustments(ctx context.Context, customerID any, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendGetRequest(ctx, AdminCustomerStoreCreditAdjustmentsPath, &http.RequestOptions{
		PathParams: map[string]any{"customer_id": customerID}, QueryParams: queryParams,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: PageOf[StoreCreditAdjustment]})
}

// CreateCustomerStoreCreditAdjustment moves a customer's store credit — POST /admin/customers/{customer_id}/store-credit/adjustments.
func (a *AdminCustomers) CreateCustomerStoreCreditAdjustment(ctx context.Context, customerID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return a.requester.SendPostRequest(ctx, AdminCustomerStoreCreditAdjustmentsPath, &http.RequestOptions{
		PathParams: map[string]any{"customer_id": customerID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: StoreCreditAdjustmentShape})
}

// normaliseFiles accepts map[string]any or [][2]any and returns the map the transport encodes.
func normaliseFiles(files any) (map[string]any, error) {
	switch shape := files.(type) {
	case nil:
		return nil, nil
	case map[string]any:
		return shape, nil
	case [][2]any:
		out := make(map[string]any, len(shape))
		for _, pair := range shape {
			name, ok := pair[0].(string)
			if !ok {
				return nil, wcerrors.NewUsageError("files: a pair's first element must be the field " +
					"name as a string")
			}
			if _, clash := out[name]; clash {
				return nil, wcerrors.NewUsageError("files: the field " + name + " appears twice; a map " +
					"cannot carry two values under one name, so the duplicate would be lost")
			}
			out[name] = pair[1]
		}
		return out, nil
	default:
		return nil, wcerrors.NewUsageError("files: expected map[string]any or [][2]any " +
			"(field, value) pairs")
	}
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "admin_customers"

// oneResponseFormat reads the optional trailing responseFormat, refusing more than one.
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

// requirePayload refuses a nil body before any I/O; the client's validator would let one through bodyless.
func requirePayload(payload any) error {
	if payload == nil || isNilValue(payload) {
		return wcerrors.NewUsageError("payload is required: this endpoint takes a body, and nil would " +
			"send the request with no body at all")
	}
	return nil
}

// isNilValue reports whether a value is a typed nil inside a non-nil any; a zero-valued struct is not.
func isNilValue(value any) bool {
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Map, reflect.Slice, reflect.Pointer, reflect.Interface, reflect.Func:
		return reflected.IsNil()
	}
	return false
}

// requireFiles refuses an upload carrying nothing — no files, or an empty container — before any I/O.
func requireFiles(files any, normalised map[string]any) error {
	if files == nil {
		return wcerrors.NewUsageError("files is required")
	}
	if len(normalised) == 0 {
		return wcerrors.NewUsageError(fmt.Sprintf("files must name the file to upload, got an empty %T", files))
	}
	return nil
}

// New binds the admin customer endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *AdminCustomers { return &AdminCustomers{requester: requester} }
