package cache

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/adrg/xdg"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/charmtone"
)

const (
	padding  = 2
	maxWidth = 80
)

type model struct {
	exit context.CancelFunc
	done bool

	active    []job
	completed []job
	queued    int

	progress progress.Model
	spinner  spinner.Model

	width, height int
}

// NewModel creates a new model with the given paths.
func newModel(exit context.CancelFunc) *model {
	m := new(model)
	m.exit = exit
	m.progress = progress.New()
	m.spinner = spinner.New(spinner.WithSpinner(spinner.MiniDot))
	return m
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

var completedText = lipgloss.NewStyle().
	Foreground(lipgloss.Color(charmtone.Pickle.Hex())).
	Render("Completed")

var errorText = lipgloss.NewStyle().
	Foreground(lipgloss.Color(charmtone.Bengal.Hex())).
	Render("Failed")

func coloredTxt(txt string, c charmtone.Key) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(c.Hex())).
		Render(txt)
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
		p := "file://" + host + path
		return ansi.SetHyperlink(p) + pretty + ansi.ResetHyperlink()
	}

	return path
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.progress.Width = min(msg.Width-padding*2-4, maxWidth)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			m.done = true
			m.exit()
			return m, tea.Quit
		}
	case state:
		var cmds []tea.Cmd
		if len(msg.completed) > len(m.completed) {
			for j := range slices.Values(msg.completed[len(m.completed):]) {
				var cmd tea.Cmd
				if j.err != nil {
					cmd = tea.Printf("✓ %s %s\n  %v", errorText, prettyPath(j.filename), j.err)
				} else {
					cmd = tea.Printf("✓ %s %s", completedText, prettyPath(j.filename))
				}
				cmds = append(cmds, cmd)
			}
		}

		m.completed = msg.completed
		m.active = msg.active
		m.queued = msg.queued

		completed := len(m.completed)
		total := completed + len(m.active) + m.queued
		perc := float64(completed) / float64(total)
		cmds = append(cmds, m.progress.SetPercent(perc))

		if msg.done {
			m.done = true
			cmds = append(cmds, tea.Quit)
		}

		return m, tea.Sequence(cmds...)
	case progress.FrameMsg:
		p, cmd := m.progress.Update(msg)
		m.progress = p.(progress.Model)
		return m, cmd
	case spinner.TickMsg:
		s, cmd := m.spinner.Update(msg)
		m.spinner = s
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	if m.done {
		return fmt.Sprintf("Done! Processed %d files\n", len(m.completed))
	}

	s := m.spinner.View()
	var buf strings.Builder
	for j := range slices.Values(m.active) {
		fmt.Fprintf(&buf, "%s %s %s\n", s, j.status, prettyPath(j.filename))
	}

	indent := strings.Repeat(" ", padding)
	fmt.Fprintf(&buf, "\n%s%s\n", indent, m.progress.View())

	return buf.String()
}
