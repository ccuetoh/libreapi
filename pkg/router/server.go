package router

import (
	"github.com/ccuetoh/libreapi/pkg/config"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Server struct {
	engine *gin.Engine
	logger *logrus.Logger
	cfg    *config.Config
}

func NewServer(cfgOpts ...config.Option) (*Server, error) {
	cfg := config.Build(cfgOpts...)

	newRelicApp, err := newNewRelic(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create NewRelic instance")
	}

	engine, err := newEngine(cfg, newRelicApp)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create Gin engine")
	}

	return &Server{
		engine: engine,
		logger: logrus.New(),
		cfg:    cfg,
	}, nil
}

func (s *Server) Start() error {
	return s.engine.Run(net.JoinHostPort(s.cfg.HTTP.Address, s.cfg.HTTP.Port))
}

func newNewRelic(cfg *config.Config) (*newrelic.Application, error) {
	return newrelic.NewApplication(
		newrelic.ConfigAppName(cfg.NewRelic.AppName),
		newrelic.ConfigLicense(cfg.NewRelic.Licence),
		newrelic.ConfigAppLogForwardingEnabled(cfg.NewRelic.LogForwardingEnabled),
	)
}

func newEngine(cfg *config.Config, newRelic *newrelic.Application) (*gin.Engine, error) {
	if cfg.HTTP.DebugEnabled {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	err := setupMiddlewares(engine, newRelic)
	if err != nil {
		return nil, err
	}

	addEndpoints(engine)

	return engine, nil
}
