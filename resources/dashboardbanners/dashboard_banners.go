package dashboardbanners

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// DashboardBanners is the dashboard banner endpoints, bound to one client.
type DashboardBanners struct {
	requester http.Requester
}

// ListDashboardBanners lists every banner, active or not — GET /admin/dashboard/banners, {"banners": [...]} in position order, positions counting from 1 without a gap.
func (b *DashboardBanners) ListDashboardBanners(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendGetRequest(ctx, DashboardBannersPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: DashboardBannerListingShape})
}

// CreateDashboardBanner creates a banner — POST /admin/dashboard/banners, answered 201 with the banner bare and no Location. The body is bare; every failing field is listed in one 400, a key already used is 400 duplicate, an unknown field is 400. It joins the Dashboard feed at once while active.
func (b *DashboardBanners) CreateDashboardBanner(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendPostRequest(ctx, DashboardBannersPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: DashboardBannerShape})
}

// GetDashboardBanner reads one banner — GET /admin/dashboard/banners/{banner_key}, where the slot takes the banner's key; the banner is bare, hidden ones included. An unknown key is 404.
func (b *DashboardBanners) GetDashboardBanner(ctx context.Context, bannerKey string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendGetRequest(ctx, DashboardBannerPath, &http.RequestOptions{
		PathParams:   map[string]any{"banner_key": bannerKey},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: DashboardBannerShape})
}

// PatchDashboardBanner changes only the fields sent — PATCH /admin/dashboard/banners/{banner_key}, answered 200 with the banner. An empty body, a blank title, a key and an unknown field are 400; an unknown banner is 404. learn_more and try_it are replaced whole, and a try_it of null removes the second button (send it in a map payload; an input leaves an unset field out). updated_at is not moved.
func (b *DashboardBanners) PatchDashboardBanner(ctx context.Context, bannerKey string, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendPatchRequest(ctx, DashboardBannerPath, &http.RequestOptions{
		PathParams: map[string]any{"banner_key": bannerKey}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: DashboardBannerShape})
}

// DeleteDashboardBanner deletes one banner — DELETE /admin/dashboard/banners/{banner_key}, 204 with no body. The banners after it move up, so the positions stay without a gap; an unknown key is 404.
func (b *DashboardBanners) DeleteDashboardBanner(ctx context.Context, bannerKey string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendDeleteRequest(ctx, DashboardBannerPath, &http.RequestOptions{
		PathParams:   map[string]any{"banner_key": bannerKey},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// GetDashboardContentInfo reads what the Dashboard shows besides the figures — GET /admin/dashboard/content_info: the active banners under banners, the support card under contact, online_store and member_label. A hidden banner is not listed.
func (b *DashboardBanners) GetDashboardContentInfo(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return b.requester.SendGetRequest(ctx, DashboardContentInfoPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: DashboardContentInfoShape})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "dashboard_banners"

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

// New binds the dashboard banner endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *DashboardBanners {
	return &DashboardBanners{requester: requester}
}
