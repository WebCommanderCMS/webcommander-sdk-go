package carts

// The admin cart and checkout path templates; the client fills {api_version}. {session_id} takes the session_id the cart start answered, and {cart_item_id} a line's cart_item_id.
const (
	AdminCartsPath = "/api/{api_version}/admin/orders/carts"

	AdminCartPath = "/api/{api_version}/admin/orders/carts/{session_id}"

	AdminCartItemsPath = "/api/{api_version}/admin/orders/carts/{session_id}/cart_items"

	AdminCartItemPath = "/api/{api_version}/admin/orders/carts/{session_id}/cart_items/{cart_item_id}"

	AdminCheckoutPath = "/api/{api_version}/admin/orders/checkout/{session_id}"

	AdminCheckoutConfirmPath = "/api/{api_version}/admin/orders/checkout/{session_id}/confirm"

	AdminCheckoutCouponPath = "/api/{api_version}/admin/orders/checkout/{session_id}/coupon"

	AdminCheckoutGiftCardPath = "/api/{api_version}/admin/orders/checkout/{session_id}/gift_card"

	AdminCheckoutDeliverySchedulesPath = "/api/{api_version}/admin/orders/checkout/{session_id}/delivery_schedules"
)

// AdminCartKey is the key every cart answer wraps its record in; the request bodies are bare.
const AdminCartKey = "cart"

// AdminCheckoutKey is the key every checkout answer wraps its record in; the request bodies are bare.
const AdminCheckoutKey = "checkout"

// AdminCartsTransactionNoHeader files the order a checkout places under this transaction number; it beats a transaction_no body field.
const AdminCartsTransactionNoHeader = "X-Transaction-No"
