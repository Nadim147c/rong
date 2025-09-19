{
  description = "A Materal You and Base16 color generator";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default-linux";
  };
  outputs = {
    self,
    systems,
    nixpkgs,
  }: let
    forAllSystems = nixpkgs.lib.genAttrs (import systems);
    pkgsFor = nixpkgs.legacyPackages;
  in {
    packages = forAllSystems (system: {
      default = pkgsFor.${system}.callPackage ./. {};
    });
    devShells = forAllSystems (system: {
      default = pkgsFor.${system}.callPackage ./shell.nix {};
    });
    homeModules = {
      rong = import ./module.nix self;
      default = self.homeModules.rong;
    };
  };
}
