package countries

import (
	"context"
	"fmt"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// Countries is the country endpoints, bound to one client.
type Countries struct {
	requester http.Requester
}

// ListCountries lists the countries — GET /admin/countries.
func (c *Countries) ListCountries(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, CountriesPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CountryPageShape})
}

// HeadCountries asks whether the collection is reachable — HEAD /admin/countries; no body, and the ETag the listing's first page carries.
func (c *Countries) HeadCountries(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendApiRequest(ctx, "HEAD", CountriesPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// GetCountry retrieves one country — GET /admin/countries/{country_id}; the slot takes the numeric id, the ISO code or the name, in any case.
func (c *Countries) GetCountry(ctx context.Context, countryID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, CountryPath, &http.RequestOptions{
		PathParams:   map[string]any{"country_id": countryID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CountryShape})
}

// HeadCountry asks whether one country is there — HEAD /admin/countries/{country_id}; no body, and the ETag its read carries.
func (c *Countries) HeadCountry(ctx context.Context, countryID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendApiRequest(ctx, "HEAD", CountryPath, &http.RequestOptions{
		PathParams:   map[string]any{"country_id": countryID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ListCountryStates lists one country's states beside the short parent country — GET /admin/countries/{country_id}/states; the same query names as ListCountries.
func (c *Countries) ListCountryStates(ctx context.Context, countryID any, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, CountryStatesPath, &http.RequestOptions{
		PathParams:  map[string]any{"country_id": countryID},
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: CountryStatePageShape})
}

// HeadCountryStates asks whether one country's states are reachable — HEAD /admin/countries/{country_id}/states; no body, and the ETag the state listing's first page carries.
func (c *Countries) HeadCountryStates(ctx context.Context, countryID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendApiRequest(ctx, "HEAD", CountryStatesPath, &http.RequestOptions{
		PathParams:   map[string]any{"country_id": countryID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// GetCountryState retrieves one state of one country — GET /admin/countries/{country_id}/states/{state_id}; stateID takes the state's id, code or name, in any case.
func (c *Countries) GetCountryState(ctx context.Context, countryID any, stateID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, CountryStatePath, &http.RequestOptions{
		PathParams:   map[string]any{"country_id": countryID, "state_id": stateID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: CountryStateShape})
}

// HeadCountryState asks whether one state of one country is there — HEAD /admin/countries/{country_id}/states/{state_id}; no body, and the ETag its read carries.
func (c *Countries) HeadCountryState(ctx context.Context, countryID any, stateID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendApiRequest(ctx, "HEAD", CountryStatePath, &http.RequestOptions{
		PathParams:   map[string]any{"country_id": countryID, "state_id": stateID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "countries"

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

// New binds the country endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *Countries {
	return &Countries{requester: requester}
}
