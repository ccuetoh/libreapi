package weather

import (
	"strings"
	"testing"
	"time"

	"github.com/ccuetoh/libreapi/internal/test"

	"github.com/stretchr/testify/assert"
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

	// WARNING: Times are ignored since parsing at different locations introduces inconsistencies
	// TODO: Add dedicated test for time
	emptyTime := time.Time{}

	for _, station := range expected {
		if !station.Operational {
			continue
		}

		// Has a bugged character in the website. Invalid for both Windows-1254 and Unicode
		if strings.Contains(station.Name, "El Huertón liceo agrícola") {
			station.Name = "El Huertón liceo agrícola"
		}

		last := emptyTime
		station.LastReport = &last
		station.Today.Minimum.Time = emptyTime
		station.Today.Maximum.Time = emptyTime
		station.Yesterday.Minimum.Time = emptyTime
		station.Yesterday.Maximum.Time = emptyTime
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

		last := emptyTime
		station.LastReport = &last
		station.Today.Minimum.Time = emptyTime
		station.Today.Maximum.Time = emptyTime
		station.Yesterday.Minimum.Time = emptyTime
		station.Yesterday.Maximum.Time = emptyTime
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
