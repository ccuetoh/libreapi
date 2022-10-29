package weather

import (
	"github.com/ccuetoh/libreapi/pkg/env"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Service interface {
	GetClimateStations() ([]*ClimateStation, error)
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

func (h *Handler) Stations() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		code := c.Query("code")

		if name != "" && code != "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"errors": gin.H{
					"conflict": "can't search both name and code at the same time",
				},
			})

			h.env.Log(c).Trace("both name and code")
			return
		}

		stations, err := h.service.GetClimateStations()
		if err != nil {
			c.JSON(500, gin.H{
				"status": "error",
				"errors": gin.H{
					"fetch": "unable to fetch data",
				},
			})

			h.env.Log(c).Errorf("unable to fetch data: %v", err)
			return
		}

		if name != "" {
			stationsFound := searchStationName(stations, name)
			if len(stationsFound) == 0 {
				c.JSON(http.StatusNotFound, gin.H{
					"status": "success",
					"data":   nil,
				})

				h.env.Log(c).Trace("ok (none)")
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data":   stationsFound,
			})

			h.env.Log(c).Trace("ok")
			return
		}

		if code != "" {
			match, found := searchStationCode(stations, code)
			if !found {
				c.JSON(http.StatusNotFound, gin.H{
					"status": "success",
					"data":   nil,
				})

				h.env.Log(c).Trace("ok (none)")
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data":   match,
			})

			h.env.Log(c).Trace("ok")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   stations,
		})

		h.env.Log(c).Trace("ok")
	}
}
