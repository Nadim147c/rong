package cache

import (
	"log/slog"
	"runtime"
	"sync"

	"github.com/Nadim147c/rong/internal/cache"
	"github.com/Nadim147c/rong/internal/ffmpeg"
	"github.com/Nadim147c/rong/internal/material"
	"github.com/spf13/cobra"
)

func init() {
	Command.Flags().Int("frames", 5, "number of frames of vidoe to process")
}

// Command is cache command
var Command = &cobra.Command{
	Use:   "cache <path...>",
	Short: "Generate color cache from a image/video",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		paths := make(chan string, 10)

		go ScanPaths(args, paths)

		var wg sync.WaitGroup

		lock := make(chan struct{}, runtime.NumCPU())
		for path := range paths {
			lock <- struct{}{}
			wg.Go(func() {
				defer func() { <-lock }()

				if cache.IsCached(path) {
					slog.Info("Skipping", "path", path, "reason", "already cached")
					return
				}

				frames, _ := cmd.Flags().GetInt("frames")
				pixels, err := ffmpeg.GetPixels(path, frames)
				if err != nil {
					slog.Error("Failed to get pixels from media", "path", path, "error", err)
					return
				}

				quantized := material.Quantize(pixels)
				if err := cache.SaveCache(path, quantized); err != nil {
					slog.Error("Failed to save cache", "path", path, "error", err)
					return
				}

				slog.Info("Successfully cached media", "path", path)
			})
		}
		close(lock)

		wg.Wait()
	},
}
