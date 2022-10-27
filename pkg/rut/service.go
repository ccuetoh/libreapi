package rut

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
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

type SIIProfile struct {
	Name       string
	Activities []Activity
}

type Activity struct {
	Name         string    `json:"name"`
	Code         int       `json:"code"`
	Category     string    `json:"category"`
	SubjectToVAT bool      `json:"subject_to_vat"`
	Date         time.Time `json:"date"`
}

func (s *DefaultService) GetProfile(rut RUT) (*SIIProfile, error) {
	code, captcha, err := s.getCaptcha()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get captcha")
	}

	form := url.Values{}
	form.Add("RUT", rut.String()[:len(rut)-1])
	form.Add("DV", VDToString(rut[len(rut)-1]))
	form.Add("PRG", "STC")
	form.Add("OPC", "NOR")
	form.Add("txt_captcha", code) // code is expected in "txt_captcha" and captcha in "txt_code"
	form.Add("txt_code", captcha)

	res, err := s.client.Post(
		"https://zeus.sii.cl/cvc_cgi/stc/getstc",
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "unable to execute request")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non ok status: %d %s", res.StatusCode, res.Status)
	}

	return parseHTML(res.Body)
}

func (s *DefaultService) getCaptcha() (code string, captcha string, err error) {
	resp, err := s.client.Get("https://zeus.sii.cl/cvc_cgi/stc/CViewCaptcha.cgi?oper=0")
	if err != nil {
		return "", "", errors.Wrap(err, "unable to execute request")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", errors.Wrap(err, "unable to read response body")
	}

	data := &struct {
		Code string `json:"txtCaptcha"`
	}{}

	err = jsoniter.Unmarshal(body, data)
	if err != nil {
		return "", "", errors.Wrap(err, "unable unmarshall response")
	}

	if len(data.Code) < 40 {
		return "", "", errors.New("captcha code is too short")
	}

	codeDecoded, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		return "", "", errors.Wrap(err, "unable to decode captcha")
	}

	return code, string(codeDecoded)[36:40], nil
}

func parseHTML(r io.ReadCloser) (*SIIProfile, error) {
	layout := "02-01-2006"

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	var name string
	doc.Find("html").Find("body").Find("div").Find("div").Each(func(i int, s *goquery.Selection) {
		if i != 4 {
			return
		}

		name = strings.TrimSpace(strings.Title(strings.ToLower(s.Text())))
	})

	var activities []Activity
	doc.Find("html").Find("body").Find("div").Find("table").Each(func(i1 int, s *goquery.Selection) {
		if i1 != 0 {
			return
		}

		var activity Activity
		s.Find("tbody").Find("tr").Each(func(i2 int, s *goquery.Selection) {
			if i2 == 0 {
				return
			}

			s.Find("td").Each(func(i3 int, s *goquery.Selection) {
				switch i3 {
				case 0:
					activity.Name = cleanHTMLText(s.Text())
				case 1:
					code, err := strconv.Atoi(cleanHTMLText(s.Text()))
					if err != nil {
						return
					}

					activity.Code = code
				case 2:
					activity.Category = cleanHTMLText(s.Text())
				case 3:
					activity.SubjectToVAT = cleanHTMLText(s.Text()) == "Si"
				case 4:
					date, err := time.Parse(layout, s.Text())
					if err != nil {
						return
					}

					activity.Date = date
				}

			})

			activities = append(activities, activity)
		})
	})

	if name == "**" {
		name = "" // No record found
	}

	return &SIIProfile{Name: name, Activities: activities}, nil
}

func cleanHTMLText(text string) string {
	return strings.Trim(strings.TrimSpace(strings.Title(strings.ToLower(text))), "\n\t")
}
