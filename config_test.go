package main

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"

	"github.com/Nadim147c/material/v2/color"
	icmd "github.com/Nadim147c/rong/v4/cmd"
	"github.com/Nadim147c/rong/v4/internal/config"
	"github.com/Nadim147c/rong/v4/internal/config/enums"
)

func writeTempFile(t *testing.T, filename string, content string) string {
	t.Helper()

	dir := t.TempDir()
	fullPath := filepath.Join(dir, filename)

	err := os.WriteFile(fullPath, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("failed to write temp file %s: %v", fullPath, err)
	}

	return fullPath
}

func TestConfig(t *testing.T) {
	testdata := []struct {
		name     string
		filename string
		text     string
		config   map[any]any
	}{
		{
			name:     "simple config",
			filename: "simple.toml",
			text: `
      [base16]
      method = "dynamic"
      `,
			config: map[any]any{
				config.Base16Method: enums.Base16MethodDynamic,
			},
		},
		{
			name:     "basic config",
			filename: "basic.toml",
			text: `
      [base16]
      method="static"
      blend=0.22
      `,
			config: map[any]any{
				config.Base16Method: enums.Base16MethodStatic,
				config.Base16Blend:  0.22,
			},
		},
		{
			name:     "basic yaml config",
			filename: "basic.yaml",
			text: `
      base16:
        method: static
        blend: 0.1
      `,
			config: map[any]any{
				config.Base16Method: enums.Base16MethodStatic,
				config.Base16Blend:  0.1,
			},
		},
		{
			name:     "custom colors config in yaml",
			filename: "custom.yaml",
			text: `
      material:
        custom:
          blend: 0.35
          colors:
            orange: '#FFA522'
            purple: '#800080'
      `,
			config: map[any]any{
				config.MaterialCustomColors: map[string]color.ARGB{
					"orange": color.ARGBFromHexMust("#FFA522"),
					"purple": color.ARGBFromHexMust("#800080"),
				},
				config.MaterialCustomBlend: 0.35,
			},
		},
	}

	for tt := range slices.Values(testdata) {
		t.Run(tt.name, func(t *testing.T) {
			cfgFile := writeTempFile(t, tt.filename, tt.text)
			_cmd := *icmd.Command
			cmd := &_cmd
			cmd.SetArgs([]string{"video", "--config", cfgFile, "/dev/null"})
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)
			cmd.ExecuteContext(t.Context())

			for key, want := range tt.config {
				rv := reflect.ValueOf(key)
				got := rv.MethodByName("Value").Call(nil)[0].Interface()
				t.Log(want, got)
				if !reflect.DeepEqual(got, want) {
					t.Errorf("viper config and watned config did not match: want=%v, got=%v", want, got)
				}
			}
		})
	}
}
