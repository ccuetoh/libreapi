package weather

import (
	"github.com/ccuetoh/libreapi/internal/test"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ccuetoh/libreapi/pkg/env"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockService struct {
	stations    []*ClimateStation
	stationsErr error
}

func (s MockService) GetClimateStations() ([]*ClimateStation, error) {
	return s.stations, s.stationsErr
}

func TestStationsOk(t *testing.T) {
	gin.SetMode(gin.TestMode)

	today, _ := time.Parse("02-01-2006", "29-10-2020")
	today = today.UTC()

	yesterday := today.AddDate(0, 0, -1).UTC()

	data := []*ClimateStation{
		{
			Code:        1,
			Name:        "test1",
			Operational: true,
			LastReport:  &today,
			Temperature: 10,
			Humidity:    20,
			PressureHPA: 30,
			Today: &ClimateReport{
				Maximum: Measurement{
					Time:        today,
					Temperature: 40,
				},
				Minimum: Measurement{
					Time:        today,
					Temperature: 50,
				},
				Precipitations: Precipitations{
					Sum: 60,
					EMA: 70,
				},
			},
			Yesterday: &ClimateReport{
				Maximum: Measurement{
					Time:        yesterday,
					Temperature: 80,
				},
				Minimum: Measurement{
					Time:        yesterday,
					Temperature: 90,
				},
				Precipitations: Precipitations{
					Sum: 100,
					EMA: 110,
				},
			},
		},
	}

	service := MockService{stations: data}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	handler.Stations()(ctx)

	assert.Equal(t, recorder.Code, http.StatusOK)
}

func TestStationsError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := MockService{stationsErr: errors.New("server is on fire")}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	handler.Stations()(ctx)

	assert.Equal(t, recorder.Code, http.StatusInternalServerError)
	test.AssertResponseBodySlice(t, recorder, nil)
}

func TestStationsInvalidQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	today, _ := time.Parse("02-01-2006", "29-10-2020")
	today = today.UTC()

	yesterday := today.AddDate(0, 0, -1).UTC()

	data := []*ClimateStation{
		{
			Code:        1,
			Name:        "test1",
			Operational: true,
			LastReport:  &today,
			Temperature: 10,
			Humidity:    20,
			PressureHPA: 30,
			Today: &ClimateReport{
				Maximum: Measurement{
					Time:        today,
					Temperature: 40,
				},
				Minimum: Measurement{
					Time:        today,
					Temperature: 50,
				},
				Precipitations: Precipitations{
					Sum: 60,
					EMA: 70,
				},
			},
			Yesterday: &ClimateReport{
				Maximum: Measurement{
					Time:        yesterday,
					Temperature: 80,
				},
				Minimum: Measurement{
					Time:        yesterday,
					Temperature: 90,
				},
				Precipitations: Precipitations{
					Sum: 100,
					EMA: 110,
				},
			},
		},
	}

	service := MockService{stations: data}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = &http.Request{}
	ctx.Request.URL, _ = url.Parse("?name=test&code=1")

	handler.Stations()(ctx)

	assert.Equal(t, recorder.Code, http.StatusBadRequest)
	test.AssertResponseBodySlice(t, recorder, nil)
}

func TestStationsNameQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	today, _ := time.Parse("02-01-2006", "29-10-2020")
	today = today.UTC()

	yesterday := today.AddDate(0, 0, -1).UTC()

	data := []*ClimateStation{
		{
			Code:        1,
			Name:        "test1",
			Operational: true,
			LastReport:  &today,
			Temperature: 10,
			Humidity:    20,
			PressureHPA: 30,
			Today: &ClimateReport{
				Maximum: Measurement{
					Time:        today,
					Temperature: 40,
				},
				Minimum: Measurement{
					Time:        today,
					Temperature: 50,
				},
				Precipitations: Precipitations{
					Sum: 60,
					EMA: 70,
				},
			},
			Yesterday: &ClimateReport{
				Maximum: Measurement{
					Time:        yesterday,
					Temperature: 80,
				},
				Minimum: Measurement{
					Time:        yesterday,
					Temperature: 90,
				},
				Precipitations: Precipitations{
					Sum: 100,
					EMA: 110,
				},
			},
		},
	}

	service := MockService{stations: data}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = &http.Request{}
	ctx.Request.URL, _ = url.Parse("?name=test")

	handler.Stations()(ctx)

	assert.Equal(t, recorder.Code, http.StatusOK)
}

func TestStationsNameQueryNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	today, _ := time.Parse("02-01-2006", "29-10-2020")
	today = today.UTC()

	yesterday := today.AddDate(0, 0, -1).UTC()

	data := []*ClimateStation{
		{
			Code:        1,
			Name:        "test1",
			Operational: true,
			LastReport:  &today,
			Temperature: 10,
			Humidity:    20,
			PressureHPA: 30,
			Today: &ClimateReport{
				Maximum: Measurement{
					Time:        today,
					Temperature: 40,
				},
				Minimum: Measurement{
					Time:        today,
					Temperature: 50,
				},
				Precipitations: Precipitations{
					Sum: 60,
					EMA: 70,
				},
			},
			Yesterday: &ClimateReport{
				Maximum: Measurement{
					Time:        yesterday,
					Temperature: 80,
				},
				Minimum: Measurement{
					Time:        yesterday,
					Temperature: 90,
				},
				Precipitations: Precipitations{
					Sum: 100,
					EMA: 110,
				},
			},
		},
	}

	service := MockService{stations: data}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = &http.Request{}
	ctx.Request.URL, _ = url.Parse("?name=santiago")

	handler.Stations()(ctx)

	assert.Equal(t, recorder.Code, http.StatusNotFound)
	test.AssertResponseBodySlice(t, recorder, nil)
}

func TestStationsCodeQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	today, _ := time.Parse("02-01-2006", "29-10-2020")
	today = today.UTC()

	yesterday := today.AddDate(0, 0, -1).UTC()

	data := []*ClimateStation{
		{
			Code:        1,
			Name:        "test1",
			Operational: true,
			LastReport:  &today,
			Temperature: 10,
			Humidity:    20,
			PressureHPA: 30,
			Today: &ClimateReport{
				Maximum: Measurement{
					Time:        today,
					Temperature: 40,
				},
				Minimum: Measurement{
					Time:        today,
					Temperature: 50,
				},
				Precipitations: Precipitations{
					Sum: 60,
					EMA: 70,
				},
			},
			Yesterday: &ClimateReport{
				Maximum: Measurement{
					Time:        yesterday,
					Temperature: 80,
				},
				Minimum: Measurement{
					Time:        yesterday,
					Temperature: 90,
				},
				Precipitations: Precipitations{
					Sum: 100,
					EMA: 110,
				},
			},
		},
	}

	service := MockService{stations: data}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = &http.Request{}
	ctx.Request.URL, _ = url.Parse("?code=1")

	handler.Stations()(ctx)

	assert.Equal(t, recorder.Code, http.StatusOK)
}

func TestStationsCodeQueryNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	today, _ := time.Parse("02-01-2006", "29-10-2020")
	today = today.UTC()

	yesterday := today.AddDate(0, 0, -1).UTC()

	data := []*ClimateStation{
		{
			Code:        1,
			Name:        "test1",
			Operational: true,
			LastReport:  &today,
			Temperature: 10,
			Humidity:    20,
			PressureHPA: 30,
			Today: &ClimateReport{
				Maximum: Measurement{
					Time:        today,
					Temperature: 40,
				},
				Minimum: Measurement{
					Time:        today,
					Temperature: 50,
				},
				Precipitations: Precipitations{
					Sum: 60,
					EMA: 70,
				},
			},
			Yesterday: &ClimateReport{
				Maximum: Measurement{
					Time:        yesterday,
					Temperature: 80,
				},
				Minimum: Measurement{
					Time:        yesterday,
					Temperature: 90,
				},
				Precipitations: Precipitations{
					Sum: 100,
					EMA: 110,
				},
			},
		},
	}

	service := MockService{stations: data}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = &http.Request{}
	ctx.Request.URL, _ = url.Parse("?code=2")

	handler.Stations()(ctx)

	assert.Equal(t, recorder.Code, http.StatusNotFound)
	test.AssertResponseBodySlice(t, recorder, nil)
}
