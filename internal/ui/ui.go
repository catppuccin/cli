package ui

import (
	"fmt"

	"github.com/catppuccin/cli/internal/utils"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)


type KeyMap struct {
  Up key.Binding
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
type tickMsg struct{}
type errMsg error

type model struct {
	textInput textinput.Model
	err       error
}

func InitialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Helix"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 20

	return model{
		textInput: ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			utils.CreateTemplate(m.textInput.Value())
		}

	// Handle errors
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd

}

func (m model) View() string {
	return fmt.Sprintf(
		"What's the project name?\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}
