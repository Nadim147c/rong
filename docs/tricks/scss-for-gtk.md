---
title: SCSS For GTK
---

# SCSS for GTK

A concise guide to using SCSS for GTK styling: install a compiler, use variables, and
compile SCSS into GTK-compatible CSS. Includes a shell script for automatic
compilation and watch mode.

## Why SCSS?

SCSS supports `@use` and variables, which are useful because—after compilation—the
output CSS becomes a single, flat file with all variables resolved. While GTK itself
supports CSS imports via `@import`, many applications (such as Wofi) do not handle
relative paths correctly.

## Install the SCSS Compiler

On Arch-based systems, install the SCSS compiler with:

```bash
sudo pacman -S dart-sass
```

## Link

Add the following line to the
[configuration](/configuration#linking-generated-files):

```toml
[links]
"colors.scss" = [
  "~/.config/<app-name>/colors.scss"
]
```

## Writing SCSS Using Variables

Create `~/.config/waybar/style.scss`:

```scss
@use "colors.scss" as *;

.module {
  color: $onBackground;
  background-color: $background;
}
```

::: warning

Backup the original `style.css` before replacing it:

```bash
mv ~/.config/waybar/style.css ~/.config/waybar/style.css.bak
```

:::

## SCSS Compilation Script

Save this as `compile-scss.sh`:

```sh
#!/bin/sh

# Usage: compile-scss.sh [--watch] path/to/file.scss

WATCH_MODE=false

if [ "$1" = "--watch" ]; then
    WATCH_MODE=true
    shift
fi

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 [--watch] path/to/file.scss"
    exit 1
fi

INPUT="$1"

if [ ! -f "$INPUT" ]; then
    echo "Error: File '$INPUT' does not exist."
    exit 1
fi

if ! command -v sass >/dev/null 2>&1; then
    echo "Error: 'sass' command not found. Please install Dart Sass."
    exit 1
fi

if $WATCH_MODE && ! command -v inotifywait >/dev/null 2>&1; then
    echo "Error: 'inotifywait' command not found. Please install inotify-tools."
    exit 1
fi

compile_scss() {
    local INPUT="$1"
    local DIR="$(dirname "$INPUT")"
    local BASE="$(basename "$INPUT" .scss)"
    local OUTPUT="$DIR/$BASE.css"
    sass --no-source-map --verbose "$INPUT":"$OUTPUT"
}

# Initial compilation
compile_scss "$INPUT"

if $WATCH_MODE; then
    echo "Watching for changes in $INPUT..."
    while true; do
        inotifywait -e close_write -e modify -e move -e create -e delete --exclude '\.css$' "$INPUT" "$(dirname "$INPUT")"
        compile_scss "$INPUT"
    done
fi
```

## Install Script and Update PATH

Make the script executable and move it to a directory in your `PATH`:

```bash
chmod +x compile-scss.sh
mkdir -p ~/.local/bin
mv compile-scss.sh ~/.local/bin/compile-scss
```

If `~/.local/bin` isn't already in your `PATH`, add it:

```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc  # or .zshrc / .profile
```

## Compile Your SCSS

To compile once:

```bash
compile-scss ~/.config/waybar/style.scss
# Output: ~/.config/waybar/style.css
```

To watch for changes and auto-compile:

```bash
compile-scss --watch ~/.config/waybar/style.scss
```
