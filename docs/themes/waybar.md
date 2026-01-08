# Waybar

Waybar is a status bar for Wayland compositors. It uses `GTK-CSS` for styling.

<!--@include: ./_gtk-issue.md-->

## Link

Add the following line to the [configuration](/configuration#linking-generated-files):

```toml{5}
[[themes]]
target = "gtk-css.css"
links = "~/.config/waybar/colors.css"
# cmds = "killall -SIGUSR2 waybar"
```

<!--@include: ./_regen.md-->

## Apply

Create a Waybar style file at `~/.config/waybar/style.css` with the following content:

```css{1,5,6}
@import "colors.scss";

.modules {
  /* ... */
  color: @on_background;
  background-color: @background;
}
```

## Reload

To reload styles automatically, enable `reload_style_on_change` in `config.jsonc`:

```jsonc{2}
{
  "reload_style_on_change": true,
  // ...
}
```

::: warning
This option may not always work reliably. As a fallback, run the following command:

```bash
killall -SIGUSR2 waybar
```

:::
