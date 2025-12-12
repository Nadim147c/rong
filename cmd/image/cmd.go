package image

import (
	"encoding/json"
	"fmt"
	"image"
	"log/slog"
	"os"

	"github.com/Nadim147c/rong/v4/internal/base16"
	"github.com/Nadim147c/rong/v4/internal/cache"
	"github.com/Nadim147c/rong/v4/internal/material"
	"github.com/Nadim147c/rong/v4/internal/models"
	"github.com/Nadim147c/rong/v4/internal/pathutil"
	"github.com/Nadim147c/rong/v4/internal/templates"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "image/jpeg" // for jpeg encoding
	_ "image/png"  // for png encoding

	_ "golang.org/x/image/webp" // for webp encoding
)

func init() {
	Command.Flags().AddFlagSet(material.Flags)
	Command.Flags().AddFlagSet(base16.Flags)
}

// Command is the image command
var Command = &cobra.Command{
	Use:   "image <image>",
	Short: "Generate colors from a image",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, _ []string) {
		viper.BindPFlags(cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		imagePath := args[0]

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		imagePath, err = pathutil.FindPath(cwd, imagePath)
		if err != nil {
			return fmt.Errorf("failed to find image path: %w", err)
		}

		slog.Info("Generating color", "from", imagePath)

		hash, err := cache.Hash(imagePath)
		if err != nil {
			return fmt.Errorf("failed to get xxh sum: %v", err)
		}

		quantized, err := cache.LoadCache(hash)
		if err != nil {
			if !os.IsNotExist(err) {
				slog.Error("Failed to load cache", "error", err)
			}

			file, err := os.Open(imagePath)
			if err != nil {
				return fmt.Errorf("failed to open image file: %w", err)
			}
			defer file.Close()

			img, _, err := image.Decode(file)
			if err != nil {
				return fmt.Errorf("failed to decode image: %w", err)
			}

			pixels := material.GetPixelsFromImage(img)
			quantized, err = material.Quantize(ctx, pixels)
			if err != nil {
				return err
			}

			if err := cache.SaveCache(hash, quantized); err != nil {
				slog.Warn("Failed to save colors to cache", "error", err)
			}
		}

		cfg, err := material.GetConfig()
		if err != nil {
			return err
		}

		colorMap, wu, err := material.GenerateFromQuantized(quantized, cfg)
		if err != nil {
			return fmt.Errorf("failed to generate colors: %w", err)
		}

		customs, err := material.GenerateCustomColors(colorMap["primary"])
		if err != nil {
			return err
		}

		based, err := base16.Generate(colorMap, wu)
		if err != nil {
			return err
		}

		output := models.NewOutput(imagePath, based, colorMap, customs)

		if viper.GetBool("json") {
			if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
				slog.Error("Failed to encode output", "error", err)
			}
		}

		if !viper.GetBool("dry-run") {
			if err := cache.SaveState(imagePath, hash, quantized); err != nil {
				slog.Warn("Failed to save colors to cache", "error", err)
			}

			return templates.Execute(output)
		}

		return nil
	},
}
