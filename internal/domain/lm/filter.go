package lm

import (
	"fmt"
	"sort"
	"strings"
)

// FilterPrices returns rows matching optional area and/or location query strings.
// Matching is case-insensitive with whitespace normalized. Empty filters keep all rows.
func FilterPrices(rows []LocationPrices, areaQuery, locationQuery string) ([]LocationPrices, error) {
	areaQ := strings.TrimSpace(areaQuery)
	locQ := strings.TrimSpace(locationQuery)

	if areaQ != "" && !hasMatchingArea(rows, areaQ) {
		return nil, &QueryValidationError{
			Code:                "unknown_area",
			Message:             fmt.Sprintf("unknown area %q", areaQ),
			AvailableAreas:      UniqueAreas(rows),
			RequestedArea:       areaQ,
			RequestedLocation:   locQ,
		}
	}

	withinArea := rows
	if areaQ != "" {
		withinArea = filterByArea(rows, areaQ)
	}

	if locQ != "" && !hasMatchingLocation(withinArea, locQ) {
		locList := UniqueLocations(withinArea)
		if areaQ == "" {
			locList = UniqueLocations(rows)
		}
		return nil, &QueryValidationError{
			Code:               "unknown_location",
			Message:            fmt.Sprintf("unknown location %q", locQ),
			AvailableAreas:     UniqueAreas(rows),
			AvailableLocations: locList,
			RequestedArea:      areaQ,
			RequestedLocation:  locQ,
		}
	}

	if areaQ == "" && locQ == "" {
		return rows, nil
	}

	var out []LocationPrices
	for _, it := range rows {
		if areaQ != "" && !areaEqual(it.Area, areaQ) {
			continue
		}
		if locQ != "" && !locationEqual(it.Location, locQ) {
			continue
		}
		out = append(out, it)
	}
	return out, nil
}

// UniqueAreas returns sorted distinct area labels from the scrape.
func UniqueAreas(rows []LocationPrices) []string {
	seen := map[string]string{}
	for _, it := range rows {
		if it.Area == "" {
			continue
		}
		k := normKey(it.Area)
		if _, ok := seen[k]; !ok {
			seen[k] = it.Area
		}
	}
	out := make([]string, 0, len(seen))
	for _, v := range seen {
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}

// UniqueLocations returns sorted distinct location (butik) names from the given slice.
func UniqueLocations(rows []LocationPrices) []string {
	seen := map[string]string{}
	for _, it := range rows {
		if it.Location == "" {
			continue
		}
		k := normKey(it.Location)
		if _, ok := seen[k]; !ok {
			seen[k] = it.Location
		}
	}
	out := make([]string, 0, len(seen))
	for _, v := range seen {
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}

func normKey(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(s)), " "))
}

func areaEqual(a, b string) bool {
	return normKey(a) == normKey(b)
}

func locationEqual(a, b string) bool {
	return normKey(a) == normKey(b)
}

func hasMatchingArea(rows []LocationPrices, areaQuery string) bool {
	q := normKey(areaQuery)
	if q == "" {
		return false
	}
	for _, it := range rows {
		if normKey(it.Area) == q {
			return true
		}
	}
	return false
}

func hasMatchingLocation(rows []LocationPrices, locQuery string) bool {
	q := normKey(locQuery)
	if q == "" {
		return false
	}
	for _, it := range rows {
		if normKey(it.Location) == q {
			return true
		}
	}
	return false
}

func filterByArea(rows []LocationPrices, areaQuery string) []LocationPrices {
	var out []LocationPrices
	for _, it := range rows {
		if areaEqual(it.Area, areaQuery) {
			out = append(out, it)
		}
	}
	return out
}
