package commands

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PopupModel struct {
	inner spinner.Model
}

func NewPopup() PopupModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return PopupModel{
		inner: s,
	}
}

func (m PopupModel) Init() tea.Cmd {
    return m.inner.Tick
}

func (m PopupModel) Update(msg tea.Msg) (PopupModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.inner, cmd = m.inner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m PopupModel) View() string {
	return m.inner.View()
}


