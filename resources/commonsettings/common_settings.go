package commonsettings

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// CommonSettingsResource is the store-wide settings and the registration form, bound to one client.
type CommonSettingsResource struct {
	requester http.Requester
}

// GetSettings reads the store's settings — GET /settings, {"settings": {group: [{"config_key", "value"}]}}; type names one group (unknown is 404), config_key one key.
// A comma list of keys answers {"settings": {}}, other names (paging included) are 400, and without type every group answers, about 2 MB, in 30 seconds or more.
func (c *CommonSettingsResource) GetSettings(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, CommonSettingsPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CommonSettingsShape})
}

// HeadSettings asks whether the settings are reachable — HEAD /settings; no body. It runs the read's query, so pass a type: an unknown one is 404.
func (c *CommonSettingsResource) HeadSettings(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendApiRequest(ctx, "HEAD", CommonSettingsPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError, ResponseFormat: format})
}

// UpdateSettings saves settings — PUT /settings with {"settings": {group: …}}, 200 and every group of the store. Only the keys sent change, and one bad entry refuses the whole request with 400.
// A write took 5 to 53 seconds, past the default 30 of TimeoutSeconds at the slow end: build the client with a longer one.
func (c *CommonSettingsResource) UpdateSettings(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendPutRequest(ctx, CommonSettingsPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CommonSettingsShape})
}

// SaveSettings saves settings — POST /settings, the PUT's body and behaviour, answered 201. A POST is not resent after a timeout, so read the keys back before retrying.
func (c *CommonSettingsResource) SaveSettings(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendPostRequest(ctx, CommonSettingsPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CommonSettingsShape})
}

// GetCreateCustomerFields reads the customer registration form — GET /create-customer-fields, the active fields in order and every field's fields_configs; any query name is 400.
func (c *CommonSettingsResource) GetCreateCustomerFields(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, CreateCustomerFieldsPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CustomerFieldsShape})
}

// HeadCreateCustomerFields asks whether the registration form is there — HEAD /create-customer-fields; no body.
func (c *CommonSettingsResource) HeadCreateCustomerFields(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendApiRequest(ctx, "HEAD", CreateCustomerFieldsPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "common_settings"

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

// New binds the common settings endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *CommonSettingsResource {
	return &CommonSettingsResource{requester: requester}
}
