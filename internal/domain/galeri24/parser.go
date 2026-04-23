package galeri24

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// ParsePricesByHeader parses https://galeri24.co.id/harga-emas and extracts the section matching the provided header,
// e.g. "ANTAM", "UBS", "GALERI 24", etc.
func ParsePricesByHeader(html []byte, header string) (PricesResponse, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return PricesResponse{}, err
	}

	root := findSectionByHeader(doc, header)
	lastUpdate := extractLastUpdate(root)
	rows := extractRows(root)

	return PricesResponse{
		LastUpdate: lastUpdate,
		Data:       rows,
	}, nil
}

// ParseAntamPricesDocument is a convenience wrapper for the "Harga ANTAM" section.
func ParseAntamPricesDocument(html []byte) (PricesResponse, error) {
	return ParsePricesByHeader(html, "ANTAM")
}

func findSectionByHeader(doc *goquery.Document, header string) *goquery.Selection {
	h := strings.ToUpper(cleanText(header))
	if h == "" {
		return doc.Selection
	}

	// Preferred: the site uses id="<VENDOR>", e.g. id="ANTAM".
	if s := doc.Find("#" + h); s.Length() > 0 {
		return s
	}

	want := "HARGA " + h
	var found *goquery.Selection
	doc.Find("*").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		if strings.ToUpper(cleanText(s.Text())) == want {
			found = s
			return false
		}
		return true
	})
	if found == nil {
		return doc.Selection
	}

	// Walk up a bit so we include the update line and the rows.
	root := found
	for i := 0; i < 6; i++ {
		p := root.Parent()
		if p.Length() == 0 {
			break
		}
		root = p
		if strings.Contains(strings.ToLower(root.Text()), "diperbarui") {
			break
		}
	}
	return root
}

func extractLastUpdate(root *goquery.Selection) time.Time {
	if root == nil || root.Length() == 0 {
		return time.Time{}
	}

	// Example: "Diperbarui Kamis, 23 April 2026"
	var rest string
	root.Find("div").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		txt := cleanText(s.Text())
		if txt == "" {
			return true
		}
		ltxt := strings.ToLower(strings.TrimSpace(txt))
		if !strings.HasPrefix(ltxt, "diperbarui ") {
			return true
		}
		rest = cleanText(strings.TrimSpace(txt[len("Diperbarui"):]))
		return false
	})
	if rest == "" {
		return time.Time{}
	}

	// The first "Diperbarui ..." node can contain extra text (it may include the whole section).
	// Extract only the "23 April 2026" date portion.
	if i := strings.Index(rest, ","); i >= 0 {
		rest = cleanText(rest[i+1:])
	}
	date := reDate.FindString(rest)
	if date == "" {
		return time.Time{}
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("WIB", 7*3600)
	}
	t, err := time.ParseInLocation("2 January 2006", date, loc)
	if err != nil {
		return time.Time{}
	}
	return t
}

var reDate = regexp.MustCompile(`\b[0-9]{1,2}\s+\p{L}+\s+[0-9]{4}\b`)

func extractRows(root *goquery.Selection) []PriceRow {
	if root == nil || root.Length() == 0 {
		return nil
	}

	// Rows are rendered as div.grid with 3 meaningful cells: weight, sell, buyback.
	var out []PriceRow
	root.Find("div.grid").Each(func(_ int, row *goquery.Selection) {
		cols := row.Children()
		if cols.Length() < 3 {
			return
		}

		w := parseFloat(cleanText(cols.Eq(0).Text()))
		sell := parseIDR(cleanText(cols.Eq(1).Text()))
		buyback := parseIDR(cleanText(cols.Eq(2).Text()))

		// Skip header rows.
		if w == 0 || sell == 0 || buyback == 0 {
			return
		}

		out = append(out, PriceRow{
			Weight:       w,
			SellPrice:    sell,
			BuybackPrice: buyback,
		})
	})
	return out
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

func parseFloat(s string) float64 {
	s = strings.ReplaceAll(s, ",", ".")
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

func cleanText(s string) string {
	s = strings.ReplaceAll(s, "\u00a0", " ")
	s = strings.TrimSpace(s)
	s = strings.Join(strings.Fields(s), " ")
	return s
}
