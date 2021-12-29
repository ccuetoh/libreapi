package economy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sahilm/fuzzy"
)

type Coin struct {
	Name          string  `json:"name"`
	Symbol        string  `json:"symbol"`
	MarketCap     int64   `json:"market_cap_usd"`
	PriceUSD      float64 `json:"price_usd"`
	PriceCLP      float64 `json:"price_clp"`
	Supply        int64   `json:"supply"`
	Volume        int64   `json:"volume_usd"`
	HourlyChange  float64 `json:"hourly_change"`
	DailyChange   float64 `json:"daily_change"`
	WeekleyChange float64 `json:"weekley_change"`
}

func GetCrypto() ([]Coin, error) {
	timeoutContext, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, err := http.NewRequestWithContext(
		timeoutContext,
		"GET",
		"https://coinmarketcap.com/all/views/all/",
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

	coins, err := ParseCryptoHTML(res.Body)
	if err != nil {
		return nil, err
	}

	return coins, nil
}

func ParseCryptoHTML(r io.ReadCloser) (coins []Coin, err error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		if i < 3 {
			// Header
			return
		}

		var name = s.Find(".cmc-table__column-name--name").Text()
		if name == "" {
			return
		}

		var symbol = s.Find(".cmc-table__cell--sort-by__symbol").Text()
		var marketCap = s.Find(".cmc-table__cell--sort-by__market-cap").Text()
		var priceUSD = s.Find(".cmc-table__cell--sort-by__price").Text()
		var supply = s.Find(".cmc-table__cell--sort-by__circulating-supply").Text()
		var volume = s.Find(".cmc-table__cell--sort-by__volume-24-h").Text()
		var hourChange = s.Find(".cmc-table__cell--sort-by__percent-change-1-h").Text()
		var dayChange = s.Find(".cmc-table__cell--sort-by__percent-change-24-h").Text()
		var weekChange = s.Find(".cmc-table__cell--sort-by__percent-change-7-d").Text()

		coins = append(coins, Coin{
			Name:          name,
			Symbol:        symbol,
			MarketCap:     usdPriceToInt64(marketCap),
			PriceUSD:      usdPriceToFloat64(priceUSD),
			PriceCLP:      usdToCLP(usdPriceToFloat64(priceUSD)),
			Supply:        supplyToInt64(supply),
			Volume:        usdPriceToInt64(volume),
			HourlyChange:  percentageVariationToFloa64(hourChange),
			DailyChange:   percentageVariationToFloa64(dayChange),
			WeekleyChange: percentageVariationToFloa64(weekChange),
		})
	})

	return coins, nil
}

func usdPriceToInt64(price string) int64 {
	price = strings.ReplaceAll(price, "$", "")
	price = strings.ReplaceAll(price, ",", "")
	price = strings.Split(price, ".")[0]

	i, _ := strconv.ParseInt(price, 10, 64)
	return i
}

func usdPriceToFloat64(price string) float64 {
	price = strings.ReplaceAll(price, "$", "")
	price = strings.ReplaceAll(price, ",", "")

	f, _ := strconv.ParseFloat(price, 64)
	return f
}

func supplyToInt64(supply string) int64 {
	supply = strings.ReplaceAll(supply, ",", "")
	supply = strings.Split(supply, " ")[0]

	i, _ := strconv.ParseInt(supply, 10, 64)
	return i
}

func percentageVariationToFloa64(variation string) float64 {
	variation = strings.ReplaceAll(variation, "%", "")
	res, _ := strconv.ParseFloat(variation, 64)
	return res
}

func searchCoin(coins []Coin, patterns string) []Coin {
	var hints []string
	for _, coin := range coins {
		hints = append(hints, fmt.Sprintf("%s %s", coin.Name, coin.Symbol))
	}

	matches := fuzzy.Find(patterns, hints)

	var results []Coin
	for _, match := range matches {
		results = append(results, coins[match.Index])
	}

	return results
}
