# Home Manager

Rong has a `Home-Manager` module. You can get this by flake.

# Input

Add Rong flake to your flake input.

```nix{8-11,25}
{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    home-manager = {
      url = "github:nix-community/home-manager";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    rong = {
      url = "github:Nadim147c/rong";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    {
      nixpkgs,
      rong,
      ...
    }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {
      homeConfigurations."<username>" = home-manager.lib.homeManagerConfiguration {
        inherit pkgs;
        modules = [
          rong.homeModules.default
          ./home
        ];
      };
    };
}
```

# Module

Now, use this module anywhere in your configuration.

```nix
{ ... }:
{
  programs.rong = {
    enable = true;
    settings = {
      dark = true;
      base16 = {
        blend = 0.5;
        method = "static";
        colors.green = "#00FF00";
      };
      material = {
        contrast = 0.0;
        platform = "phone";
        variant = "tonal_spot";
        version = "2025";
      };
      post-cmds."pywalfox.json" = /* bash */ "pywalfox --verbose update";
      installs = {
        "quickshell.json" = "${xdg.stateHome}/quickshell/colors.json";
      };
      links = {
        "hyprland.conf" = "~/.config/hypr/colors.conf";
        "colors.lua" = "~/.config/wezterm/colors.lua";
        "spicetify-sleek.ini" = "~/.config/spicetify/Themes/Sleek/color.ini";
        "kitty-full.conf" = "~/.config/kitty/colors.conf";
        "pywalfox.json" = "~/.cache/wal/colors.json";
        "qtct.colors" = [
          "~/.config/qt5ct/colors/rong.colors"
          "~/.config/qt6ct/colors/rong.colors"
        ];
      };
    };

    # Create or overwrite templates
    templates."cava.ini" = /* ini */ ''
      [color]
      ; background = '{{ .Background }}'
      background = 'default'

      ; gradient = 0
      gradient = 1
      gradient_color_1 = '{{ .Color1 }}'
      gradient_color_2 = '{{ .Color2 }}'
      gradient_color_3 = '{{ .Color3 }}'
      gradient_color_4 = '{{ .Color4 }}'
      gradient_color_5 = '{{ .Color5 }}'
      gradient_color_6 = '{{ .Color6 }}'

      ; horizontal_gradient = 0
      horizontal_gradient = 1
      horizontal_gradient_color_1 = '{{ .Color1 }}'
      horizontal_gradient_color_2 = '{{ .Color2 }}'
      horizontal_gradient_color_3 = '{{ .Color3 }}'
      horizontal_gradient_color_4 = '{{ .Color4 }}'
      horizontal_gradient_color_5 = '{{ .Color5 }}'
      horizontal_gradient_color_6 = '{{ .Color6 }}'
    '';
  };
}
```
