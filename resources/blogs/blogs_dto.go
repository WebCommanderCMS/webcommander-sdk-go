package blogs

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// BlogCategoryReference is a category a blog is filed under; a write names it {"id": …}, a bare number is 500.
type BlogCategoryReference struct {
	ID           *int64
	Name         *string
	IsDisposable *bool
	IsInTrash    *bool
}

// BlogImage is the cover image as a read answers it, {} when there is none; Link's query string differs on every read.
type BlogImage struct {
	Title           *string
	AlternativeText *string
	FileName        *string
	Link            *string
}

// BlogImageInput is the cover image a write uploads; without Base64 the stored image is left as it was.
type BlogImageInput struct {
	FileName        string
	Base64          string `wire:"base64"`
	Title           string
	AlternativeText string
}

// BlogAuthor is created_by; {} on a row the store recorded no author for.
type BlogAuthor struct {
	ID   *int64
	Name *string
}

// BlogAudienceMember is one customer, by internal_id, or one customer group a restricted blog is shown to.
type BlogAudienceMember struct {
	ID   *int64
	Name *string
}

// BlogAudience is selected_customers once it names anyone; the store sends [] when it names no one, read as nil.
type BlogAudience struct {
	Customers []*BlogAudienceMember
	Groups    []*BlogAudienceMember
}

// BlogAudienceInput is the audience a write names: customer internal_id values and customer group ids.
type BlogAudienceInput struct {
	Customers []int64
	Groups    []int64
}

// BlogMetaTag is one meta_tag entry, the same two fields on a write and a read.
type BlogMetaTag struct {
	Name  *string
	Value *string
}

// BlogSeoConfigs is seo_configs, an object holding the meta_tag list.
type BlogSeoConfigs struct {
	MetaTag []*BlogMetaTag
}

// ResourceMeta is the one group_meta entry: what the collection is, whether it is writable and where.
type ResourceMeta struct {
	Kind     *string
	Writable *bool
	Href     *string
}

// GroupMeta is the group_meta block the listing, a row read and a PATCH answer carry.
type GroupMeta struct {
	Blogs *ResourceMeta
}

// Blog is one blog, the same eighteen fields on a listing row, a row read and every write's answer.
type Blog struct {
	ID                *int64
	UUID              *string
	Title             *string
	URL               *string
	Content           *string
	Date              *string
	IsPublished       *bool
	Visibility        *string
	VisibleTo         *string
	Categories        []*BlogCategoryReference
	CreatedOn         *string
	LastUpdatedOn     *string
	IsDisposable      *bool
	IsInTrash         *bool
	Image             *BlogImage
	CreatedBy         *BlogAuthor
	SelectedCustomers *BlogAudience
	SeoConfigs        *BlogSeoConfigs
	// GroupMeta is nil on a create's and a replace's answer.
	GroupMeta *GroupMeta
	// FieldMetadata is nil unless the call asked for it.
	FieldMetadata any
}

// EnvelopeKey is the key a row read and every write's answer wrap the record in; a listing row has none.
func (Blog) EnvelopeKey() string { return BlogKey }

// BlogInput is the writable form, for the create, the replace and the patch; it declares no url, which field_metadata marks read-only.
type BlogInput struct {
	Title       string
	Content     string
	Date        string
	IsPublished *bool
	Visibility  string
	VisibleTo   string
	Categories  []*BlogCategoryReference
	Image       *BlogImageInput
	SeoConfigs  *BlogSeoConfigs
	// SelectedCustomers names customers by internal_id; a customer_id there is dropped without a refusal.
	SelectedCustomers *BlogAudienceInput
}

// EnvelopeKey is the wrapper all three writes require.
func (BlogInput) EnvelopeKey() string { return BlogKey }

// ToMap renders the request body.
func (b BlogInput) ToMap() map[string]any { return dto.ToMap(b) }

// Pagination is the listing's ten-key pagination block; Limit is the size the store applied, capped at 20.
type Pagination struct {
	Total        *int64
	Limit        *int64
	Offset       *int64
	Count        *int64
	CurrentPage  *int64
	TotalPages   *int64
	HasNext      *bool
	HasPrevious  *bool
	PreviousPage *string
	NextPage     *string
}

// BlogPage is one page of the listing: the rows, group_meta, the paging block and any schema.
type BlogPage struct {
	Blogs     []*Blog
	GroupMeta *GroupMeta
	// Pagination is nil when the body carried none.
	Pagination    *Pagination
	FieldMetadata any
}

// BlogCount is the count sub-path's body, {"blogs": {"count": N}}.
type BlogCount struct {
	Count *int64
}

// EnvelopeKey is the plural key the count record arrives under.
func (BlogCount) EnvelopeKey() string { return BlogsKey }

// BlogComment is one comment, the same twelve fields in the listing and on its own path; PostTitle follows the blog's current title.
type BlogComment struct {
	ID            *int64
	Status        *string
	Name          *string
	Email         *string
	PostTitle     *string
	Content       *string
	Spam          *bool
	Likes         *int64
	TotalReply    *int64
	CreatedOn     *string
	LastUpdatedOn *string
	// Replies is the raw replies list.
	Replies []any
}

// BlogCommentList is the comments sub-path's body: the blog's id beside every comment it holds, unpaged; no envelope.
type BlogCommentList struct {
	ID       *int64
	Comments []*BlogComment
}

// BlogCommentEnvelope is one comment's body, from its own path and each comment write: the blog's id beside the comment.
type BlogCommentEnvelope struct {
	ID      *int64
	Comment *BlogComment
}

// BlogCommentDraft is the comment a create sends: Content required, at most 1000 characters; Name and Email may be left out.
type BlogCommentDraft struct {
	Content string
	Name    string
	Email   string
}

// BlogCommentInput is the comment create's body, {"blog": {"comment": {…}}}; without both wrappers it is 404.
type BlogCommentInput struct {
	Comment *BlogCommentDraft
}

// EnvelopeKey is the outer wrapper a comment create takes.
func (BlogCommentInput) EnvelopeKey() string { return BlogKey }

func (b BlogCommentInput) ToMap() map[string]any { return dto.ToMap(b) }

// BlogShape is the typed view of a body shaped {"blog": {…}}.
func BlogShape(decoded any) any { return dto.Read[Blog](decoded) }

// BlogPageShape is the typed view of the listing body.
func BlogPageShape(decoded any) any { return dto.Read[BlogPage](decoded) }

// BlogCountShape is the typed view of the count sub-path's body.
func BlogCountShape(decoded any) any { return dto.Read[BlogCount](decoded) }

// BlogCommentListShape is the typed view of the comments sub-path's body.
func BlogCommentListShape(decoded any) any { return dto.Read[BlogCommentList](decoded) }

// BlogCommentEnvelopeShape is the typed view of one comment's body.
func BlogCommentEnvelopeShape(decoded any) any { return dto.Read[BlogCommentEnvelope](decoded) }
