package seoconfiguration

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/customers"
)

// SeoConfigurationSettings is the SEO configuration settings endpoints: a facade over three sections and two collections.
type SeoConfigurationSettings struct {
	requester http.Requester
}

// GetSeoConfiguration reads the whole area — GET /admin/settings/seo_configuration.
func (s *SeoConfigurationSettings) GetSeoConfiguration(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, SeoConfigurationPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: SeoConfigurationShape})
}

// PatchSeoConfiguration writes one or more sections — PATCH /admin/settings/seo_configuration; omitted fields keep their values.
func (s *SeoConfigurationSettings) PatchSeoConfiguration(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, SeoConfigurationPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: SeoConfigurationShape})
}

// PutSeoConfiguration replaces the sections a body names — PUT /admin/settings/seo_configuration; send each section whole.
func (s *SeoConfigurationSettings) PutSeoConfiguration(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPutRequest(ctx, SeoConfigurationPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: SeoConfigurationShape})
}

// HeadSeoConfiguration asks whether the area is reachable — HEAD /admin/settings/seo_configuration.
func (s *SeoConfigurationSettings) HeadSeoConfiguration(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", SeoConfigurationPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// GetSeoConfigurationSection reads one section — GET /admin/settings/seo_configuration/{section}.
func (s *SeoConfigurationSettings) GetSeoConfigurationSection(ctx context.Context, section string, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, SeoConfigurationSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, QueryParams: queryParams,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: SeoConfigurationShape})
}

// PatchSeoConfigurationSection writes one section — PATCH /admin/settings/seo_configuration/{section}, under its own key.
func (s *SeoConfigurationSettings) PatchSeoConfigurationSection(ctx context.Context, section string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, SeoConfigurationSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: SeoConfigurationShape})
}

// PutSeoConfigurationSection replaces one section — PUT /admin/settings/seo_configuration/{section}; send it back whole.
func (s *SeoConfigurationSettings) PutSeoConfigurationSection(ctx context.Context, section string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPutRequest(ctx, SeoConfigurationSectionPath, &http.RequestOptions{
		PathParams: map[string]any{"section": section}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: SeoConfigurationShape})
}

// GetSeoConfigurationDocumentRoot lists the served files — GET /admin/settings/seo_configuration/document_root; limit caps at 100.
func (s *SeoConfigurationSettings) GetSeoConfigurationDocumentRoot(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, SeoConfigurationDocumentRootPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: customers.PageOf[DocumentRootFile]})
}

// HeadSeoConfigurationDocumentRoot asks whether the served files are reachable — HEAD /admin/settings/seo_configuration/document_root; no ETag.
func (s *SeoConfigurationSettings) HeadSeoConfigurationDocumentRoot(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", SeoConfigurationDocumentRootPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// UploadSeoConfigurationDocumentRoot adds a file — POST /admin/settings/seo_configuration/document_root, as a part named file.
func (s *SeoConfigurationSettings) UploadSeoConfigurationDocumentRoot(ctx context.Context, files any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
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
	return s.requester.SendPostRequest(ctx, SeoConfigurationDocumentRootPath, &http.RequestOptions{
		Files: normalised, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: DocumentRootFileShape})
}

// DeleteSeoConfigurationDocumentRootFile stops serving one file — DELETE /admin/settings/seo_configuration/document_root/{config_key}.
func (s *SeoConfigurationSettings) DeleteSeoConfigurationDocumentRootFile(ctx context.Context, configKey any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendDeleteRequest(ctx, SeoConfigurationDocumentRootFilePath, &http.RequestOptions{
		PathParams: map[string]any{"config_key": configKey}, RaiseOnError: raiseOnError,
		ResponseFormat: format})
}

// GetSeoConfigurationRedirects lists the redirects under items — GET /admin/settings/seo_configuration/redirects; limit caps at 100.
func (s *SeoConfigurationSettings) GetSeoConfigurationRedirects(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, SeoConfigurationRedirectsPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: customers.PageOf[Redirect]})
}

// HeadSeoConfigurationRedirects asks whether the redirects are reachable — HEAD /admin/settings/seo_configuration/redirects; no ETag.
func (s *SeoConfigurationSettings) HeadSeoConfigurationRedirects(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", SeoConfigurationRedirectsPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CreateSeoConfigurationRedirect adds a redirect — POST /admin/settings/seo_configuration/redirects; source_url and destination_url are both required.
func (s *SeoConfigurationSettings) CreateSeoConfigurationRedirect(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPostRequest(ctx, SeoConfigurationRedirectsPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: RedirectShape})
}

// GetSeoConfigurationRedirect reads one redirect — GET /admin/settings/seo_configuration/redirects/{redirect_id}.
func (s *SeoConfigurationSettings) GetSeoConfigurationRedirect(ctx context.Context, redirectID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendGetRequest(ctx, SeoConfigurationRedirectPath, &http.RequestOptions{
		PathParams: map[string]any{"redirect_id": redirectID}, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: RedirectShape})
}

// HeadSeoConfigurationRedirect asks whether one redirect is there — HEAD /admin/settings/seo_configuration/redirects/{redirect_id}; no ETag, and 404 for an unknown id.
func (s *SeoConfigurationSettings) HeadSeoConfigurationRedirect(ctx context.Context, redirectID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendApiRequest(ctx, "HEAD", SeoConfigurationRedirectPath, &http.RequestOptions{
		PathParams: map[string]any{"redirect_id": redirectID}, RaiseOnError: raiseOnError,
		ResponseFormat: format})
}

// PatchSeoConfigurationRedirect changes one redirect — PATCH /admin/settings/seo_configuration/redirects/{redirect_id}.
func (s *SeoConfigurationSettings) PatchSeoConfigurationRedirect(ctx context.Context, redirectID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendPatchRequest(ctx, SeoConfigurationRedirectPath, &http.RequestOptions{
		PathParams: map[string]any{"redirect_id": redirectID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: RedirectShape})
}

// DeleteSeoConfigurationRedirect removes one redirect — DELETE /admin/settings/seo_configuration/redirects/{redirect_id}.
func (s *SeoConfigurationSettings) DeleteSeoConfigurationRedirect(ctx context.Context, redirectID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return s.requester.SendDeleteRequest(ctx, SeoConfigurationRedirectPath, &http.RequestOptions{
		PathParams: map[string]any{"redirect_id": redirectID}, RaiseOnError: raiseOnError,
		ResponseFormat: format})
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

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "seo_configuration"

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

// New binds the SEO configuration settings endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *SeoConfigurationSettings {
	return &SeoConfigurationSettings{requester: requester}
}
