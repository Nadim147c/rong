# Hyprland

Hyprland is a dynamic tiling Wayland compositor with modern features and flexible
configuration. It uses custom configuration language called **Hyprlang**.

## Link

Add the following line to the [configuration](/configuration#linking-generated-files):

```toml{3}
[links]
# ...
"hyprland.conf" = "~/.config/hypr/colors.conf"
```

## Apply

In your main config file (`~/.config/hypr/hyprland.conf`), source the theme file and
define your layout:

```kdl{1,5,6}
source = colors.conf

general {
    # ...
    col.active_border = $primary $secondary 45deg
    col.inactive_border = $on_primary $on_secondary 45deg
    # ...
}
```

## Reload

While Hyprland auto reloads, you might need to run:

```bash
hyprctl reload
```
