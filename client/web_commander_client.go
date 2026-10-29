package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net"
	neturl "net/url"
	"os"
	"reflect"
	"slices"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/auth"
	wccfg "github.com/WebCommanderCMS/webcommander-sdk-go/v4/config"
	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/redact"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/aliases"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/blogcategories"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/blogs"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/carts"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/categories"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/cities"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/commonsettings"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/countries"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/currencies"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/customergroups"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/customers"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/discountprofiles"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/email"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/inventorylocations"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/ordercustomfields"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/orderflags"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/orders"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/pages"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/paymentgateways"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/products"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/settings/animationeffects"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/settings/customerportal"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/settings/ecommercesettings"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/settings/pagesettings"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/settings/pricedisplaysettings"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/settings/productlisting"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/settings/searchsettings"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/settings/seoconfiguration"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/settings/shippingsettings"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/settings/sitesettings"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/shippingclasses"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/shippingmethods"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/taxcodes"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/taxcountrytaxes"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/taxcustomprofiles"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/taxprofiles"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/taxrules"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/taxsettings"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/taxzones"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/transactions"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/trash"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/resources/variations"
)

type WebCommanderClient struct {
	config        wccfg.SdkConfig
	transport     *http.HTTPTransport
	tokenProvider *auth.TokenProvider
	log           func(string)

	// ownsTransport is true when no WithTransport was applied, so Close releases the transport.
	ownsTransport bool

	// closeMutex guards logFile, which Close read-modify-writes.
	closeMutex sync.Mutex
	logFile    *os.File // the LogFile sink, released by Close; nil when unset

	// closed is the lifetime latch Close sets and refuseIfClosed reads.
	closed atomic.Bool

	// lifetime is the window between a send being accepted and that send being finished.
	lifetime sync.RWMutex

	adminCustomers       lazy[*customers.AdminCustomers]
	seoConfiguration     lazy[*seoconfiguration.SeoConfigurationSettings]
	searchResultPage     lazy[*pagesettings.SearchResultPageSettings]
	categoryPage         lazy[*pagesettings.CategoryPageSettings]
	shoppingCart         lazy[*pagesettings.ShoppingCartSettings]
	productPage          lazy[*pagesettings.ProductPageSettings]
	checkoutPage         lazy[*pagesettings.CheckoutPageSettings]
	priceDisplaySettings lazy[*pricedisplaysettings.PriceDisplaySettingsResource]
	animationEffects     lazy[*animationeffects.AnimationEffectsSettings]
	productListing       lazy[*productlisting.ProductListingSettingsResource]
	customerPortal       lazy[*customerportal.CustomerPortalSettings]
	siteSettings         lazy[*sitesettings.SiteSettingsResource]
	shippingSettings     lazy[*shippingsettings.ShippingSettingsResource]
	searchSettings       lazy[*searchsettings.SearchSettingsResource]
	ecommerceSettings    lazy[*ecommercesettings.EcommerceSettingsResource]
	orderFlags           lazy[*orderflags.OrderFlags]
	orders               lazy[*orders.Orders]
	inventoryLocations   lazy[*inventorylocations.InventoryLocations]
	products             lazy[*products.Products]
	aliases              lazy[*aliases.Aliases]
	carts                lazy[*carts.Carts]
	blogs                lazy[*blogs.Blogs]
	orderCustomFields    lazy[*ordercustomfields.OrderCustomFields]
	blogCategories       lazy[*blogcategories.BlogCategories]
	categories           lazy[*categories.Categories]
	cities               lazy[*cities.Cities]
	commonSettings       lazy[*commonsettings.CommonSettingsResource]
	countries            lazy[*countries.Countries]
	currencies           lazy[*currencies.Currencies]
	paymentGateways      lazy[*paymentgateways.PaymentGateways]
	shippingClasses      lazy[*shippingclasses.ShippingClasses]
	shippingMethods      lazy[*shippingmethods.ShippingMethods]
	taxCodes             lazy[*taxcodes.TaxCodes]
	taxProfiles          lazy[*taxprofiles.TaxProfiles]
	taxRules             lazy[*taxrules.TaxRules]
	taxZones             lazy[*taxzones.TaxZones]
	taxSettings          lazy[*taxsettings.TaxSettingsResource]
	taxCountryTaxes      lazy[*taxcountrytaxes.TaxCountryTaxes]
	taxCustomProfiles    lazy[*taxcustomprofiles.TaxCustomProfiles]
	discountProfiles     lazy[*discountprofiles.DiscountProfiles]
	customerGroups       lazy[*customergroups.CustomerGroups]
	variations           lazy[*variations.VariationsResource]
	transactions         lazy[*transactions.Transactions]
	trash                lazy[*trash.Trash]
	email                lazy[*email.Email]
	pages                lazy[*pages.Pages]
}

// lazy is a resource built once, on first use.
type lazy[T any] struct {
	once  sync.Once
	value T
}

func (l *lazy[T]) get(build func() T) T {
	l.once.Do(func() { l.value = build() })
	return l.value
}

const (
	// LoggerName prefixes each log-file line, followed by the escaped client name.
	LoggerName = "webcommander"

	apiVersionParam   = "api_version"
	tokenFailureError = "unauthorized"

	// authorizationHeader is the canonical spelling the bearer token is always sent under.
	authorizationHeader = "Authorization"

	// noResponseStatus stands where a status would go; text, so it cannot look like a real code.
	noResponseStatus = "no response"

	// The record levels. LogLevel gates which of them the log hook receives.
	levelDebug   = "DEBUG"
	levelInfo    = "INFO"
	levelWarning = "WARNING"
	levelError   = "ERROR"
)

