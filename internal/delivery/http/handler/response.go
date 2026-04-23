package handler

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Standard API envelope used by all endpoints.
type PricesEnvelope struct {
	LastUpdate string         `json:"last_update"`
	Source     string         `json:"source"`
	Data       []LocationData `json:"data"`
}

type LocationData struct {
	Location string       `json:"location"`
	Product  string       `json:"product,omitempty"`
	Area     string       `json:"area"`
	Prices   []PriceEntry `json:"prices"`
}

type PriceEntry struct {
	Gramasi   float64 `json:"gramasi"`
	BuyPrice  int     `json:"buy_price"`
	SellPrice int     `json:"sell_price"`
	Stock     int     `json:"stock"`
	SoldOut   bool    `json:"sold_out"`
}

func timeToRFC3339String(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

func indonesiaIfEmpty(s string) string {
	if strings.TrimSpace(s) == "" {
		return "Indonesia"
	}
	return s
}

var reLeadingNumber = regexp.MustCompile(`^[\s]*([0-9]+(?:[.,][0-9]+)?)`)

func parseLeadingFloat(s string) float64 {
	m := reLeadingNumber.FindStringSubmatch(strings.TrimSpace(s))
	if len(m) < 2 {
		return 0
	}
	v := strings.ReplaceAll(m[1], ",", ".")
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0
	}
	return f
}
