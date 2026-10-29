package trash

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// Trash is the trash endpoints, bound to one client.
type Trash struct {
	requester http.Requester
}

// ListTrash lists trashed rows, newest change first — GET /admin/trash; type (one or a comma list), search, sort and order, limit (above 100 is 400), offset, page, removed_from and removed_to; any other name is 400.
func (t *Trash) ListTrash(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendGetRequest(ctx, TrashPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: TrashPageShape})
}

// HeadTrash asks whether the listing is reachable — HEAD /admin/trash; no body. It runs the listing's query, so pass a type as on the GET.
func (t *Trash) HeadTrash(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendApiRequest(ctx, "HEAD", TrashPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ListTrashTypes lists the types the trash holds — GET /admin/trash/types, snake_case names such as category, customer and blog_post; any query name is 400.
func (t *Trash) ListTrashTypes(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendGetRequest(ctx, TrashTypesPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: TrashTypeListShape})
}

// HeadTrashTypes asks whether the type list is reachable — HEAD /admin/trash/types; no body.
func (t *Trash) HeadTrashTypes(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendApiRequest(ctx, "HEAD", TrashTypesPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// RestoreTrashItem restores one trashed row — POST /admin/trash/{trash_type}/{trash_item_id}/restore with no body, answered 200 with the row under item; a live row is 409, an id the type does not hold 404.
func (t *Trash) RestoreTrashItem(ctx context.Context, trashType string, trashItemID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendPostRequest(ctx, TrashItemRestorePath, &http.RequestOptions{
		PathParams:   map[string]any{"trash_type": trashType, "trash_item_id": trashItemID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: TrashItemShape})
}

// DeleteTrashItem deletes one trashed row for good — DELETE /admin/trash/{trash_type}/{trash_item_id}; 204. A live row is 409; once purged, a second delete is 404.
func (t *Trash) DeleteTrashItem(ctx context.Context, trashType string, trashItemID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendDeleteRequest(ctx, TrashItemPath, &http.RequestOptions{
		PathParams:   map[string]any{"trash_type": trashType, "trash_item_id": trashItemID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// BulkRestoreTrash restores up to 100 rows — POST /admin/trash/bulk-restore with {"items": [{"type", "id"}]}; 200 when every item is restored, 207 with the same body when any fails.
func (t *Trash) BulkRestoreTrash(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendPostRequest(ctx, TrashBulkRestorePath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: TrashBulkRestoreResultShape})
}

// BulkDeleteTrash deletes up to 100 trashed rows for good — POST /admin/trash/bulk-delete, the bulk restore's body; 200 when every item is deleted, 207 with the same body when any fails.
func (t *Trash) BulkDeleteTrash(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return t.requester.SendPostRequest(ctx, TrashBulkDeletePath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: TrashBulkDeleteResultShape})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "trash"

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

// New binds the trash endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *Trash {
	return &Trash{requester: requester}
}
