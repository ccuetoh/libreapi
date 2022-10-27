package economy

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type DefaultService struct {
	client *http.Client
}

func NewDefaultService() *DefaultService {
	return &DefaultService{
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (s *DefaultService) GetIndicators() (*Indicators, error) {
	res, err := s.client.Get("https://si3.bcentral.cl/Indicadoressiete/secure/Indicadoresdiarios.aspx")
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non ok status: %d %s", res.StatusCode, res.Status)
	}

	indicators, err := parseIndicatorsHTML(res.Body)
	if err != nil {
		return nil, err
	}

	return indicators, nil
}

func (s *DefaultService) GetCurrencies() ([]*Currency, error) {
	url, err := s.getDailyCurrenciesURL()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get url")
	}

	res, err := s.client.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "unable to execute request")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	currencies, err := parseCurrenciesHTML(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse html")
	}

	return currencies, nil
}
