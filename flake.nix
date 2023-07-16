{
  description = "Nix flake for gotchet";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, flake-utils, nixpkgs }:
  flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = nixpkgs.legacyPackages.${system};
    in rec {
      packages = flake-utils.lib.flattenTree rec {
        default = gotchet;
        gotchet = with pkgs; buildGoModule rec {
          name = "gotchet";
          src = ./.;

          subPackages = ["cmd/gotchet"];
          vendorSha256 = "sha256-Ia9s5bCVdcG6QijEcA3h5IkEVPsLf/kzV1UBElk1lLQ=";
        };
      };
      devShells.default = with pkgs; mkShell {
        hardeningDisable = [ "fortify" ];
        buildInputs = [
          go
        ];
      };
    }
  );
}
