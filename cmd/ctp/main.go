package main

import (
	"github.com/caarlos0/log"
	"github.com/catppuccin/cli/internal/cli"
	"github.com/charmbracelet/lipgloss"
)

func init() {
	log.Styles = [...]lipgloss.Style{
		log.DebugLevel: lipgloss.NewStyle().Foreground(lipgloss.Color("#a6adc8")).Bold(true), // Subtext0
		log.InfoLevel:  lipgloss.NewStyle().Foreground(lipgloss.Color("#89b4fa")).Bold(true), // Blue
		log.WarnLevel:  lipgloss.NewStyle().Foreground(lipgloss.Color("#f9e2af")).Bold(true), // Yellow
		log.ErrorLevel: lipgloss.NewStyle().Foreground(lipgloss.Color("#f38ba8")).Bold(true), // Red
		log.FatalLevel: lipgloss.NewStyle().Foreground(lipgloss.Color("#f38ba8")).Bold(true), // Red
	}
}

func main() {
	cli.Execute()
}
