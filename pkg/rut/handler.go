package rut

import (
	"github.com/ccuetoh/libreapi/pkg/env"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Service interface {
	GetProfile(rut RUT) (*SIIProfile, error)
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

func (h *Handler) Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		rutStr := c.Query("rut")
		if rutStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"data": gin.H{
					"rut": "no rut was provided",
				},
			})

			h.env.Log(c).Trace("no rut")
			return
		}

		rut, err := ParseRUT(rutStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"data": gin.H{
					"rut": "the provided rut is invalid",
				},
			})

			h.env.Log(c).Tracef("invalid rut: %v", err)
			return
		}

		valid := rut.IsValid()
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"valid": valid,
				"rut":   rut.PrettyString(),
			},
		})

		h.env.Log(c).Trace("ok")
	}
}

func (h *Handler) VD() gin.HandlerFunc {
	return func(c *gin.Context) {
		rutStr := c.Query("rut")
		if rutStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"data": gin.H{
					"rut": "no rut was provided",
				},
			})

			h.env.Log(c).Trace("no rut")
			return
		}

		rut, err := ParseRUT(rutStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"data": gin.H{
					"rut": "the provided rut is invalid",
				},
			})

			h.env.Log(c).Tracef("invalid rut: %v", err)
			return
		}

		digit := rut.CalculateVD(false)
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"digit": VDToString(digit),
				"rut":   append(rut, digit).PrettyString(),
			},
		})

		h.env.Log(c).Trace("ok")
	}
}

func (h *Handler) Activity() gin.HandlerFunc {
	return func(c *gin.Context) {
		rutStr := c.Query("rut")
		if rutStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"data": gin.H{
					"rut": "no rut was provided",
				},
			})

			h.env.Log(c).Trace("no rut")
			return
		}

		rut, err := ParseRUT(rutStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"data": gin.H{
					"rut": "the provided rut is invalid",
				},
			})

			h.env.Log(c).Tracef("invalid rut: %v", err)
			return
		}

		details, err := h.service.GetProfile(rut)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "unable to fetch the data",
			})

			h.env.Log(c).Errorf("unable to fetch data: %v", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"rut":        rut.PrettyString(),
				"name":       details.Name,
				"activities": details.Activities,
			},
		})

		h.env.Log(c).Trace("ok")
	}
}
