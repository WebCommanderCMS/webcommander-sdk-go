package webcommander

import (
	"context"
	"iter"
	nethttp "net/http"
	"time"

	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/auth"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/client"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/config"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
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

func SDKVersion() string { return config.SDKVersion() }

// Version is SDKVersion under a shorter name — the same value.
func Version() string { return config.Version() }

// Ptr returns a pointer to v, for the optional/tri-state fields — e.g. wc.Ptr(true).
func Ptr[T any](v T) *T { return &v }

type (
	AuthConfig     = config.AuthConfig
	SdkConfig      = config.SdkConfig
	TokenCacheMode = config.TokenCacheMode
)

// SelectCredentials picks one credentials block out of a list by position or by client name.
func SelectCredentials(credentials []any, instance any) (any, error) {
	return config.SelectCredentials(credentials, instance)
}

func ReadClientName(block any) string { return config.ReadClientName(block) }

// SdkConfigFromCredentials builds one validated config; overrides are keyed by snake_case field names.
func SdkConfigFromCredentials(credentials []any, instance any, overrides map[string]any) (SdkConfig, error) {
	return config.SdkConfigFromCredentials(credentials, instance, overrides)
}

const (
	UseDefaultTokenCache      = config.UseDefaultTokenCache
	TokenCachePlatformDefault = config.TokenCachePlatformDefault
	TokenCacheMemoryOnly      = config.TokenCacheMemoryOnly
	// TokenCacheExplicitPath uses the file named by AuthConfig.TokenCachePath.
	TokenCacheExplicitPath = config.TokenCacheExplicitPath
	PrintFormatPretty      = config.PrintFormatPretty
	PrintFormatCompact     = config.PrintFormatCompact

	ResponseFormatMap     = config.ResponseFormatMap
	ResponseFormatJSON    = config.ResponseFormatJSON
	ResponseFormatDTO     = config.ResponseFormatDTO
	DefaultResponseFormat = config.DefaultResponseFormat
	// ResponseFormatEnvVar names the environment variable read when neither the call nor the config sets a response format.
	ResponseFormatEnvVar = config.ResponseFormatEnvVar
)

func IsResponseFormat(value string) bool { return config.IsResponseFormat(value) }

// ResponseFormatNames returns a fresh copy of the three legal response formats.
func ResponseFormatNames() []string { return config.ResponseFormatNames() }

type (
	// SdkError is the one error type the SDK returns.
	SdkError  = wcerrors.SdkError
	ErrorType = wcerrors.ErrorType
	// RequestContext is the request record an SdkError carries; *RequestInfo implements it.
	RequestContext = wcerrors.RequestContext
	// EndpointNamer is the optional half of RequestContext: it names the endpoint template.
	EndpointNamer            = wcerrors.EndpointNamer
	InsecureTransportWarning = wcerrors.InsecureTransportWarning
)

var (
	// ErrConfig matches configuration errors (errors.Is(err, ErrConfig)).
	ErrConfig = wcerrors.ErrConfig
	ErrUsage  = wcerrors.ErrUsage
	// ErrAttribute matches lookups of an unknown name; it is also a usage error (errors.Is both).
	ErrAttribute = wcerrors.ErrAttribute
	ErrAuth      = wcerrors.ErrAuth
)

// Build errors with these constructors: a hand-built &SdkError{} matches no sentinel.

func NewAttributeError(message string) *SdkError { return wcerrors.NewAttributeError(message) }

func NewConfigError(message string) *SdkError { return wcerrors.NewConfigError(message) }

func NewUsageError(message string) *SdkError { return wcerrors.NewUsageError(message) }

func NewAuthError(message string) *SdkError { return wcerrors.NewAuthError(message) }

func NewAuthResponseError(message string, statusCode int, errorType ErrorType, rawResponse any,
	responseHeaders map[string]string, requestInfo RequestContext, messages []string) *SdkError {
	return wcerrors.NewAuthResponseError(message, statusCode, errorType, rawResponse,
		responseHeaders, requestInfo, messages)
}

func NewInsecureTransportWarning(message string) *InsecureTransportWarning {
	return wcerrors.NewInsecureTransportWarning(message)
}

const (
	ErrorTypeNone             = wcerrors.ErrorTypeNone
	ErrorTypeRedirect         = wcerrors.ErrorTypeRedirect
	ErrorTypeValidation       = wcerrors.ErrorTypeValidation
	ErrorTypeAuthentication   = wcerrors.ErrorTypeAuthentication
	ErrorTypeAuthorization    = wcerrors.ErrorTypeAuthorization
	ErrorTypeNotFound         = wcerrors.ErrorTypeNotFound
	ErrorTypeMethodNotAllowed = wcerrors.ErrorTypeMethodNotAllowed
	ErrorTypeConflict         = wcerrors.ErrorTypeConflict
	ErrorTypeRateLimit        = wcerrors.ErrorTypeRateLimit
	ErrorTypeServer           = wcerrors.ErrorTypeServer
	ErrorTypeNetwork          = wcerrors.ErrorTypeNetwork
	ErrorTypeTimeout          = wcerrors.ErrorTypeTimeout
	ErrorTypeUnknown          = wcerrors.ErrorTypeUnknown
)

const (
	Redacted = redact.Redacted
	// TooDeep replaces a structure nested deeper than MaxRedactionDepth.
	TooDeep = redact.TooDeep
	// Circular replaces a container that contains itself.
	Circular = redact.Circular
	// MaxRedactionDepth is how deep RedactValue walks before writing TooDeep.
	MaxRedactionDepth = redact.MaxRedactionDepth
)

func ShouldRedact(key string) bool { return redact.ShouldRedact(key) }

// RedactValue walks a value and masks every secret-named entry, returning a fresh copy.
func RedactValue(value any) any { return redact.RedactValue(value) }

// RedactMapping masks the secret-named entries of a header map, returning a fresh copy.
func RedactMapping(headers map[string]string) map[string]string { return redact.RedactMapping(headers) }

// RedactText masks the secret-named query parameters and userinfo inside one string.
func RedactText(value string) string { return redact.RedactText(value) }

// SensitiveKeyMarkers returns a fresh copy of the substrings that make a key name secret.
func SensitiveKeyMarkers() []string { return redact.SensitiveKeyMarkers() }

type (
	ApiResponse      = http.ApiResponse
	RequestInfo      = http.RequestInfo
	RequestOptions   = http.RequestOptions
	TransportRequest = http.TransportRequest
	// Requester is the subset of the client that resource modules call.
	Requester     = http.Requester
	SortDirection = http.SortDirection
	// DTOShape builds the typed view of one response body.
	DTOShape = http.DTOShape
	// DTOPayload is a typed request body: anything that can hand over the map the SDK encodes.
	DTOPayload = http.DTOPayload
)

func BuildURL(baseURL, path string) (string, error) { return http.BuildURL(baseURL, path) }

// SubstitutePathParams fills a path template's {placeholders} from a map.
func SubstitutePathParams(pathTemplate string, pathParams map[string]any) (string, error) {
	return http.SubstitutePathParams(pathTemplate, pathParams)
}

func BuildQueryParams(params map[string]any) (map[string]any, error) {
	return http.BuildQueryParams(params)
}

// NewRequestInfo is the redacting constructor for RequestInfo.
func NewRequestInfo(clientName, apiVersion, httpMethod, url string, queryParams map[string]any,
	body any, headers map[string]string, attemptCount, maxAttempts int, elapsedSeconds float64) *RequestInfo {
	return http.NewRequestInfo(clientName, apiVersion, httpMethod, url, queryParams, body, headers,
		attemptCount, maxAttempts, elapsedSeconds)
}

// DefaultRetryableHTTPCodes is the default retryable set for a safe method; a fresh copy each call.
func DefaultRetryableHTTPCodes() []int { return config.DefaultRetryableCodes() }

// DefaultUnsafeRetryableHTTPCodes is empty, so no write is replayed unless a caller opts in.
func DefaultUnsafeRetryableHTTPCodes() []int { return config.DefaultUnsafeRetryableCodes() }

const (
	SortAsc  = http.SortAsc
	SortDesc = http.SortDesc
)

type (
	HTTPTransport   = http.HTTPTransport
	TransportOption = http.TransportOption
	RetryPolicy     = http.RetryPolicy
)

// MaxMultipartBytes caps the multipart body assembled in memory: 64 MiB.
const MaxMultipartBytes = http.MaxMultipartBytes

// MaxRequestBodyBytes caps the final wire body of any request, whichever kind it is: 64 MiB.
const MaxRequestBodyBytes = http.MaxRequestBodyBytes

const MaxCABundleBytes = http.MaxCABundleBytes

const MaxTokenCacheBytes = auth.MaxTokenCacheBytes

func NewHTTPTransport(sdkConfig SdkConfig, options ...TransportOption) (*HTTPTransport, error) {
	return http.NewHTTPTransport(sdkConfig, options...)
}

func WithHTTPClient(httpClient *nethttp.Client) TransportOption {
	return http.WithHTTPClient(httpClient)
}

// WithProgress wires the transport's per-attempt progress hook.
func WithProgress(onProgress func(string)) TransportOption { return http.WithProgress(onProgress) }

// WithTransportRetryPolicy is http.WithRetryPolicy, named apart from the client's WithRetryPolicy.
func WithTransportRetryPolicy(retry *RetryPolicy) TransportOption {
	return http.WithRetryPolicy(retry)
}

// NewRetryPolicy constructs a retry policy; a nil code list takes the matching default set.
func NewRetryPolicy(maxAttempts int, backoffSeconds, backoffMaxSeconds float64,
	retryable, unsafeRetryable []int) (*RetryPolicy, error) {
	return http.NewRetryPolicy(maxAttempts, backoffSeconds, backoffMaxSeconds, retryable, unsafeRetryable)
}

// SnapshotJSONValue deep-copies one JSON-shaped value at any depth; RequestOptions.Frozen uses it.
func SnapshotJSONValue(value any) any { return http.SnapshotJSONValue(value) }

func ReadCappedFile(path string, limitBytes int) ([]byte, error) {
	return http.ReadCappedFile(path, limitBytes)
}

type (
	TokenRecord = auth.TokenRecord
	TokenStore  = auth.TokenStore
	// ContextTokenStore is a TokenStore whose lock wait can be abandoned when the context is done.
	ContextTokenStore = auth.ContextTokenStore
	// TokenRefreshScope is the store operations under a held refresh lock.
	TokenRefreshScope = auth.TokenRefreshScope
	// NestedTokenRefreshScope is a TokenRefreshScope re-entrant over the lock it already holds.
	NestedTokenRefreshScope = auth.NestedTokenRefreshScope
	MemoryTokenStore        = auth.MemoryTokenStore
	// FileTokenStore is token storage in a JSON file, shared across processes under an OS lock.
	FileTokenStore = auth.FileTokenStore
	TokenProvider  = auth.TokenProvider
)

func NewMemoryTokenStore() *MemoryTokenStore { return auth.NewMemoryTokenStore() }

// NewFileTokenStore builds token storage in a JSON file, under DefaultLockWait.
func NewFileTokenStore(filePath, namespace string) (*FileTokenStore, error) {
	return auth.NewFileTokenStore(filePath, namespace)
}

func NewFileTokenStoreWithLockWait(filePath, namespace string, lockWait time.Duration) (*FileTokenStore, error) {
	return auth.NewFileTokenStoreWithLockWait(filePath, namespace, lockWait)
}

func NewTokenProvider(sdkConfig SdkConfig, transport *HTTPTransport, store TokenStore) *TokenProvider {
	return auth.NewTokenProvider(sdkConfig, transport, store)
}

func BuildCredentialsNamespace(baseURL, clientID, clientSecret, authString, grantType string) string {
	return auth.BuildCredentialsNamespace(baseURL, clientID, clientSecret, authString, grantType)
}

// BuildTokenPath is the token endpoint path for an API version ("/api/v4/oauth2/token").
func BuildTokenPath(apiVersion string) string { return auth.BuildTokenPath(apiVersion) }

// BuildDefaultTokenCachePath is the per-user cache path used when no TokenCachePath is named.
func BuildDefaultTokenCachePath(filename string) (string, error) {
	return auth.BuildDefaultTokenCachePath(filename)
}

const DefaultTokenCacheFilename = auth.DefaultTokenCacheFilename

type (
	WebCommanderClient = client.WebCommanderClient
	ClientOption       = client.ClientOption
	ClientRegistry     = client.ClientRegistry
	// BoundResource is a resource family with a client applied; type-assert Value to the resource.
	BoundResource   = client.BoundResource
	ResourceFactory = client.ResourceFactory
	PaginateOptions = client.PaginateOptions
)

func NewClient(sdkConfig SdkConfig, options ...ClientOption) (*WebCommanderClient, error) {
	return client.NewClient(sdkConfig, options...)
}

func WithTransport(transport *HTTPTransport) ClientOption { return client.WithTransport(transport) }

func WithTokenStore(store TokenStore) ClientOption { return client.WithTokenStore(store) }

// WithTokenProvider injects a built TokenProvider, so clients can share one provider and cache.
func WithTokenProvider(provider *TokenProvider) ClientOption {
	return client.WithTokenProvider(provider)
}

func WithRetryPolicy(retry *RetryPolicy) ClientOption { return client.WithRetryPolicy(retry) }

func WithLogger(log func(string)) ClientOption { return client.WithLogger(log) }

func IterPages(ctx context.Context, requester Requester, path string,
	options PaginateOptions) iter.Seq2[*ApiResponse, error] {
	return client.IterPages(ctx, requester, path, options)
}

// IterRecords walks every record of a listing across every page, each element untouched.
func IterRecords(ctx context.Context, requester Requester, path, recordKey string,
	options PaginateOptions) iter.Seq2[any, error] {
	return client.IterRecords(ctx, requester, path, recordKey, options)
}

func NewClientRegistry(credentials ...any) (*ClientRegistry, error) {
	return client.NewClientRegistry(credentials...)
}

func RegisterResourceModule(resourceName string, factory ResourceFactory) error {
	return client.RegisterResourceModule(resourceName, factory)
}

func RegisteredResourceNames() []string { return client.RegisteredResourceNames() }

// RegisterInstances registers credential blocks by client name for later BuildClient calls.
func RegisterInstances(credentials ...any) error { return client.RegisterInstances(credentials...) }

func BuildClient(instance any, overrides ...func(*SdkConfig)) (*WebCommanderClient, error) {
	return client.BuildClient(instance, overrides...)
}

func ListInstanceNames() ([]string, error) { return client.ListInstanceNames() }

func GetRegisteredInstances() ([]SdkConfig, error) { return client.GetRegisteredInstances() }

// RenderExchange renders one exchange as the console would print it.
func RenderExchange(response *ApiResponse, body any, files map[string]any,
	sdkConfig SdkConfig, formEncoded bool) string {
	return client.RenderExchange(response, body, files, sdkConfig, formEncoded)
}

func BuildCurlCommand(response *ApiResponse, payload any, files map[string]any,
	sdkConfig SdkConfig, formEncoded bool) string {
	return client.BuildCurlCommand(response, payload, files, sdkConfig, formEncoded)
}

// AdminCustomers is the admin customer endpoints bound to a client — the primary operation surface.
type AdminCustomers = customers.AdminCustomers

// DTOFields is the wire names a DTO type declares, in declaration order.
func DTOFields(value any) []string { return dto.Fields(value) }

type (
	// Customer is one customer, in either the list form or the detail form.
	Customer = customers.Customer
	Address  = customers.Address
	// Reference is an embedded {id, name, code} lookup value.
	Reference = customers.Reference
	// CustomProperty is one entry of a customer's custom_properties; a write replaces the whole list.
	CustomProperty = customers.CustomProperty
	// CustomerInput is the writable form of a Customer; a field left at its zero is omitted.
	CustomerInput         = customers.CustomerInput
	StoreCreditAdjustment = customers.StoreCreditAdjustment
	// StoreCreditAdjustmentInput is the writable form of one; its body carries no envelope.
	StoreCreditAdjustmentInput = customers.StoreCreditAdjustmentInput
	Pagination                 = customers.Pagination
	// Page is one page of a listing: records on Items, paging block beside them.
	Page[T any] = customers.Page[T]
)

// PageOf is a list endpoint's DTOShape, told its record type by instantiation — PageOf[Customer].
func PageOf[T any](decoded any) any { return customers.PageOf[T](decoded) }

func NewAdminCustomers(requester Requester) *AdminCustomers { return customers.New(requester) }

// CustomerShape is the DTOShape for a body shaped {"customer": {…}}, for the raw-verb call.
func CustomerShape(decoded any) any { return customers.CustomerShape(decoded) }

func StoreCreditAdjustmentShape(decoded any) any {
	return customers.StoreCreditAdjustmentShape(decoded)
}

// The path templates for the raw-verb call.
const (
	AdminCustomersPath                  = customers.AdminCustomersPath
	AdminCustomerPath                   = customers.AdminCustomerPath
	AdminCustomersCountPath             = customers.AdminCustomersCountPath
	AdminCustomersCheckEmailPath        = customers.AdminCustomersCheckEmailPath
	AdminCustomersExportPath            = customers.AdminCustomersExportPath
	AdminCustomersImportPath            = customers.AdminCustomersImportPath
	AdminCustomersBulkResetPasswordPath = customers.AdminCustomersBulkResetPasswordPath
	AdminCustomerAddressesPath          = customers.AdminCustomerAddressesPath
	AdminCustomerResetPasswordPath      = customers.AdminCustomerResetPasswordPath
	AdminCustomerChangePasswordPath     = customers.AdminCustomerChangePasswordPath

	AdminCustomerStoreCreditAdjustmentsPath = customers.AdminCustomerStoreCreditAdjustmentsPath

	// AdminCustomersTransactionNoHeader is the header CreateCustomer sends transactionNo under.
	AdminCustomersTransactionNoHeader = customers.AdminCustomersTransactionNoHeader
)

// SeoConfigurationSettings is the SEO configuration settings endpoints bound to a client.
type SeoConfigurationSettings = seoconfiguration.SeoConfigurationSettings

type (
	// SeoConfiguration is the area, in either form the API sends it: the facade's, and one section's.
	SeoConfiguration = seoconfiguration.SeoConfiguration
	// GlobalStoreSettings through URLSettings are the three sections, each its own envelope.
	GlobalStoreSettings = seoconfiguration.GlobalStoreSettings
	CustomCode          = seoconfiguration.CustomCode
	URLSettings         = seoconfiguration.URLSettings
	// SeoTemplates is one content type's block of GlobalStoreSettings.
	SeoTemplates = seoconfiguration.SeoTemplates
	// CustomCodeBlock is the header or footer half of CustomCode.
	CustomCodeBlock = seoconfiguration.CustomCodeBlock
	// Redirect is one redirect row, and the typed body a create accepts wrapped.
	Redirect = seoconfiguration.Redirect
	// RedirectsPreview and DocumentRootPreview are the collections inlined for reading, no pagination.
	RedirectsPreview    = seoconfiguration.RedirectsPreview
	DocumentRootPreview = seoconfiguration.DocumentRootPreview
	// DocumentRootFile is one file served from the document root; ConfigKey is what the delete takes.
	DocumentRootFile = seoconfiguration.DocumentRootFile
	// GroupMeta and SectionMeta are what the facade declares the area holds; only the facade sends them.
	GroupMeta   = seoconfiguration.GroupMeta
	SectionMeta = seoconfiguration.SectionMeta
)

// NewSeoConfigurationSettings binds the resource to a requester of your own.
func NewSeoConfigurationSettings(requester Requester) *SeoConfigurationSettings {
	return seoconfiguration.New(requester)
}

// SeoConfigurationShape is the DTOShape for the facade's body and for one section's.
func SeoConfigurationShape(decoded any) any { return seoconfiguration.SeoConfigurationShape(decoded) }

// RedirectShape is the DTOShape for a body shaped {"redirect": {…}}.
func RedirectShape(decoded any) any { return seoconfiguration.RedirectShape(decoded) }

// DocumentRootFileShape is the DTOShape for a body shaped {"document_root_file": {…}}.
func DocumentRootFileShape(decoded any) any {
	return seoconfiguration.DocumentRootFileShape(decoded)
}

// SeoConfigurationSections returns a fresh copy of the three values {section} takes.
func SeoConfigurationSections() []string { return seoconfiguration.SeoConfigurationSections() }

const (
	SeoConfigurationPath        = seoconfiguration.SeoConfigurationPath
	SeoConfigurationSectionPath = seoconfiguration.SeoConfigurationSectionPath

	SeoConfigurationDocumentRootPath     = seoconfiguration.SeoConfigurationDocumentRootPath
	SeoConfigurationDocumentRootFilePath = seoconfiguration.SeoConfigurationDocumentRootFilePath

	SeoConfigurationRedirectsPath = seoconfiguration.SeoConfigurationRedirectsPath
	SeoConfigurationRedirectPath  = seoconfiguration.SeoConfigurationRedirectPath
)

// SearchResultPageSettings is the search result page settings endpoints bound to a client.
type SearchResultPageSettings = pagesettings.SearchResultPageSettings

type (
	// SearchResultPage is the area, in either form the API sends it: the facade's, and one section's.
	SearchResultPage = pagesettings.SearchResultPage
	// ProductInformation through Reinitialise are the three sections, each its own envelope.
	ProductInformation = pagesettings.ProductInformation
	Seo                = pagesettings.Seo
	Reinitialise       = pagesettings.Reinitialise
	// ProductInteractionButtons through DisplayPrice are the six blocks inside ProductInformation.
	ProductInteractionButtons    = pagesettings.ProductInteractionButtons
	NumberOfSearchesPerResult    = pagesettings.NumberOfSearchesPerResult
	SetProductDisplayPreferences = pagesettings.SetProductDisplayPreferences
	SetProductSorting            = pagesettings.SetProductSorting
	ConfigureSortingOptions      = pagesettings.ConfigureSortingOptions
	DisplayPrice                 = pagesettings.DisplayPrice
	// ProductPromotion is one promotion row, and the typed body every write to that collection takes.
	ProductPromotion = pagesettings.ProductPromotion
	// ProductPromotionPreview is the collection's marker block on the facade — not its rows.
	ProductPromotionPreview = pagesettings.ProductPromotionPreview
	// SearchResultPageGroupMeta and SearchResultPageSectionMeta are the area's GroupMeta and SectionMeta.
	SearchResultPageGroupMeta   = pagesettings.GroupMeta
	SearchResultPageSectionMeta = pagesettings.SectionMeta
)

// NewSearchResultPageSettings binds the resource to a requester of your own.
func NewSearchResultPageSettings(requester Requester) *SearchResultPageSettings {
	return pagesettings.New(requester)
}

// SearchResultPageShape is the DTOShape for the facade's body and for one section's.
func SearchResultPageShape(decoded any) any { return pagesettings.SearchResultPageShape(decoded) }

// ProductPromotionShape is the DTOShape for a body shaped {"product_promotion": {…}}.
func ProductPromotionShape(decoded any) any { return pagesettings.ProductPromotionShape(decoded) }

// SearchResultPageSections returns a fresh copy of the three values {section} takes.
func SearchResultPageSections() []string { return pagesettings.SearchResultPageSections() }

const (
	SearchResultPagePath        = pagesettings.SearchResultPagePath
	SearchResultPageSectionPath = pagesettings.SearchResultPageSectionPath

	SearchResultPageProductPromotionsPath      = pagesettings.SearchResultPageProductPromotionsPath
	SearchResultPageProductPromotionPath       = pagesettings.SearchResultPageProductPromotionPath
	SearchResultPageProductPromotionActivePath = pagesettings.SearchResultPageProductPromotionActivePath
)

// CategoryPageSettings is the category page settings endpoints bound to a client.
type CategoryPageSettings = pagesettings.CategoryPageSettings

type (
	// CategoryPage is the area, in either form the API sends it: the facade's, and one section's.
	CategoryPage = pagesettings.CategoryPage
	// CategoryInformation and CategoryPageProductInformation are the two sections, each its own envelope.
	CategoryInformation            = pagesettings.CategoryInformation
	CategoryPageProductInformation = pagesettings.CategoryPageProductInformation
	// CategoryHeading and Subcategory are the two blocks inside CategoryInformation.
	CategoryHeading = pagesettings.CategoryHeading
	Subcategory     = pagesettings.Subcategory
	// InteractionButtons through SortOptions are the blocks inside CategoryPageProductInformation.
	InteractionButtons   = pagesettings.InteractionButtons
	Price                = pagesettings.Price
	VariationCombination = pagesettings.VariationCombination
	VariationType        = pagesettings.VariationType
	Rating               = pagesettings.Rating
	ViewSwitcher         = pagesettings.ViewSwitcher
	DefaultDisplay       = pagesettings.DefaultDisplay
	ColumnsPerRow        = pagesettings.ColumnsPerRow
	ProductLoad          = pagesettings.ProductLoad
	DefaultSort          = pagesettings.DefaultSort
	SortOptions          = pagesettings.SortOptions
	// CategoryPageGroupMeta and CategoryPageSectionMeta are the area's GroupMeta and SectionMeta.
	CategoryPageGroupMeta   = pagesettings.CategoryPageGroupMeta
	CategoryPageSectionMeta = pagesettings.CategoryPageSectionMeta
)

// NewCategoryPageSettings binds the resource to a requester of your own.
func NewCategoryPageSettings(requester Requester) *CategoryPageSettings {
	return pagesettings.NewCategoryPageSettings(requester)
}

// CategoryPageShape is the DTOShape for the facade's body and for one section's.
func CategoryPageShape(decoded any) any { return pagesettings.CategoryPageShape(decoded) }

// CategoryPageSections returns a fresh copy of the two values {section} takes.
func CategoryPageSections() []string { return pagesettings.CategoryPageSections() }

const (
	CategoryPagePath        = pagesettings.CategoryPagePath
	CategoryPageSectionPath = pagesettings.CategoryPageSectionPath
)

// ShoppingCartSettings is the shopping cart settings endpoints bound to a client.
type ShoppingCartSettings = pagesettings.ShoppingCartSettings

type (
	// ShoppingCart is the area, in either form the API sends it: the facade's, and one section's.
	ShoppingCart = pagesettings.ShoppingCart
	// CartPage and FlyingCart are the two sections, each its own envelope.
	CartPage   = pagesettings.CartPage
	FlyingCart = pagesettings.FlyingCart
	// ShoppingCartGroupMeta and ShoppingCartSectionMeta are the area's GroupMeta and SectionMeta.
	ShoppingCartGroupMeta   = pagesettings.ShoppingCartGroupMeta
	ShoppingCartSectionMeta = pagesettings.ShoppingCartSectionMeta
)

// NewShoppingCartSettings binds the resource to a requester of your own.
func NewShoppingCartSettings(requester Requester) *ShoppingCartSettings {
	return pagesettings.NewShoppingCartSettings(requester)
}

// ShoppingCartShape is the DTOShape for the facade's body and for one section's.
func ShoppingCartShape(decoded any) any { return pagesettings.ShoppingCartShape(decoded) }

// ShoppingCartSections returns a fresh copy of the two values {section} takes.
func ShoppingCartSections() []string { return pagesettings.ShoppingCartSections() }

const (
	ShoppingCartPath        = pagesettings.ShoppingCartPath
	ShoppingCartSectionPath = pagesettings.ShoppingCartSectionPath
)

// ProductPageSettings is the product page settings endpoints bound to a client.
type ProductPageSettings = pagesettings.ProductPageSettings

type (
	// ProductPage is the area, in either form the API sends it: the facade's, and one section's.
	ProductPage = pagesettings.ProductPage
	// Configuration through QuickView are the seven sections, each its own envelope.
	Configuration    = pagesettings.Configuration
	Variations       = pagesettings.Variations
	ImageAndVideo    = pagesettings.ImageAndVideo
	UpSell           = pagesettings.UpSell
	RelatedProducts  = pagesettings.RelatedProducts
	ReviewAndRatings = pagesettings.ReviewAndRatings
	QuickView        = pagesettings.QuickView
	// ActionButtons through ReviewCapabilities are the blocks the four nesting sections hold.
	ActionButtons            = pagesettings.ActionButtons
	ProductDetails           = pagesettings.ProductDetails
	SocialSharing            = pagesettings.SocialSharing
	SetProductProperties     = pagesettings.SetProductProperties
	LabelsAndDisplay         = pagesettings.LabelsAndDisplay
	VariationLayout          = pagesettings.VariationLayout
	ProductPageVariationType = pagesettings.ProductPageVariationType
	ImageDisplay             = pagesettings.ImageDisplay
	ProductVideo             = pagesettings.ProductVideo
	ReviewCapabilities       = pagesettings.ReviewCapabilities
	// ProductListing is the listing the up-sell and related-products sections each copy.
	ProductListing = pagesettings.ProductListing
	// ProductPageGroupMeta and ProductPageSectionMeta are the area's GroupMeta and SectionMeta.
	ProductPageGroupMeta   = pagesettings.ProductPageGroupMeta
	ProductPageSectionMeta = pagesettings.ProductPageSectionMeta
)

// NewProductPageSettings binds the resource to a requester of your own.
func NewProductPageSettings(requester Requester) *ProductPageSettings {
	return pagesettings.NewProductPageSettings(requester)
}

// ProductPageShape is the DTOShape for the facade's body and for one section's.
func ProductPageShape(decoded any) any { return pagesettings.ProductPageShape(decoded) }

// ProductPageSections returns a fresh copy of the seven values {section} takes.
func ProductPageSections() []string { return pagesettings.ProductPageSections() }

const (
	ProductPagePath        = pagesettings.ProductPagePath
	ProductPageSectionPath = pagesettings.ProductPageSectionPath
)

// CheckoutPageSettings is the checkout page settings endpoints bound to a client.
type CheckoutPageSettings = pagesettings.CheckoutPageSettings

type (
	// CheckoutPage is the area, in either form the API sends it: the facade's, and one section's.
	CheckoutPage = pagesettings.CheckoutPage
	// BillingAddress through Others are the six sections; the two address sections are read-only.
	BillingAddress   = pagesettings.BillingAddress
	ShippingAddress  = pagesettings.ShippingAddress
	CustomerCheckout = pagesettings.CustomerCheckout
	DeliveryOptions  = pagesettings.DeliveryOptions
	OrderSummary     = pagesettings.OrderSummary
	Others           = pagesettings.Others
	// AddressForm is the declaration both address sections copy, and AddressField its list entry.
	AddressForm  = pagesettings.AddressForm
	AddressField = pagesettings.AddressField
	// SignUp through TermsReference are the blocks the two nesting sections hold.
	SignUp            = pagesettings.SignUp
	GuestCheckout     = pagesettings.GuestCheckout
	OrderComment      = pagesettings.OrderComment
	TermsAndCondition = pagesettings.TermsAndCondition
	TermsReference    = pagesettings.TermsReference
	// CheckoutPageGroupMeta and CheckoutPageSectionMeta are the area's GroupMeta and SectionMeta.
	CheckoutPageGroupMeta   = pagesettings.CheckoutPageGroupMeta
	CheckoutPageSectionMeta = pagesettings.CheckoutPageSectionMeta
)

// NewCheckoutPageSettings binds the resource to a requester of your own.
func NewCheckoutPageSettings(requester Requester) *CheckoutPageSettings {
	return pagesettings.NewCheckoutPageSettings(requester)
}

// CheckoutPageShape is the DTOShape for the facade's body and for one section's.
func CheckoutPageShape(decoded any) any { return pagesettings.CheckoutPageShape(decoded) }

// CheckoutPageSections returns a fresh copy of the six values {section} takes.
func CheckoutPageSections() []string { return pagesettings.CheckoutPageSections() }

// CheckoutPageWritableSections returns a fresh copy of the four of those six a PATCH may name.
func CheckoutPageWritableSections() []string { return pagesettings.CheckoutPageWritableSections() }

const (
	CheckoutPagePath        = pagesettings.CheckoutPagePath
	CheckoutPageSectionPath = pagesettings.CheckoutPageSectionPath
)

// PriceDisplaySettingsResource is the price display settings endpoints bound to a client.
type PriceDisplaySettingsResource = pricedisplaysettings.PriceDisplaySettingsResource

type (
	// PriceDisplaySettings is the area, in either form the API sends it: the facade's, and one section's.
	PriceDisplaySettings = pricedisplaysettings.PriceDisplaySettings
	// OrderItem and PriceDisplaySettingsOrderSummary are the two sections, each its own envelope.
	OrderItem = pricedisplaysettings.OrderItem
	// PriceDisplaySettingsOrderSummary is the order_summary section.
	PriceDisplaySettingsOrderSummary = pricedisplaysettings.OrderSummary
	// DisplayScopes through PricingInfo are the blocks the two sections hold.
	DisplayScopes    = pricedisplaysettings.DisplayScopes
	ScopedToggle     = pricedisplaysettings.ScopedToggle
	PlainToggle      = pricedisplaysettings.PlainToggle
	Subtotal         = pricedisplaysettings.Subtotal
	Discount         = pricedisplaysettings.Discount
	ShippingCost     = pricedisplaysettings.ShippingCost
	ShippingDiscount = pricedisplaysettings.ShippingDiscount
	Handling         = pricedisplaysettings.Handling
	PricingInfo      = pricedisplaysettings.PricingInfo
	// PriceDisplaySettingsGroupMeta and PriceDisplaySettingsSectionMeta are the area's GroupMeta pair.
	PriceDisplaySettingsGroupMeta   = pricedisplaysettings.GroupMeta
	PriceDisplaySettingsSectionMeta = pricedisplaysettings.SectionMeta
)

// NewPriceDisplaySettings binds the resource to a requester of your own.
func NewPriceDisplaySettings(requester Requester) *PriceDisplaySettingsResource {
	return pricedisplaysettings.New(requester)
}

// PriceDisplaySettingsShape is the DTOShape for the facade's body and for one section's.
func PriceDisplaySettingsShape(decoded any) any {
	return pricedisplaysettings.PriceDisplaySettingsShape(decoded)
}

// PriceDisplaySettingsSections returns a fresh copy of the two values {section} takes.
func PriceDisplaySettingsSections() []string {
	return pricedisplaysettings.PriceDisplaySettingsSections()
}

const (
	PriceDisplaySettingsPath        = pricedisplaysettings.PriceDisplaySettingsPath
	PriceDisplaySettingsSectionPath = pricedisplaysettings.PriceDisplaySettingsSectionPath
)

// AnimationEffectsSettings is the animation effects settings endpoints bound to a client.
type AnimationEffectsSettings = animationeffects.AnimationEffectsSettings

type (
	// AnimationEffects is the area, in either form the API sends it: the facade's, and one section's.
	AnimationEffects = animationeffects.AnimationEffects
	// PreloaderDesign through PageSkeleton are the five sections, each its own envelope.
	PreloaderDesign  = animationeffects.PreloaderDesign
	WidgetTransition = animationeffects.WidgetTransition
	PageTransition   = animationeffects.PageTransition
	ImageLazyEffect  = animationeffects.ImageLazyEffect
	PageSkeleton     = animationeffects.PageSkeleton
	// CustomLoader is the read-only loader block, and PreloaderPage the entry behind its one list.
	CustomLoader  = animationeffects.CustomLoader
	PreloaderPage = animationeffects.PreloaderPage
	// AnimationEffectsGroupMeta and AnimationEffectsSectionMeta are the area's GroupMeta and SectionMeta.
	AnimationEffectsGroupMeta   = animationeffects.GroupMeta
	AnimationEffectsSectionMeta = animationeffects.SectionMeta
)

// NewAnimationEffectsSettings binds the resource to a requester of your own.
func NewAnimationEffectsSettings(requester Requester) *AnimationEffectsSettings {
	return animationeffects.New(requester)
}

// AnimationEffectsShape is the DTOShape for the facade's body and for one section's.
func AnimationEffectsShape(decoded any) any { return animationeffects.AnimationEffectsShape(decoded) }

// AnimationEffectsSections returns a fresh copy of the five values {section} takes.
func AnimationEffectsSections() []string { return animationeffects.AnimationEffectsSections() }

const (
	AnimationEffectsPath        = animationeffects.AnimationEffectsPath
	AnimationEffectsSectionPath = animationeffects.AnimationEffectsSectionPath
)

// ProductListingSettingsResource is the product listing settings endpoints bound to a client.
type ProductListingSettingsResource = productlisting.ProductListingSettingsResource

type (
	// ProductListingSettings is the area, in the one form the API sends it: the facade's.
	ProductListingSettings = productlisting.ProductListing
	// ProductListingProductInformation through SortableItems are the four sections; none has an envelope.
	ProductListingProductInformation = productlisting.ProductInformation
	Load                             = productlisting.Load
	Sorting                          = productlisting.Sorting
	SortableItems                    = productlisting.SortableItems
	// ProductListingInteractionButtons through SortableOptions are the blocks the four sections hold.
	ProductListingInteractionButtons   = productlisting.InteractionButtons
	ListingPrice                       = productlisting.ListingPrice
	Reviews                            = productlisting.Reviews
	ProductListingVariationCombination = productlisting.VariationCombination
	ProductListingPagination           = productlisting.Pagination
	LoadMore                           = productlisting.LoadMore
	SortableOptions                    = productlisting.SortableOptions
	// ProductListingGroupMeta is the area's GroupMeta.
	ProductListingGroupMeta = productlisting.GroupMeta
	// ProductListingSectionMeta is the area's SectionMeta, the one group_meta entry.
	ProductListingSectionMeta = productlisting.SectionMeta
)

// NewProductListingSettings binds the resource to a requester of your own.
func NewProductListingSettings(requester Requester) *ProductListingSettingsResource {
	return productlisting.New(requester)
}

// ProductListingShape is the DTOShape for the facade's body, the one shape this area sends.
func ProductListingShape(decoded any) any { return productlisting.ProductListingShape(decoded) }

// ProductListingSections returns a fresh copy of the four section names a PATCH body may carry.
func ProductListingSections() []string { return productlisting.ProductListingSections() }

const ProductListingPath = productlisting.ProductListingPath

// CustomerPortalSettings is the customer portal settings endpoints bound to a client.
type CustomerPortalSettings = customerportal.CustomerPortalSettings

type (
	// CustomerPortal is the area, in either form the API sends it: the facade's, and one section's.
	CustomerPortal = customerportal.CustomerPortal
	// Portal through Sso are the four writable sections, each its own envelope.
	Portal       = customerportal.Portal
	Registration = customerportal.Registration
	Login        = customerportal.Login
	Sso          = customerportal.Sso
	// ProfilePage and RegistrationFields are read-only sections: read and typed, never written.
	ProfilePage        = customerportal.ProfilePage
	RegistrationFields = customerportal.RegistrationFields
	// RegistrationTerms through AssignToCustomer are the blocks those sections hold.
	RegistrationTerms  = customerportal.RegistrationTerms
	FacebookLogin      = customerportal.FacebookLogin
	GoogleLogin        = customerportal.GoogleLogin
	B2bConfiguration   = customerportal.B2bConfiguration
	AlternateLogin     = customerportal.AlternateLogin
	LoginCaptcha       = customerportal.LoginCaptcha
	AzureAd            = customerportal.AzureAd
	AzureAdB2c         = customerportal.AzureAdB2c
	AzureAdDataMapping = customerportal.AzureAdDataMapping
	ManageMyAccount    = customerportal.ManageMyAccount
	MyOrders           = customerportal.MyOrders
	MyOrderFlags       = customerportal.MyOrderFlags
	AssignToCustomer   = customerportal.AssignToCustomer
	// CustomerGroup through RegistrationField are the row types this area's lists carry.
	CustomerGroup        = customerportal.CustomerGroup
	ProfileField         = customerportal.ProfileField
	ProfileEntry         = customerportal.ProfileEntry
	ProfileOverviewEntry = customerportal.ProfileOverviewEntry
	RegistrationField    = customerportal.RegistrationField
	// CustomerPortalGroupMeta and CustomerPortalSectionMeta are the area's GroupMeta and SectionMeta.
	CustomerPortalGroupMeta   = customerportal.GroupMeta
	CustomerPortalSectionMeta = customerportal.SectionMeta
)

// NewCustomerPortalSettings binds the resource to a requester of your own.
func NewCustomerPortalSettings(requester Requester) *CustomerPortalSettings {
	return customerportal.New(requester)
}

// CustomerPortalShape is the DTOShape for the facade's body and for one section's.
func CustomerPortalShape(decoded any) any { return customerportal.CustomerPortalShape(decoded) }

// CustomerPortalSections returns a fresh copy of the six values {section} takes.
func CustomerPortalSections() []string { return customerportal.CustomerPortalSections() }

// CustomerPortalWritableSections returns a fresh copy of the four sections a PATCH may name.
func CustomerPortalWritableSections() []string {
	return customerportal.CustomerPortalWritableSections()
}

const (
	CustomerPortalPath        = customerportal.CustomerPortalPath
	CustomerPortalSectionPath = customerportal.CustomerPortalSectionPath
)

// SiteSettingsResource is the site settings endpoints bound to a client.
type SiteSettingsResource = sitesettings.SiteSettingsResource

type (
	// SiteSettings is the area, in either form the API sends it: the facade's, and one section's.
	SiteSettings = sitesettings.SiteSettings
	// Administration through StoreLogo are the eight sections, each its own envelope.
	Administration = sitesettings.Administration
	BrandKit       = sitesettings.BrandKit
	DateTime       = sitesettings.DateTime
	DomainSettings = sitesettings.DomainSettings
	EmailSettings  = sitesettings.EmailSettings
	Favicon        = sitesettings.Favicon
	StoreDetails   = sitesettings.StoreDetails
	StoreLogo      = sitesettings.StoreLogo
	// AutoSuggest through StoreLogoVariants are the blocks those sections hold.
	AutoSuggest          = sitesettings.AutoSuggest
	AutoSuggestAddresses = sitesettings.AutoSuggestAddresses
	Captcha              = sitesettings.Captcha
	BrandColors          = sitesettings.BrandColors
	AdminDateTime        = sitesettings.AdminDateTime
	SeparableDateTime    = sitesettings.SeparableDateTime
	Maintenance          = sitesettings.Maintenance
	SecurityCode         = sitesettings.SecurityCode
	Ssl                  = sitesettings.Ssl
	SslCertificateFile   = sitesettings.SslCertificateFile
	SslKeyFile           = sitesettings.SslKeyFile
	DefaultEmail         = sitesettings.DefaultEmail
	DefaultDomain        = sitesettings.DefaultDomain
	EmailProvider        = sitesettings.EmailProvider
	OauthProvider        = sitesettings.OauthProvider
	ReplyTo              = sitesettings.ReplyTo
	SendGrid             = sitesettings.SendGrid
	Smtp                 = sitesettings.Smtp
	SmtpAuthentication   = sitesettings.SmtpAuthentication
	StoreAddress         = sitesettings.StoreAddress
	StoreContact         = sitesettings.StoreContact
	StoreLogoVariants    = sitesettings.StoreLogoVariants
	// Region and MediaImage are shared blocks: Region by two addresses, MediaImage by both media sections.
	Region     = sitesettings.Region
	MediaImage = sitesettings.MediaImage
	// SiteSettingsGroupMeta and SiteSettingsSectionMeta are the area's GroupMeta and SectionMeta.
	SiteSettingsGroupMeta   = sitesettings.GroupMeta
	SiteSettingsSectionMeta = sitesettings.SectionMeta
)

// NewSiteSettings binds the resource to a requester of your own.
func NewSiteSettings(requester Requester) *SiteSettingsResource {
	return sitesettings.New(requester)
}

// SiteSettingsShape is the DTOShape for the facade's body and for one section's.
func SiteSettingsShape(decoded any) any { return sitesettings.SiteSettingsShape(decoded) }

// SiteSettingsSections returns a fresh copy of the eight values {section} takes.
func SiteSettingsSections() []string { return sitesettings.SiteSettingsSections() }

// SiteSettingsMediaSections returns a fresh copy of the two image sections this SDK never writes.
func SiteSettingsMediaSections() []string { return sitesettings.SiteSettingsMediaSections() }

// SiteSettingsWritableSections returns a fresh copy of the six sections a PATCH may name here.
func SiteSettingsWritableSections() []string { return sitesettings.SiteSettingsWritableSections() }

const (
	SiteSettingsPath        = sitesettings.SiteSettingsPath
	SiteSettingsSectionPath = sitesettings.SiteSettingsSectionPath
)

// ShippingSettingsResource is the shipping settings endpoints bound to a client.
type ShippingSettingsResource = shippingsettings.ShippingSettingsResource

type (
	// ShippingSettings is the area, in the one form the API sends it: the facade's; the area is flat.
	ShippingSettings = shippingsettings.ShippingSettings
	// ProfileReference is what each of the three reference fields holds: an id and its name, or null.
	ProfileReference = shippingsettings.ProfileReference
	// ShippingSettingsGroupMeta is the area's GroupMeta.
	ShippingSettingsGroupMeta = shippingsettings.GroupMeta
	// AreaMeta is the area's one metadata entry, two fields wide with no href.
	AreaMeta = shippingsettings.AreaMeta
)

// NewShippingSettings binds the resource to a requester of your own.
func NewShippingSettings(requester Requester) *ShippingSettingsResource {
	return shippingsettings.New(requester)
}

// ShippingSettingsShape is the DTOShape for the facade's body, the one shape this area sends.
func ShippingSettingsShape(decoded any) any { return shippingsettings.ShippingSettingsShape(decoded) }

// ShippingSettingsFields returns a fresh copy of the seven field names a PATCH body may carry.
func ShippingSettingsFields() []string { return shippingsettings.ShippingSettingsFields() }

// ShippingSettingsReferenceFields returns a fresh copy of the three that hold an id/name record.
func ShippingSettingsReferenceFields() []string {
	return shippingsettings.ShippingSettingsReferenceFields()
}

const ShippingSettingsPath = shippingsettings.ShippingSettingsPath

// SearchSettingsResource is the search settings endpoints bound to a client.
type SearchSettingsResource = searchsettings.SearchSettingsResource

type (
	// SearchSettings is the area, in either form the API sends it: the facade's, and one section's.
	SearchSettings = searchsettings.SearchSettings

	SearchableContent = searchsettings.SearchableContent
	SearchFields      = searchsettings.SearchFields
	Synonyms          = searchsettings.Synonyms

	SearchableContentGroup   = searchsettings.SearchableContentGroup
	ProductCategoriesContent = searchsettings.ProductCategoriesContent
	ExcludedRecord           = searchsettings.ExcludedRecord
	SearchFieldPriority      = searchsettings.SearchFieldPriority
	SynonymGroup             = searchsettings.SynonymGroup

	// ReindexStatus is the read-only sub-resource beside the sections; its body has no wrapper.
	ReindexStatus = searchsettings.ReindexStatus

	// SearchSettingsGroupMeta and SearchSettingsSectionMeta are the area's GroupMeta and SectionMeta.
	SearchSettingsGroupMeta   = searchsettings.GroupMeta
	SearchSettingsSectionMeta = searchsettings.SectionMeta
)

// NewSearchSettings binds the resource to a requester of your own.
func NewSearchSettings(requester Requester) *SearchSettingsResource {
	return searchsettings.New(requester)
}

// SearchSettingsShape is the DTOShape for the facade's body and for one section's.
func SearchSettingsShape(decoded any) any { return searchsettings.SearchSettingsShape(decoded) }

// ReindexStatusShape is the DTOShape for the reindex_status body, the one shape that route sends.
func ReindexStatusShape(decoded any) any { return searchsettings.ReindexStatusShape(decoded) }

// SearchSettingsSections returns a fresh copy of the three values {section} takes.
func SearchSettingsSections() []string { return searchsettings.SearchSettingsSections() }

// SearchSettingsReindexingSection is the one section not to write: it triggers an irreversible reindex.
const SearchSettingsReindexingSection = searchsettings.SearchSettingsReindexingSection

const (
	SearchSettingsPath              = searchsettings.SearchSettingsPath
	SearchSettingsSectionPath       = searchsettings.SearchSettingsSectionPath
	SearchSettingsReindexStatusPath = searchsettings.SearchSettingsReindexStatusPath
)

// EcommerceSettingsResource is the ecommerce settings endpoints bound to a client.
type EcommerceSettingsResource = ecommercesettings.EcommerceSettingsResource

type (
	// EcommerceSettings is the area, in either form the API sends it: the facade's, and one section's.
	EcommerceSettings = ecommercesettings.EcommerceSettings

	ContinueShopping   = ecommercesettings.ContinueShopping
	ProductInventory   = ecommercesettings.ProductInventory
	Pricing            = ecommercesettings.Pricing
	ProductRestriction = ecommercesettings.ProductRestriction
	OrderFilters       = ecommercesettings.OrderFilters
	LocationSettings   = ecommercesettings.LocationSettings

	// EcommerceSettingsGroupMeta and EcommerceSettingsSectionMeta are the area's GroupMeta pair.
	EcommerceSettingsGroupMeta   = ecommercesettings.GroupMeta
	EcommerceSettingsSectionMeta = ecommercesettings.SectionMeta
)

// NewEcommerceSettings binds the resource to a requester of your own.
func NewEcommerceSettings(requester Requester) *EcommerceSettingsResource {
	return ecommercesettings.New(requester)
}

// EcommerceSettingsShape is the DTOShape for the facade's body and for one section's.
func EcommerceSettingsShape(decoded any) any {
	return ecommercesettings.EcommerceSettingsShape(decoded)
}

// EcommerceSettingsAreaKeys returns a fresh copy of the seven keys the facade carries, in store order.
func EcommerceSettingsAreaKeys() []string { return ecommercesettings.EcommerceSettingsAreaKeys() }

// EcommerceSettingsSections returns a fresh copy of the six values {section} takes, all writable.
func EcommerceSettingsSections() []string { return ecommercesettings.EcommerceSettingsSections() }

// EcommerceSettingsCollectionKey is the one area key that is no section: a read-only collection link.
const EcommerceSettingsCollectionKey = ecommercesettings.EcommerceSettingsCollectionKey

const (
	EcommerceSettingsPath        = ecommercesettings.EcommerceSettingsPath
	EcommerceSettingsSectionPath = ecommercesettings.EcommerceSettingsSectionPath
)

// OrderFlags is the order flag endpoints bound to a client.
type OrderFlags = orderflags.OrderFlags

type (
	// OrderFlag is one flag, a listing row or a row path's {"order_flag": {…}}.
	OrderFlag = orderflags.OrderFlag
	// OrderFlagInput is the writable form of one; the update replaces, so an omitted field is blanked.
	OrderFlagInput = orderflags.OrderFlagInput
	// OrderFlagPage is one page of the listing: rows on OrderFlags, the paging block beside them.
	OrderFlagPage = orderflags.OrderFlagPage
	// OrderFlagsPagination is the order flag listing's five-key pagination block.
	OrderFlagsPagination = orderflags.Pagination
)

// NewOrderFlags binds the resource to a requester of your own.
func NewOrderFlags(requester Requester) *OrderFlags { return orderflags.New(requester) }

// OrderFlagShape is the DTOShape for a body shaped {"order_flag": {…}}, for the raw-verb call.
func OrderFlagShape(decoded any) any { return orderflags.OrderFlagShape(decoded) }

// OrderFlagPageShape is the DTOShape for the listing body.
func OrderFlagPageShape(decoded any) any { return orderflags.OrderFlagPageShape(decoded) }

// OrderFlagsQueryParams returns a fresh copy of the four names the listing accepts; others are 400.
func OrderFlagsQueryParams() []string { return orderflags.OrderFlagsQueryParams() }

const (
	OrderFlagsPath = orderflags.OrderFlagsPath
	OrderFlagPath  = orderflags.OrderFlagPath

	// OrderFlagsKey wraps the listing's rows, and OrderFlagKey one row and both writes.
	OrderFlagsKey = orderflags.OrderFlagsKey
	OrderFlagKey  = orderflags.OrderFlagKey
)

// Orders is the order endpoints bound to a client.
type Orders = orders.Orders

type (
	// Order is one order in sections: a row read, a PUT or PATCH answer, and a listing row.
	Order = orders.Order
	// FlatOrder is one order with its totals at the top: a create, a cancel, a status move, a per-customer row.
	FlatOrder = orders.FlatOrder
	// OrderCustomer is the customer section, also sent as customer_summary; OrderAddress one address.
	OrderCustomer = orders.OrderCustomer
	OrderAddress  = orders.OrderAddress
	// OrdersOrderItem is one order line.
	OrdersOrderItem = orders.OrderItem
	// OrderShipping, OrderDiscounts, OrderChargesTotals, OrderPayment and OrderFulfillment are the other sections.
	OrderShipping      = orders.OrderShipping
	OrderDiscounts     = orders.OrderDiscounts
	OrderChargesTotals = orders.OrderChargesTotals
	OrderPayment       = orders.OrderPayment
	OrderFulfillment   = orders.OrderFulfillment
	// OrdersReference is a country or a state on an order.
	OrdersReference = orders.Reference
	// OrdersGroupMeta is the order group_meta, each entry an OrdersSectionMeta.
	OrdersGroupMeta   = orders.GroupMeta
	OrdersSectionMeta = orders.SectionMeta
	// OrderPage is one page of the listing: rows on Orders, the ten-key block and group_meta beside them.
	OrderPage = orders.OrderPage
	// OrdersPagination is the order listing's ten-key pagination block.
	OrdersPagination = orders.Pagination
	// OrderCount is the count sub-path's body.
	OrderCount = orders.OrderCount
	// CustomerOrders is one customer's orders, flat.
	CustomerOrders = orders.CustomerOrders
	// OrderAddresses is the addresses sub-path's answer.
	OrderAddresses = orders.OrderAddresses
	// OrdersOrderComment is one order comment.
	OrdersOrderComment = orders.OrderComment
	OrderCommentList   = orders.OrderCommentList
	// OrdersOrderCustomFields is the custom_fields sub-path's body.
	OrdersOrderCustomFields = orders.OrderCustomFields
	// PaymentReceipt is what recording a payment answers, its customer a PaymentReceiptCustomer.
	PaymentReceipt         = orders.PaymentReceipt
	PaymentReceiptCustomer = orders.PaymentReceiptCustomer
	// OrderRefund and PaymentRefund are what the two refunds answer.
	OrderRefund   = orders.OrderRefund
	PaymentRefund = orders.PaymentRefund
	// OrderRefundPage is an order's refund listing: one OrderRefundRecord per refund, each with its OrderRefundRecordLine entries.
	OrderRefundPage       = orders.OrderRefundPage
	OrderRefundRecord     = orders.OrderRefundRecord
	OrderRefundRecordLine = orders.OrderRefundRecordLine
	// InvoiceDispatch is what sending an invoice answers.
	InvoiceDispatch = orders.InvoiceDispatch
	// ShipmentStatus is the shipment-status body, one ShipmentTracking per shipment.
	ShipmentStatus   = orders.ShipmentStatus
	ShipmentTracking = orders.ShipmentTracking
	// OrderShipment is one shipment over ShipmentItem lines; OrderShipmentList an order's shipments.
	OrderShipment     = orders.OrderShipment
	ShipmentItem      = orders.ShipmentItem
	OrderShipmentList = orders.OrderShipmentList

	// OrderInput is the create body, OrderWithPaymentInput the create with its InitialPaymentInput.
	OrderInput            = orders.OrderInput
	OrderWithPaymentInput = orders.OrderWithPaymentInput
	InitialPaymentInput   = orders.InitialPaymentInput
	// OrderAddressInput is an address to write, whole; OrderProductInput one line to order.
	OrderAddressInput = orders.OrderAddressInput
	OrderProductInput = orders.OrderProductInput
	// OrderUpdateInput and OrderPatchInput are the PUT and PATCH bodies, both in the order envelope.
	OrderUpdateInput = orders.OrderUpdateInput
	OrderPatchInput  = orders.OrderPatchInput
	// OrderAddressesInput, OrderCancelInput and OrderCommentInput are bare records.
	OrderAddressesInput = orders.OrderAddressesInput
	OrderCancelInput    = orders.OrderCancelInput
	OrderCommentInput   = orders.OrderCommentInput
	// OrderPaymentInput is in the payment envelope with transaction_no beside it.
	OrderPaymentInput = orders.OrderPaymentInput
	// OrderRefundInput and PaymentRefundInput are in the refund envelope; RefundLineInput one line refunded.
	OrderRefundInput   = orders.OrderRefundInput
	PaymentRefundInput = orders.PaymentRefundInput
	RefundLineInput    = orders.RefundLineInput
	// OrderShipmentInput is a shipment write, ShipmentItemInput one line shipped.
	OrderShipmentInput = orders.OrderShipmentInput
	ShipmentItemInput  = orders.ShipmentItemInput
	// OrderStatusInput moves an order's status; CustomOrderStatusInput adds a custom one.
	OrderStatusInput       = orders.OrderStatusInput
	CustomOrderStatusInput = orders.CustomOrderStatusInput
)

// NewOrders binds the resource to a requester of your own.
func NewOrders(requester Requester) *Orders { return orders.New(requester) }

// OrderShape is the DTOShape for a body shaped {"order": {…}} in sections.
func OrderShape(decoded any) any { return orders.OrderShape(decoded) }

// OrderPageShape is the DTOShape for the listing body.
func OrderPageShape(decoded any) any { return orders.OrderPageShape(decoded) }

// OrderCountShape is the DTOShape for the count sub-path's body.
func OrderCountShape(decoded any) any { return orders.OrderCountShape(decoded) }

// FlatOrderShape is the DTOShape for a body shaped {"order": {…}} in the flat shape.
func FlatOrderShape(decoded any) any { return orders.FlatOrderShape(decoded) }

// CustomerOrdersShape is the DTOShape for one customer's orders.
func CustomerOrdersShape(decoded any) any { return orders.CustomerOrdersShape(decoded) }

// OrderAddressesShape is the DTOShape for the addresses sub-path's body.
func OrderAddressesShape(decoded any) any { return orders.OrderAddressesShape(decoded) }

// OrdersOrderCommentShape is the DTOShape for a body shaped {"comment": {…}}.
func OrdersOrderCommentShape(decoded any) any { return orders.OrderCommentShape(decoded) }

// OrderCommentListShape is the DTOShape for the comments sub-path's body.
func OrderCommentListShape(decoded any) any { return orders.OrderCommentListShape(decoded) }

// OrdersOrderCustomFieldsShape is the DTOShape for the custom_fields sub-path's body.
func OrdersOrderCustomFieldsShape(decoded any) any { return orders.OrderCustomFieldsShape(decoded) }

// PaymentReceiptShape is the DTOShape for a body shaped {"payment": {…}}.
func PaymentReceiptShape(decoded any) any { return orders.PaymentReceiptShape(decoded) }

// OrderRefundShape is the DTOShape for an order refund's answer.
func OrderRefundShape(decoded any) any { return orders.OrderRefundShape(decoded) }

// OrderRefundPageShape is the DTOShape for an order's refund listing.
func OrderRefundPageShape(decoded any) any { return orders.OrderRefundPageShape(decoded) }

// PaymentRefundShape is the DTOShape for a payment refund's answer.
func PaymentRefundShape(decoded any) any { return orders.PaymentRefundShape(decoded) }

// InvoiceDispatchShape is the DTOShape for the send-invoice answer.
func InvoiceDispatchShape(decoded any) any { return orders.InvoiceDispatchShape(decoded) }

// ShipmentStatusShape is the DTOShape for the shipment-status body.
func ShipmentStatusShape(decoded any) any { return orders.ShipmentStatusShape(decoded) }

// OrderShipmentShape is the DTOShape for a body shaped {"shipment": {…}}.
func OrderShipmentShape(decoded any) any { return orders.OrderShipmentShape(decoded) }

// OrderShipmentListShape is the DTOShape for the shipments sub-path's body.
func OrderShipmentListShape(decoded any) any { return orders.OrderShipmentListShape(decoded) }

// OrdersSortFields returns a fresh copy of the ten values the listing's sort names.
func OrdersSortFields() []string { return orders.OrdersSortFields() }

const (
	OrdersPath                  = orders.OrdersPath
	OrdersCountPath             = orders.OrdersCountPath
	OrdersWithPaymentPath       = orders.OrdersWithPaymentPath
	OrderPath                   = orders.OrderPath
	OrderAddressesPath          = orders.OrderAddressesPath
	OrderCancelPath             = orders.OrderCancelPath
	OrderCommentsPath           = orders.OrderCommentsPath
	OrderPaymentsPath           = orders.OrderPaymentsPath
	OrderRefundsPath            = orders.OrderRefundsPath
	OrderShipmentStatusPath     = orders.OrderShipmentStatusPath
	OrderShipmentsPath          = orders.OrderShipmentsPath
	OrderShipmentPath           = orders.OrderShipmentPath
	OrderUpdateStatusPath       = orders.OrderUpdateStatusPath
	OrdersOrderCustomFieldsPath = orders.OrderCustomFieldsPath
	OrderSendInvoicePath        = orders.OrderSendInvoicePath
	PaymentRefundsPath          = orders.PaymentRefundsPath
	OrderStatusAddPath          = orders.OrderStatusAddPath
	CustomerOrdersPath          = orders.CustomerOrdersPath

	// OrdersKey wraps the listing's rows, a customer's orders and the count record; OrderKey one order.
	OrdersKey = orders.OrdersKey
	OrderKey  = orders.OrderKey
	// The comments and one comment, the shipments and one shipment, and the payment wrapper and receipt key.
	OrderCommentsKey  = orders.OrderCommentsKey
	OrderCommentKey   = orders.OrderCommentKey
	OrderShipmentsKey = orders.OrderShipmentsKey
	OrderShipmentKey  = orders.OrderShipmentKey
	OrderPaymentKey   = orders.OrderPaymentKey

	// OrdersPageSize is the page size when limit is not sent and the ceiling on one that is.
	OrdersPageSize = orders.OrdersPageSize
	// OrdersIdempotencyKeyHeader is the header CreatePaymentRefund sends its key under.
	OrdersIdempotencyKeyHeader = orders.OrdersIdempotencyKeyHeader
	// OrdersTransactionNoHeader is the header CreateOrder sends transactionNo under.
	OrdersTransactionNoHeader = orders.OrdersTransactionNoHeader
)

// Aliases is the domain alias endpoints bound to a client.
type Aliases = aliases.Aliases

type (
	// Alias is one alias, a listing row or a create's {"alias": {…}}: the same four fields.
	Alias = aliases.Alias
	// AliasInput is the create body, {"alias": "<hostname>"}: a bare field, no envelope.
	AliasInput = aliases.AliasInput
	// AliasPage is one page of the listing: rows on Aliases, the five-key block beside them.
	AliasPage = aliases.AliasPage
	// AliasesPagination is the alias listing's five-key pagination block.
	AliasesPagination = aliases.Pagination
)

// NewAliases binds the resource to a requester of your own.
func NewAliases(requester Requester) *Aliases { return aliases.New(requester) }

// AliasShape is the DTOShape for a body shaped {"alias": {…}}, for the raw-verb call.
func AliasShape(decoded any) any { return aliases.AliasShape(decoded) }

// AliasPageShape is the DTOShape for the listing body.
func AliasPageShape(decoded any) any { return aliases.AliasPageShape(decoded) }

// AliasesQueryParams returns a fresh copy of the four names the listing accepts.
func AliasesQueryParams() []string { return aliases.AliasesQueryParams() }

// AliasesPaginationKeys returns a fresh copy of the five keys the listing's pagination block carries.
func AliasesPaginationKeys() []string { return aliases.AliasesPaginationKeys() }

const (
	AliasesPath = aliases.AliasesPath
	AliasPath   = aliases.AliasPath

	// AliasesKey wraps the listing's rows; AliasKey is the create body's one field and the new row's key.
	AliasesKey = aliases.AliasesKey
	AliasKey   = aliases.AliasKey

	// AliasesPageSize is the page size when limit is not sent and the ceiling on one that is.
	AliasesPageSize = aliases.AliasesPageSize
)

// Carts is the cart endpoints bound to a client.
type Carts = carts.Carts

type (
	// Cart is one cart.
	Cart                = carts.Cart
	CartTotals          = carts.CartTotals
	CartShipping        = carts.CartShipping
	CartItem            = carts.CartItem
	CartItemPricing     = carts.CartItemPricing
	CartItemProduct     = carts.CartItemProduct
	CartDiscount        = carts.CartDiscount
	CartMessages        = carts.CartMessages
	CartDiscountMessage = carts.CartDiscountMessage
	// CartsGroupMeta and CartsSectionMeta are the read's group_meta block and its one entry.
	CartsGroupMeta   = carts.GroupMeta
	CartsSectionMeta = carts.SectionMeta
	// CartSession is the session write's answer, {"sessionId": …}.
	CartSession = carts.CartSession
	// CartOrder is the order a cart became, its Status read from beside the order envelope.
	CartOrder         = carts.CartOrder
	CartOrderCustomer = carts.CartOrderCustomer
	CartOrderLine     = carts.CartOrderLine
	CartOrderPayment  = carts.CartOrderPayment
	CartOrderAddress  = carts.CartOrderAddress
	// CartsReference is a country or a state on the order a cart became.
	CartsReference      = carts.Reference
	CartItemInput       = carts.CartItemInput
	CartItemUpdateInput = carts.CartItemUpdateInput
	CartAddressInput    = carts.CartAddressInput
	CartSessionInput    = carts.CartSessionInput
	CartOrderInput      = carts.CartOrderInput
)

// NewCarts binds the resource to a requester of your own.
func NewCarts(requester Requester) *Carts { return carts.New(requester) }

// CartShape is the DTOShape for a body shaped {"cart": {…}}, for the raw-verb call.
func CartShape(decoded any) any { return carts.CartShape(decoded) }

// CartSessionShape is the DTOShape for the session write's answer.
func CartSessionShape(decoded any) any { return carts.CartSessionShape(decoded) }

// CartOrderShape is the DTOShape for the order write's answer, {"status": …, "order": {…}}.
func CartOrderShape(decoded any) any { return carts.CartOrderShape(decoded) }

const (
	CartsPath     = carts.CartsPath
	CartPath      = carts.CartPath
	CartItemPath  = carts.CartItemPath
	CartOrderPath = carts.CartOrderPath

	CartKey        = carts.CartKey
	CartSessionKey = carts.CartSessionKey
	CartOrderKey   = carts.CartOrderKey

	// CartsTransactionNoHeader is the header CreateOrderFromCart sends transactionNo under.
	CartsTransactionNoHeader = carts.CartsTransactionNoHeader
)

// Blogs is the blog endpoints, the comment sub-collection included, bound to a client.
type Blogs = blogs.Blogs

type (
	// Blog is one blog, the same eighteen fields on a listing row, a row read and every write's answer.
	Blog = blogs.Blog
	// BlogInput is the writable form, in the blog envelope, for the create, the replace and the patch.
	BlogInput = blogs.BlogInput
	// BlogCategoryReference is a category a blog is filed under; a write names it by ID.
	BlogCategoryReference = blogs.BlogCategoryReference
	// BlogImage is the cover image a read answers; BlogImageInput the one a write uploads.
	BlogImage      = blogs.BlogImage
	BlogImageInput = blogs.BlogImageInput
	// BlogAuthor is created_by.
	BlogAuthor = blogs.BlogAuthor
	// BlogAudience is selected_customers once it names anyone, over BlogAudienceMember; BlogAudienceInput is its writable form.
	BlogAudience       = blogs.BlogAudience
	BlogAudienceMember = blogs.BlogAudienceMember
	BlogAudienceInput  = blogs.BlogAudienceInput
	// BlogSeoConfigs is seo_configs, holding BlogMetaTag entries, the same on a write and a read.
	BlogSeoConfigs = blogs.BlogSeoConfigs
	BlogMetaTag    = blogs.BlogMetaTag
	// BlogPage is one page of the listing: rows on Blogs, group_meta and pagination beside them.
	BlogPage = blogs.BlogPage
	// BlogCount is the count sub-path's body.
	BlogCount = blogs.BlogCount
	// BlogComment is one comment; BlogCommentList and BlogCommentEnvelope carry the blog's id beside it.
	BlogComment         = blogs.BlogComment
	BlogCommentList     = blogs.BlogCommentList
	BlogCommentEnvelope = blogs.BlogCommentEnvelope
	// BlogCommentInput is the comment create's body, {"blog": {"comment": …}}, the comment a BlogCommentDraft.
	BlogCommentInput = blogs.BlogCommentInput
	BlogCommentDraft = blogs.BlogCommentDraft
	// BlogsPagination is the blog listing's ten-key pagination block.
	BlogsPagination = blogs.Pagination
	// BlogsGroupMeta is the blog collection's group_meta, its one entry a BlogsResourceMeta.
	BlogsGroupMeta    = blogs.GroupMeta
	BlogsResourceMeta = blogs.ResourceMeta
)

// NewBlogs binds the resource to a requester of your own.
func NewBlogs(requester Requester) *Blogs { return blogs.New(requester) }

// BlogShape is the DTOShape for a body shaped {"blog": {…}}.
func BlogShape(decoded any) any { return blogs.BlogShape(decoded) }

// BlogPageShape is the DTOShape for the listing body.
func BlogPageShape(decoded any) any { return blogs.BlogPageShape(decoded) }

// BlogCountShape is the DTOShape for the count sub-path's body.
func BlogCountShape(decoded any) any { return blogs.BlogCountShape(decoded) }

// BlogCommentListShape is the DTOShape for the comments sub-path's body.
func BlogCommentListShape(decoded any) any { return blogs.BlogCommentListShape(decoded) }

// BlogCommentEnvelopeShape is the DTOShape for one comment's body.
func BlogCommentEnvelopeShape(decoded any) any { return blogs.BlogCommentEnvelopeShape(decoded) }

// BlogsPaginationKeys returns a fresh copy of the ten keys the listing's pagination block carries.
func BlogsPaginationKeys() []string { return blogs.BlogsPaginationKeys() }

const (
	BlogsPath               = blogs.BlogsPath
	BlogsCountPath          = blogs.BlogsCountPath
	BlogPath                = blogs.BlogPath
	BlogCommentsPath        = blogs.BlogCommentsPath
	BlogCommentPath         = blogs.BlogCommentPath
	BlogCommentApprovePath  = blogs.BlogCommentApprovePath
	BlogCommentMarkSpamPath = blogs.BlogCommentMarkSpamPath

	// BlogsKey wraps the listing's rows and the count record; BlogKey one row and every write.
	BlogsKey = blogs.BlogsKey
	BlogKey  = blogs.BlogKey
	// BlogsCountKey is the count record's one key; BlogCommentsKey and BlogCommentKey sit beside the blog's id.
	BlogsCountKey   = blogs.BlogsCountKey
	BlogCommentsKey = blogs.BlogCommentsKey
	BlogCommentKey  = blogs.BlogCommentKey

	// BlogsPageSize is the page size when limit is not sent and the ceiling on one that is.
	BlogsPageSize = blogs.BlogsPageSize
)

// OrderCustomFields is the order custom field endpoints bound to a client.
type OrderCustomFields = ordercustomfields.OrderCustomFields

type (
	// OrderCustomField is one field, a listing row or a row path's {"order_custom_field": {…}}.
	OrderCustomField = ordercustomfields.OrderCustomField
	// OrderCustomFieldTimeslotInput is the validate_timeslot body, a bare record with no envelope.
	OrderCustomFieldTimeslotInput = ordercustomfields.OrderCustomFieldTimeslotInput
	// OrderCustomFieldPage is one page of the listing: rows on OrderCustomFields, paging and group_meta beside them.
	OrderCustomFieldPage = ordercustomfields.OrderCustomFieldPage
	// OrderCustomFieldsPagination is the order custom field listing's ten-key pagination block.
	OrderCustomFieldsPagination = ordercustomfields.Pagination
	// OrderCustomFieldsGroupMeta is the listing's group_meta, and OrderCustomFieldsResourceMeta its one entry.
	OrderCustomFieldsGroupMeta    = ordercustomfields.GroupMeta
	OrderCustomFieldsResourceMeta = ordercustomfields.ResourceMeta
)

// NewOrderCustomFields binds the resource to a requester of your own.
func NewOrderCustomFields(requester Requester) *OrderCustomFields {
	return ordercustomfields.New(requester)
}

// OrderCustomFieldShape is the DTOShape for a body shaped {"order_custom_field": {…}}.
func OrderCustomFieldShape(decoded any) any {
	return ordercustomfields.OrderCustomFieldShape(decoded)
}

// OrderCustomFieldPageShape is the DTOShape for the listing body.
func OrderCustomFieldPageShape(decoded any) any {
	return ordercustomfields.OrderCustomFieldPageShape(decoded)
}

// OrderCustomFieldsPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func OrderCustomFieldsPaginationKeys() []string {
	return ordercustomfields.OrderCustomFieldsPaginationKeys()
}

