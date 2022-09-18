package ui

import (

	// "os"
	// "os/exec"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


type ExecModel struct {
	textInput textinput.Model
	err       error
}

func NewExecModel() ExecModel {
	ti := textinput.New()
	ti.Placeholder = "Helix"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 20

	return ExecModel{
		textInput: ti,
		err:       nil,
	}
}

func (m ExecModel) Init() tea.Cmd {
	return textinput.Blink
}


func (m ExecModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			// save value to global so it doesn't get lost
			// or you can wrap it as a tea.Msg and send it to the spinnerView to get handled
			ExecName = m.textInput.Value()
			return models[execView+1], models[execView+1].Init()
		}

	case errMsg:
		m.err = msg
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m ExecModel) View() string {
	// lipgloss will format the layout for you
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"What's the executable name?",
		m.textInput.View(),
		"(esc to quit)",
	)
}
