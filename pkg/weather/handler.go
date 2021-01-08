package weather

import "github.com/gin-gonic/gin"


func StationsHandler(ctx *gin.Context) {
	stationName := ctx.Query("name")
	stationCode := ctx.Query("code")

	if stationName != "" && stationCode != "" {
		ctx.JSON(400, gin.H{
			"status": "error",
			"errors": gin.H{
				"conflict": "can't search both name and code at the same time",
			},
		})
		return
	}

	stations, err := GetClimateStations()
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": "error",
			"errors": gin.H{
				"fetch error": "unable to fetch the data",
			},
		})
		return
	}

	if stationName != "" {
		stationsFound := searchStationName(stations, stationName)
		if len(stationsFound) == 0 {
			ctx.JSON(404, gin.H{
				"status": "success",
				"data": nil,
			})
			return
		}

		ctx.JSON(200, gin.H{
			"status": "success",
			"data":   stationsFound,
		})
		return
	}

	if stationCode != "" {
		match, found := searchStationCode(stations, stationCode)
		if !found {
			ctx.JSON(404, gin.H{
				"status": "success",
				"data": nil,
			})
			return
		}

		ctx.JSON(200, gin.H{
			"status": "success",
			"data":   match,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status": "success",
		"data":   stations,
	})
}
