---
title: Configuration
---

# Configuration

Configuration can be loaded from a variety of sources, including files ( `json`,
`toml`, `yaml`, `yml`, `properties`, `props`, `prop`, `hcl`, `tfvars`, `dotenv`,
`env`, `ini`), command-line flags, and environment variables (prefixed with
`RONG_`). Rong will automatically merge values from all sources. For sake of
simplicity, I'm only going to show `toml` configuration. Checkout the
[examples](#examples) for different format.

## Locations

By default, the `rong` configuration file is located at:

```sh
"$XDG_CONFIG_HOME/rong" # Usually `~/.config/rong`
```

The file can be in any of the supported format.

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

[base16]
method = "dynamic"
```

## Configuration Fields

### Core Settings

- `dark`: Generate dark color palette (`true`) or light palette (`false`).
- `dry-run`: Generate colors without applying templates.
- `json`: Print generated colors as JSON to stdout.
- `log-file`: File path to save logs.
- `quiet`: Suppress all log output.
- `verbose`: Verbose logging level (0-3, where 3 is most verbose).
- `frames`: Number of frames to process for videos.
- `worker`: Number of thread for process caching.
- `preview-format`: Format generated thumbnail for videos.

### Material You Settings

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

- `version`: Material You specification version (`"2021"` or `"2025"`).
- `platform`: Target platform (`"phone"` or `"watch"`).
- `contrast`: Contrast adjustment (-1.0 to 1.0).
- `custom`: Custom color configuration.
  - `blend`: Ratio of color blending with `Material You` primary color.
  - `colors`: A list of color names with there values.

Here is an example:

```toml
[material]
contrast = 0.0
platform = "phone"
variant = "tonal_spot"
version = "2025"
custom.blend = 0.35

[material.custom.colors]
orange = "#FFA500"
purple = "#800080"
```

### Base16 Settings

The `[base16]` section controls Base16 color generation:

- `method`: Color generation method (`"static"` or `"dynamic"`).
- `blend`: Blend ratio toward the primary color (0.0 to 1.0).
- `colors`: Source colors for base16 color generation (all hex colors):
  - `black`, `blue`, `cyan`, `green`, `magenta`, `red`, `white`, `yellow`.

Here is an example:

```toml
[base16]
blend = 0.35
method = "static"
colors.red = "#FF0000"
```

### Themes Settings

You can use `themes` to copy/install file and run any command afterward.

The `[[themes]]` section controls copy and running command after generating colors.

- `target` (required): Name of target template (without `.tmpl` extension).
- `links`: A path or a list of path to **hardlink** or **copy** the template.
- `installs`: A path or a list of path to atomically **install** the template.
- `cmds`: A command or a list of command to run after `links` and `installs`.

```toml
[[themes]] # Make sure to double square brackets here
target = "spicetify-sleek.ini"
links = "~/.config/spicetify/Themes/Sleek/color.ini"
cmds = """
spicetify watch -s 2>&1 | sed '/Reloaded Spotify/q'
"""
```

The configuration above will link/copy the builtin `spicetify-sleek.ini` template to
`~/.config/spicetify/Themes/Sleek/color.ini` then run
`spicetify watch -s 2>&1 | sed '/Reloaded Spotify/q'`.

### Single copy/link/cmd action

You also use _independent_ `links`, `installs` or `cmds` action. This is useful when
you just need copy or install.

#### Syntax

```toml{3}
[links]
"template_name.ext" = "destination_path"
"template_name.ext" = ["path1", "path2", ...]
```

For some themes, you might need to use `[installs]`. The structure exactly same as
`[links]`. The difference is that fill will be _installed_ by atomic copy. This
ensure apps don't see incomplete theme files.

```toml
[installs]
"quickshell.json" = "~/.local/state/quickshell/colors.json"
```

`cmds` are commands to run after which will run after `links` and `installs`.

```toml
[cmds]
"hyprland.conf" = "hyprctl reload"
"kitty-full.conf" = "pidof kitty | xargs -r kill -SIGUSR1"
```

## Examples

**TOML Configuration** `~/.config/rong/config.toml`

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

[[themes]]
cmds = "pidof cava | xargs -r kill -SIGUSR2"
links = "~/.config/cava/config"
target = "cava.ini"

[[themes]]
cmds = "hyprctl reload"
links = "~/.config/hypr/colors.conf"
target = "hyprland.conf"

[[themes]]
cmds = "pidof kitty | xargs -r kill -SIGUSR1"
links = "~/.config/kitty/colors.conf"
target = "kitty-full.conf"

[installs]
"quickshell.json" = "~/.local/state/quickshell/colors.json"

[links]
"gtk.css" = [
  "~/.config/gtk-3.0/gtk.css",
  "~/.config/gtk-4.0/gtk.css"
]
```

**YAML Configuration** `~/.config/rong/config.yaml`

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

themes:
  - cmds: pidof cava | xargs -r kill -SIGUSR2
    links: ~/.config/cava/config
    target: cava.ini
  - cmds: hyprctl reload
    links: ~/.config/hypr/colors.conf
    target: hyprland.conf
  - cmds: pidof kitty | xargs -r kill -SIGUSR1
    links: ~/.config/kitty/colors.conf
    target: kitty-full.conf
```

**JSON Configuration** `~/.config/rong/config.json`

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
