# Spotify

Discord doesn't support theming by default. You need to modify the Spotify client to
apply a material theme. This guide is for [spicetify](https://spicetify.app/) mod.

::: danger STOP
Using a modified client violates Spotify’s
[Terms of Service](https://www.spotify.com/legal/) and may result in account
suspension or a permanent ban. Proceed at your own risk — I am not responsible for
any actions taken against your account. However, I yet to see anyone getting banned
only for theming.
:::

## Install

Follow the installation
[guide](https://spicetify.app/docs/getting-started#linux-and-macos) of spicetify and
continue

## Link

Add the following line to the [configuration](/configuration#linking-generated-files):

```toml{3}
[links]
# ...
"spicetify-sleek.ini" = "~/.config/spicetify/Themes/Sleek/color.ini"
```

## Apply

Download `Sleek` theme from spicetify GitHub repository:

```bash
mkdir -p ~/.config/spicetify/Themes/Sleek/
curl -L "https://github.com/spicetify/spicetify-themes/raw/refs/heads/master/Sleek/user.css" \
  -o ~/.config/spicetify/Themes/Sleek/user.css
```

Set the new theme and apply changes:

```bash
spicetify config current_theme Sleek
spicetify config color_scheme rong
spicetify apply
```

## Reload

Reloading Spotify is somewhat complicated. You need to run `spicetify apply` and
restart Spotify to apply new generated theme. However, it is possible by running
spicetify in watch mode (`spicetify watch -s`).

### Reload with Hyprland

To auto reload Spotify on Hyprland compositor, make these changes to your Hyprland
configuration

```ini
$music = spotify # [!code --]
$music = spicetify watch -s # [!code ++]

# To start Spotify on 4th workspace with watch mode on start
exec-once = [workspace 4 silent] $music

# To start Spotify on with $mainMod (usually SUPER) + M shortcut
bind = $mainMod, M, exec, $music
```
