package rut

import (
	"testing"

	"github.com/ccuetoh/libreapi/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestGetProfile(t *testing.T) {
	service := NewDefaultService()

	rut, err := parseRUT("3.632.455-4", false)
	assert.NoError(t, err)

	profile, err := service.GetProfile(rut)
	assert.NoError(t, err)
	assert.NotNil(t, profile)
}

func TestParseProfileHTMLOk(t *testing.T) {
	page, err := test.LoadHTML("activities_ok")
	if err != nil {
		t.Fatalf("unable to load test case html: %v", err)
	}

	var expected *SIIProfile
	err = test.LoadJSON("activities_ok", &expected)
	if err != nil {
		t.Fatalf("unable to load test case json: %v", err)
	}

	got, err := parseActivitiesHTML(page)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestParseProfileHTMLBadDate(t *testing.T) {
	page, err := test.LoadHTML("activities_bad_date")
	if err != nil {
		t.Fatalf("unable to load test case html: %v", err)
	}

	got, err := parseActivitiesHTML(page)
	assert.Error(t, err)
	assert.Equal(t, (*SIIProfile)(nil), got)
}

func TestParseProfileHTMLBadCode(t *testing.T) {
	page, err := test.LoadHTML("activities_bad_code")
	if err != nil {
		t.Fatalf("unable to load test case html: %v", err)
	}

	got, err := parseActivitiesHTML(page)
	assert.Error(t, err)
	assert.Equal(t, (*SIIProfile)(nil), got)
}

func TestParseCurrenciesHTMLInvalidReader(t *testing.T) {
	page, err := test.LoadHTML("activities_bad_date")
	if err != nil {
		t.Fatalf("unable to load test case html: %v", err)
	}

	page.Close()

	got, err := parseActivitiesHTML(page)
	assert.Error(t, err)
	assert.Equal(t, (*SIIProfile)(nil), got)
}
