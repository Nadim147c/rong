---
title: SCSS For GTK
---

# SCSS For GTK

A concise guide to using SCSS for GTK styling, including compiler installation,
variable usage, and compiling SCSS to GTK-compatible CSS. Includes a shell script for
easy compilation and watch mode.

## Install SCSS Compiler

On Arch-based systems, install the SCSS compiler using:

```bash
sudo pacman -S ruby-sass
```

## Link

Add the following line to the [configuration](/configuration#linking-generated-files):

```toml{5}
[links]
# ...
"colors.scss" = [
  # ...
  "~/.config/<app-name>/colors.scss"
]
```

## Writing SCSS Using Variables

Create a `style.scss` file in your configuration directory. For waybar, it should be
in `~/.config/waybar/style.scss`.

::: warning STOP
Create a backup for the original `style.css` file.

```bash
mv ~/.config/waybar/style.css ~/.config/waybar/style.css.bak
```

:::

```scss
@import "colors";

.module {
  color: $onBackground;
  background-color: $background;
}
```

## Create SCSS Compilation Script

Save the following script as `compile-scss.sh`:

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
    echo "Error: 'sass' command not found. Please install Ruby Sass."
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
    local TMPFILE="$(mktemp "$DIR/$BASE.XXXXXX.css")"

    # Compile SCSS to a temporary file
    if sass --no-cache --sourcemap=none "$INPUT" "$TMPFILE"; then
        # Use install to safely move the file (preserves permissions, atomic)
        install -m 644 "$TMPFILE" "$OUTPUT"
        echo "Compiled $INPUT -> $OUTPUT"
        rm -f "$TMPFILE"
    else
        echo "SCSS compilation failed."
        rm -f "$TMPFILE"
        return 1
    fi
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

## Make the Script Executable and Accessible

```bash
chmod +x compile-scss.sh
mkdir -p ~/.local/bin
mv compile-scss.sh ~/.local/bin/compile-scss
```

Make sure `~/.local/bin` is in your `PATH`. Add this to your `.bashrc`, `.zshrc`,
or `.profile` if needed:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

## Compile Your SCSS

Run once:

```bash
compile-scss path/to/style.scss
# It compiles path/to/style.scss -> path/to/style.css
```

For waybar:

```bash
compile-scss ~/.config/waybar/style.scss
# It compiles ~/.config/waybar/style.scss -> ~/.config/waybar/style.css
```

::: tip
You can use `--watch` flag to watch for change and automatically compile.

```bash
compile-scss --watch path/to/style.scss
```

:::
