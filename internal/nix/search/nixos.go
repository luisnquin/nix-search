package nix_search

import (
	"context"

	"github.com/luisnquin/nix-search/internal/nix"
)

type nixosOptionResponseItem struct {
	Name        string  `json:"option_name"`
	Description string  `json:"option_description"`
	Example     *string `json:"option_example"`
	Default     string  `json:"option_default"`
	Source      *string `json:"option_source"`
}

type nixosOptionsResponse = searchResponse[nixosOptionResponseItem]

func (c Client) SearchNixOSOptions(ctx context.Context, channelBranch, searchTerm string, maxCount int) ([]*nix.Option, error) {
	response, err := c.searchNixOSOptions(ctx, channelBranch, searchTerm, maxCount)
	if err != nil {
		return nil, err
	}

	options := make([]*nix.Option, len(response.Hits.Items))

	for i, item := range response.Hits.Items {
		options[i] = &nix.Option{
			Name:        item.Source.Name,
			Description: item.Source.Description,
			Example:     item.Source.Example,
			Default:     item.Source.Default,
			Source:      item.Source.Source,
		}
	}

	return options, nil
}

func (c Client) searchNixOSOptions(ctx context.Context, channelBranch, searchTerm string, maxCount int) (nixosOptionsResponse, error) {
	esClient, err := c.prepareElasticSearchClient(channelBranch)
	if err != nil {
		return searchResponse[nixosOptionResponseItem]{}, err
	}

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
								{
									"wildcard": {
										"option_name": {
											"value": "*{{ . | lower_case }}*", "case_insensitive": true
										}
									}
								}
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

	response, err := esClient.Search(
		esClient.Search.WithFrom(0),
		esClient.Search.WithSize(maxCount),
		esClient.Search.WithSort("_score:desc", "option_name:desc"),
		esClient.Search.WithBody(prepareQuery(query, searchTerm)),
	)
	if err != nil {
		return nixosOptionsResponse{}, err
	}

	defer response.Body.Close()

	if response.IsError() {
		return nixosOptionsResponse{}, handleSearchErrorResponse(response.Body)
	}

	return parseSearchResponse[nixosOptionResponseItem](response.Body)
}
