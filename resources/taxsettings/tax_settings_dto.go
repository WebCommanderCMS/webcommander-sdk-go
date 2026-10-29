package taxsettings

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// TaxSettings is the nine keys the facade sends, flat; DefaultTaxProfile follows the default tax profile.
type TaxSettings struct {
	ConfigurationType  *string
	SystemDefaultTax   *bool
	DefaultCountry     *string
	CustomTaxProfileID *int64
	DefaultTaxProfile  *int64
	ShowPriceWithTax   *bool
	IsPriceWithTax     *bool
	DefaultTaxCode     *string
	TaxMessage         *string
}

// EnvelopeKey is the key the read answers under, and the write's answer too; the write takes another.
func (TaxSettings) EnvelopeKey() string { return TaxSettingsKey }

// TaxSettingsInput is the PUT's body, wrapped in settings; the write merges.
type TaxSettingsInput struct {
	ConfigurationType string
	// Pointers, so an unset flag or number is omitted rather than cleared by the merge.
	SystemDefaultTax   *bool
	DefaultCountry     string
	CustomTaxProfileID *int64
	DefaultTaxProfile  *int64
	ShowPriceWithTax   *bool
	IsPriceWithTax     *bool
	DefaultTaxCode     string
	TaxMessage         string
}

// EnvelopeKey is the settings wrapper, not the key the read answers under.
func (TaxSettingsInput) EnvelopeKey() string { return TaxSettingsWriteKey }

// ToMap renders the request body.
func (t TaxSettingsInput) ToMap() map[string]any { return dto.ToMap(t) }

// TaxSettingsShape is the typed view of a body shaped {"tax_settings": {…}}.
func TaxSettingsShape(decoded any) any { return dto.Read[TaxSettings](decoded) }
