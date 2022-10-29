package economy

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/PuerkitoBio/goquery"
	"github.com/sahilm/fuzzy"
)

var iso4217 = map[string]string{
	"Baht tailandés":                   "THB",
	"Balboa panameño":                  "PAB",
	"Bolívar fuerte venezolano":        "VEF",
	"Boliviano":                        "BOB",
	"Colón costarricense":              "CRC",
	"Corona Checa":                     "CZK",
	"Corona Danesa":                    "DKK",
	"Corona islandesa":                 "ISK",
	"Corona noruega":                   "NOK",
	"Corona sueca":                     "SEK",
	"DEG":                              "DEG",
	"Dírham de Emiratos Árabes Unidos": "AED",
	"Dírham Marroquí":                  "MAD",
	"Dólar australiano":                "AUD",
	"Dólar canadiense":                 "CAD",
	"Dólar de bermudas":                "BMD",
	"Dólar de Islas Caimán":            "KYD",
	"Dólar de las Bahamas":             "BSD",
	"Dólar singapurense":               "SGD",
	"Dolár fiyiano":                    "FJD",
	"Dólar hongkonés":                  "HKD",
	"Dólar neozelandés":                "NZD",
	"Dólar taiwanés":                   "TWD",
	"Euro":                             "EUR",
	"Forint húngaro":                   "HUF",
	"Franco de la Polinesia Francesa":  "XPF",
	"Franco suizo":                     "CHF",
	"Guaraní  paraguayo":               "PYG",
	"Hryvnia ucraniano":                "UAH",
	"Leu rumano":                       "RON",
	"Libra egipcia":                    "EGP",
	"Libra esterlina":                  "GBP",
	"Nueva lira turca":                 "TRY",
	"Nuevo sol peruano":                "PEN",
	"Peso argentino":                   "ARS",
	"Peso colombiano":                  "COP",
	"Peso cubano":                      "CUP",
	"Peso de República Dominicana":     "DOP",
	"Peso filipino":                    "PHP",
	"Peso mexicano":                    "MXN",
	"Peso uruguayo":                    "UYU",
	"Quetzal guatemalteco":             "GTQ",
	"Rand sudafricano":                 "ZAR",
	"Real brasileño":                   "BRL",
	"Rial iraní":                       "IRR",
	"Rial saudita":                     "SAR",
	"Ringgit malasio":                  "MYR",
	"Riyal Catarí":                     "QAR",
	"Rublo ruso":                       "RUB",
	"Rupia de Indonesia":               "IDR",
	"Rupia india":                      "INR",
	"Rupia pakistaní":                  "PKR",
	"Shekel israelí":                   "ILS",
	"Tenge de Kazajstán":               "KZT",
	"Won coreano":                      "KRW",
	"Yen":                              "JPY",
	"Yuan":                             "CNY",
	"Zloty polaco":                     "PLN",
}

type Currency struct {
	Name         string  `json:"name"`
	ISO4217      string  `json:"iso4217"`
	ExchangeRate float64 `json:"exchange_rate"`
}

func parseCurrenciesHTML(r io.ReadCloser) (currencies []*Currency, err error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	doc.Find("tr").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		name := strings.TrimSpace(s.Children().Get(0).FirstChild.Data)

		rateStr := s.Children().Get(1).FirstChild.Data

		var rate float64
		rate, err = parseCurrency(rateStr)
		if err != nil {
			err = fmt.Errorf("unable to parse currency '%s' with value '%s", name, rateStr)
			return false
		}

		code, exists := iso4217[name]
		if !exists {
			err = fmt.Errorf("unknown currency '%s', no matching iso4217 code", name)
			return false
		}

		currencies = append(currencies, &Currency{
			Name:         name,
			ISO4217:      code,
			ExchangeRate: rate,
		})

		return true
	})

	if err != nil {
		return nil, err
	}

	return
}

func filterCurrencies(currencies []*Currency, filter string) []*Currency {
	var hints []string
	for _, currency := range currencies {
		hints = append(hints, fmt.Sprintf("%s %s", removeTilde(currency.Name), currency.ISO4217))
	}

	matches := fuzzy.Find(removeTilde(filter), hints)

	var results []*Currency
	for _, match := range matches {
		results = append(results, currencies[match.Index])
	}

	return results
}

func removeTilde(text string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, text)
	return result
}
