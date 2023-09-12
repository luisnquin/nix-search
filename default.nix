{
  pkgs ? import <nixpkgs> {},
  elasticSearchMappingVersion,
  nixosChannels,
}:
pkgs.buildGoModule rec {
  name = "nix-search";
  src = builtins.path {
    inherit name;
    path = ./.;
  };

  preConfigure = let
    internal-config = {
      nix =
        {
          flakes_id = "group-manual";
          sources = {
            home_manager_options = {
              url = "https://mipmip.github.io/home-manager-option-search/data/options.json";
            };

            elastic_search = {
              url = "https://search.nixos.org/backend";
              username = "aWVSALXpZv";
              password = "X8gPHnzL52wFEekuxsfQ9cSh";
              mapping_version = toString elasticSearchMappingVersion;
            };
          };
        }
        // nixosChannels;
    };
  in ''
    echo '${builtins.toJSON internal-config}' > $PWD/internal/config/internal.config.json
  '';

  vendorHash = "sha256-UBqN1SGbi9aQdFJZie2gV5oQimm5s8lFVNIFyptE6Qk=";
  doCheck = false; # TODO: remove
}
