package dashboardbanners

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// BannerStep is one step of an open_modal banner's walkthrough.
type BannerStep struct {
	Title       *string
	Description *string
}

// BannerLearnMore is what the banner's learn-more button does, read and written in one type. Action is open_link, with an https URL, or open_modal, with a Title, a Subtitle and at least one of Steps. A patch replaces the whole block, so an open_modal block sent over an open_link one drops the URL.
type BannerLearnMore struct {
	Action   *string
	Label    *string
	URL      *string
	Title    *string
	Subtitle *string
	Steps    []*BannerStep
}

// BannerTryIt is the optional second button; both fields are needed.
type BannerTryIt struct {
	Label *string
	URL   *string
}

// DashboardBanner is one banner. Key is the banner_key of the other calls. Position counts from 1 over every banner, hidden ones included, and Active false hides it from the Dashboard feed only. TryIt is nil when there is no second button.
type DashboardBanner struct {
	Key         *string
	Position    *int64
	Active      *bool
	Icon        *string
	Title       *string
	Description *string
	Tags        []*string
	LearnMore   *BannerLearnMore
	TryIt       *BannerTryIt
	CreatedAt   *string
	UpdatedAt   *string
}

// DashboardBannerListing is the management list: every banner, active or not, in position order.
type DashboardBannerListing struct {
	Banners []*DashboardBanner
}

// DashboardBannerInput is the write body, bare. A create needs Key (a-z, 0-9 and _, starting with a letter, up to 64 characters, unique), Title (1 to 80 characters), Description (1 to 200) and LearnMore; Tags holds at most 3, Active defaults to true and Position to the end of the list. A patch sends only what changes and cannot send Key. Icon is stored as sent. Position runs from 1 to one past the banner count, and the banners after it shift down.
type DashboardBannerInput struct {
	Key         string
	Icon        string
	Title       string
	Description string
	Tags        []string
	LearnMore   *BannerLearnMore
	TryIt       *BannerTryIt
	Position    *int64
	Active      *bool
}

// ToMap renders the request body.
func (i DashboardBannerInput) ToMap() map[string]any { return dto.ToMap(i) }

// BannerContact is the support card of the Dashboard.
type BannerContact struct {
	Email         *string
	Phone         *string
	LiveChat      *string
	HelpCenterURL *string
}

// DashboardContentInfo is the Dashboard's feed: only the active banners, in position order, the support card, whether the store is an online store and the label of the membership.
type DashboardContentInfo struct {
	Banners     []*DashboardBanner
	Contact     *BannerContact
	OnlineStore *bool
	MemberLabel *string
}

// DashboardBannerShape is the typed view of one banner, a bare body.
func DashboardBannerShape(decoded any) any { return dto.Read[DashboardBanner](decoded) }

// DashboardBannerListingShape is the typed view of the management list.
func DashboardBannerListingShape(decoded any) any { return dto.Read[DashboardBannerListing](decoded) }

// DashboardContentInfoShape is the typed view of the Dashboard feed.
func DashboardContentInfoShape(decoded any) any { return dto.Read[DashboardContentInfo](decoded) }
