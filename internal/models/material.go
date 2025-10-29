package models

import (
	"github.com/Nadim147c/material/color"
)

// Material contains all material colors
type Material struct {
	Background              FormatedColor `json:"background"                 toml:"background"`
	Error                   FormatedColor `json:"error"                      toml:"error"`
	ErrorContainer          FormatedColor `json:"error_container"            toml:"error_container"`
	InverseOnSurface        FormatedColor `json:"inverse_on_surface"         toml:"inverse_on_surface"`
	InversePrimary          FormatedColor `json:"inverse_primary"            toml:"inverse_primary"`
	InverseSurface          FormatedColor `json:"inverse_surface"            toml:"inverse_surface"`
	OnBackground            FormatedColor `json:"on_background"              toml:"on_background"`
	OnError                 FormatedColor `json:"on_error"                   toml:"on_error"`
	OnErrorContainer        FormatedColor `json:"on_error_container"         toml:"on_error_container"`
	OnPrimary               FormatedColor `json:"on_primary"                 toml:"on_primary"`
	OnPrimaryContainer      FormatedColor `json:"on_primary_container"       toml:"on_primary_container"`
	OnPrimaryFixed          FormatedColor `json:"on_primary_fixed"           toml:"on_primary_fixed"`
	OnPrimaryFixedVariant   FormatedColor `json:"on_primary_fixed_variant"   toml:"on_primary_fixed_variant"`
	OnSecondary             FormatedColor `json:"on_secondary"               toml:"on_secondary"`
	OnSecondaryContainer    FormatedColor `json:"on_secondary_container"     toml:"on_secondary_container"`
	OnSecondaryFixed        FormatedColor `json:"on_secondary_fixed"         toml:"on_secondary_fixed"`
	OnSecondaryFixedVariant FormatedColor `json:"on_secondary_fixed_variant" toml:"on_secondary_fixed_variant"`
	OnSurface               FormatedColor `json:"on_surface"                 toml:"on_surface"`
	OnSurfaceVariant        FormatedColor `json:"on_surface_variant"         toml:"on_surface_variant"`
	OnTertiary              FormatedColor `json:"on_tertiary"                toml:"on_tertiary"`
	OnTertiaryContainer     FormatedColor `json:"on_tertiary_container"      toml:"on_tertiary_container"`
	OnTertiaryFixed         FormatedColor `json:"on_tertiary_fixed"          toml:"on_tertiary_fixed"`
	OnTertiaryFixedVariant  FormatedColor `json:"on_tertiary_fixed_variant"  toml:"on_tertiary_fixed_variant"`
	Outline                 FormatedColor `json:"outline"                    toml:"outline"`
	OutlineVariant          FormatedColor `json:"outline_variant"            toml:"outline_variant"`
	Primary                 FormatedColor `json:"primary"                    toml:"primary"`
	PrimaryContainer        FormatedColor `json:"primary_container"          toml:"primary_container"`
	PrimaryFixed            FormatedColor `json:"primary_fixed"              toml:"primary_fixed"`
	PrimaryFixedDim         FormatedColor `json:"primary_fixed_dim"          toml:"primary_fixed_dim"`
	Scrim                   FormatedColor `json:"scrim"                      toml:"scrim"`
	Secondary               FormatedColor `json:"secondary"                  toml:"secondary"`
	SecondaryContainer      FormatedColor `json:"secondary_container"        toml:"secondary_container"`
	SecondaryFixed          FormatedColor `json:"secondary_fixed"            toml:"secondary_fixed"`
	SecondaryFixedDim       FormatedColor `json:"secondary_fixed_dim"        toml:"secondary_fixed_dim"`
	Shadow                  FormatedColor `json:"shadow"                     toml:"shadow"`
	Surface                 FormatedColor `json:"surface"                    toml:"surface"`
	SurfaceBright           FormatedColor `json:"surface_bright"             toml:"surface_bright"`
	SurfaceContainer        FormatedColor `json:"surface_container"          toml:"surface_container"`
	SurfaceContainerHigh    FormatedColor `json:"surface_container_high"     toml:"surface_container_high"`
	SurfaceContainerHighest FormatedColor `json:"surface_container_highest"  toml:"surface_container_highest"`
	SurfaceContainerLow     FormatedColor `json:"surface_container_low"      toml:"surface_container_low"`
	SurfaceContainerLowest  FormatedColor `json:"surface_container_lowest"   toml:"surface_container_lowest"`
	SurfaceDim              FormatedColor `json:"surface_dim"                toml:"surface_dim"`
	SurfaceTint             FormatedColor `json:"surface_tint"               toml:"surface_tint"`
	SurfaceVariant          FormatedColor `json:"surface_variant"            toml:"surface_variant"`
	Tertiary                FormatedColor `json:"tertiary"                   toml:"tertiary"`
	TertiaryContainer       FormatedColor `json:"tertiary_container"         toml:"tertiary_container"`
	TertiaryFixed           FormatedColor `json:"tertiary_fixed"             toml:"tertiary_fixed"`
	TertiaryFixedDim        FormatedColor `json:"tertiary_fixed_dim"         toml:"tertiary_fixed_dim"`
}

// NewMaterial return material color type
func NewMaterial(colorMap map[string]color.ARGB) Material {
	get := func(name string) FormatedColor {
		if dc, ok := colorMap[name]; ok {
			return NewFormatedColor(dc)
		}
		return NewFormatedColor(0) // default ARGB (fully transparent black)
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