const (
	OrderCustomFieldsPath                = ordercustomfields.OrderCustomFieldsPath
	OrderCustomFieldsCountPath           = ordercustomfields.OrderCustomFieldsCountPath
	OrderCustomFieldPath                 = ordercustomfields.OrderCustomFieldPath
	OrderCustomFieldValidateTimeslotPath = ordercustomfields.OrderCustomFieldValidateTimeslotPath

	// OrderCustomFieldsKey wraps the rows and keys the one group_meta entry, OrderCustomFieldKey one row, OrderCustomFieldsCountKey the count.
	OrderCustomFieldsKey      = ordercustomfields.OrderCustomFieldsKey
	OrderCustomFieldKey       = ordercustomfields.OrderCustomFieldKey
	OrderCustomFieldsCountKey = ordercustomfields.OrderCustomFieldsCountKey

	// OrderCustomFieldsPageSize is the page size when limit is not sent, and the ceiling on one that is.
	OrderCustomFieldsPageSize = ordercustomfields.OrderCustomFieldsPageSize
)

// BlogCategories is the blog category endpoints bound to a client.
type BlogCategories = blogcategories.BlogCategories

type (
	// BlogCategory is one category, a listing row or a row path's {"blog_category": {…}}: eleven fields.
	BlogCategory = blogcategories.BlogCategory
	// BlogCategoryImage is a category's image, {} until one is uploaded.
	BlogCategoryImage = blogcategories.BlogCategoryImage
	// BlogCategorySeoConfig is one seo_configs entry, which no write here sets.
	BlogCategorySeoConfig = blogcategories.BlogCategorySeoConfig
	// BlogCategoryInput is the body both writes take, wrapper included.
	BlogCategoryInput = blogcategories.BlogCategoryInput
	// BlogCategoryImageInput is an image upload, BlogCategoryBlog a blog the category is linked to.
	BlogCategoryImageInput = blogcategories.BlogCategoryImageInput
	BlogCategoryBlog       = blogcategories.BlogCategoryBlog
	// BlogCategoryPage is one page of the listing: rows on BlogCategories, the ten-key block beside them.
	BlogCategoryPage = blogcategories.BlogCategoryPage
	// BlogCategoriesPagination is the blog category listing's ten-key pagination block.
	BlogCategoriesPagination = blogcategories.Pagination
)

