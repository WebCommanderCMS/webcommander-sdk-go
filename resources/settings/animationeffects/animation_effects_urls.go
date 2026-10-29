package animationeffects

// The path templates for the animation effects endpoints; the client fills {api_version}.
const (
	AnimationEffectsPath = "/api/{api_version}/admin/settings/animation_effects"

	AnimationEffectsSectionPath = "/api/{api_version}/admin/settings/animation_effects/{section}"
)

// animationEffectsSections are the five values {section} takes, in facade order; all are writable.
var animationEffectsSections = [5]string{
	"preloader_design",
	"widget_transition",
	"page_transition",
	"image_lazy_effect",
	"page_skeleton",
}

// AnimationEffectsSections returns a fresh copy of the five values {section} takes.
func AnimationEffectsSections() []string {
	return append(make([]string, 0, len(animationEffectsSections)), animationEffectsSections[:]...)
}
