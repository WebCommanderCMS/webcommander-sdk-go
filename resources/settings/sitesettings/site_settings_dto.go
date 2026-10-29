package sitesettings

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// Region is a country or a state as a read sends it; a write may send a bare id, code or name instead.
type Region struct {
	ID   *int64
	Code *string
	Name *string
}

// AutoSuggestAddresses is the auto_suggest.addresses block: suggest addresses as typed, then validate.
type AutoSuggestAddresses struct {
	Enabled  *bool
	Validate *bool
}

// AutoSuggest is the auto_suggest block — address and suburb suggestions.
type AutoSuggest struct {
	Addresses     *AutoSuggestAddresses
	SuburbsCities *bool
}

// Captcha is the captcha block; PrivateKey and PrivateKeyV3 are secret and read back as ***.
type Captcha struct {
	Enabled             *bool
	MinimumScore        *float64
	PrivateKey          *string
	PrivateKeyTurnstile *string
	// Tagged: the mechanical spelling would ask the store for private_key_v_3.
	PrivateKeyV3       *string `wire:"private_key_v3"`
	PublicKey          *string
	PublicKeyTurnstile *string
	PublicKeyV3        *string `wire:"public_key_v3"`
	Type               *string
	UseOwn             *bool
	Version            *string
}

// Administration is the administration section; MaxPrecision is bounded 1 to 9, and Page404 must name a page.
type Administration struct {
	AutoSuggest      *AutoSuggest
	Captcha          *Captcha
	DefaultCountry   *Region
	DefaultState     *Region
	EcommerceEnabled *bool
	MaxPrecision     *int64
	Page403          *string
	Page404          *string
	UnitLength       *string
	UnitWeight       *string
}

// EnvelopeKey is the section's own name, also one of the two wrappers a write takes.
func (Administration) EnvelopeKey() string { return "administration" }

// ToMap renders the request body.
func (a Administration) ToMap() map[string]any { return dto.ToMap(a) }

// BrandColors is the brand_kit.colors block: six roles, each a hex colour, else 400 invalid_color.
type BrandColors struct {
	Accent     *string
	Background *string
	Primary    *string
	Secondary  *string
	Surface    *string
	Text       *string
}

// BrandKit is the brand_kit section — the site-wide palette.
type BrandKit struct {
	Colors      *BrandColors
	Enabled     *bool
	PaletteName *string
}

// EnvelopeKey is the section's own name, also one of the two wrappers a write takes.
func (BrandKit) EnvelopeKey() string { return "brand_kit" }

func (b BrandKit) ToMap() map[string]any { return dto.ToMap(b) }

// AdminDateTime is the date_time.admin block, the only audience block that carries FirstDayOfWeek.
type AdminDateTime struct {
	DateFormat     *string
	TimeFormat     *string
	TimeZone       *string
	FirstDayOfWeek *string
}

// SeparableDateTime is the date_time.email and date_time.user blocks: two formats, a time zone and a Separate flag.
type SeparableDateTime struct {
	DateFormat *string
	TimeFormat *string
	TimeZone   *string
	Separate   *bool
}

// DateTime is the date_time section: formats and zones for the admin area, emails and the storefront.
type DateTime struct {
	Admin        *AdminDateTime
	Email        *SeparableDateTime
	HideTimeZone *bool
	User         *SeparableDateTime
}

// EnvelopeKey is the section's own name, also one of the two wrappers a write takes.
func (DateTime) EnvelopeKey() string { return "date_time" }

func (d DateTime) ToMap() map[string]any { return dto.ToMap(d) }

// Maintenance is the maintenance block: whether the storefront is closed, and what it says meanwhile.
type Maintenance struct {
	Enabled *bool
	Message *string
}

// SecurityCode is the security_code block, the site lock; Code is secret and reads back masked.
type SecurityCode struct {
	Code       *string
	Enabled    *bool
	NumberOnly *bool
}

// SslKeyFile is the ssl.key_file block: what the store holds for the private key, never the key itself.
type SslKeyFile struct {
	FileName *string
	Uploaded *bool
}

// SslCertificateFile is the ssl.certificate_file block: SslKeyFile's fields plus subject and expiry.
type SslCertificateFile struct {
	FileName  *string
	Uploaded  *bool
	ExpiresAt *string
	Subject   *string
}

// Ssl is the ssl block; its read keys and the write-only *_base64 keys share no name.
type Ssl struct {
	CertificateFile *SslCertificateFile
	KeyFile         *SslKeyFile
}

// DomainSettings is the domain_settings section; Aliases is read-only and managed elsewhere.
type DomainSettings struct {
	Aliases        []string
	HTTPSEnabled   *bool
	Maintenance    *Maintenance
	SecurityCode   *SecurityCode
	Ssl            *Ssl
	WebsiteAddress *string
}

// EnvelopeKey is the section's own name, also one of the two wrappers a write takes.
func (DomainSettings) EnvelopeKey() string { return "domain_settings" }

func (d DomainSettings) ToMap() map[string]any { return dto.ToMap(d) }

// EmailProvider is one entry of available; ConnectURL is sent only for the oauth entries.
type EmailProvider struct {
	Key        *string
	Label      *string
	Type       *string
	ConnectURL *string
}

// ReplyTo is a reply block, the same three fields under default.domains and under sendgrid.
type ReplyTo struct {
	Different *bool
	ToEmail   *string
	ToName    *string
}

