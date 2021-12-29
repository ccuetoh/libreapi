package weather

import (
	"context"
	"fmt"
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
	Code        int           `json:"code"`
	Name        string        `json:"name"`
	Operational bool          `json:"operational"`
	LastReport  time.Time     `json:"last_report"`
	Temperature float64       `json:"temperature"`
	Humidity    float64       `json:"humidity"`
	PresureHPas float64       `json:"presure_hpas"`
	Today       ClimateReport `json:"today"`
	Yesterday   ClimateReport `json:"yesterday"`
}

type Precipitations struct {
	Sum float64 `json:"sum"`
	EMA float64 `json:"ema"`
}

type ClimateReport struct {
	Maximum        ClimateInstance `json:"maximum"`
	Minimum        ClimateInstance `json:"minimum"`
	Precipitations Precipitations  `json:"precipitations"`
}

type ClimateInstance struct {
	Time        time.Time `json:"time"`
	Temperature float64   `json:"temperature"`
}

func GetClimateStations() ([]ClimateStation, error) {
	timeoutContext, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, err := http.NewRequestWithContext(
		timeoutContext,
		"GET",
		"https://climatologia.meteochile.gob.cl/application/diario/climatDiarioRecienteEmas/",
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "LibreAPI")
	req.Header.Set("Accept", "*/*")

	var client = http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %s", res.Status)
	}

	stations, err := ParseClimateHTML(res.Body)
	if err != nil {
		return nil, err
	}

	return stations, nil
}

func ParseClimateHTML(r io.ReadCloser) (stations []ClimateStation, err error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	doc.Find(".table").Find("tr").Each(func(i int, s *goquery.Selection) {
		if i < 4 {
			// Header
			return
		}

		station := ClimateStation{Operational: true}
		s.Find("td").Each(func(i2 int, s2 *goquery.Selection) {
			if s2.Find("strike").Length() > 0 {
				station.Operational = false
			}

			if !station.Operational {
				return
			}

			content := strings.TrimSpace(s2.Text())
			if content == "." {
				// Sometimes information is missing and a dot is the value
				return
			}

			switch i2 {
			case 1:
				station.Code, _ = strconv.Atoi(content)
			case 2:
				station.Name = content
			case 3:
				station.LastReport = todayHourToTime(content)
			case 4:
				station.Temperature = tempToFloat64(content)
			case 5:
				station.Humidity = humidityToFloat64(content)
			case 6:
				station.PresureHPas = presureToFloat64(content)
			case 7:
				station.Today.Maximum.Temperature = tempToFloat64(content)
			case 8:
				station.Today.Maximum.Time = todayHourToTime(content)
			case 9:
				station.Today.Minimum.Temperature = tempToFloat64(content)
			case 10:
				station.Today.Minimum.Time = todayHourToTime(content)
			case 11:
				station.Yesterday.Maximum.Temperature = tempToFloat64(content)
			case 12:
				station.Yesterday.Maximum.Time = yesterdayHourToTime(content)
			case 13:
				station.Yesterday.Minimum.Temperature = tempToFloat64(content)
			case 14:
				station.Yesterday.Minimum.Time = yesterdayHourToTime(content)
			case 15:
				station.Today.Precipitations.Sum = precipitationsToFloat64(content)
			case 16:
				station.Today.Precipitations.EMA = precipitationsToFloat64(content)
			case 17:
				station.Yesterday.Precipitations.Sum = precipitationsToFloat64(content)
			case 18:
				station.Yesterday.Precipitations.EMA = precipitationsToFloat64(content)
			}
		})

		stations = append(stations, station)
	})

	return stations, nil
}

func todayHourToTime(hour string) time.Time {
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
	res, _ := strconv.ParseFloat(temp, 64)
	return res
}

func humidityToFloat64(hum string) float64 {
	res, _ := strconv.ParseFloat(hum, 64)
	return res / 100
}

func presureToFloat64(p string) float64 {
	p = strings.ReplaceAll(p, ",", "")
	res, _ := strconv.ParseFloat(p, 64)
	return res
}

func precipitationsToFloat64(pre string) float64 {
	if pre == "s/p" || pre == "." {
		return 0
	}

	res, _ := strconv.ParseFloat(pre, 64)
	return res
}

func searchStationName(stations []ClimateStation, patterns string) []ClimateStation {
	var hints []string
	for _, station := range stations {
		hints = append(hints, fmt.Sprintf("%s", removeTilde(station.Name)))
	}

	matches := fuzzy.Find(removeTilde(patterns), hints)

	var results []ClimateStation
	for _, match := range matches {
		results = append(results, stations[match.Index])
	}

	return results
}

func searchStationCode(stations []ClimateStation, code string) (ClimateStation, bool) {
	for _, station := range stations {
		if strconv.Itoa(station.Code) == code {
			return station, true
		}
	}

	return ClimateStation{}, false
}

func removeTilde(text string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, text)
	return result
}
