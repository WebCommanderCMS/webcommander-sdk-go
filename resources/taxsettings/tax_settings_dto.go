package taxsettings

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// TaxRegion is the country tax the store charges by default, as tax_region reads it.
type TaxRegion struct {
	CountryTaxID *int64
	CountryCode  *string
	CountryName  *string
}

// TaxProfileReference is a custom tax profile as a read names it: its id and name, or null.
type TaxProfileReference struct {
	ID   *int64
	Name *string
}

// TaxSettingsAreaMeta is one group_meta entry: what the area is, whether it can be written, and its path.
type TaxSettingsAreaMeta struct {
	Kind     *string
	Writable *bool
	Href     *string
}

// TaxSettingsGroupMeta is the group_meta block: the settings area and the two tax collections beside it.
type TaxSettingsGroupMeta struct {
	TaxSettings *TaxSettingsAreaMeta
	TaxProfiles *TaxSettingsAreaMeta
	TaxRules    *TaxSettingsAreaMeta
}

// TaxSettings is the nine keys under tax_settings, with GroupMeta and, when asked for, FieldMetadata.
type TaxSettings struct {
	TaxRateForCustomers  *string
	TaxRegion            *TaxRegion
	CustomTaxProfile     *TaxProfileReference
	IsPriceShownWithTax  *bool
	PricesInCatalog      *string
	CalculatedBy         *string
	IsOutOfRegionTaxFree *bool
	TaxMessage           *string
	BasePriceRounding    *string
	GroupMeta            *TaxSettingsGroupMeta
	FieldMetadata        any
}

// EnvelopeKey is the key the read answers under, and the write's answer too.
func (TaxSettings) EnvelopeKey() string { return TaxSettingsKey }

// TaxRegionInput is tax_region as a write takes it: a CountryCode (with an optional StateCode) or a CountryTaxID.
type TaxRegionInput struct {
	CountryTaxID *int64
	CountryCode  string
	StateCode    string
}

// TaxSettingsInput is the body PUT and PATCH take, wrapped in tax_settings; both merge.
// The profile is written as CustomTaxProfileID; calculated_by is read-only.
type TaxSettingsInput struct {
	TaxRateForCustomers string
	TaxRegion           *TaxRegionInput
	CustomTaxProfileID  *int64
	// Pointers, so an unset flag is omitted rather than cleared by the merge.
	IsPriceShownWithTax  *bool
	PricesInCatalog      string
	IsOutOfRegionTaxFree *bool
	TaxMessage           string
	BasePriceRounding    string
}

// EnvelopeKey is the tax_settings wrapper, the same key the read answers under.
func (TaxSettingsInput) EnvelopeKey() string { return TaxSettingsWriteKey }

// ToMap renders the request body.
func (t TaxSettingsInput) ToMap() map[string]any { return dto.ToMap(t) }

// TaxSettingsShape is the typed view of a body shaped {"tax_settings": {…}}.
func TaxSettingsShape(decoded any) any { return dto.Read[TaxSettings](decoded) }
