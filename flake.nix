{
  description = "Nix flake for gotchet";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    flockenzeit.url = "github:balsoft/Flockenzeit";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, flake-utils, flockenzeit, nixpkgs }:
  flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = nixpkgs.legacyPackages.${system};
      gotchetVersion = "0.1.3";
    in rec {
      packages = flake-utils.lib.flattenTree rec {
        default = gotchet-cli;
        gotchet-cli = with pkgs; buildGoModule rec {
          pname = "gotchet-cli";
          version = gotchetVersion;
          src = ./.;

          CGO_ENABLED = 0;

          patchPhase = ''
            mkdir -p internal/report/dist
            cp "${gotchet-frontend}/index.html" internal/report/dist
          '';

          ldflags = let 
            pkgPath = "github.com/alexbakker/gotchet/cmd/gotchet/cmd";
          in [
            "-X ${pkgPath}.versionNumber=${version}"
            "-X ${pkgPath}.versionRevision=${self.shortRev or "dirty"}"
            "-X ${pkgPath}.versionRevisionTime=${toString self.lastModified}"
          ];

          subPackages = ["cmd/gotchet"];
          vendorSha256 = "sha256-w+R+iCnm8q5TjCS6Ov/M8kOP4sVJCBWCeKCq6sxNl+w=";
        };
        gotchet-docker = with pkgs; dockerTools.buildImage {
          name = "gotchet";
          tag = "latest";
          created = flockenzeit.lib.ISO-8601 self.lastModified;
          copyToRoot = gotchet-cli;
          config = {
            Entrypoint = [ "/bin/gotchet" ];
          };
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
      apps = rec {
        default = gotchet;
        gotchet = flake-utils.lib.mkApp {
          drv = self.packages.${system}.gotchet-cli;
          name = "gotchet";
          exePath = "/bin/gotchet";
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
