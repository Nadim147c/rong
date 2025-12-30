package cache

import (
	"context"
	"log/slog"
	"runtime"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	Command.Flags().Int("frames", 5, "number of frames of video to process")
	Command.Flags().
		Duration("duration", 5*time.Second, "maxium number of duration to process")
	Command.Flags().
		Int("workers", runtime.NumCPU(), "number of concurrent workers")
}

// Command is cache command.
var Command = &cobra.Command{
	Use:   "cache <path...>",
	Short: "Generate color cache from images or videos",
	Example: `
# Cache from a video
rong cache path/to/video.mkv

# Cache from a image
rong cache path/to/image.webp

# Cache all png in a directory
rong cache path/to/*.png

# Recursively cache all image and video in a directory
rong cache path/to/directory
  `,
	Args: cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())

		states := make(chan state)
		model := newModel(cancel)
		p := tea.NewProgram(model)

		var wg sync.WaitGroup

		wg.Go(func() {
			_, err := p.Run()
			if err != nil {
				slog.Info("Tui program failed", "error", err)
			}
		})
		wg.Go(func() { cacheRec(ctx, args, states) })
		for state := range states {
			p.Send(state)
		}
		wg.Wait()
		return nil
	},
}
