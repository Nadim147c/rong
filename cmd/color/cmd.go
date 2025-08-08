package color

import (
	"encoding/json"
	"os"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
	"github.com/Nadim147c/material/palettes"
	"github.com/Nadim147c/rong/internal/base16"
	"github.com/Nadim147c/rong/internal/config"
	"github.com/Nadim147c/rong/internal/models"
	"github.com/Nadim147c/rong/templates"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "image/jpeg" // for jpeg encoding
	_ "image/png"  // for png encoding

	_ "golang.org/x/image/webp" // for webp encoding
)

func init() {
	Command.Flags().Bool("dark", false, "generate dark color palette")
	Command.Flags().String("variant", string(dynamic.TonalSpot), "variant to use (e.g., tonal_spot, vibrant, expressive)")
	Command.Flags().Float64("contrast", 0.0, "contrast adjustment (-1.0 to 1.0)")
	Command.Flags().String("platform", string(dynamic.Phone), "target platform (phone or watch)")
	Command.Flags().Int("version", int(dynamic.V2021), "version of the theme (2021 or 2025)")
	Command.Flags().Bool("dry-run", false, "generate colors without applying templates")
}

// Command is the image command
var Command = &cobra.Command{
	Use:   "color [flags] <image>",
	Short: "Generate colors from a color",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, _ []string) {
		viper.BindPFlags(cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		source, err := color.ARGBFromHex(args[0])
		if err != nil {
			return err
		}

		primary := palettes.NewFromARGB(source)

		config, err := config.GetGeneratorConfig()
		if err != nil {
			return err
		}
		scheme := dynamic.NewDynamicScheme(source.ToHct(),
			config.Variant, config.Constrast, config.Dark,
			config.Platform, config.Version, primary,
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

		if jsonFlag, _ := cmd.Flags().GetBool("json"); jsonFlag {
			json.NewEncoder(os.Stdout).Encode(output)
		}

		if dry, _ := cmd.Flags().GetBool("dry-run"); !dry {
			templates.Execute(output)
		}
		return nil
	},
}