// NewBlogCategories binds the resource to a requester of your own.
func NewBlogCategories(requester Requester) *BlogCategories { return blogcategories.New(requester) }

// BlogCategoryShape is the DTOShape for a body shaped {"blog_category": {…}}, for the raw-verb call.
func BlogCategoryShape(decoded any) any { return blogcategories.BlogCategoryShape(decoded) }

// BlogCategoryPageShape is the DTOShape for the listing body.
func BlogCategoryPageShape(decoded any) any { return blogcategories.BlogCategoryPageShape(decoded) }

// BlogCategoriesPaginationKeys returns a fresh copy of the ten keys the listing's pagination block carries.
func BlogCategoriesPaginationKeys() []string { return blogcategories.BlogCategoriesPaginationKeys() }

const (
	BlogCategoriesPath = blogcategories.BlogCategoriesPath
	BlogCategoryPath   = blogcategories.BlogCategoryPath

	// BlogCategoriesKey wraps the listing's rows, and BlogCategoryKey one row and both writes.
	BlogCategoriesKey = blogcategories.BlogCategoriesKey
	BlogCategoryKey   = blogcategories.BlogCategoryKey

	// BlogCategoriesPageSize is the page size when limit is not sent and the ceiling on one that is.
	BlogCategoriesPageSize = blogcategories.BlogCategoriesPageSize
)

