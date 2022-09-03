package ui

import (
// "github.com/muesli/reflow/indent"
  tea "github.com/charmbracelet/bubbletea"
)

type spinnerMsg string

var repoName string

func makeSpinner(val string) tea.Cmd {
	// Returns a type that will have the program update to a spinner
	return func() tea.Msg {
	  return spinnerMsg(val)
	}
}

type ui struct {
  Current  tea.Model
}

func InitialUi() ui {
  return ui{
    Current: InitialModel(),
  }
}

func (m ui) Init() tea.Cmd {
  return m.Current.Init()
}

func (m ui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmd tea.Cmd
  commands := []tea.Cmd{}

  switch msg := msg.(type) {
  case spinnerMsg:
    m.Current = InitialModelSpinner()
    repoName = string(msg)
    commands = append(commands, m.Current.Init)
  }
  m.Current, cmd =  m.Current.Update(msg)
  commands = append(commands, cmd)
  return m, tea.Batch(commands...)
}

func (m ui) View() string {
  return m.Current.View()
}
