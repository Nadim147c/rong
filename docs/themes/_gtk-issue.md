---
pagefind-indexed: false
---

::: details Common Issue with GTK-CSS {close}
GTK's CSS implementation supports many features of standard (vanilla) CSS used in web
development; however, some features are not available or behave differently. The
level of support can also vary depending on the specific application or widget.

One commonly used feature in GTK CSS is color variable via the `@<name>` syntax
(e.g., `color: @mycolor;`). This approach, while convenient, is not part of
the official CSS specification and may trigger syntax errors or warnings in many text
editors.

There are two common ways to handle this:

1. **Ignore editor warnings**: You can safely ignore these errors.
2. **Use SCSS preprocessing**: Define your variables in SCSS and compile it into
   standard CSS with hardcoded values. This avoids editor errors while maintaining
   maintainability in your source stylesheets. [Learn more](/tricks/scss-for-gtk).

:::
