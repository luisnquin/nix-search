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
    config =
      {
        elastic_search = {
          host = "https://search.nixos.org/backend";
          username = "aWVSALXpZv";
          password = "X8gPHnzL52wFEekuxsfQ9cSh";
          mapping_version = toString elasticSearchMappingVersion;
        };

        home_manager = {
          data_url = "https://mipmip.github.io/home-manager-option-search/data/options.json";
        };
      }
      // nixosChannels;
  in ''
    echo '${builtins.toJSON config}' > $PWD/internal/config/internal.config.json
  '';

  vendorHash = "sha256-UBqN1SGbi9aQdFJZie2gV5oQimm5s8lFVNIFyptE6Qk=";
  doCheck = false; # TODO: remove
}
