package router

import (
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
	"net/http"

	"github.com/ccuetoh/libreapi/pkg"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func setupMiddlewares(e *gin.Engine, newRelicApp *newrelic.Application) error {
	e.Use(nrgin.Middleware(newRelicApp))

	e.Use(gzip.Gzip(gzip.DefaultCompression))
	e.Use(cors.Default())
	e.Use(serverInfoMiddleware)

	// useBanlist(e, banlist)

	err := useRateLimiter(e, "1-S")
	if err != nil {
		return errors.Wrap(err, "rate limiter")
	}

	return nil
}

func useBanlist(e *gin.Engine, banlist []string) {
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

func serverInfoMiddleware(c *gin.Context) {
	c.Header("Server", libreapi.Descriptor)
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
