package nix

type Option struct {
	// The name of the NixOS option.
	Name string `json:"name"`
	// The description value of the NixOS option.
	Description string `json:"description"`
	// A usage example for the NixOS option.
	Example *string `json:"example"`
	// The default value of the NixOS option when not declared
	// but root triggered.
	Default string `json:"default"`
	// The place in https://github.com/NixOS/nixpkgs where
	// the NixOS option has been declared.
	Source *string `json:"source"`
}
