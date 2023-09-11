package nix_search

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/luisnquin/nix-search/internal/config"
	"github.com/stoewer/go-strcase"
)

type (
	Client struct {
		config *config.Config

		homeManager homeManager
	}

	homeManager struct {
		data homeManagerOptionsData
		mu   *sync.Mutex
		*sync.Once
	}
)

const (
	CLIENT_TIMEOUT = time.Second * 20

	DEFAULT_SEARCH_QUERY_SIZE             = 50
	DEFAULT_SEARCH_QUERY_AGGREGATION_SIZE = 20
)

func NewClient(ctx context.Context, config *config.Config) (Client, error) {
	return Client{
		homeManager: homeManager{
			Once: new(sync.Once),
			mu:   new(sync.Mutex),
		},
		config: config,
	}, nil
}

func (c Client) prepareElasticSearchClient(channel string) (*elasticsearch.Client, error) {
	return elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("%s/%s", c.config.Internal.NixOSElasticSearch.Host, channel),
		},
		Username: c.config.Internal.NixOSElasticSearch.Username,
		Password: c.config.Internal.NixOSElasticSearch.Password,
	})
}

func getQueryTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"lower_case": strings.ToLower,
		"snake_case": strcase.SnakeCase,
		"kebab_case": strcase.KebabCase,
	}
}

func prepareQuery(rawQuery, searchTerm string) io.Reader {
	tpl := template.New("query").Funcs(getQueryTemplateFuncs())
	tpl = template.Must(tpl.Parse(rawQuery))

	var b bytes.Buffer

	if err := tpl.Execute(&b, searchTerm); err != nil {
		panic(err)
	}

	return &b
}
