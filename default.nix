{
  pkgs ? import <nixpkgs> {},
  nixElasticSearch ? null,
}:
pkgs.buildGoModule rec {
  name = "nix-search";
  src = builtins.path {
    inherit name;
    path = ./.;
  };

  preConfigure = let
    settingsPath = "$PWD/internal/config/settings.json";
    hasStrValue = s: s != null && builtins.isString s && s != "";

    jsonValueReplacer = jsonPath: jsonValue: "${pkgs.jq}/bin/jq -r '${jsonPath} |= ${jsonValue}' ${settingsPath} | ${pkgs.moreutils}/bin/sponge ${settingsPath}";
  in
    if nixElasticSearch != null
    then ''
      ${
        if hasStrValue nixElasticSearch.defaultChannel
        then jsonValueReplacer ".nix.default_channel" "\"${nixElasticSearch.defaultChannel}\""
        else ""
      }

      ${
        if (builtins.length nixElasticSearch.channels) != 0
        then jsonValueReplacer ".nix.channels" (builtins.toJSON nixElasticSearch.channels)
        else ""
      }

      ${
        if hasStrValue nixElasticSearch.mappingVersion
        then jsonValueReplacer ".nix.sources.elastic_search.mapping_version" "\"${nixElasticSearch.mappingVersion}\""
        else ""
      }
    ''
    else "";

  vendorHash = "sha256-08GOTAgpJE3pEcTRb+3hQX1UkSPegoug33GOAkBIYKo=";
  doCheck = true;
}
