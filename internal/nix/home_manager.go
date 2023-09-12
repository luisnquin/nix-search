package nix

import (
	"encoding/json"
	"fmt"
	"strings"
)

type HomeManagerOption struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Note        string `json:"note"`
	Type        string `json:"type"`
	Default     string `json:"default"`
	Example     string `json:"example"`
	DeclaredBy  string `json:"declared_by"`
}

type otherOption HomeManagerOption

func (opt HomeManagerOption) String() string {
	return fmt.Sprintf("\n %s - %s\n Example: %s\n Default: %s\n", opt.Type, opt.Description, opt.Example, opt.Default)
}

func (option *HomeManagerOption) UnmarshalJSON(data []byte) error {
	var opt otherOption

	err := json.Unmarshal(data, &opt)
	if err != nil {
		return err
	}

	option.Description = strings.TrimSpace(opt.Description)
	option.DeclaredBy = strings.TrimSpace(opt.DeclaredBy)
	option.Default = strings.TrimSpace(opt.Default)
	option.Example = strings.TrimSpace(opt.Example)
	option.Title = strings.TrimSpace(opt.Title)
	option.Type = strings.TrimSpace(opt.Title)
	option.Note = strings.TrimSpace(opt.Note)

	return nil
}
