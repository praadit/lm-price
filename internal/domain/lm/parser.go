package lm

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ParsePrices parses lm.txt HTML into structured rows.
func ParsePrices(html []byte) ([]LocationPrices, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}
	return parseFromDoc(doc), nil
}

// MarshalPricesJSON encodes rows as JSON (same shape as the HTTP API).
func MarshalPricesJSON(rows []LocationPrices) ([]byte, error) {
	return json.Marshal(rows)
}

// PricesHTMLToJSON parses HTML then returns JSON bytes.
func PricesHTMLToJSON(html []byte) ([]byte, error) {
	rows, err := ParsePrices(html)
	if err != nil {
		return nil, err
	}
	return MarshalPricesJSON(rows)
}

func parseFromDoc(doc *goquery.Document) []LocationPrices {
	var out []LocationPrices
	doc.Find("div.tab-pane").Each(func(_ int, pane *goquery.Selection) {
		table := pane.Find("table").First()
		if table.Length() == 0 {
			return
		}

		thead := table.Find("thead")
		product := cleanText(thead.Find("tr").Eq(0).Text())
		area := cleanText(thead.Find("tr").Eq(1).Text())

		butikNames := headerButikNames(thead)
		if len(butikNames) == 0 {
			return
		}

		byCol := make([][]Price, len(butikNames))
		table.Find("tbody tr").Each(func(_ int, tr *goquery.Selection) {
			tds := tr.Find("td")
			if tds.Length() == 0 {
				return
			}

			gramasiStr := cleanText(tds.Eq(0).Text())
			gramasi, err := strconv.ParseFloat(strings.ReplaceAll(gramasiStr, ",", "."), 64)
			if err != nil {
				return
			}

			for j := range butikNames {
				idx := j + 1
				if idx >= tds.Length() {
					break
				}
				q := cellToQuote(gramasi, tds.Eq(idx))
				byCol[j] = append(byCol[j], q)
			}
		})

		for j, name := range butikNames {
			if name == "" {
				continue
			}
			out = append(out, LocationPrices{
				Location: name,
				Product:  product,
				Area:     area,
				Prices:   byCol[j],
			})
		}
	})
	return out
}

func headerButikNames(thead *goquery.Selection) []string {
	var columns []string
	thead.Find("tr").Last().Find("th").Each(func(_ int, th *goquery.Selection) {
		txt := cleanText(th.Clone().Children().Remove().End().Text())
		if txt != "" && !strings.EqualFold(txt, "Gramasi") {
			columns = append(columns, txt)
		}
	})

	if len(columns) == 0 {
		thead.Find("tr").Each(func(_ int, tr *goquery.Selection) {
			var rowCols []string
			tr.Find("th").Each(func(_ int, th *goquery.Selection) {
				txt := cleanText(th.Clone().Children().Remove().End().Text())
				if txt != "" && !strings.EqualFold(txt, "Gramasi") {
					rowCols = append(rowCols, txt)
				}
			})
			if len(rowCols) > len(columns) {
				columns = rowCols
			}
		})
	}

	return columns
}

func cellToQuote(gramasi float64, td *goquery.Selection) Price {
	priceIDR, stockPtr, soldOut := parsePriceStockCell(td)
	stock := 0
	if stockPtr != nil {
		stock = *stockPtr
	}
	if soldOut {
		stock = 0
	}
	return Price{
		Gramasi: gramasi,
		Price:   priceIDR,
		Stock:   stock,
		SoldOut: soldOut,
	}
}

var reDigits = regexp.MustCompile(`[0-9]+`)

func parsePriceStockCell(td *goquery.Selection) (priceIDR int, stock *int, soldOut bool) {
	raw := cleanText(td.Text())
	if raw == "" {
		return 0, nil, false
	}

	parts := strings.Split(raw, "Stock:")
	pricePart := strings.TrimSpace(parts[0])
	priceIDR = parseIDR(pricePart)

	if len(parts) < 2 {
		return priceIDR, nil, false
	}

	stockPart := cleanText(parts[1])
	if stockPart == "" {
		return priceIDR, nil, false
	}
	if strings.EqualFold(stockPart, "Sold Out") {
		return priceIDR, nil, true
	}

	d := reDigits.FindString(stockPart)
	if d == "" {
		return priceIDR, nil, false
	}
	n, err := strconv.Atoi(d)
	if err != nil {
		return priceIDR, nil, false
	}
	return priceIDR, &n, false
}

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
