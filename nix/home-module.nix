self: {
  config,
  lib,
  pkgs,
  ...
}: let
  inherit
    (lib)
    types
    mkOption
    mkMerge
    mkEnableOption
    mkIf
    ;
  inherit
    (types)
    attrsOf
    bool
    either
    enum
    float
    int
    listOf
    nullOr
    package
    path
    str
    ;

  pkg = self.packages.${pkgs.stdenv.hostPlatform.system}.default;
  cfg = config.programs.rong;
  format = pkgs.formats.json {};
in {
  options.programs.rong = {
    enable = mkEnableOption "Enable rong color generator";

    package = mkOption {
      type = package;
      default = pkg;
      description = "Rong package to use";
    };

    wallpaper = mkOption {
      type = nullOr path;
      default = null;
      description = "Wallpaper to use for generating color";
    };

    templates = mkOption {
      type = attrsOf str;
      default = {};
      description = "Templates to use";
    };

    settings = {
      dark = mkOption {
        type = bool;
        default = false;
        description = "Generate dark color palette";
      };

      preview-format = mkOption {
        type = str;
        default = false;
        description = "Format of the image preview generate for video input";
      };

      dry-run = mkOption {
        type = bool;
        default = false;
        description = "Generate colors without applying templates";
      };

      json = mkOption {
        type = bool;
        default = false;
        description = "Print generated colors as JSON";
      };

      log_file = mkOption {
        type = nullOr str;
        default = null;
        description = "File to save logs";
      };

      quiet = mkOption {
        type = bool;
        default = false;
        description = "Suppress all logs";
      };

      frames = mkOption {
        type = int;
        default = 5;
        description = "Number of frames to process for video";
      };

      verbose = mkOption {
        type = int;
        default = 0;
        description = "Enable verbose logging (0-3)";
        example = "2";
      };

      # Base16 subtree matching the schema
      base16 = {
        blend = mkOption {
          type = float;
          default = 0.5;
          description = "Blend ratio toward the primary color (0..1)";
        };

        colors = mkOption {
          type = attrsOf str;
          default = {
            black = "#000000";
            blue = "#0044FF";
            cyan = "#008080";
            green = "#008000";
            magenta = "#800080";
            red = "#800000";
            white = "#C0C0C0";
            yellow = "#808000";
          };
          description = ''
            Source colors for base16 color generation. Expected hex strings
            like "#RRGGBB".
          '';
        };

        method = mkOption {
          type = enum ["static" "dynamic"];
          default = "static";
          description = "Color generation method";
        };
      };

      # Material subtree (note: renamed from 'materail' to 'material')
      material = {
        custom = {
          blend = mkOption {
            type = bool;
            default = true;
            description = "Whether or not blend custom colors toward primary";
          };
          ratio = mkOption {
            type = float;
            default = 0.3;
            description = "Whether or not blend custom colors toward primary";
          };
          colors = mkOption {
            type = attrsOf str;
            default = 0.3;
            description = "Whether or not blend custom colors toward primary";
          };
        };

        contrast = mkOption {
          type = float;
          default = 0.0;
          description = "Contrast adjustment (-1.0 .. 1.0)";
        };

        platform = mkOption {
          type = enum ["phone" "watch"];
          default = "phone";
          description = "Target platform";
        };

        variant = mkOption {
          type = enum [
            "monochrome"
            "expressive"
            "vibrant"
            "neutral"
            "fidelity"
            "rainbow"
            "tonal_spot"
            "content"
            "fruit_salad"
          ];
          default = "tonal_spot";
          description = "Color variant to use";
        };

        version = mkOption {
          type = enum ["2021" "2025"];
          default = "2025";
          description = "Version of the theme (2021 or 2025)";
        };
      };

      links = mkOption {
        type = attrsOf (either str (listOf str));
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

    home.activation.generateRongThemes = lib.mkIf (cfg.wallpaper != null) (
      let
        rong = "${cfg.package}/bin/rong";
        state = "${config.xdg.stateHome}/rong/state.json";
      in
        lib.hm.dag.entryAfter ["checkLinkTargets"] ''
          if [ -f "${state}" ]; then
            run --silence ${rong} regen $VERBOSE_ARG
          else
            run --silence ${rong} video $VERBOSE_ARG "${cfg.wallpaper}"
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
            "rong/templates/${name}".text = text;
          })
          cfg.templates)
      ))
    ];
  };
}
