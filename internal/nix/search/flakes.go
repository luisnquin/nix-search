package nix_search

type flakeResolved struct {
	Type  string `json:"type"`
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}
