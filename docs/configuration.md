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
- `dry-run`: Generate colors without applying templates
- `json`: Print generated colors as JSON to stdout
- `log-file`: File path to save logs
- `quiet`: Suppress all log output
- `verbose`: Verbose logging level (0-3, where 3 is most verbose)
- `frames`: Number of frames to process for videos
- `preview-format`: Format generated thumbnail for videos

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

## Copying Generated Files

The `[links]` section tells `rong` where to hardlink/copy each generated theme file.
You can assign a single path or an array of paths if you want the same file
copied/linked to multiple locations.

### Syntax

```toml{3}
[links]
"template_name.ext" = "destination_path"
"template_name.ext" = ["path1", "path2", ...]
```

::: info IMPORTANT
Each key must match the name of a template in the theme templates' output directory.
Usually, `~/.local/state/rong`.
:::

For some themes, you might need to use `[installs]`. The structure exactly same as
`[links]`. The difference is that fill will be _installed_ by atomic copy. This
ensure apps don't see incomplete theme files.

```toml
[installs]
"quickshell.json" = "~/.local/state/quickshell/colors.json"
```

### Notes

- Paths support `~` for your home directory.
- Existing files at the destination will be overwritten or replaced by symlinks.

## Post Commands

The `post-cmds` section tells `rong` to run some command in POSIX shell. For example,
`pidof kitty | xargs -r kill -SIGUSR1` to reload kitty terminal.

```toml
[post-cmds]
"hyprland.conf" = "hyprctl reload"
"kitty-full.conf" = "pidof kitty | xargs -r kill -SIGUSR1"
"spicetify-sleek.ini" = """
# If spotify is already running in debug mode
timeout 2s spicetify watch -s >/dev/null 2>&1 || true
"""
```

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

[installs]
"quickshell.json" = "~/.local/state/quickshell/colors.json"

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

post-cmds:
  cava.ini: pidof cava | xargs -r kill -SIGUSR2
  kitty-full.conf: pidof kitty | xargs -r kill -SIGUSR1
  colors.tmux: |
    tmux source-file ~/.config/tmux/tmux.conf || true
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
