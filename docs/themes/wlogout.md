# Wlogout

Wlogout is a customizable logout menu for Wayland environments. It uses `GTK-CSS` for
styling.

<!--@include: ./_gtk-issue.md-->

## Link

Add the following line to the [configuration](/configuration#linking-generated-files):

```toml{5}
[links]
# ...
"gtk-css.css" = [
  # ...
  "~/.config/wlogout/style.css"
]
```

<!--@include: ./_regen.md-->

## Apply

Create the Wlogout style file at `~/.config/wlogout/style.css` with the following
content:

```css{2}
button {
  background-color: @background;
}
```

## Reload

Wlogout isn't a long running application. Thus, Wlogout automatically loads newly
generated theme.
