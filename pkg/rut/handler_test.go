package rut

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/ccuetoh/libreapi/internal/test"
	"github.com/ccuetoh/libreapi/pkg/env"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

type MockService struct {
	profile    *SIIProfile
	profileErr error
}

func (s MockService) GetProfile(_ RUT) (*SIIProfile, error) {
	return s.profile, s.profileErr
}

func TestValidateNoRut(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := MockService{}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	handler.Validate()(ctx)

	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}

func TestValidateInvalid(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := MockService{}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = &http.Request{}
	ctx.Request.URL, _ = url.Parse("?rut=1231231")

	handler.Validate()(ctx)

	assert.Equal(t, recorder.Code, http.StatusOK)
	test.AssertResponseBody(t, recorder, struct {
		RUT   string `json:"rut"`
		Valid bool   `json:"valid"`
	}{"123123-1", false})
}

func TestValidateValid(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := MockService{}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = &http.Request{}
	ctx.Request.URL, _ = url.Parse("?rut=5.126.663-3")

	handler.Validate()(ctx)

	assert.Equal(t, recorder.Code, http.StatusOK)
	test.AssertResponseBody(t, recorder, struct {
		RUT   string `json:"rut"`
		Valid bool   `json:"valid"`
	}{"5126663-3", true})
}

func TestValidateBadRUT(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := MockService{}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = &http.Request{}
	ctx.Request.URL, _ = url.Parse("?rut=asdasd")

	handler.Validate()(ctx)

	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}

func TestVDOk(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := MockService{}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = &http.Request{}
	ctx.Request.URL, _ = url.Parse("?rut=5.811.892")

	handler.VD()(ctx)

	assert.Equal(t, recorder.Code, http.StatusOK)
	test.AssertResponseBody(t, recorder, struct {
		RUT   string `json:"rut"`
		Digit string `json:"digit"`
	}{"5811892-3", "3"})
}

func TestVDNoRut(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := MockService{}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = &http.Request{}
	ctx.Request.URL, _ = url.Parse("?rut=")

	handler.VD()(ctx)

	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}

func TestVDBadRut(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := MockService{}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = &http.Request{}
	ctx.Request.URL, _ = url.Parse("?rut=asdasdasd")

	handler.VD()(ctx)

	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}

func TestActivityOk(t *testing.T) {
	gin.SetMode(gin.TestMode)

	profile := &SIIProfile{
		Name:       "Eduardo Alfredo Juan Bernardo Frei Ruiz–Tagle",
		Activities: []*Activity{},
	}

	service := MockService{profile: profile}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = &http.Request{}
	ctx.Request.URL, _ = url.Parse("?rut=4100738-9")

	handler.Activity()(ctx)

	assert.Equal(t, recorder.Code, http.StatusOK)
	test.AssertResponseBody(t, recorder, gin.H{
		"rut":        "4100738-9",
		"name":       profile.Name,
		"activities": []Activity{},
	})
}

func TestActivityError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := MockService{profileErr: errors.New("server is on fire")}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = &http.Request{}
	ctx.Request.URL, _ = url.Parse("?rut=4100738-9")

	handler.Activity()(ctx)

	assert.Equal(t, recorder.Code, http.StatusInternalServerError)
}

func TestActivityNoRut(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := MockService{}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	handler.Activity()(ctx)

	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}

func TestActivityBadRut(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := MockService{}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = &http.Request{}
	ctx.Request.URL, _ = url.Parse("?rut=asdasdasd")

	handler.Activity()(ctx)

	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}

func TestGenerateOk(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := MockService{}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	handler.Generate()(ctx)

	assert.Equal(t, recorder.Code, http.StatusOK)

	var buf bytes.Buffer
	reader := io.TeeReader(recorder.Body, &buf)

	respBody, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("unable to read response data: %v", err)
	}

	resp := struct {
		Data struct {
			Digits string `json:"digits"`
			VD     string `json:"vd"`
		} `json:"data"`
	}{}

	err = jsoniter.Unmarshal(respBody, &resp)
	if err != nil {
		t.Fatalf("unable to unmarhsall response body: %v", err)
	}

	rut, err := parseRUT(fmt.Sprintf("%s-%s", resp.Data.Digits, resp.Data.VD), false)
	assert.NoError(t, err)
	assert.True(t, rut.IsValid())
}

func TestGenerateOkRange(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := MockService{}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	const min = 5000000
	const max = 6000000

	ctx.Request = &http.Request{}
	ctx.Request.URL, _ = url.Parse(fmt.Sprintf("?min=%d&max=%d", min, max))

	handler.Generate()(ctx)

	assert.Equal(t, recorder.Code, http.StatusOK)

	var buf bytes.Buffer
	reader := io.TeeReader(recorder.Body, &buf)

	respBody, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("unable to read response data: %v", err)
	}

	resp := struct {
		Data struct {
			Digits string `json:"digits"`
			VD     string `json:"vd"`
		} `json:"data"`
	}{}

	err = jsoniter.Unmarshal(respBody, &resp)
	if err != nil {
		t.Fatalf("unable to unmarhsall response body: %v", err)
	}

	digitsNum, err := strconv.Atoi(resp.Data.Digits)
	assert.NoError(t, err)
	assert.LessOrEqual(t, digitsNum, max)
	assert.GreaterOrEqual(t, digitsNum, min)

	rut, err := parseRUT(fmt.Sprintf("%s-%s", resp.Data.Digits, resp.Data.VD), false)
	assert.NoError(t, err)
	assert.True(t, rut.IsValid())
}

func TestGenerateOkBadRange(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := MockService{}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	const min = 7000000
	const max = 6000000

	ctx.Request = &http.Request{}
	ctx.Request.URL, _ = url.Parse(fmt.Sprintf("?min=%d&max=%d", min, max))

	handler.Generate()(ctx)

	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}

func TestGenerateOkBadMin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := MockService{}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = &http.Request{}
	ctx.Request.URL, _ = url.Parse("?min=asdasd")

	handler.Generate()(ctx)

	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}

func TestGenerateOkBadMax(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := MockService{}
	handler := NewHandler(env.NewTestEnv(), service)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = &http.Request{}
	ctx.Request.URL, _ = url.Parse("?max=asdasd")

	handler.Generate()(ctx)

	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}
