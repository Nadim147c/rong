package models

import (
	"github.com/Nadim147c/material/color"
)

// Material contains all material colors
type Material struct {
	Background              ColorValue `json:"background" toml:"background"`
	Error                   ColorValue `json:"error" toml:"error"`
	ErrorContainer          ColorValue `json:"error_container" toml:"error_container"`
	InverseOnSurface        ColorValue `json:"inverse_on_surface" toml:"inverse_on_surface"`
	InversePrimary          ColorValue `json:"inverse_primary" toml:"inverse_primary"`
	InverseSurface          ColorValue `json:"inverse_surface" toml:"inverse_surface"`
	OnBackground            ColorValue `json:"on_background" toml:"on_background"`
	OnError                 ColorValue `json:"on_error" toml:"on_error"`
	OnErrorContainer        ColorValue `json:"on_error_container" toml:"on_error_container"`
	OnPrimary               ColorValue `json:"on_primary" toml:"on_primary"`
	OnPrimaryContainer      ColorValue `json:"on_primary_container" toml:"on_primary_container"`
	OnPrimaryFixed          ColorValue `json:"on_primary_fixed" toml:"on_primary_fixed"`
	OnPrimaryFixedVariant   ColorValue `json:"on_primary_fixed_variant" toml:"on_primary_fixed_variant"`
	OnSecondary             ColorValue `json:"on_secondary" toml:"on_secondary"`
	OnSecondaryContainer    ColorValue `json:"on_secondary_container" toml:"on_secondary_container"`
	OnSecondaryFixed        ColorValue `json:"on_secondary_fixed" toml:"on_secondary_fixed"`
	OnSecondaryFixedVariant ColorValue `json:"on_secondary_fixed_variant" toml:"on_secondary_fixed_variant"`
	OnSurface               ColorValue `json:"on_surface" toml:"on_surface"`
	OnSurfaceVariant        ColorValue `json:"on_surface_variant" toml:"on_surface_variant"`
	OnTertiary              ColorValue `json:"on_tertiary" toml:"on_tertiary"`
	OnTertiaryContainer     ColorValue `json:"on_tertiary_container" toml:"on_tertiary_container"`
	OnTertiaryFixed         ColorValue `json:"on_tertiary_fixed" toml:"on_tertiary_fixed"`
	OnTertiaryFixedVariant  ColorValue `json:"on_tertiary_fixed_variant" toml:"on_tertiary_fixed_variant"`
	Outline                 ColorValue `json:"outline" toml:"outline"`
	OutlineVariant          ColorValue `json:"outline_variant" toml:"outline_variant"`
	Primary                 ColorValue `json:"primary" toml:"primary"`
	PrimaryContainer        ColorValue `json:"primary_container" toml:"primary_container"`
	PrimaryFixed            ColorValue `json:"primary_fixed" toml:"primary_fixed"`
	PrimaryFixedDim         ColorValue `json:"primary_fixed_dim" toml:"primary_fixed_dim"`
	Scrim                   ColorValue `json:"scrim" toml:"scrim"`
	Secondary               ColorValue `json:"secondary" toml:"secondary"`
	SecondaryContainer      ColorValue `json:"secondary_container" toml:"secondary_container"`
	SecondaryFixed          ColorValue `json:"secondary_fixed" toml:"secondary_fixed"`
	SecondaryFixedDim       ColorValue `json:"secondary_fixed_dim" toml:"secondary_fixed_dim"`
	Shadow                  ColorValue `json:"shadow" toml:"shadow"`
	Surface                 ColorValue `json:"surface" toml:"surface"`
	SurfaceBright           ColorValue `json:"surface_bright" toml:"surface_bright"`
	SurfaceContainer        ColorValue `json:"surface_container" toml:"surface_container"`
	SurfaceContainerHigh    ColorValue `json:"surface_container_high" toml:"surface_container_high"`
	SurfaceContainerHighest ColorValue `json:"surface_container_highest" toml:"surface_container_highest"`
	SurfaceContainerLow     ColorValue `json:"surface_container_low" toml:"surface_container_low"`
	SurfaceContainerLowest  ColorValue `json:"surface_container_lowest" toml:"surface_container_lowest"`
	SurfaceDim              ColorValue `json:"surface_dim" toml:"surface_dim"`
	SurfaceTint             ColorValue `json:"surface_tint" toml:"surface_tint"`
	SurfaceVariant          ColorValue `json:"surface_variant" toml:"surface_variant"`
	Tertiary                ColorValue `json:"tertiary" toml:"tertiary"`
	TertiaryContainer       ColorValue `json:"tertiary_container" toml:"tertiary_container"`
	TertiaryFixed           ColorValue `json:"tertiary_fixed" toml:"tertiary_fixed"`
	TertiaryFixedDim        ColorValue `json:"tertiary_fixed_dim" toml:"tertiary_fixed_dim"`
}

// MaterialFromMap return material color type
func MaterialFromMap(dcs map[string]color.ARGB) Material {
	get := func(name string) ColorValue {
		if dc, ok := dcs[name]; ok {
			return NewColorValue(dc)
		}
		return NewColorValue(0) // default ARGB (fully transparent black)
	}

	return Material{
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
