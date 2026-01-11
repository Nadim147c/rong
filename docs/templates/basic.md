---
title: Go Templates Basics
---

<span v-pre>

# Basics

Go templates use`{{ ... }}` to embed expressions, variables, or functions.

## Data Context (`.`)

The dot `.` refers to the current context. In `rong`, this includes:

- `.Colors` — a list of all material colors
- `.Image` — the source image or video path
- `.<ColorName>` — individual color values (e.g., `.Background`)

```go
{{ .Background }}  {{/* Prints a hex color like "RRGGBB" */}}
```

## Comments

Use `{{/* ... */}}` for comments.

```go
{{/* This is a comment */}}
```

## Loops

Use `range` to iterate over a slice or map.

```go
{{ range .Colors }}
  {{ .Name.Snake }} = "{{ .Color }}"
{{ end }}
```

You can also use an index variable:

```go
{{ range $i, $color := .Colors }}
  {{ $i }}: {{ $color.Name.Camel }}
{{ end }}
```

## Conditionals

Use `if`, `else`, and `end` to control logic.

```go
{{ if .Background }}
  Background: {{ .Background }}
{{ else }}
  No background color provided.
{{ end }}
```

You can nest conditions:

```go
{{ if and .Background .Image }}
  Background and image are both available.
{{ end }}
```

## Functions

Standard functions like `index`, `print`, and `printf` are supported.
`rong` also adds:

- **`upper`** – Converts a string to uppercase.
- **`lower`** – Converts a string to lowercase.
- **`replace`** – Replaces all occurrences of a substring with another substring.
- **`parse`** – Parses a color.
- **`chroma`** – Adjusts the chroma (color intensity) of a color.
- **`tone`** – Adjusts the tone (lightness/darkness) of a color.
- **`blend`** – Blend to colors with given ratio.
- **`quote`** – Wraps the given value in double quotes as a string.
- **`json`** – Converts a value to its JSON string representation or `"null"` on error.
- **`sprig`** - A big list of templates functions: https://masterminds.github.io/sprig.

Examples:

```go
{{ (blend "#FF0000" "#00FF00" 0.5).RGB }} // blend red and grean and put it RGB
{{ blend .Primary "#00FF00" 0.7 }} // blend primary toward grean
```

```go
{{ tone .Primary 70 }}
```

```go
{{ chroma .Primary 70 }}
```

```go
{ "image": {{ .Image | json }} }
```

```go
{{ replace .Title " " "_" | upper }}
```

```go
{{ quote .Background }}
```

```go
{{ printf "%s: %s" .Name .Value }}
```

## Pipelining

You can chain multiple functions using the pipe `|`.

```go
{{ .Primary | lower | quote }}
```

</span>
