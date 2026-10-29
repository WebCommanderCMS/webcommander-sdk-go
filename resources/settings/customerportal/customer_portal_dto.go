package customerportal

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// CustomerGroup is one entry of default_groups or approval_groups as a read sends it; a write sends ids.
type CustomerGroup struct {
	ID   *int64
	Name *string
}

// Portal is the portal section: the one switch that turns the customer portal on.
type Portal struct {
	Enabled *bool
}

// EnvelopeKey is the section's own name, also one of the two wrappers a write takes.
func (Portal) EnvelopeKey() string { return "portal" }

// ToMap renders the request body.
func (p Portal) ToMap() map[string]any { return dto.ToMap(p) }

// RegistrationTerms is the terms block; enabling terms while Reference is null is 400.
type RegistrationTerms struct {
	Enabled    *bool
	Text       *string
	AnchorText *string
	Type       *string
	Reference  *string
}

// FacebookLogin is the facebook block; AppID is required when enabled.
type FacebookLogin struct {
	Enabled *bool
	AppID   *string
}

// GoogleLogin is the google block; ClientID is required when enabled.
type GoogleLogin struct {
	Enabled  *bool
	ClientID *string
}

// B2bConfiguration is the b2b_configuration block: four booleans for what a new registrant becomes.
type B2bConfiguration struct {
	Enabled              *bool
	CreateAsOrganisation *bool
	CreateAsCustomer     *bool
	OrganisationSelected *bool
}

// AlternateLogin is the alternate_login block: third-party login, and whether it offers an email field.
type AlternateLogin struct {
	Enabled          *bool
	EnableEmailField *bool
}

// Registration is the registration section: who may register, as what, and what they must accept.
type Registration struct {
	Enabled                  *bool
	CloseRegistrationMessage *string
	AssignDefaultGroup       *bool
	DefaultGroups            []*CustomerGroup
	RequiresApproval         *bool
	ApprovalGroups           []*CustomerGroup
	Terms                    *RegistrationTerms
	Facebook                 *FacebookLogin
	Google                   *GoogleLogin
	B2bConfiguration         *B2bConfiguration
	AlternateLogin           *AlternateLogin
}

// EnvelopeKey is the section's own name, also one of the two wrappers a write takes.
func (Registration) EnvelopeKey() string { return "registration" }

func (r Registration) ToMap() map[string]any { return dto.ToMap(r) }

// LoginCaptcha is the captcha block shown after failed logins; FailCount is bounded.
type LoginCaptcha struct {
	Enabled   *bool
	FailCount *int64
}

// Login is the login section: whether customers may sign in and what the sign-in page offers.
type Login struct {
	Enabled                               *bool
	LoginDisabledMessage                  *string
	AwaitingApprovalMessage               *string
	PasswordRecovery                      *bool
	RegistrationLink                      *bool
	Captcha                               *LoginCaptcha
	EnableLandingPageAfterSuccessfulLogin *bool
	LandingPageAfterSuccessfulLogin       *string
}

// EnvelopeKey is the section's own name, also one of the two wrappers a write takes.
func (Login) EnvelopeKey() string { return "login" }

func (l Login) ToMap() map[string]any { return dto.ToMap(l) }

// AzureAdDataMapping is the azure_ad.data_mapping block; every entry is a string but CountryID.
type AzureAdDataMapping struct {
	AddressLine1 *string
	AddressLine2 *string
	PostCode     *string
	Phone        *string
	Mobile       *string
	CountryID    *int64
	State        *string
	City         *string
}

// AzureAd is the azure_ad block, the SAML provider.
type AzureAd struct {
	DataMapping     *AzureAdDataMapping
	Identifier      *string
	LoginURL        *string
	ReplyURLPostfix *string
}

// AzureAdB2c is the azure_ad_b2c block, the B2C provider; the SDK redacts its ClientSecret.
type AzureAdB2c struct {
	TenantName      *string
	PolicyName      *string
	ClientID        *string
	ClientSecret    *string
	ReplyURLPrefix  *string
	ReplyURLPostfix *string
}

// Sso is the sso section; turning Enabled on demands the whole azure_ad_b2c block, secret included.
type Sso struct {
	Enabled           *bool
	RequireSiteAccess *bool
	ServiceProvider   *string
	AzureAd           *AzureAd
	AzureAdB2c        *AzureAdB2c
}

