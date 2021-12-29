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
)

var DollarToCLP float64

type EconomicIndicators struct {
	UF        float64 `json:"uf"`
	IVP       float64 `json:"ivp"`
	Dollar    float64 `json:"dollar"`
	Euro      float64 `json:"euro"`
	ITCNM     float64 `json:"itcnm"`
	OztSilver float64 `json:"ozt_silver"`
	OztGold   float64 `json:"ozt_gold"`
	LbCopper  float64 `json:"lb_copper"`
}

func GetBancoCentralIndicators() (EconomicIndicators, error) {
	timeoutContext, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	req, err := http.NewRequestWithContext(
		timeoutContext,
		"GET",
		"https://si3.bcentral.cl/Indicadoressiete/secure/Indicadoresdiarios.aspx",
		nil,
	)
	if err != nil {
		return EconomicIndicators{}, err
	}

	var client = http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return EconomicIndicators{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return EconomicIndicators{}, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	indicators, err := ParseIndicatorsHTML(res.Body)
	if err != nil {
		return EconomicIndicators{}, err
	}

	DollarToCLP = indicators.Dollar

	return indicators, nil
}

func ParseIndicatorsHTML(r io.ReadCloser) (EconomicIndicators, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return EconomicIndicators{}, err
	}

	var indicators EconomicIndicators
	doc.Find("label#lblValor1_1").Each(func(i int, selection *goquery.Selection) {
		if i != 0 {
			return
		}

		indicators.UF = parseChileanFloat(selection.Text())
	})

	doc.Find("label#lblValor1_2").Each(func(i int, selection *goquery.Selection) {
		if i != 0 {
			return
		}

		indicators.IVP = parseChileanFloat(selection.Text())
	})

	doc.Find("label#lblValor1_3").Each(func(i int, selection *goquery.Selection) {
		if i != 0 {
			return
		}

		indicators.Dollar = parseChileanFloat(selection.Text())
	})

	// For some reason 4 was skipped
	doc.Find("label#lblValor1_5").Each(func(i int, selection *goquery.Selection) {
		if i != 0 {
			return
		}

		indicators.Euro = parseChileanFloat(selection.Text())
	})

	doc.Find("label#lblValor1_5").Each(func(i int, selection *goquery.Selection) {
		if i != 0 {
			return
		}

		indicators.ITCNM = parseChileanFloat(selection.Text())
	})

	doc.Find("label#lblValor2_3").Each(func(i int, selection *goquery.Selection) {
		if i != 0 {
			return
		}

		indicators.OztGold = parseChileanFloat(selection.Text())
	})

	doc.Find("label#lblValor2_4").Each(func(i int, selection *goquery.Selection) {
		if i != 0 {
			return
		}

		indicators.OztSilver = parseChileanFloat(selection.Text())
	})

	doc.Find("label#lblValor2_5").Each(func(i int, selection *goquery.Selection) {
		if i != 0 {
			return
		}

		indicators.LbCopper = parseChileanFloat(selection.Text())
	})

	return indicators, nil
}

func parseChileanFloat(floatStr string) float64 {
	floatStr = strings.ReplaceAll(floatStr, ".", "")
	floatStr = strings.ReplaceAll(floatStr, ",", ".")
	res, _ := strconv.ParseFloat(floatStr, 64)

	return res
}

func usdToCLP(usd float64) float64 {
	if DollarToCLP == 0 {
		_, err := GetBancoCentralIndicators()
		if err != nil {
			return 0
		}
	}

	return DollarToCLP * usd
}
