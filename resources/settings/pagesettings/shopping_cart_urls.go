package pagesettings

// The path templates for the shopping cart settings endpoints; the client fills {api_version}.
const (
	ShoppingCartPath = "/api/{api_version}/admin/settings/page_settings/shopping_cart"

	ShoppingCartSectionPath = "/api/{api_version}/admin/settings/page_settings/shopping_cart/{section}"
)

// shoppingCartSections are the two values {section} takes; both are writable.
var shoppingCartSections = [2]string{"cart_page", "flying_cart"}

// ShoppingCartSections returns a fresh copy of the two values {section} takes.
func ShoppingCartSections() []string {
	return append(make([]string, 0, len(shoppingCartSections)), shoppingCartSections[:]...)
}
