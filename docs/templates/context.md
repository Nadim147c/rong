---
title: Templates Context
---

<span v-pre>

# Template Context

The following data is passed to Go's template engine to render the themes files.

## Top-Level Fields

| Field Path  | Type       | Example Template     | Output Example   |
| ----------- | ---------- | -------------------- | ---------------- |
| `.Image`    | `string`   | `{{ .Image }}`       | `"material.png"` |
| `.Material` | `Material` | N/A (use sub-fields) |                  |
| `.Colors`   | `[]Color`  | See looping section  |                  |

## Material Color Fields

All fields are of type `ColorValue`.

| Field Name                | Example Template                 | Output Example |
| ------------------------- | -------------------------------- | -------------- |
| `Background`              | `{{ .Background }}`              | `121212`       |
| `Error`                   | `{{ .Error }}`                   | `B3261E`       |
| `ErrorContainer`          | `{{ .ErrorContainer }}`          | `F9DEDC`       |
| `InverseOnSurface`        | `{{ .InverseOnSurface }}`        | `322F29`       |
| `InversePrimary`          | `{{ .InversePrimary }}`          | `D0BCFF`       |
| `InverseSurface`          | `{{ .InverseSurface }}`          | `E4E2E6`       |
| `OnBackground`            | `{{ .OnBackground }}`            | `E0E0E0`       |
| `OnError`                 | `{{ .OnError }}`                 | `FFFFFF`       |
| `OnErrorContainer`        | `{{ .OnErrorContainer }}`        | `410E0B`       |
| `OnPrimary`               | `{{ .OnPrimary }}`               | `FFFFFF`       |
| `OnPrimaryContainer`      | `{{ .OnPrimaryContainer }}`      | `3700B3`       |
| `OnPrimaryFixed`          | `{{ .OnPrimaryFixed }}`          | `21005D`       |
| `OnPrimaryFixedVariant`   | `{{ .OnPrimaryFixedVariant }}`   | `000000`       |
| `OnSecondary`             | `{{ .OnSecondary }}`             | `FFFFFF`       |
| `OnSecondaryContainer`    | `{{ .OnSecondaryContainer }}`    | `1D192B`       |
| `OnSecondaryFixed`        | `{{ .OnSecondaryFixed }}`        | `332D41`       |
| `OnSecondaryFixedVariant` | `{{ .OnSecondaryFixedVariant }}` | `4A4458`       |
| `OnSurface`               | `{{ .OnSurface }}`               | `E0E0E0`       |
| `OnSurfaceVariant`        | `{{ .OnSurfaceVariant }}`        | `C4C7C5`       |
| `OnTertiary`              | `{{ .OnTertiary }}`              | `FFFFFF`       |
| `OnTertiaryContainer`     | `{{ .OnTertiaryContainer }}`     | `31111D`       |
| `OnTertiaryFixed`         | `{{ .OnTertiaryFixed }}`         | `492532`       |
| `OnTertiaryFixedVariant`  | `{{ .OnTertiaryFixedVariant }}`  | `633B48`       |
| `Outline`                 | `{{ .Outline }}`                 | `938F99`       |
| `OutlineVariant`          | `{{ .OutlineVariant }}`          | `C4C7C5`       |
| `Primary`                 | `{{ .Primary }}`                 | `BB86FC`       |
| `PrimaryContainer`        | `{{ .PrimaryContainer }}`        | `6200EE`       |
| `PrimaryFixed`            | `{{ .PrimaryFixed }}`            | `EADDFF`       |
| `PrimaryFixedDim`         | `{{ .PrimaryFixedDim }}`         | `D0BCFF`       |
| `Scrim`                   | `{{ .Scrim }}`                   | `000000`       |
| `Secondary`               | `{{ .Secondary }}`               | `03DAC6`       |
| `SecondaryContainer`      | `{{ .SecondaryContainer }}`      | `018786`       |
| `SecondaryFixed`          | `{{ .SecondaryFixed }}`          | `E8DEF8`       |
| `SecondaryFixedDim`       | `{{ .SecondaryFixedDim }}`       | `CCC2DC`       |
| `Shadow`                  | `{{ .Shadow }}`                  | `000000`       |
| `Surface`                 | `{{ .Surface }}`                 | `121212`       |
| `SurfaceBright`           | `{{ .SurfaceBright }}`           | `373737`       |
| `SurfaceContainer`        | `{{ .SurfaceContainer }}`        | `211F26`       |
| `SurfaceContainerHigh`    | `{{ .SurfaceContainerHigh }}`    | `2B2930`       |
| `SurfaceContainerHighest` | `{{ .SurfaceContainerHighest }}` | `36343B`       |
| `SurfaceContainerLow`     | `{{ .SurfaceContainerLow }}`     | `1D1B20`       |
| `SurfaceContainerLowest`  | `{{ .SurfaceContainerLowest }}`  | `0F0D13`       |
| `SurfaceDim`              | `{{ .SurfaceDim }}`              | `121212`       |
| `SurfaceTint`             | `{{ .SurfaceTint }}`             | `BB86FC`       |
| `SurfaceVariant`          | `{{ .SurfaceVariant }}`          | `49454F`       |
| `Tertiary`                | `{{ .Tertiary }}`                | `03DAC6`       |
| `TertiaryContainer`       | `{{ .TertiaryContainer }}`       | `3700B3`       |
| `TertiaryFixed`           | `{{ .TertiaryFixed }}`           | `FFD8E4`       |
| `TertiaryFixedDim`        | `{{ .TertiaryFixedDim }}`        | `EFB8C8`       |

