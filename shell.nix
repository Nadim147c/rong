{
  pkgs ? import <nixpkgs> { },
}:
pkgs.mkShell {
  name = "rong";
  # Get dependencies from the main package
  inputsFrom = [ (pkgs.callPackage ./nix/package.nix { }) ];
  # Additional tooling
  buildInputs = with pkgs; [
    go
    gopls
    golangci-lint
    gotestsum

    (pkgs.callPackage ./nix/go-enum.nix { })

    just
    just-lsp

    bun
    prettier

    nixfmt-tree
    nixfmt

    harper
  ];
}
