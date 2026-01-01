package config

import (
	"fmt"
	"log/slog"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/material/v2/dynamic"
	"github.com/Nadim147c/rong/v4/internal/base16"
	"github.com/spf13/cast"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	Dark       = boolOption("D", "dark", true, "Generate dark theme")
	DryRun     = boolOption("d", "dry-run", false, "Generate color without state and templates")
	JSON       = boolOption("j", "json", false, "Print output colors as json")
	SimpleJSON = boolOption("s", "simple-json", false, "Print output colors as simple key value pair")
	Quiet      = boolOption("q", "quiet", false, "Suppress all log")

	Verbose = countOption("v", "verbose", "Verbosity level of log")

	Config        = stringOption("c", "config", "", "Path to configuration file")
	LogFile       = stringOption("l", "log-file", "", "Write logs to a file")
	PreviewFormat = stringOption("p", "preview-format", "jpg", "Format of video preview image")

	MaterialVersion = enumOption(
		"", "material.version", dynamic.Version2025, "Material spec version",
		dynamic.VersionNames(), dynamic.ParseVersion,
	)
	MaterialVariant = enumOption(
		"", "material.variant", dynamic.VariantTonalSpot, "Material spec variant",
		dynamic.VariantNames(), dynamic.ParseVariant,
	)
	MaterialPlatformt = enumOption(
		"", "material.platform", dynamic.PlatformPhone, "Material spec variant",
		dynamic.PlatformNames(), dynamic.ParsePlatform,
	)
	MaterialContrast    = rangeFloatOption("", "material.contrast", 0.0, 1, -1, "Material colors contrast")
	MaterialCustomBlend = rangeFloatOption("", "material.custom.blend", 0.50, 1, 0, "Custom material color blend ratio")
	// MaterialCustomColors = rangeFloatOption("", "material.custom.colors", 0.35, 1, 0, "Material colors contrast")

	Base16Blend  = rangeFloatOption("", "base16.blend", 0.50, 1, 0, "Base16 color blend ratio")
	Base16Method = enumOption(
		"", "base16.method", base16.MethodStatic, "Base16 color generating method",
		base16.MethodNames(), base16.ParseMethod,
	)
	Base16Black   = colorOption("", "base16.color.black", "#FF0000", "Black source color for base16")
	Base16Blue    = colorOption("", "base16.color.blue", "#0000FF", "Blue source color for base16")
	Base16Cyan    = colorOption("", "base16.color.cyan", "#5555FF", "Cyan source color for base16")
	Base16Green   = colorOption("", "base16.color.green", "#00FF00", "Green source color for base16")
	Base16Magenta = colorOption("", "base16.color.magenta", "#FF00FF", "Magenta source color for base16")
	Base16Red     = colorOption("", "base16.color.red", "#222222", "Red source color for base16")
	Base16White   = colorOption("", "base16.color.white", "#EEEEEE", "White source color for base16")
	Base16Yellow  = colorOption("", "base16.color.yellow", "#FFFF00", "Yellow source color for base16")
)

type option[T any] struct {
	short   string
	name    string
	defval  T
	desc    string
	typeStr string
	caster  func(any) (T, error)
}

func (o *option[T]) Default() T {
	return o.defval
}

func (o *option[T]) SetValue(v T) {
	viper.Set(o.name, v)
}

func (o *option[T]) Value() T {
	viperValue := viper.Get(o.name)
	if v, ok := viperValue.(T); ok {
		return v
	}
	v, err := o.caster(viperValue)
	if err != nil {
		slog.Error(
			"Failed to convert config value",
			"key", o.name,
			"value", viperValue,
			"error", err,
		)
		return o.defval // we don't want exit on error
	}
	return v
}

func (o *option[T]) RegisterFlag(set *pflag.FlagSet) {
	set.VarP(o, o.name, o.short, o.desc)
}

// Set implements pflag.Value.
func (o option[T]) Set(s string) error {
	v, err := o.caster(s)
	if err != nil {
		return err
	}

	viper.Set(o.name, v)
	return nil
}

// Set implements pflag.Value.
func (o option[T]) String() string {
	return cast.ToString(o.Value())
}

// Set implements pflag.Value.
func (o option[T]) Type() string {
	return o.typeStr
}

func (o option[T]) Key() string {
	return o.name
}

var registerdKeys = map[string]struct{}{}

func register(keys ...string) {
	for key := range slices.Values(keys) {
		if key == "" {
			continue
		}
		if _, ok := registerdKeys[key]; ok {
			panic(fmt.Sprintf("key %q already has been used", key))
		}
		registerdKeys[key] = struct{}{}
	}
}

