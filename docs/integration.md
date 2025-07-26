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

## Post-Change Hook

Some applications may require a manual reload after the theme has changed. You can
handle this with a `post_hooks` function, defined at the top of your script for easy
access and modification.

```bash
post_hooks() {
    # Compile SCSS files for Waybar (see GTK/SCSS theming docs)
    compile-scss ~/.config/waybar/style.scss && killall -v -SIGUSR2 waybar

    # Reload dunst without resetting pause level
    local dunst_level=$(dunstctl get-pause-level)
    dunstctl reload && dunstctl set-pause-level "$dunst_level"

    # Update Pywalfox
    pywalfox --verbose update

    # Reload Hyprland config
    hyprctl reload
}
```

## Final Script

Now put everything together into a single executable script—e.g.,
`~/.local/bin/wallpaper.sh`:

```bash
!/bin/bash

WALLPAPER=$1
WALLPAPER_DIR="$HOME/Pictures/Wallpapers/"

# Post-configuration hook
post_hooks() {
    compile-scss ~/.config/waybar/style.scss && killall -v -SIGUSR2 waybar

    local dunst_level=$(dunstctl get-pause-level)
    dunstctl reload && dunstctl set-pause-level "$dunst_level"

    pywalfox --verbose update
    hyprctl reload
}

if [ -z "$WALLPAPER" ]; then
  WALLPAPER=$(find "$WALLPAPER_DIR" -type f \( \
      -iname "*.jpg" -o \
      -iname "*.jpeg" -o \
      -iname "*.png" -o \
      -iname "*.webp" \) | shuf -n1)
fi

if [ -z "$WALLPAPER" ]; then
  echo "ERROR: No wallpaper found"
  exit 1
fi

# Apply wallpaper using swww
swww img \
  --transition-duration 2 \
  --transition-bezier ".09,.91,.52,.93" \
  --transition-fps 60 \
  --invert-y \
  "$WALLPAPER" &

# Generate color scheme
rong image "$WALLPAPER"

post_hooks
```

Make it executable:

```bash
chmod +x ~/.local/bin/wallpaper.sh
```

You're good to go. Run the script with or without an argument to set the theme and
wallpaper:

```bash
wallpaper.sh                    # random wallpaper
wallpaper.sh path/to/image.png  # specific wallpaper
```
