package economy

import (
	"github.com/ccuetoh/libreapi/pkg/env"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Service interface {
	GetIndicators() (*Indicators, error)
	GetCurrencies() ([]*Currency, error)
}

type Handler struct {
	env     *env.Env
	service Service
}

func NewHandler(env *env.Env, service Service) *Handler {
	return &Handler{
		env:     env,
		service: service,
	}
}

func (h *Handler) Indicators() gin.HandlerFunc {
	return func(c *gin.Context) {
		indicators, err := h.service.GetIndicators()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "unable to get data",
			})

			h.env.Log(c).Errorf("unable to fecth data: %v", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   indicators,
		})

		h.env.Log(c).Trace("ok")
	}
}

func (h *Handler) Currencies() gin.HandlerFunc {
	return func(c *gin.Context) {
		currencies, err := h.service.GetCurrencies()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "unable to get data",
			})

			h.env.Log(c).Errorf("unable to fecth data: %v", err)
			return
		}

		filter := c.Query("name")
		if filter == "" {
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data":   currencies,
			})

			h.env.Log(c).Trace("ok")
			return
		}

		currencies = filterCurrencies(currencies, filter)
		if len(currencies) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"status": "success",
				"data":   nil,
			})

			h.env.Log(c).Trace("ok (none matched)")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   currencies,
		})

		h.env.Log(c).Trace("ok")
		return
	}
}
