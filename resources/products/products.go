package products

import (
	"context"
	"fmt"
	"reflect"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// Products is the product, review import and stock endpoints, bound to one client.
type Products struct {
	requester http.Requester
}

// ListProducts lists products, the most recently updated first — GET /admin/products; each row is the full record GetProduct answers, and a trashed product is not listed. limit is 1 to 100 (default 20, above 100 is 400), with offset or a 1-based page. search_text matches name or sku; name, sku, ids (repeat the name for more; a comma list is 500), parent (a category id), product_type, stock (in, out or low), is_available, is_inventory_enabled, created_by, price_from, price_to, created_from, created_to, updated_from and updated_to narrow it. sort takes name, isAvailable, sku, parent, basePrice, isInventoryEnabled, createdBy, created, updated or id, with dir; field_metadata adds a schema. Any other name is 400.
func (p *Products) ListProducts(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendGetRequest(ctx, ProductsPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ProductPageShape})
}

// CreateProduct creates a product — POST /admin/products as a bare record, answered 201; name and base_price are required, and sku must be free or left to the store.
func (p *Products) CreateProduct(ctx context.Context, payload any, transactionNo string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	var headers map[string]string
	if transactionNo != "" {
		headers = map[string]string{ProductsTransactionNoHeader: transactionNo}
	}
	return p.requester.SendPostRequest(ctx, ProductsPath, &http.RequestOptions{
		JSONBody: payload, Headers: headers, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ProductShape})
}

// GetProduct retrieves one product — GET /admin/products/{product_id}; field_metadata=true adds a field schema beside product, and any other name is 400 unknown_parameter.
func (p *Products) GetProduct(ctx context.Context, productID any, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendGetRequest(ctx, ProductPath, &http.RequestOptions{
		PathParams:  map[string]any{"product_id": productID},
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ProductShape})
}

// UpdateProduct updates one product — PUT /admin/products/{product_id}; a bare record merged at every level; it must carry name, sku, base_price, product_type, available, is_on_sale, available_on_date_range, available_on_weekdays, available_for and selected_customers (400 required otherwise), and fields left out keep their values. {"product": {…}} is 400 unknown_field.
func (p *Products) UpdateProduct(ctx context.Context, productID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendPutRequest(ctx, ProductPath, &http.RequestOptions{
		PathParams: map[string]any{"product_id": productID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ProductShape})
}

// PatchProduct merges into one product — PATCH /admin/products/{product_id}; the same bare record the PUT takes, and a product in the trash is 409.
func (p *Products) PatchProduct(ctx context.Context, productID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendPatchRequest(ctx, ProductPath, &http.RequestOptions{
		PathParams: map[string]any{"product_id": productID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ProductShape})
}

// DeleteProduct moves a product to the trash — DELETE /admin/products/{product_id}; 204, it still reads back with status trashed, and a second delete is 404.
func (p *Products) DeleteProduct(ctx context.Context, productID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendDeleteRequest(ctx, ProductPath, &http.RequestOptions{
		PathParams:   map[string]any{"product_id": productID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ImportProductReviews adds reviews to a product — POST /admin/products/{product_id}/reviews with {"product": {"reviews": [...]}}; 201 with every review the product holds.
func (p *Products) ImportProductReviews(ctx context.Context, productID any, payload any, transactionNo string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	var headers map[string]string
	if transactionNo != "" {
		headers = map[string]string{ProductsTransactionNoHeader: transactionNo}
	}
	return p.requester.SendPostRequest(ctx, ProductReviewsPath, &http.RequestOptions{
		PathParams: map[string]any{"product_id": productID}, JSONBody: payload,
		Headers: headers, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ProductReviewsShape})
}

// GetProductStock reads one product's stock at each location — GET /admin/products/{product_id}/stock; the store ignores every query name here, so the method takes none.
func (p *Products) GetProductStock(ctx context.Context, productID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendGetRequest(ctx, ProductStockPath, &http.RequestOptions{
		PathParams:   map[string]any{"product_id": productID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ProductStockShape})
}

// UpdateProductStock sets the product's stock — PUT /admin/products/{product_id}/stock; locations replaces the whole set, so a location left out is unassigned whatever it holds.
func (p *Products) UpdateProductStock(ctx context.Context, productID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendPutRequest(ctx, ProductStockPath, &http.RequestOptions{
		PathParams: map[string]any{"product_id": productID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ProductStockShape})
}

// DeleteProductStockLocation unassigns one location from the product — DELETE /admin/products/{product_id}/stock/locations/{inventory_location_id}.
func (p *Products) DeleteProductStockLocation(ctx context.Context, productID any, inventoryLocationID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendDeleteRequest(ctx, ProductStockLocationPath, &http.RequestOptions{
		PathParams:   map[string]any{"product_id": productID, "inventory_location_id": inventoryLocationID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "products"

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

// New binds the product endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *Products {
	return &Products{requester: requester}
}
