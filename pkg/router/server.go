package router

import (
	"context"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"net"
	"time"

	"github.com/ccuetoh/libreapi/pkg/config"
	"github.com/ccuetoh/libreapi/pkg/economy"
	"github.com/ccuetoh/libreapi/pkg/rut"
	"github.com/ccuetoh/libreapi/pkg/weather"

	"github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	contextNrLogrus "github.com/newrelic/go-agent/v3/integrations/logcontext-v2/nrlogrus"
	"github.com/newrelic/go-agent/v3/integrations/nrlogrus"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Server struct {
	engine   *gin.Engine
	logger   *logrus.Logger
	cfg      *config.Config
	newRelic *newrelic.Application
}

func NewServer(cfgOpts ...config.Option) (*Server, error) {
	cfg := config.Build(cfgOpts...)
	logger := logrus.New()

	newRelicApp, err := newNewRelic(cfg, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create NewRelic instance")
	}

	formatter := contextNrLogrus.NewFormatter(newRelicApp, &logrus.TextFormatter{})
	logger.SetFormatter(formatter)

	engine := newEngine(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create Gin engine")
	}

	server := &Server{
		engine:   engine,
		logger:   logger,
		cfg:      cfg,
		newRelic: newRelicApp,
	}

	setupMiddlewares(server)
	addEndpoints(server)

	return server, nil
}

func (s *Server) Start() error {
	return s.engine.Run(net.JoinHostPort(s.cfg.HTTP.Address, s.cfg.HTTP.Port))
}

func (s *Server) logEntry(c *gin.Context) *logrus.Entry {
	txn := nrgin.Transaction(c)
	ctx := newrelic.NewContext(context.Background(), txn)

	ip := c.ClientIP()
	if s.cfg.HTTP.ProxyClientIPHeader != "" {
		ip = c.GetHeader(s.cfg.HTTP.ProxyClientIPHeader)
	}

	fields := logrus.Fields{
		"client":      c.GetHeader("User-Agent"),
		"method":      c.Request.Method,
		"request_uri": c.Request.RequestURI,
		"ip":          ip,
	}

	return s.logger.WithContext(ctx).WithFields(fields)
}

func newNewRelic(cfg *config.Config, logger *logrus.Logger) (*newrelic.Application, error) {
	return newrelic.NewApplication(
		newrelic.ConfigAppName(cfg.NewRelic.AppName),
		newrelic.ConfigLicense(cfg.NewRelic.Licence),
		newrelic.ConfigAppLogForwardingEnabled(cfg.NewRelic.LogForwardingEnabled),
		func(config *newrelic.Config) {
			logrus.SetLevel(logrus.DebugLevel)
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

	rutGroup := server.engine.Group("/rut")
	rutGroup.GET("/validate", cache.CacheByRequestURI(store, time.Hour), rut.ValidateHandler)
	rutGroup.GET("/digit", cache.CacheByRequestURI(store, time.Hour), rut.DigitHandler)
	rutGroup.GET("/activities", cache.CacheByRequestURI(store, time.Hour), rut.SIIActivityHandler)

	economyGroup := server.engine.Group("/economy")
	economyGroup.GET("/indicators", cache.CacheByRequestPath(store, time.Hour), economy.BancoCentralIndicatorsHandler)
	economyGroup.GET("/crypto", cache.CacheByRequestURI(store, time.Minute*5), economy.BancoCentralIndicatorsHandler)
	economyGroup.GET("/currencies", cache.CacheByRequestURI(store, time.Hour), economy.CurrencyHandler)

	weatherGroup := server.engine.Group("/weather")
	weatherGroup.GET("/stations", cache.CacheByRequestURI(store, time.Minute*30), weather.StationsHandler)
}
