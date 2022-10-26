package rut

import (
	"github.com/ccuetoh/libreapi/pkg/env"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ValidateHandler(env *env.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		rutStr := c.Query("rut")
		if rutStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"errors": gin.H{
					"rut": "no rut was provided",
				},
			})

			env.Log(c).Trace("no rut")
			return
		}

		rut, err := ParseRUT(rutStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"errors": gin.H{
					"rut": "the provided rut is invalid",
				},
			})

			env.Log(c).Tracef("invalid rut: %v", err)
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

		env.Log(c).Trace("ok")
	}
}

func DigitHandler(env *env.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		rutStr := c.Query("rut")
		if rutStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"errors": gin.H{
					"rut": "no rut was provided",
				},
			})

			env.Log(c).Trace("no rut")
			return
		}

		rut, err := ParseRUT(rutStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"errors": gin.H{
					"rut": "the provided rut is invalid",
				},
			})

			env.Log(c).Tracef("invalid rut: %v", err)
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

		env.Log(c).Trace("ok")
	}
}

func SIIActivityHandler(env *env.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		rutStr := c.Query("rut")
		if rutStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"errors": gin.H{
					"rut": "no rut was provided",
				},
			})

			env.Log(c).Trace("no rut")
			return
		}

		rut, err := ParseRUT(rutStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"errors": gin.H{
					"rut": "the provided rut is invalid",
				},
			})

			env.Log(c).Tracef("invalid rut: %v", err)
			return
		}

		details, err := GetSIIDetails(rut)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"errors": gin.H{
					"fetch error": "unable to fetch the data",
				},
			})

			env.Log(c).Errorf("unable to fetch data: %v", err)
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

		env.Log(c).Trace("ok")
	}
}
