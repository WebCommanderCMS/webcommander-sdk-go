package animationeffects

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// CustomLoader is the uploaded loader that goes with enable_custom_preloader; it is read-only.
type CustomLoader struct {
	Enabled  *bool
	FileName *string
	FileURL  *string
}

// PreloaderPage is one entry of preloader_custom_page as a read sends it; a write takes page ids.
type PreloaderPage struct {
	ID   *int64
	Name *string
}

// PreloaderDesign is the preloader_design section: the loader shown while a page loads; its depends_on rules are enforced.
type PreloaderDesign struct {
	EnablePreloader            *bool
	EnableCustomPreloader      *bool
	PreloaderPreset            *string
	PreloaderSize              *int64
	PreloaderBackgroundOpacity *float64
	BackgroundPreloader        *string
	PreloaderHideCondition     *string
	PreloaderShowTimeCondition *string
	PreloaderHideAfterTime     *float64
	PreloaderPage              *string
	// The one list this area carries; a write takes plain page ids in its place.
	PreloaderCustomPage []*PreloaderPage
	CustomLoader        *CustomLoader
}

// EnvelopeKey is the section's own name, the only wrapper a write takes.
func (PreloaderDesign) EnvelopeKey() string { return "preloader_design" }

// ToMap renders the request body.
func (p PreloaderDesign) ToMap() map[string]any { return dto.ToMap(p) }

// WidgetTransition is the widget_transition section — whether widgets fade in, and over how long.
type WidgetTransition struct {
	WidgetSmoothRender         *bool
	WidgetSmoothRenderDuration *int64
}

// EnvelopeKey is the section's own name, the only wrapper a write takes.
func (WidgetTransition) EnvelopeKey() string { return "widget_transition" }

func (w WidgetTransition) ToMap() map[string]any { return dto.ToMap(w) }

// PageTransition is the page_transition section — the animation between pages and when it plays.
type PageTransition struct {
	PageTransitionAnimationType    *string
	TransitionAnimationDuration    *float64
	AnimationBeforeEnteringNewPage *bool
	AnimationAfterEnteringNewPage  *bool
}

// EnvelopeKey is the section's own name, the only wrapper a write takes.
func (PageTransition) EnvelopeKey() string { return "page_transition" }

func (p PageTransition) ToMap() map[string]any { return dto.ToMap(p) }

// ImageLazyEffect is the image_lazy_effect section: which pages lazy-load images, and the placeholder.
type ImageLazyEffect struct {
	PreloaderImageLazyEffectPreset     *string
	ProductDetailsPageImageLazyEffect  *bool
	FlyingCartImageLazyEffect          *bool
	ProductListImageLazyEffect         *bool
	PropertyListImageLazyEffect        *bool
	PropertyDetailsPageImageLazyEffect *bool
	BlogListImageLazyEffect            *bool
	BlogDetailsPageImageLazyEffect     *bool
	LazyLoaderSize                     *int64
	BackgroundLazyLoader               *string
	ContentImageLazyLoadBackground     *string
}

// EnvelopeKey is the section's own name, the only wrapper a write takes.
func (ImageLazyEffect) EnvelopeKey() string { return "image_lazy_effect" }

func (i ImageLazyEffect) ToMap() map[string]any { return dto.ToMap(i) }

// PageSkeleton is the page_skeleton section — the one-field section of this area.
type PageSkeleton struct {
	CheckoutPageSkeleton *bool
}

// EnvelopeKey is the section's own name, the only wrapper a write takes.
func (PageSkeleton) EnvelopeKey() string { return "page_skeleton" }

func (p PageSkeleton) ToMap() map[string]any { return dto.ToMap(p) }

// SectionMeta is one group_meta entry: what a section is, whether it is writable, where it is managed.
type SectionMeta struct {
	Kind     *string
	Writable *bool
	Href     *string
}

// GroupMeta is the group_meta block, one SectionMeta per section; only the facade sends it.
type GroupMeta struct {
	PreloaderDesign  *SectionMeta
	WidgetTransition *SectionMeta
	PageTransition   *SectionMeta
	ImageLazyEffect  *SectionMeta
	PageSkeleton     *SectionMeta
}

// AnimationEffects is the facade's body or one section's; ToMap renders the facade PATCH, not a section's.
type AnimationEffects struct {
	PreloaderDesign  *PreloaderDesign
	WidgetTransition *WidgetTransition
	PageTransition   *PageTransition
	ImageLazyEffect  *ImageLazyEffect
	PageSkeleton     *PageSkeleton
	GroupMeta        *GroupMeta
	// FieldMetadata is the raw field_metadata block, keyed by section on the facade and by field on a section.
	FieldMetadata any
}

func (AnimationEffects) EnvelopeKey() string { return "animation_effects" }

func (a AnimationEffects) ToMap() map[string]any { return dto.ToMap(a) }

// AnimationEffectsShape is the typed view of the facade's body and of one section's.
func AnimationEffectsShape(decoded any) any { return dto.Read[AnimationEffects](decoded) }