## ColorValue Representations

You can use these suffixes on any `ColorValue` field (Material fields or `.Colors` elements):

| Representation      | Example Template                   | Output Example              |
| ------------------- | ---------------------------------- | --------------------------- |
| `HexRGB`            | `{{ .Primary.HexRGB }}`            | `BB86FC`                    |
| `TrimmedHexRGB`     | `{{ .Error.TrimmedHexRGB }}`       | `B3261E`                    |
| `HexRGBA`           | `{{ .OnSurface.HexRGBA }}`         | `E0E0E0FF`                  |
| `TrimmedHexRGBA`    | `{{ .Background.TrimmedHexRGBA }}` | `121212FF`                  |
| `RGB`               | `{{ .Primary.RGB }}`               | `rgb(187, 134, 252)`        |
| `TrimmedRGB`        | `{{ .Primary.TrimmedRGB }}`        | `187, 134, 252`             |
| `RGBA`              | `{{ .Primary.RGBA }}`              | `rgba(187, 134, 252, 255)`  |
| `TrimmedRGBA`       | `{{ .Primary.TrimmedRGBA }}`       | `187, 134, 252, 255`        |
| `LinearRGB`         | `{{ .Primary.LinearRGB }}`         | `rgb(0.73, 0.53, 0.99)`     |
| `TrimmedLinearRGB`  | `{{ .Primary.TrimmedLinearRGB }}`  | `0.73, 0.53, 0.99`          |
| `LinearRGBA`        | `{{ .Primary.LinearRGBA }}`        | `rgba(0.73, 0.53, 0.99, 1)` |
| `TrimmedLinearRGBA` | `{{ .Primary.TrimmedLinearRGBA }}` | `0.73, 0.53, 0.99, 1`       |
| `Red`               | `{{ .Primary.Red }}`               | `187`                       |
| `Green`             | `{{ .Primary.Green }}`             | `134`                       |
| `Blue`              | `{{ .Primary.Blue }}`              | `252`                       |
| `Alpha`             | `{{ .Primary.Alpha }}`             | `255`                       |

## Color Name Cases (for `.Colors` elements)

| Case     | Example Template     | Output Example |
| -------- | -------------------- | -------------- |
| `Snake`  | `{{ .Name.Snake }}`  | `on_primary`   |
| `Camel`  | `{{ .Name.Camel }}`  | `onPrimary`    |
| `Kebab`  | `{{ .Name.Kebab }}`  | `on-primary`   |
| `Pascal` | `{{ .Name.Pascal }}` | `OnPrimary`    |

## Colors Looping Example

```go
{{ range .Colors }}
## {{ .Name.Pascal }} Color
- All cases:
  Snake: `{{ .Name.Snake }}`,
  Camel: `{{ .Name.Camel }}`,
  Kebab: `{{ .Name.Kebab }}`,
  Pascal: `{{ .Name.Pascal }}`
- Formats:
  Hex: {{ .Color.HexRGB }}
  Trimmed RGB: {{ .Color.TrimmedRGB }}
  RGBA: {{ .Color.RGBA }}
  Linear: {{ .Color.LinearRGB }}
  Components: R={{ .Color.Red }} G={{ .Color.Green }} B={{ .Color.Blue }} A={{ .Color.Alpha }}
{{ end }}
```

## Example Template

```go
{{/* Rong theme example */}}

primary = "{{ .Primary }}"
background = "{{ .Background }}"
```

::: tip

- Use `{{-` and `-}}` to trim whitespace.
- Run `rong` to test and preview your template.

:::

</span>
