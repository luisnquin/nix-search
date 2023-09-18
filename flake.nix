{
  inputs = {
    nixos-org-configurations = {
      url = "github:NixOS/nixos-org-configurations";
      flake = false;
    };

    nixos-search = {
      url = "github:NixOS/nixos-search";
      flake = false;
    };

    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = inputs:
    with inputs;
      flake-utils.lib.eachDefaultSystem
      (
        system: let
          pkgs = import nixpkgs {
            inherit system;
          };

          inherit (nixpkgs) lib;

          nixElasticSearch = let
            allChannels = (import "${nixos-org-configurations}/channels.nix").channels;
            filteredChannels =
              lib.filterAttrs
              (
                n: v:
                  builtins.elem v.status ["rolling" "beta" "stable" "deprecated"]
                  && lib.hasPrefix "nixos-" n
                  && v ? variant
                  && v.variant == "primary"
              )
              allChannels;
          in {
            mappingVersion = toString (lib.fileContents "${nixos-search}/VERSION");

            channels =
              lib.mapAttrsToList
              (
                n: v: {
                  id = lib.removePrefix "nixos-" n;
                  status = v.status;
                  jobset =
                    builtins.concatStringsSep
                    "/"
                    (lib.init (lib.splitString "/" v.job));
                  branch = n;
                }
              )
              filteredChannels;
            defaultChannel =
              builtins.head
              (
                builtins.sort (e1: e2: ! (builtins.lessThan e1 e2))
                (
                  builtins.map
                  (lib.removePrefix "nixos-")
                  (
                    builtins.attrNames
                    (lib.filterAttrs (_: v: v.status == "stable") filteredChannels)
                  )
                )
              );
          };
        in rec {
          packages.default = packages.nix-search;
          packages.nix-search = import ./default.nix {
            inherit pkgs nixElasticSearch;
          };
        }
      );
}
