package cache

import (
	"context"
	"fmt"
	"strings"

	"github.com/Nadim147c/rong/internal/cache"
	"github.com/Nadim147c/rong/internal/ffmpeg"
	"github.com/Nadim147c/rong/internal/material"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/charmtone"

	"github.com/gabriel-vasile/mimetype"
)

// jobState represents the state of an individual file processing job
type jobState int

const (
	jobWaiting jobState = iota
	jobExtracting
	jobQuantizing
	jobSaving
	jobGeneratingPreview
	jobCompleted
	jobFailed
)

type jobDone struct {
	Name string
	Err  error
}

// job represents an individual file processing job
type job struct {
	ctx      context.Context
	filename string
	frames   int

	State     jobState
	Error     error
	Completed bool
}

// newJob creates a new job with a spinner
func newJob(ctx context.Context, filename string, frames int) *job {
	return &job{
		ctx:      ctx,
		filename: filename,
		frames:   frames,
		State:    jobWaiting,
	}
}

func (j *job) processFile() error {
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
	j.State = jobExtracting
	pixels, err := ffmpeg.GetPixels(j.ctx, j.filename, j.frames)
	if err != nil {
		return err
	}

	// Update state for quantization
	j.State = jobQuantizing
	quantized, err := material.Quantize(j.ctx, pixels)
	if err != nil {
		return err
	}

	// Update state for saving
	j.State = jobSaving
	if err := cache.SaveCache(hash, quantized); err != nil {
		return err
	}

	if isVideo {
		// Update state for preview generation
		j.State = jobGeneratingPreview
		if _, err := cache.GetPreview(j.filename, hash); err != nil {
			return err
		}
	}

	return nil
}

// Init is Init
func (j *job) Init() tea.Cmd {
	return func() tea.Msg {
		return jobDone{Name: j.filename, Err: j.processFile()}
	}
}

func (j *job) Update(msg tea.Msg) {
	if done, ok := msg.(jobDone); ok {
		j.Completed = true
		j.State = jobCompleted
		if done.Err != nil {
			j.State = jobFailed
			j.Error = done.Err
		}
	}
}

// This does breaks the tea.Model
func (j *job) View(spiner string) string {
	icon := spiner
	var status string
	var c charmtone.Key

	switch j.State {
	case jobWaiting:
		icon = "○"
		status = "Queued"
		c = charmtone.Squid
	case jobExtracting:
		status = "Extracting frames"
		c = charmtone.Zinc
	case jobQuantizing:
		status = "Quantizing colors"
		c = charmtone.Cherry
	case jobSaving:
		status = "Saving cache"
		c = charmtone.Grape
	case jobGeneratingPreview:
		status = "Generating preview"
		c = charmtone.Yam
	case jobCompleted:
		icon = "✓"
		status = "Completed"
	case jobFailed:
		icon = "✗"
		status = fmt.Sprintf("Failed: %v", j.Error)
		c = charmtone.Charcoal
	}

	status = lipgloss.NewStyle().
		Foreground(lipgloss.Color(c.Hex())).
		Bold(true).
		Render(status)

	return fmt.Sprintf("%s %s %s", icon, status, j.filename)
}
