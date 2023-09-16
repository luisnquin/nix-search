package nix_search

import (
	"context"
	"strings"

	"github.com/luisnquin/nix-search/internal/nix"
	"github.com/microcosm-cc/bluemonday"
	"github.com/samber/lo"
)

type nixPackageResponseItem struct {
	Type          string   `json:"type"`
	Attr          string   `json:"package_attr_name"`
	AttrSet       string   `json:"package_attr_set"`
	Name          string   `json:"package_pname"`
	Version       string   `json:"package_pversion"`
	Platforms     []string `json:"package_platforms"`
	Outputs       []string `json:"package_outputs"`
	DefaultOutput string   `json:"package_default_output"`
	Programs      []string `json:"package_programs"`
	License       []struct {
		URL      string `json:"url"`
		FullName string `json:"fullName"`
	} `json:"package_license"`
	LicenseSet  []string `json:"package_license_set"`
	Maintainers []struct {
		Name   *string `json:"name"`
		Github *string `json:"github"`
		Email  string  `json:"email"`
	} `json:"package_maintainers"`
	MaintainersSet  []string    `json:"package_maintainers_set"`
	Description     string      `json:"package_description"`
	LongDescription string      `json:"package_longDescription"`
	Hydra           interface{} `json:"package_hydra"`
	System          string      `json:"package_system"`
	Homepage        []string    `json:"package_homepage"`
	Position        *string     `json:"package_position"`
}

type nixPackagesResponse = searchResponse[nixPackageResponseItem]

func (c Client) SearchPackages(ctx context.Context, channelBranch, searchTerm string, maxCount int) ([]*nix.Package, error) {
	response, err := c.searchPackages(ctx, channelBranch, searchTerm, maxCount)
	if err != nil {
		return nil, err
	}

	r := strings.NewReplacer("\n", " ")
	sp := bluemonday.StrictPolicy()

	pkgs := make([]*nix.Package, len(response.Hits.Items))

	for i, item := range response.Hits.Items {
		var homepage *string
		if len(item.Source.Homepage) > 0 {
			homepage = &item.Source.Homepage[0]
		}

		var license *nix.PackageLicense
		if len(item.Source.License) > 0 {
			license = &nix.PackageLicense{
				URL:      item.Source.License[0].URL,
				FullName: item.Source.License[0].FullName,
			}
		}

		maintainers := make([]*nix.PackageMaintainer, len(item.Source.Maintainers))
		for j, m := range item.Source.Maintainers {
			maintainers[j] = &nix.PackageMaintainer{
				Name:   m.Name,
				GitHub: m.Github,
				Email:  m.Email,
			}
		}

		pkgs[i] = &nix.Package{
			Name:               item.Source.Name,
			Pname:              item.Source.Attr,
			Description:        item.Source.Description,
			LongDescription:    strings.TrimSpace(r.Replace(sp.Sanitize(item.Source.LongDescription))),
			Version:            item.Source.Version,
			Set:                lo.ToPtr(item.Source.AttrSet),
			Programs:           item.Source.Programs,
			DefaultOutput:      item.Source.DefaultOutput,
			Outputs:            item.Source.Outputs,
			Platforms:          item.Source.Platforms,
			System:             item.Source.System,
			Homepage:           homepage,
			License:            license,
			Maintainers:        maintainers,
			RepositoryPosition: item.Source.Position,
			Query:              nix.PackageQuery{Score: item.Score},
		}
	}

	return pkgs, nil
}

func (c Client) searchPackages(ctx context.Context, channelBranch, searchTerm string, maxCount int) (nixPackagesResponse, error) {
	const query = // AQL
	`
	{
        "query": {
            "bool": {
                "filter": [
                    { 
                        "term": { 
                            "type": {
                                "value": "package", "_name": "filter_packages"
                            }
                        }
                    },
                    {
                        "bool": {
                            "must": [
                                { "bool": { "should": [] } },
                                { "bool": { "should": [] } },
                                { "bool": { "should": [] } },
                                { "bool": { "should": [] } }
                            ]
                        }
                    }
                ],
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
                                            "package_attr_name^9",
                                            "package_attr_name.*^5.3999999999999995",
                                            "package_programs^9",
                                            "package_programs.*^5.3999999999999995",
                                            "package_pname^6",
                                            "package_pname.*^3.5999999999999996",
                                            "package_description^1.3",
                                            "package_description.*^0.78",
                                            "package_longDescription^1",
                                            "package_longDescription.*^0.6",
                                            "flake_name^0.5",
                                            "flake_name.*^0.3"
                                        ]
                                    }
                                },
                                {
                                    "wildcard": {
                                        "package_attr_name": { "value": "*{{ . | snake_case }}*", "case_insensitive": true }
                                    }
                                },
                                {
                                    "wildcard": {
                                        "package_attr_name": { "value": "*{{ . | kebab_case }}*", "case_insensitive": true }
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

	esClient, err := c.prepareElasticSearchClient(channelBranch)
	if err != nil {
		return nixPackagesResponse{}, err
	}

	response, err := esClient.Search(
		esClient.Search.WithFrom(0),
		esClient.Search.WithSize(maxCount),
		esClient.Search.WithSort("_score:desc", "package_attr_name:desc", "package_pversion:desc"),
		esClient.Search.WithBody(prepareQuery(query, searchTerm)),
	)
	if err != nil {
		return nixPackagesResponse{}, err
	}

	defer response.Body.Close()

	if response.IsError() {
		return nixPackagesResponse{}, handleSearchErrorResponse(response.Body)
	}

	return parseSearchResponse[nixPackageResponseItem](response.Body)
}
