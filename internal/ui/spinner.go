package ui

import (
	"github.com/catppuccin/cli/internal/utils"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SpinnerParent struct {
	spinner spinner.Model
}

type spinnerMsg int

func GetRepoName() tea.Msg {
	utils.CloneTemplate(RepoName)
	utils.InitTemplate(RepoName, ExecName, LinuxLoc, MacLoc, WindowsLoc)
	return spinnerMsg(1)
}

func NewSpinnerParent() *SpinnerParent {
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return &SpinnerParent{spinner: spin}
}

func (m SpinnerParent) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, GetRepoName)
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
	case spinnerMsg:
		// Cloning completed, quit.
		Cloned = true
		return m, tea.Quit
	}
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m SpinnerParent) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, "Setting up repo...", m.spinner.View())
}
