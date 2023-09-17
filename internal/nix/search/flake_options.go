package nix_search

import (
	"context"
	"strings"

	"github.com/luisnquin/nix-search/internal/nix"
	"github.com/microcosm-cc/bluemonday"
)

type flakeOptionResponseItem struct {
	FlakeName         string        `json:"flake_name"`
	FlakeDescription  string        `json:"flake_description"`
	FlakeResolved     flakeResolved `json:"flake_resolved"`
	Revision          string        `json:"revision"`
	OptionName        string        `json:"option_name"`
	OptionDescription string        `json:"option_description"`
	OptionSource      *string       `json:"option_source"`
	OptionDefault     string        `json:"option_default"`
	OptionExample     *string       `json:"option_example"`
	OptionType        string        `json:"option_type"`
}

type flakeOptionsResponse = searchResponse[flakeOptionResponseItem]

func (c Client) SearchFlakeOptions(ctx context.Context, flakesBranchId, searchTerm string, maxCount int) ([]*nix.FlakeOption, error) {
	response, err := c.searchFlakeOptions(ctx, flakesBranchId, searchTerm, maxCount)
	if err != nil {
		return nil, err
	}

	sp := bluemonday.StrictPolicy()
	replacer := strings.NewReplacer("\n", " ")

	options := make([]*nix.FlakeOption, len(response.Hits.Items))

	for i, item := range response.Hits.Items {
		options[i] = &nix.FlakeOption{
			Flake: &nix.FlakeMetadata{
				Name:        item.Source.FlakeName,
				Description: item.Source.FlakeDescription,
				Origin: nix.FlakeOrigin{
					Owner: item.Source.FlakeResolved.Owner,
					Repo:  item.Source.FlakeResolved.Repo,
					Type:  item.Source.FlakeResolved.Type,
				},
			},
			Revision: item.Source.Revision,
			Option: &nix.Option{
				Name:        item.Source.OptionName,
				Description: replacer.Replace(sp.Sanitize(item.Source.OptionDescription)),
				Example:     item.Source.OptionExample,
				Default:     item.Source.OptionDefault,
				Source:      item.Source.OptionSource,
				Type:        item.Source.OptionType,
			},
		}
	}

	return options, nil
}

func (c Client) searchFlakeOptions(ctx context.Context, flakesBranchId, searchTerm string, maxCount int) (flakeOptionsResponse, error) {
	const query = // AQL
	`
	{
		"query": {
			"bool": {
				"filter": [{ "term": { "type": { "value": "option", "_name": "filter_options" } } }],
				"must": [
					{
						"dis_max": {
							"tie_breaker": 0.7,
							"queries": [
								{
									"multi_match": {
										"type": "cross_fields",
										"query": "{{ . }}",
										"analyzer": "whitespace",
										"auto_generate_synonyms_phrase_query": false,
										"operator": "and",
										"_name": "multi_match_{{ . }}",
										"fields": [
											"option_name^6",
											"option_name.*^3.5999999999999996",
											"option_description^1",
											"option_description.*^0.6",
											"flake_name^0.5",
											"flake_name.*^0.3"
										]
									}
								},
								{ "wildcard": { "option_name": { "value": "*{{ . }}*", "case_insensitive": true } } }
							]
						}
					}
				]
			}
		}
	}
	`

	if maxCount <= 0 || maxCount > MAX_RESULTS_COUNT {
		maxCount = MAX_RESULTS_COUNT
	}

	esClient, err := c.prepareElasticSearchClient(flakesBranchId)
	if err != nil {
		return flakeOptionsResponse{}, err
	}

	response, err := esClient.Search(
		esClient.Search.WithFrom(0),
		esClient.Search.WithSize(maxCount),
		esClient.Search.WithSort("_score:desc", "option_name:desc"),
		esClient.Search.WithBody(prepareQuery(query, searchTerm)),
	)
	if err != nil {
		return flakeOptionsResponse{}, err
	}

	defer response.Body.Close()

	if response.IsError() {
		return flakeOptionsResponse{}, handleSearchErrorResponse(response.Body)
	}

	return parseSearchResponse[flakeOptionResponseItem](response.Body)
}
