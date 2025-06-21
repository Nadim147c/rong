# Rong (রং)

> [!CAUTION]
> This is highly experimental 🧪 and expect breaking changes

**Rong** is a CLI that extracts a **Material You** color palette from an image or
video, and applies it across your system using configurable template mappings.

---

### 🔧 Configuration (`config.toml`)

Rong uses a `[link]` table to map template outputs to destination files, allowing you
to theme multiple apps consistently from a single palette.

```toml
[link]
"colors.lua" = "~/.config/wezterm/colors.lua"
"colors.scss" = "~/.config/eww/colors.scss"
"spicetify-sleek.ini" = "~/.config/spicetify/Themes/Sleek/color.ini"
"kitty.conf" = "~/.config/kitty/colors.conf"
"pywalfox.json" = "~/.cache/wal/colors.json"
"gtk-css.css" = [
  "~/.config/waybar/colors.css",
  "~/.config/swaync/colors.css",
  "~/.config/wlogout/colors.css",
]
"midnight-discord.css" = [
  "~/.config/equibop/settings/quickCss.css",
  "~/.config/vesktop/settings/quickCss.css",
]
```

### 🔧 Templates (`~/.config/rong/templates/colors.gotmpl`)

Rong uses go's [text/template](https://pkg.go.dev/text/template) to generate theme
files.

```gotemplate
// Auto-generated Material You color scheme
$primary: {{ .Primary }};
$on-primary: {{ .OnPrimary }};
```

This will output SCSS-compatible variables like:

```scss
$primary: #6750A4;
$on-primary: #FFFFFF;
...
```

For more example checkout built-in [templates](./templates/built-in/).

### 🎨 Features

- Extracts color palettes from **images** or **videos**
- Generates color output using customizable **templates**
- Automatically writes output to multiple **target files**
- Easily theme apps like **WezTerm**, **Eww**, **Kitty**, **Spicetify**, **Waybar**,
  and **Discord mods**

**Example:**

```sh
rong image ~/Pictures/wallpaper.jpg
```

## License and Credit

Rong is licensed under [GNU GPL-3.0](./LICENSE). This cli uses [material](https://github.com/Nadim147c/material)
for generating colors.

### Thanks to

- [Matugen](https://github.com/InioX/matugen-themes/): For some of the templates.
- [Material Color Utilities](https://github.com/material-foundation/material-color-utilities): For the material color algorithm.
