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

          CGO_ENABLED = 0;

          patchPhase = ''
            mkdir -p internal/report/dist
            cp "${gotchet-frontend}/index.html" internal/report/dist
          '';

          subPackages = ["cmd/gotchet"];
          vendorSha256 = "sha256-Ia9s5bCVdcG6QijEcA3h5IkEVPsLf/kzV1UBElk1lLQ=";
        };
        gotchet-frontend = with pkgs; mkYarnPackage rec {
          name = "gotchet-frontend";
          src = ./internal/report;
          packageJSON = src + "/package.json";
          yarnLock = src + "/yarn.lock";

          buildPhase = ''
            export HOME=$(mktemp -d)
            yarn --offline build
          '';

          installPhase = ''
            mkdir -p $out
            cp -r deps/${name}/dist/* $out
          '';

          doDist = false;
        };
      };
      devShells.default = with pkgs; mkShell {
        hardeningDisable = [ "fortify" ];
        buildInputs = [
          go

          nodejs-18_x
          yarn
        ];
      };
    }
  );
}
