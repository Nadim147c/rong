# SwayNC

SwayNC is a notification daemon for Wayland compositors. It supports `GTK-CSS` for
theming, allowing you to style notifications css.

<!--@include: ./_gtk-issue.md-->

## Link

Add the following line to the [configuration](/configuration#linking-generated-files):

```toml{5}
[[themes]]
target = "gtk-css.css"
links = "~/.config/swaync/style.css"
cmds = "swaync-client --reload-css"
```

<!--@include: ./_regen.md-->

## Apply

Create the SwayNC style file at `~/.config/swaync/style.css` with the following content:

```css{2,3}
.notification {
  color: @on_background;
  background-color: @background;
}
```

## Reload

Reloading via `--reload-css` may not work sometimes. Restarting the service manually might be required.

```bash
pkill swaync
setid swaync & disown
```
