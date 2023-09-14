package nix_search

import (
	"context"
	"strings"

	"github.com/luisnquin/nix-search/internal/nix"
	"github.com/microcosm-cc/bluemonday"
)

type (
	flakeOptionResponseItem struct {
		FlakeName         string        `json:"flake_name"`
		FlakeDescription  string        `json:"flake_description"`
		FlakeResolved     flakeResolved `json:"flake_resolved"`
		Revision          string        `json:"revision"`
		OptionName        string        `json:"option_name"`
		OptionDescription string        `json:"option_description"`
		OptionSource      *string       `json:"option_source"`
		OptionDefault     string        `json:"option_default"`
		OptionExample     *string       `json:"option_example"`
	}

	flakePackageResponseItem struct {
		FlakeName            string        `json:"flake_name"`
		FlakeDescription     string        `json:"flake_description"`
		FlakeResolved        flakeResolved `json:"flake_resolved"`
		Revision             string        `json:"revision"`
		PackageAttr          string        `json:"package_attr_name"`
		PackageAttrSet       *string       `json:"package_attr_set"`
		PackageName          string        `json:"package_pname"`
		PackageDescription   string        `json:"package_description"`
		PackageLongDesc      *string       `json:"package_longDescription"`
		PackageVersion       string        `json:"package_pversion"`
		PackagePlatforms     []string      `json:"package_platforms"`
		PackageOutputs       []string      `json:"package_outputs"`
		PackageDefaultOutput string        `json:"package_default_output"`
		PackagePrograms      []string      `json:"package_programs"`
		PackageLicense       []struct {
			URL      string `json:"url"`
			FullName string `json:"fullName"`
		} `json:"package_license"`
		PackageMaintainers []struct {
			Name   *string `json:"name"`
			GitHub *string `json:"github"`
			Email  string  `json:"email"`
		} `json:"package_maintainers"`
		PackageSystem   string   `json:"package_system"`
		PackageHomePage []string `json:"package_homepage"`
		PackagePosition *string  `json:"package_position"`
	}

	flakeResolved struct {
		Type  string `json:"type"`
		Owner string `json:"owner"`
		Repo  string `json:"repo"`
	}
)

type (
	flakeOptionsResponse  = searchResponse[flakeOptionResponseItem]
	flakePackagesResponse = searchResponse[flakePackageResponseItem]
)

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
			},
		}
	}

	return options, nil
}

func (c Client) SearchFlakePackages(ctx context.Context, flakesBranchId, searchTerm string, maxCount int) ([]*nix.FlakePackage, error) {
	response, err := c.searchFlakePackages(ctx, flakesBranchId, searchTerm, maxCount)
	if err != nil {
		return nil, err
	}

	packages := make([]*nix.FlakePackage, len(response.Hits.Items))

	for i, item := range response.Hits.Items {
		var homepage *string
		if len(item.Source.PackageHomePage) > 0 {
			homepage = &item.Source.PackageHomePage[0]
		}

		var license *nix.PackageLicense
		if len(item.Source.PackageLicense) > 0 {
			license = &nix.PackageLicense{
				URL:      item.Source.PackageLicense[0].URL,
				FullName: item.Source.PackageLicense[0].FullName,
			}
		}

		maintainers := make([]*nix.PackageMaintainer, len(item.Source.PackageMaintainers))
		for j, m := range item.Source.PackageMaintainers {
			maintainers[j] = &nix.PackageMaintainer{
				Name:   m.Name,
				GitHub: m.GitHub,
				Email:  m.Email,
			}
		}

		packages[i] = &nix.FlakePackage{
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
			Package: &nix.Package{
				Name:               item.Source.PackageName,
				Pname:              item.Source.PackageAttr,
				Description:        item.Source.PackageDescription,
				LongDescription:    item.Source.PackageLongDesc,
				Version:            item.Source.PackageVersion,
				Set:                item.Source.PackageAttrSet,
				Programs:           item.Source.PackagePrograms,
				DefaultOutput:      item.Source.PackageDefaultOutput,
				Outputs:            item.Source.PackageOutputs,
				Platforms:          item.Source.PackagePlatforms,
				System:             item.Source.PackageSystem,
				Homepage:           homepage,
				License:            license,
				Maintainers:        maintainers,
				RepositoryPosition: item.Source.PackagePosition,
				Query:              nix.PackageQuery{Score: item.Score},
			},
		}
	}

	return packages, nil
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

func (c Client) searchFlakePackages(ctx context.Context, flakesBranchId, searchTerm string, maxCount int) (flakePackagesResponse, error) {
	const query = // AQL
	`
	{
		"query": {
			"bool": {
				"filter": [
					{ "term": { "type": { "value": "package", "_name": "filter_packages" } } },
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
										"_name": "multi_match_{{ . | snake_case }}",
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
										"package_attr_name": { "value": "*{{ . }}*", "case_insensitive": true }
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

	esClient, err := c.prepareElasticSearchClient(flakesBranchId)
	if err != nil {
		return flakePackagesResponse{}, err
	}

	response, err := esClient.Search(
		esClient.Search.WithFrom(0),
		esClient.Search.WithSize(maxCount),
		esClient.Search.WithSort("_score:desc", "option_name:desc"),
		esClient.Search.WithBody(prepareQuery(query, searchTerm)),
	)
	if err != nil {
		return flakePackagesResponse{}, err
	}

	defer response.Body.Close()

	if response.IsError() {
		return flakePackagesResponse{}, handleSearchErrorResponse(response.Body)
	}

	return parseSearchResponse[flakePackageResponseItem](response.Body)
}
