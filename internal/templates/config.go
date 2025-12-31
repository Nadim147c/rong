package templates

import (
	"log/slog"
	"slices"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func getConfig(key string) (map[string][]string, error) {
	rawCfg := viper.Get(key)
	if rawCfg == nil {
		return map[string][]string{}, nil // user has not specified config
	}
	return cast.ToStringMapStringSliceE(rawCfg)
}

func convertThemes(links, installs, cmds map[string][]string) {
	themesCfg := viper.Get("themes")
	if themesCfg == nil {
		return
	}

	themes, err := cast.ToSliceE(themesCfg)
	if err != nil {
		slog.Error("failed to convert themes config")
		return
	}

	for theme := range slices.Values(themes) {
		block, err := cast.ToStringMapE(theme)
		if err != nil {
			slog.Error("failed to convert theme block to map", "error", err)
			continue
		}

		var target string

		if t, ok := block["target"]; !ok || t == "" {
			slog.Error("theme block does not have target")
			continue
		} else {
			target = cast.ToString(t)
		}

		if s, ok := block["installs"]; ok {
			installs[target] = append(installs[target], toStringSlice(s)...)
		}
		if s, ok := block["links"]; ok {
			links[target] = append(links[target], toStringSlice(s)...)
		}
		if s, ok := block["cmds"]; ok {
			cmds[target] = append(cmds[target], toStringSlice(s)...)
		}
	}
}

func toStringSlice(s any) []string {
	switch s := s.(type) {
	case string:
		return []string{s}
	case []string:
		return s
	}
	return []string{}
}
