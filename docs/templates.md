---
title: Theme Templates
---

# Theme Templates

Rong uses Go’s [text/template](https://pkg.go.dev/text/template) to render themes
from Material colors. Go's templates provide simple yet powerful syntax to generate
templates.

## Built-In

Rong has a list of built-in templates for common formats. Checks output of built-in
templates:

```bash
rong color teal
ls ~/.local/state/rong/
```

For example, there is a built-in template for **bash**. You can use that in your
shell script:

```bash
source ~/.local/state/rong/colors.bash
echo "$PRIMARY"
```

## Custom Templates

You can create your own templates use [Golang's templates syntax](./templates/basic).
Templates are rendered using structured color data; you can find the available
variables in the [template context](./templates/context).

---

Once you're familiar with the template syntax and execution data, follow these
steps:

1. **Write the template**

   Create your file in `~/.config/rong/templates/` using the `.tmpl` extension.

   Example:

   ```bash
   mkdir -p ~/.config/rong/templates
   echo '
   0: {{ .Primary }}
   1: {{ .Secondary }}
   2: {{ .Tertiary }}
   ' > ~/.config/rong/templates/my_theme.ext.tmpl
   ```

2. **Generate the theme**

   Generate themes from a color, image or video:

   ```bash
   rong color cyan
   ```

   Inspect the output:

   ```bash
   cat ~/.local/state/rong/my_theme.ext
   ```

   ```text
   0: #D19488
   1: #E7BDB5
   2: #FFF8F3
   ```

3. **Link it**

   Add an entry in `config` to copy/link/install the theme to your desired location.
   See [configuration](./configuration#Links).
