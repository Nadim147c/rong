---
title: Configuration
---

# Configuration

This document explains how to configure `rong` using the `config.toml` file.
This file controls how `rong` generates themes, which variant to use, and where to
link the generated files.

## Locations

The `rong` configuration file should be`$XDG_CONFIG_HOME/rong/config.toml` (Usually
`~/.config/rong/config.toml`). The configuration file is in `toml` format.

::: tip
Use `rong --config /path/to/config.toml` to load a custom config file.
:::

## Basic Structure

Here is a minimal example of a `config.toml`:

```toml
variant = "expressive"
version = 2021
dark = true
```

### Fields

- `variant`: Defines the Material You color variant. Possible variant:
  <table>
    <tbody>
      <tr>
        <td>monochrome</td>
        <td>expressive</td>
        <td>vibrant</td>
      </tr>
      <tr>
        <td>neutral</td>
        <td>fidelity</td>
        <td>rainbow</td>
      </tr>
      <tr>
        <td>tonal_spot</td>
        <td>content</td>
        <td>fruit_salad</td>
      </tr>
    </tbody>
  </table>

- `version`: Material You specification version. (`2021` or `2025`)
- `light`: Whether to generate a light theme (`true`) or dark theme (`false`).

## Linking Generated Files

The `[links]` section tells `rong` where to copy/hardlink each generated theme file.
You can either assign a single path or an array of paths if you want the same file
copied/linked to multiple locations.

### Syntax

```toml{3}
[links]
"template_name.ext" = "destination_path"
"template_name.ext" = ["path1", "path2", ...]
```

::: info IMPORTANT
Each key must match the name of a template in the theme templates directory.
:::

### Notes

- All paths support `~` for your home directory.
- Existing files at the destination will be overwritten or replaced by symlinks.
- Make sure you’ve named your template files exactly as the keys in `[links]`.

### Example

```toml{6,14,15}
variant = "expressive"
version = 2021
light = false

[links]
"hyprland.conf" = "~/.config/hypr/colors.conf"
"colors.lua" = "~/.config/wezterm/colors.lua"
"spicetify-sleek.ini" = "~/.config/spicetify/Themes/Sleek/color.ini"
"kitty.conf" = "~/.config/kitty/colors.conf"
"pywalfox.json" = "~/.cache/wal/colors.json"

"colors.scss" = [ "~/.config/eww/colors.scss" ]
"qtct.conf" = [
  "~/.config/qt5ct/colors/rong.conf",
  "~/.config/qt6ct/colors/rong.conf"
]
"gtk.css" = [
  "~/.config/gtk-3.0/gtk.css",
  "~/.config/gtk-4.0/gtk.css"
]
"gtk-css.css" = [
  "~/.config/wlogout/colors.css"
]
"midnight-discord.css" = [
  "~/.config/vesktop/settings/quickCss.css"
]
```
