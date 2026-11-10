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

// CountProducts counts products — GET /admin/products/count; {"total": n, "count": n}. It takes the list's filters (search_text, name, sku, ids, parent, product_type, stock, is_available, is_inventory_enabled, created_by, price_from, price_to and the four date bounds) and refuses its paging and sorting, and status, as 400 unknown_parameter. A trashed product is not counted.
func (p *Products) CountProducts(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendGetRequest(ctx, ProductsCountPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ProductCountShape})
}

// CheckProductSku asks whether a SKU is free — GET /admin/products/check-sku. The match ignores case and sees a SKU held by a trashed product or by a variant; excludeProductID is the id of a product that may hold it, as when that product is being edited, and nil leaves it out. Any other query name, and SKU in capitals, is 400 unknown_parameter.
func (p *Products) CheckProductSku(ctx context.Context, sku string, excludeProductID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if sku == "" {
		return nil, wcerrors.NewUsageError("sku is required")
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	queryParams := map[string]any{"sku": sku}
	if excludeProductID != nil {
		queryParams["exclude_product_id"] = excludeProductID
	}
	return p.requester.SendGetRequest(ctx, ProductsCheckSkuPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ProductSkuCheckShape})
}

// ListSpecialProducts lists a special selection of products — GET /admin/products/special. type is required: onSale, featured, new, topSelling, lastBought or lastViewed. max (1 to 100, default 10) and offset page it. The last two are per customer and need customerId (the customer_id or internal_id of a customer) or customer_email; customer_id is not read, and an unknown customer is 400 unknown_reference. A name it does not know is ignored.
func (p *Products) ListSpecialProducts(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendGetRequest(ctx, ProductsSpecialPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ProductSpecialPageShape})
}

// BulkDeleteProducts moves several products to the trash — POST /admin/products/bulk-delete, answered 200 even when an id failed. The body is {"ids": [...]} (the collection sends product_ids, which the store answers 400 required for ids); an empty list is 400. Each id that is not a product is counted in failed and listed in errors.
func (p *Products) BulkDeleteProducts(ctx context.Context, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendPostRequest(ctx, ProductsBulkDeletePath, &http.RequestOptions{
		JSONBody: payload, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ProductsBulkDeleteResultShape})
}

// ListProductReviews lists a product's reviews — GET /admin/products/{product_id}/reviews; {"reviews": [...]} with no paging block. status, limit and offset are in the collection and this store ignores them, as it does any other name.
func (p *Products) ListProductReviews(ctx context.Context, productID any, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendGetRequest(ctx, ProductReviewsPath, &http.RequestOptions{
		PathParams:  map[string]any{"product_id": productID},
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ProductReviewListShape})
}

