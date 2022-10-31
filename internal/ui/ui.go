package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

var p *tea.Program

func Run() {
	p := tea.NewProgram(NewInitialModel())
	if err := p.Start(); err != nil {
		panic(err)
	}
	if Cloned {
		fmt.Println("Finished!")
	}
}

// I put all the globals here :shrug:
var (
	RepoName   string
	ExecName   string
	Cloned     bool // Planning to use this to determine when to exit the spinner when the repo is cloned.
	LinuxLoc   string
	MacLoc     string
	WindowsLoc string // Bruh
)

type errMsg error
