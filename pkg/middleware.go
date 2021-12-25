package libreapi

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"net/http"
)

func useCORSAllowAllMiddleware(e *gin.Engine) {
	e.Use(cors.Default())
}

func useGzipMiddleware(e *gin.Engine, compression int) {
	e.Use(gzip.Gzip(compression))
}

func useBanlistMiddleware(e *gin.Engine, banlist []string) {
	e.Use(func(c *gin.Context) {
		for _, ip := range banlist {
			if ip == c.ClientIP() {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"status": "fail",
					"errors": gin.H{
						"remote": "remote is forbidden",
					},
				})

				return
			}
		}

		c.Next()
	})
}

func addServerInfoMiddleware(e *gin.Engine) {
	e.Use(func(c *gin.Context) {
		c.Header("Server", Descriptor)
	})
}

func useRateLimiter(e *gin.Engine, rate string) error {
	rate2, err := limiter.NewRateFromFormatted(rate)
	if err != nil {
		return err
	}

	store := memory.NewStore()

	instance := limiter.New(store, rate2)

	e.Use(mgin.NewMiddleware(instance))

	return nil
}
