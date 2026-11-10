package orderpreviews

// The path templates; the client fills {api_version}. {preview_session_id} takes the session_id the start answered.
const (
	OrderPreviewsPath = "/api/{api_version}/admin/orders/preview"

	OrderPreviewPath = "/api/{api_version}/admin/orders/preview/{preview_session_id}"

	OrderPreviewConfirmPath = "/api/{api_version}/admin/orders/preview/{preview_session_id}/confirm"

	OrderPreviewCouponPath = "/api/{api_version}/admin/orders/preview/{preview_session_id}/coupon"

	OrderPreviewGiftCardPath = "/api/{api_version}/admin/orders/preview/{preview_session_id}/gift-card"
)

// OrderPreviewKey is the key every preview answer wraps its record in; the request bodies are bare.
const OrderPreviewKey = "order"