var tokenFailureMessages = map[string]bool{"invalid_access_token": true, "invalid_request": true}

var logRecordPrefixes = map[string]string{
	levelDebug: "DEBUG ", levelInfo: "INFO ", levelWarning: "WARN ", levelError: "ERROR ",
}

func logAllows(config wccfg.SdkConfig, level string) bool {
	return wccfg.LogLevelRank(level) >= wccfg.LogLevelRank(config.LogLevel)
}

type clientOptions struct {
	transport     *http.HTTPTransport
	tokenStore    auth.TokenStore
	tokenProvider *auth.TokenProvider
	retryPolicy   *http.RetryPolicy
	log           func(string)
}

// ClientOption configures a client at construction — the injection points, mainly for testing.
type ClientOption func(*clientOptions)

// WithTransport injects a transport in place of the one the config would build.
func WithTransport(t *http.HTTPTransport) ClientOption {
	return func(applied *clientOptions) { applied.transport = t }
}

func WithTokenStore(s auth.TokenStore) ClientOption {
	return func(applied *clientOptions) { applied.tokenStore = s }
}

// WithTokenProvider injects a built token provider, so several clients can share one token cache.
func WithTokenProvider(p *auth.TokenProvider) ClientOption {
	return func(applied *clientOptions) { applied.tokenProvider = p }
}

// WithRetryPolicy injects a retry policy (built with NewRetryPolicy) into the default transport.
func WithRetryPolicy(p *http.RetryPolicy) ClientOption {
	return func(applied *clientOptions) { applied.retryPolicy = p }
}

// WithLogger installs a plain log hook: a func(string) called with one line per exchange.
func WithLogger(f func(string)) ClientOption { return func(applied *clientOptions) { applied.log = f } }

func NewClient(config wccfg.SdkConfig, options ...ClientOption) (*WebCommanderClient, error) {
	// Validated raw, before Normalised: normaliseAPIVersion is not idempotent.
	if err := config.Validate(); err != nil {
		return nil, err
	}
	normalisedConfig := config.Normalised()
	// Resolved here, so a later change to WC_RESPONSE_FORMAT cannot reach this client.
	normalisedConfig.ResponseFormat = normalisedConfig.ResponseFormatValue()
	applied := &clientOptions{}
	for _, option := range options {
		option(applied)
	}
	if applied.transport != nil && applied.retryPolicy != nil {
		return nil, wcerrors.NewUsageError("pass WithTransport or WithRetryPolicy, not both")
	}

	logFn, logFile, err := attachLogFile(normalisedConfig, applied.log)
	if err != nil {
		return nil, err
	}
	closeOnError := func() {
		if logFile != nil {
			_ = logFile.Close()
		}
	}

	warnAboutInsecureTransport(normalisedConfig, logFn)

	transport := applied.transport

	ownsTransport := transport == nil
	if transport == nil {
		transport, err = http.NewHTTPTransport(normalisedConfig,
			http.WithProgress(progressFunc(normalisedConfig, logFn)),
			http.WithRetryPolicy(applied.retryPolicy))
		if err != nil {
			closeOnError()
			return nil, err
		}
	}
	tokenProvider := applied.tokenProvider
	if tokenProvider == nil {
		store := applied.tokenStore
		if store == nil {
			store, err = auth.BuildDefaultTokenStore(normalisedConfig)
			if err != nil {
				closeOnError()
				return nil, err
			}
		}
		tokenProvider = auth.NewTokenProvider(normalisedConfig, transport, store)
	}
	return &WebCommanderClient{config: normalisedConfig, transport: transport, ownsTransport: ownsTransport,
		tokenProvider: tokenProvider, log: logFn, logFile: logFile}, nil
}

// loggerSafe flattens a client name into one log-name segment, so no instance nests under another.
func loggerSafe(clientName string) string {
	var escaped strings.Builder
	for index := 0; index < len(clientName); index++ {
		c := clientName[index]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9',
			c == '-', c == '_', c == '~':
			escaped.WriteByte(c)
		default:
			fmt.Fprintf(&escaped, "%%%02X", c)
		}
	}
	return escaped.String()
}

// attachLogFile opens config.LogFile for append and returns a hook that also writes timestamped lines there.
func attachLogFile(config wccfg.SdkConfig, base func(string)) (func(string), *os.File, error) {
	if config.LogFile == "" {
		return base, nil, nil
	}
	// 0600, like the token cache: the exchange lines logged here carry customer ids in their URLs.
	f, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		// The path and the OS error go to CauseText, not the message callers are asked to retain.
		refusal := wcerrors.NewConfigError(fmt.Sprintf(
			"logFile: %s could not be opened for appending", redact.DescribePath(config.LogFile)))
		refusal.Cause = err
		refusal.CauseText = fmt.Sprintf("%s: %v", config.LogFile, err)
		return nil, nil, refusal
	}
	name := loggerSafe(config.Auth.ClientName)
	sink := func(line string) {
		fmt.Fprintf(f, "%s %s.%s %s\n", time.Now().Format("2006-01-02 15:04:05"), LoggerName, name, line)
		if base != nil {
			base(line)
		}
	}
	return sink, f, nil
}

