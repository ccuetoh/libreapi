package economy

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetIndicators(t *testing.T) {
	service := NewDefaultService()
	service.client.Timeout = time.Minute

	indicators, err := service.GetIndicators()
	assert.NoError(t, err)
	assert.NotNil(t, indicators)
}

func TestGetCurrencies(t *testing.T) {
	service := NewDefaultService()
	service.client.Timeout = time.Minute

	indicators, err := service.GetCurrencies()
	assert.NoError(t, err)
	assert.NotNil(t, indicators)
}
