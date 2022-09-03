package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SpinnerParent struct {
	spinner spinner.Model
}

func NewSpinnerParent() *SpinnerParent {
	spin := spinner.New()
	return &SpinnerParent{spinner: spin}
}

func (m SpinnerParent) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m SpinnerParent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m SpinnerParent) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, "Setting up repo...", m.spinner.View())
}
