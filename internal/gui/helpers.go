package gui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/luisnquin/nix-search/internal/config"
)

func getRenderedText[T any](name, branch, tplText string, items []T) (string, error) {
	data, err := structSliceToMapSlice(items)
	if err != nil {
		return "", err
	}

	transformSource := func(source string) string { // TODO: branch as optional parameter
		return fmt.Sprintf("https://github.com/NixOS/nixpkgs/blob/%s/%s",
			branch, strings.Replace(source, ":", "#L", -1))
	}

	tpl := template.New(name).Funcs(template.FuncMap{
		"transform_source": transformSource,
	})
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

func transform_source(channel config.NixChannel, source string) string {
	return fmt.Sprintf("https://github.com/NixOS/nixpkgs/blob/%s/%s",
		channel.Branch, strings.Replace(source, ":", "#L", -1))
}
