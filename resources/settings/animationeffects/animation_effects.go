package animationeffects

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// AnimationEffectsSettings is the animation effects endpoints, bound to one client: a facade over five sections.
type AnimationEffectsSettings struct {
	requester http.Requester
}

// GetAnimationEffects reads the whole area — GET /admin/settings/animation_effects.
func (s *AnimationEffectsSettings) GetAnimationEffects(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, AnimationEffectsPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: AnimationEffectsShape})
}

// PatchAnimationEffects writes one or more sections — PATCH /admin/settings/animation_effects, wrapped in animation_effects.
func (s *AnimationEffectsSettings) PatchAnimationEffects(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, AnimationEffectsPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: AnimationEffectsShape})
}

// HeadAnimationEffects asks whether the area is reachable — HEAD /admin/settings/animation_effects.
func (s *AnimationEffectsSettings) HeadAnimationEffects(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", AnimationEffectsPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// GetAnimationEffectsSection reads one section — GET /admin/settings/animation_effects/{section}, under the section's key.
func (s *AnimationEffectsSettings) GetAnimationEffectsSection(ctx context.Context, section string, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, AnimationEffectsSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, QueryParams: queryParams,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: AnimationEffectsShape})
}

// PatchAnimationEffectsSection writes one section — PATCH /admin/settings/animation_effects/{section}, wrapped in the section's key.
func (s *AnimationEffectsSettings) PatchAnimationEffectsSection(ctx context.Context, section string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, AnimationEffectsSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: AnimationEffectsShape})
}

// UploadAnimationEffectsCustomLoader uploads the custom preloader image — POST /admin/settings/animation_effects/custom_loader, as a part named file (svg, png, jpg, jpeg or gif; another type is 400 invalid_file, and no file 400). It answers 201 {"custom_loader": {"enabled": true, "file_name", "file_url"}}.
func (s *AnimationEffectsSettings) UploadAnimationEffectsCustomLoader(ctx context.Context, files any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	normalised, err := normaliseFiles(files)
	if err != nil {
		return nil, err
	}
	if err := requireFiles(files, normalised); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPostRequest(ctx, AnimationEffectsCustomLoaderPath, &http.RequestOptions{
		Files: normalised, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CustomLoaderSectionShape})
}

// DeleteAnimationEffectsCustomLoader removes the custom preloader image — DELETE /admin/settings/animation_effects/custom_loader; 204, also when there is none, and the loader then reads enabled false with null file fields.
func (s *AnimationEffectsSettings) DeleteAnimationEffectsCustomLoader(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendDeleteRequest(ctx, AnimationEffectsCustomLoaderPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// normaliseFiles accepts map[string]any or [][2]any and returns the map the transport encodes.
func normaliseFiles(files any) (map[string]any, error) {
	switch shape := files.(type) {
	case nil:
		return nil, nil
	case map[string]any:
		return shape, nil
	case [][2]any:
		out := make(map[string]any, len(shape))
		for _, pair := range shape {
			name, ok := pair[0].(string)
			if !ok {
				return nil, wcerrors.NewUsageError("files: a pair's first element must be the field " +
					"name as a string")
			}
			if _, clash := out[name]; clash {
				return nil, wcerrors.NewUsageError("files: the field " + name + " appears twice; a map " +
					"cannot carry two values under one name, so the duplicate would be lost")
			}
			out[name] = pair[1]
		}
		return out, nil
	default:
		return nil, wcerrors.NewUsageError("files: expected map[string]any or [][2]any " +
			"(field, value) pairs")
	}
}

// requireFiles refuses an upload carrying nothing — no files, or an empty container — before any I/O.
func requireFiles(files any, normalised map[string]any) error {
	if files == nil {
		return wcerrors.NewUsageError("files is required")
	}
	if len(normalised) == 0 {
		return wcerrors.NewUsageError(fmt.Sprintf("files must name the file to upload, got an empty %T", files))
	}
	return nil
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "animation_effects"

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

// New binds the animation effects endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *AnimationEffectsSettings {
	return &AnimationEffectsSettings{requester: requester}
}
