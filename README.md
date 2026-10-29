# WebCommander Go SDK

A Go client for the WebCommander v4 REST API, with OAuth2 token management, retries, secret
redaction and no dependencies outside the standard library.

## Requirements

- Go 1.24 or later.
- API credentials for a WebCommander instance: base URL, client ID, client secret, redirect URI
  and auth string.

## Installation

```sh
go get github.com/WebCommanderCMS/webcommander-sdk-go/v4
```

Keep the `/v4` suffix: Go requires it on the module path of a major version 2 or later.

Import the root package, which re-exports the public surface:

```go
import wc "github.com/WebCommanderCMS/webcommander-sdk-go/v4"
```

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	wc "github.com/WebCommanderCMS/webcommander-sdk-go/v4"
)

func main() {
	client, err := wc.NewClient(wc.SdkConfig{
		Auth: wc.AuthConfig{
			ClientName:   "default",
			BaseURL:      os.Getenv("WC_BASE_URL"),
			ClientID:     os.Getenv("WC_CLIENT_ID"),
			ClientSecret: os.Getenv("WC_CLIENT_SECRET"),
			RedirectURI:  os.Getenv("WC_REDIRECT_URI"),
			AuthString:   os.Getenv("WC_AUTH_STRING"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	response, err := client.AdminCustomers().GetCustomers(context.Background(), map[string]any{"limit": 20}, nil)
	if err != nil {
		log.Fatal(err)
	}
	if !response.Succeeded() {
		log.Fatal(response.StatusCode, response.GetErrorMessages())
	}
	for _, row := range response.Map().(map[string]any)["customers"].([]any) {
		customer := row.(map[string]any)
		fmt.Println(customer["customer_id"], customer["email"])
	}
}
```

- The first call fetches a token, caches it and renews it before it expires.
- A failed call returns a response with `Succeeded()` false and a `nil` error unless you set
  `RaiseOnError`.
- `Close` releases the connection pool and any log file, and every later call on the client is
  refused. Build one client and share it across goroutines rather than one per request.
- The SDK never reads credentials from the environment or from `.env`. The example reads them in
  your code.

## Configuration

`AuthConfig` holds the credentials; the rest of `SdkConfig` is behaviour, and every behaviour field
has a default. The settings you are most likely to change:

| Setting | Default | Meaning |
|---|---|---|
| `RaiseOnError` | `false` | Also return an `*SdkError` when a call fails. |
| `ResponseFormat` | `"map"` | What `response.Data` holds: `"map"`, `"json"` or `"dto"`. |
| `TimeoutSeconds` | `30` | Time limit for one attempt. |
| `MaxRetries` | `1` | Total attempts, so `1` means no retry. |
| `LogLevel` | `"WARNING"` | Lowest level passed to the log hook and `LogFile`. |
| `PrintApiInfo` | `false` | Print method, URL, status and timing for each call. |

```go
client, err := wc.NewClient(wc.SdkConfig{Auth: auth, MaxRetries: 3, RaiseOnError: true})
```

To register several instances once and pick one by name:

```go
wc.RegisterInstances(productionAuth, stagingAuth)
client, err := wc.BuildClient("staging")
```

## Usage

### Bound and explicit calls

Each resource is reached through an accessor on the client, or bound to any `wc.Requester` with
its `wc.New…` constructor. Both give the same type:

```go
client.AdminCustomers().GetCustomer(ctx, 1192, nil)      // bound to the client
wc.NewAdminCustomers(client).GetCustomer(ctx, 1192, nil) // explicit requester
```

Every method takes a `context.Context` first, then its path values and payload, then a
`raiseOnError *bool` (`nil` defers to the config) and an optional response format. It returns
`(*wc.ApiResponse, error)`.

For an endpoint the SDK has no method for, use `client.SendGetRequest`, `SendPostRequest`,
`SendPutRequest`, `SendPatchRequest`, `SendDeleteRequest` or `SendApiRequest` with a path template
and `*wc.RequestOptions`. The path templates are exported (`wc.AdminCustomersPath` and the rest).

### Payload forms

A body is a map, a raw JSON string or a payload DTO. A map and a DTO send the same bytes; a string
is sent exactly as written:

```go
client.AdminCustomers().CreateCustomer(ctx, map[string]any{
	"customer": map[string]any{"first_name": "Ada", "email": "ada@example.com"},
}, "", nil)
client.AdminCustomers().CreateCustomer(ctx, `{"customer":{"first_name":"Ada","email":"ada@example.com"}}`, "", nil)
client.AdminCustomers().CreateCustomer(ctx, wc.CustomerInput{FirstName: "Ada", Email: "ada@example.com"}, "", nil)
```

The `""` argument is the transaction number, sent as `X-Transaction-No` when not empty. Query
parameter and body field names go to the API as you write them.

### Response forms

The body is available in four forms:

```go
response, err := client.AdminCustomers().GetCustomer(ctx, 1192, nil)

response.Data   // the form selected by ResponseFormat; the decoded body by default
response.Map()  // always the decoded body: map[string]any, []any, string or []byte
response.JSON() // always the body text as received
response.DTO()  // always the typed view, or nil where the endpoint declares no type
```

Choose the form per call, per client, or with the `WC_RESPONSE_FORMAT` environment variable:

```go
response, err := client.AdminCustomers().GetCustomers(ctx, nil, nil, wc.ResponseFormatDTO)
page := response.Data.(*wc.Page[wc.Customer])
```

JSON numbers decode as `json.Number`, so integer ids never become floats.

### Pagination

`wc.IterPages` and `wc.IterRecords` follow `pagination.next_page` for you:

```go
options := wc.PaginateOptions{QueryParams: map[string]any{"limit": 20}}
for customer, err := range wc.IterRecords(ctx, client, wc.AdminCustomersPath, "customers", options) {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(customer.(map[string]any)["customer_id"])
}
```

Most listings cap `limit` at 20.

## Error handling

By default a failed call returns its response and a `nil` error. Check `Succeeded()`, then read
the status and messages:

```go
response, _ := client.AdminCustomers().GetCustomer(ctx, 1192, nil)
if !response.Succeeded() {
	fmt.Println(response.StatusCode, response.ErrorType, response.GetErrorMessages())
}
```

With `RaiseOnError` (in the config, or `wc.Ptr(true)` per call) a failed call also returns an
`*SdkError`, the one error type the SDK returns. Classify it with `errors.Is` and `errors.As`:

```go
response, err := client.AdminCustomers().GetCustomer(ctx, 1192, wc.Ptr(true))
var sdkError *wc.SdkError
if errors.As(err, &sdkError) {
	fmt.Println(sdkError.StatusCode, sdkError.ErrorType, sdkError.GetErrorMessages())
}
if errors.Is(err, wc.ErrAuth) {
	// credentials rejected
}
```

A `200` whose body carries `{"isSuccess": false}` counts as a failure. Configuration mistakes
(`wc.ErrConfig`) and caller mistakes caught before any I/O (`wc.ErrUsage`) are always returned as
errors.

## Resources

The client exposes 46 resources with 349 operations. The [API reference](https://newlify.mintlify.app)
lists the method, path and behaviour of each one.

| Resource | Client accessor | Operations |
|---|---|---:|
| Customers | `client.AdminCustomers()` | 15 |
| Aliases | `client.Aliases()` | 4 |
| Animation effects settings | `client.AnimationEffectsSettings()` | 5 |
| Blog categories | `client.BlogCategories()` | 5 |
| Blogs | `client.Blogs()` | 12 |
| Carts | `client.Carts()` | 7 |
| Categories | `client.Categories()` | 10 |
| Category page settings | `client.CategoryPageSettings()` | 5 |
| Checkout page settings | `client.CheckoutPageSettings()` | 5 |
| Cities | `client.Cities()` | 2 |
| Common settings | `client.CommonSettings()` | 6 |
| Countries | `client.Countries()` | 8 |
| Currencies | `client.Currencies()` | 11 |
| Customer groups | `client.CustomerGroups()` | 9 |
| Customer portal settings | `client.CustomerPortalSettings()` | 5 |
| Discount profiles | `client.DiscountProfiles()` | 7 |
| Ecommerce settings | `client.EcommerceSettings()` | 5 |
| Email | `client.Email()` | 1 |
| Inventory locations | `client.InventoryLocations()` | 11 |
| Order custom fields | `client.OrderCustomFields()` | 4 |
| Order flags | `client.OrderFlags()` | 7 |
| Orders | `client.Orders()` | 25 |
| Pages | `client.Pages()` | 15 |
| Payment gateways | `client.PaymentGateways()` | 8 |
| Price display settings | `client.PriceDisplaySettings()` | 5 |
| Product listing settings | `client.ProductListingSettings()` | 3 |
| Product page settings | `client.ProductPageSettings()` | 5 |
| Products | `client.Products()` | 10 |
| Search result page settings | `client.SearchResultPageSettings()` | 13 |
| Search settings | `client.SearchSettings()` | 9 |
| SEO configuration | `client.SeoConfigurationSettings()` | 18 |
| Shipping classes | `client.ShippingClasses()` | 9 |
| Shipping methods | `client.ShippingMethods()` | 9 |
| Shipping settings | `client.ShippingSettings()` | 3 |
| Shopping cart settings | `client.ShoppingCartSettings()` | 5 |
| Site settings | `client.SiteSettings()` | 5 |
| Tax codes | `client.TaxCodes()` | 7 |
| Tax country taxes | `client.TaxCountryTaxes()` | 5 |
| Tax custom profiles | `client.TaxCustomProfiles()` | 5 |
| Tax profiles | `client.TaxProfiles()` | 7 |
| Tax rules | `client.TaxRules()` | 7 |
| Tax settings | `client.TaxSettings()` | 2 |
| Tax zones | `client.TaxZones()` | 7 |
| Transactions | `client.Transactions()` | 6 |
| Trash | `client.Trash()` | 8 |
| Variations | `client.Variations()` | 12 |

Each resource lives in its own package under `resources/` with its typed views and path
constants; the root package re-exports them. `client.ListResourceNames()` returns the registry
keys (`"admin_customers"`, `"orders"` and so on).

## Documentation

- API reference: <https://newlify.mintlify.app>

Every exported name has a doc comment; run `go doc github.com/WebCommanderCMS/webcommander-sdk-go/v4`
from your project to browse them.

## License

MIT. See [LICENSE](LICENSE).
