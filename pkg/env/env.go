package env

import (
	"context"
	"github.com/ccuetoh/libreapi/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type Env struct {
	Logger   *logrus.Logger
	Cfg      *config.Config
	NewRelic *newrelic.Application
}

func (e *Env) Log(c *gin.Context) *logrus.Entry {
	if e.NewRelic == nil {
		return e.Logger.WithContext(context.Background())
	}

	txn := nrgin.Transaction(c)
	ctx := newrelic.NewContext(context.Background(), txn)

	return e.Logger.WithContext(ctx)
}

func NewTestEnv() *Env {
	return &Env{
		Logger: logrus.New(),
		Cfg:    config.Default(),
	}
}
