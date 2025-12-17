self:
{
  config,
  lib,
  pkgs,
  ...
}:
let
  inherit (lib)
    types
    mkOption
    mkMerge
    mkEnableOption
    mkIf
    ;
  inherit (types)
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
  format = pkgs.formats.json { };
in
{
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
      default = { };
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
          type = enum [
            "static"
            "dynamic"
          ];
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
          type = enum [
            "phone"
            "watch"
          ];
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
          type = enum [
            "2021"
            "2025"
          ];
          default = "2025";
          description = "Version of the theme (2021 or 2025)";
        };
      };

      installs = mkOption {
        type = attrsOf (either str (listOf str));
        default = { };
        example = ''
          {
            "quickshell.json" = "~/.config/quickshell/material.json";
          }
        '';
        description = ''
          Map of theme files to target paths or list of paths.

          These are installed using **atomic copy**, ensuring that applications
          never read partially-written or incomplete theme files. This method is
          safer for programs that load configuration very fast — for example,
          *Quickshell*, which can break if a theme file is replaced mid-read.

          Use `installs` when the target application requires atomic,
          fully-written files and cannot tolerate partial updates.
        '';
      };

      links = mkOption {
        type = attrsOf (either str (listOf str));
        default = { };
        example = ''
          {
            "hyprland.conf" = "~/.config/hypr/colors.conf";
            "colors.lua" = "~/.config/wezterm/colors.lua";
            "spicetify-sleek.ini" = "~/.config/spicetify/Themes/Sleek/color.ini";
          }
        '';
        description = ''
          Map of theme files to target paths or list of paths.

          These are installed using **hardlinks** whenever possible. If a
          hardlink already exists, only the file timestamps are updated —
          avoiding unnecessary writes and reducing filesystem churn. This method
          is efficient and ideal for most applications.

          Users should **prefer `links`** because it results in less disk I/O,
          and programs like break when reading files that are updated via atomic
          replacement. Only use `installs` when atomicity is required.
        '';
      };
    };
  };

  config = mkIf cfg.enable {
    home.packages = mkIf (cfg.package != null) [ cfg.package ];

    home.activation.generateRongThemes = lib.mkIf (cfg.wallpaper != null) (
      let
        rong = "${cfg.package}/bin/rong";
        state = "${config.xdg.stateHome}/rong/state.json";
      in
      lib.hm.dag.entryAfter [ "checkLinkTargets" ] ''
        if [ -f "${state}" ]; then
          run --silence ${rong} regen $VERBOSE_ARG
        else
          run --silence ${rong} video $VERBOSE_ARG "${cfg.wallpaper}"
        fi
      ''
    );

    xdg.configFile = mkMerge [
      (mkIf (cfg.settings != { }) {
        "rong/config.json" = {
          source = format.generate "rong.json" cfg.settings;
        };
      })

      (mkIf (cfg.templates != { }) (
        mkMerge (
          lib.mapAttrsToList (name: text: {
            "rong/templates/${name}".text = text;
          }) cfg.templates
        )
      ))
    ];
  };
}
