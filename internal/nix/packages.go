package nix

// The metadata of a Nix package.
type (
	Package struct {
		// The package name.
		Name string `json:"name"`
		// The package name used as identifier.
		Pname string `json:"pname"`
		// The package description.
		Description     string  `json:"description"`
		LongDescription *string `json:"long_description"`
		// The package version.
		Version string `json:"version"`
		// The set to which the package belongs to.
		Set *string `json:"set"`
		// The programs provided by the package.
		Programs []string `json:"programs"`
		// The default output of the package.
		DefaultOutput string `json:"default_output"`
		// The list of package outputs.
		Outputs []string `json:"outputs"`
		// The list of supported platforms.
		Platforms []string `json:"platforms"`
		System    string   `json:"system"`
		// The homepage of the package.
		Homepage *string `json:"homepage"`
		// The license that the package is licensed under.
		License *PackageLicense `json:"licenses"`
		// The list of person who maintains the Nix package.
		Maintainers []*PackageMaintainer `json:"maintainers"`
		// The place in https://github.com/NixOS/nixpkgs where
		// the package has been declared.
		RepositoryPosition *string      `json:"repo_position"`
		Query              PackageQuery `json:"query"`
	}

	// Type               string             `json:"type"`

	PackageQuery struct {
		Score float64 `json:"score"`
	}

	PackageLicense struct {
		URL      string `json:"url"`
		FullName string `json:"full_name"`
	}

	PackageMaintainer struct {
		Name   *string `json:"name"`
		GitHub *string `json:"github"`
		Email  string  `json:"email"`
	}
)
