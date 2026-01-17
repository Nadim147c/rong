---
title: GTK Applications
---

# GTK Applications

A large numbers of for Linux uses applications GTK (GIMP Toolkit) for UI. You can
apply theme by editing `gtk.css`. Apply this method will make affect most GTK based
applications with some exception.

## Link

Add the following line to the [configuration](/configuration#linking-generated-files):

```toml{3-4}
[links]
"gtk.css" = [
  "~/.config/gtk-3.0/gtk.css",
  "~/.config/gtk-4.0/gtk.css",
]
```

<!--@include: ./_regen.md-->

## Reload

Restarting your GTK based applications will apply the material theme. If you're
changing the theme color then run:

> Replace `[mode]` with `light` or `dark`.

```bash
gsettings set org.gnome.desktop.interface gtk-theme ""
gsettings set org.gnome.desktop.interface gtk-theme adw-gtk3-[mode]
```
