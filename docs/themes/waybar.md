# Waybar

Waybar is a status bar for Wayland compositors. It uses `gtk-css` for styling,
allowing the use of `@define-color` via a built-in `gtk-css.css` file.

## Link

Add the following line to the [configuration](/configuration#linking-generated-files):

```toml{5}
[links]
# ...
"gtk-css.css" = [
  # ...
  "~/.config/waybar/colors.css"
]
```

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
