package email

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// Email is the email endpoint, bound to one client.
type Email struct {
	requester http.Requester
}

// SendEmail sends an email — POST /admin/email/send with a bare camelCase body, answered 200 once every message is sent: recipients and subject, with htmlContent or a template identifier and its data. pluginUUID is sent as the plugin_uuid header a template send needs and switches the call to plugin auth, so an admin token is 401; empty sends none. Past 500 addresses a day per sender the store answers 429 until 00:00 UTC; a body over 1048576 bytes is 413. A send the mail server fails is 502, status failed, each address in Rejected.
func (e *Email) SendEmail(ctx context.Context, payload any, pluginUUID string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	var headers map[string]string
	if pluginUUID != "" {
		headers = map[string]string{EmailPluginUUIDHeader: pluginUUID}
	}
	return e.requester.SendPostRequest(ctx, EmailSendPath, &http.RequestOptions{
		JSONBody: payload, Headers: headers, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: EmailSendResultShape})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "email"

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

// New binds the email endpoint to any requester (the WebCommanderClient).
func New(requester http.Requester) *Email {
	return &Email{requester: requester}
}
