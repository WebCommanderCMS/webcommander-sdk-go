package variations

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// VariationsResource is the variation type and product variant endpoints, bound to one client.
type VariationsResource struct {
	requester http.Requester
}

// ListVariations lists the types by ascending id — GET /admin/variations; search, standard, limit (1 to 100), offset or page, and field_metadata, any other name is 400.
func (v *VariationsResource) ListVariations(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return v.requester.SendGetRequest(ctx, VariationsPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: VariationPageShape})
}

// HeadVariations asks whether the collection is reachable — HEAD /admin/variations; no body.
func (v *VariationsResource) HeadVariations(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return v.requester.SendApiRequest(ctx, "HEAD", VariationsPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CreateVariation creates a type — POST /admin/variations, 201; {"variation": {...}} with name, standard and at least one option.
func (v *VariationsResource) CreateVariation(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return v.requester.SendPostRequest(ctx, VariationsPath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: VariationShape})
}

// GetVariation retrieves one type — GET /admin/variations/{variation_id}; field_metadata adds the schema beside it.
func (v *VariationsResource) GetVariation(ctx context.Context, variationID any, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return v.requester.SendGetRequest(ctx, VariationPath, &http.RequestOptions{
		PathParams:  map[string]any{"variation_id": variationID},
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: VariationShape})
}

// HeadVariation asks whether one type is there — HEAD /admin/variations/{variation_id}; no body, the row read's ETag.
func (v *VariationsResource) HeadVariation(ctx context.Context, variationID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return v.requester.SendApiRequest(ctx, "HEAD", VariationPath, &http.RequestOptions{
		PathParams:   map[string]any{"variation_id": variationID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// UpdateVariation replaces one type — PUT /admin/variations/{variation_id}; name and options are required, and an option left out is deleted.
func (v *VariationsResource) UpdateVariation(ctx context.Context, variationID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return v.requester.SendPutRequest(ctx, VariationPath, &http.RequestOptions{
		PathParams: map[string]any{"variation_id": variationID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: VariationShape})
}

// PatchVariation merges into one type — PATCH /admin/variations/{variation_id}; the options sent are added, or changed where an id names one, and the rest kept.
func (v *VariationsResource) PatchVariation(ctx context.Context, variationID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return v.requester.SendPatchRequest(ctx, VariationPath, &http.RequestOptions{
		PathParams: map[string]any{"variation_id": variationID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: VariationShape})
}

// DeleteVariation deletes one type — DELETE /admin/variations/{variation_id}, 204; a type with options is 409 unless cascade is true, which deletes them with it.
func (v *VariationsResource) DeleteVariation(ctx context.Context, variationID any, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return v.requester.SendDeleteRequest(ctx, VariationPath, &http.RequestOptions{
		PathParams:  map[string]any{"variation_id": variationID},
		QueryParams: queryParams, RaiseOnError: raiseOnError, ResponseFormat: format})
}

// DeleteVariationOption deletes one option of a type — DELETE /admin/variations/{variation_id}/options/{variation_option_id}, 204; a type's last option is 409.
func (v *VariationsResource) DeleteVariationOption(ctx context.Context, variationID any, variationOptionID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return v.requester.SendDeleteRequest(ctx, VariationOptionPath, &http.RequestOptions{
		PathParams:   map[string]any{"variation_id": variationID, "variation_option_id": variationOptionID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ListProductVariations reads one product's variants, grouped and unpaged — GET /admin/products/{product_id}/variations; a tax-inclusive product's prices read net of tax.
func (v *VariationsResource) ListProductVariations(ctx context.Context, productID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return v.requester.SendGetRequest(ctx, ProductVariationsPath, &http.RequestOptions{
		PathParams:   map[string]any{"product_id": productID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ProductVariationsShape})
}

// ImportProductVariations adds or changes a product's variants by key — POST /admin/products/{product_id}/variations, answered 201; a product_key naming other types than the product's is 409.
func (v *VariationsResource) ImportProductVariations(ctx context.Context, productID any, payload any, transactionNo string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	var headers map[string]string
	if transactionNo != "" {
		headers = map[string]string{VariationsTransactionNoHeader: transactionNo}
	}
	return v.requester.SendPostRequest(ctx, ProductVariationsPath, &http.RequestOptions{
		PathParams: map[string]any{"product_id": productID}, JSONBody: payload,
		Headers: headers, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ProductVariationsImportResultShape})
}

// UpdateProductVariations sets fields on several variants of one product — PATCH /admin/products/{product_id}/variations/bulk as a bare record; response is slim or full.
func (v *VariationsResource) UpdateProductVariations(ctx context.Context, productID any, payload any, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return v.requester.SendPatchRequest(ctx, ProductVariationsBulkPath, &http.RequestOptions{
		PathParams: map[string]any{"product_id": productID}, QueryParams: queryParams,
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ProductVariationsBulkResultShape})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "variations"

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

func requirePayload(payload any) error {
	if payload == nil || isNilValue(payload) {
		return wcerrors.NewUsageError("payload is required: this endpoint takes a body, and nil would " +
			"send the request with no body at all")
	}
	return nil
}

func isNilValue(value any) bool {
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Map, reflect.Slice, reflect.Pointer, reflect.Interface, reflect.Func:
		return reflected.IsNil()
	}
	return false
}

// New binds the variation type and product variant endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *VariationsResource {
	return &VariationsResource{requester: requester}
}
