---
title: Getting Started
---

# Getting Started

To get started [install](#installation) `rong` and [generate](#generate-colors)
colors.

## Installation

#### Arch Linux (AUR)

```bash
yay -S rong
```

#### Nix (flake)

```bash
nix run github:Nadim147c/rong -- --help
```

For better nix integration check [nix page](./nix).

#### Go Install

> Needs: `go`.

```bash
go install github.com/Nadim147c/rong/v5@latest
```

#### Manual from source

> Needs: `go`, `coreutils` and `just`.

```bash
# cd my-build-dir
git clone https://github.com/Nadim147c/rong.git
cd rong
just build-install
```

---

To ensure you've properly installed `rong`, run:

```bash
rong --help
```

If you see the help menu, then you've successfully installed `rong`.

```
A Material You color generator from an image or video.

Usage:
  rong [command]

Available Commands:
  ...
```

## Generate Colors

- To extract Material You compatible colors from an image:

  ```bash
  rong image /path/to/image
  ```

- To extract colors from a video:
  ```bash
  rong video /path/to/video
  ```

::: tip

If you want to use both video and image, you can use the `video` command.

```bash
rong video /path/to/image/or/video
```

:::

Generated colors will be used to generate theme files using templates. These
generated files will be stored in
[`<user-state-dir>/rong`](https://specifications.freedesktop.org/basedir-spec/latest/#variables)
(usually `~/.local/state/rong`):

```bash
$ ls ~/.local/state/rong/
colors.scss  gtk-css.css  midnight-discord.css  spicetify-sleek.ini
colors.css   dunstrc      hyprland.conf         pywalfox.json
colors.lua   ghostty      image.txt             qtct.conf
colors.nu    gtk.css      kitty.conf            rofi.rasi
```

**Rong** has a list of built-in templates for commonly used formats. You can also
create your own theme templates. See the [templates page](./templates.md).
