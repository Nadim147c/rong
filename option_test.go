package main

import (
	"io"
	"slices"
	"testing"

	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/material/v2/dynamic"
	icmd "github.com/Nadim147c/rong/v4/cmd"
	"github.com/Nadim147c/rong/v4/internal/config"
	"github.com/Nadim147c/rong/v4/internal/config/enums"
	shlex "github.com/carapace-sh/carapace-shlex"
	"github.com/spf13/viper"
)

func fatal(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func TestOptions(t *testing.T) {
	prefix := "image --config=false"

	testdata := []struct {
		name    string
		command string
		key     string
		val     any
	}{
		{"2025 version", "--material.version 2025", config.MaterialVersion.Key(), dynamic.Version2025},
		{"2021 version", "--material.version 2021", config.MaterialVersion.Key(), dynamic.Version2021},
		{"monochrome variant", "--material.variant monochrome", config.MaterialVariant.Key(), dynamic.VariantMonochrome},
		{"material platform phone", "--material.platform phone", config.MaterialPlatformt.Key(), dynamic.PlatformPhone},
		{"material platform watch", "--material.platform watch", config.MaterialPlatformt.Key(), dynamic.PlatformWatch},
		{"material contrast positive", "--material.contrast 0.5", config.MaterialContrast.Key(), float64(0.5)},
		{"material contrast negative", "--material.contrast -0.5", config.MaterialContrast.Key(), float64(-0.5)},
		{"custom material blend", "--material.custom.blend 0.75", config.MaterialCustomBlend.Key(), float64(0.75)},
		{"base16 blend", "--base16.blend 0.25", config.Base16Blend.Key(), float64(0.25)},
		{"base16 method static", "--base16.method static", config.Base16Method.Key(), enums.Base16MethodStatic},
		{"base16 method dynamic", "--base16.method dynamic", config.Base16Method.Key(), enums.Base16MethodDynamic},
		{"base16 black color", "--base16.colors.black '#000000'", config.Base16Black.Key(), color.ARGBFromHexMust("#000000")},
		{"base16 white color", "--base16.colors.white '#FFFFFF'", config.Base16White.Key(), color.ARGBFromHexMust("#FFFFFF")},
	}
	for tt := range slices.Values(testdata) {
		t.Run(tt.name, func(t *testing.T) {
			command, err := shlex.Split(prefix + " " + tt.command)
			fatal(t, err)
			_cmd := *icmd.Command
			cmd := &_cmd
			cmd.SetArgs(command.Strings())
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)
			cmd.ExecuteContext(t.Context())

			viperValue := viper.Get(tt.key)
			if viperValue != tt.val {
				t.Errorf("configuration did not match: want (%v) got (%v)", viperValue, tt.val)
			}
		})
	}
}