// emittedTransportWarnings dedupes the stderr copy of each warning per process; the log hook gets all.
var (
	emittedTransportWarningsMu sync.Mutex
	emittedTransportWarnings   = map[string]bool{}
)

// warnAboutInsecureTransport warns when TLS verification is off or the base URL is http:// to a remote host.
func warnAboutInsecureTransport(config wccfg.SdkConfig, log func(string)) {
	emit := func(message string) {
		w := wcerrors.NewInsecureTransportWarning(message)
		if log != nil && logAllows(config, levelWarning) {
			log(logRecordPrefixes[levelWarning] + w.Error())
		}
		emittedTransportWarningsMu.Lock()
		seen := emittedTransportWarnings[message]
		emittedTransportWarnings[message] = true
		emittedTransportWarningsMu.Unlock()
		if !seen {
			fmt.Fprintln(os.Stderr, "InsecureTransportWarning: "+w.Error())
		}
	}
	if !config.VerifySslEnabled() {
		emit("TLS verification is disabled (verifySsl=false); responses cannot be trusted")
	}
	u, err := neturl.Parse(config.Auth.BaseURL)
	if err != nil || u.Scheme != "http" || isLoopbackHost(u.Hostname()) {
		return
	}
	emit("baseUrl uses http://, so every request to a host that is not this machine crosses the " +
		"network in cleartext — the bearer token, and the credentials in the token request itself. " +
		"Use https:// unless this is a local development instance")
}

var loopbackHostnames = map[string]bool{"localhost": true, "localhost.localdomain": true}

func isLoopbackHost(host string) bool {
	if loopbackHostnames[strings.ToLower(host)] {
		return true
	}
	if ip := net.ParseIP(host); ip != nil {
		return ip.IsLoopback()
	}
	return false
}

// Config returns a deep copy of the configuration this client was built with.
func (c *WebCommanderClient) Config() wccfg.SdkConfig { return c.config.DeepCopy() }

// AdminCustomers is built and cached on first use.
func (c *WebCommanderClient) AdminCustomers() *customers.AdminCustomers {
	return c.adminCustomers.get(func() *customers.AdminCustomers { return customers.New(c) })
}

// SeoConfigurationSettings is built and cached on first use.
func (c *WebCommanderClient) SeoConfigurationSettings() *seoconfiguration.SeoConfigurationSettings {
	return c.seoConfiguration.get(func() *seoconfiguration.SeoConfigurationSettings { return seoconfiguration.New(c) })
}

// SearchResultPageSettings is built and cached on first use.
func (c *WebCommanderClient) SearchResultPageSettings() *pagesettings.SearchResultPageSettings {
	return c.searchResultPage.get(func() *pagesettings.SearchResultPageSettings { return pagesettings.New(c) })
}

// CategoryPageSettings is built and cached on first use.
func (c *WebCommanderClient) CategoryPageSettings() *pagesettings.CategoryPageSettings {
	return c.categoryPage.get(func() *pagesettings.CategoryPageSettings { return pagesettings.NewCategoryPageSettings(c) })
}

// ShoppingCartSettings is built and cached on first use.
func (c *WebCommanderClient) ShoppingCartSettings() *pagesettings.ShoppingCartSettings {
	return c.shoppingCart.get(func() *pagesettings.ShoppingCartSettings { return pagesettings.NewShoppingCartSettings(c) })
}

// ProductPageSettings is built and cached on first use.
func (c *WebCommanderClient) ProductPageSettings() *pagesettings.ProductPageSettings {
	return c.productPage.get(func() *pagesettings.ProductPageSettings { return pagesettings.NewProductPageSettings(c) })
}

// CheckoutPageSettings is built and cached on first use.
func (c *WebCommanderClient) CheckoutPageSettings() *pagesettings.CheckoutPageSettings {
	return c.checkoutPage.get(func() *pagesettings.CheckoutPageSettings { return pagesettings.NewCheckoutPageSettings(c) })
}

// PriceDisplaySettings is built and cached on first use.
func (c *WebCommanderClient) PriceDisplaySettings() *pricedisplaysettings.PriceDisplaySettingsResource {
	return c.priceDisplaySettings.get(func() *pricedisplaysettings.PriceDisplaySettingsResource { return pricedisplaysettings.New(c) })
}

// AnimationEffectsSettings is built and cached on first use.
func (c *WebCommanderClient) AnimationEffectsSettings() *animationeffects.AnimationEffectsSettings {
	return c.animationEffects.get(func() *animationeffects.AnimationEffectsSettings { return animationeffects.New(c) })
}

// ProductListingSettings is built and cached on first use.
func (c *WebCommanderClient) ProductListingSettings() *productlisting.ProductListingSettingsResource {
	return c.productListing.get(func() *productlisting.ProductListingSettingsResource { return productlisting.New(c) })
}

// CustomerPortalSettings is built and cached on first use.
func (c *WebCommanderClient) CustomerPortalSettings() *customerportal.CustomerPortalSettings {
	return c.customerPortal.get(func() *customerportal.CustomerPortalSettings { return customerportal.New(c) })
}

// SiteSettings is built and cached on first use.
func (c *WebCommanderClient) SiteSettings() *sitesettings.SiteSettingsResource {
	return c.siteSettings.get(func() *sitesettings.SiteSettingsResource { return sitesettings.New(c) })
}

