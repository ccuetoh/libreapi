package libreapi

import (
	"compress/gzip"
	"github.com/CamiloHernandez/libreapi/pkg/weather"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io"
	"os"
	"time"

	"github.com/CamiloHernandez/libreapi/pkg/economy"
	"github.com/CamiloHernandez/libreapi/pkg/rut"
)

type TLSPaths struct {
	CertificatePath string
	KeyPath         string
}

func Start(port int, certs ...TLSPaths) error {
	f, err := os.Create("libreapi.log")
	if err != nil {
		panic(errors.Wrap(err, "log file"))
	}

	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	s := getDefaultHTTPServer(port)

	var tlsPaths TLSPaths
	useTLS := false
	if len(certs) > 0 {
		tlsPaths = certs[0]
		useTLS = true
	}

	if useTLS {
		tlsConfig := getSSLLabAConfig()
		s.TLSConfig = tlsConfig
	}

	useGzipMiddleware(engine, gzip.DefaultCompression)
	useCORSAllowAllMiddleware(engine)
	addServerInfoMiddleware(engine)

	err = useRateLimiter(engine, "3-S")
	if err != nil {
		panic(errors.Wrap(err, "rate limiter"))
	}

	addEndpoints(engine)
	s.Handler = engine

	if useTLS {
		return s.ListenAndServeTLS(tlsPaths.CertificatePath, tlsPaths.KeyPath)
	} else {
		return s.ListenAndServe()
	}
}

func addEndpoints(e *gin.Engine) {
	day := time.Hour * 24
	store := persistence.NewInMemoryStore(day)

	e.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	rutGroup := e.Group("/rut")
	rutGroup.GET("/validate", cache.CachePage(store, day, rut.ValidateHandler))
	rutGroup.GET("/digit", cache.CachePage(store, day, rut.DigitHandler))
	rutGroup.GET("/activities", cache.CachePage(store, day, rut.SIIActivityHandler))

	economyGroup := e.Group("/economy")
	economyGroup.GET("/indicators", cache.CachePage(store, time.Hour, economy.BancoCentraIndicatorsHandler))
	economyGroup.GET("/crypto", cache.CachePage(store, time.Minute * 15, economy.CryptoHandler))
	economyGroup.GET("/currencies", cache.CachePage(store, day/2, economy.CurrencyHandler))

	weatherGrop := e.Group("/weather")
	weatherGrop.GET("/stations", cache.CachePage(store, time.Minute * 30, weather.StationsHandler))
}
