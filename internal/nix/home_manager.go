package nix

type HomeManagerOption struct {
	// The option name (e.g. programs.zsh.enable)
	Title string `json:"title"`
	// The option description.
	Description string `json:"description"`
	// An additional maintainer note about the option.
	Note string `json:"note"`
	// The expected(s) data type(s) of the option.
	Type string `json:"type"`
	// The default value of the option if allowed.
	Default string `json:"default"`
	// The example value of the option.
	Example string `json:"example"`
	// The repository file where the option was declared.
	Position string `json:"declared_by"`
}
