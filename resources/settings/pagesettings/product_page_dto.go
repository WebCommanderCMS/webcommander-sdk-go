package pagesettings

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// ActionButtons is the actions_buttons block: which buttons a product page carries, all booleans.
type ActionButtons struct {
	ShowProductActionButtons   *bool
	DisplayAddToCartButton     *bool
	DisplayBuyNowButton        *bool
	DisplayOnHover             *bool
	DisplayAddToCompareButton  *bool
	DisplayAddToWishlistButton *bool
}

// ProductDetails is the details block: which parts of the product's own description are shown.
type ProductDetails struct {
	DisplayDescription              *bool
	DisplaySalePercentage           *bool
	UseStrikeThroughOnPreviousPrice *bool
	DisplayTaxWithExpectToPayPrice  *bool
}

// SocialSharing is the social_sharing block: share buttons, the networks they cover and their title.
type SocialSharing struct {
	EnableShareOnSocialMedia *bool
	Facebook                 *bool
	Twitter                  *bool
	Pinterest                *bool
	Messenger                *bool
	Linkedin                 *bool
	Email                    *bool
	Snapchat                 *bool
	Telegram                 *bool
	Reddit                   *bool
	Instagram                *bool
	Youtube                  *bool
	Tiktok                   *bool
	Threads                  *bool
	Whatsapp                 *bool
	Bluesky                  *bool
	Mastodon                 *bool
	CopyLink                 *bool
	DisplayTitle             *bool
	Title                    *string
}

// SetProductProperties is the set_product_properties block: which physical properties a page lists.
type SetProductProperties struct {
	DisplayHeight            *bool
	DisplayWidth             *bool
	DisplayLength            *bool
	DisplayWeight            *bool
	DisplayAddToGiftRegistry *bool
}

// Configuration is the configuration section — what the product page itself shows, in four blocks.
type Configuration struct {
	ActionsButtons       *ActionButtons
	Details              *ProductDetails
	SocialSharing        *SocialSharing
	SetProductProperties *SetProductProperties
}

// EnvelopeKey is the section's own name, one of the two wrappers a write takes; the facade's key is the other.
func (Configuration) EnvelopeKey() string { return "configuration" }

// ToMap renders the request body.
func (c Configuration) ToMap() map[string]any { return dto.ToMap(c) }

// LabelsAndDisplay is the labels_and_display block: how a variation is labelled and drawn.
type LabelsAndDisplay struct {
	DisableVariationLabelForPrice           *bool
	DisplayProductImageAsVariationIcons     *bool
	DisplaySelectedVariationValue           *bool
	EnableColourTypeVariationImageDisplay   *bool
	UseStrikeThroughForOutOfStockVariations *bool
	EnableSingleClickVariationLoading       *bool
}

// ProductPageVariationType is one entry of layout.variation_combinations; take ids from allowed_values.
type ProductPageVariationType struct {
	ID   *int64
	Name *string
}

// VariationLayout is the layout block: how variations are laid out and which price is shown.
type VariationLayout struct {
	SelectVariationView          *bool
	SelectVariationViewType      *string
	DisplayVariationCombinations *bool
	// An element is nil where the reader would not coerce the entry.
	VariationCombinations     []*ProductPageVariationType
	VariationPriceDisplayMode *bool
	VariationPriceDisplayType *string
}

// Variations is the variations section — how a product's variations are labelled and laid out.
type Variations struct {
	LabelsAndDisplay *LabelsAndDisplay
	Layout           *VariationLayout
}

// EnvelopeKey is the section's own name, one of the two wrappers a write takes; the facade's key is the other.
func (Variations) EnvelopeKey() string { return "variations" }

func (v Variations) ToMap() map[string]any { return dto.ToMap(v) }

// ImageDisplay is the image_display block: how product images are laid out and zoomed.
type ImageDisplay struct {
	EnableImageSlider            *bool
	ProductThumbnailVerticalView *bool
	ProductImageView             *bool
	ProductImageViewType         *string
	ImageZoomType                *bool
	ZoomType                     *string
}

// ProductVideo is the video block: whether a video is previewed, and where its thumbnail sits.
type ProductVideo struct {
	EnableVideoPreview                *bool
	DisplayVideoThumbnailInFirstOrder *bool
}

// ImageAndVideo is the image_and_video section — one field of its own beside two blocks.
type ImageAndVideo struct {
	AllowSingleFileUploadPerProduct *bool
	ImageDisplay                    *ImageDisplay
	Video                           *ProductVideo
}

// EnvelopeKey is the section's own name, one of the two wrappers a write takes; the facade's key is the other.
func (ImageAndVideo) EnvelopeKey() string { return "image_and_video" }

