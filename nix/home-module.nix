self: {
  config,
  lib,
  pkgs,
  ...
}: let
  inherit (lib) types mkOption mkMerge mkEnableOption mkIf;
  pkg = self.packages.${pkgs.stdenv.hostPlatform.system}.default;
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
      type = types.attrs;
      default = {};
      description = "Rong settings";
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
