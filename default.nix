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

  vendorHash = "sha256-UBqN1SGbi9aQdFJZie2gV5oQimm5s8lFVNIFyptE6Qk=";
  doCheck = false; # TODO: remove
}
