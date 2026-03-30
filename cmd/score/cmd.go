package score

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"os"

	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/material/v2/score"
	"github.com/Nadim147c/rong/v5/internal/cache"
	"github.com/Nadim147c/rong/v5/internal/config"
	"github.com/Nadim147c/rong/v5/internal/ffmpeg"
	"github.com/Nadim147c/rong/v5/internal/material"
	"github.com/Nadim147c/rong/v5/internal/pathutil"
	"github.com/spf13/cobra"
)

// Command is the video command.
var Command = &cobra.Command{
	Use:   "score <media>",
	Short: "Generate colors from a video",
	Example: `
# Generate from a video
rong video path/to/video.mkv

# Generate from a image
rong video path/to/image.webp

# Get generate colors as json
rong video path/to/image.mp4 --dry-run --json | jq
  `,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		videoPath := args[0]

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		videoPath, err = pathutil.FindPath(cwd, videoPath)
		if err != nil {
			return fmt.Errorf("failed to find image path: %w", err)
		}

		slog.Info("Generating color", "from", videoPath)

		hash, err := cache.Hash(videoPath)
		if err != nil {
			return fmt.Errorf("failed to get xxh sum: %w", err)
		}

		quantized, err := cache.LoadCache(hash)
		if err != nil {
			if !os.IsNotExist(err) {
				slog.Error("Failed to load cache", "error", err)
			}

			frames := config.FFmpegFrames.Value()
			duration := config.FFmpegDuration.Value().Seconds()
			pixels, err := ffmpeg.GetPixels(ctx, videoPath, frames, duration)
			if err != nil {
				return fmt.Errorf("failed to get pixels from media: %w", err)
			}
			quantized, err = material.Quantize(ctx, pixels)
			if err != nil {
				return err
			}

			if err := cache.SaveCache(hash, quantized); err != nil {
				slog.Warn("Failed to save colors to cache", "error", err)
			}
		}

		slog.Info("Generating colors from source")

		colors := score.Score(quantized.Celebi, score.WithFilter(), score.WithLimit(5))

		for i, color := range colors {
			lab := color.ToOkLab()
			lab.L = 50
			colors[i] = lab.ToARGB()
		}

		if config.MergeThreshold.Value() != 0 {
			colors = mergeCloseColors(colors, config.MergeThreshold.Value())
		}

		return json.NewEncoder(cmd.OutOrStdout()).Encode(colors)
	},
}

func distance(a, b color.OkLab) float64 {
	dL, dA, dB := a.L-b.L, a.A-b.A, a.B-b.B
	return math.Cbrt(dL*dL + dA*dA + dB*dB)
}

func mergeCloseColors(colors []color.ARGB, threshold float64) []color.ARGB {
	var merged []color.ARGB
	for _, c := range colors {
		cLab := c.ToOkLab()
		isUnique := true

		for _, m := range merged {
			mLab := m.ToOkLab()
			dist := distance(cLab, mLab)
			if dist < threshold {
				isUnique = false
				break
			}
		}

		if isUnique {
			merged = append(merged, c)
		}
	}
	return merged
}
