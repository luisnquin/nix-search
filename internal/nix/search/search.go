package nix_search

import (
	"encoding/json"
	"io"
)

type (
	searchResponse[responseItem any] struct {
		Took     uint32 `json:"took"`
		TimedOut bool   `json:"timed_out"`
		// Shards   any    `json:"_shards"`
		Hits searchHits[responseItem] `json:"hits"`
	}

	searchHits[responseItem any] struct {
		Total struct {
			Value    uint32 `json:"value"`
			Relation string `json:"eq"`
		} `json:"total"`
		MaxScore *float64                        `json:"max_score"`
		Items    []*searchHitsItem[responseItem] `json:"hits"`
	}

	searchHitsItem[responseItem any] struct {
		Index          string        `json:"_index"`
		Type           string        `json:"_type"`
		ID             string        `json:"_id"`
		Score          float64       `json:"_score"`
		Source         responseItem  `json:"_source"`
		Sort           []interface{} `json:"sort"`
		MatchedQueries []string      `json:"matched_queries"`
	}
)

func parseSearchResponse[responseItem any](r io.Reader) (searchResponse[responseItem], error) {
	var response searchResponse[responseItem]

	if err := json.NewDecoder(r).Decode(&response); err != nil {
		return searchResponse[responseItem]{}, err
	}

	return response, nil
}
