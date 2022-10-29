package rut

import (
	"fmt"
	"github.com/ccuetoh/libreapi/pkg/env"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
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

		rut, err := parseRUT(rutStr, false)
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
				"rut":   rut.String(),
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

		rut, err := parseRUT(rutStr, true)
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

		vd := rut.calculateVD()
		rut.VD = vd

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"digit": vd.String(),
				"rut":   rut.String(),
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

		rut, err := parseRUT(rutStr, false)
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

		profile, err := h.service.GetProfile(rut)
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
				"rut":        rut.String(),
				"name":       profile.Name,
				"activities": profile.Activities,
			},
		})

		h.env.Log(c).Trace("ok")
	}
}

func (h *Handler) Generate() gin.HandlerFunc {
	return func(c *gin.Context) {
		min := 500000
		max := 25000000

		var err error
		minParam := c.Query("min")
		if minParam != "" {
			min, err = strconv.Atoi(minParam)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status": "fail",
					"data": gin.H{
						"min": "min must be numeric",
					},
				})

				h.env.Log(c).Trace("bad min")
				return
			}
		}

		maxParam := c.Query("max")
		if maxParam != "" {
			max, err = strconv.Atoi(maxParam)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status": "fail",
					"data": gin.H{
						"max": "max must be numeric",
					},
				})

				h.env.Log(c).Trace("bad max")
				return
			}
		}

		if min >= max {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"data": gin.H{
					"range": fmt.Sprintf("min must be lower than max (%d)", max),
				},
			})

			h.env.Log(c).Trace("invalid range")
			return
		}

		rut, _ := generateRUT(min, max)

		var digits strings.Builder
		for _, d := range rut.Digits {
			digits.WriteString(strconv.Itoa(int(d)))
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"rut":    rut.String(),
				"vd":     rut.VD.String(),
				"digits": digits.String(),
			},
		})

		h.env.Log(c).Trace("ok")
	}
}
