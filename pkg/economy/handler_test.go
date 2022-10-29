package economy

import (
	"github.com/ccuetoh/libreapi/internal/test"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ccuetoh/libreapi/pkg/env"
)

type MockService struct {
	indicators    *Indicators
	indicatorsErr error
	currencies    []*Currency
	currenciesErr error
}

func (s MockService) GetIndicators() (*Indicators, error) {
	return s.indicators, s.indicatorsErr
}

func (s MockService) GetCurrencies() ([]*Currency, error) {
	return s.currencies, s.currenciesErr
}

func TestIndicatorsOk(t *testing.T) {
	gin.SetMode(gin.TestMode)

	data := &Indicators{
		UF:        1,
		IVP:       2,
		Dollar:    3,
		Euro:      4,
		ITCNM:     5,
		OztSilver: 6,
		OztGold:   7,
		LbCopper:  8,
	}

	service := MockService{
		indicators: data,
	}

	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	handler.Indicators()(ctx)

	assert.Equal(t, recorder.Code, http.StatusOK)
	test.AssertResponseBody(t, recorder, data)
}

func TestIndicatorsError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := MockService{
		indicatorsErr: errors.New("server is on fire"),
	}

	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	handler.Indicators()(ctx)

	assert.Equal(t, recorder.Code, http.StatusInternalServerError)
}

func TestCurrenciesOk(t *testing.T) {
	gin.SetMode(gin.TestMode)

	data := []*Currency{
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

	service := MockService{
		currencies: data,
	}

	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	handler.Currencies()(ctx)

	assert.Equal(t, recorder.Code, http.StatusOK)
	test.AssertResponseBodySlice(t, recorder, data)
}

func TestCurrenciesFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	data := []*Currency{
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

	service := MockService{
		currencies: data,
	}

	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = &http.Request{}
	ctx.Request.URL, _ = url.Parse("?name=euro")

	handler.Currencies()(ctx)

	assert.Equal(t, recorder.Code, http.StatusOK)
	test.AssertResponseBodySlice(t, recorder, []*Currency{
		{
			Name:         "Euro",
			ISO4217:      "EUR",
			ExchangeRate: 1.0018,
		},
	})
}

func TestCurrenciesFilterNoneMatched(t *testing.T) {
	gin.SetMode(gin.TestMode)

	data := []*Currency{
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

	service := MockService{
		currencies: data,
	}

	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = &http.Request{}
	ctx.Request.URL, _ = url.Parse("?name=dolar")

	handler.Currencies()(ctx)

	assert.Equal(t, recorder.Code, http.StatusNotFound)
	test.AssertResponseBodySlice(t, recorder, nil)
}

func TestCurrenciesError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := MockService{
		currenciesErr: errors.New("server is on fire"),
	}

	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	handler.Currencies()(ctx)

	assert.Equal(t, recorder.Code, http.StatusInternalServerError)
}
