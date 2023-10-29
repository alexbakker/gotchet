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
      gotchetVersion = "0.2.0";
    in rec {
      packages = flake-utils.lib.flattenTree rec {
        default = gotchet-cli;
        gotchet-cli = with pkgs; buildGo121Module rec {
          pname = "gotchet-cli";
          version = gotchetVersion;
          src = ./.;

          subPackages = ["cmd/gotchet"];
          vendorSha256 = "sha256-gtBvrOpuh0WeGwyZfDlAMo4FujIFGVi09pgHl0VYqyM=";

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

          nativeBuildInputs = [ installShellFiles ];

          postInstall = ''
            for shell in bash fish zsh; do
              installShellCompletion --cmd gotchet --$shell <($out/bin/gotchet completion $shell)
            done
          '';

          checkPhase = ''
            go test -v $(go list ./... | grep -v /test)
            go test -json -v=test2json ./test/ -run GenerateFakeTree > test_output.json || true
            go run github.com/alexbakker/gotchet/cmd/gotchet generate -i test_output.json -o report.html
            ls -lah report.html
          '';
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
          go_1_21

          nodejs-18_x
          yarn
        ];
      };
    }
  );
}
