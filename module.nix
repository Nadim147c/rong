{
  config,
  lib,
  pkgs,
  ...
}: let
  inherit (lib) types mkOption mkEnableOption mkIf;
  cfg = config.programs.rong;
in {
  options.programs.rong = {
    enable = mkEnableOption "Enable rong color generator";

    options = {
      variant = mkOption {
        type = types.str;
        default = "expressive";
        description = "Theme variant";
      };

      version = mkOption {
        type = types.int;
        default = 2021;
        description = "Theme version";
      };

      light = mkOption {
        type = types.bool;
        default = false;
        description = "Light or dark theme";
      };

      links = mkOption {
        type = types.attrsOf (types.either types.str (types.listOf types.str));
        default = {};
        description = "Map of theme files to target paths or list of paths.";
      };
    };
  };

  config = mkIf cfg.enable {
    home.file.".config/rong/config.yml".text = pkgs.lib.generators.toYAML cfg.options;
  };
}
