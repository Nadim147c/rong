rong: {
  config,
  lib,
  pkgs,
  ...
}: let
  inherit (lib) types mkOption mkMerge mkEnableOption mkIf;
  pkg = pkgs.callPackage ./. {};
  cfg = config.programs.rong;
  format = pkgs.formats.json {};
in {
  options.programs.rong = {
    enable = mkEnableOption "Enable rong color generator";

    package = mkOption {
      type = types.package;
      default = pkg;
      description = "Rong package to use";
    };

    wallpaper = mkOption {
      type = types.nullOr types.path;
      default = null;
      description = "Wallpaper to use for generating color";
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
        description = "Map of theme files to target paths or list of paths";
      };
    };
  };

  config = mkIf cfg.enable {
    home.packages = mkIf (cfg.package != null) [cfg.package];

    home.activation.generateThemes = lib.mkIf (cfg.wallpaper != null) (
      let
        rong = "${cfg.package}/bin/rong";
        state = "${config.xdg.stateHome}/rong/image.txt";
      in
        lib.hm.dag.entryAfter ["writeBoundary"] ''
          if [ -f "${state}" ] && [ -f "$(cat "${state}")" ]; then
            ${rong} video "$(cat "${state}")"
          else
            ${rong} video "${cfg.wallpaper}"
          fi
        ''
    );
    xdg.configFile = mkMerge [
      (mkIf (cfg.settings != {}) {
        "rong/config.json" = {
          source = format.generate "rong.json" cfg.settings;
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
