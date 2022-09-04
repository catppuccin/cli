package ui

// I'll work on this please but I need help with the io.writer function.
import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"io"
	"time"
)

const (
	padding  = 2
	maxWidth = 80
)

type ProgressParent struct {
	progress progress.Model
}

type tickMsg time.Time

func NewProgressParent() *ProgressParent {
	prog := progress.New()
	prog.Width = maxWidth
	return &ProgressParent{progress: prog}
}

func (m ProgressParent) Init() tea.Cmd {
	return tickCmd()
}
func (m ProgressParent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	}
	return m, nil // Placeholder code for now.
}

func (m ProgressParent) Write(io io.Writer) tea.Cmd {
	var cmd tea.Cmd // Placeholder for now while I figure out this part.
	return cmd
}
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m ProgressParent) View() string {
	// Implementation will come soon. Stay tuned.
	return "" // Placeholder code for now.
}