// EnvelopeKey is the section's own name, also one of the two wrappers a write takes.
func (Sso) EnvelopeKey() string { return "sso" }

func (s Sso) ToMap() map[string]any { return dto.ToMap(s) }

// ProfileField is one row of an address or account form on the profile page; Key is read-only.
type ProfileField struct {
	Key    *string
	Active *bool
	Order  *int64
}

// ProfileEntry is one row of a named list on the profile page: ProfileField's fields and a Label.
type ProfileEntry struct {
	Key    *string
	Active *bool
	Order  *int64
	Label  *string
}

// ProfileOverviewEntry is one row of overview; ShowProductInfo is sent on some rows only.
type ProfileOverviewEntry struct {
	Key             *string
	Active          *bool
	Order           *int64
	Label           *string
	ShowProductInfo *bool
}

// ManageMyAccount is the manage_my_account block: the account screen's sections and form fields.
type ManageMyAccount struct {
	Sections           []*ProfileEntry
	AccountInformation []*ProfileField
	BillingAddress     []*ProfileField
	ShippingAddress    []*ProfileField
}

// MyOrderFlags is the my_orders.flags block: six booleans.
type MyOrderFlags struct {
	EnableDuePayment     *bool
	ShowTrackMyOrders    *bool
	DisplayTrackInColumn *bool
	ShowOrderChannel     *bool
	ShowPayNow           *bool
	AllowOrderShare      *bool
}

// MyOrders is the my_orders block — what the orders tab shows and which order lists it offers.
type MyOrders struct {
	Flags    *MyOrderFlags
	Sections []*ProfileEntry
}

// AssignToCustomer is the assign_to_customer block: the assigned-pages tab and its pagination.
type AssignToCustomer struct {
	Items          []*ProfileEntry
	ShowPagination *string
	ItemPerPage    *int64
}

// ProfilePage is the profile_page section, read-only: every tab of the customer's own page.
type ProfilePage struct {
	PageHTML         *string
	Tabs             []*ProfileEntry
	Overview         []*ProfileOverviewEntry
	ManageMyAccount  *ManageMyAccount
	MyOrders         *MyOrders
	MyEntitlements   []*ProfileEntry
	AssignToCustomer *AssignToCustomer
}

// EnvelopeKey is the section's own name, which is the key a read of this section arrives under.
func (ProfilePage) EnvelopeKey() string { return "profile_page" }

func (p ProfilePage) ToMap() map[string]any { return dto.ToMap(p) }

// RegistrationField is one row of the registration form's field catalog.
type RegistrationField struct {
	Key          *string
	Label        *string
	DefaultLabel *string
	Active       *bool
	Required     *bool
	System       *bool
	Order        *int64
	Placeholder  *string
}

// RegistrationFields is the registration_fields section, read-only: the registration form's field catalog.
type RegistrationFields struct {
	Fields []*RegistrationField
}

// EnvelopeKey is the section's own name, which is the key a read of this section arrives under.
func (RegistrationFields) EnvelopeKey() string { return "registration_fields" }

func (r RegistrationFields) ToMap() map[string]any { return dto.ToMap(r) }

// SectionMeta is one group_meta entry: what a section is, whether it is writable, where it is managed.
type SectionMeta struct {
	Kind     *string
	Writable *bool
	Href     *string
}

// GroupMeta is the group_meta block, one SectionMeta per section; only the facade sends it.
type GroupMeta struct {
	Portal             *SectionMeta
	Registration       *SectionMeta
	Login              *SectionMeta
	Sso                *SectionMeta
	ProfilePage        *SectionMeta
	RegistrationFields *SectionMeta
}

// CustomerPortal is the facade's body or one section's; ToMap renders the facade PATCH, not a section's.
type CustomerPortal struct {
	Portal             *Portal
	Registration       *Registration
	Login              *Login
	Sso                *Sso
	ProfilePage        *ProfilePage
	RegistrationFields *RegistrationFields
	GroupMeta          *GroupMeta
	// FieldMetadata is the raw field_metadata block.
	FieldMetadata any
}

func (CustomerPortal) EnvelopeKey() string { return "customer_portal" }

func (c CustomerPortal) ToMap() map[string]any { return dto.ToMap(c) }

// CustomerPortalShape is the typed view of the facade's body and of one section's.
func CustomerPortalShape(decoded any) any { return dto.Read[CustomerPortal](decoded) }
