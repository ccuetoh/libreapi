package server

import (
	"github.com/ccuetoh/libreapi/pkg/config"
	"github.com/ccuetoh/libreapi/pkg/env"
	contextNrLogrus "github.com/newrelic/go-agent/v3/integrations/logcontext-v2/nrlogrus"
	"net"
	"time"

	"github.com/ccuetoh/libreapi/pkg/economy"
	"github.com/ccuetoh/libreapi/pkg/rut"
	"github.com/ccuetoh/libreapi/pkg/weather"

	"github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrlogrus"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Server struct {
	engine *gin.Engine
	env    *env.Env
}

func NewServer(cfgOpts ...config.Option) (*Server, error) {
	cfg := config.Build(cfgOpts...)
	logger := logrus.New()

	var newRelicApp *newrelic.Application
	if cfg.NewRelic.Licence != "" {
		var err error
		newRelicApp, err = newNewRelic(cfg, logger)
		if err != nil {
			return nil, errors.Wrap(err, "unable to create NewRelic instance")
		}

		formatter := contextNrLogrus.NewFormatter(newRelicApp, &logrus.TextFormatter{})
		logger.SetFormatter(formatter)
	}

	server := &Server{
		engine: newEngine(cfg),
		env: &env.Env{
			Logger:   logger,
			Cfg:      cfg,
			NewRelic: newRelicApp,
		},
	}

	setupMiddlewares(server)
	addEndpoints(server)

	return server, nil
}

func (s *Server) Start() error {
	return s.engine.Run(net.JoinHostPort(s.env.Cfg.HTTP.Address, s.env.Cfg.HTTP.Port))
}

func newNewRelic(cfg *config.Config, logger *logrus.Logger) (*newrelic.Application, error) {
	return newrelic.NewApplication(
		newrelic.ConfigAppName(cfg.NewRelic.AppName),
		newrelic.ConfigLicense(cfg.NewRelic.Licence),
		newrelic.ConfigAppLogForwardingEnabled(cfg.NewRelic.LogForwardingEnabled),
		func(config *newrelic.Config) {
			logger.SetLevel(logrus.InfoLevel)
			config.Logger = nrlogrus.Transform(logger)
		},
	)
}

func newEngine(cfg *config.Config) *gin.Engine {
	if cfg.HTTP.DebugEnabled {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	return gin.New()
}

func addEndpoints(server *Server) {
	store := persist.NewMemoryStore(time.Minute)

	server.engine.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	rutHandler := rut.NewHandler(server.env, rut.NewDefaultService())
	rutGroup := server.engine.Group("/rut")

	rutGroup.GET("/random", rutHandler.Generate())

	rutGroup.Use(cache.CacheByRequestURI(store, time.Hour))
	rutGroup.GET("/validate", rutHandler.Validate())
	rutGroup.GET("/digit", rutHandler.VD())
	rutGroup.GET("/activities", rutHandler.Activity())

	economyHandler := economy.NewHandler(server.env, economy.NewDefaultService())
	economyGroup := server.engine.Group("/economy")

	economyGroup.Use(cache.CacheByRequestURI(store, time.Minute*5))
	economyGroup.GET("/indicators", economyHandler.Indicators())
	economyGroup.GET("/currencies", economyHandler.Currencies())

	weatherHandler := weather.NewHandler(server.env, weather.NewService())
	weatherGroup := server.engine.Group("/weather")

	weatherGroup.Use(cache.CacheByRequestURI(store, time.Minute*5))
	weatherGroup.GET("/stations", weatherHandler.Stations())
}
