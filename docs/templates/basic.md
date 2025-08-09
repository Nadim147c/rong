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

- `upper`, `lower` — changes string case
- `chroma` - adjust saturation (E.g. `{{ chroma .Primary 50 }}`)
- `tone` - adjust brightness (E.g. `{{ tone .Primary 75 }}`)
- `replace` — replaces parts of a string
- `quote` — wraps a string in quotes
- `json` — outputs JSON-formatted data

Examples:

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
