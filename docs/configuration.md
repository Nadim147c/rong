---
title: Configuration
---

# Configuration

Configuration can be loaded from a variety of sources, including files ( `json`,
`toml`, `yaml`, `yml`, `properties`, `props`, `prop`, `hcl`, `tfvars`, `dotenv`,
`env`, `ini`), command-line flags, and environment variables (prefixed with
`RONG_`). Rong will automatically merge values from all sources.

## Locations

By default, the `rong` configuration file is located at:

```sh
"$XDG_CONFIG_HOME/rong/config.toml"
```

(Usually `~/.config/rong/config.toml`.)
The file can be in `json`, `toml`, `yaml`, or `yml` format.

::: tip
Use `rong --config /path/to/config` to load a custom config file. The extension
determines the format.
:::

## Basic Structure

Here is a minimal example of a `toml` config:

```toml
dark = true

[material]
variant = "tonal_spot"
version = "2025"
```

Or the equivalent in `yaml`:

```yaml
dark: true
material:
  variant: tonal_spot
  version: "2025"
```

## Configuration Fields

### Core Settings

- `dark`: Generate dark color palette (`true`) or light palette (`false`)
- `dry_run`: Generate colors without applying templates
- `json`: Print generated colors as JSON to stdout
- `log_file`: File path to save logs
- `quiet`: Suppress all log output
- `verbose`: Verbose logging level (0-3, where 3 is most verbose)

### Material Design Settings

The `[material]` section controls Material You color generation:

- `variant`: Color variant to use. Possible variants:

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

- `version`: Material You specification version (`"2021"` or `"2025"`)
- `platform`: Target platform (`"phone"` or `"watch"`)
- `contrast`: Contrast adjustment (-1.0 to 1.0)

### Base16 Settings

The `[base16]` section controls Base16 color generation:

- `method`: Color generation method (`"static"` or `"dynamic"`)
- `blend`: Blend ratio toward the primary color (0.0 to 1.0)
- `colors`: Source colors for base16 color generation (all hex colors):
  - `black`, `blue`, `cyan`, `green`, `magenta`, `red`, `white`, `yellow`

## Linking Generated Files

The `[links]` section tells `rong` where to copy/hardlink each generated theme file.
You can assign a single path or an array of paths if you want the same file
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

- Paths support `~` for your home directory.
- Existing files at the destination will be overwritten or replaced by symlinks.
- Template file names in `[links]` must match exactly with the template directory.

## Examples

**TOML Configuration**

```toml
dark = true
verbose = 1

[material]
variant = "expressive"
version = "2025"
contrast = 0.2

[base16]
method = "dynamic"
blend = 0.7
colors.black = "#0a0a0a"
colors.blue = "#1e90ff"

[links]
"hyprland.conf" = "~/.config/hypr/colors.conf"
"colors.lua" = "~/.config/wezterm/colors.lua"
"gtk.css" = [
  "~/.config/gtk-3.0/gtk.css",
  "~/.config/gtk-4.0/gtk.css"
]
```

**YAML Configuration**

```yaml
dark: true
dry_run: false
verbose: 2

material:
  variant: vibrant
  version: "2025"
  platform: phone
  contrast: 0.1

base16:
  method: static
  blend: 0.5
  colors:
    black: "#000000"
    blue: "#0044FF"
    cyan: "#008080"
    green: "#008000"

links:
  "kitty.conf": "~/.config/kitty/colors.conf"
  "gtk.css":
    - "~/.config/gtk-3.0/gtk.css"
    - "~/.config/gtk-4.0/gtk.css"
```

**JSON Configuration**

```json
{
  "dark": true,
  "json": false,
  "material": {
    "variant": "monochrome",
    "version": "2021",
    "contrast": 0.0
  },
  "base16": {
    "method": "dynamic",
    "blend": 0.6,
    "colors": {
      "red": "#800000",
      "green": "#008000"
    }
  },
  "links": {
    "colors.lua": "~/.config/wezterm/colors.lua",
    "qtct.conf": [
      "~/.config/qt5ct/colors/rong.conf",
      "~/.config/qt6ct/colors/rong.conf"
    ]
  }
}
```

**Command Line Examples**

```bash
# Generate dark theme with expressive variant
rong --dark --material.variant expressive image path/to/image

# Dry run to preview colors as JSON
rong --dry-run --json image path/to/image

# Debug output with custom base16 colors
rong -vvv --base16.colors.red "#ff0000" image path/to/image
```

**Environment Variables**

```bash
RONG_DARK=true rong image path/to/image
RONG_MATERIAL_VARIANT=vibrant rong image path/to/image
RONG_BASE16_BLEND=0.8 rong image path/to/image
```
