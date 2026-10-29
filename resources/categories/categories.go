package categories

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// Categories is the category endpoints, under /api/{api_version}/admin, bound to one client.
type Categories struct {
	requester http.Requester
}

// ListCategories lists the categories, newest update first — GET /admin/categories; limit caps at 20, unknown names are ignored.
func (c *Categories) ListCategories(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, CategoriesPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CategoryPageShape})
}

// CountCategories counts the categories — GET /admin/categories/count; name, searchText and parent narrow it, trashed rows are not counted.
func (c *Categories) CountCategories(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, CategoriesCountPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CategoryCountShape})
}

// CheckCategoryName asks whether a name is free — GET /admin/categories/check-name; exact and case-insensitive over every level and the trash.
func (c *Categories) CheckCategoryName(ctx context.Context, name string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, CategoriesCheckNamePath, &http.RequestOptions{
		QueryParams: map[string]any{"name": name}, RaiseOnError: raiseOnError,
		ResponseFormat: format})
}

// CreateCategory creates a category — POST /admin/categories, wrapped in category; a top-level row's name, trash included, is 409. transactionNo is sent as X-Transaction-No and files the category under that number; empty sends none.
func (c *Categories) CreateCategory(ctx context.Context, payload any, transactionNo string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	var headers map[string]string
	if transactionNo != "" {
		headers = map[string]string{CategoriesTransactionNoHeader: transactionNo}
	}
	return c.requester.SendPostRequest(ctx, CategoriesPath, &http.RequestOptions{
		JSONBody: payload, Headers: headers, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CategoryShape})
}

// GetCategory retrieves one category — GET /admin/categories/{category_id}; a trashed one still reads 200 with in_trash true.
func (c *Categories) GetCategory(ctx context.Context, categoryID any, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, CategoryPath, &http.RequestOptions{
		PathParams: map[string]any{"category_id": categoryID}, QueryParams: queryParams,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CategoryShape})
}

// UpdateCategory replaces one category — PUT /admin/categories/{category_id}; omitted fields reset but name, sku, url, images and five tax_shipping fields.
func (c *Categories) UpdateCategory(ctx context.Context, categoryID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendPutRequest(ctx, CategoryPath, &http.RequestOptions{
		PathParams: map[string]any{"category_id": categoryID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CategoryShape})
}

// PatchCategory merges fields into one category — PATCH /admin/categories/{category_id}, across as many sections as the body names.
func (c *Categories) PatchCategory(ctx context.Context, categoryID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendPatchRequest(ctx, CategoryPath, &http.RequestOptions{
		PathParams: map[string]any{"category_id": categoryID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CategoryShape})
}

// DeleteCategory moves one category to the trash — DELETE /admin/categories/{category_id}; 204, and a second delete is 404.
func (c *Categories) DeleteCategory(ctx context.Context, categoryID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendDeleteRequest(ctx, CategoryPath, &http.RequestOptions{
		PathParams:   map[string]any{"category_id": categoryID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// GetCategorySection reads one section — GET /admin/categories/{category_id}/{section}, answered as {"category": {<section>: {…}}}.
func (c *Categories) GetCategorySection(ctx context.Context, categoryID any, section string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, CategorySectionPath, &http.RequestOptions{
		PathParams:   map[string]any{"category_id": categoryID, "section": section},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CategoryShape})
}

// PatchCategorySection merges into one section — PATCH /admin/categories/{category_id}/{section}, wrapped as {"category": {<section>: {…}}}.
func (c *Categories) PatchCategorySection(ctx context.Context, categoryID any, section string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendPatchRequest(ctx, CategorySectionPath, &http.RequestOptions{
		PathParams: map[string]any{"category_id": categoryID, "section": section}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CategoryShape})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "categories"

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

// New binds the category endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *Categories {
	return &Categories{requester: requester}
}
