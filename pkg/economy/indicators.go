package economy

import (
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

type Indicators struct {
	UF        float64 `json:"uf"`
	IVP       float64 `json:"ivp"`
	Dollar    float64 `json:"dollar"`
	Euro      float64 `json:"euro"`
	ITCNM     float64 `json:"itcnm"`
	OztSilver float64 `json:"ozt_silver"`
	OztGold   float64 `json:"ozt_gold"`
	LbCopper  float64 `json:"lb_copper"`
}

func parseIndicatorsHTML(r io.ReadCloser) (*Indicators, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create document")
	}

	indicators := &Indicators{}
	indicators.UF, err = parseCurrency(doc.Find("label#lblValor1_1").Text())
	if err != nil {
		return nil, errors.Wrap(err, "invalid value")
	}

	indicators.IVP, err = parseCurrency(doc.Find("label#lblValor1_2").Text())
	if err != nil {
		return nil, errors.Wrap(err, "invalid value")
	}

	indicators.Dollar, err = parseCurrency(doc.Find("label#lblValor1_3").Text())
	if err != nil {
		return nil, errors.Wrap(err, "invalid value")
	}

	// 4 is skipped
	indicators.Euro, err = parseCurrency(doc.Find("label#lblValor1_5").Text())
	if err != nil {
		return nil, errors.Wrap(err, "invalid value")
	}

	// 6 is skipped
	indicators.ITCNM, err = parseCurrency(doc.Find("label#lblValor1_7").Text())
	if err != nil {
		return nil, errors.Wrap(err, "invalid value")
	}

	indicators.OztGold, err = parseCurrency(doc.Find("label#lblValor2_3").Text())
	if err != nil {
		return nil, errors.Wrap(err, "invalid value")
	}

	indicators.OztSilver, err = parseCurrency(doc.Find("label#lblValor2_4").Text())
	if err != nil {
		return nil, errors.Wrap(err, "invalid value")
	}

	indicators.LbCopper, err = parseCurrency(doc.Find("label#lblValor2_5").Text())
	if err != nil {
		return nil, errors.Wrap(err, "invalid value")
	}

	return indicators, nil
}

func parseCurrency(s string) (float64, error) {
	return strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(s, ".", ""), ",", "."), 64)
}
