package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//type errMsg error

type modelSpinner struct {
	spinner  spinner.Model
	quitting bool
	err      error
}

func InitialModelSpinner() modelSpinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return modelSpinner{spinner: s}
}
func (j modelSpinner) Init() tea.Cmd {
	return j.spinner.Tick
}

func (j modelSpinner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			j.quitting = true
			return j, tea.Quit
		default:
			return j, nil
		}

	case errMsg:
		j.err = msg
		return j, nil

	default:
		var cmd tea.Cmd
		j.spinner, cmd = j.spinner.Update(msg)
		return j, cmd
	}

}
func (j modelSpinner) View() string {
	if j.err != nil {
		return j.err.Error()
	}
	str := fmt.Sprintf("\n\n   %s Loading forever...press q to quit\n\n", j.spinner.View())
	//fmt.Printf()
	if j.quitting {
		return str + "\n"
	}
	return str
}
