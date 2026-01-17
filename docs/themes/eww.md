# Eww (ElKawor Wacky Widget)

Eww is a highly customizable widget system for Wayland/X11. It supports CSS or SCSS
for theming. SCSS is recommended since it has more features and flexibility.

## Configuration

Add the following line to the [configuration](/configuration#linking-generated-files):

```toml
[[themes]]
target = "colors.css"
links = "~/.config/eww/colors.scss"
cmds = "eww reload"
```

<!--@include: ./_regen.md-->

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
