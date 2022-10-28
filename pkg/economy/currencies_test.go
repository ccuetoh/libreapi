package economy

import (
	"testing"

	"github.com/ccuetoh/libreapi/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestParseCurrenciesHTMLOk(t *testing.T) {
	page, err := test.LoadHTML("currencies_ok")
	if err != nil {
		t.Fatalf("unable to load test case html: %v", err)
	}

	var expected []*Currency
	err = test.LoadJSON("currencies_ok", &expected)
	if err != nil {
		t.Fatalf("unable to load test case json: %v", err)
	}

	got, err := parseCurrenciesHTML(page)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestParseCurrenciesHTMLBadValue(t *testing.T) {
	page, err := test.LoadHTML("currencies_bad_value")
	if err != nil {
		t.Fatalf("unable to load test case html: %v", err)
	}

	got, err := parseCurrenciesHTML(page)
	assert.Error(t, err)
	assert.Equal(t, ([]*Currency)(nil), got)
}

func TestParseCurrenciesHTMLUnknown(t *testing.T) {
	page, err := test.LoadHTML("currencies_unknown")
	if err != nil {
		t.Fatalf("unable to load test case html: %v", err)
	}

	got, err := parseCurrenciesHTML(page)
	assert.Error(t, err)
	assert.Equal(t, ([]*Currency)(nil), got)
}

func TestParseCurrenciesHTMLInvalidReader(t *testing.T) {
	page, err := test.LoadHTML("currencies_unknown")
	if err != nil {
		t.Fatalf("unable to load test case html: %v", err)
	}

	page.Close()

	got, err := parseCurrenciesHTML(page)
	assert.Error(t, err)
	assert.Equal(t, ([]*Currency)(nil), got)
}

func TestFilterCurrencies(t *testing.T) {
	currencies := []*Currency{
		{
			Name:         "Euro",
			ISO4217:      "EUR",
			ExchangeRate: 1.0018,
		},
		{
			Name:         "Rupia india",
			ISO4217:      "INR",
			ExchangeRate: 82.4950,
		},
		{
			Name:         "Rupia pakistaní",
			ISO4217:      "PKR",
			ExchangeRate: 221.2003,
		},
	}

	assert.Equal(t, ([]*Currency)(nil), filterCurrencies(currencies, "Dolar"))
	assert.Equal(t, []*Currency{
		{
			Name:         "Euro",
			ISO4217:      "EUR",
			ExchangeRate: 1.0018,
		},
	}, filterCurrencies(currencies, "eur"))
	assert.Equal(t, []*Currency{
		{
			Name:         "Rupia india",
			ISO4217:      "INR",
			ExchangeRate: 82.4950,
		},
		{
			Name:         "Rupia pakistaní",
			ISO4217:      "PKR",
			ExchangeRate: 221.2003,
		},
	}, filterCurrencies(currencies, "rupia"))
}

func TestRemoveTilde(t *testing.T) {
	assert.Equal(t, "aeiou", removeTilde("áéíóú"))
	assert.Equal(t, "que", removeTilde("qué"))
	assert.Equal(t, "ahi", removeTilde("ahí"))
	assert.Equal(t, "a b c d e f g", removeTilde("á b ć d é f ǵ"))
}
