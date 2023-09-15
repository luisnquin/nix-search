package internal

import _ "embed"

const (
	PROGRAM_NAME = "nix-search"

	PROGRAM_DESCRIPTION = "Search nix stuff in your terminal"
)

const DEFAULT_VERSION = "development"

//go:embed help.tpl
var helpTpl string

func GetHelpTemplate() string {
	return helpTpl
}
