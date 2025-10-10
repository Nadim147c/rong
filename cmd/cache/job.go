package cache

import (
	"context"
	"fmt"
	"image/color"
	"strings"
	"sync"
	"time"

	"github.com/Nadim147c/rong/internal/cache"
	"github.com/Nadim147c/rong/internal/ffmpeg"
	"github.com/Nadim147c/rong/internal/material"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	ansi "github.com/charmbracelet/lipgloss/v2"

	"github.com/gabriel-vasile/mimetype"
)

// jobState represents the state of an individual file processing job
type jobState int

const (
	jobWaiting jobState = iota
	jobProcessing
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
	spinner  spinner.Model
	frames   int

	State     jobState
	Error     error
	Completed bool

	mu sync.Mutex
}

// newJob creates a new job with a spinner
func newJob(ctx context.Context, filename string, frames int) *job {
	s := spinner.New(spinner.WithSpinner(spinner.MiniDot))
	s.Spinner.FPS = time.Second / 30
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	return &job{
		ctx:      ctx,
		filename: filename,
		frames:   frames,
		State:    jobWaiting,
		spinner:  s,
	}
}

func (j *job) setState(s jobState) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.State = s
}

func (j *job) processFile() error {
	ctx := context.Background()

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
	j.setState(jobExtracting)
	pixels, err := ffmpeg.GetPixels(ctx, j.filename, j.frames)
	if err != nil {
		return err
	}

	// Update state for quantization
	j.setState(jobQuantizing)
	quantized, err := material.Quantize(ctx, pixels)
	if err != nil {
		return err
	}

	// Update state for saving
	j.setState(jobSaving)
	if err := cache.SaveCache(hash, quantized); err != nil {
		return err
	}

	if isVideo {
		// Update state for preview generation
		j.setState(jobGeneratingPreview)
		if _, err := cache.GetPreview(j.filename, hash); err != nil {
			return err
		}
	}

	return nil
}

// Init is Init
func (j *job) Init() tea.Cmd {
	return tea.Batch(j.spinner.Tick, func() tea.Msg {
		return jobDone{
			Name: j.filename,
			Err:  j.processFile(),
		}
	})
}

func (j *job) Update(msg tea.Msg) (*job, tea.Cmd) {
	switch msg := msg.(type) {
	case jobDone:
		j.Completed = true
		j.State = jobCompleted
		if msg.Err != nil {
			j.State = jobFailed
			j.Error = msg.Err
		}
		return j, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		j.spinner, cmd = j.spinner.Update(msg)
		return j, cmd
	}
	return j, nil
}

func (j *job) View() string {
	var icon, status string
	var c color.Color

	switch j.State {
	case jobWaiting:
		icon = "○"
		status = "Queued"
		c = ansi.BrightBlack
	case jobProcessing:
		icon = j.spinner.View()
		status = "Processing"
		c = ansi.Blue
	case jobExtracting:
		icon = j.spinner.View()
		status = "Extracting frames"
		c = ansi.BrightCyan
	case jobQuantizing:
		icon = j.spinner.View()
		status = "Quantizing colors"
		c = ansi.Blue
	case jobSaving:
		icon = j.spinner.View()
		status = "Saving cache"
		c = ansi.BrightGreen
	case jobGeneratingPreview:
		icon = j.spinner.View()
		status = "Generating preview"
		c = ansi.Magenta
	case jobCompleted:
		icon = "✓"
		status = "Completed"
		c = ansi.Green
	case jobFailed:
		icon = "✗"
		status = fmt.Sprintf("Failed: %v", j.Error)
		c = ansi.Red
	}

	status = ansi.NewStyle().Foreground(c).Bold(true).Render(status)

	return fmt.Sprintf("%s %s %s", icon, status, j.filename)
}
