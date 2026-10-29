package blogcategories

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// BlogCategoryImage is a category's image, {} until one is uploaded; the query string on Link changes on every read.
type BlogCategoryImage struct {
	Title           *string
	AlternativeText *string
	FileName        *string
	Link            *string
}

// BlogCategorySeoConfig is one seo_configs entry; no write here sets the list, so a created category's is empty.
type BlogCategorySeoConfig struct {
	ID        *int64
	Value     *string
	Type      *string
	ConfigKey *string
}

// BlogCategory is one blog category, the same eleven fields on a listing row, a row path and every write's answer.
type BlogCategory struct {
	ID            *int64
	UUID          *string
	Name          *string
	URL           *string
	Description   *string
	IsDisposable  *bool
	IsInTrash     *bool
	Image         *BlogCategoryImage
	CreatedOn     *string
	LastUpdatedOn *string
	SeoConfigs    []*BlogCategorySeoConfig
}

// EnvelopeKey is the key a row path and every write's answer wrap the record in.
func (BlogCategory) EnvelopeKey() string { return BlogCategoryKey }

// BlogCategoryImageInput is an image upload; Title and AlternativeText are stored only beside a Base64.
type BlogCategoryImageInput struct {
	FileName        string
	Base64          string `wire:"base64"`
	Title           string
	AlternativeText string
}

// BlogCategoryBlog is a blog the category is linked to, by id.
type BlogCategoryBlog struct {
	ID int64
}

// BlogCategoryInput is the body both writes take; an update blanks an unset Description and keeps the url, image and blog links.
type BlogCategoryInput struct {
	Name        string
	URL         string
	Description string
	Image       *BlogCategoryImageInput
	Blogs       []*BlogCategoryBlog
}

// EnvelopeKey is the wrapper both writes require.
func (BlogCategoryInput) EnvelopeKey() string { return BlogCategoryKey }

// ToMap renders the request body.
func (b BlogCategoryInput) ToMap() map[string]any { return dto.ToMap(b) }

// Pagination is the listing's ten-key pagination block; Limit is the size the store applied.
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

// BlogCategoryPage is one page of the listing; field_metadata adds no key here.
type BlogCategoryPage struct {
	BlogCategories []*BlogCategory
	Pagination     *Pagination
}

// BlogCategoryShape is the typed view of a body shaped {"blog_category": {…}}.
func BlogCategoryShape(decoded any) any { return dto.Read[BlogCategory](decoded) }

// BlogCategoryPageShape is the typed view of the listing body.
func BlogCategoryPageShape(decoded any) any { return dto.Read[BlogCategoryPage](decoded) }
