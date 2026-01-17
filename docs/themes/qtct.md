---
title: QT Applications
---

# QT Applications (With `qt6ct-kde`)

QT is one of the most popular toolkits for creating Unix GUI applications. We can use
`qt6ct-kde` to theme KDE applications like Dolphin.

::: info
[qt6ct-kde](https://aur.archlinux.org/packages/qt6ct-kde) is a patch for
[qt6ct](https://www.opencode.net/trialuser/qt6ct) for KDE Applications.
:::

## Install `qt6ct-kde`

You can install `qt6ct-kde` from AUR (Arch User Repository).

```bash
yay -S qt6ct-kde
```

## Link

Add the following line to the
[configuration](/configuration#linking-generated-files):

```toml{2}
[links]
"qtct.colors" = "~/.config/qt6ct/colors/rong.colors"
```

<!--@include: ./_regen.md-->

## Settings

Now, open `qt6ct-kde` by running `qt6ct` or using a launcher to search for "QT6 settings".
Set `Color scheme:` to `Rong (KColorScheme)`.

## Environment Variable

Now, run any QT application with `QT_QPA_PLATFORMTHEME=qt6ct`. For example:

```bash
QT_QPA_PLATFORMTHEME=qt6ct dolphin
```

For convenience, you can set `QT_QPA_PLATFORMTHEME=qt6ct` in your global Environment
Variables. For Hyprland:

```ini
env = QT_QPA_PLATFORMTHEME,qt6ct
```

## Reload

The `qt6ct` settings auto-reload on generation, but you have to restart individual
applications for them to apply the new theme.
