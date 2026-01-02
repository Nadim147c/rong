package color

import (
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/material/v2/dynamic"
	"github.com/Nadim147c/material/v2/palettes"
	"github.com/Nadim147c/rong/v4/internal/base16"
	"github.com/Nadim147c/rong/v4/internal/config"
	"github.com/Nadim147c/rong/v4/internal/config/enums"
	"github.com/Nadim147c/rong/v4/internal/material"
	"github.com/Nadim147c/rong/v4/internal/models"
	"github.com/Nadim147c/rong/v4/internal/templates"
	"github.com/spf13/cobra"

	_ "image/jpeg" // for jpeg encoding
	_ "image/png"  // for png encoding

	_ "golang.org/x/image/webp" // for webp encoding
)

// Command is the color command.
var Command = &cobra.Command{
	Use:   "color <color>",
	Short: "Generate colors from a color",
	Example: `
# Using color name
rong color hot_pink

# Using #RGB format
rong color '#F00'

# Using #RRGGBB format
rong color '#00FF00'

# Get generate colors as json
rong color green --dry-run --json | jq
  `,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		name := strings.ToLower(args[0])
		source, ok := Names[name]
		if !ok {
			src, err := color.ARGBFromHex(name)
			if err != nil {
				return err
			}
			source = src
		}

		slog.Info("Generating color", "from", source)

		primary := palettes.NewFromARGB(source)

		cfg := material.GetConfig()

		scheme := dynamic.NewDynamicScheme(source.ToHct(),
			cfg.Variant, cfg.Constrast, cfg.Dark,
			cfg.Platform, cfg.Version, primary,
			nil, nil, nil, nil, nil,
		)

		dcs := scheme.ToColorMap()

		colorMap := map[string]color.ARGB{}
		for key, value := range dcs {
			if value != nil {
				colorMap[key] = value.GetArgb(scheme)
			}
		}

		customs, err := material.GenerateCustomColors(colorMap["primary"])
		if err != nil {
			return err
		}

		// dynamic base16 generation is not possible with single source color
		config.Base16Method.SetValue(enums.Base16MethodStatic)
		based, err := base16.Generate(colorMap, nil)
		if err != nil {
			return err
		}

		output := models.NewOutput("", based, colorMap, customs)

		if config.JSON.Value() {
			err := json.NewEncoder(cmd.OutOrStdout()).Encode(output)
			if err != nil {
				slog.Error("Failed to encode output", "error", err)
			}
		}

		if config.SimpleJSON.Value() {
			err := models.WriteSimpleJSON(cmd.OutOrStdout(), output)
			if err != nil {
				slog.Error("Failed to encode output", "error", err)
			}
		}

		if config.DryRun.Value() {
			return nil
		}

		return templates.Execute(ctx, output)
	},
}