// DefaultDomain is one entry of default.domains: a platform sending domain and this store's sender on it.
type DefaultDomain struct {
	Configure   *bool
	Domain      *string
	EmailPrefix *string
	Reply       *ReplyTo
	SenderName  *string
}

// DefaultEmail is the default block: the platform's own sending domains.
type DefaultEmail struct {
	Domains []*DefaultDomain
}

// OauthProvider is the gmail and office365 blocks, set through the interactive connect_url flow.
type OauthProvider struct {
	ConnectURL *string
}

// SendGrid is the sendgrid block; APIKey is secret and reads back as ***, so never echo a read into a write.
type SendGrid struct {
	APIKey      *string
	Reply       *ReplyTo
	SenderEmail *string
	SenderName  *string
}

// SmtpAuthentication is the smtp.authentication block; Password is secret and reads back as ***.
type SmtpAuthentication struct {
	Enabled  *bool
	Password *string
	Username *string
}

// Smtp is the smtp block: the server the store sends through and who it sends as.
type Smtp struct {
	Authentication *SmtpAuthentication
	Encryption     *string
	Host           *string
	Port           *int64
	SenderEmail    *string
	SenderName     *string
}

// EmailSettings is the email_settings section; write only the field you move, so masked secrets stay put.
type EmailSettings struct {
	Available []*EmailProvider
	Default   *DefaultEmail
	Gmail     *OauthProvider
	// Tagged: the mechanical spelling would ask the store for office_365.
	Office365       *OauthProvider `wire:"office365"`
	Sendgrid        *SendGrid
	ServiceProvider *string
	Smtp            *Smtp
}

// EnvelopeKey is the section's own name, also one of the two wrappers a write takes.
func (EmailSettings) EnvelopeKey() string { return "email_settings" }

func (e EmailSettings) ToMap() map[string]any { return dto.ToMap(e) }

// MediaImage is the image record of either media section; Height, Width and SizeBytes are not always sent.
type MediaImage struct {
	ContentType *string
	FileName    *string
	Height      *int64
	SizeBytes   *int64
	URL         *string
	Width       *int64
}

// Favicon is the favicon section, kind media; both fields are read-only, so it is read, never written.
type Favicon struct {
	Enabled *bool
	Image   *MediaImage
}

// EnvelopeKey is the section's own name, which is the key a read of this section arrives under.
func (Favicon) EnvelopeKey() string { return "favicon" }

func (f Favicon) ToMap() map[string]any { return dto.ToMap(f) }

// StoreLogoVariants is the store_logo.variants block — the logos for dark backgrounds and the sticky header.
type StoreLogoVariants struct {
	// Both read null while the main logo is used; a set variant's shape has not been measured.
	DarkBackground any
	StickyHeader   any
}

// StoreLogo is the store_logo section, kind media: read, never written, like Favicon.
type StoreLogo struct {
	Image    *MediaImage
	Variants *StoreLogoVariants
}

// EnvelopeKey is the section's own name, which is the key a read of this section arrives under.
func (StoreLogo) EnvelopeKey() string { return "store_logo" }

func (s StoreLogo) ToMap() map[string]any { return dto.ToMap(s) }

// StoreAddress is the store_details.address block; AddressLine1 is required.
type StoreAddress struct {
	AddressLine1 *string
	AddressLine2 *string
	City         *string
	Country      *Region
	PostCode     *string
	State        *Region
}

// StoreContact is the store_details.contact block — the three numbers the store publishes.
type StoreContact struct {
	Fax    *string
	Mobile *string
	Phone  *string
}

// StoreDetails is the store_details section; CompanyName and Email are required, and ABN takes eleven digits.
type StoreDetails struct {
	ABN         *string
	Address     *StoreAddress
	CompanyName *string
	Contact     *StoreContact
	Email       *string
}

// EnvelopeKey is the section's own name, also one of the two wrappers a write takes.
func (StoreDetails) EnvelopeKey() string { return "store_details" }

func (s StoreDetails) ToMap() map[string]any { return dto.ToMap(s) }

// SectionMeta is one group_meta entry: what a section is, whether it is writable, where it is managed.
type SectionMeta struct {
	Kind     *string
	Writable *bool
	Href     *string
}

// GroupMeta is the group_meta block, one SectionMeta per section; only the facade sends it.
type GroupMeta struct {
	Administration *SectionMeta
	BrandKit       *SectionMeta
	DateTime       *SectionMeta
	DomainSettings *SectionMeta
	EmailSettings  *SectionMeta
	Favicon        *SectionMeta
	StoreDetails   *SectionMeta
	StoreLogo      *SectionMeta
}

// SiteSettings is the facade's body or one section's; ToMap renders the facade PATCH, not a section's.
type SiteSettings struct {
	Administration *Administration
	BrandKit       *BrandKit
	DateTime       *DateTime
	DomainSettings *DomainSettings
	EmailSettings  *EmailSettings
	Favicon        *Favicon
	StoreDetails   *StoreDetails
	StoreLogo      *StoreLogo
	GroupMeta      *GroupMeta
	// FieldMetadata is the raw field_metadata block.
	FieldMetadata any
}

func (SiteSettings) EnvelopeKey() string { return "site_settings" }

func (s SiteSettings) ToMap() map[string]any { return dto.ToMap(s) }

// SiteSettingsShape is the typed view of the facade's body and of one section's.
func SiteSettingsShape(decoded any) any { return dto.Read[SiteSettings](decoded) }
