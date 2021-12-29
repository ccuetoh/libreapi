package economy

import (
	"context"
	"fmt"

	"io"
	"net/http"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/sahilm/fuzzy"
)

var ISO4217Dict = map[string]string{
	"Baht tailandés":                   "THB",
	"Balboa panameño":                  "PAB",
	"Bolívar fuerte venezolano":        "VEF",
	"Boliviano":                        "BOB",
	"Colón costarricense":              "CRC",
	"Corona checa":                     "CZK",
	"Corona danesa":                    "DKK",
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
	"Guaraní paraguayo":                "PYG",
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
	"Rand surafricano":                 "ZAR",
	"Real Brasileño":                   "BRL",
	"Rial iraní":                       "IRR",
	"Rial saudita":                     "SAR",
	"Ringgit malasio":                  "MYR",
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

func GetCurrencies() ([]Currency, error) {
	url, err := GetCurrenciesDailyURL()
	if err != nil {
		return nil, err
	}

	timeoutContext, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, err := http.NewRequestWithContext(
		timeoutContext,
		"GET",
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	var client = http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	currencies, err := ParseCurrenciesHTML(res.Body)
	if err != nil {
		return nil, err
	}

	return currencies, nil
}

func GetCurrenciesDailyURL() (string, error) {
	resp, err := http.Get("https://si3.bcentral.cl/Indicadoressiete/secure/IndicadoresDiarios.aspx")
	if err != nil {
		return "", errors.Wrap(err, "unable to execute request")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	var url string
	doc.Find("#hypLnk1_11").Each(func(i int, s *goquery.Selection) {
		if i != 0 {
			return
		}

		url, _ = s.Attr("href")
	})

	if url == "" {
		return "", errors.New("no url found")
	}

	return "https://si3.bcentral.cl/Indicadoressiete/secure/" + url, nil
}

func ParseCurrenciesHTML(r io.ReadCloser) (currencies []Currency, err error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		var name string
		var exchange string

		s.Find("td").Each(func(i int, s *goquery.Selection) {
			switch i {
			case 0:
				name = strings.TrimSpace(s.Text())
			case 1:
				exchange = s.Text()
			}
		})

		currencies = append(currencies, Currency{
			Name:         name,
			ISO4217:      currencyNameToISO4217(name),
			ExchangeRate: parseChileanFloat(exchange),
		})
	})

	return currencies, nil
}

func searchCurrency(currencies []Currency, patterns string) []Currency {
	var hints []string
	for _, currency := range currencies {
		hints = append(hints, fmt.Sprintf("%s %s", removeTilde(currency.Name), currency.ISO4217))
	}

	matches := fuzzy.Find(removeTilde(patterns), hints)

	var results []Currency
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

func currencyNameToISO4217(name string) string {
	iso4217, _ := ISO4217Dict[name]
	return iso4217
}
