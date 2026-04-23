package antaremas

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// ParsePricesDocument parses the Antaremas "Ukuran / Harga Beli" table and the closest "Terakhir Diperbarui" timestamp.
func ParsePricesDocument(html []byte) (PricesResponse, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return PricesResponse{}, err
	}

	tableSel := findUkuranHargaBeliTable(doc)
	rows := parseTableRows(tableSel)
	lastUpdate := extractLastUpdateNear(tableSel)

	return PricesResponse{
		LastUpdate: lastUpdate,
		Data:       rows,
	}, nil
}

func findUkuranHargaBeliTable(doc *goquery.Document) *goquery.Selection {
	var found *goquery.Selection
	doc.Find("table").EachWithBreak(func(_ int, t *goquery.Selection) bool {
		ths := t.Find("thead th")
		if ths.Length() < 2 {
			return true
		}
		h0 := cleanText(ths.Eq(0).Text())
		h1 := cleanText(ths.Eq(1).Text())
		if strings.EqualFold(h0, "Ukuran") && strings.EqualFold(h1, "Harga Beli") {
			found = t
			return false
		}
		return true
	})
	if found == nil {
		return doc.Selection // empty
	}
	return found
}

func parseTableRows(table *goquery.Selection) []PriceRow {
	if table == nil || table.Length() == 0 {
		return nil
	}
	var out []PriceRow
	table.Find("tbody tr").Each(func(_ int, tr *goquery.Selection) {
		tds := tr.Find("td")
		if tds.Length() < 2 {
			return
		}
		size := cleanText(tds.Eq(0).Text())
		price := parseIDR(cleanText(tds.Eq(1).Text()))
		if size == "" || price == 0 {
			return
		}
		out = append(out, PriceRow{
			Size:     size,
			BuyPrice: price,
		})
	})
	return out
}

// Example text: "Terakhir Diperbarui 23 April 2026 09:30 WIB"
var reLastUpdate = regexp.MustCompile(`(?i)Terakhir\s+Diperbarui\s+(.+)$`)

func extractLastUpdateNear(table *goquery.Selection) time.Time {
	if table == nil || table.Length() == 0 {
		return time.Time{}
	}

	// The page often contains multiple "Terakhir Diperbarui ..." nodes in the same section (desktop/mobile variants).
	// We parse all of them in the table's nearest ancestor container and take the latest parsed instant.
	ancestor := table
	for i := 0; i < 6; i++ {
		if ancestor.Parent().Length() == 0 {
			break
		}
		ancestor = ancestor.Parent()
	}

	var best time.Time
	ancestor.Find("*").Each(func(_ int, s *goquery.Selection) {
		txt := cleanText(s.Text())
		if !strings.Contains(strings.ToLower(txt), "terakhir diperbarui") {
			return
		}
		m := reLastUpdate.FindStringSubmatch(txt)
		if len(m) < 2 {
			return
		}
		// Only parse the portion after the "Terakhir Diperbarui" prefix.
		if t := parseIndonesianDateTime(m[1]); !t.IsZero() && t.After(best) {
			best = t
		}
	})
	return best
}

func parseIndonesianDateTime(s string) time.Time {
	s = cleanText(s)
	if s == "" {
		return time.Time{}
	}

	// Antaremas sometimes uses 09.00 or 09:30.
	s = normalizeClockSeparators(s)

	// Many strings end with WIB/WITA/WIT.
	s = stripIndonesianTZSuffix(s)

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("WIB", 7*3600)
	}

	// Examples:
	// - 23 April 2026 09:30
	// - 23 April 2026 09:00
	layouts := []string{
		"2 January 2006 15:04",
		"2 January 2006 15:04:05",
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, loc); err == nil {
			return t
		}
	}
	return time.Time{}
}

func normalizeClockSeparators(s string) string {
	// Only normalize the time portion (replace '.' with ':') so we don't affect "Rp. 1.234.000" style text.
	parts := strings.Fields(s)
	if len(parts) == 0 {
		return s
	}
	last := parts[len(parts)-1]
	if strings.Count(last, ".") == 1 && reDigits.MatchString(last) && strings.Contains(last, ".") {
		parts[len(parts)-1] = strings.ReplaceAll(last, ".", ":")
		return strings.Join(parts, " ")
	}
	return s
}

func stripIndonesianTZSuffix(s string) string {
	s = strings.TrimSpace(s)
	for _, tz := range []string{" WIB", " WITA", " WIT"} {
		if len(s) >= len(tz) {
			suf := s[len(s)-len(tz):]
			if strings.EqualFold(suf, tz) {
				return strings.TrimSpace(s[:len(s)-len(tz)])
			}
		}
	}
	return s
}

var reDigits = regexp.MustCompile(`[0-9]+`)

func parseIDR(s string) int {
	d := strings.Join(reDigits.FindAllString(s, -1), "")
	if d == "" {
		return 0
	}
	n, err := strconv.Atoi(d)
	if err != nil {
		return 0
	}
	return n
}

func cleanText(s string) string {
	s = strings.ReplaceAll(s, "\u00a0", " ")
	s = strings.TrimSpace(s)
	s = strings.Join(strings.Fields(s), " ")
	return s
}