func (i ImageAndVideo) ToMap() map[string]any { return dto.ToMap(i) }

// ProductListing is the product listing both the up_sell and related_products sections carry.
type ProductListing struct {
	NumberOfProducts                 *int64
	AllowViewersToChangeDisplayMode  *bool
	GridView                         *bool
	ListView                         *bool
	TabularView                      *bool
	SelectDisplayType                *string
	MaximumNumberOfColumnsPerRow     *int64
	MaximumNumberOfColumnsPerRowList *int64
	ShortDescription                 *bool
	DisplayOnHover                   *bool
	DisplayAddToCartButton           *bool
	DisplayAddToCompareButton        *bool
	DisplayAddToWishlistButton       *bool
	Price                            *bool
	ShowExpectToPayPrice             *bool
	DisplayRating                    *bool
	DisplayReviewCount               *bool
}

// UpSell is the up_sell section, a ProductListing under its own name and envelope.
type UpSell ProductListing

// EnvelopeKey is the section's own name, one of the two wrappers a write takes; the facade's key is the other.
func (UpSell) EnvelopeKey() string { return "up_sell" }

func (u UpSell) ToMap() map[string]any { return dto.ToMap(u) }

// RelatedProducts is the related_products section, a ProductListing under its own name and envelope.
type RelatedProducts ProductListing

// EnvelopeKey is the section's own name, one of the two wrappers a write takes; the facade's key is the other.
func (RelatedProducts) EnvelopeKey() string { return "related_products" }

func (r RelatedProducts) ToMap() map[string]any { return dto.ToMap(r) }

// ReviewCapabilities is the capabilities block: what the review plugin offers, read-only.
type ReviewCapabilities struct {
	Available *bool
}

// ReviewAndRatings is the review_and_ratings section: where reviews appear and which version renders them.
type ReviewAndRatings struct {
	Capabilities      *ReviewCapabilities
	EnableReviewInTab *bool
	EnableNewVersion  *bool
}

// EnvelopeKey is the section's own name, one of the two wrappers a write takes; the facade's key is the other.
func (ReviewAndRatings) EnvelopeKey() string { return "review_and_ratings" }

func (r ReviewAndRatings) ToMap() map[string]any { return dto.ToMap(r) }

// QuickView is the quick_view section: the quick-view popup and what it shows.
type QuickView struct {
	ProductQuickView           *bool
	DisplayProductName         *bool
	DisplaySku                 *bool
	DisplayProductRating       *bool
	DisplayProductStock        *bool
	DisplayProductVariation    *bool
	DisplayProductImage        *bool
	DisplayVariationImage      *bool
	DisplayProductSummary      *bool
	DisplayAddToCartButton     *bool
	DisplayAddToCompareButton  *bool
	DisplayAddToWishlistButton *bool
}

// EnvelopeKey is the section's own name, one of the two wrappers a write takes; the facade's key is the other.
func (QuickView) EnvelopeKey() string { return "quick_view" }

func (q QuickView) ToMap() map[string]any { return dto.ToMap(q) }

// ProductPageSectionMeta is one group_meta entry; the plugin-backed one adds Plugin and Available.
type ProductPageSectionMeta struct {
	Kind      *string
	Writable  *bool
	Plugin    *string
	Available *bool
	UILabel   *string
	Href      *string
}

// ProductPageGroupMeta is the group_meta block, one entry per section; only the facade sends it.
type ProductPageGroupMeta struct {
	Configuration    *ProductPageSectionMeta
	Variations       *ProductPageSectionMeta
	ImageAndVideo    *ProductPageSectionMeta
	UpSell           *ProductPageSectionMeta
	RelatedProducts  *ProductPageSectionMeta
	ReviewAndRatings *ProductPageSectionMeta
	QuickView        *ProductPageSectionMeta
}

// ProductPage is the facade's body or one section's, both under the facade's key; ToMap renders either PATCH.
type ProductPage struct {
	Configuration    *Configuration
	Variations       *Variations
	ImageAndVideo    *ImageAndVideo
	UpSell           *UpSell
	RelatedProducts  *RelatedProducts
	ReviewAndRatings *ReviewAndRatings
	QuickView        *QuickView
	GroupMeta        *ProductPageGroupMeta
	// FieldMetadata is the raw field_metadata block.
	FieldMetadata any
}

func (ProductPage) EnvelopeKey() string { return "product_page" }

func (p ProductPage) ToMap() map[string]any { return dto.ToMap(p) }

// ProductPageShape is the typed view of the facade's body and of one section's.
func ProductPageShape(decoded any) any { return dto.Read[ProductPage](decoded) }
