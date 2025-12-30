{ pkgs, ... }:
pkgs.mkShell {
  name = "rong";
  # Get dependencies from the main package
  inputsFrom = [ (pkgs.callPackage ./package.nix { }) ];
  # Additional tooling
  buildInputs = with pkgs; [
    go
    gopls
    golangci-lint
    gotestsum

    just
    just-lsp

    bun
    prettier

    nixfmt-tree
    nixfmt-rfc-style
  ];
}
