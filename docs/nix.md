# Home Manager

Rong has a experimental `Home-Manager` module. You can get this by flake.

::: warning
As Rong is experimental and module also have experimental issue.
:::

# Input

Add Rong flake to you flake input.

```nix{9,10,23}
{
    inputs = {
        nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
        home-manager = {
            url = "github:nix-community/home-manager";
            inputs.nixpkgs.follows = "nixpkgs";
        };

        rong.url = "github:Nadim147c/rong";
        rong.inputs.nixpkgs.follows = "nixpkgs";
    };

    outputs = {
        rong,
        ...
    }: let
        system = "x86_64-linux";
        pkgs = import nixpkgs { inherit system; };
    in {
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
{...}: {
    programs.rong = {
        enable = true;
        settings = {
            variant = "expressive";
            version = 2021;
            dark = true;
            links = {
                "hyprland.conf" = "~/.config/hypr/colors.conf";
                "colors.lua" = "~/.config/wezterm/colors.lua";
                "spicetify-sleek.ini" = "~/.config/spicetify/Themes/Sleek/color.ini";
                "kitty-full.conf" = "~/.config/kitty/colors.conf";
                "pywalfox.json" = "~/.cache/wal/colors.json";
                "gtk-css.css" = "~/.config/wlogout/colors.css";
                "rofi.rasi" = "~/.config/rofi/config.rasi";
                "ghostty" = "~/.config/ghostty/colors";
                "dunstrc" = "~/.config/dunst/dunstrc";
                "cava.ini" = "~/.config/cava/themes/rong.ini";

                "gtk.css" = [
                    "~/.config/gtk-3.0/gtk.css"
                    "~/.config/gtk-4.0/gtk.css"
                ];
            };
        };


        # Create or overwrite templates
        templates = {
            "cava.ini" = /* gotmpl */ ''
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
    };
}
```
