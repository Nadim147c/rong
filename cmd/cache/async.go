package cache

import (
	"context"
	"slices"
	"strings"
	"sync"

	"github.com/Nadim147c/rong/v4/internal/cache"
	"github.com/Nadim147c/rong/v4/internal/ffmpeg"
	"github.com/Nadim147c/rong/v4/internal/material"
	"github.com/gabriel-vasile/mimetype"
	"github.com/spf13/viper"
)

type job struct {
	filename string
	status   string
	err      error
}

func (j *job) process(
	ctx context.Context,
	update func(),
	frames int,
	duration float64,
) error {
	hash, err := cache.Hash(j.filename)
	if err != nil {
		return err
	}

	mtype, err := mimetype.DetectFile(j.filename)
	if err != nil {
		return err
	}
	isVideo := strings.HasPrefix(mtype.String(), "video")

	if cache.IsCached(hash, isVideo) {
		return nil
	}

	// Update state for extraction
	j.status = "Extracting pixels"
	update()
	pixels, err := ffmpeg.GetPixels(ctx, j.filename, frames, duration)
	if err != nil {
		return err
	}

	// Update state for quantization
	j.status = "Quantizing colors"
	update()
	quantized, err := material.Quantize(ctx, pixels)
	if err != nil {
		return err
	}

	// Update state for saving
	j.status = "Saving cache"
	update()
	if err := cache.SaveCache(hash, quantized); err != nil {
		return err
	}

	if isVideo {
		j.status = "Creating preview"
		update()
		if _, err := cache.GetPreview(j.filename, hash); err != nil {
			return err
		}
	}

	return nil
}

// This does breaks the tea.Model
type state struct {
	done      bool
	queued    int
	completed []job
	active    []job
}

func copyJobs(src []*job) []job {
	dst := make([]job, len(src))
	for i, v := range src {
		dst[i] = *v
	}
	return dst
}

func cacheRec(ctx context.Context, inputs []string, ch chan<- state) {
	defer close(ch)

	paths := make(chan string, 100)
	go func() {
		find(ctx, inputs, paths)
		close(paths)
	}()

	workers := viper.GetInt("workers")
	if workers == 0 {
		workers = 4
	}
	frames := viper.GetInt("frames")
	if frames == 0 {
		frames = 4
	}
	duration := viper.GetDuration("duration").Seconds()
	if duration == 0 {
		duration = 5
	}

	var mu sync.Mutex
	var wg sync.WaitGroup

	var completed []*job
	var active []*job

	update := func() {
		mu.Lock()
		defer mu.Unlock()
		var state state
		state.active = copyJobs(active)
		state.completed = copyJobs(completed)
		state.queued = len(paths)
		ch <- state
	}

	for range workers {
		wg.Go(func() {
			for path := range paths {
				select {
				case <-ctx.Done():
					return
				default:
				}
				j := &job{filename: path}

				mu.Lock()
				active = append(active, j)
				mu.Unlock()
				update()

				j.err = j.process(ctx, update, frames, duration)

				mu.Lock()
				active = slices.DeleteFunc(active, func(x *job) bool {
					return x.filename == j.filename
				})
				completed = append(completed, j)
				mu.Unlock()
				update()
			}
		})
	}

	wg.Wait()

	var state state
	state.done = true
	state.active = copyJobs(active)
	state.completed = copyJobs(completed)
	state.queued = len(paths)
	ch <- state
}
