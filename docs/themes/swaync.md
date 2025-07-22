# SwayNC

SwayNC is a notification daemon for Wayland compositors. It supports `GTK-CSS` for
theming, allowing you to style notifications css.

<!--@include: ./_gtk-issue.md-->

## Link

Add the following line to the [configuration](/configuration#linking-generated-files):

```toml{5}
[links]
# ...
"gtk-css.css" = [
  # ...
  "~/.config/swaync/style.css"
]
```

## Apply

Create the SwayNC style file at `~/.config/swaync/style.css` with the following content:

```css{2,3}
.notification {
  color: @on_background;
  background-color: @background;
}
```

## Reload

To apply style changes without restarting the daemon, run:

```bash
swaync-client --reload-css
```

::: warning
Reloading via `--reload-css` may not work sometimes. Restarting the service manually might be required.

```bash
pkill swaync
swaync &
```

:::
