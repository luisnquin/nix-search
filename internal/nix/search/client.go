package nix_search

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/luisnquin/nix-search/internal/config"
	"github.com/luisnquin/nix-search/internal/nix"
	"github.com/stoewer/go-strcase"
)

type (
	Client struct {
		config *config.Config
		store  store
	}

	store struct {
		homeManagerShell homeManagerShell
	}

	homeManagerShell struct {
		options []*nix.Option
		*sync.Once
	}
)

const (
	CLIENT_TIMEOUT = time.Second * 20

	DEFAULT_SEARCH_QUERY_SIZE             = 50
	DEFAULT_SEARCH_QUERY_AGGREGATION_SIZE = 20
)

func NewClient(config *config.Config) *Client {
	return &Client{
		config: config,
		store: store{
			homeManagerShell: homeManagerShell{
				Once: new(sync.Once),
			},
		},
	}
}

func (c Client) prepareElasticSearchClient(indexBranch string) (*elasticsearch.Client, error) {
	esConfig := c.config.Internal.Nix.Sources.ElasticSearch

	indexParts := []string{INDEX_PREFIX, esConfig.MappingVersion, indexBranch}
	index := strings.Join(indexParts, "-")

	return elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("%s/%s", esConfig.URL, index),
		},
		Username: esConfig.Username,
		Password: esConfig.Password,
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

func doGET(ctx context.Context, url string) (*http.Response, error) {
	httpClient := http.Client{Timeout: CLIENT_TIMEOUT}

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	return httpClient.Do(r)
}