// GetProductReview reads one review under review — GET /admin/products/{product_id}/reviews/{review_id}. reviewID is a review's id; an unknown one is 404.
func (p *Products) GetProductReview(ctx context.Context, productID any, reviewID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendGetRequest(ctx, ProductReviewPath, &http.RequestOptions{
		PathParams:   map[string]any{"product_id": productID, "review_id": reviewID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ProductReviewShape})
}

// PatchProductReview changes one review — PATCH /admin/products/{product_id}/reviews/{review_id}, a bare body answered 200 with the review under review. status is pending (what an imported review without one reads) or approved; rejected, which the collection uses, is 400 invalid_enum. rating and comment change too.
func (p *Products) PatchProductReview(ctx context.Context, productID any, reviewID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendPatchRequest(ctx, ProductReviewPath, &http.RequestOptions{
		PathParams: map[string]any{"product_id": productID, "review_id": reviewID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ProductReviewShape})
}

// DeleteProductReview deletes one review for good; 204, and a second delete is 404. DELETE /admin/products/{product_id}/reviews/{review_id}
func (p *Products) DeleteProductReview(ctx context.Context, productID any, reviewID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendDeleteRequest(ctx, ProductReviewPath, &http.RequestOptions{
		PathParams:   map[string]any{"product_id": productID, "review_id": reviewID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// CreateProductInventory sets a product's stock figure — POST /admin/products/{product_id}/inventories, answered 201. The body is {"product": {"product_pricing_and_stock": {"track_inventory": true, "available_stock": n}}}; low_level_stock sets the low level. It records a history row and a variant keeps its own figure. The header is stored as the product's transaction_no when it has none, and the answer echoes the value from before the write.
func (p *Products) CreateProductInventory(ctx context.Context, productID any, payload any, transactionNo string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
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
	return p.requester.SendPostRequest(ctx, ProductInventoriesPath, &http.RequestOptions{
		PathParams: map[string]any{"product_id": productID}, JSONBody: payload,
		Headers: headers, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: ProductInventoryResultShape})
}

// CreateProductImage stores images on a product — POST /admin/products/{product_id}/image, answered 201 with the files just stored. The body is JSON, {"product": {"image_and_video": {"images": [{"filename": ..., "base64": ...}]}}}; a multipart body is 415.
func (p *Products) CreateProductImage(ctx context.Context, productID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendPostRequest(ctx, ProductImagePath, &http.RequestOptions{
		PathParams: map[string]any{"product_id": productID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ProductMediaUploadedShape})
}

// CreateProductMedia stores one file on a product — POST /admin/products/{product_id}/media, answered 201 with the stored image's row, flat, not under product. A bare JSON body of type, filename and content (base64 or a data: URI). Only type images is accepted; videos and specification_file are 400 not_supported, and position is 400 unknown_field.
func (p *Products) CreateProductMedia(ctx context.Context, productID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendPostRequest(ctx, ProductMediaPath, &http.RequestOptions{
		PathParams: map[string]any{"product_id": productID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ProductMediaImageShape})
}

// DeleteProductMediaImage deletes one image from a product; 204, and a second delete is 404. imageID is an image's id on the product. DELETE /admin/products/{product_id}/media/images/{image_id}
func (p *Products) DeleteProductMediaImage(ctx context.Context, productID any, imageID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendDeleteRequest(ctx, ProductMediaImagePath, &http.RequestOptions{
		PathParams:   map[string]any{"product_id": productID, "image_id": imageID},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// UpdateProductMediaOrder sets the order of a product's images — PUT /admin/products/{product_id}/media/order, a bare {"type": "images", "ids": [...]} naming every image once (400 invalid_set otherwise). The first id becomes primary, and the answer is {"images": [...]} in the new order.
func (p *Products) UpdateProductMediaOrder(ctx context.Context, productID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendPutRequest(ctx, ProductMediaOrderPath, &http.RequestOptions{
		PathParams: map[string]any{"product_id": productID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ProductMediaOrderShape})
}

// PatchProductMediaImage moves one image to a 1-based position — PATCH /admin/products/{product_id}/media/images/{image_id}; a bare {"position": n} is the only body (any other is 400 invalid_body), and the image is answered as one flat row. imageID is an image's id on the product.
func (p *Products) PatchProductMediaImage(ctx context.Context, productID any, imageID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendPatchRequest(ctx, ProductMediaImagePath, &http.RequestOptions{
		PathParams: map[string]any{"product_id": productID, "image_id": imageID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ProductMediaImageShape})
}

// GetProductSeo reads a product's SEO, analytics, custom code, social and AI visibility settings under seo — GET /admin/products/{product_id}/seo. It has no HEAD (404) and ignores field_metadata.
func (p *Products) GetProductSeo(ctx context.Context, productID any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendGetRequest(ctx, ProductSeoPath, &http.RequestOptions{
		PathParams:   map[string]any{"product_id": productID},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ProductSeoShape})
}

// PatchProductStock changes a product's stock — PATCH /admin/products/{product_id}/stock; a bare body of tracking, locations, note and transaction_no, at least one of them (a wrapper or a flat available_stock is 400 unknown_field). The answer holds the stock now and the adjustment it recorded, null when the figure did not change.
func (p *Products) PatchProductStock(ctx context.Context, productID any, payload any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	if err := requirePayload(payload); err != nil {
		return nil, err
	}
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendPatchRequest(ctx, ProductStockPath, &http.RequestOptions{
		PathParams: map[string]any{"product_id": productID}, JSONBody: payload,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ProductStockPatchResultShape})
}

// ListProductStockAdjustments lists a product's stock history, newest first — GET /admin/products/{product_id}/stock/adjustments, paged by limit (default 20) and offset. Each row reads transaction_no null, also for a change sent with one; HEAD is 404.
func (p *Products) ListProductStockAdjustments(ctx context.Context, productID any, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return p.requester.SendGetRequest(ctx, ProductStockAdjustmentsPath, &http.RequestOptions{
		PathParams: map[string]any{"product_id": productID}, QueryParams: queryParams,
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: ProductStockAdjustmentPageShape})
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
