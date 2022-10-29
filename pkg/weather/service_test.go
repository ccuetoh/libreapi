package weather

import (
	"github.com/ccuetoh/libreapi/internal/test"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestGetStations(t *testing.T) {
	service := NewDefaultService()
	service.client.Timeout = time.Minute

	stations, err := service.GetClimateStations()
	assert.NoError(t, err)
	assert.NotNil(t, stations)
}

func TestParseStationHTMLOk(t *testing.T) {
	page, err := test.LoadHTML("stations_ok")
	if err != nil {
		t.Fatalf("unable to load test case html: %v", err)
	}

	var expected []*ClimateStation
	err = test.LoadJSON("stations_ok", &expected)
	if err != nil {
		t.Fatalf("unable to load test case json: %v", err)
	}

	// Time location is parsed differently
	loc, _ := time.LoadLocation("America/Santiago")
	for _, station := range expected {
		if !station.Operational {
			continue
		}

		// Has a bugged character in the website. Invalid for both Windows-1254 and Unicode
		if strings.Contains(station.Name, "El Huertón liceo agrícola") {
			station.Name = "El Huertón liceo agrícola"
		}

		last := station.LastReport.In(loc)
		station.LastReport = &last
		station.Today.Minimum.Time = station.Today.Minimum.Time.In(loc)
		station.Today.Maximum.Time = station.Today.Maximum.Time.In(loc)
		station.Yesterday.Minimum.Time = station.Yesterday.Minimum.Time.In(loc)
		station.Yesterday.Maximum.Time = station.Yesterday.Maximum.Time.In(loc)
	}

	got, err := parseClimateHTML(page)
	assert.NoError(t, err)

	for _, station := range got {
		if !station.Operational {
			continue
		}

		// Has a bugged character in the website. Invalid for both Windows-1254 and Unicode
		if strings.Contains(station.Name, "El Huertón liceo agrícola") {
			station.Name = "El Huertón liceo agrícola"
		}

		last := station.LastReport.In(loc)
		station.LastReport = &last
		station.Today.Minimum.Time = station.Today.Minimum.Time.In(loc)
		station.Today.Maximum.Time = station.Today.Maximum.Time.In(loc)
		station.Yesterday.Minimum.Time = station.Yesterday.Minimum.Time.In(loc)
		station.Yesterday.Maximum.Time = station.Yesterday.Maximum.Time.In(loc)
	}

	assert.Equal(t, expected, got)
}

func TestParseStationsHTMLInvalidReader(t *testing.T) {
	page, err := test.LoadHTML("stations_ok")
	if err != nil {
		t.Fatalf("unable to load test case html: %v", err)
	}

	page.Close()

	got, err := parseClimateHTML(page)
	assert.Error(t, err)
	assert.Equal(t, ([]*ClimateStation)(nil), got)
}
