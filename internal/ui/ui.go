package ui

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

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

type tickMsg struct{}
type errMsg error

type confirm struct {
	choices  []string
	cursor   int
	selected int
}

type model struct {
	textInput textinput.Model
	confirm   confirm
	err       error
}

func InitialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Helix"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 20
	con := confirm{
		choices:  []string{"yes", "no"},
		selected: 0,
	}

	return model{
		textInput: ti,
		confirm:   con,
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
			//utils.CreateTemplate(m.textInput.Value())
			EnterVal := m.textInput.Value()
			fmt.Printf("Enterval %T\n", EnterVal)
			c := exec.Command("clear")
			c.Stdout = os.Stdout
			err := c.Run()
			if err != nil {
				return nil, nil
			}
			return m, tea.Quit
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
