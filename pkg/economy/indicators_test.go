package economy

import (
	"testing"

	"github.com/ccuetoh/libreapi/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestParseIndicatorsHTML(t *testing.T) {
	page, err := test.LoadHTML("indicators_ok")
	if err != nil {
		t.Fatalf("unable to load test case html: %v", err)
	}

	expected := &Indicators{}
	err = test.LoadJSON("indicators_ok", expected)
	if err != nil {
		t.Fatalf("unable to load test case json: %v", err)
	}

	got, err := parseIndicatorsHTML(page)

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestParseIndicatorsHTMLInvalidReader(t *testing.T) {
	page, err := test.LoadHTML("currencies_ok")
	if err != nil {
		t.Fatalf("unable to load test case html: %v", err)
	}

	page.Close()

	got, err := parseIndicatorsHTML(page)
	assert.Error(t, err)
	assert.Equal(t, (*Indicators)(nil), got)
}
