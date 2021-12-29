package rut

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/json-iterator/go"
	errs "github.com/pkg/errors"
)

var (
	ErrNoCaptchaCode       = errors.New("unexpected response: no captcha code")
	ErrInvalidCaptchaCode  = errors.New("unexpected response: captcha code is not a string")
	ErrTooShortCaptchaCode = errors.New("unexpected response: captcha code is too short")
)

type SIIDetail struct {
	Name       string
	Activities []EconomicActivity
}

type EconomicActivity struct {
	Name         string    `json:"name"`
	Code         int       `json:"code"`
	Category     string    `json:"category"`
	SubjectToVAT bool      `json:"subject_to_vat"`
	Date         time.Time `json:"date"`
}

func GetCaptcha() (code string, captcha string, err error) {
	resp, err := http.Get("https://zeus.sii.cl/cvc_cgi/stc/CViewCaptcha.cgi?oper=0")
	if err != nil {
		return "", "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var data map[string]interface{}
	err = jsoniter.Unmarshal(body, &data)
	if err != nil {
		return "", "", err
	}

	codeIntf, ok := data["txtCaptcha"]
	if !ok {
		return "", "", ErrNoCaptchaCode
	}

	code, ok = codeIntf.(string)
	if !ok {
		return "", "", ErrInvalidCaptchaCode
	}

	if len(code) < 40 {
		return "", "", ErrTooShortCaptchaCode
	}

	codeDecoded, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		return "", "", errs.Wrap(err, "unable to decode captcha")
	}

	return code, string(codeDecoded)[36:40], nil
}

func GetSIIDetails(r RUT) (SIIDetail, error) {
	code, captcha, err := GetCaptcha()
	if err != nil {
		return SIIDetail{}, err
	}

	form := url.Values{}
	form.Add("RUT", r.String()[:len(r)-1])
	form.Add("DV", VDToString(r[len(r)-1]))
	form.Add("PRG", "STC")
	form.Add("OPC", "NOR")
	form.Add("txt_captcha", code) // For some reason this is expected inverted:
	form.Add("txt_code", captcha) // code in "txt_captcha" and captcha in "txt_code"

	timeoutContext, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, err := http.NewRequestWithContext(
		timeoutContext,
		"POST",
		"https://zeus.sii.cl/cvc_cgi/stc/getstc",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return SIIDetail{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var client = http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return SIIDetail{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return SIIDetail{}, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	detail, err := ParseHTML(res.Body)
	if err != nil {
		return SIIDetail{}, err
	}

	return detail, nil
}

func ParseHTML(r io.ReadCloser) (SIIDetail, error) {
	layout := "02-01-2006"

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return SIIDetail{}, err
	}

	var name string
	doc.Find("html").Find("body").Find("div").Find("div").Each(func(i int, s *goquery.Selection) {
		if i != 4 {
			return
		}

		name = strings.TrimSpace(strings.Title(strings.ToLower(s.Text())))
	})

	var activities []EconomicActivity
	doc.Find("html").Find("body").Find("div").Find("table").Each(func(i1 int, s *goquery.Selection) {
		if i1 != 0 {
			return
		}

		var activity EconomicActivity
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

	return SIIDetail{Name: name, Activities: activities}, nil
}

func cleanHTMLText(text string) string {
	return strings.Trim(strings.TrimSpace(strings.Title(strings.ToLower(text))), "\n\t")
}
