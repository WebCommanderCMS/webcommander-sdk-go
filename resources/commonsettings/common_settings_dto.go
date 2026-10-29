package commonsettings

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// Setting is one entry of a settings group, for a write's lists; a read's entries are maps with the same keys.
type Setting struct {
	ConfigKey string
	// Value is a pointer, so an empty text is sent rather than omitted; a number written is stored as text.
	Value *string
}

// CommonSettings is what the read and both writes answer: Settings maps a group name to a list of {"config_key", "value"} maps, every value text and a masked one "***". The notifications group alone is a name to bool map, so the groups stay untyped.
type CommonSettings struct {
	Settings any
}

// CommonSettingsInput is the list form both writes take, a group name to its Setting entries. A group may instead be {config_key: value} or {name: {"config_key", "value"}}; send those as a map or a JSON string.
type CommonSettingsInput struct {
	Settings map[string][]*Setting
}

// ToMap renders the request body; dto.ToMap does not walk a map's values, so each group's entries are rendered here.
func (i CommonSettingsInput) ToMap() map[string]any {
	groups := make(map[string]any, len(i.Settings))
	for group, entries := range i.Settings {
		rendered := make([]any, len(entries))
		for at, entry := range entries {
			rendered[at] = dto.ToMap(entry)
		}
		groups[group] = rendered
	}
	return map[string]any{CommonSettingsKey: groups}
}

// CustomerFieldConfigs is the registration form's fields_configs: six keys for each of the store's 20 fields. Order reads as text; Active and Required are booleans here, though the settings read answers them as text.
type CustomerFieldConfigs struct {
	FirstNameOrder            *string
	FirstNameKey              *string
	FirstNameLabel            *string
	FirstNamePlaceholder      *string
	FirstNameActive           *bool
	FirstNameRequired         *bool
	LastNameOrder             *string
	LastNameKey               *string
	LastNameLabel             *string
	LastNamePlaceholder       *string
	LastNameActive            *bool
	LastNameRequired          *bool
	AddressLine1Order         *string
	AddressLine1Key           *string
	AddressLine1Label         *string
	AddressLine1Placeholder   *string
	AddressLine1Active        *bool
	AddressLine1Required      *bool
	AddressLine2Order         *string
	AddressLine2Key           *string
	AddressLine2Label         *string
	AddressLine2Placeholder   *string
	AddressLine2Active        *bool
	AddressLine2Required      *bool
	PostCodeOrder             *string
	PostCodeKey               *string
	PostCodeLabel             *string
	PostCodePlaceholder       *string
	PostCodeActive            *bool
	PostCodeRequired          *bool
	CityOrder                 *string
	CityKey                   *string
	CityLabel                 *string
	CityPlaceholder           *string
	CityActive                *bool
	CityRequired              *bool
	PhoneOrder                *string
	PhoneKey                  *string
	PhoneLabel                *string
	PhonePlaceholder          *string
	PhoneActive               *bool
	PhoneRequired             *bool
	MobileOrder               *string
	MobileKey                 *string
	MobileLabel               *string
	MobilePlaceholder         *string
	MobileActive              *bool
	MobileRequired            *bool
	FaxOrder                  *string
	FaxKey                    *string
	FaxLabel                  *string
	FaxPlaceholder            *string
	FaxActive                 *bool
	FaxRequired               *bool
	EmailOrder                *string
	EmailKey                  *string
	EmailLabel                *string
	EmailPlaceholder          *string
	EmailActive               *bool
	EmailRequired             *bool
	ConfirmEmailOrder         *string
	ConfirmEmailKey           *string
	ConfirmEmailLabel         *string
	ConfirmEmailPlaceholder   *string
	ConfirmEmailActive        *bool
	ConfirmEmailRequired      *bool
	PasswordOrder             *string
	PasswordKey               *string
	PasswordLabel             *string
	PasswordPlaceholder       *string
	PasswordActive            *bool
	PasswordRequired          *bool
	RetypePasswordOrder       *string
	RetypePasswordKey         *string
	RetypePasswordLabel       *string
	RetypePasswordPlaceholder *string
	RetypePasswordActive      *bool
	RetypePasswordRequired    *bool
	CountryOrder              *string
	CountryKey                *string
	CountryLabel              *string
	CountryPlaceholder        *string
	CountryActive             *bool
	CountryRequired           *bool
	CustomerTypeOrder         *string
	CustomerTypeKey           *string
	CustomerTypeLabel         *string
	CustomerTypePlaceholder   *string
	CustomerTypeActive        *bool
	CustomerTypeRequired      *bool
	SexOrder                  *string
	SexKey                    *string
	SexLabel                  *string
	SexPlaceholder            *string
	SexActive                 *bool
	SexRequired               *bool
	ABNOrder                  *string
	ABNKey                    *string
	ABNLabel                  *string
	ABNPlaceholder            *string
	ABNActive                 *bool
	ABNRequired               *bool
	ABNBranchOrder            *string
	ABNBranchKey              *string
	ABNBranchLabel            *string
	ABNBranchPlaceholder      *string
	ABNBranchActive           *bool
	ABNBranchRequired         *bool
	CompanyNameOrder          *string
	CompanyNameKey            *string
	CompanyNameLabel          *string
	CompanyNamePlaceholder    *string
	CompanyNameActive         *bool
	CompanyNameRequired       *bool
	DateOfBirthOrder          *string
	DateOfBirthKey            *string
	DateOfBirthLabel          *string
	DateOfBirthPlaceholder    *string
	DateOfBirthActive         *bool
	DateOfBirthRequired       *bool
}

// CustomerFields is the registration form: Fields lists the active fields' keys in their order, and FieldsConfigs configures every field, inactive ones included.
type CustomerFields struct {
	Fields        []*string
	FieldsConfigs *CustomerFieldConfigs
}

// CommonSettingsShape is the typed view of a body shaped {"settings": {…}}.
func CommonSettingsShape(decoded any) any { return dto.Read[CommonSettings](decoded) }

// CustomerFieldsShape is the typed view of the registration form's body.
func CustomerFieldsShape(decoded any) any { return dto.Read[CustomerFields](decoded) }
