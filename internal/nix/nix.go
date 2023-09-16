package nix

// Source names.
const (
	FLAKE_PACKAGES = "flake-packages"
	FLAKE_OPTIONS  = "flake-options"
	NIXOS_OPTIONS  = "nixos-options"
	NIX_PACKAGES   = "nix-packages"
	HOME_OPTIONS   = "home-options"
)

const (
	CHANNEL_STATUS_STABLE       = "stable"
	CHANNEL_STATUS_ROLLING      = "rolling"
	CHANNEL_STATUS_UNMAINTAINED = "unmaintained"
)

func GetSourceNames() []string {
	return []string{
		FLAKE_PACKAGES,
		FLAKE_OPTIONS,
		NIXOS_OPTIONS,
		NIX_PACKAGES,
		HOME_OPTIONS,
	}
}