// ShippingSettings is built and cached on first use.
func (c *WebCommanderClient) ShippingSettings() *shippingsettings.ShippingSettingsResource {
	return c.shippingSettings.get(func() *shippingsettings.ShippingSettingsResource { return shippingsettings.New(c) })
}

// SearchSettings is built and cached on first use.
func (c *WebCommanderClient) SearchSettings() *searchsettings.SearchSettingsResource {
	return c.searchSettings.get(func() *searchsettings.SearchSettingsResource { return searchsettings.New(c) })
}

// EcommerceSettings is built and cached on first use.
func (c *WebCommanderClient) EcommerceSettings() *ecommercesettings.EcommerceSettingsResource {
	return c.ecommerceSettings.get(func() *ecommercesettings.EcommerceSettingsResource { return ecommercesettings.New(c) })
}

// OrderFlags is built and cached on first use.
func (c *WebCommanderClient) OrderFlags() *orderflags.OrderFlags {
	return c.orderFlags.get(func() *orderflags.OrderFlags { return orderflags.New(c) })
}

// Orders is built and cached on first use.
func (c *WebCommanderClient) Orders() *orders.Orders {
	return c.orders.get(func() *orders.Orders { return orders.New(c) })
}

// InventoryLocations is built and cached on first use.
func (c *WebCommanderClient) InventoryLocations() *inventorylocations.InventoryLocations {
	return c.inventoryLocations.get(func() *inventorylocations.InventoryLocations { return inventorylocations.New(c) })
}

// Products is built and cached on first use.
func (c *WebCommanderClient) Products() *products.Products {
	return c.products.get(func() *products.Products { return products.New(c) })
}

// Aliases is built and cached on first use.
func (c *WebCommanderClient) Aliases() *aliases.Aliases {
	return c.aliases.get(func() *aliases.Aliases { return aliases.New(c) })
}

// Carts is built and cached on first use.
func (c *WebCommanderClient) Carts() *carts.Carts {
	return c.carts.get(func() *carts.Carts { return carts.New(c) })
}

// Blogs is built and cached on first use.
func (c *WebCommanderClient) Blogs() *blogs.Blogs {
	return c.blogs.get(func() *blogs.Blogs { return blogs.New(c) })
}

// OrderCustomFields is built and cached on first use.
func (c *WebCommanderClient) OrderCustomFields() *ordercustomfields.OrderCustomFields {
	return c.orderCustomFields.get(func() *ordercustomfields.OrderCustomFields { return ordercustomfields.New(c) })
}

// BlogCategories is built and cached on first use.
func (c *WebCommanderClient) BlogCategories() *blogcategories.BlogCategories {
	return c.blogCategories.get(func() *blogcategories.BlogCategories { return blogcategories.New(c) })
}

// Categories is built and cached on first use.
func (c *WebCommanderClient) Categories() *categories.Categories {
	return c.categories.get(func() *categories.Categories { return categories.New(c) })
}

// Cities is built and cached on first use.
func (c *WebCommanderClient) Cities() *cities.Cities {
	return c.cities.get(func() *cities.Cities { return cities.New(c) })
}

// CommonSettings is built and cached on first use.
func (c *WebCommanderClient) CommonSettings() *commonsettings.CommonSettingsResource {
	return c.commonSettings.get(func() *commonsettings.CommonSettingsResource { return commonsettings.New(c) })
}

// Countries is built and cached on first use.
func (c *WebCommanderClient) Countries() *countries.Countries {
	return c.countries.get(func() *countries.Countries { return countries.New(c) })
}

// Currencies is built and cached on first use.
func (c *WebCommanderClient) Currencies() *currencies.Currencies {
	return c.currencies.get(func() *currencies.Currencies { return currencies.New(c) })
}

// PaymentGateways is built and cached on first use.
func (c *WebCommanderClient) PaymentGateways() *paymentgateways.PaymentGateways {
	return c.paymentGateways.get(func() *paymentgateways.PaymentGateways { return paymentgateways.New(c) })
}

// ShippingClasses is built and cached on first use.
func (c *WebCommanderClient) ShippingClasses() *shippingclasses.ShippingClasses {
	return c.shippingClasses.get(func() *shippingclasses.ShippingClasses { return shippingclasses.New(c) })
}

// ShippingMethods is built and cached on first use.
func (c *WebCommanderClient) ShippingMethods() *shippingmethods.ShippingMethods {
	return c.shippingMethods.get(func() *shippingmethods.ShippingMethods { return shippingmethods.New(c) })
}

// TaxCodes is built and cached on first use.
func (c *WebCommanderClient) TaxCodes() *taxcodes.TaxCodes {
	return c.taxCodes.get(func() *taxcodes.TaxCodes { return taxcodes.New(c) })
}

// TaxProfiles is built and cached on first use.
func (c *WebCommanderClient) TaxProfiles() *taxprofiles.TaxProfiles {
	return c.taxProfiles.get(func() *taxprofiles.TaxProfiles { return taxprofiles.New(c) })
}

// TaxRules is built and cached on first use.
func (c *WebCommanderClient) TaxRules() *taxrules.TaxRules {
	return c.taxRules.get(func() *taxrules.TaxRules { return taxrules.New(c) })
}

// TaxZones is built and cached on first use.
func (c *WebCommanderClient) TaxZones() *taxzones.TaxZones {
	return c.taxZones.get(func() *taxzones.TaxZones { return taxzones.New(c) })
}

