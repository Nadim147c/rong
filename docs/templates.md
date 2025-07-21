---
title: Theme Templates
---

# Theme Templates

Rong uses Go’s [text/template](https://pkg.go.dev/text/template) to render themes
from Material colors. Go's templates provide simple yet powerful syntax to generate
templates.

## Built-in

Rong has big enough list of built-in templates for various. After
[generating](./getting-started#Generate-Colors) colors, `rong` will execute these
templates and put them in
[`$XDG_STATE_DIR/rong/`](https://specifications.freedesktop.org/basedir-spec/latest/#variables)
(usually `~/.local/state/rong/`). Run:

```bash
ls ~/.local/state/rong/
```

This command will show list of generate theme file ready to be used in your desired
applications. To automatically copy/link these theme files to any desired location,
check out [links](./configuration#linking-generated-files) in configuration files

## Custom Templates

As built-in templates might not be enough for all kinds of use case, you can write
you own template. To create template you have to use Go’s `text/template` syntax. You
can learn about go templates from [offical docs](https://pkg.go.dev/text/template) or
[here](./templates/basic). Templates are executed with a data called execution
context. Learn about rong's templates execution context from
[here](./templates/context). After learn about the template follow these steps:

1. **Create a file**

   Save it in the `rong` templates directory:
   `$XDG_CONFIG_HOME/rong/templates/` (usually `~/.config/rong/templates/`)
   File extension must be `.tmpl`
   Example: `~/.config/rong/templates/mytemplate.tmpl`

2. **Write the template**

   Use Go template syntax. Save the file.

3. **Generate the theme**

   Run `rong <image|video>`
   Output will be in `$XDG_STATE_HOME/rong/` (usually `~/.local/state/rong/`)
   The output file will match your template filename, e.g., `mytemplate`

4. **Link it**

   Add an entry in `config.toml` to link the theme to your desired location. More
   info: [Configuration](./configuration#Links).