// InventoryLocations is the inventory location endpoints bound to a client.
type InventoryLocations = inventorylocations.InventoryLocations

type (
	// InventoryLocation is one location, in every form the API sends it.
	InventoryLocation = inventorylocations.InventoryLocation
	// InventoryLocationInput is the writable form of one, with no envelope; the update replaces and the patch merges.
	InventoryLocationInput = inventorylocations.InventoryLocationInput
	// InventoryLocationStatusInput is the status sub-path's body; it accepts active only.
	InventoryLocationStatusInput = inventorylocations.InventoryLocationStatusInput
	// InventoryLocationPage is one page of the listing: rows on InventoryLocations, paging beside them.
	InventoryLocationPage = inventorylocations.InventoryLocationPage
	// InventoryLocationStock is the stock sub-path's body.
	InventoryLocationStock        = inventorylocations.InventoryLocationStock
	InventoryLocationStockItem    = inventorylocations.InventoryLocationStockItem
	InventoryLocationStockProduct = inventorylocations.InventoryLocationStockProduct
	// InventoryLocationsPagination is the inventory location listing's five-key pagination block.
	InventoryLocationsPagination = inventorylocations.Pagination
)

// NewInventoryLocations binds the resource to a requester of your own.
func NewInventoryLocations(requester Requester) *InventoryLocations {
	return inventorylocations.New(requester)
}

