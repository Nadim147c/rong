package material

import "github.com/Nadim147c/goyou/color"

type MaterialColor struct {
	Background              color.ARGB `json:"background"`
	Error                   color.ARGB `json:"error"`
	ErrorContainer          color.ARGB `json:"error_container"`
	InverseOnSurface        color.ARGB `json:"inverse_on_surface"`
	InversePrimary          color.ARGB `json:"inverse_primary"`
	InverseSurface          color.ARGB `json:"inverse_surface"`
	OnBackground            color.ARGB `json:"on_background"`
	OnError                 color.ARGB `json:"on_error"`
	OnErrorContainer        color.ARGB `json:"on_error_container"`
	OnPrimary               color.ARGB `json:"on_primary"`
	OnPrimaryContainer      color.ARGB `json:"on_primary_container"`
	OnPrimaryFixed          color.ARGB `json:"on_primary_fixed"`
	OnPrimaryFixedVariant   color.ARGB `json:"on_primary_fixed_variant"`
	OnSecondary             color.ARGB `json:"on_secondary"`
	OnSecondaryContainer    color.ARGB `json:"on_secondary_container"`
	OnSecondaryFixed        color.ARGB `json:"on_secondary_fixed"`
	OnSecondaryFixedVariant color.ARGB `json:"on_secondary_fixed_variant"`
	OnSurface               color.ARGB `json:"on_surface"`
	OnSurfaceVariant        color.ARGB `json:"on_surface_variant"`
	OnTertiary              color.ARGB `json:"on_tertiary"`
	OnTertiaryContainer     color.ARGB `json:"on_tertiary_container"`
	OnTertiaryFixed         color.ARGB `json:"on_tertiary_fixed"`
	OnTertiaryFixedVariant  color.ARGB `json:"on_tertiary_fixed_variant"`
	Outline                 color.ARGB `json:"outline"`
	OutlineVariant          color.ARGB `json:"outline_variant"`
	Primary                 color.ARGB `json:"primary"`
	PrimaryContainer        color.ARGB `json:"primary_container"`
	PrimaryFixed            color.ARGB `json:"primary_fixed"`
	PrimaryFixedDim         color.ARGB `json:"primary_fixed_dim"`
	Scrim                   color.ARGB `json:"scrim"`
	Secondary               color.ARGB `json:"secondary"`
	SecondaryContainer      color.ARGB `json:"secondary_container"`
	SecondaryFixed          color.ARGB `json:"secondary_fixed"`
	SecondaryFixedDim       color.ARGB `json:"secondary_fixed_dim"`
	Shadow                  color.ARGB `json:"shadow"`
	Surface                 color.ARGB `json:"surface"`
	SurfaceBright           color.ARGB `json:"surface_bright"`
	SurfaceContainer        color.ARGB `json:"surface_container"`
	SurfaceContainerHigh    color.ARGB `json:"surface_container_high"`
	SurfaceContainerHighest color.ARGB `json:"surface_container_highest"`
	SurfaceContainerLow     color.ARGB `json:"surface_container_low"`
	SurfaceContainerLowest  color.ARGB `json:"surface_container_lowest"`
	SurfaceDim              color.ARGB `json:"surface_dim"`
	SurfaceTint             color.ARGB `json:"surface_tint"`
	SurfaceVariant          color.ARGB `json:"surface_variant"`
	Tertiary                color.ARGB `json:"tertiary"`
	TertiaryContainer       color.ARGB `json:"tertiary_container"`
	TertiaryFixed           color.ARGB `json:"tertiary_fixed"`
	TertiaryFixedDim        color.ARGB `json:"tertiary_fixed_dim"`
}

func MaterialColorFromMap(dcs map[string]color.ARGB) MaterialColor {
	get := func(name string) color.ARGB {
		if dc, ok := dcs[name]; ok {
			return dc
		}
		return 0 // default ARGB (fully transparent black)
	}

	return MaterialColor{
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
