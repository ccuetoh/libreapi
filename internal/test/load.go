package test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func LoadHTML(caseName string) (io.ReadCloser, error) {
	filename := fmt.Sprintf("../../internal/test/data/%s/%s.html", caseName, caseName)
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func LoadJSON(caseName string, v any) error {
	filename := fmt.Sprintf("../../internal/test/data/%s/%s.json", caseName, caseName)
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}
