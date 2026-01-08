package config

import (
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/Nadim147c/material/v2/color"
	"github.com/carapace-sh/carapace"
	"github.com/spf13/cast"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// option represents a configuration option with a specific type.
type option[T any] struct {
	short   string
	key     string
	defval  T
	desc    string
	typeStr string
	caster  func(any) (T, error)
}

// Default returns the default value of the option.
func (o *option[T]) Default() T {
	return o.defval
}

// SetValue sets the current value in the configuration store.
func (o *option[T]) SetValue(v T) {
	viper.Set(o.key, v)
}

func isBool(a any) bool {
	if _, ok := a.(bool); ok {
		return true
	}
	s := cast.ToString(a)
	return s == "true" || s == "false"
}

// Value retrieves the current value from the configuration store. Returns the
// default value if conversion fails.
func (o *option[T]) Value() T {
	viperValue := viper.Get(o.key)
	slog.Debug("Config key accessed", "key", o.Key(), "raw-value", viperValue)
	if v, ok := viperValue.(T); ok {
		return v
	}

	if o.Key() == Base16Blend.Key() && isBool(viperValue) {
		slog.Error("Config option has been changed. Use 0 to disable blending.",
			"option", o.Key(),
			"used", viperValue,
			"want", "a float between 0 to 1",
		)
		viperValue = Base16Blend.Default()
	}

	v, err := o.caster(viperValue)
	if err != nil {
		slog.Error(
			"Failed to convert config value",
			"key", o.key,
			"value", viperValue,
			"error", err,
		)
		return o.defval // Don't exit on error, return default
	}
	slog.Debug("Config value returned", "key", o.Key(), "value", v)
	return v
}

// RegisterFlag registers the option with a flag set.
func (o *option[T]) RegisterFlag(set *pflag.FlagSet) {
	set.VarP(o, o.key, o.short, o.desc)
}

// Set implements pflag.Value interface.
// Converts string to option type and stores it.
func (o *option[T]) Set(s string) error {
	v, err := o.caster(s)
	if err != nil {
		return err
	}

	viper.Set(o.key, v)
	return nil
}

// String implements pflag.Value interface.
// WARNING: This method is only for pflag interface compatibility.
// Do not use for general string representation.
func (o *option[T]) String() string {
	return cast.ToString(o.Value())
}

// Type returns the type description for pflag.
func (o *option[T]) Type() string {
	return o.typeStr
}

// Key returns the configuration key name.
func (o *option[T]) Key() string {
	return o.key
}

// Track registered keys to prevent duplicates.
var registeredKeys = map[string]struct{}{}

// register tracks used keys to prevent duplicates.
func register(keys ...string) {
	for _, key := range keys {
		if key == "" {
			continue
		}
		if _, ok := registeredKeys[key]; ok {
			panic(fmt.Sprintf("key %q already has been used", key))
		}
		registeredKeys[key] = struct{}{}
	}
}

// newOption creates a generic configuration option.
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
		key:     key,
		short:   short,
		defval:  defval,
		desc:    desc,
		typeStr: typeName,
		caster:  caster,
	}
}

// boolOption represents a boolean configuration option with special handling for flags.
type boolOption struct{ *option[bool] }

