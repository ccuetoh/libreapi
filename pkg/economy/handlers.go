package economy

import "github.com/gin-gonic/gin"

func BancoCentraIndicatorsHandler(ctx *gin.Context) {
	indicators, err := GetBancoCentralIndicators()
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
		"data": indicators,
	})
}

func CryptoHandler(ctx *gin.Context) {
	coins, err := GetCrypto()
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": "error",
			"errors": gin.H{
				"fetch error": "unable to fetch the data",
			},
		})
		return
	}

	coinName := ctx.Query("name")
	if coinName != "" {
		coinsFound := searchCoin(coins, coinName)
		if len(coinsFound) == 0 {
			ctx.JSON(404, gin.H{
				"status": "success",
				"data": nil,
			})
			return
		}

		ctx.JSON(200, gin.H{
			"status": "success",
			"data": coinsFound,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status": "success",
		"data": coins,
	})
}

func CurrencyHandler(ctx *gin.Context) {
	currencies, err := GetCurrencies()
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": "error",
			"errors": gin.H{
				"fetch error": "unable to fetch the data",
			},
		})
		return
	}

	currencyName := ctx.Query("name")
	if currencyName != "" {
		currenciesFound := searchCurrency(currencies, currencyName)
		if len(currenciesFound) == 0 {
			ctx.JSON(404, gin.H{
				"status": "success",
				"data": nil,
			})
			return
		}

		ctx.JSON(200, gin.H{
			"status": "success",
			"data":   currenciesFound,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status": "success",
		"data":   currencies,
	})
}