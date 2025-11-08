package cache

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Nadim147c/rong/internal/cache"
	"github.com/Nadim147c/rong/internal/ffmpeg"
	"github.com/Nadim147c/rong/internal/material"
	"github.com/adrg/xdg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	lipgloss2 "github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
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

func coloredTxt(txt string, c charmtone.Key) string {
	return lipgloss2.NewStyle().Bold(true).Foreground(c).Render(txt)
}

func prettyPath(path string) string {
	host, err := os.Hostname()
	if err != nil {
		return path
	}

	ud := xdg.UserDirs

	pathMap := map[string]string{
		// XDG Base Dirs
		xdg.Home:       coloredTxt("HOME", charmtone.Violet),
		xdg.ConfigHome: coloredTxt("CONFIG", charmtone.Charple),
		xdg.CacheHome:  coloredTxt("CACHE", charmtone.Zest),
		xdg.DataHome:   coloredTxt("DATA", charmtone.Julep),
		xdg.StateHome:  coloredTxt("STATE", charmtone.Tuna),
		xdg.RuntimeDir: coloredTxt("RUNTIME", charmtone.Butter),

		// User Dirs
		ud.Desktop:     coloredTxt("DESKTOP", charmtone.Malibu),
		ud.Documents:   coloredTxt("DOCUMENTS", charmtone.Coral),
		ud.Download:    coloredTxt("DOWNLOADS", charmtone.Zest),
		ud.Music:       coloredTxt("MUSIC", charmtone.Uni),
		ud.Pictures:    coloredTxt("PICTURES", charmtone.Blush),
		ud.PublicShare: coloredTxt("PUBLIC", charmtone.Bok),
		ud.Templates:   coloredTxt("TEMPLATES", charmtone.Dolly),
		ud.Videos:      coloredTxt("VIDEOS", charmtone.Turtle),
	}

	var prefix string
	for userPath := range pathMap {
		if strings.HasPrefix(path, userPath) && len(userPath) > len(prefix) {
			prefix = userPath
		}
	}

	if short, ok := pathMap[prefix]; ok {
		pretty := strings.Replace(path, prefix, short, 1)
		return ansi.SetHyperlink(
			"file://"+host+path,
		) + pretty + ansi.ResetHyperlink()
	}

	return path
}

// This does breaks the tea.Model
func (j *job) View(spiner string) string {
	icon := spiner
	var status string
	var c charmtone.Key

	width := 18

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
		width = 9 // size of the word 'Completed'
	case jobFailed:
		icon = "✗"
		status = fmt.Sprintf("Failed: %v", j.Error)
		width = len(status) // error length
		c = charmtone.Charcoal
	}

	status = lipgloss.NewStyle().
		Width(width).
		Foreground(lipgloss.Color(c.Hex())).
		Bold(true).
		Render(status)

	return fmt.Sprintf("%s %s %s", icon, status, prettyPath(j.filename))
}
