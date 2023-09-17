package nix

type Option struct {
	// The name of the NixOS option.
	Name string `json:"name"`
	// The description value of the NixOS option.
	Description string `json:"description"`

	LongDescription string `json:"long_description"`
	// The expected(s) data type(s) of the option.
	Type string `json:"type"`
	// The default value of the NixOS option when not declared
	// but root triggered.
	Default string `json:"default"`
	// A usage example for the NixOS option.
	Example *string `json:"example"`
	// The place in https://github.com/NixOS/nixpkgs where
	// the NixOS option has been declared.
	Source *string `json:"source"`
}
