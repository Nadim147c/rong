package models

import (
	"encoding/json"
	"strings"
	"unicode"

	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/rong/v4/internal/material"
)

// CustomColors is map of user defined custom colors
type CustomColors map[string]FormatedColor

func pascalToSnake(s string) string {
	var out []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				out = append(out, '_')
			}
			out = append(out, unicode.ToLower(r))
		} else {
			out = append(out, r)
		}
	}
	return string(out)
}

func snakeToPascal(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}

// MarshalJSON json implements json.Marshaller
func (c CustomColors) MarshalJSON() ([]byte, error) {
	tmp := make(map[string]FormatedColor, len(c))
	for k, v := range c {
		tmp[pascalToSnake(k)] = v
	}
	return json.Marshal(tmp)
}

// UnmarshalJSON json implements json.Unmarshaller
func (c *CustomColors) UnmarshalJSON(data []byte) error {
	tmp := map[string]FormatedColor{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	out := make(CustomColors, len(tmp))
	for snakeKey, v := range tmp {
		out[snakeToPascal(snakeKey)] = v
	}

	*c = out
	return nil
}

// Material contains all material colors
type Material struct {
	Custom                  CustomColors  `json:"custom"`
	Background              FormatedColor `json:"background"`
	Error                   FormatedColor `json:"error"`
	ErrorContainer          FormatedColor `json:"error_container"`
	InverseOnSurface        FormatedColor `json:"inverse_on_surface"`
	InversePrimary          FormatedColor `json:"inverse_primary"`
	InverseSurface          FormatedColor `json:"inverse_surface"`
	OnBackground            FormatedColor `json:"on_background"`
	OnError                 FormatedColor `json:"on_error"`
	OnErrorContainer        FormatedColor `json:"on_error_container"`
	OnPrimary               FormatedColor `json:"on_primary"`
	OnPrimaryContainer      FormatedColor `json:"on_primary_container"`
	OnPrimaryFixed          FormatedColor `json:"on_primary_fixed"`
	OnPrimaryFixedVariant   FormatedColor `json:"on_primary_fixed_variant"`
	OnSecondary             FormatedColor `json:"on_secondary"`
	OnSecondaryContainer    FormatedColor `json:"on_secondary_container"`
	OnSecondaryFixed        FormatedColor `json:"on_secondary_fixed"`
	OnSecondaryFixedVariant FormatedColor `json:"on_secondary_fixed_variant"`
	OnSurface               FormatedColor `json:"on_surface"`
	OnSurfaceVariant        FormatedColor `json:"on_surface_variant"`
	OnTertiary              FormatedColor `json:"on_tertiary"`
	OnTertiaryContainer     FormatedColor `json:"on_tertiary_container"`
	OnTertiaryFixed         FormatedColor `json:"on_tertiary_fixed"`
	OnTertiaryFixedVariant  FormatedColor `json:"on_tertiary_fixed_variant"`
	Outline                 FormatedColor `json:"outline"`
	OutlineVariant          FormatedColor `json:"outline_variant"`
	Primary                 FormatedColor `json:"primary"`
	PrimaryContainer        FormatedColor `json:"primary_container"`
	PrimaryFixed            FormatedColor `json:"primary_fixed"`
	PrimaryFixedDim         FormatedColor `json:"primary_fixed_dim"`
	Scrim                   FormatedColor `json:"scrim"`
	Secondary               FormatedColor `json:"secondary"`
	SecondaryContainer      FormatedColor `json:"secondary_container"`
	SecondaryFixed          FormatedColor `json:"secondary_fixed"`
	SecondaryFixedDim       FormatedColor `json:"secondary_fixed_dim"`
	Shadow                  FormatedColor `json:"shadow"`
	Surface                 FormatedColor `json:"surface"`
	SurfaceBright           FormatedColor `json:"surface_bright"`
	SurfaceContainer        FormatedColor `json:"surface_container"`
	SurfaceContainerHigh    FormatedColor `json:"surface_container_high"`
	SurfaceContainerHighest FormatedColor `json:"surface_container_highest"`
	SurfaceContainerLow     FormatedColor `json:"surface_container_low"`
	SurfaceContainerLowest  FormatedColor `json:"surface_container_lowest"`
	SurfaceDim              FormatedColor `json:"surface_dim"`
	SurfaceTint             FormatedColor `json:"surface_tint"`
	SurfaceVariant          FormatedColor `json:"surface_variant"`
	Tertiary                FormatedColor `json:"tertiary"`
	TertiaryContainer       FormatedColor `json:"tertiary_container"`
	TertiaryFixed           FormatedColor `json:"tertiary_fixed"`
	TertiaryFixedDim        FormatedColor `json:"tertiary_fixed_dim"`
}

// NewMaterial return material color type
func NewMaterial(
	colorMap map[string]color.ARGB,
	customColors map[string]material.CustomColor,
) Material {
	get := func(name string) FormatedColor {
		if dc, ok := colorMap[name]; ok {
			return NewFormatedColor(dc)
		}
		return NewFormatedColor(0) // default ARGB (fully transparent black)
	}

	custom := make(CustomColors, len(customColors)*4)
	for name, col := range customColors {
		name := toCamelCase(strings.ToLower(name), true)
		custom[name] = NewFormatedColor(col.Color)
		custom["On"+name] = NewFormatedColor(col.Color)
		custom[name+"Container"] = NewFormatedColor(col.Color)
		custom["On"+name+"Container"] = NewFormatedColor(col.Color)
	}

	return Material{
		Custom:                  custom,
		Background:              get("background"),
		Error:                   get("error"),
		ErrorContainer:          get("error_container"),
		InverseOnSurface:        get("inverse_on_surface"),
		InversePrimary:          get("inverse_primary"),
		InverseSurface:          get("inverse_surface"),
		OnBackground:            get("on_background"),
		OnError:                 get("on_error"),
		OnErrorContainer:        get("on_error_container"),
		OnPrimary:               get("on_primary"),
		OnPrimaryContainer:      get("on_primary_container"),
		OnPrimaryFixed:          get("on_primary_fixed"),
		OnPrimaryFixedVariant:   get("on_primary_fixed_variant"),
		OnSecondary:             get("on_secondary"),
		OnSecondaryContainer:    get("on_secondary_container"),
		OnSecondaryFixed:        get("on_secondary_fixed"),
		OnSecondaryFixedVariant: get("on_secondary_fixed_variant"),
		OnSurface:               get("on_surface"),
		OnSurfaceVariant:        get("on_surface_variant"),
		OnTertiary:              get("on_tertiary"),
		OnTertiaryContainer:     get("on_tertiary_container"),
		OnTertiaryFixed:         get("on_tertiary_fixed"),
		OnTertiaryFixedVariant:  get("on_tertiary_fixed_variant"),
		Outline:                 get("outline"),
		OutlineVariant:          get("outline_variant"),
		Primary:                 get("primary"),
		PrimaryContainer:        get("primary_container"),
		PrimaryFixed:            get("primary_fixed"),
		PrimaryFixedDim:         get("primary_fixed_dim"),
		Scrim:                   get("scrim"),
		Secondary:               get("secondary"),
		SecondaryContainer:      get("secondary_container"),
		SecondaryFixed:          get("secondary_fixed"),
		SecondaryFixedDim:       get("secondary_fixed_dim"),
		Shadow:                  get("shadow"),
		Surface:                 get("surface"),
		SurfaceBright:           get("surface_bright"),
		SurfaceContainer:        get("surface_container"),
		SurfaceContainerHigh:    get("surface_container_high"),
		SurfaceContainerHighest: get("surface_container_highest"),
		SurfaceContainerLow:     get("surface_container_low"),
		SurfaceContainerLowest:  get("surface_container_lowest"),
		SurfaceDim:              get("surface_dim"),
		SurfaceTint:             get("surface_tint"),
		SurfaceVariant:          get("surface_variant"),
		Tertiary:                get("tertiary"),
		TertiaryContainer:       get("tertiary_container"),
		TertiaryFixed:           get("tertiary_fixed"),
		TertiaryFixedDim:        get("tertiary_fixed_dim"),
	}
}
