package test

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"testing"
)

func AssertResponseBody(t *testing.T, recorder *httptest.ResponseRecorder, expect any) {
	assertResponse(t, recorder, expect, make(map[string]interface{}))
}

func AssertResponseBodySlice(t *testing.T, recorder *httptest.ResponseRecorder, expect any) {
	assertResponse(t, recorder, expect, make([]map[interface{}]interface{}, 0))
}

func assertResponse(t *testing.T, recorder *httptest.ResponseRecorder, expect any, expectData any) {
	var buf bytes.Buffer
	reader := io.TeeReader(recorder.Body, &buf)

	respBody, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("unable to read response data: %v", err)
	}

	resp := struct {
		Data interface{} `json:"data"`
	}{}

	err = jsoniter.Unmarshal(respBody, &resp)
	if err != nil {
		t.Fatalf("unable to unmarhsall response body: %v", err)
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &expectData,
		TagName: "json",
	})
	if err != nil {
		t.Fatalf("unable to create decoder: %v", err)
	}

	err = decoder.Decode(expect)
	if err != nil {
		t.Fatalf("unable to marshall expected data: %v", err)
	}

	got := fmt.Sprint(resp.Data)
	if got == "<nil>" {
		got = "[]"
	}

	assert.Equal(t, fmt.Sprint(expectData), got)
}