// InventoryLocationShape is the DTOShape for a body shaped {"inventory_location": {…}}.
func InventoryLocationShape(decoded any) any {
	return inventorylocations.InventoryLocationShape(decoded)
}

// InventoryLocationPageShape is the DTOShape for the listing body.
func InventoryLocationPageShape(decoded any) any {
	return inventorylocations.InventoryLocationPageShape(decoded)
}

// InventoryLocationStockShape is the DTOShape for the stock sub-path's body.
func InventoryLocationStockShape(decoded any) any {
	return inventorylocations.InventoryLocationStockShape(decoded)
}

// InventoryLocationsQueryParams returns a fresh copy of the three names the listing accepts.
func InventoryLocationsQueryParams() []string {
	return inventorylocations.InventoryLocationsQueryParams()
}

// InventoryLocationsUnroutedMethodBody is the 200 body the store answers for an unrouted method.
func InventoryLocationsUnroutedMethodBody() map[string]any {
	return inventorylocations.UnroutedMethodBody()
}

const (
	InventoryLocationsPath       = inventorylocations.InventoryLocationsPath
	InventoryLocationPath        = inventorylocations.InventoryLocationPath
	InventoryLocationStockPath   = inventorylocations.InventoryLocationStockPath
	InventoryLocationDefaultPath = inventorylocations.InventoryLocationDefaultPath
	InventoryLocationStatusPath  = inventorylocations.InventoryLocationStatusPath

	// InventoryLocationsKey wraps the listing's rows, InventoryLocationKey one row and InventoryLocationStockKey the stock rows.
	InventoryLocationsKey     = inventorylocations.InventoryLocationsKey
	InventoryLocationKey      = inventorylocations.InventoryLocationKey
	InventoryLocationStockKey = inventorylocations.InventoryLocationStockKey

	// InventoryLocationsPageSize is the page size applied when no limit is sent.
	InventoryLocationsPageSize = inventorylocations.InventoryLocationsPageSize
)

// Products is the product, review import and stock endpoints bound to a client.
type Products = products.Products

type (
	// Product is one product, the row read's {"product": {…}} and every row write's answer.
	Product = products.Product
	// ProductInput is the create's, the PUT's and the PATCH's body, a bare record.
	ProductInput = products.ProductInput
	// ProductReviews is the review import's answer; ProductReviewsImportInput its body, ProductReviewInput one review.
	ProductReviews            = products.ProductReviews
	ProductReviewsImportInput = products.ProductReviewsImportInput
	ProductReviewInput        = products.ProductReviewInput
	// ProductStock is the stock sub-path's body, read and written, and ProductStockLocation one location's row.
	ProductStock          = products.ProductStock
	ProductStockLocation  = products.ProductStockLocation
	ProductStockVariation = products.ProductStockVariation
	// ProductStockInput is the stock write's body, ProductStockLocationInput one of its rows.
	ProductStockInput         = products.ProductStockInput
	ProductStockLocationInput = products.ProductStockLocationInput
	// ProductUrls through ProductCombinationPriceAndStock are the blocks a product carries.
	ProductUrls                     = products.ProductUrls
	ProductNamedReference           = products.ProductNamedReference
	ProductParent                   = products.ProductParent
	ProductSelectedCustomers        = products.ProductSelectedCustomers
	ProductWeekday                  = products.ProductWeekday
	ProductAvailableOnWeekdays      = products.ProductAvailableOnWeekdays
	ProductProperties               = products.ProductProperties
	ProductTaxSettings              = products.ProductTaxSettings
	ProductVariationPriceDisplay    = products.ProductVariationPriceDisplay
	ProductWarehouseStock           = products.ProductWarehouseStock
	ProductPricingAndStock          = products.ProductPricingAndStock
	ProductImage                    = products.ProductImage
	ProductImageAndVideo            = products.ProductImageAndVideo
	ProductAdvanceInformation       = products.ProductAdvanceInformation
	ProductAdvanced                 = products.ProductAdvanced
	ProductLink                     = products.ProductLink
	ProductBundle                   = products.ProductBundle
	ProductSeoConfig                = products.ProductSeoConfig
	ProductTagValue                 = products.ProductTagValue
	ProductTag                      = products.ProductTag
	ProductCustomProperty           = products.ProductCustomProperty
	ProductReview                   = products.ProductReview
	ProductVariationTypeOption      = products.ProductVariationTypeOption
	ProductVariationType            = products.ProductVariationType
	ProductCombinationOptionType    = products.ProductCombinationOptionType
	ProductCombinationOption        = products.ProductCombinationOption
	ProductCombinationImage         = products.ProductCombinationImage
	ProductCombinationPriceAndStock = products.ProductCombinationPriceAndStock
	ProductCombinationDetails       = products.ProductCombinationDetails
	ProductCombination              = products.ProductCombination
	ProductVariation                = products.ProductVariation
	// ProductsProductVideo is a product's video block.
	ProductsProductVideo = products.ProductVideo
	// ProductsPage is one page of the listing: full records on Products, the five-key block beside them.
	ProductsPage = products.ProductPage
	// ProductsPagination is the product listing's five-key pagination block.
	ProductsPagination = products.Pagination
	// ProductSelectedCustomer through ProductLinks are the blocks the record's block-shaped copy carries.
	ProductSelectedCustomer          = products.ProductSelectedCustomer
	ProductSelectedCustomersInput    = products.ProductSelectedCustomersInput
	ProductFlag                      = products.ProductFlag
	ProductFlags                     = products.ProductFlags
	ProductDimension                 = products.ProductDimension
	ProductIdentifiers               = products.ProductIdentifiers
	ProductWebtool                   = products.ProductWebtool
	ProductReviewer                  = products.ProductReviewer
	ProductCategory                  = products.ProductCategory
	ProductMediaImage                = products.ProductMediaImage
	ProductMediaVideo                = products.ProductMediaVideo
	ProductMedia                     = products.ProductMedia
	ProductWeekdayWindow             = products.ProductWeekdayWindow
	ProductScheduleWeekdays          = products.ProductScheduleWeekdays
	ProductSchedule                  = products.ProductSchedule
	ProductPasswordProtection        = products.ProductPasswordProtection
	ProductAvailableForPurchase      = products.ProductAvailableForPurchase
	ProductAvailableForPurchaseInput = products.ProductAvailableForPurchaseInput
	ProductOnSale                    = products.ProductOnSale
	ProductQuantityMultiple          = products.ProductQuantityMultiple
	ProductExpectToPay               = products.ProductExpectToPay
	ProductAdvancedPricing           = products.ProductAdvancedPricing
	ProductPricing                   = products.ProductPricing
	ProductInventoryTracking         = products.ProductInventoryTracking
	ProductCombinedStockLimit        = products.ProductCombinedStockLimit
	ProductOrdering                  = products.ProductOrdering
	ProductRestrictions              = products.ProductRestrictions
	ProductPriceIncludesTax          = products.ProductPriceIncludesTax
	ProductTax                       = products.ProductTax
	ProductShipping                  = products.ProductShipping
	ProductLinks                     = products.ProductLinks
	// ProductsProductInventory and ProductsProductRestriction are the block's inventory and one restriction.
	ProductsProductInventory   = products.ProductInventory
	ProductsProductRestriction = products.ProductRestriction
)

// NewProducts binds the resource to a requester of your own.
func NewProducts(requester Requester) *Products {
	return products.New(requester)
}

// ProductShape is the DTOShape for a body shaped {"product": {…}}.
func ProductShape(decoded any) any {
	return products.ProductShape(decoded)
}

// ProductsPageShape is the DTOShape for the listing body, {"products": […], "pagination": {…}}.
func ProductsPageShape(decoded any) any {
	return products.ProductPageShape(decoded)
}

// ProductReviewsShape is the DTOShape for the review import's answer.
func ProductReviewsShape(decoded any) any {
	return products.ProductReviewsShape(decoded)
}

// ProductStockShape is the DTOShape for the stock sub-path's body.
func ProductStockShape(decoded any) any {
	return products.ProductStockShape(decoded)
}

// ProductQueryParams returns a fresh copy of the one name the row read accepts.
func ProductQueryParams() []string {
	return products.ProductQueryParams()
}

const (
	ProductsPath             = products.ProductsPath
	ProductPath              = products.ProductPath
	ProductReviewsPath       = products.ProductReviewsPath
	ProductStockPath         = products.ProductStockPath
	ProductStockLocationPath = products.ProductStockLocationPath

	// ProductKey is the key a product and the review import answer under, ProductStockLocationsKey the stock rows' key.
	ProductKey               = products.ProductKey
	ProductStockLocationsKey = products.ProductStockLocationsKey

	// ProductsTransactionNoHeader is the header CreateProduct and ImportProductReviews send transactionNo under.
	ProductsTransactionNoHeader = products.ProductsTransactionNoHeader
)

// Categories is the category endpoints bound to a client.
type Categories = categories.Categories

type (
	// Category is one category: a listing row, a row path's {"category": {…}}, a write's answer, or one section's.
	Category = categories.Category
	// CategoryGeneral through CategoryAssignProducts are the six sections as read; CategoryReference is every {id, name}.
	CategoryGeneral                = categories.CategoryGeneral
	CategoryDescription            = categories.CategoryDescription
	CategoryImagesMedia            = categories.CategoryImagesMedia
	CategoryAvailabilityVisibility = categories.CategoryAvailabilityVisibility
	CategorySelectedCustomers      = categories.CategorySelectedCustomers
	CategoryTaxShipping            = categories.CategoryTaxShipping
	CategoryAssignProducts         = categories.CategoryAssignProducts
	CategoryProduct                = categories.CategoryProduct
	CategoryReference              = categories.CategoryReference
	// CategoryInput is the body every write takes, wrapper included, over the section inputs below.
	CategoryInput                       = categories.CategoryInput
	CategoryGeneralInput                = categories.CategoryGeneralInput
	CategoryImagesMediaInput            = categories.CategoryImagesMediaInput
	CategoryImage                       = categories.CategoryImage
	CategoryAvailabilityVisibilityInput = categories.CategoryAvailabilityVisibilityInput
	CategorySelectedCustomersInput      = categories.CategorySelectedCustomersInput
	CategoryAssignProductsInput         = categories.CategoryAssignProductsInput
	// CategoryCount is the count sub-path's body — the plural key wrapping a record.
	CategoryCount = categories.CategoryCount
	// CategoriesPage is one page of the category listing.
	CategoriesPage = categories.CategoryPage
	// CategoriesPagination, CategoriesGroupMeta and CategoriesSectionMeta are the category listing's blocks.
	CategoriesPagination  = categories.Pagination
	CategoriesGroupMeta   = categories.GroupMeta
	CategoriesSectionMeta = categories.SectionMeta
)

// NewCategories binds the resource to a requester of your own.
func NewCategories(requester Requester) *Categories { return categories.New(requester) }

// CategoryShape is the DTOShape for a body shaped {"category": {…}}, for the raw-verb call.
func CategoryShape(decoded any) any { return categories.CategoryShape(decoded) }

// CategoriesPageShape is the DTOShape for the listing body.
func CategoriesPageShape(decoded any) any { return categories.CategoryPageShape(decoded) }

// CategoryCountShape is the DTOShape for the count sub-path's body.
func CategoryCountShape(decoded any) any { return categories.CategoryCountShape(decoded) }

// CategorySections returns a fresh copy of the six values {section} takes.
func CategorySections() []string { return categories.CategorySections() }

// CategoriesPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func CategoriesPaginationKeys() []string { return categories.CategoriesPaginationKeys() }

const (
	CategoriesPath          = categories.CategoriesPath
	CategoriesCountPath     = categories.CategoriesCountPath
	CategoriesCheckNamePath = categories.CategoriesCheckNamePath
	CategoryPath            = categories.CategoryPath
	CategorySectionPath     = categories.CategorySectionPath

	// CategoriesKey wraps the rows and the count's record, CategoryKey one row and every write.
	CategoriesKey      = categories.CategoriesKey
	CategoryKey        = categories.CategoryKey
	CategoriesCountKey = categories.CategoriesCountKey

	// CategoriesPageSize is the page size when limit is not sent and the ceiling on one that is.
	CategoriesPageSize = categories.CategoriesPageSize
	// CategoriesTransactionNoHeader is the header CreateCategory sends transactionNo under.
	CategoriesTransactionNoHeader = categories.CategoriesTransactionNoHeader
)

// Cities is the city endpoints bound to a client; it reads only, under /app.
type Cities = cities.Cities

type (
	// City is one listing row.
	City = cities.City
	// CityState is the state a row carries.
	CityState = cities.CityState
	// CityCountry is the country a row carries.
	CityCountry = cities.CityCountry
	// CityPage is one page of the listing: rows on Cities, paging beside them.
	CityPage = cities.CityPage
	// CitiesPagination is the listing's five-key pagination block.
	CitiesPagination = cities.Pagination
)

// NewCities binds the resource to a requester of your own.
func NewCities(requester Requester) *Cities { return cities.New(requester) }

// CityPageShape is the DTOShape for the listing body.
func CityPageShape(decoded any) any { return cities.CityPageShape(decoded) }

const (
	CitiesPath = cities.CitiesPath

	// CitiesKey wraps the rows.
	CitiesKey = cities.CitiesKey

	// CitiesPageSize is the page size when limit is not sent, and the largest limit accepted.
	CitiesPageSize = cities.CitiesPageSize
)

// CommonSettingsResource is the store-wide settings and the registration form bound to a client, outside /admin.
type CommonSettingsResource = commonsettings.CommonSettingsResource

type (
	// CommonSettings is what the settings read and both writes answer, groups keyed by name.
	CommonSettings = commonsettings.CommonSettings
	// CommonSettingsInput is the list form both writes take, over Setting.
	CommonSettingsInput = commonsettings.CommonSettingsInput
	Setting             = commonsettings.Setting
	// CustomerFields is the registration form, its configuration on CustomerFieldConfigs.
	CustomerFields       = commonsettings.CustomerFields
	CustomerFieldConfigs = commonsettings.CustomerFieldConfigs
)

// NewCommonSettings binds the resource to a requester of your own.
func NewCommonSettings(requester Requester) *CommonSettingsResource {
	return commonsettings.New(requester)
}

// CommonSettingsShape is the DTOShape for a body shaped {"settings": {…}}.
func CommonSettingsShape(decoded any) any { return commonsettings.CommonSettingsShape(decoded) }

// CustomerFieldsShape is the DTOShape for the registration form's body.
func CustomerFieldsShape(decoded any) any { return commonsettings.CustomerFieldsShape(decoded) }

const (
	CommonSettingsPath       = commonsettings.CommonSettingsPath
	CreateCustomerFieldsPath = commonsettings.CreateCustomerFieldsPath

	// CommonSettingsKey wraps the groups; CustomerFieldsKey and CustomerFieldConfigsKey are the registration form's two keys.
	CommonSettingsKey       = commonsettings.CommonSettingsKey
	CustomerFieldsKey       = commonsettings.CustomerFieldsKey
	CustomerFieldConfigsKey = commonsettings.CustomerFieldConfigsKey
)

// Countries is the country endpoints bound to a client; it reads only.
type Countries = countries.Countries

type (
	// Country is one country, a listing row or a row path's {"country": {…}}.
	Country = countries.Country
	// CountryReference is the short parent {id, code, name} beside a state listing or read.
	CountryReference = countries.CountryReference
	// CountryState is one state, a listing row or a state path's body with its parent on Country.
	CountryState = countries.CountryState
	// CountryPage is one page of the country listing: rows on Countries, paging beside them.
	CountryPage = countries.CountryPage
	// CountryStatePage is one page of a country's states: the parent, rows on States, paging beside them.
	CountryStatePage = countries.CountryStatePage
	// CountriesPagination is the two listings' five-key pagination block.
	CountriesPagination = countries.Pagination
)

