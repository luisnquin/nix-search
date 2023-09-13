package nix

import (
	"fmt"
)

type HomeManagerOption struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Note        string `json:"note"`
	Type        string `json:"type"`
	Default     string `json:"default"`
	Example     string `json:"example"`
	Position    string `json:"declared_by"`
}

func (opt HomeManagerOption) String() string {
	return fmt.Sprintf("\n %s - %s\n %s\n Example: %s\n Default: %s\n",
		opt.Title, opt.Description, opt.Position, opt.Example, opt.Default)
}
