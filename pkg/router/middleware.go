package router

import (
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"time"

	"github.com/ccuetoh/libreapi/pkg"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func setupMiddlewares(server *Server) {
	server.engine.Use(gin.Recovery())
	server.engine.Use(nrgin.Middleware(server.newRelic))
	server.engine.Use(loggingMiddleware(server))
	server.engine.Use(gzip.Gzip(gzip.DefaultCompression))
	server.engine.Use(cors.Default())
	server.engine.Use(serverInfoMiddleware())
	server.engine.Use(rateLimiterMiddleware(server.cfg.HTTP.ProxyClientIPHeader))
}

func loggingMiddleware(server *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if server.cfg.HTTP.ProxyClientIPHeader != "" {
			ip = c.GetHeader(server.cfg.HTTP.ProxyClientIPHeader)
		}

		server.logEntry(c).Infof("%s %d %s %s", c.Request.Method, c.Writer.Status(), ip, c.Request.RequestURI)
	}
}

func serverInfoMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Server", libreapi.Descriptor)
	}
}

func rateLimiterMiddleware(ipHeader string) gin.HandlerFunc {
	rate := limiter.Rate{
		Period: time.Minute,
		Limit:  20,
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate)

	if ipHeader == "" {
		return mgin.NewMiddleware(instance)

	}

	return mgin.NewMiddleware(instance, mgin.WithKeyGetter(
		func(c *gin.Context) string {
			return c.GetHeader(ipHeader)
		}))
}