// NewCountries binds the resource to a requester of your own.
func NewCountries(requester Requester) *Countries { return countries.New(requester) }

// CountryShape is the DTOShape for a body shaped {"country": {…}}.
func CountryShape(decoded any) any { return countries.CountryShape(decoded) }

// CountryPageShape is the DTOShape for the country listing body.
func CountryPageShape(decoded any) any { return countries.CountryPageShape(decoded) }

// CountryStateShape is the DTOShape for a body shaped {"country": {…}, "state": {…}}.
func CountryStateShape(decoded any) any { return countries.CountryStateShape(decoded) }

// CountryStatePageShape is the DTOShape for a state listing body.
func CountryStatePageShape(decoded any) any { return countries.CountryStatePageShape(decoded) }

// CountriesSortFields returns a fresh copy of the three names sort takes on both listings.
func CountriesSortFields() []string { return countries.CountriesSortFields() }

const (
	CountriesPath     = countries.CountriesPath
	CountryPath       = countries.CountryPath
	CountryStatesPath = countries.CountryStatesPath
	CountryStatePath  = countries.CountryStatePath

	// CountriesKey wraps the rows, CountryKey one country and a state's parent, CountryStatesKey a state listing's rows and CountryStateKey one state.
	CountriesKey     = countries.CountriesKey
	CountryKey       = countries.CountryKey
	CountryStatesKey = countries.CountryStatesKey
	CountryStateKey  = countries.CountryStateKey

	// CountriesPageSize is the page size when limit is not sent, CountriesMaxLimit the largest numeric one and CountriesLimitAll the value that answers every row.
	CountriesPageSize = countries.CountriesPageSize
	CountriesMaxLimit = countries.CountriesMaxLimit
	CountriesLimitAll = countries.CountriesLimitAll
)

// Currencies is the currency endpoints bound to a client.
type Currencies = currencies.Currencies

type (
	// Currency is one currency, a listing row, a row path's {"currency": {…}} or a write's answer.
	Currency = currencies.Currency
	// CurrencyCountry is the country a row is filed under.
	CurrencyCountry = currencies.CurrencyCountry
	// CurrencyInput is the writable form for the create, the replace and the merge, with no envelope.
	CurrencyInput = currencies.CurrencyInput
	// CurrencyPage is one page of the listing: rows on Currencies, paging and any schema beside them.
	CurrencyPage = currencies.CurrencyPage
	// CurrencyCount is the count sub-path's bare {"total": N}.
	CurrencyCount = currencies.CurrencyCount
	// CurrencyBaseChange is the base write's body, the new base and the CurrencyReference that held the flag.
	CurrencyBaseChange = currencies.CurrencyBaseChange
	CurrencyReference  = currencies.CurrencyReference
	// CurrenciesPagination is the currency listing's five-key pagination block.
	CurrenciesPagination = currencies.Pagination
)

// NewCurrencies binds the resource to a requester of your own.
func NewCurrencies(requester Requester) *Currencies { return currencies.New(requester) }

// CurrencyShape is the DTOShape for a body shaped {"currency": {…}}.
func CurrencyShape(decoded any) any { return currencies.CurrencyShape(decoded) }

// CurrencyPageShape is the DTOShape for the listing body.
func CurrencyPageShape(decoded any) any { return currencies.CurrencyPageShape(decoded) }

// CurrencyCountShape is the DTOShape for the count sub-path's body.
func CurrencyCountShape(decoded any) any { return currencies.CurrencyCountShape(decoded) }

// CurrencyBaseChangeShape is the DTOShape for the base write's body.
func CurrencyBaseChangeShape(decoded any) any { return currencies.CurrencyBaseChangeShape(decoded) }

// CurrenciesSortFields returns a fresh copy of the ten names sort takes.
func CurrenciesSortFields() []string { return currencies.CurrenciesSortFields() }

// CurrencyRoundingValues returns a fresh copy of the three values rounding takes.
func CurrencyRoundingValues() []string { return currencies.CurrencyRoundingValues() }

const (
	CurrenciesPath      = currencies.CurrenciesPath
	CurrenciesCountPath = currencies.CurrenciesCountPath
	CurrencyPath        = currencies.CurrencyPath
	CurrencyBasePath    = currencies.CurrencyBasePath

	// CurrenciesKey wraps the rows, CurrencyKey one row, CurrencyPreviousBaseKey the base write's former holder and CurrenciesCountKey the count.
	CurrenciesKey           = currencies.CurrenciesKey
	CurrencyKey             = currencies.CurrencyKey
	CurrencyPreviousBaseKey = currencies.CurrencyPreviousBaseKey
	CurrenciesCountKey      = currencies.CurrenciesCountKey

	// CurrenciesPageSize is the page size when limit is not sent, CurrenciesPageSizeCap the largest limit.
	CurrenciesPageSize    = currencies.CurrenciesPageSize
	CurrenciesPageSizeCap = currencies.CurrenciesPageSizeCap
	// CurrenciesTransactionNoHeader is the header CreateCurrency sends transactionNo under.
	CurrenciesTransactionNoHeader = currencies.CurrenciesTransactionNoHeader
)

// PaymentGateways is the payment gateway endpoints bound to a client.
type PaymentGateways = paymentgateways.PaymentGateways

type (
	// PaymentGateway is one gateway, a listing row, a row path's {"payment_gateway": {…}} or a write's answer.
	PaymentGateway = paymentgateways.PaymentGateway
	// PaymentGatewayReference is an {id, name} pair: a zone, or an audience member.
	PaymentGatewayReference       = paymentgateways.PaymentGatewayReference
	PaymentGatewayVisibility      = paymentgateways.PaymentGatewayVisibility
	PaymentGatewayAllowedFor      = paymentgateways.PaymentGatewayAllowedFor
	PaymentGatewaySurcharge       = paymentgateways.PaymentGatewaySurcharge
	PaymentGatewayMessages        = paymentgateways.PaymentGatewayMessages
	PaymentGatewaySettings        = paymentgateways.PaymentGatewaySettings
	PaymentGatewayPayPalStandard  = paymentgateways.PaymentGatewayPayPalStandard
	PaymentGatewayPayPalClassic   = paymentgateways.PaymentGatewayPayPalClassic
	PaymentGatewayToggle          = paymentgateways.PaymentGatewayToggle
	PaymentGatewayCardProcessor   = paymentgateways.PaymentGatewayCardProcessor
	PaymentGatewayCardLogo        = paymentgateways.PaymentGatewayCardLogo
	PaymentGatewayGooglePay       = paymentgateways.PaymentGatewayGooglePay
	PaymentGatewayApplePay        = paymentgateways.PaymentGatewayApplePay
	PaymentGatewayExpressPayPal   = paymentgateways.PaymentGatewayExpressPayPal
	PaymentGatewayVisibilityInput = paymentgateways.PaymentGatewayVisibilityInput
	// PaymentGatewayInput is the writable form for the PUT and the PATCH, with no envelope.
	PaymentGatewayInput = paymentgateways.PaymentGatewayInput
	// PaymentGatewayPage is one page of the listing: rows on PaymentGateways, beside the paging.
	PaymentGatewayPage = paymentgateways.PaymentGatewayPage
	// PaymentGatewayCount is the count sub-path's bare {"total": N}.
	PaymentGatewayCount = paymentgateways.PaymentGatewayCount
	// PaymentGatewaysPagination is the gateway listing's five-key pagination block.
	PaymentGatewaysPagination = paymentgateways.Pagination
)

// NewPaymentGateways binds the resource to a requester of your own.
func NewPaymentGateways(requester Requester) *PaymentGateways { return paymentgateways.New(requester) }

// PaymentGatewayShape is the DTOShape for a body shaped {"payment_gateway": {…}}.
func PaymentGatewayShape(decoded any) any { return paymentgateways.PaymentGatewayShape(decoded) }

// PaymentGatewayPageShape is the DTOShape for the listing body.
func PaymentGatewayPageShape(decoded any) any {
	return paymentgateways.PaymentGatewayPageShape(decoded)
}

// PaymentGatewayCountShape is the DTOShape for the count sub-path's body.
func PaymentGatewayCountShape(decoded any) any {
	return paymentgateways.PaymentGatewayCountShape(decoded)
}

const (
	PaymentGatewaysPath      = paymentgateways.PaymentGatewaysPath
	PaymentGatewaysCountPath = paymentgateways.PaymentGatewaysCountPath
	PaymentGatewayPath       = paymentgateways.PaymentGatewayPath

	// PaymentGatewaysKey wraps the rows, PaymentGatewayKey one row and PaymentGatewaysCountKey the count.
	PaymentGatewaysKey      = paymentgateways.PaymentGatewaysKey
	PaymentGatewayKey       = paymentgateways.PaymentGatewayKey
	PaymentGatewaysCountKey = paymentgateways.PaymentGatewaysCountKey

	// PaymentGatewaysPageSize is the page size when limit is not sent, PaymentGatewaysPageSizeCap the largest limit accepted.
	PaymentGatewaysPageSize    = paymentgateways.PaymentGatewaysPageSize
	PaymentGatewaysPageSizeCap = paymentgateways.PaymentGatewaysPageSizeCap
)

// ShippingClasses is the shipping class endpoints bound to a client.
type ShippingClasses = shippingclasses.ShippingClasses

type (
	// ShippingClass is one class, a listing row or a row path's {"shipping_class": {…}}: four fields.
	ShippingClass = shippingclasses.ShippingClass
	// ShippingClassInput is the body both writes take, wrapper included: Name and Description.
	ShippingClassInput = shippingclasses.ShippingClassInput
	// ShippingClassPage is one page of the listing, with the schema block field_metadata=true adds.
	ShippingClassPage = shippingclasses.ShippingClassPage
	// ShippingClassCount is the count sub-path's body — the plural key wrapping a record.
	ShippingClassCount = shippingclasses.ShippingClassCount
	// ShippingClassesPagination is the shipping class listing's ten-key pagination block.
	ShippingClassesPagination = shippingclasses.Pagination
)

// NewShippingClasses binds the resource to a requester of your own.
func NewShippingClasses(requester Requester) *ShippingClasses {
	return shippingclasses.New(requester)
}

// ShippingClassShape is the DTOShape for a body shaped {"shipping_class": {…}}, for the raw-verb call.
func ShippingClassShape(decoded any) any { return shippingclasses.ShippingClassShape(decoded) }

// ShippingClassPageShape is the DTOShape for the listing body.
func ShippingClassPageShape(decoded any) any {
	return shippingclasses.ShippingClassPageShape(decoded)
}

// ShippingClassCountShape is the DTOShape for the count sub-path's body.
func ShippingClassCountShape(decoded any) any {
	return shippingclasses.ShippingClassCountShape(decoded)
}

// ShippingClassesQueryParams returns a fresh copy of the four names the listing accepts.
func ShippingClassesQueryParams() []string {
	return shippingclasses.ShippingClassesQueryParams()
}

// ShippingClassesPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func ShippingClassesPaginationKeys() []string {
	return shippingclasses.ShippingClassesPaginationKeys()
}

// ShippingClassesUnroutedMethodBody is the 200 body the store answers for an unrouted method.
func ShippingClassesUnroutedMethodBody() map[string]any {
	return shippingclasses.UnroutedMethodBody()
}

const (
	ShippingClassesPath      = shippingclasses.ShippingClassesPath
	ShippingClassesCountPath = shippingclasses.ShippingClassesCountPath
	ShippingClassPath        = shippingclasses.ShippingClassPath
	ShippingClassDefaultPath = shippingclasses.ShippingClassDefaultPath

	// ShippingClassesKey wraps the rows and the count's record, ShippingClassKey one row and both writes.
	ShippingClassesKey      = shippingclasses.ShippingClassesKey
	ShippingClassKey        = shippingclasses.ShippingClassKey
	ShippingClassesCountKey = shippingclasses.ShippingClassesCountKey

	// ShippingClassesPageSize is the page size when limit is not sent and the ceiling on one that is.
	ShippingClassesPageSize = shippingclasses.ShippingClassesPageSize
)

// ShippingMethods is the shipping method endpoints bound to a client.
type ShippingMethods = shippingmethods.ShippingMethods

type (
	// ShippingMethod is one method, a listing row or a row path's {"shipping_method": {…}}.
	ShippingMethod = shippingmethods.ShippingMethod
	// ShippingMethodInput is the writable form of one, shared by both writes; the store reads only Name.
	ShippingMethodInput = shippingmethods.ShippingMethodInput
	// ShippingMethodPage is one page of the listing, with the schema block field_metadata=true adds.
	ShippingMethodPage = shippingmethods.ShippingMethodPage
	// ShippingMethodCount is the count sub-path's body — the plural key wrapping a record.
	ShippingMethodCount = shippingmethods.ShippingMethodCount
	// ShippingMethodsPagination is the shipping method listing's ten-key pagination block.
	ShippingMethodsPagination = shippingmethods.Pagination
)

// NewShippingMethods binds the resource to a requester of your own.
func NewShippingMethods(requester Requester) *ShippingMethods {
	return shippingmethods.New(requester)
}

// ShippingMethodShape is the DTOShape for a body shaped {"shipping_method": {…}}.
func ShippingMethodShape(decoded any) any { return shippingmethods.ShippingMethodShape(decoded) }

// ShippingMethodPageShape is the DTOShape for the listing body.
func ShippingMethodPageShape(decoded any) any {
	return shippingmethods.ShippingMethodPageShape(decoded)
}

// ShippingMethodCountShape is the DTOShape for the count sub-path's body.
func ShippingMethodCountShape(decoded any) any {
	return shippingmethods.ShippingMethodCountShape(decoded)
}

// ShippingMethodsQueryParams returns a fresh copy of the four names the listing accepts; others are 400.
func ShippingMethodsQueryParams() []string {
	return shippingmethods.ShippingMethodsQueryParams()
}

// ShippingMethodsPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func ShippingMethodsPaginationKeys() []string {
	return shippingmethods.ShippingMethodsPaginationKeys()
}

// ShippingMethodsUnroutedMethodBody is the 200 body the store answers for an unrouted method.
func ShippingMethodsUnroutedMethodBody() map[string]any {
	return shippingmethods.UnroutedMethodBody()
}

const (
	ShippingMethodsPath       = shippingmethods.ShippingMethodsPath
	ShippingMethodsCountPath  = shippingmethods.ShippingMethodsCountPath
	ShippingMethodPath        = shippingmethods.ShippingMethodPath
	ShippingMethodDefaultPath = shippingmethods.ShippingMethodDefaultPath

	// ShippingMethodsKey wraps the rows and the count's record, ShippingMethodKey one row and both writes.
	ShippingMethodsKey      = shippingmethods.ShippingMethodsKey
	ShippingMethodKey       = shippingmethods.ShippingMethodKey
	ShippingMethodsCountKey = shippingmethods.ShippingMethodsCountKey

	// ShippingMethodsPageSize is the page size when limit is not sent, and the ceiling on one that is.
	ShippingMethodsPageSize = shippingmethods.ShippingMethodsPageSize
)

// TaxCodes is the tax code endpoints bound to a client.
type TaxCodes = taxcodes.TaxCodes

type (
	// TaxCode is one tax code, the same eight keys on a listing row and on a row path.
	TaxCode = taxcodes.TaxCode
	// TaxCodeInput is the writable form of one; an update replaces Priority and ResolverType.
	TaxCodeInput = taxcodes.TaxCodeInput
	// TaxCodePage is one page of the listing: rows on TaxCodes, the paging block beside them.
	TaxCodePage = taxcodes.TaxCodePage
	// TaxCodesPagination is the tax code listing's ten-key pagination block.
	TaxCodesPagination = taxcodes.Pagination
)

// NewTaxCodes binds the resource to a requester of your own.
func NewTaxCodes(requester Requester) *TaxCodes { return taxcodes.New(requester) }

// TaxCodeShape is the DTOShape for a body shaped {"tax_code": {…}}, for the raw-verb call.
func TaxCodeShape(decoded any) any { return taxcodes.TaxCodeShape(decoded) }

// TaxCodePageShape is the DTOShape for the listing body.
func TaxCodePageShape(decoded any) any { return taxcodes.TaxCodePageShape(decoded) }

// TaxCodesQueryParams returns a fresh copy of the three names the listing acts on.
func TaxCodesQueryParams() []string { return taxcodes.TaxCodesQueryParams() }

// TaxCodesPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func TaxCodesPaginationKeys() []string { return taxcodes.TaxCodesPaginationKeys() }

// TaxCodesUnroutedMethodBody is the 200 body the store answers for an unrouted method.
func TaxCodesUnroutedMethodBody() map[string]any { return taxcodes.UnroutedMethodBody() }

const (
	TaxCodesPath = taxcodes.TaxCodesPath
	TaxCodePath  = taxcodes.TaxCodePath

	// TaxCodesKey wraps the listing's rows, and TaxCodeKey one row and both writes.
	TaxCodesKey = taxcodes.TaxCodesKey
	TaxCodeKey  = taxcodes.TaxCodeKey

	// TaxCodesPageSize is the default page size; the listing answers up to 100.
	TaxCodesPageSize = taxcodes.TaxCodesPageSize
)

// TaxProfiles is the tax profile endpoints bound to a client.
type TaxProfiles = taxprofiles.TaxProfiles

type (
	// TaxProfile is one tax profile, the same seven keys on a listing row and on a row path.
	TaxProfile = taxprofiles.TaxProfile
	// TaxProfileRule is one tax rule a profile nests, over TaxProfileRuleCode and TaxProfileRuleZone.
	TaxProfileRule     = taxprofiles.TaxProfileRule
	TaxProfileRuleCode = taxprofiles.TaxProfileRuleCode
	TaxProfileRuleZone = taxprofiles.TaxProfileRuleZone
	// TaxProfileInput is the writable form of one; the update merges, and Default is exclusive and one-way.
	TaxProfileInput = taxprofiles.TaxProfileInput
	// TaxProfilePage is one page of the listing: rows on TaxProfiles, the paging block beside them.
	TaxProfilePage = taxprofiles.TaxProfilePage
	// TaxProfilesPagination is the tax profile listing's ten-key pagination block.
	TaxProfilesPagination = taxprofiles.Pagination
)

// NewTaxProfiles binds the resource to a requester of your own.
func NewTaxProfiles(requester Requester) *TaxProfiles { return taxprofiles.New(requester) }

// TaxProfileShape is the DTOShape for a body shaped {"tax_profile": {…}}, for the raw-verb call.
func TaxProfileShape(decoded any) any { return taxprofiles.TaxProfileShape(decoded) }

// TaxProfilePageShape is the DTOShape for the listing body.
func TaxProfilePageShape(decoded any) any { return taxprofiles.TaxProfilePageShape(decoded) }

