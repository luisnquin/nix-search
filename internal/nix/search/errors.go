package nix_search

import (
	"encoding/json"
	"fmt"
	"io"
)

type (
	ErrorResponseInSearchOperation struct {
		Error struct {
			RootCause []ErrorCausesItem `json:"root_cause"`
			Type      string            `json:"type"`
			Reason    string            `json:"reason"`
			Line      int               `json:"line"`
			Column    int               `json:"column"`
		} `json:"error"`
		Status int `json:"status"`
	}

	ErrorCausesItem struct {
		Type   string `json:"type"`
		Reason string `json:"reason"`
		Line   int    `json:"line"`
		Column int    `json:"col"`
	}
)

func handleSearchErrorResponse(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	var response ErrorResponseInSearchOperation

	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}

	return fmt.Errorf(`Elastic Search API response: {"status": %d, "type": "%s", "reason": "%s"}`,
		response.Status, response.Error.Type, response.Error.Reason)
}
