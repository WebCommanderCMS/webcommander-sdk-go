package blogs

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// Blogs is the blog and blog comment endpoints, bound to one client.
type Blogs struct {
	requester http.Requester
}

// ListBlogs lists the blogs, latest date first — GET /admin/blogs; limit is capped at 20, and an unknown name is 400.
func (b *Blogs) ListBlogs(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendGetRequest(ctx, BlogsPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: BlogPageShape})
}

// CountBlogs reads how many blogs are live — GET /admin/blogs/count; narrowed by the listing's filters, not its paging.
func (b *Blogs) CountBlogs(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendGetRequest(ctx, BlogsCountPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: BlogCountShape})
}

// CreateBlog creates a blog — POST /admin/blogs, answered 201; title is required and unique among live and trashed blogs.
func (b *Blogs) CreateBlog(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendPostRequest(ctx, BlogsPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: BlogShape})
}

// GetBlog retrieves one blog — GET /admin/blogs/{blog_id}; field_metadata=true adds the schema, and a trashed blog is 404.
func (b *Blogs) GetBlog(ctx context.Context, blogID any, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendGetRequest(ctx, BlogPath, &http.RequestOptions{
		PathParams:  map[string]any{"blog_id": blogID},
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: BlogShape})
}

// UpdateBlog replaces one blog — PUT /admin/blogs/{blog_id}; title and visibility are required, and every other field left out is reset but the image and url.
func (b *Blogs) UpdateBlog(ctx context.Context, blogID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendPutRequest(ctx, BlogPath, &http.RequestOptions{
		PathParams: map[string]any{"blog_id": blogID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: BlogShape})
}

// PatchBlog changes only the fields sent — PATCH /admin/blogs/{blog_id}; the answer carries group_meta beside the row.
func (b *Blogs) PatchBlog(ctx context.Context, blogID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendPatchRequest(ctx, BlogPath, &http.RequestOptions{
		PathParams: map[string]any{"blog_id": blogID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: BlogShape})
}

// DeleteBlog moves one blog to the trash — DELETE /admin/blogs/{blog_id}, 204; its title stays taken and its comments readable.
func (b *Blogs) DeleteBlog(ctx context.Context, blogID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendDeleteRequest(ctx, BlogPath, &http.RequestOptions{
		PathParams:   map[string]any{"blog_id": blogID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ListBlogComments reads every comment on one blog, oldest first and unpaged — GET /admin/blogs/{blog_id}/comments; no query parameter has an effect.
func (b *Blogs) ListBlogComments(ctx context.Context, blogID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendGetRequest(ctx, BlogCommentsPath, &http.RequestOptions{
		PathParams:   map[string]any{"blog_id": blogID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: BlogCommentListShape})
}

// CreateBlogComment adds a pending comment — POST /admin/blogs/{blog_id}/comments, 201; the body is {"blog": {"comment": {…}}}.
func (b *Blogs) CreateBlogComment(ctx context.Context, blogID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendPostRequest(ctx, BlogCommentsPath, &http.RequestOptions{
		PathParams: map[string]any{"blog_id": blogID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: BlogCommentEnvelopeShape})
}

// GetBlogComment retrieves one comment — GET /admin/blogs/{blog_id}/comments/{blog_comment_id}; the top-level id is the blog's.
func (b *Blogs) GetBlogComment(ctx context.Context, blogID any, blogCommentID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendGetRequest(ctx, BlogCommentPath, &http.RequestOptions{
		PathParams:   map[string]any{"blog_id": blogID, "blog_comment_id": blogCommentID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: BlogCommentEnvelopeShape})
}

// ApproveBlogComment approves one comment — POST .../comments/{blog_comment_id}/approve; no body, 201, and spam is left as it was.
func (b *Blogs) ApproveBlogComment(ctx context.Context, blogID any, blogCommentID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendPostRequest(ctx, BlogCommentApprovePath, &http.RequestOptions{
		PathParams:   map[string]any{"blog_id": blogID, "blog_comment_id": blogCommentID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: BlogCommentEnvelopeShape})
}

// MarkBlogCommentSpam marks one comment as spam, setting status and spam — POST .../comments/{blog_comment_id}/mark-spam; no body, 201.
func (b *Blogs) MarkBlogCommentSpam(ctx context.Context, blogID any, blogCommentID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendPostRequest(ctx, BlogCommentMarkSpamPath, &http.RequestOptions{
		PathParams:   map[string]any{"blog_id": blogID, "blog_comment_id": blogCommentID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: BlogCommentEnvelopeShape})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "blogs"

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

// New binds the blog endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *Blogs {
	return &Blogs{requester: requester}
}