// TaxProfilesQueryParams returns a fresh copy of the three names the listing acts on.
func TaxProfilesQueryParams() []string { return taxprofiles.TaxProfilesQueryParams() }

// TaxProfilesPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func TaxProfilesPaginationKeys() []string { return taxprofiles.TaxProfilesPaginationKeys() }

// TaxProfilesUnroutedMethodBody is the 200 body the store answers for an unrouted method.
func TaxProfilesUnroutedMethodBody() map[string]any { return taxprofiles.UnroutedMethodBody() }

const (
	TaxProfilesPath = taxprofiles.TaxProfilesPath
	TaxProfilePath  = taxprofiles.TaxProfilePath

	// TaxProfilesKey wraps the listing's rows, and TaxProfileKey one row and both writes.
	TaxProfilesKey = taxprofiles.TaxProfilesKey
	TaxProfileKey  = taxprofiles.TaxProfileKey

	// TaxProfilesPageSize is the default page size and the largest page the listing answers.
	TaxProfilesPageSize = taxprofiles.TaxProfilesPageSize
)

// TaxRules is the tax rule endpoints bound to a client.
type TaxRules = taxrules.TaxRules

type (
	// TaxRule is one tax rule; TaxRuleCode and TaxRuleZone are the read-only projections it nests.
	TaxRule     = taxrules.TaxRule
	TaxRuleCode = taxrules.TaxRuleCode
	TaxRuleZone = taxrules.TaxRuleZone
	// TaxRuleInput is the writable form of one: two of the row's ten fields; the update merges.
	TaxRuleInput = taxrules.TaxRuleInput
	// TaxRulePage is one page of the listing: rows on TaxRules, the paging block beside them.
	TaxRulePage = taxrules.TaxRulePage
	// TaxRulesPagination is the tax rule listing's ten-key pagination block.
	TaxRulesPagination = taxrules.Pagination
)

// NewTaxRules binds the resource to a requester of your own.
func NewTaxRules(requester Requester) *TaxRules { return taxrules.New(requester) }

// TaxRuleShape is the DTOShape for a body shaped {"tax_rule": {…}}, for the raw-verb call.
func TaxRuleShape(decoded any) any { return taxrules.TaxRuleShape(decoded) }

// TaxRulePageShape is the DTOShape for the listing body.
func TaxRulePageShape(decoded any) any { return taxrules.TaxRulePageShape(decoded) }

// TaxRulesQueryParams returns a fresh copy of the three names the listing acts on.
func TaxRulesQueryParams() []string { return taxrules.TaxRulesQueryParams() }

// TaxRulesPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func TaxRulesPaginationKeys() []string { return taxrules.TaxRulesPaginationKeys() }

// TaxRulesUnroutedMethodBody is the 200 body the store answers for an unrouted method.
func TaxRulesUnroutedMethodBody() map[string]any { return taxrules.UnroutedMethodBody() }

const (
	TaxRulesPath = taxrules.TaxRulesPath
	TaxRulePath  = taxrules.TaxRulePath

	// TaxRulesKey wraps the listing's rows, and TaxRuleKey one row and both writes.
	TaxRulesKey = taxrules.TaxRulesKey
	TaxRuleKey  = taxrules.TaxRuleKey

	// TaxRulesPageSize is the default page size; the listing answers up to 100.
	TaxRulesPageSize = taxrules.TaxRulesPageSize
)

// TaxZones is the tax zone endpoints bound to a client.
type TaxZones = taxzones.TaxZones

type (
	// TaxZone is one tax zone; PostCode reads as the flat list a write sends.
	TaxZone = taxzones.TaxZone
	// TaxZoneInput is the writable form of one; both writes require Name, and the update replaces.
	TaxZoneInput = taxzones.TaxZoneInput
	// TaxZonePage is one page of the listing: rows on TaxZones, the paging block beside them.
	TaxZonePage = taxzones.TaxZonePage
	// TaxZonesCountry and TaxZonesState are a place as a zone nests it; a zone write sends the record.
	TaxZonesCountry = taxzones.TaxCountry
	TaxZonesState   = taxzones.TaxState
	// TaxZonesPagination is the tax zone listing's ten-key pagination block.
	TaxZonesPagination = taxzones.Pagination
)

// NewTaxZones binds the resource to a requester of your own.
func NewTaxZones(requester Requester) *TaxZones { return taxzones.New(requester) }

// TaxZoneShape is the DTOShape for a body shaped {"tax_zone": {…}}, for the raw-verb call.
func TaxZoneShape(decoded any) any { return taxzones.TaxZoneShape(decoded) }

// TaxZonePageShape is the DTOShape for the listing body.
func TaxZonePageShape(decoded any) any { return taxzones.TaxZonePageShape(decoded) }

// TaxZonesQueryParams returns a fresh copy of the three names the listing acts on.
func TaxZonesQueryParams() []string { return taxzones.TaxZonesQueryParams() }

// TaxZonesPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func TaxZonesPaginationKeys() []string { return taxzones.TaxZonesPaginationKeys() }

// TaxZonesUnroutedMethodBody is the 200 body the store answers for an unrouted method.
func TaxZonesUnroutedMethodBody() map[string]any { return taxzones.UnroutedMethodBody() }

const (
	TaxZonesPath = taxzones.TaxZonesPath
	TaxZonePath  = taxzones.TaxZonePath

	// TaxZonesKey wraps the listing's rows, and TaxZoneKey one row and both writes.
	TaxZonesKey = taxzones.TaxZonesKey
	TaxZoneKey  = taxzones.TaxZoneKey

	// TaxZonesPageSize is the default page size and the largest page the listing answers.
	TaxZonesPageSize = taxzones.TaxZonesPageSize
)

// TaxSettingsResource is the flat tax settings facade bound to a client; it serves no HEAD.
type TaxSettingsResource = taxsettings.TaxSettingsResource

type (
	// TaxSettings is the facade's nine keys, flat; there are no sections.
	TaxSettings = taxsettings.TaxSettings
	// TaxSettingsInput is the PUT's body, enveloped as settings; the write merges.
	TaxSettingsInput = taxsettings.TaxSettingsInput
)

// NewTaxSettings binds the resource to a requester of your own.
func NewTaxSettings(requester Requester) *TaxSettingsResource {
	return taxsettings.New(requester)
}

// TaxSettingsShape is the DTOShape for a body shaped {"tax_settings": {…}}, for the raw-verb call.
func TaxSettingsShape(decoded any) any { return taxsettings.TaxSettingsShape(decoded) }

const (
	TaxSettingsPath = taxsettings.TaxSettingsPath

	// TaxSettingsKey is the key the read answers under; TaxSettingsWriteKey, a different one, wraps the PUT.
	TaxSettingsKey      = taxsettings.TaxSettingsKey
	TaxSettingsWriteKey = taxsettings.TaxSettingsWriteKey

	// TaxSettingsUnknownPathMessage is the tax router's answer for a path or method it does not route.
	TaxSettingsUnknownPathMessage = taxsettings.TaxSettingsUnknownPathMessage
)

// TaxCountryTaxes is the country tax endpoints bound to a client; it serves no HEAD.
type TaxCountryTaxes = taxcountrytaxes.TaxCountryTaxes

type (
	// CountryTax is one country tax; the row path and the update add StateRules where the listing has StateRulesCount.
	CountryTax          = taxcountrytaxes.CountryTax
	CountryTaxStateRule = taxcountrytaxes.CountryTaxStateRule
	// CountryTaxInput is the writable form of one; the create requires CountryCode, which no read sends, and TransactionNo sits beside the wrapper.
	CountryTaxInput = taxcountrytaxes.CountryTaxInput
	// CountryTaxPage is one page of the listing: rows on CountryTaxes, group_meta and the paging block beside them.
	CountryTaxPage = taxcountrytaxes.CountryTaxPage
	// TaxCountryTaxesCountry is a country tax's nested country, resolved from the create's country_code.
	TaxCountryTaxesCountry = taxcountrytaxes.TaxCountry
	// TaxCountryTaxesState is a state rule's nested state.
	TaxCountryTaxesState = taxcountrytaxes.TaxState
	// TaxCountryTaxesGroupMeta is the listing's group_meta block, each entry a TaxCountryTaxesResourceMeta.
	TaxCountryTaxesGroupMeta    = taxcountrytaxes.GroupMeta
	TaxCountryTaxesResourceMeta = taxcountrytaxes.ResourceMeta
	// TaxCountryTaxesPagination is the country tax listing's ten-key pagination block.
	TaxCountryTaxesPagination = taxcountrytaxes.Pagination
)

// NewTaxCountryTaxes binds the resource to a requester of your own.
func NewTaxCountryTaxes(requester Requester) *TaxCountryTaxes {
	return taxcountrytaxes.New(requester)
}

// CountryTaxShape is the DTOShape for a body shaped {"country_tax": {…}}, for the raw-verb call.
func CountryTaxShape(decoded any) any { return taxcountrytaxes.CountryTaxShape(decoded) }

// CountryTaxPageShape is the DTOShape for the listing body.
func CountryTaxPageShape(decoded any) any { return taxcountrytaxes.CountryTaxPageShape(decoded) }

// TaxCountryTaxesQueryParams returns a fresh copy of the four names the listing acts on.
func TaxCountryTaxesQueryParams() []string { return taxcountrytaxes.TaxCountryTaxesQueryParams() }

// TaxCountryTaxesPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func TaxCountryTaxesPaginationKeys() []string {
	return taxcountrytaxes.TaxCountryTaxesPaginationKeys()
}

// TaxCountryTaxesUnroutedMethodBody is the 200 body the legacy underscore row path returns for an unrouted method.
func TaxCountryTaxesUnroutedMethodBody() map[string]any {
	return taxcountrytaxes.UnroutedMethodBody()
}

const (
	TaxCountryTaxesPath = taxcountrytaxes.TaxCountryTaxesPath
	TaxCountryTaxPath   = taxcountrytaxes.TaxCountryTaxPath

	// TaxCountryTaxesKey wraps the listing's rows, and TaxCountryTaxKey one row and both writes.
	TaxCountryTaxesKey = taxcountrytaxes.TaxCountryTaxesKey
	TaxCountryTaxKey   = taxcountrytaxes.TaxCountryTaxKey

	// TaxCountryTaxesPageSize is the page size when no limit is sent; limit takes 1 to 100.
	TaxCountryTaxesPageSize = taxcountrytaxes.TaxCountryTaxesPageSize

	// TaxCountryTaxesUnknownPathMessage opens the tax router's 404 for a method the listing does not route.
	TaxCountryTaxesUnknownPathMessage = taxcountrytaxes.TaxCountryTaxesUnknownPathMessage

	// TaxCountryTaxesTransactionNoHeader is the header CreateTaxCountryTax sends transactionNo under.
	TaxCountryTaxesTransactionNoHeader = taxcountrytaxes.TaxCountryTaxesTransactionNoHeader
)

// TaxCustomProfiles is the custom tax profile endpoints bound to a client; it serves no HEAD.
type TaxCustomProfiles = taxcustomprofiles.TaxCustomProfiles

type (
	// CustomProfile is one custom tax profile; the row path adds Rules where the listing has RulesCount.
	CustomProfile          = taxcustomprofiles.CustomProfile
	CustomProfileRule      = taxcustomprofiles.CustomProfileRule
	CustomProfileReference = taxcustomprofiles.CustomProfileReference
	// CustomProfileInput is the writable form of one; the create requires Name, and TransactionNo sits beside the wrapper.
	CustomProfileInput = taxcustomprofiles.CustomProfileInput
	// CustomProfilePage is one page of the listing: rows on CustomProfiles, group_meta and the paging block beside them.
	CustomProfilePage = taxcustomprofiles.CustomProfilePage
	// TaxCustomProfilesGroupMeta is the listing's group_meta block, each entry a TaxCustomProfilesResourceMeta.
	TaxCustomProfilesGroupMeta    = taxcustomprofiles.GroupMeta
	TaxCustomProfilesResourceMeta = taxcustomprofiles.ResourceMeta
	// TaxCustomProfilesPagination is the custom tax profile listing's ten-key pagination block.
	TaxCustomProfilesPagination = taxcustomprofiles.Pagination
)

// NewTaxCustomProfiles binds the resource to a requester of your own.
func NewTaxCustomProfiles(requester Requester) *TaxCustomProfiles {
	return taxcustomprofiles.New(requester)
}

// CustomProfileShape is the DTOShape for a body shaped {"custom_profile": {…}}.
func CustomProfileShape(decoded any) any { return taxcustomprofiles.CustomProfileShape(decoded) }

// CustomProfilePageShape is the DTOShape for the listing body.
func CustomProfilePageShape(decoded any) any {
	return taxcustomprofiles.CustomProfilePageShape(decoded)
}

// TaxCustomProfilesQueryParams returns a fresh copy of the four names the listing acts on.
func TaxCustomProfilesQueryParams() []string {
	return taxcustomprofiles.TaxCustomProfilesQueryParams()
}

// TaxCustomProfilesPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func TaxCustomProfilesPaginationKeys() []string {
	return taxcustomprofiles.TaxCustomProfilesPaginationKeys()
}

// TaxCustomProfilesUnroutedMethodBody is the 200 body the legacy underscore row path returns for an unrouted method.
func TaxCustomProfilesUnroutedMethodBody() map[string]any {
	return taxcustomprofiles.UnroutedMethodBody()
}

const (
	TaxCustomProfilesPath = taxcustomprofiles.TaxCustomProfilesPath
	TaxCustomProfilePath  = taxcustomprofiles.TaxCustomProfilePath

	// TaxCustomProfilesKey wraps the listing's rows, and TaxCustomProfileKey one row and both writes.
	TaxCustomProfilesKey = taxcustomprofiles.TaxCustomProfilesKey
	TaxCustomProfileKey  = taxcustomprofiles.TaxCustomProfileKey

	// TaxCustomProfilesPageSize is the page size when no limit is sent; limit takes 1 to 100.
	TaxCustomProfilesPageSize = taxcustomprofiles.TaxCustomProfilesPageSize

	// TaxCustomProfilesUnknownPathMessage opens the tax router's 404 for a method the listing does not route.
	TaxCustomProfilesUnknownPathMessage = taxcustomprofiles.TaxCustomProfilesUnknownPathMessage

	// TaxCustomProfilesTransactionNoHeader is the header CreateTaxCustomProfile sends transactionNo under.
	TaxCustomProfilesTransactionNoHeader = taxcustomprofiles.TaxCustomProfilesTransactionNoHeader
)

// DiscountProfiles is the discount profile endpoints bound to a client; it serves no HEAD.
type DiscountProfiles = discountprofiles.DiscountProfiles

type (
	// DiscountProfile is one profile, a listing row or a row path's {"discount_profile": {…}}.
	DiscountProfile = discountprofiles.DiscountProfile
	// DiscountProfileInput is the writable form of one, with no envelope; the update replaces.
	DiscountProfileInput = discountprofiles.DiscountProfileInput
	// AmountDetails, ShippingDetails and DiscountProfilesProductDetails are the three details blocks.
	AmountDetails   = discountprofiles.AmountDetails
	ShippingDetails = discountprofiles.ShippingDetails
	// DiscountProfilesProductDetails is a profile's product details block.
	DiscountProfilesProductDetails = discountprofiles.ProductDetails
	// CustomerCoupon is one generated code on a profile's customer_coupons.
	CustomerCoupon = discountprofiles.CustomerCoupon
	// DiscountProfilePage is one page of the listing: rows on DiscountProfiles, paging and group_meta beside them.
	DiscountProfilePage = discountprofiles.DiscountProfilePage
	// DiscountProfilesGroupMeta and DiscountProfilesResourceMeta are the listing's group_meta and its one entry.
	DiscountProfilesGroupMeta    = discountprofiles.GroupMeta
	DiscountProfilesResourceMeta = discountprofiles.ResourceMeta
	// DiscountProfileCoupons is the create-coupons answer, and DiscountProfileCouponsInput its body.
	DiscountProfileCoupons      = discountprofiles.DiscountProfileCoupons
	DiscountProfileCouponsInput = discountprofiles.DiscountProfileCouponsInput
	// DiscountCouponCustomerInput is the assign-customer body.
	DiscountCouponCustomerInput = discountprofiles.DiscountCouponCustomerInput
	// DiscountProfilesPagination is the discount profile listing's ten-key pagination block.
	DiscountProfilesPagination = discountprofiles.Pagination
)

// NewDiscountProfiles binds the resource to a requester of your own.
func NewDiscountProfiles(requester Requester) *DiscountProfiles {
	return discountprofiles.New(requester)
}

// DiscountProfileShape is the DTOShape for a body shaped {"discount_profile": {…}}, for the raw-verb call.
func DiscountProfileShape(decoded any) any { return discountprofiles.DiscountProfileShape(decoded) }

// DiscountProfilePageShape is the DTOShape for the listing body.
func DiscountProfilePageShape(decoded any) any {
	return discountprofiles.DiscountProfilePageShape(decoded)
}

// DiscountProfileCouponsShape is the DTOShape for the create-coupons body.
func DiscountProfileCouponsShape(decoded any) any {
	return discountprofiles.DiscountProfileCouponsShape(decoded)
}

// DiscountProfilesQueryParams returns a fresh copy of the eight names the listing reads; others are ignored.
func DiscountProfilesQueryParams() []string {
	return discountprofiles.DiscountProfilesQueryParams()
}

// DiscountProfilesSortFields returns a fresh copy of the three fields sort accepts.
func DiscountProfilesSortFields() []string { return discountprofiles.DiscountProfilesSortFields() }

// DiscountProfileTypes returns a fresh copy of the four values type takes.
func DiscountProfileTypes() []string { return discountprofiles.DiscountProfileTypes() }

// DiscountDetailsKeys returns a fresh copy of the three details block keys, one per discount_details_type.
func DiscountDetailsKeys() []string { return discountprofiles.DiscountDetailsKeys() }

// DiscountProfilesPaginationKeys returns a fresh copy of the ten keys the pagination block carries.
func DiscountProfilesPaginationKeys() []string {
	return discountprofiles.DiscountProfilesPaginationKeys()
}

const (
	DiscountProfilesPath              = discountprofiles.DiscountProfilesPath
	DiscountProfilePath               = discountprofiles.DiscountProfilePath
	DiscountProfileCreateCouponsPath  = discountprofiles.DiscountProfileCreateCouponsPath
	DiscountCouponsAssignCustomerPath = discountprofiles.DiscountCouponsAssignCustomerPath

	// DiscountProfilesKey wraps the listing's rows, DiscountProfileKey one row and DiscountProfileCouponCodesKey the codes.
	DiscountProfilesKey           = discountprofiles.DiscountProfilesKey
	DiscountProfileKey            = discountprofiles.DiscountProfileKey
	DiscountProfileCouponCodesKey = discountprofiles.DiscountProfileCouponCodesKey

	// DiscountProfilesPageSize is the page size when limit is not sent and the ceiling on one that is.
	DiscountProfilesPageSize = discountprofiles.DiscountProfilesPageSize
)

