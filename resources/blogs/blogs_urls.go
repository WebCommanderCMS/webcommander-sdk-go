package blogs

// The path templates; the client fills {api_version}.
const (
	BlogsPath = "/api/{api_version}/admin/blogs"

	BlogsCountPath = "/api/{api_version}/admin/blogs/count"

	BlogPath = "/api/{api_version}/admin/blogs/{blog_id}"

	BlogCommentsPath = "/api/{api_version}/admin/blogs/{blog_id}/comments"

	BlogCommentPath = "/api/{api_version}/admin/blogs/{blog_id}/comments/{blog_comment_id}"

	BlogCommentApprovePath = "/api/{api_version}/admin/blogs/{blog_id}/comments/{blog_comment_id}/approve"

	BlogCommentMarkSpamPath = "/api/{api_version}/admin/blogs/{blog_id}/comments/{blog_comment_id}/mark-spam"
)

// BlogsKey is the key the listing's rows arrive under; the count sub-path reuses it for a record.
const BlogsKey = "blogs"

// BlogKey is the key one row arrives under, and the wrapper every write takes, a comment's included.
const BlogKey = "blog"

// BlogsCountKey is the one key the count sub-path's record carries, inside BlogsKey.
const BlogsCountKey = "count"

// BlogCommentsKey is the key a blog's comments arrive under, beside the blog's own id.
const BlogCommentsKey = "comments"

// BlogCommentKey is the key one comment arrives under, beside the blog's own id.
const BlogCommentKey = "comment"

// BlogsPageSize is the default page size, and the ceiling on a limit that is sent.
const BlogsPageSize = 20

var blogsPaginationKeys = [10]string{
	"count", "current_page", "has_next", "has_previous", "limit",
	"next_page", "offset", "previous_page", "total", "total_pages",
}

// BlogsPaginationKeys returns a fresh copy of the ten keys the listing's pagination block carries.
func BlogsPaginationKeys() []string {
	return append(make([]string, 0, len(blogsPaginationKeys)), blogsPaginationKeys[:]...)
}
