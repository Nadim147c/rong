# Hyprlock

Hyprland's GPU-accelerated screen locking utility. Just like Hyprland, it also uses
Hyprlang for it's configuration.

## Link

Add the following line to the [configuration](/configuration#linking-generated-files):

```toml{3}
[links]
# ...
"hyprland.conf" = "~/.config/hypr/colors.conf"
```

<!--@include: ./_regen.md-->

## Apply

In your main config file (`~/.config/hypr/hyprlock.conf`), source the theme file and
define your layout:

```txt{1,4,5,10}
source = colors.conf

background {
    path = $image
    color = $background
    # ...
}

input-field {
    placeholder_text = <i><span foreground="$primary_hex">Password</span></i>
    # ...
}
```

## Reload

You don't need to reload Hyprlock as it will auto load the new theme on next startup.
