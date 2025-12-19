# Integration

The most common use case for `rong` is generating colors from wallpaper. Here we will
to go through creating a simple script. This script will randomly select a image, set
that image as wallpaper and generate colors using `rong`.

## The script

1.  Let's start with listing all image out `~/Pictures/Wallpapers/` unix `find`
    command:

    ```bash
    find ~/Pictures/Wallpapers/ -type f \( \
        -iname "*.jpg" -o \
        -iname "*.jpeg" -o \
        -iname "*.png" -o \
        -iname "*.webp" \)
    ```

    Alternatively, you can use [fd](https://github.com/sharkdp/fd):

    ```bash
    fd -tf -e .jpg -e .jpeg -e .png -e .webp . ~/Pictures/Wallpapers
    ```

    This finds all `jpg`, `jpeg`, `png`, and `webp` files in the specified directory.

2.  We can select a random image by `shuf -n1`. Sometimes, we might want to set a
    specific wallpaper instead of random one. Here this script checks if any
    wallpaper path is specified. If not than it chooses a random wallpaper.

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

3.  Once you have your wallpaper file, you can set it using your preferred wallpaper
    daemon. Here's an example using `swww`:

    ```bash
    # Apply wallpaper with swww
    swww img \
      --transition-duration 2 \
      --transition-bezier ".09,.91,.52,.93" \
      --transition-fps 60 \
      "$WALLPAPER"
    ```

4.  Finally, use the wallpaper to generate a new color scheme:

    ```bash
    rong image "$WALLPAPER"
    ```

    ::: info IMPORTANT
    For video (animated) wallpapers, use `rong video "$WALLPAPER"`.
    :::

    This step generates the theme and applies templates to their configured destinations.

5.  With everything wired up, we can finish the script. Now, we can save the script
    to any directory which is listed in out `PATH` path variable. For example, we
    will save this script in `~/.local/bin/chwall`.

    ::: warning IMPORTANT
    Make sure is `~/.local/bin` in your `PATH` variable:

    ```bash
    printenv PATH | grep -q ~/.local/bin && echo "YES" || echo "NO"
    ```

    :::

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
    setsid swww img \
      --transition-duration 2 \
      --transition-bezier ".09,.91,.52,.93" \
      --transition-fps 60 \
      --invert-y \
      "$WALLPAPER" &
    disown # Prevent swww from stopping if script exists too early

    # Generate colors and apply templates
    rong image "$WALLPAPER"
    ```

    Make it executable:

    ```bash
    chmod +x ~/.local/bin/chwall
    ```

    ### Usage

    Run the script with or without an argument:

    ```bash
    chwall                    # random wallpaper
    chwall path/to/image.png  # specific wallpaper
    ```

## Post-Change Commands

Some applications need to be reloaded **after** their [config](./configuration.md) files are updated.
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
