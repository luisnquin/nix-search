package nix

type HomeManagerOption struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Note        string `json:"note"`
	Type        string `json:"type"`
	Default     string `json:"default"`
	Example     string `json:"example"`
	Position    string `json:"declared_by"`
}
