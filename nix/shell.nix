{pkgs ? import <nixpkgs> {}}:
pkgs.mkShell {
  name = "rong";
  # Get dependencies from the main package
  inputsFrom = [(pkgs.callPackage ./package.nix {})];
  # Additional tooling
  buildInputs = with pkgs; [
    go
    gopls
    gofumpt
    golines
    revive
    gnumake
    findutils

    bun
    prettier

    alejandra
  ];
}