// TaxSettings is built and cached on first use.
func (c *WebCommanderClient) TaxSettings() *taxsettings.TaxSettingsResource {
	return c.taxSettings.get(func() *taxsettings.TaxSettingsResource { return taxsettings.New(c) })
}

// TaxCountryTaxes is built and cached on first use.
func (c *WebCommanderClient) TaxCountryTaxes() *taxcountrytaxes.TaxCountryTaxes {
	return c.taxCountryTaxes.get(func() *taxcountrytaxes.TaxCountryTaxes { return taxcountrytaxes.New(c) })
}

// TaxCustomProfiles is built and cached on first use.
func (c *WebCommanderClient) TaxCustomProfiles() *taxcustomprofiles.TaxCustomProfiles {
	return c.taxCustomProfiles.get(func() *taxcustomprofiles.TaxCustomProfiles { return taxcustomprofiles.New(c) })
}

// DiscountProfiles is built and cached on first use.
func (c *WebCommanderClient) DiscountProfiles() *discountprofiles.DiscountProfiles {
	return c.discountProfiles.get(func() *discountprofiles.DiscountProfiles { return discountprofiles.New(c) })
}

// CustomerGroups is built and cached on first use.
func (c *WebCommanderClient) CustomerGroups() *customergroups.CustomerGroups {
	return c.customerGroups.get(func() *customergroups.CustomerGroups { return customergroups.New(c) })
}

// Variations is built and cached on first use.
func (c *WebCommanderClient) Variations() *variations.VariationsResource {
	return c.variations.get(func() *variations.VariationsResource { return variations.New(c) })
}

// Transactions is built and cached on first use.
func (c *WebCommanderClient) Transactions() *transactions.Transactions {
	return c.transactions.get(func() *transactions.Transactions { return transactions.New(c) })
}

// Trash is built and cached on first use.
func (c *WebCommanderClient) Trash() *trash.Trash {
	return c.trash.get(func() *trash.Trash { return trash.New(c) })
}

// Email is built and cached on first use.
func (c *WebCommanderClient) Email() *email.Email {
	return c.email.get(func() *email.Email { return email.New(c) })
}

// Pages is built and cached on first use.
func (c *WebCommanderClient) Pages() *pages.Pages {
	return c.pages.get(func() *pages.Pages { return pages.New(c) })
}

// ListResourceNames returns every resource reachable on this client, e.g. ["admin_customers"].
func (c *WebCommanderClient) ListResourceNames() []string {
	return RegisteredResourceNames()
}

// SendApiRequest sends one request and returns the response — every argument on the options.
func (c *WebCommanderClient) SendApiRequest(ctx context.Context, httpMethod, path string, options *http.RequestOptions) (*http.ApiResponse, error) {
	finished, err := c.acceptSend()
	if err != nil {
		return nil, err
	}
	defer finished()

	if err := http.RefuseIfBudgetSpent(ctx, "before the request was prepared"); err != nil {
		return nil, err
	}
	// Frozen before anything reads it, so later caller writes cannot reach this call.
	var requestOptions http.RequestOptions
	if options != nil {
		requestOptions = options.Frozen()
	}
	method, err := http.NormaliseMethod(httpMethod)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(path) == "" {
		return nil, wcerrors.NewUsageError(fmt.Sprintf("path must name an endpoint, got '%s'", path))
	}
	responseFormat, err := c.resolveResponseFormat(requestOptions.ResponseFormat)
	if err != nil {
		return nil, err
	}

	headers := c.mergeRequestHeaders(requestOptions.Headers)
	if err := http.CheckRequestHeaders(headers); err != nil {
		return nil, err
	}
	resolvedPath, err := http.SubstitutePathParams(path, c.seedApiVersion(path, requestOptions.PathParams))
	if err != nil {
		return nil, err
	}
	url, err := http.BuildURL(c.config.Auth.BaseURL, resolvedPath)
	if err != nil {
		return nil, err
	}
	queryParams, err := http.BuildQueryParams(requestOptions.QueryParams)
	if err != nil {
		return nil, err
	}
	if err := rejectDuplicatedQueryParams(url, queryParams); err != nil {
		return nil, err
	}

	// Judged and encoded only after the target is settled: the order is load-bearing.
	wireJSONBody, err := validatePayload(&requestOptions)
	if err != nil {
		return nil, err
	}

	// Taken once here, so the report describes what was sent, not what the caller's map holds later.
	reportedBody := redact.RedactValue(requestOptions.ReportedBody())
	// The curl line needs the body kind: the transport, not the recorded headers, sets Content-Type.
	formEncoded := requestOptions.JSONBody == nil && requestOptions.FormData != nil

	materialised, err := http.TransportRequest{
		HTTPMethod:         method,
		URL:                url,
		EndpointTemplate:   endpointTemplateOf(path),
		QueryParams:        queryParams,
		JSONBody:           wireJSONBody,
		FormData:           requestOptions.FormData,
		Files:              requestOptions.Files,
		Headers:            headers,
		RequestInfoHeaders: redact.RedactMapping(headers),
		RequestInfoBody:    redact.RedactValue(summarisePayload(reportedBody, requestOptions.Files)),
		ResponseFormat:     responseFormat,
		DTOShape:           requestOptions.DTOShape,
	}.Materialised()
	if err != nil {
		return nil, err
	}
	reportedFiles := materialised.Files

	// Read again: a deadline that expired during preparation must not now buy a token and send.
	if err := http.RefuseIfBudgetSpent(ctx, "while the request was being prepared"); err != nil {
		return nil, err
	}

	accessToken := ""
	if !requestOptions.Unauthenticated {
		accessToken, err = c.tokenProvider.GetAccessToken(ctx)
		if err != nil {
			return nil, err
		}
	}

	response, err := c.transport.SendHTTPRequest(ctx, withBearerToken(materialised, accessToken))
	if err != nil {
		c.reportFailedSend(err, reportedBody, reportedFiles, formEncoded)
		return nil, err
	}

	if !requestOptions.Unauthenticated && meansExpiredToken(response, method) {
		// Reported before the replay, so a 401 -> refresh -> 200 sequence shows the 401.
		c.reportExchange(response, reportedBody, reportedFiles, formEncoded, levelWarning)
		fresh, refreshErr := c.tokenProvider.RefreshAccessToken(ctx, accessToken)
		if refreshErr != nil {
			return nil, refreshErr
		}
		response, err = c.transport.SendHTTPRequest(ctx, withBearerToken(materialised, fresh))
		if err != nil {
			c.reportFailedSend(err, reportedBody, reportedFiles, formEncoded)
			return nil, err
		}
	}
	c.reportExchange(response, reportedBody, reportedFiles, formEncoded, levelInfo)

	raise := c.config.RaiseOnError
	if requestOptions.RaiseOnError != nil {
		raise = *requestOptions.RaiseOnError
	}
	if raise {
		if statusError := response.RaiseForStatus(); statusError != nil {
			return response, statusError
		}
	}
	return response, nil
}

