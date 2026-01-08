---
title: Getting Started
---

# Getting Started

To get started [install](#installation) `rong` and [generate](#generate-colors)
colors.

## Installation

You need `go` to install `rong`. Run:

```bash
go install github.com/Nadim147c/rong/v5@latest
```

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

To extract Material You-compatible colors from an image:

```bash
rong image /path/to/image
```

To extract colors from a video:

```bash
rong video /path/to/video
```

::: info

This command internally uses `ffmpeg` to extract frames. Only **5 frames** are
sampled by default to ensure a balance between performance and accuracy. These frames
are **evenly distributed** across the video duration—not just the first 5.

:::

The generated colors will be used to create theme files using built-in templates (or
user-defined templates). These files will be stored in
[`$XDG_STATE_DIR/rong/`](https://specifications.freedesktop.org/basedir-spec/latest/#variables)
(usually `~/.local/state/rong/`):

```bash
$ ls ~/.local/state/rong/
colors.scss  gtk-css.css  midnight-discord.css  spicetify-sleek.ini
colors.css   dunstrc      hyprland.conf         pywalfox.json
colors.lua   ghostty      image.txt             qtct.conf
colors.nu    gtk.css      kitty.conf            rofi.rasi
```

See the [templates page](./templates.md).

::: tip

If you're not sure whether the file is an image or video—or if you want to use
both—you can use the `video` command, as `ffmpeg` supports both image and video
inputs:

```bash
rong video /path/to/image/or/video
```

:::
