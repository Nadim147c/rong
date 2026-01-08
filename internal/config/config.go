package config

import (
	"runtime"
	"time"

	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/material/v2/dynamic"
	"github.com/Nadim147c/rong/v5/internal/config/enums"
	"github.com/carapace-sh/carapace"
)

// Global configuration options.
var (
	CarapaceAction = carapace.ActionMap{
		Config.Key():  carapace.ActionFiles(),
		LogFile.Key(): carapace.ActionFiles(),
	}

	Dark       = newBoolOption("D", "dark", true, "Generate a dark color theme")
	DryRun     = newBoolOption("d", "dry-run", false, "Generate colors without writing templates")
	JSON       = newBoolOption("j", "json", false, "Output generated colors as JSON")
	SimpleJSON = newBoolOption("s", "simple-json", false, "Output colors as simple key-value JSON")
	Quiet      = newBoolOption("q", "quiet", false, "Disable all logs")

	Verbose = newCountOption("v", "verbose", "Increase log verbosity level")

	Config        = newStringOption("c", "config", "", "Path to a configuration file")
	LogFile       = newStringOption("l", "log-file", "", "Write logs to the specified file")
	PreviewFormat = newEnumOption(
		"p", "preview-format", enums.PreviewFormatJpg, "Output format for preview image",
		enums.PreviewFormatNames(), enums.ParsePreviewFormat,
	)

	FFmpegFrames   = newIntOption("", "frames", 5, "Number of frames to process with ffmpeg")
	FFmpegDuration = newDurationOption("", "duration", 5*time.Second, "Maximum ffmpeg processing duration")
	Workers        = newIntOption("", "workers", runtime.GOMAXPROCS(runtime.NumCPU()), "Number of worker threads to use")

	MaterialVersion = newEnumOption(
		"", "material.version", dynamic.Version2025, "Material Design specification version",
		dynamic.VersionNames(), dynamic.ParseVersion,
	)
	MaterialVariant = newEnumOption(
		"", "material.variant", dynamic.VariantTonalSpot, "Material color generation variant",
		dynamic.VariantNames(), dynamic.ParseVariant,
	)
	MaterialPlatformt = newEnumOption(
		"", "material.platform", dynamic.PlatformPhone, "Target Material platform",
		dynamic.PlatformNames(), dynamic.ParsePlatform,
	)
	MaterialContrast     = newRangeFloatOption("", "material.contrast", 0.0, 1, -1, "Adjust contrast of Material colors")
	MaterialCustomBlend  = newRangeFloatOption("", "material.custom.blend", 0.50, 1, 0, "Blend ratio for custom Material colors")
	MaterialCustomColors = newKvOption("", "material.custom.colors", nil, "Add custom Material colors", "color", color.ARGBFromHex)

	Base16Blend  = newRangeFloatOption("", "base16.blend", 0.50, 1, 0, "Blend ratio for Base16 color generation")
	Base16Method = newEnumOption(
		"", "base16.method", enums.Base16MethodStatic, "Base16 color generation method",
		enums.Base16MethodNames(), enums.ParseBase16Method,
	)
	Base16Black   = newColorOption("", "base16.colors.black", "#222222", "Base16 black source color")
	Base16Blue    = newColorOption("", "base16.colors.blue", "#0000FF", "Base16 blue source color")
	Base16Cyan    = newColorOption("", "base16.colors.cyan", "#5555FF", "Base16 cyan source color")
	Base16Green   = newColorOption("", "base16.colors.green", "#00FF00", "Base16 green source color")
	Base16Magenta = newColorOption("", "base16.colors.magenta", "#FF00FF", "Base16 magenta source color")
	Base16Red     = newColorOption("", "base16.colors.red", "#FF0000", "Base16 red source color")
	Base16White   = newColorOption("", "base16.colors.white", "#EEEEEE", "Base16 white source color")
	Base16Yellow  = newColorOption("", "base16.colors.yellow", "#FFFF00", "Base16 yellow source color")
)
