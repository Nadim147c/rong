---
title: GTK Applications
---

# GTK Applications

A large numbers of for linux uses applications GTK (GIMP Toolkit) for UI. You can
apply theme by editing `gtk.css`. Apply this method will make effect most gtk based
applications with some exception.

::: details Here some of the applications that will effected by this: {close}

| Name           | Toolkit | Name                   | Toolkit |
| -------------- | ------- | ---------------------- | ------- |
| Thunar         | GTK 3   | Gedit                  | GTK 3   |
| GNOME Terminal | GTK 3   | GNOME Files (Nautilus) | GTK 4   |
| GNOME Calendar | GTK 4   | GNOME Maps             | GTK 4   |
| GNOME Weather  | GTK 4   | GNOME Calculator       | GTK 4   |
| GNOME Music    | GTK 4   | GNOME Photos           | GTK 4   |
| GNOME Contacts | GTK 4   | Evince                 | GTK 3   |
| Eye of GNOME   | GTK 3   | Geary                  | GTK 4   |
| Fractal        | GTK 4   | Builder                | GTK 4   |
| Loupe          | GTK 4   | Snapshot               | GTK 4   |
| Celluloid      | GTK 4   | Foliate                | GTK 3   |
| Lollypop       | GTK 3   | Blueman                | GTK 3   |
| NM-Applet      | GTK 3   | Pavucontrol            | GTK 3   |
| Mousepad       | GTK 3   | Atril                  | GTK 3   |
| Pluma          | GTK 3   | Engrampa               | GTK 3   |
| Xed            | GTK 3   | Xreader                | GTK 3   |
| Galculator     | GTK 3   | LXTerminal             | GTK 3   |

:::

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

Restarting you're gtk based applications will apply the material theme. If you're
changing the theme color then run:

> Replace `[mode]` with `light` or `dark`.

```bash
gsettings set org.gnome.desktop.interface gtk-theme ""
gsettings set org.gnome.desktop.interface gtk-theme adw-gtk3-[mode]
```
