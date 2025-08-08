rong: {
  config,
  lib,
  pkgs,
  ...
}: let
  inherit (lib) types mkOption mkMerge mkEnableOption mkIf;
  pkg = pkgs.callPackage ./. {};
  cfg = config.programs.rong;
  format = pkgs.formats.yaml {};
in {
  options.programs.rong = {
    enable = mkEnableOption "Enable rong color generator";

    package = mkOption {
      type = types.package;
      default = pkg;
      description = "Rong package to use";
    };

    templates = mkOption {
      type = types.attrsOf types.str;
      default = {};
      description = "Templates to use";
    };

    settings = {
      variant = mkOption {
        type = types.enum [
          "monochrome"
          "neutral"
          "tonal_spot"
          "vibrant"
          "expressive"
          "fidelity"
          "content"
          "rainbow"
          "fruit_salad"
        ];
        default = "expressive";
        description = "Theme variant";
      };

      version = mkOption {
        type = types.enum [2021 2025];
        default = 2021;
        example = "2021";
        description = "Theme version";
      };

      dark = mkOption {
        type = types.bool;
        default = false;
        description = "Light or dark theme";
      };

      links = mkOption {
        type = types.attrsOf (types.either types.str (types.listOf types.str));
        default = {};
        example = ''
          {
            "hyprland.conf" = "~/.config/hypr/colors.conf";
            "colors.lua" = "~/.config/wezterm/colors.lua";
            "spicetify-sleek.ini" = "~/.config/spicetify/Themes/Sleek/color.ini";
            "kitty.conf" = "~/.config/kitty/colors.conf";
            "pywalfox.json" = "~/.cache/wal/colors.json";
            "rofi.rasi" = "~/.config/rofi/config.rasi";
            "ghostty" = "~/.config/ghostty/colors";
            "dunstrc" = "~/.config/dunst/dunstrc";
          }
        '';
        description = "Map of theme files to target paths or list of paths.";
      };
    };
  };

  config = mkIf cfg.enable {
    home.packages = mkIf (cfg.package != null) [cfg.package];
    xdg.configFile = mkMerge [
      (mkIf (cfg.settings != {}) {
        "rong/config.yaml" = {
          source = format.generate "rong.yaml" cfg.settings;
        };
      })

      (mkIf (cfg.templates != {}) (
        mkMerge (lib.mapAttrsToList (name: text: {
            "rong/templates/${name}" = {
              text = text;
            };
          })
          cfg.templates)
      ))
    ];
  };
}
