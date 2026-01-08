# Hyprland

Hyprland is a dynamic tiling Wayland compositor with modern features and flexible
configuration. It uses custom configuration language called **Hyprlang**.

## Configuration

Add the following line to the [configuration](/configuration#linking-generated-files):

```toml{3}
[[themes]]
target = "hyprland.conf"
links = "~/.config/hypr/colors.conf"
cmds = "hyprctl reload"
```

<!--@include: ./_regen.md-->

## Apply

In your main config file (`~/.config/hypr/hyprland.conf`), source the theme file and
define your layout:

```txty1,5,6y
source = colors.conf

general {
    # ...
    col.active_border = $primary $secondary 45deg
    col.inactive_border = $on_primary $on_secondary 45deg
    # ...
}
```