// CustomerGroups is the customer group endpoints bound to a client; it serves no HEAD and no PUT.
type CustomerGroups = customergroups.CustomerGroups

type (
	// CustomerGroupsCustomerGroup is one group, a listing row or {"customer_group": {…}}; the facade's CustomerGroup is the customer portal's.
	CustomerGroupsCustomerGroup = customergroups.CustomerGroup
	// CustomerGroupInput is the writable form, wrapped in customer_group, for the create and the patch.
	CustomerGroupInput = customergroups.CustomerGroupInput
	// CustomerGroupShippingProfile through CustomerGroupProductDisplay are a group's nested blocks; the last four are written as read.
	CustomerGroupShippingProfile    = customergroups.CustomerGroupShippingProfile
	CustomerGroupBackgroundImage    = customergroups.CustomerGroupBackgroundImage
	CustomerGroupAmountSetting      = customergroups.CustomerGroupAmountSetting
	CustomerGroupPaymentTerms       = customergroups.CustomerGroupPaymentTerms
	CustomerGroupOnAccountPurchase  = customergroups.CustomerGroupOnAccountPurchase
	CustomerGroupPurchasePreference = customergroups.CustomerGroupPurchasePreference
	CustomerGroupTaxConfiguration   = customergroups.CustomerGroupTaxConfiguration
	CustomerGroupProductDisplay     = customergroups.CustomerGroupProductDisplay
	// CustomerGroupMember is one member, an entry of Customers or the member read's {"customer": {…}}.
	CustomerGroupMember = customergroups.CustomerGroupMember
	// CustomerGroupPage is one page of the listing: rows on CustomerGroups, paging beside them.
	CustomerGroupPage = customergroups.CustomerGroupPage
	// CustomerGroupCount is the count's bare {"total": N}.
	CustomerGroupCount = customergroups.CustomerGroupCount
	// CustomerGroupAssignment is the assign answer, and CustomerGroupAssignInput its body.
	CustomerGroupAssignment  = customergroups.CustomerGroupAssignment
	CustomerGroupAssignInput = customergroups.CustomerGroupAssignInput
	// CustomerGroupCustomersInput is the member write's body, replacing the whole list.
	CustomerGroupCustomersInput = customergroups.CustomerGroupCustomersInput
	// CustomerGroupsPagination is the customer group listing's five-key pagination block.
	CustomerGroupsPagination = customergroups.Pagination
)

// NewCustomerGroups binds the resource to a requester of your own.
func NewCustomerGroups(requester Requester) *CustomerGroups {
	return customergroups.New(requester)
}

// CustomerGroupShape is the DTOShape for a body shaped {"customer_group": {…}}.
func CustomerGroupShape(decoded any) any { return customergroups.CustomerGroupShape(decoded) }

// CustomerGroupPageShape is the DTOShape for the listing body.
func CustomerGroupPageShape(decoded any) any { return customergroups.CustomerGroupPageShape(decoded) }

// CustomerGroupCountShape is the DTOShape for the count body.
func CustomerGroupCountShape(decoded any) any { return customergroups.CustomerGroupCountShape(decoded) }

// CustomerGroupMemberShape is the DTOShape for the member read.
func CustomerGroupMemberShape(decoded any) any {
	return customergroups.CustomerGroupMemberShape(decoded)
}

// CustomerGroupAssignmentShape is the DTOShape for the assign answer.
func CustomerGroupAssignmentShape(decoded any) any {
	return customergroups.CustomerGroupAssignmentShape(decoded)
}

// CustomerGroupsSortFields returns a fresh copy of the six names sort accepts.
func CustomerGroupsSortFields() []string { return customergroups.CustomerGroupsSortFields() }

const (
	CustomerGroupsPath         = customergroups.CustomerGroupsPath
	CustomerGroupsCountPath    = customergroups.CustomerGroupsCountPath
	CustomerGroupsAssignPath   = customergroups.CustomerGroupsAssignPath
	CustomerGroupPath          = customergroups.CustomerGroupPath
	CustomerGroupCustomersPath = customergroups.CustomerGroupCustomersPath
	CustomerGroupCustomerPath  = customergroups.CustomerGroupCustomerPath

	CustomerGroupsKey         = customergroups.CustomerGroupsKey
	CustomerGroupKey          = customergroups.CustomerGroupKey
	CustomerGroupCustomersKey = customergroups.CustomerGroupCustomersKey
	CustomerGroupCustomerKey  = customergroups.CustomerGroupCustomerKey
	CustomerGroupsCountKey    = customergroups.CustomerGroupsCountKey

	// CustomerGroupsPageSize is the page size when limit is not sent, and CustomerGroupsPageSizeCap the largest limit.
	CustomerGroupsPageSize    = customergroups.CustomerGroupsPageSize
	CustomerGroupsPageSizeCap = customergroups.CustomerGroupsPageSizeCap
)

// VariationsResource is the variation type and product variant endpoints bound to a client.
type VariationsResource = variations.VariationsResource

type (
	// Variation is one type, a listing row or a row path's {"variation": {…}}; VariationOption one option and VariationOptionImage an image option's picture.
	Variation            = variations.Variation
	VariationOption      = variations.VariationOption
	VariationOptionImage = variations.VariationOptionImage
	// VariationPage is one page of the listing: rows on Variations, paging and any schema beside them.
	VariationPage = variations.VariationPage
	// VariationsPagination is the variation listing's five-key pagination block.
	VariationsPagination = variations.Pagination
	// VariationInput is the type writes' body, wrapped in "variation"; VariationOptionInput and VariationOptionImageInput are its parts.
	VariationInput            = variations.VariationInput
	VariationOptionInput      = variations.VariationOptionInput
	VariationOptionImageInput = variations.VariationOptionImageInput
	// ProductVariations is the product read's body, ProductVariationGroup one group and ProductVariant one variant.
	ProductVariations     = variations.ProductVariations
	ProductVariationGroup = variations.ProductVariationGroup
	ProductVariant        = variations.ProductVariant
	// ProductVariantPricing through ProductVariantBadges are a variant's blocks.
	ProductVariantPricing     = variations.ProductVariantPricing
	ProductVariantExpectToPay = variations.ProductVariantExpectToPay
	ProductVariantInventory   = variations.ProductVariantInventory
	ProductVariantOrdering    = variations.ProductVariantOrdering
	ProductVariantDimensions  = variations.ProductVariantDimensions
	ProductVariantBadges      = variations.ProductVariantBadges
	// ProductVariationsBulkResult is the bulk write's answer, ProductVariantUpdate one written variant and ProductVariantOptionReference one of its options.
	ProductVariationsBulkResult   = variations.ProductVariationsBulkResult
	ProductVariantUpdate          = variations.ProductVariantUpdate
	ProductVariantOptionReference = variations.ProductVariantOptionReference
	// ProductVariationsBulkInput is the bulk write's body, a bare record.
	ProductVariationsBulkInput = variations.ProductVariationsBulkInput
	// ProductVariationsImportInput is the import's bare body, ProductVariationGroupInput one group and ProductVariantInput one variant.
	ProductVariationsImportInput = variations.ProductVariationsImportInput
	ProductVariationGroupInput   = variations.ProductVariationGroupInput
	ProductVariantInput          = variations.ProductVariantInput
	// ProductVariationsImportResult is the import's answer.
	ProductVariationsImportResult = variations.ProductVariationsImportResult
)

// NewVariations binds the resource to a requester of your own.
func NewVariations(requester Requester) *VariationsResource {
	return variations.New(requester)
}

// VariationShape is the DTOShape for a body shaped {"variation": {…}}.
func VariationShape(decoded any) any { return variations.VariationShape(decoded) }

// VariationPageShape is the DTOShape for the listing body.
func VariationPageShape(decoded any) any { return variations.VariationPageShape(decoded) }

// ProductVariationsShape is the DTOShape for the product read's body.
func ProductVariationsShape(decoded any) any { return variations.ProductVariationsShape(decoded) }

// ProductVariationsBulkResultShape is the DTOShape for the bulk write's answer.
func ProductVariationsBulkResultShape(decoded any) any {
	return variations.ProductVariationsBulkResultShape(decoded)
}

// ProductVariationsImportResultShape is the DTOShape for the import's answer.
func ProductVariationsImportResultShape(decoded any) any {
	return variations.ProductVariationsImportResultShape(decoded)
}

// VariationsQueryParams returns a fresh copy of the six names the listing accepts; any other is 400.
func VariationsQueryParams() []string { return variations.VariationsQueryParams() }

// VariationStandards returns a fresh copy of the three standards a type takes.
func VariationStandards() []string { return variations.VariationStandards() }

// ProductVariationsBulkFields returns a fresh copy of the six fields the bulk write sets.
func ProductVariationsBulkFields() []string { return variations.ProductVariationsBulkFields() }

// ProductVariationsResponseModes returns a fresh copy of the bulk write's two response values.
func ProductVariationsResponseModes() []string { return variations.ProductVariationsResponseModes() }

const (
	VariationsPath            = variations.VariationsPath
	VariationPath             = variations.VariationPath
	VariationOptionPath       = variations.VariationOptionPath
	ProductVariationsPath     = variations.ProductVariationsPath
	ProductVariationsBulkPath = variations.ProductVariationsBulkPath

	// VariationsKey wraps the listing's rows and the product read's groups, VariationKey one type and every type write.
	VariationsKey = variations.VariationsKey
	VariationKey  = variations.VariationKey
	// ProductVariantsKey wraps a group's variants; the bulk answer lists its variants under the other two.
	ProductVariantsKey          = variations.ProductVariantsKey
	ProductVariationsUpdatedKey = variations.ProductVariationsUpdatedKey
	ProductVariationsFailedKey  = variations.ProductVariationsFailedKey

	// VariationsPageSize is the page size when limit is not sent, and VariationsPageSizeCap the ceiling on one that is.
	VariationsPageSize    = variations.VariationsPageSize
	VariationsPageSizeCap = variations.VariationsPageSizeCap

	// VariationsTransactionNoHeader is the header ImportProductVariations sends transactionNo under.
	VariationsTransactionNoHeader = variations.VariationsTransactionNoHeader
)

// Transactions is the transaction-number endpoints bound to a client; it reads only.
type Transactions = transactions.Transactions

type (
	// TransactionListEntry is one listing row: a number, its record count and its read path.
	TransactionListEntry = transactions.TransactionListEntry
	// TransactionTypeCount is one by_type row of a number's read.
	TransactionTypeCount = transactions.TransactionTypeCount
	// Transaction is one number's bare read: its record count and its 15 by_type rows.
	Transaction = transactions.Transaction
	// TransactionItem is one record under a number, the detail=slim shape.
	TransactionItem = transactions.TransactionItem
	// TransactionPage is one page of the number listing: rows on Transactions, paging beside them.
	TransactionPage = transactions.TransactionPage
	// TransactionItemPage is one page of a number's records, beside the number.
	TransactionItemPage = transactions.TransactionItemPage
	// TransactionsPagination is the two listings' five-key pagination block.
	TransactionsPagination = transactions.Pagination
)

// NewTransactions binds the resource to a requester of your own.
func NewTransactions(requester Requester) *Transactions { return transactions.New(requester) }

// TransactionShape is the DTOShape for a number's bare read body.
func TransactionShape(decoded any) any { return transactions.TransactionShape(decoded) }

// TransactionPageShape is the DTOShape for the number listing body.
func TransactionPageShape(decoded any) any { return transactions.TransactionPageShape(decoded) }

// TransactionItemPageShape is the DTOShape for the items listing body.
func TransactionItemPageShape(decoded any) any { return transactions.TransactionItemPageShape(decoded) }

const (
	TransactionsPath     = transactions.TransactionsPath
	TransactionPath      = transactions.TransactionPath
	TransactionItemsPath = transactions.TransactionItemsPath

	// TransactionsKey wraps the number listing's rows and TransactionItemsKey a number's records of one type.
	TransactionsKey     = transactions.TransactionsKey
	TransactionItemsKey = transactions.TransactionItemsKey

	// TransactionsPageSize is the page size when limit is not sent, TransactionsMaxLimit the largest one either listing applies.
	TransactionsPageSize = transactions.TransactionsPageSize
	TransactionsMaxLimit = transactions.TransactionsMaxLimit
)

// Trash is the trash endpoints bound to a client.
type Trash = trash.Trash

type (
	// TrashItem is one trashed row and a restore's answer; TrashPage one page of the listing, rows on Items.
	TrashItem = trash.TrashItem
	TrashPage = trash.TrashPage
	// TrashPagination is the trash listing's five-key pagination block.
	TrashPagination = trash.Pagination
	// TrashTypeList is the type list's body over TrashType.
	TrashTypeList = trash.TrashTypeList
	TrashType     = trash.TrashType
	// TrashBulkRestoreResult and TrashBulkDeleteResult are the bulk calls' rows, each failed item a TrashBulkError.
	TrashBulkRestoreResult = trash.TrashBulkRestoreResult
	TrashBulkDeleteResult  = trash.TrashBulkDeleteResult
	TrashBulkError         = trash.TrashBulkError
	// TrashBulkInput is the body both bulk calls take, with no envelope, over TrashBulkItemInput.
	TrashBulkInput     = trash.TrashBulkInput
	TrashBulkItemInput = trash.TrashBulkItemInput
)

// NewTrash binds the resource to a requester of your own.
func NewTrash(requester Requester) *Trash { return trash.New(requester) }

// TrashItemShape is the DTOShape for a restore's answer.
func TrashItemShape(decoded any) any { return trash.TrashItemShape(decoded) }

// TrashPageShape is the DTOShape for the listing body.
func TrashPageShape(decoded any) any { return trash.TrashPageShape(decoded) }

// TrashTypeListShape is the DTOShape for the type list's body.
func TrashTypeListShape(decoded any) any { return trash.TrashTypeListShape(decoded) }

// TrashBulkRestoreResultShape is the DTOShape for the bulk restore's body.
func TrashBulkRestoreResultShape(decoded any) any { return trash.TrashBulkRestoreResultShape(decoded) }

// TrashBulkDeleteResultShape is the DTOShape for the bulk delete's body.
func TrashBulkDeleteResultShape(decoded any) any { return trash.TrashBulkDeleteResultShape(decoded) }

const (
	TrashPath            = trash.TrashPath
	TrashTypesPath       = trash.TrashTypesPath
	TrashItemPath        = trash.TrashItemPath
	TrashItemRestorePath = trash.TrashItemRestorePath
	TrashBulkRestorePath = trash.TrashBulkRestorePath
	TrashBulkDeletePath  = trash.TrashBulkDeletePath

	// TrashItemsKey wraps the listing's rows and a bulk body's items, TrashTypesKey the type list.
	TrashItemsKey = trash.TrashItemsKey
	TrashTypesKey = trash.TrashTypesKey

	// TrashPageSize is the page size when limit is not sent, TrashPageSizeCap the largest limit accepted.
	TrashPageSize    = trash.TrashPageSize
	TrashPageSizeCap = trash.TrashPageSizeCap
)

// Email is the email endpoint bound to a client.
type Email = email.Email

type (
	// EmailInput is the send body, bare and camelCase, with no envelope.
	EmailInput = email.EmailInput
	// EmailSendResult is the send's verdict: status and the accepted, rejected, cc and bcc addresses.
	EmailSendResult = email.EmailSendResult
)

// NewEmail binds the resource to a requester of your own.
func NewEmail(requester Requester) *Email { return email.New(requester) }

// EmailSendResultShape is the DTOShape for the send's body.
func EmailSendResultShape(decoded any) any { return email.EmailSendResultShape(decoded) }

const (
	EmailSendPath = email.EmailSendPath

	// EmailPluginUUIDHeader is the header SendEmail sends its pluginUUID under.
	EmailPluginUUIDHeader = email.EmailPluginUUIDHeader
)

// Pages is the page endpoints bound to a client.
type Pages = pages.Pages

type (
	// PagesPage is one page, a listing row, a row path's {"page": {…}} or a write's answer; Page is the generic listing page.
	PagesPage = pages.Page
	// PageReference is an {id, name} pair: a layout, an access-list member, a navigation or the author.
	PageReference          = pages.PageReference
	PagePasswordProtection = pages.PagePasswordProtection
	PageMenuPlacement      = pages.PageMenuPlacement
	PagePublishing         = pages.PagePublishing
	// PageInput is the writable form for the create, the replace and the merge, wrapped in page.
	PageInput = pages.PageInput
	// PageURLInput is the URL write's bare {"url": …}.
	PageURLInput = pages.PageURLInput
	// PageListing is one page of the listing: rows on Pages, beside the paging.
	PageListing = pages.PageListing
	// PageCount is the count sub-path's bare {"total": N}.
	PageCount = pages.PageCount
	// PageDeleteBlockers is the delete check's body, each blocker a PageDeleteBlocker.
	PageDeleteBlockers = pages.PageDeleteBlockers
	PageDeleteBlocker  = pages.PageDeleteBlocker
	// PagesPagination is the page listing's five-key pagination block.
	PagesPagination = pages.Pagination
)

// NewPages binds the resource to a requester of your own.
func NewPages(requester Requester) *Pages { return pages.New(requester) }

// PageShape is the DTOShape for a body shaped {"page": {…}}.
func PageShape(decoded any) any { return pages.PageShape(decoded) }

// PageListingShape is the DTOShape for the listing body.
func PageListingShape(decoded any) any { return pages.PageListingShape(decoded) }

// PageCountShape is the DTOShape for the count sub-path's body.
func PageCountShape(decoded any) any { return pages.PageCountShape(decoded) }

// PageDeleteBlockersShape is the DTOShape for the delete check's body.
func PageDeleteBlockersShape(decoded any) any { return pages.PageDeleteBlockersShape(decoded) }

const (
	PagesPath              = pages.PagesPath
	PagesCountPath         = pages.PagesCountPath
	PagePath               = pages.PagePath
	PageURLPath            = pages.PageURLPath
	PageLandingPath        = pages.PageLandingPath
	PageDeleteBlockersPath = pages.PageDeleteBlockersPath
	PageDuplicatePath      = pages.PageDuplicatePath

	// PagesKey wraps the listing's rows, PageKey one row and the writes, PagesCountKey the count.
	PagesKey      = pages.PagesKey
	PageKey       = pages.PageKey
	PagesCountKey = pages.PagesCountKey

	// PagesPageSize is the page size when limit is not sent, PagesPageSizeCap the largest number limit takes.
	PagesPageSize    = pages.PagesPageSize
	PagesPageSizeCap = pages.PagesPageSizeCap

	// PagesIfMatchHeader is the header UpdatePage, PatchPage and DeletePage send ifMatch as.
	PagesIfMatchHeader = pages.PagesIfMatchHeader
)
