package nix

import "fmt"

type (
	FlakePackage struct {
		Flake    *FlakeMetadata `json:"flake"`
		Revision string         `json:"revision"`
		*Package `json:"package"`
	}

	FlakeOption struct {
		Flake    *FlakeMetadata `json:"flake"`
		Revision string         `json:"revision"`
		*Option  `json:"option"`
	}

	FlakeMetadata struct {
		Name        string      `json:"name"`
		Description string      `json:"description"`
		Origin      FlakeOrigin `json:"origin"`
	}

	FlakeOrigin struct {
		Type  string `json:"type"`
		Owner string `json:"owner"`
		Repo  string `json:"repo"`
	}
)

func (fp FlakePackage) String() string {
	return fmt.Sprintf(`{"flake": "%s", "revision": "%s", "package": "%s"}`, fp.Flake.Name, fp.Revision, fp.Package.Name)
}

func (fo FlakeOption) String() string {
	return fmt.Sprintf(`{"flake": "%s", "revision": "%s", "option": "%s"}`, fo.Flake.Name, fo.Revision, fo.Option.Name)
}
