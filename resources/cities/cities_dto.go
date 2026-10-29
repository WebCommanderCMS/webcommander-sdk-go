package cities

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// CityState is the state a city row carries; ID is the one a country's states listing answers, and Default follows the administration section's default_state.
type CityState struct {
	ID      *int64
	Code    *string
	Name    *string
	Active  *bool
	Default *bool
}

// CityCountry is the country a city row carries; ID is the one the countries listing answers.
type CityCountry struct {
	ID   *int64
	Code *string
	Name *string
}

// City is one listing row.
type City struct {
	ID       *int64
	Name     *string
	PostCode *string
	State    *CityState
	Country  *CityCountry
}

// Pagination is the listing's five-key block: Total matched, Limit the page size applied, Offset the start.
type Pagination struct {
	Total       *int64
	Limit       *int64
	Offset      *int64
	HasPrevious *bool
	HasNext     *bool
}

// CityPage is one page of the city listing: the rows and the pagination block.
type CityPage struct {
	Cities     []*City
	Pagination *Pagination
}

// CityPageShape is the typed view of the listing body.
func CityPageShape(decoded any) any { return dto.Read[CityPage](decoded) }
