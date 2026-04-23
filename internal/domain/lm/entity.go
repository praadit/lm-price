package lm

import "time"

// PricesResponse is the JSON envelope for parsed LM rows from the upstream document.
type PricesResponse struct {
	LastUpdate time.Time        `json:"last_update,omitempty"`
	Data       []LocationPrices `json:"data"`
}

// LocationPrices is one butik row in the public JSON shape.
type LocationPrices struct {
	Location string  `json:"location"`
	Product  string  `json:"product,omitempty"`
	Area     string  `json:"area,omitempty"`
	Prices   []Price `json:"prices"`
}

// Price is a single gram / price / stock row for one location.
type Price struct {
	Gramasi float64 `json:"gramasi"`
	Price   int     `json:"price"`
	Stock   int     `json:"stock"`
	SoldOut bool    `json:"sold_out"`
}

// QueryValidationError is returned when area or location query params are not in the scrape.
type QueryValidationError struct {
	Code               string   `json:"code"`
	Message            string   `json:"message"`
	AvailableAreas     []string `json:"available_areas,omitempty"`
	AvailableLocations []string `json:"available_locations,omitempty"`
	RequestedArea      string   `json:"requested_area,omitempty"`
	RequestedLocation  string   `json:"requested_location,omitempty"`
}

func (e *QueryValidationError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}
