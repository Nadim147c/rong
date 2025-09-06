package material

import (
	"image"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
	"github.com/spf13/pflag"
)

// GeneratorConfig is configuration used to generate colors
type GeneratorConfig struct {
	Variant   dynamic.Variant  `config:"variant" check:"enum='monochrome,neutral,tonal_spot,vibrant,expressive,fidelity,content,rainbow,fruit_salad'"`
	Platform  dynamic.Platform `config:"platform" check:"enum='phone,watch'"`
	Version   dynamic.Version  `config:"version" check:"enum='2021,2025'"`
	Dark      bool             `config:"dark"`
	Constrast float64          `config:"Constrast" check:"min=-1,max=1"`
}

// GeneratorFlags are the flags used for generating colors
var GeneratorFlags = pflag.NewFlagSet("generate", pflag.ContinueOnError)

func init() {
	GeneratorFlags.Bool("dark", false, "generate dark color palette")
	GeneratorFlags.Bool("dry-run", false, "generate colors without applying templates")
	GeneratorFlags.Bool("json", false, "print generated colors as json")
	GeneratorFlags.Float64("contrast", 0.0, "contrast adjustment (-1.0 to 1.0)")
	GeneratorFlags.Int("version", int(dynamic.V2021), "version of the theme (2021 or 2025)")
	GeneratorFlags.String("platform", string(dynamic.Phone), "target platform (phone or watch)")
	GeneratorFlags.String("variant", string(dynamic.TonalSpot), "variant to use (e.g., tonal_spot, vibrant, expressive)")
}

// GetPixelsFromImage returns pixels from image.Imaget interface
func GetPixelsFromImage(img image.Image) []color.ARGB {
	bounds := img.Bounds()
	pixels := make([]color.ARGB, 0, bounds.Dx()*bounds.Dy())

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			argb := color.ARGBFromInterface(c)
			pixels = append(pixels, argb)
		}
	}

	return pixels
}

// GenerateFromImage colors from an image.Image
func GenerateFromImage(
	img image.Image,
	cfg GeneratorConfig,
) (Colors, []color.ARGB, error) {
	pixels := GetPixelsFromImage(img)
	return GenerateFromPixels(pixels, cfg)
}
