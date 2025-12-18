# Integration

Generating colors isn't enough—you need to **apply** the theme as well. This often
involves running multiple commands and optionally changing your wallpaper during the
process. The best way to manage this is by using a shell script.

## Wallpaper Selection

If you want to use your wallpaper for color generation, you'll first need to get the
path to the image. For this, we'll use the following command to search for wallpapers
in `~/Pictures/Wallpapers/` (you can change this path later):

```bash
find ~/Pictures/Wallpapers/ -type f \( \
    -iname "*.jpg" -o \
    -iname "*.jpeg" -o \
    -iname "*.png" -o \
    -iname "*.webp" \)
```

Alternatively, you can use [fd](https://github.com/sharkdp/fd); `fd
'\.(jpg|jpeg|png|webp)' ~/Pictures/Wallpapers --type f`.

This finds all `jpg`, `jpeg`, `png`, and `webp` files in the specified directory.

To randomly select one:

```bash
WALLPAPER=$1
if [ -z "$WALLPAPER" ]; then
  WALLPAPER=$(find ~/Pictures/Wallpapers/ -type f \( \
      -iname "*.jpg" -o \
      -iname "*.jpeg" -o \
      -iname "*.png" -o \
      -iname "*.webp" \) | shuf -n1)
fi
```

## Wallpaper Change

Once you have your wallpaper file, you can set it using your preferred wallpaper
daemon. Here's an example using `swww`:

```bash
# Apply wallpaper with swww
swww img \
  --transition-duration 2 \
  --transition-bezier ".09,.91,.52,.93" \
  --transition-fps 60 \
  --invert-y \
  "$WALLPAPER" &
```

## Generate Colors

Next, use the wallpaper to generate a new color scheme:

```bash
rong image "$WALLPAPER"
```

::: info IMPORTANT
For video (animated) wallpapers, use `rong video "$WALLPAPER"`.
:::

This step generates the theme and applies templates to their configured destinations.

## Post-Change Commands

Some applications need to be reloaded **after** their config files are updated.
Instead of handling this in your shell script, you can now define commands that run
automatically after a template is rendered and copied.

This is done via a `[post-cmds]` section in your configuration:

```toml
[post-cmds]
"cava.ini"        = "pidof cava | xargs -r kill -SIGUSR2"
"hyprland.conf"   = "hyprctl reload"
"kitty-full.conf" = """
  PID=$(pidof kitty)
  if [[ -n "$PID" ]]; then
  | kill -SIGUSR1 $PID
  fi
"""
"some-other-template" = [
  "first command",
  "second command",
]
```

Each key is the name of a generated template file, and the value is the command that
should be executed once that file is updated. This keeps reload logic close to the
config it affects and removes the need for custom post-hook scripts.

## Final Script

With everything wired up, the shell script becomes very small and focused.
For example, `~/.local/bin/wallpaper.sh`:

```bash
#!/usr/bin/env bash

WALLPAPER=$1
WALLPAPER_DIR="$HOME/Pictures/Wallpapers/"

if [ -z "$WALLPAPER" ]; then
  WALLPAPER=$(find "$WALLPAPER_DIR" -type f \( \
    -iname "*.jpg" -o \
    -iname "*.jpeg" -o \
    -iname "*.png" -o \
    -iname "*.webp" \
    \) | shuf -n1)
fi

if [ -z "$WALLPAPER" ]; then
  echo "ERROR: No wallpaper found"
  exit 1
fi

# Apply wallpaper
(
  exec setid swww img \
    --transition-duration 2 \
    --transition-bezier ".09,.91,.52,.93" \
    --transition-fps 60 \
    --invert-y \
    "$WALLPAPER" &
  disown
) & # Prevent swww from stopping if script exists too early

# Generate colors and apply templates
rong image "$WALLPAPER"
```

Make it executable:

```bash
chmod +x ~/.local/bin/wallpaper.sh
```

## Usage

Run the script with or without an argument:

```bash
wallpaper.sh                    # random wallpaper
wallpaper.sh path/to/image.png  # specific wallpaper
```

The wallpaper is applied, colors are generated, templates are updated, and any
configured post-commands are executed automatically.
