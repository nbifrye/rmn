{
  description = "rmn - Redmine CLI";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        packages.default = pkgs.buildGoModule {
          pname = "rmn";
          version = self.shortRev or self.dirtyShortRev or "dev";
          src = self;

          # Run `nix build` once to get the correct hash from the error message,
          # then replace this value.
          vendorHash = null;

          ldflags = [
            "-s" "-w"
            "-X main.version=${self.shortRev or "dev"}"
            "-X main.commit=${self.shortRev or "none"}"
            "-X main.date=1970-01-01T00:00:00Z"
          ];

          subPackages = [ "cmd/rmn" ];

          meta = with pkgs.lib; {
            description = "CLI tool for interacting with Redmine";
            homepage = "https://github.com/nbifrye/rmn";
            license = licenses.mit;
            mainProgram = "rmn";
          };
        };
      }
    );
}
