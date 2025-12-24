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
continue.

## Link

Add the following line to the [configuration](/configuration#linking-generated-files):

```toml{3}
[links]
# ...
"spicetify-sleek.ini" = "~/.config/spicetify/Themes/Sleek/color.ini"
```

<!--@include: ./_regen.md-->

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

You need to run `spicetify apply` and restart Spotify to apply new generated theme.
Live theme reload is possible by starting spicetify in watch mode
(`spicetify watch -s`).

::: info

For flatpak version of spotify, `spicetify watch -s` **might** fail to (re)start
spotify. In that case, you can do the following step.

- Start spotify:

```bash
flatpak run com.spotify.Client \
  --remote-debugging-port=9222 \
  --remote-allow-origins='*'
```

- Live reload after generating colors color:

```bash
spicetify watch -s 2>&1 | sed '/Reloaded Spotify/q'
```

> `sed '/Reloaded Spotify/q'` exits whenever it match output `Reloaded Spotify`.

- Final `config.yaml`:

```yaml
links:
  spicetify-sleek.ini: ~/.config/spicetify/Themes/Sleek/color.ini
post-cmds:
  spicetify-sleek.ini: |
    spicetify watch -s 2>&1 | sed '/Reloaded Spotify/q'
```

:::
