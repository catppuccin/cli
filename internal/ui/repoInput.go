package ui

import (
	"strings"

	// "os"
	// "os/exec"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// you don't need these keymaps but they can be helpful for generating the help
// menu for you
type KeyMap struct {
	Up   key.Binding
	Down key.Binding
}

var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),        // actual keybindings
		key.WithHelp("↑/k", "move up"), // corresponding help text
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "move down"),
	),
}

type InitialModel struct {
	textInput textinput.Model
	err       error
}

func NewInitialModel() InitialModel {
	ti := textinput.New()
	ti.Placeholder = "Helix"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 20

	return InitialModel{
		textInput: ti,
		err:       nil,
	}
}

func (m InitialModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m InitialModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			// save value to global so it doesn't get lost
			// or you can wrap it as a tea.Msg and send it to the spinnerView to get handled
			RepoName = strings.TrimSpace(m.textInput.Value()) // QOL addition to remove spaces so that the directory formed is `helix` and
			// not `helix\ /`
			if RepoName == "" {
				return m, nil
			}
			m := NewExecModel(strings.ToLower(RepoName))
			return m, m.Init()
		}

	case errMsg:
		m.err = msg
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m InitialModel) View() string {
	// lipgloss will format the layout for you
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"What's the project name?",
		m.textInput.View(),
		"(esc to quit)",
	)
}
