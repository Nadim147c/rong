package cache

import (
	"context"
	"fmt"
	"runtime"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func init() {
	Command.Flags().Int("frames", 5, "number of frames of video to process")
	Command.Flags().Int("workers", runtime.NumCPU(), "number of concurrent workers")
}

type processingDoneMsg struct{}

const (
	padding  = 2
	maxWidth = 80
)

type model struct {
	ctx             context.Context
	frames, workers int

	active []*job
	queue  queue

	progress progress.Model

	width     int
	height    int
	completed int
	total     int
	done      bool
}

// NewModel creates a new model with the given paths
func newModel(ctx context.Context, paths []string, frames, workers int) model {
	m := model{
		total:    len(paths),
		frames:   frames,
		workers:  workers,
		progress: progress.New(),
	}
	for path := range slices.Values(paths) {
		j := newJob(ctx, path, frames)
		if len(m.active) < workers {
			m.active = append(m.active, j)
		} else {
			m.queue.Enqueue(j)
		}
	}

	return m
}

func (m *model) selectJob(filename string) *job {
	for j := range slices.Values(m.active) {
		if j.filename == filename {
			return j
		}
	}
	panic("invalid job name")
}

// Init is Init
func (m model) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0, m.workers)
	for j := range slices.Values(m.active) {
		cmds = append(cmds, j.Init())
	}
	return tea.Batch(cmds...)
}

// Update updates Update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.progress.Width = min(msg.Width-padding*2-4, maxWidth)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}
	case processingDoneMsg:
		m.done = true
		return m, tea.Quit
	case jobDone:
		m.completed++
		j := m.selectJob(msg.Name)
		jobModel, _ := j.Update(msg)
		var cmds []tea.Cmd

		cmds = append(cmds, tea.Println(jobModel.View()))

		m.active = slices.DeleteFunc(m.active, func(e *job) bool {
			return e.filename == j.filename
		})

		progCmd := m.progress.SetPercent(float64(m.completed) / float64(m.total))
		cmds = append(cmds, progCmd)

		if len(m.active) == 0 && m.queue.Size() == 0 {
			cmds = append(cmds, func() tea.Msg { return processingDoneMsg{} })
			return m, tea.Sequence(cmds...)
		}

		if len(m.active) < m.workers {
			newJob, ok := m.queue.Dequeue()
			if ok {
				m.active = append(m.active, newJob)
				cmds = append(cmds, newJob.Init())
			}
		}

		return m, tea.Sequence(cmds...)
	case progress.FrameMsg:
		p, cmd := m.progress.Update(msg)
		m.progress = p.(progress.Model)
		return m, cmd
	case spinner.TickMsg:
		cmds := make([]tea.Cmd, 0, len(m.active))
		for i, j := range m.active {
			var cmd tea.Cmd
			m.active[i], cmd = j.Update(msg)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

// View views View
func (m model) View() string {
	if m.done {
		return fmt.Sprintf("Done! Processed %d/%d files\n", m.completed, m.total)
	}

	var buf strings.Builder
	for _, job := range m.active {
		fmt.Fprintln(&buf, job.View())
	}

	for _, job := range m.queue.NextN(m.height - m.workers - 3) {
		fmt.Fprintln(&buf, job.View())
	}

	fmt.Fprintf(&buf, "\n%s%s\n", strings.Repeat(" ", padding), m.progress.View())

	return buf.String()
}

// Command is cache command
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
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		frames, _ := cmd.Flags().GetInt("frames")
		workers, _ := cmd.Flags().GetInt("workers")

		// Collect all paths first
		paths, err := ScanPaths(ctx, args)
		if err != nil {
			return err
		}

		if len(paths) == 0 {
			fmt.Println("No files found to process")
			return nil
		}

		model := newModel(ctx, paths, frames, workers)
		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("error running program: %w", err)
		}

		return nil
	},
}
