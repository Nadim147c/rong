{pkgs ? import <nixpkgs> {}}:
pkgs.mkShell {
  name = "rong";
  # Get dependencies from the main package
  inputsFrom = [(pkgs.callPackage ./default.nix {})];
  # Additional tooling
  buildInputs = with pkgs; [
    go
    gopls
    gofumpt
    revive
    gnumake
    bun
    prettier
  ];
}
