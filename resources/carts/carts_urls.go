package carts

// The path templates; the client fills {api_version}. The cart paths sit under /carts, not /admin.
const (
	CartsPath = "/api/{api_version}/carts"

	CartPath = "/api/{api_version}/carts/{session_id}"

	CartItemPath = "/api/{api_version}/carts/{session_id}/cart_items/{cart_item_id}"

	// CartOrderPath is the one cart path under /admin.
	CartOrderPath = "/api/{api_version}/admin/orders/carts/{session_id}"
)

// CartKey is the key a cart arrives under; no write takes it as a wrapper.
const CartKey = "cart"

// CartSessionKey is the key the session write answers under, spelled as the store spells it.
const CartSessionKey = "sessionId"

// CartOrderKey is the key the order a cart became arrives under, beside status.
const CartOrderKey = "order"

// CartsTransactionNoHeader files the order a cart becomes under this transaction number; it beats a transaction_no body field.
const CartsTransactionNoHeader = "X-Transaction-No"
