package libreapi

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/CamiloHernandez/libreapi/pkg/economy"
	"github.com/CamiloHernandez/libreapi/pkg/rut"
	"github.com/CamiloHernandez/libreapi/pkg/weather"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type TLSPaths struct {
	CertificatePath string
	KeyPath         string
}

func Start(port int, certs ...TLSPaths) error {
	f, err := os.Create("libreapi.log")
	if err != nil {
		return errors.Wrap(err, "log file")
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

	var banlist []string
	binContent, err := ioutil.ReadFile("./banlist.txt")
	if err != nil {
		fmt.Println("Unable to read banlist")
	} else {
		banData := strings.ReplaceAll(string(binContent), "\r", "")
		for _, line := range strings.Split(banData, "\n") {
			if len(line) == 0 {
				continue
			}

			if line[0] != '#' {
				banlist = append(banlist, line)
			}
		}
	}

	useGzipMiddleware(engine, gzip.DefaultCompression)
	useCORSAllowAllMiddleware(engine)
	useBanlistMiddleware(engine, banlist)
	addServerInfoMiddleware(engine)

	err = useRateLimiter(engine, "1-S")
	if err != nil {
		return errors.Wrap(err, "rate limiter")
	}

	addEndpoints(engine)
	s.Handler = engine

	if useTLS {
		return s.ListenAndServeTLS(tlsPaths.CertificatePath, tlsPaths.KeyPath)
	}

	return s.ListenAndServe()
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
	economyGroup.GET("/indicators", cache.CachePage(store, time.Hour, economy.BancoCentralIndicatorsHandler))
	economyGroup.GET("/crypto", cache.CachePage(store, time.Minute*15, economy.CryptoHandler))
	economyGroup.GET("/currencies", cache.CachePage(store, day/2, economy.CurrencyHandler))

	weatherGroup := e.Group("/weather")
	weatherGroup.GET("/stations", cache.CachePage(store, time.Minute*30, weather.StationsHandler))
}
