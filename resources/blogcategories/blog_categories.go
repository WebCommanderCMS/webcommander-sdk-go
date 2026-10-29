package blogcategories

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// BlogCategories is the blog category endpoints under /api/{api_version}/admin, bound to one client.
type BlogCategories struct {
	requester http.Requester
}

// ListBlogCategories lists the categories, ordered by name — GET /admin/blog_categories; limit (capped at 20), offset, page, sort, dir and field_metadata, any other name is 400.
func (b *BlogCategories) ListBlogCategories(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendGetRequest(ctx, BlogCategoriesPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: BlogCategoryPageShape})
}

// CreateBlogCategory creates a category — POST /admin/blog_categories, wrapped in blog_category, 201; name is required and unique, and a url another row holds is stored with -1 appended.
func (b *BlogCategories) CreateBlogCategory(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendPostRequest(ctx, BlogCategoriesPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: BlogCategoryShape})
}

// GetBlogCategory retrieves one category — GET /admin/blog_categories/{blog_category_id}.
func (b *BlogCategories) GetBlogCategory(ctx context.Context, blogCategoryID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendGetRequest(ctx, BlogCategoryPath, &http.RequestOptions{
		PathParams:   map[string]any{"blog_category_id": blogCategoryID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: BlogCategoryShape})
}

// UpdateBlogCategory updates a category — PUT /admin/blog_categories/{blog_category_id}, wrapped as the create is; name is required, an omitted description is blanked, and url and image are kept.
func (b *BlogCategories) UpdateBlogCategory(ctx context.Context, blogCategoryID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendPutRequest(ctx, BlogCategoryPath, &http.RequestOptions{
		PathParams: map[string]any{"blog_category_id": blogCategoryID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: BlogCategoryShape})
}

// DeleteBlogCategory deletes a category — DELETE /admin/blog_categories/{blog_category_id}; 204, and the blogs linked to it are kept.
func (b *BlogCategories) DeleteBlogCategory(ctx context.Context, blogCategoryID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendDeleteRequest(ctx, BlogCategoryPath, &http.RequestOptions{
		PathParams:   map[string]any{"blog_category_id": blogCategoryID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "blog_categories"

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

// New binds the blog category endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *BlogCategories { return &BlogCategories{requester: requester} }