// SendGetRequest sends a GET.
func (c *WebCommanderClient) SendGetRequest(ctx context.Context, path string, options *http.RequestOptions) (*http.ApiResponse, error) {
	return c.SendApiRequest(ctx, "GET", path, options)
}

func (c *WebCommanderClient) SendPostRequest(ctx context.Context, path string, options *http.RequestOptions) (*http.ApiResponse, error) {
	return c.SendApiRequest(ctx, "POST", path, options)
}

func (c *WebCommanderClient) SendPutRequest(ctx context.Context, path string, options *http.RequestOptions) (*http.ApiResponse, error) {
	return c.SendApiRequest(ctx, "PUT", path, options)
}

func (c *WebCommanderClient) SendPatchRequest(ctx context.Context, path string, options *http.RequestOptions) (*http.ApiResponse, error) {
	return c.SendApiRequest(ctx, "PATCH", path, options)
}

func (c *WebCommanderClient) SendDeleteRequest(ctx context.Context, path string, options *http.RequestOptions) (*http.ApiResponse, error) {
	return c.SendApiRequest(ctx, "DELETE", path, options)
}

// Close releases the LogFile handle and the connections of a transport this client built.
func (c *WebCommanderClient) Close() (err error) {
	// Latched first, so a send racing this Close is refused rather than handed a closing transport.
	c.closed.Store(true)
	// Then the lifetime write lock, granted once every accepted send has returned.
	c.lifetime.Lock()
	defer c.lifetime.Unlock()
	defer func() {
		c.closeMutex.Lock()
		handle := c.logFile
		c.logFile = nil // a second Close is a no-op, not a spurious error
		c.closeMutex.Unlock()
		if handle == nil {
			return
		}
		if closeErr := handle.Close(); closeErr != nil && err == nil {
			refusal := wcerrors.NewConfigError(fmt.Sprintf(
				"logFile: %s failed to flush its last records to disk",
				redact.DescribePath(c.config.LogFile)))
			refusal.Cause = closeErr
			refusal.CauseText = fmt.Sprintf("%s: %v", c.config.LogFile, closeErr)
			err = refusal
		}
	}()
	if c.ownsTransport {
		c.transport.Close()
	}
	return nil
}

// acceptSend refuses a closed client, else takes the lifetime read lock the send holds until it returns.
func (c *WebCommanderClient) acceptSend() (func(), error) {
	if err := c.refuseIfClosed(); err != nil {
		return nil, err
	}
	c.lifetime.RLock()
	if err := c.refuseIfClosed(); err != nil {
		c.lifetime.RUnlock()
		return nil, err
	}
	return c.lifetime.RUnlock, nil
}

func (c *WebCommanderClient) refuseIfClosed() error {
	if !c.closed.Load() {
		return nil
	}
	return wcerrors.NewUsageError(fmt.Sprintf(
		"client '%s' is closed: Close() has already run, so this request was not sent — "+
			"build a new client rather than reusing a closed one",
		c.config.Auth.ClientName))
}

// withBearerToken applies one attempt's bearer token: the only change a stale-token replay may make.
func withBearerToken(materialised http.TransportRequest, accessToken string) http.TransportRequest {
	if accessToken == "" {
		return materialised
	}
	headers := make(map[string]string, len(materialised.Headers)+1)
	for name, value := range materialised.Headers {
		// Any caller spelling of Authorization is replaced, not sent beside the SDK's own.
		if !strings.EqualFold(name, authorizationHeader) {
			headers[name] = value
		}
	}
	headers[authorizationHeader] = "Bearer " + accessToken
	out := materialised
	out.Headers = headers
	out.RequestInfoHeaders = redact.RedactMapping(headers)
	return out
}

