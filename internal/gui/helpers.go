package gui

import (
	"bytes"
	"encoding/json"
	"text/template"
)

func getRenderedText[T any](name, tplText string, items []T) (string, error) {
	data, err := structSliceToMapSlice(items)
	if err != nil {
		return "", err
	}

	tpl := template.New(name)
	tpl = template.Must(tpl.Parse(tplText))

	var b bytes.Buffer

	if err := tpl.Execute(&b, data); err != nil {
		return "", err
	}

	return b.String(), nil
}

func structSliceToMapSlice[T any](items []T) ([]map[string]any, error) {
	var b bytes.Buffer

	if err := json.NewEncoder(&b).Encode(items); err != nil {
		return nil, err
	}

	var results []map[string]any

	if err := json.NewDecoder(&b).Decode(&results); err != nil {
		return nil, err
	}

	return results, nil
}
