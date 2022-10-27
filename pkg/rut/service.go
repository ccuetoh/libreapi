package rut

import (
	"encoding/base64"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	Activities []*Activity
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

	codeDecoded, err := base64.StdEncoding.DecodeString(data.Code)
	if err != nil {
		return "", "", errors.Wrap(err, "unable to decode captcha")
	}

	if len(codeDecoded) < 40 {
		return "", "", errors.New("captcha code is too short")
	}

	return data.Code, string(codeDecoded)[36:40], nil
}

func parseHTML(r io.ReadCloser) (*SIIProfile, error) {
	layout := "02-01-2006"

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create document")
	}

	name := clean(doc.Find("#contenedor > div:nth-child(4)"))
	if name == "**" {
		// No record found
		return &SIIProfile{}, nil
	}

	var activities []*Activity
	doc.Find("table.tabla:nth-child(27) > tbody:nth-child(1) > tr:nth-child(2)").Each(func(_ int, s *goquery.Selection) {
		code, err := strconv.Atoi(clean(s.Find("td:nth-child(2)")))
		if err != nil {
			return
		}

		date, err := time.Parse(layout, clean(s.Find("td:nth-child(5)")))
		if err != nil {
			return
		}

		activities = append(activities, &Activity{
			Name:         clean(s.Find("td:nth-child(1)")),
			Code:         code,
			Category:     clean(s.Find("td:nth-child(3)")),
			SubjectToVAT: clean(s.Find("td:nth-child(4)")) == "Si",
			Date:         date,
		})
	})

	return &SIIProfile{Name: name, Activities: activities}, nil
}

func clean(s *goquery.Selection) string {
	caser := cases.Title(language.LatinAmericanSpanish)
	return caser.String(strings.Trim(s.Text(), "\n\t "))
}