// resolveResponseFormat picks the per-call format if given, otherwise the client's own.
func (c *WebCommanderClient) resolveResponseFormat(perCall string) (string, error) {
	if perCall == "" {
		return c.config.ResponseFormat, nil
	}
	format := strings.ToLower(strings.TrimSpace(perCall))
	if !wccfg.IsResponseFormat(format) {
		return "", wcerrors.NewUsageError(fmt.Sprintf(
			"responseFormat: expected one of %s, got '%s'",
			strings.Join(wccfg.ResponseFormatNames(), ", "), perCall))
	}
	return format, nil
}

// mergeRequestHeaders merges SDK, config and caller headers case-insensitively; the caller's values win.
func (c *WebCommanderClient) mergeRequestHeaders(callerHeaders map[string]string) map[string]string {
	canonical := map[string]string{}   // lower -> canonical spelling
	byLowerName := map[string]string{} // lower -> value

	putHeader(canonical, byLowerName, "Accept", "application/json")
	putHeader(canonical, byLowerName, "User-Agent", c.config.UserAgent)
	for _, name := range slices.Sorted(maps.Keys(c.config.DefaultHeaders)) {
		putHeader(canonical, byLowerName, name, c.config.DefaultHeaders[name])
	}
	for _, name := range slices.Sorted(maps.Keys(callerHeaders)) {
		putHeader(canonical, byLowerName, name, callerHeaders[name])
	}
	merged := make(map[string]string, len(byLowerName))
	for lower, value := range byLowerName {
		merged[canonical[lower]] = value
	}
	return merged
}

func putHeader(canonical, byLowerName map[string]string, name, value string) {
	lower := strings.ToLower(name)
	// The last spelling seen wins; only the three SDK-owned names below are restored to canonical form.
	canonical[lower] = name
	switch lower {
	case "authorization":
		canonical[lower] = authorizationHeader
	case "accept":
		canonical[lower] = "Accept"
	case "user-agent":
		canonical[lower] = "User-Agent"
	}
	byLowerName[lower] = value
}

func (c *WebCommanderClient) seedApiVersion(path string, pathParams map[string]any) map[string]any {
	resolved := make(map[string]any, len(pathParams)+1)
	maps.Copy(resolved, pathParams)
	if slices.Contains(http.FindPathParamNames(path), apiVersionParam) {
		if _, ok := resolved[apiVersionParam]; !ok {
			resolved[apiVersionParam] = c.config.APIVersion
		}
	}
	return resolved
}

// validatePayload judges the payload before the token grant, so a bad body costs no OAuth2 round trip.
func validatePayload(requestOptions *http.RequestOptions) (any, error) {
	// A typed payload is a view over the map form, never a second wire format.
	if payload, isDTO := requestOptions.JSONBody.(http.DTOPayload); isDTO {
		rendered := payload.ToMap()
		if rendered == nil {
			return nil, wcerrors.NewUsageError(fmt.Sprintf(
				"jsonBody: %T rendered no body at all — a DTO payload must return the map to send, "+
					"and nil would send the request with no body", requestOptions.JSONBody))
		}

		requestOptions.JSONBody = http.SnapshotJSONValue(rendered)
	}
	if requestOptions.JSONBody != nil && requestOptions.FormData != nil {
		return nil, wcerrors.NewUsageError("pass jsonBody or formData, not both — they are different content types " +
			"and only one can be sent")
	}
	if requestOptions.JSONBody != nil && requestOptions.Files != nil {
		return nil, wcerrors.NewUsageError("pass jsonBody or files, not both — a multipart body replaces the JSON " +
			"one, so the JSON would be discarded unsent")
	}
	var wireJSONBody any
	// A string body is JSON the caller wrote, verified here and then sent byte for byte.
	if raw, isRawJSON := requestOptions.JSONBody.(string); isRawJSON {
		if err := validateRawJSONBody(raw); err != nil {
			return nil, err
		}
		wireJSONBody = raw
	} else if requestOptions.JSONBody != nil {
		switch reflect.ValueOf(requestOptions.JSONBody).Kind() {
		case reflect.Map, reflect.Slice, reflect.Array:
		default:
			return nil, wcerrors.NewUsageError(fmt.Sprintf(
				"jsonBody must be a mapping or list, got %T", requestOptions.JSONBody))
		}

		encoded, err := http.EncodeJSONBody(requestOptions.JSONBody)
		if err != nil {
			return nil, http.RefuseUnencodableJSONBody(err)
		}
		wireJSONBody = string(encoded)
	}
	// The same pre-flight for a multipart body, so the grant is not spent on one that cannot be built.
	if err := http.CheckFiles(requestOptions.Files); err != nil {
		return nil, err
	}
	if len(requestOptions.Files) > 0 {
		// Only a multipart body puts formData names into part headers, where a line break could forge one.
		if err := http.CheckFormDataNames(requestOptions.FormData); err != nil {
			return nil, err
		}
	}
	return wireJSONBody, nil
}

func validateRawJSONBody(jsonBody string) error {
	var decoded any
	if err := json.Unmarshal([]byte(jsonBody), &decoded); err != nil {
		// The byte offset, not the decoder's words, which quote the caller's own document.
		refusal := wcerrors.NewUsageError(fmt.Sprintf(
			"jsonBody: the text is not valid JSON%s", jsonSyntaxPosition(err)))
		refusal.Cause = err
		refusal.CauseText = err.Error()
		return refusal
	}
	switch decoded.(type) {
	case map[string]any, []any:
		return nil
	}
	return wcerrors.NewUsageError(fmt.Sprintf(
		"jsonBody: a JSON string must hold an object or an array, got %s", http.JSONKind(decoded)))
}

