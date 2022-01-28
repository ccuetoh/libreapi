package economy

import (
	"context"
	"fmt"
	"io/ioutil"

	"io"
	"net/http"
	"time"

	"github.com/json-iterator/go"
	"github.com/sahilm/fuzzy"
)

type Coin struct {
	Name         string  `json:"name"`
	Symbol       string  `json:"symbol"`
	MarketCap    int64   `json:"market_cap_usd"`
	PriceUSD     float64 `json:"price_usd"`
	PriceCLP     float64 `json:"price_clp"`
	Supply       int64   `json:"supply"`
	Volume       int64   `json:"volume_usd"`
	HourlyChange float64 `json:"hourly_change"`
	DailyChange  float64 `json:"daily_change"`
	WeeklyChange float64 `json:"weekly_change"`
}

func GetCrypto() ([]Coin, error) {
	timeoutContext, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, err := http.NewRequestWithContext(
		timeoutContext,
		"GET",
		"https://api.coinmarketcap.com/data-api/v3/cryptocurrency/listing?cryptoType=all&tagType=all&limit=1000",
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

	coins, err := ParseCryptoJSON(res.Body)
	if err != nil {
		return nil, err
	}

	return coins, nil
}

func ParseCryptoJSON(r io.ReadCloser) (coins []Coin, err error) {
	dataBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	type ResponseData struct {
		Data struct {
			CryptoCurrencyList []struct {
				Name        string  `json:"name"`
				Symbol      string  `json:"symbol"`
				IsActive    int     `json:"isActive"`
				TotalSupply float64 `json:"totalSupply"`
				Quotes      []struct {
					Price            float64 `json:"price"`
					Volume24H        float64 `json:"volume24h"`
					MarketCap        float64 `json:"marketCap"`
					PercentChange1H  float64 `json:"percentChange1h"`
					PercentChange24H float64 `json:"percentChange24h"`
					PercentChange7D  float64 `json:"percentChange7d"`
					PercentChange30D float64 `json:"percentChange30d"`
				} `json:"quotes"`
			} `json:"cryptoCurrencyList"`
		} `json:"data"`
	}

	var data ResponseData
	err = jsoniter.Unmarshal(dataBytes, &data)
	if err != nil {
		return nil, err
	}

	for _, coin := range data.Data.CryptoCurrencyList {
		if coin.IsActive != 1 || len(coin.Quotes) < 1 {
			continue
		}

		coins = append(coins, Coin{
			Name:         coin.Name,
			Symbol:       coin.Symbol,
			MarketCap:    int64(coin.Quotes[0].MarketCap),
			PriceUSD:     coin.Quotes[0].Price,
			PriceCLP:     usdToCLP(coin.Quotes[0].Price),
			Supply:       int64(coin.TotalSupply),
			Volume:       int64(coin.Quotes[0].Volume24H),
			HourlyChange: coin.Quotes[0].PercentChange1H,
			DailyChange:  coin.Quotes[0].PercentChange24H,
			WeeklyChange: coin.Quotes[0].PercentChange7D,
		})
	}

	return coins, nil
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
