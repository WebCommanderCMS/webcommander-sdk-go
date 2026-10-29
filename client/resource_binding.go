package client

import (
	"maps"
	"slices"
	"strings"
	"sync"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
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

// ResourceFactory builds one resource family over a client.
type ResourceFactory func(requester http.Requester) any

// resourceModules is the registry, guarded because registration is process-wide and may race.
var (
	resourceModulesMutex sync.RWMutex
	resourceModules      = map[string]ResourceFactory{
		customers.ResourceName: func(requester http.Requester) any { return customers.New(requester) },
		orderflags.ResourceName: func(requester http.Requester) any {
			return orderflags.New(requester)
		},
		inventorylocations.ResourceName: func(requester http.Requester) any {
			return inventorylocations.New(requester)
		},
		products.ResourceName: func(requester http.Requester) any {
			return products.New(requester)
		},
		orders.ResourceName: func(requester http.Requester) any {
			return orders.New(requester)
		},
		aliases.ResourceName: func(requester http.Requester) any {
			return aliases.New(requester)
		},
		carts.ResourceName: func(requester http.Requester) any {
			return carts.New(requester)
		},
		blogs.ResourceName: func(requester http.Requester) any {
			return blogs.New(requester)
		},
		ordercustomfields.ResourceName: func(requester http.Requester) any {
			return ordercustomfields.New(requester)
		},
		blogcategories.ResourceName: func(requester http.Requester) any {
			return blogcategories.New(requester)
		},
		categories.ResourceName: func(requester http.Requester) any {
			return categories.New(requester)
		},
		cities.ResourceName: func(requester http.Requester) any {
			return cities.New(requester)
		},
		commonsettings.ResourceName: func(requester http.Requester) any {
			return commonsettings.New(requester)
		},
		countries.ResourceName: func(requester http.Requester) any {
			return countries.New(requester)
		},
		currencies.ResourceName: func(requester http.Requester) any {
			return currencies.New(requester)
		},
		paymentgateways.ResourceName: func(requester http.Requester) any {
			return paymentgateways.New(requester)
		},
		shippingclasses.ResourceName: func(requester http.Requester) any {
			return shippingclasses.New(requester)
		},
		shippingmethods.ResourceName: func(requester http.Requester) any {
			return shippingmethods.New(requester)
		},
		pagesettings.ResourceName: func(requester http.Requester) any {
			return pagesettings.New(requester)
		},
		pagesettings.CategoryPageResourceName: func(requester http.Requester) any {
			return pagesettings.NewCategoryPageSettings(requester)
		},
		pagesettings.ShoppingCartResourceName: func(requester http.Requester) any {
			return pagesettings.NewShoppingCartSettings(requester)
		},
		pagesettings.ProductPageResourceName: func(requester http.Requester) any {
			return pagesettings.NewProductPageSettings(requester)
		},
		pagesettings.CheckoutPageResourceName: func(requester http.Requester) any {
			return pagesettings.NewCheckoutPageSettings(requester)
		},
		pricedisplaysettings.ResourceName: func(requester http.Requester) any {
			return pricedisplaysettings.New(requester)
		},
		animationeffects.ResourceName: func(requester http.Requester) any {
			return animationeffects.New(requester)
		},
		productlisting.ResourceName: func(requester http.Requester) any {
			return productlisting.New(requester)
		},
		customerportal.ResourceName: func(requester http.Requester) any {
			return customerportal.New(requester)
		},
		sitesettings.ResourceName: func(requester http.Requester) any {
			return sitesettings.New(requester)
		},
		shippingsettings.ResourceName: func(requester http.Requester) any {
			return shippingsettings.New(requester)
		},
		seoconfiguration.ResourceName: func(requester http.Requester) any {
			return seoconfiguration.New(requester)
		},
		searchsettings.ResourceName: func(requester http.Requester) any {
			return searchsettings.New(requester)
		},
		ecommercesettings.ResourceName: func(requester http.Requester) any {
			return ecommercesettings.New(requester)
		},
		taxcodes.ResourceName: func(requester http.Requester) any {
			return taxcodes.New(requester)
		},
		taxprofiles.ResourceName: func(requester http.Requester) any {
			return taxprofiles.New(requester)
		},
		taxrules.ResourceName: func(requester http.Requester) any {
			return taxrules.New(requester)
		},
		taxzones.ResourceName: func(requester http.Requester) any {
			return taxzones.New(requester)
		},
		taxsettings.ResourceName: func(requester http.Requester) any {
			return taxsettings.New(requester)
		},
		taxcountrytaxes.ResourceName: func(requester http.Requester) any {
			return taxcountrytaxes.New(requester)
		},
		taxcustomprofiles.ResourceName: func(requester http.Requester) any {
			return taxcustomprofiles.New(requester)
		},
		discountprofiles.ResourceName: func(requester http.Requester) any {
			return discountprofiles.New(requester)
		},
		customergroups.ResourceName: func(requester http.Requester) any {
			return customergroups.New(requester)
		},
		variations.ResourceName: func(requester http.Requester) any {
			return variations.New(requester)
		},
		transactions.ResourceName: func(requester http.Requester) any {
			return transactions.New(requester)
		},
		trash.ResourceName: func(requester http.Requester) any {
			return trash.New(requester)
		},
		email.ResourceName: func(requester http.Requester) any {
			return email.New(requester)
		},
		pages.ResourceName: func(requester http.Requester) any {
			return pages.New(requester)
		},
	}
)

// RegisterResourceModule makes a resource family reachable as one of a client's resources.
func RegisterResourceModule(resourceName string, factory ResourceFactory) error {
	if err := validateResourceName(resourceName); err != nil {
		return err
	}
	if factory == nil {
		return wcerrors.NewUsageError("RegisterResourceModule: " + resourceName +
			" has no factory — nothing would be reachable under that name")
	}
	resourceModulesMutex.Lock()
	defer resourceModulesMutex.Unlock()
	if _, found := resourceModules[resourceName]; found {
		// Any re-registration is refused: Go cannot compare function values, so a repeat is a clash.
		return wcerrors.NewUsageError("RegisterResourceModule: " + resourceName +
			" is already registered — one name cannot answer to two resources")
	}
	resourceModules[resourceName] = factory
	return nil
}

// validateResourceName refuses a name no lookup could ever reach.
func validateResourceName(resourceName string) error {
	if resourceName == "" {
		return wcerrors.NewUsageError("RegisterResourceModule: the resource name is blank")
	}
	for _, r := range resourceName {
		if !(r >= 'a' && r <= 'z') && !(r >= '0' && r <= '9') && r != '_' {
			return wcerrors.NewUsageError("RegisterResourceModule: " + resourceName +
				" is not lower snake_case, which is how a resource name is spelled")
		}
	}
	return nil
}

// RegisteredResourceNames returns every registered resource name, sorted.
func RegisteredResourceNames() []string {
	resourceModulesMutex.RLock()
	defer resourceModulesMutex.RUnlock()
	return registeredNamesLocked()
}

// ResolveResourceModule returns the factory registered under a name.
func ResolveResourceModule(resourceName string) (ResourceFactory, error) {
	resourceModulesMutex.RLock()
	defer resourceModulesMutex.RUnlock()
	factory, found := resourceModules[resourceName]
	if !found {
		return nil, wcerrors.NewAttributeError("no resource named " + resourceName +
			"; registered: " + strings.Join(registeredNamesLocked(), ", "))
	}
	return factory, nil
}

// registeredNamesLocked lists the names with the lock already held.
func registeredNamesLocked() []string {
	return slices.Sorted(maps.Keys(resourceModules))
}

// BoundResource is a resource family with a client already applied, constructed once.
type BoundResource struct {
	// Name is the registered resource name, e.g. "admin_customers".
	Name string
	// Value is the constructed family; assert it to its concrete type, or use the client's typed accessor.
	Value any
}

// NewBoundResource constructs the resource registered under a name over one requester.
func NewBoundResource(requester http.Requester, resourceName string) (*BoundResource, error) {
	factory, err := ResolveResourceModule(resourceName)
	if err != nil {
		return nil, err
	}
	return &BoundResource{Name: resourceName, Value: factory(requester)}, nil
}