// RegisterFlag registers boolean flag with "no-" counterpart for negation.
func (c *boolOption) RegisterFlag(set *pflag.FlagSet) {
	f := set.VarPF(c, c.key, c.short, c.desc)
	f.NoOptDefVal = "<bool>"

	// Add "no-" prefix flag for negation
	noflag := "no-" + c.key
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

// Set implements pflag.Value for boolean flag.
func (c *boolOption) Set(s string) error {
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

// String returns empty string for boolean flags.
// WARNING: This method is only for pflag interface compatibility.
func (c *boolOption) String() string { return "" }

// newBoolOption creates a new boolean configuration option.
func newBoolOption(short, key string, defval bool, desc string) *boolOption {
	return &boolOption{
		option: newOption(short, key, defval, desc, "bool", cast.ToBoolE),
	}
}

// newStringOption creates a new string configuration option.
func newStringOption(short, key string, defval string, desc string) *option[string] {
	return newOption(short, key, defval, desc, "string", cast.ToStringE)
}

// newStringOption creates a new string configuration option.
func newIntOption(short, key string, defval int, desc string) *option[int] {
	return newOption(short, key, defval, desc, "int", cast.ToIntE)
}

// castDuration casts the type to time.duration. The differenace between
// castDuration and cast.ToDurationE is the default duration unit is second in
// castDuration where cast.ToDurationE uses nanoseocond. It make sense to go dev
// to use nanoseocond but everyone else expect second as default time.
func castDuration(a any) (time.Duration, error) {
	s, ok := a.(string)
	if !ok {
		return 0, fmt.Errorf("failed to convert %v to duration", a) //nolint
	}
	if !strings.ContainsAny(s, "nsuµmh") {
		s += "s"
	}
	return time.ParseDuration(s)
}

// newStringOption creates a new string configuration option.
func newDurationOption(short, key string, defval time.Duration, desc string) *option[time.Duration] {
	return newOption(short, key, defval, desc, "duration", castDuration)
}

// formatFloat formats float with 2 decimal precision.
func formatFloat(v float64) string {
	// round to 2 decimals first
	v = math.Round(v*100) / 100
	return strconv.FormatFloat(v, 'f', -1, 64)
}

// newRangeFloatOption creates a float option with range validation.
func newRangeFloatOption(short, key string, defval float64, up, low float64, desc string) *option[float64] {
	description := fmt.Sprintf("%s. range [%s, %s].", desc, formatFloat(low), formatFloat(up))
	return newOption(short, key, defval, description, "float", func(a any) (float64, error) {
		f, err := cast.ToFloat64E(a)
		if err != nil {
			return 0, err
		}
		if up < f || f < low {
			return 0, fmt.Errorf("value %g must be between %g and %g", f, low, up) //nolint
		}
		return f, nil
	})
}

// colorOption represents a color configuration option.
type colorOption struct{ *option[color.ARGB] }

// RegisterFlag registers count flag with special handling for incremental
// counting.
func (c *colorOption) RegisterFlag(set *pflag.FlagSet) {
	set.VarP(c, c.key, c.short, c.desc)
}

// String returns the color as an ANSI-formatted string
// WARNING: This method is only for pflag interface compatibility.
// For general use, prefer Value().String() or Value().AnsiBg().
func (c colorOption) String() string {
	col := c.Value()
	return col.AnsiFg(col.String())
}

// newColorOption creates a new color configuration option.
func newColorOption(short, key, defval, desc string) *colorOption { //nolint
	return &colorOption{
		option: newOption(short, key, color.ARGBFromHexMust(defval), desc, "color", func(a any) (color.ARGB, error) {
			s, err := cast.ToStringE(a)
			if err != nil {
				return 0, err
			}
			return color.ARGBFromHex(s)
		}),
	}
}

// countOption represents a counter configuration option (e.g., verbosity
// levels).
type countOption struct{ *option[int] }

// RegisterFlag registers count flag with special handling for incremental
// counting.
func (c *countOption) RegisterFlag(set *pflag.FlagSet) {
	f := set.VarPF(c, c.key, c.short, c.desc)
	f.NoOptDefVal = "<count>"
}

// Set implements pflag.Value for count flags.
func (c *countOption) Set(s string) error {
	if s == "<count>" {
		// Increment the current value
		c.SetValue(viper.GetInt(c.key) + 1)
		return nil
	}
	i, err := c.caster(s)
	if err != nil {
		return err
	}
	c.SetValue(i)
	return nil
}

// newCountOption creates a new count configuration option.
func newCountOption(short, key string, desc string) *countOption {
	return &countOption{
		option: newOption(short, key, 0, desc, "...count", cast.ToIntE),
	}
}

// enumOption represents an enumeration configuration option.
type enumOption[T any] struct {
	*option[T]
	names []string
}

// newEnumOption creates a new enumeration configuration option.
func newEnumOption[T any](
	short, key string,
	defval T,
	desc string,
	names []string,
	parser func(s string) (T, error),
) *enumOption[T] {
	description := fmt.Sprintf("%s. choice: (%s).", desc, joinWithOr(names))
	CarapaceAction[key] = carapace.ActionValues(names...)
	return &enumOption[T]{
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

// joinWithOr joins items with Oxford comma and "or".
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

type kvOption[T any] struct {
	*option[map[string]T]
	caster func(string) (T, error)
}

func (o *kvOption[T]) RegisterFlag(set *pflag.FlagSet) {
	set.VarP(o, o.key, o.short, o.desc)
}

// Set sets a value from name=value.
func (o *kvOption[T]) Set(s string) error {
	name, value, found := strings.Cut(s, "=")
	if !found {
		return fmt.Errorf("Failed to convert %q to map of string, %s", s, o.typeStr) //nolint
	}

	v, err := o.caster(value)
	if err != nil {
		return err
	}

	viper.Set(o.key+"."+name, v)
	return nil
}

// newKvOption creates a new enumeration configuration option.
func newKvOption[T any](
	short, key string,
	defval map[string]T,
	desc string,
	typeStr string,
	parser func(s string) (T, error),
) *kvOption[T] {
	return &kvOption[T]{
		caster: parser,
		option: newOption(short, key, defval, desc, "name="+typeStr, func(a any) (map[string]T, error) {
			rawData, err := cast.ToStringMapE(a)
			if err != nil {
				return nil, err
			}

			if v, ok := any(rawData).(map[string]T); ok {
				return v, nil
			}

			m := make(map[string]T, len(rawData))
			for key, value := range rawData {
				v, err := parser(cast.ToString(value))
				if err != nil {
					return nil, err
				}
				m[key] = v
			}
			return m, nil
		}),
	}
}
