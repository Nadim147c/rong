# Eww (ElKawor Wecky Widget)

Eww is a highly customizable widget system for Wayland/X11. It supports CSS or SCSS
for theming.

## Link

Add the following line to the [configuration](/configuration#linking-generated-files):

```toml{5}
[links]
# ...
"colors.scss" = [
  # ...
  "~/.config/eww/colors.scss"
]
```

## Apply

Create the Eww style file at `~/.config/eww/eww.scss` and import your color
variables:

```scss{1,8,9}
@import "colors.scss";

* {
  all: unset;
}

.widget {
  background-color: $background;
  color: $onBackground;
}
```

## Reload

Eww auto reloads styles on change. But if you're generating for the first time you
need to run:

```bash
eww reload
```
