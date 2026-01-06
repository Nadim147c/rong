package config

import (
	"runtime"
	"time"

	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/material/v2/dynamic"
	"github.com/Nadim147c/rong/v4/internal/config/enums"
	"github.com/carapace-sh/carapace"
)

// Global configuration options.
var (
	CarapaceAction = carapace.ActionMap{
		Config.Key():  carapace.ActionFiles(),
		LogFile.Key(): carapace.ActionFiles(),
	}

	Dark       = newBoolOption("D", "dark", true, "Generate dark theme")
	DryRun     = newBoolOption("d", "dry-run", false, "Generate color without state and templates")
	JSON       = newBoolOption("j", "json", false, "Print output colors as json")
	SimpleJSON = newBoolOption("s", "simple-json", false, "Print output colors as simple key value pair")
	Quiet      = newBoolOption("q", "quiet", false, "Suppress all log")

	Verbose = newCountOption("v", "verbose", "Verbosity level of log")

	Config        = newStringOption("c", "config", "", "Path to configuration file")
	LogFile       = newStringOption("l", "log-file", "", "Write logs to a file")
	PreviewFormat = newEnumOption(
		"p", "preview-format", enums.PreviewFormatJpg, "Format of video preview image",
		enums.PreviewFormatNames(), enums.ParsePreviewFormat,
	)

	FFmpegFrames   = newIntOption("", "frames", 5, "Number of frames ffmpeg should process")
	FFmpegDuration = newDurationOption("", "duration", 5*time.Second, "Max number of time ffmpeg should process")
	Workers        = newIntOption("", "workers", runtime.GOMAXPROCS(runtime.NumCPU()), "Number for concurate thread to use")

	MaterialVersion = newEnumOption(
		"", "material.version", dynamic.Version2025, "Material spec version",
		dynamic.VersionNames(), dynamic.ParseVersion,
	)
	MaterialVariant = newEnumOption(
		"", "material.variant", dynamic.VariantTonalSpot, "Material spec variant",
		dynamic.VariantNames(), dynamic.ParseVariant,
	)
	MaterialPlatformt = newEnumOption(
		"", "material.platform", dynamic.PlatformPhone, "Material spec variant",
		dynamic.PlatformNames(), dynamic.ParsePlatform,
	)
	MaterialContrast     = newRangeFloatOption("", "material.contrast", 0.0, 1, -1, "Material colors contrast")
	MaterialCustomBlend  = newRangeFloatOption("", "material.custom.blend", 0.50, 1, 0, "Custom material color blend ratio")
	MaterialCustomColors = newKvOption("", "material.custom.colors", nil, "Custom material colors", "color", color.ARGBFromHex)

	Base16Blend  = newRangeFloatOption("", "base16.blend", 0.50, 1, 0, "Base16 color blend ratio")
	Base16Method = newEnumOption(
		"", "base16.method", enums.Base16MethodStatic, "Base16 color generating method",
		enums.Base16MethodNames(), enums.ParseBase16Method,
	)
	Base16Black   = newColorOption("", "base16.colors.black", "#FF0000", "Black source color for base16")
	Base16Blue    = newColorOption("", "base16.colors.blue", "#0000FF", "Blue source color for base16")
	Base16Cyan    = newColorOption("", "base16.colors.cyan", "#5555FF", "Cyan source color for base16")
	Base16Green   = newColorOption("", "base16.colors.green", "#00FF00", "Green source color for base16")
	Base16Magenta = newColorOption("", "base16.colors.magenta", "#FF00FF", "Magenta source color for base16")
	Base16Red     = newColorOption("", "base16.colors.red", "#222222", "Red source color for base16")
	Base16White   = newColorOption("", "base16.colors.white", "#EEEEEE", "White source color for base16")
	Base16Yellow  = newColorOption("", "base16.colors.yellow", "#FFFF00", "Yellow source color for base16")
)