// jsonSyntaxPosition names the byte offset where a JSON document stopped parsing, when the decoder says.
func jsonSyntaxPosition(err error) string {
	var syntax *json.SyntaxError
	if errors.As(err, &syntax) {
		return fmt.Sprintf(" — it stops parsing at byte %d", syntax.Offset)
	}
	var unmarshalType *json.UnmarshalTypeError
	if errors.As(err, &unmarshalType) {
		return fmt.Sprintf(" — a %s where a value was expected, at byte %d",
			unmarshalType.Value, unmarshalType.Offset)
	}
	return ""
}

func rejectDuplicatedQueryParams(url string, queryParams map[string]any) error {
	if len(queryParams) == 0 {
		return nil
	}
	// Parsed with net/url, so a '?' in the fragment cannot fake a query and encoded names still clash.
	parsed, err := neturl.Parse(url)
	if err != nil {
		return nil // the URL builder has already vetted the origin
	}
	// ParseQuery decodes what it can; the decoded names are still the honest comparison set.
	existing, _ := neturl.ParseQuery(parsed.RawQuery)
	if len(existing) == 0 {
		return nil
	}
	var clash []string
	for k := range queryParams {
		if _, ok := existing[k]; ok {
			clash = append(clash, k)
		}
	}
	sort.Strings(clash)
	if len(clash) > 0 {
		return wcerrors.NewUsageError(fmt.Sprintf("queryParams %s already appear in the path's query string",
			strings.Join(clash, ", ")))
	}
	return nil
}

// meansExpiredToken reports whether a 401 means the access token rather than the request.
func meansExpiredToken(response *http.ApiResponse, method string) bool {
	if response.StatusCode != 401 {
		return false
	}
	// Map(), not Data: a JSON or DTO presentation must not stop a stale token being refreshed.
	body, ok := response.Map().(map[string]any)
	if !ok {
		return !http.IsUnsafeMethod(method)
	}
	if normalisedReason(body["error"]) == tokenFailureError {
		return true
	}
	return tokenFailureMessages[normalisedReason(body["message"])]
}

func normalisedReason(value any) string {
	if value == nil {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(fmt.Sprint(value)))
}

func progressFunc(config wccfg.SdkConfig, log func(string)) func(string) {
	return func(message string) {
		if config.PrintApiInfo && !config.PrintOnlyOnError {
			fmt.Println(message)
		}
		// An info record, so LogLevel "INFO" gets the timeline without touching stdout.
		if log != nil && logAllows(config, levelInfo) {
			log(message)
		}
	}
}

func (c *WebCommanderClient) reportFailedSend(err error, reportedBody any,
	files map[string]any, formEncoded bool) {
	var sdkError *wcerrors.SdkError
	if !errors.As(err, &sdkError) {
		return
	}
	info, ok := sdkError.RequestInfo.(*http.RequestInfo)
	if !ok || info == nil {
		return
	}
	failed := &http.ApiResponse{
		Headers:     map[string]string{},
		RequestInfo: info,
	}
	c.reportExchange(failed, reportedBody, files, formEncoded, levelError)
}

func (c *WebCommanderClient) reportExchange(response *http.ApiResponse, reportedBody any,
	files map[string]any, formEncoded bool, level string) {
	wantsPrint := c.config.PrintApiInfo || c.config.PrintRequestPayload || c.config.PrintResponseBody ||
		c.config.PrintHeaders || c.config.PrintCurl
	// Reported even without a request record: Redacted absorbs a nil one.
	if wantsPrint && !(c.config.PrintOnlyOnError && response.Succeeded()) {
		fmt.Println(renderExchange(response, reportedBody, files, formEncoded, c.config))
	}

	if c.log != nil {
		if logAllows(c.config, level) {
			// Redacted(), not the raw record: a host-built &RequestInfo{} never met the redacting constructor.
			info := response.RequestInfo.Redacted()
			// Named as SdkError.Error() names a call: by endpoint template, else by the resolved path's shape.
			c.log(fmt.Sprintf("%s%s %s -> %s in %ss (attempt %d/%d)", logRecordPrefixes[level],
				info.HTTPMethod, wcerrors.EndpointLabel(info.EndpointTemplate, info.URL),
				statusLabel(response.StatusCode),
				wcerrors.FormatSeconds(info.ElapsedSeconds), info.AttemptCount, info.MaxAttempts))
		}
		if c.config.LogBodies && logAllows(c.config, levelDebug) {
			// RenderJSONValue, not %v, which would write raw newlines and terminal escapes into the log.
			c.log(fmt.Sprintf("%spayload=%s response=%s", logRecordPrefixes[levelDebug],
				RenderJSONValue(redact.RedactValue(reportedBody), true),
				RenderJSONValue(redact.RedactValue(response.Map()), true)))
		}
	}
}

// endpointTemplateOf is the path as the caller wrote it, placeholders unresolved, for error rendering.
func endpointTemplateOf(path string) string {
	lowered := strings.ToLower(strings.TrimSpace(path))
	if strings.HasPrefix(lowered, "http://") || strings.HasPrefix(lowered, "https://") {
		return ""
	}
	return path
}
