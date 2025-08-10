---
title: Configuration
---

# Configuration

Configuration can be loaded from a variety of sources, including files (`json`,
`toml`, `yaml`, `yml`, `hcl`), command-line flags, and environment variables
(prefixed with `RONG_`). Rong will automatically merge values from all sources.

## Locations

By default, the `rong` configuration file is located at:

```sh
"$XDG_CONFIG_HOME/rong/config.toml"
```

(Usually `~/.config/rong/config.toml`.)
The file can be in `json`, `toml`, `yaml`, `yml`, or `hcl` format.

::: tip
Use `rong --config /path/to/config` to load a custom config file. The extension determines the format.
:::

## Basic Structure

Here is a minimal example of a `toml` config:

```toml
variant = "expressive"
version = 2021
dark = true
```

Or the equivalent in `yaml`:

```yaml
variant: expressive
version: 2025
dark: true
```

### Fields

- `variant`: Defines the Material You color variant. Possible variants:

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

- `version`: Material You specification version (`2021` or `2025`)

- `dark`: Whether to generate a dark theme (`true`) or light theme (`false`).

::: info
If both `dark` and `light` are specified from different sources, the last one in the priority order above takes effect.
:::

## Linking Generated Files

The `[links]` section tells `rong` where to copy/hardlink each generated theme file. You can assign a single path or an array of paths if you want the same file copied/linked to multiple locations.

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

- Paths support `~` for your home directory.
- Existing files at the destination will be overwritten or replaced by symlinks.
- Template file names in `[links]` must match exactly with the template directory.

### Example

**FLAGS**

```bash
rong --version 2025 image path/to/image
```

**ENV**

```bash
RONG_DARK=true rong image path/to/image
```

**TOML**

```toml{6,14,15}
variant = "expressive"
version = 2021
dark = true

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

**YAML**

```yaml
variant: vibrant
version: 2025
dark: false
links:
  "kitty.conf": "~/.config/kitty/colors.conf"
  "gtk.css":
    - "~/.config/gtk-3.0/gtk.css"
    - "~/.config/gtk-4.0/gtk.css"
```

**JSON**

```json
{
  "variant": "monochrome",
  "version": 2021,
  "dark": true,
  "links": {
    "colors.lua": "~/.config/wezterm/colors.lua",
    "qtct.conf": [
      "~/.config/qt5ct/colors/rong.conf",
      "~/.config/qt6ct/colors/rong.conf"
    ]
  }
}
```
