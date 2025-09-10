package color

import (
	"encoding/json"
	"log/slog"
	"os"
	"strings"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
	"github.com/Nadim147c/material/palettes"
	"github.com/Nadim147c/rong/internal/base16"
	"github.com/Nadim147c/rong/internal/material"
	"github.com/Nadim147c/rong/internal/models"
	"github.com/Nadim147c/rong/templates"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "image/jpeg" // for jpeg encoding
	_ "image/png"  // for png encoding

	_ "golang.org/x/image/webp" // for webp encoding
)

func init() {
	Command.Flags().AddFlagSet(material.GeneratorFlags)
}

// Command is the image command
var Command = &cobra.Command{
	Use:   "color [flags] <image>",
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
	PreRun: func(cmd *cobra.Command, _ []string) {
		viper.BindPFlags(cmd.Flags())
	},
	RunE: func(_ *cobra.Command, args []string) error {
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

		cfg, err := material.GetConfig()
		if err != nil {
			return err
		}

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

		fg, bg := colorMap["on_background"], colorMap["background"]
		based := base16.GenerateRandom(fg, bg)

		output := models.NewOutput("", based, colorMap)

		if viper.GetBool("json") {
			json.NewEncoder(os.Stdout).Encode(output)
		}

		if !viper.GetBool("dry-run") {
			return templates.Execute(output)
		}

		return nil
	},
}
