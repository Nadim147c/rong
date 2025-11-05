# Kitty

kitty is a free and open-source GPU-accelerated terminal emulator for Linux, macOS,
and some BSD distributions. Focused on performance and features, kitty is written in
a mix of C and Python programming languages. It provides GPU support.

## Install

Install kitty from your preferred package manager or from
[here](https://sw.kovidgoyal.net/kitty/binary/).

## Link

Add the following lines to the
[configuration](/configuration#linking-generated-files):

- With base16 terminal colors

```toml{4-6}
[links]
"kitty-full.conf" = "~/.config/kitty/colors.conf"
```

- Without base16 terminal colors

```toml{4-6}
[links]
"kitty.conf" = "~/.config/kitty/colors.conf"
```

<!--@include: ./_regen.md-->

## Apply

Add the following line to your `kitty.conf`:

```bash
include colors.conf
```

## Reload

You can live reload kitty theme via sending signal:

```bash
pidof kitty | xargs kill -SIGUSR1
```