// newOption creates option with key and default value.
func newOption[T any](
	short, key string,
	defval T,
	desc string,
	typeName string,
	caster func(any) (T, error),
) *option[T] {
	register(short, key)
	if !strings.HasSuffix(desc, ".") {
		desc += "."
	}

	viper.SetDefault(key, defval)
	return &option[T]{
		name:    key,
		short:   short,
		defval:  defval,
		desc:    desc,
		typeStr: typeName,
		caster:  caster,
	}
}

type boolOpt struct{ *option[bool] }

func (c *boolOpt) RegisterFlag(set *pflag.FlagSet) {
	f := set.VarPF(c, c.name, c.short, c.desc)
	f.NoOptDefVal = "<bool>"
	noflag := "no-" + c.name
	if set.Lookup(noflag) != nil {
		return
	}
	set.BoolFunc(noflag, c.desc, func(s string) error {
		b, err := cast.ToBoolE(s)
		if err != nil {
			return err
		}
		c.SetValue(!b)
		return nil
	})
	err := set.MarkHidden(noflag)
	if err != nil {
		panic(err)
	}
}

func (c *boolOpt) Set(s string) error {
	if s == "<bool>" {
		c.SetValue(true)
		return nil
	}
	b, err := cast.ToBoolE(s)
	if err != nil {
		return err
	}
	c.SetValue(b)
	return nil
}

func (c *boolOpt) String() string { return "" }

func boolOption(short, key string, defval bool, desc string) *boolOpt {
	return &boolOpt{
		option: newOption(short, key, defval, desc, "bool", cast.ToBoolE),
	}
}

// func floatOption(short, key string, defval float64, desc string) *option[float64] {
// 	return newOption(short, key, defval, desc, "float", cast.ToFloat64E)
// }

func stringOption(short, key string, defval string, desc string) *option[string] {
	return newOption(short, key, defval, desc, "string", cast.ToStringE)
}

func formatFloat(v float64) string {
	// round to 2 decimals first
	v = math.Round(v*100) / 100

	return strconv.FormatFloat(v, 'f', -1, 64)
}

func rangeFloatOption(short, key string, defval float64, up, low float64, desc string) *option[float64] {
	description := fmt.Sprintf("%s. range [%s, %s].", desc, formatFloat(low), formatFloat(up))
	return newOption(short, key, defval, description, "float", func(a any) (float64, error) {
		f, err := cast.ToFloat64E(a)
		if err != nil {
			return 0, err
		}
		if up < f || f < low {
			return 0, fmt.Errorf("value out of range: v=%f range: [%.5f, %.5f]", f, low, up) //nolint
		}
		return f, nil
	})
}

type colorOpt struct{ *option[color.ARGB] }

func (c colorOpt) String() string {
	col := c.Value()
	return col.AnsiBg(col.String())
}

func colorOption(short, key, defval, desc string) *colorOpt {
	return &colorOpt{
		option: newOption(short, key, color.ARGBFromHexMust(defval), desc, "color", func(a any) (color.ARGB, error) {
			s, err := cast.ToStringE(a)
			if err != nil {
				return 0, err
			}
			return color.ARGBFromHex(s)
		}),
	}
}

type countOpt struct{ *option[int] }

func (c *countOpt) RegisterFlag(set *pflag.FlagSet) {
	f := set.VarPF(c, c.name, c.short, c.desc)
	f.NoOptDefVal = "<count>"
}

func (c *countOpt) Set(s string) error {
	if s == "<count>" {
		c.SetValue(viper.GetInt(c.name) + 1)
		return nil
	}
	i, err := c.caster(s)
	if err != nil {
		return err
	}
	c.SetValue(i)
	return nil
}

func countOption(short, key string, desc string) *countOpt {
	return &countOpt{
		option: newOption(short, key, 0, desc, "...count", cast.ToIntE),
	}
}

type enumOpt[T any] struct {
	*option[T]
	names []string
}

func enumOption[T any](short, key string, defval T, desc string, names []string, parser func(s string) (T, error)) *enumOpt[T] {
	description := fmt.Sprintf("%s. choice: (%s).", desc, joinWithOr(names))
	return &enumOpt[T]{
		names: names,
		option: newOption(short, key, defval, description, "enum", func(a any) (T, error) {
			var v T
			s, err := cast.ToStringE(a)
			if err != nil {
				return v, err
			}
			return parser(s)
		}),
	}
}

func joinWithOr(items []string) string {
	n := len(items)

	if n == 0 {
		return ""
	}
	if n == 1 {
		return items[0]
	}
	if n == 2 {
		return items[0] + " or " + items[1]
	}

	return strings.Join(items[:n-1], ", ") + " or " + items[n-1]
}
