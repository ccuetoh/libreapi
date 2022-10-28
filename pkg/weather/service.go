package weather

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/sahilm/fuzzy"
)

type ClimateStation struct {
	Code        int            `json:"code"`
	Name        string         `json:"name"`
	Operational bool           `json:"operational"`
	LastReport  *time.Time     `json:"last_report,omitempty"`
	Temperature float64        `json:"temperature,omitempty"`
	Humidity    float64        `json:"humidity,omitempty"`
	PressureHPA float64        `json:"pressure_hpa,omitempty"`
	Today       *ClimateReport `json:"today,omitempty"`
	Yesterday   *ClimateReport `json:"yesterday,omitempty"`
}

type Precipitations struct {
	Sum float64 `json:"sum"`
	EMA float64 `json:"ema"`
}

type ClimateReport struct {
	Maximum        Measurement    `json:"maximum"`
	Minimum        Measurement    `json:"minimum"`
	Precipitations Precipitations `json:"precipitations"`
}

type Measurement struct {
	Time        time.Time `json:"time"`
	Temperature float64   `json:"temperature"`
}

type DefaultService struct {
	client *http.Client
}

func NewService() *DefaultService {
	return &DefaultService{
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (s *DefaultService) GetClimateStations() ([]*ClimateStation, error) {
	res, err := s.client.Get("https://climatologia.meteochile.gob.cl/application/diario/climatDiarioRecienteEmas/")
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %s", res.Status)
	}

	stations, err := parseClimateHTML(res.Body)
	if err != nil {
		return nil, err
	}

	return stations, nil
}

func parseClimateHTML(r io.ReadCloser) (stations []*ClimateStation, err error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	doc.Find(".table-bordered > tbody:nth-child(1) > tr").Each(func(i int, row *goquery.Selection) {
		if i < 3 {
			// Headers
			return
		}

		if strings.TrimSpace(row.Find("td:nth-child(4)").Text()) == "." {
			station := &ClimateStation{
				Name:        cleanName(row.Find("td:nth-child(3)").Text()),
				Operational: false,
			}

			station.Code, _ = strconv.Atoi(strings.TrimSpace(row.Find("td:nth-child(2)").Text()))
			stations = append(stations, station)
			return
		}

		station := &ClimateStation{
			Operational: true,
			Name:        cleanName(row.Find("td:nth-child(3)").Text()),
			Temperature: tempToFloat64(row.Find("td:nth-child(5)").Text()),
			Humidity:    humidityToFloat64(row.Find("td:nth-child(6)").Text()),
			PressureHPA: pressureToFloat64(row.Find("td:nth-child(7)").Text()),
			Today: &ClimateReport{
				Maximum: Measurement{
					Time:        todayHourToTime(row.Find("td:nth-child(9)").Text()),
					Temperature: tempToFloat64(row.Find("td:nth-child(8)").Text()),
				},
				Minimum: Measurement{
					Time:        todayHourToTime(row.Find("td:nth-child(11)").Text()),
					Temperature: tempToFloat64(row.Find("td:nth-child(10)").Text()),
				},
				Precipitations: Precipitations{
					Sum: precipitationsToFloat64(row.Find("td:nth-child(16)").Text()),
					EMA: precipitationsToFloat64(row.Find("td:nth-child(17)").Text()),
				},
			},
			Yesterday: &ClimateReport{
				Maximum: Measurement{
					Time:        yesterdayHourToTime(row.Find("td:nth-child(13)").Text()),
					Temperature: tempToFloat64(row.Find("td:nth-child(12)").Text()),
				},
				Minimum: Measurement{
					Time:        yesterdayHourToTime(row.Find("td:nth-child(15)").Text()),
					Temperature: tempToFloat64(row.Find("td:nth-child(14)").Text()),
				},
				Precipitations: Precipitations{
					Sum: precipitationsToFloat64(row.Find("td:nth-child(18)").Text()),
					EMA: precipitationsToFloat64(row.Find("td:nth-child(19)").Text()),
				},
			},
		}

		lastReport := todayHourToTime(row.Find("td:nth-child(4)").Text())
		station.LastReport = &lastReport
		station.Code, _ = strconv.Atoi(strings.TrimSpace(row.Find("td:nth-child(2)").Text()))

		stations = append(stations, station)
	})

	return stations, nil
}

func todayHourToTime(hour string) time.Time {
	hour = strings.TrimSpace(hour)

	if hour == "." || hour == "" {
		return time.Time{}
	}

	loc, _ := time.LoadLocation("America/Santiago")

	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	sections := strings.Split(hour, ":")
	h, _ := strconv.Atoi(sections[0])
	m, _ := strconv.Atoi(sections[1])

	return t.Add((time.Duration(h) * time.Hour) + (time.Duration(m) * time.Minute))
}

func yesterdayHourToTime(hour string) time.Time {
	today := todayHourToTime(hour)
	return today.AddDate(0, 0, -1)
}

func tempToFloat64(temp string) float64 {
	temp = strings.TrimSpace(temp)

	if temp == "." {
		return 0
	}

	res, _ := strconv.ParseFloat(temp, 64)
	return res
}

func humidityToFloat64(hum string) float64 {
	hum = strings.TrimSpace(hum)

	if hum == "." {
		return 0
	}

	res, _ := strconv.ParseFloat(hum, 64)
	return res / 100
}

func pressureToFloat64(p string) float64 {
	p = strings.TrimSpace(p)

	if p == "." {
		return 0
	}

	p = strings.ReplaceAll(p, ",", "")
	res, _ := strconv.ParseFloat(p, 64)
	return res
}

func precipitationsToFloat64(pre string) float64 {
	pre = strings.TrimSpace(pre)

	if pre == "s/p" || pre == "." {
		return 0
	}

	res, _ := strconv.ParseFloat(pre, 64)
	return res
}

func removeTilde(text string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, text)
	return result
}

func cleanName(name string) string {
	// Station names are in Windows-1252 encoding. Transform to UTF-8
	var bufSlice []byte
	buf := bytes.NewBuffer(bufSlice)

	w1252Transformer := transform.NewWriter(buf, charmap.Windows1252.NewEncoder())
	w1252Transformer.Write([]byte(name))
	w1252Transformer.Close()

	return strings.TrimSpace(buf.String())
}

func searchStationName(stations []*ClimateStation, patterns string) []*ClimateStation {
	var hints []string
	for _, station := range stations {
		hints = append(hints, fmt.Sprintf("%s", removeTilde(station.Name)))
	}

	matches := fuzzy.Find(removeTilde(patterns), hints)

	var results []*ClimateStation
	for _, match := range matches {
		results = append(results, stations[match.Index])
	}

	return results
}

func searchStationCode(stations []*ClimateStation, code string) (*ClimateStation, bool) {
	for _, station := range stations {
		if strconv.Itoa(station.Code) == code {
			return station, true
		}
	}

	return nil, false
}
