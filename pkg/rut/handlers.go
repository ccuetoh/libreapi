package rut

import (
	"github.com/gin-gonic/gin"
)

func ValidateHandler(ctx *gin.Context) {
	rut := ctx.Query("rut")
	if rut == "" {
		ctx.JSON(400, gin.H{
			"status": "fail",
			"errors": gin.H{
				"rut": "no rut was provided",
			},
		})
		return
	}

	rut2, err := ParseRUT(rut)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": "fail",
			"errors": gin.H{
				"rut": "the provided rut is invalid",
			},
		})
		return
	}

	valid := rut2.IsValid()
	ctx.JSON(200, gin.H{
		"status": "success",
		"data": gin.H{
			"valid": valid,
			"rut":   rut2.PrettyString(),
		},
	})
}

func DigitHandler(ctx *gin.Context) {
	rut := ctx.Query("rut")
	if rut == "" {
		ctx.JSON(400, gin.H{
			"status": "fail",
			"errors": gin.H{
				"rut": "no rut was provided",
			},
		})
		return
	}

	rut2, err := ParseRUT(rut)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": "fail",
			"errors": gin.H{
				"rut": "the provided rut is invalid",
			},
		})
		return
	}

	digit := rut2.CalculateVD(false)
	ctx.JSON(200, gin.H{
		"status": "success",
		"data": gin.H{
			"digit": VDToString(digit),
			"rut":   append(rut2, digit).PrettyString(),
		},
	})
}

func SIIActivityHandler(ctx *gin.Context) {
	rut := ctx.Query("rut")
	if rut == "" {
		ctx.JSON(400, gin.H{
			"status": "fail",
			"errors": gin.H{
				"rut": "no rut was provided",
			},
		})
		return
	}

	rut2, err := ParseRUT(rut)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": "fail",
			"errors": gin.H{
				"rut": "the provided rut is invalid",
			},
		})
		return
	}

	details, err := GetSIIDetails(rut2)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": "error",
			"errors": gin.H{
				"fetch error": "unable to fetch the data",
			},
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status": "success",
		"data": gin.H{
			"rut":        rut2.PrettyString(),
			"name":       details.Name,
			"activities": details.Activities,
		},
	})
}
