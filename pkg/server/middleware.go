package server

import (
	"time"

	"github.com/ccuetoh/libreapi/pkg"
	"github.com/ccuetoh/libreapi/pkg/env"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func setupMiddlewares(server *Server) {
	server.engine.Use(gin.Recovery())

	if server.env.NewRelic != nil {
		server.engine.Use(nrgin.Middleware(server.env.NewRelic))
	}

	server.engine.Use(
		loggingMiddleware(server.env),
		gzip.Gzip(gzip.DefaultCompression),
		cors.Default(),
		serverInfoMiddleware(),
		rateLimiterMiddleware(server.env.Cfg.HTTP.ProxyClientIPHeader))
}

func loggingMiddleware(env *env.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		ip := c.ClientIP()
		if env.Cfg.HTTP.ProxyClientIPHeader != "" {
			ip = c.GetHeader(env.Cfg.HTTP.ProxyClientIPHeader)
		}

		env.Log(c).Infof("%s %d %s %s", c.Request.Method, c.Writer.Status(), ip, c.Request.RequestURI)
	}
}

func serverInfoMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Server", libreapi.Descriptor)
	}
}

func rateLimiterMiddleware(ipHeader string) gin.HandlerFunc {
	rate := limiter.Rate{
		Period: time.Hour,
		Limit:  1500,
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
