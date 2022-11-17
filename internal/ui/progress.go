package ui

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/log"
	"github.com/catppuccin/cli/internal/utils"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-git/go-git/v5"
)

type progressMsg float64 // Progress

type finishMsg string // Finish

type GitProgress struct {
	Progress int
}

func NewProgressBar() ProgressWrapper {
	return ProgressWrapper{
		progress: progress.New(),
	}
}

// CloneRepo clones a repo into the specified location.
func CloneRepo(stagePath string, repo string) string {
	org := utils.GetEnv("ORG_OVERRIDE", "catppuccin")
	gitProgress := GitProgress{
		Progress: 0,
	}
	_, err := git.PlainClone(stagePath, false, &git.CloneOptions{
		URL:      fmt.Sprintf("https://github.com/%s/%s.git", org, repo),
		Progress: gitProgress,
	})
	if err != nil {
		// fmt.Println(err)
		log.WithError(err).Error("failed to clone repo")

	}
	return stagePath
}

func StartClone(repo string) tea.Cmd {
	return func() tea.Msg {
		CloneRepo(utils.GetTemplateDir(repo), "template")
		return nil
	}
}

// Regex to get the percentage in a subgroup
var re *regexp.Regexp = regexp.MustCompile(`Compressing objects: \s*(\d*)%`)
var reEnd = regexp.MustCompile(`Resolving deltas: \s*(\d*)%`)

// Write intercepts the content and updates the percentage
func (g GitProgress) Write(raw []byte) (n int, err error) {
	data := string(raw)
	matches := re.FindStringSubmatch(data)
	if len(matches) > 1 {
		deltaMatch := reEnd.FindStringSubmatch(data) // Please delete this. This does not work at all. :sadge:
		if !strings.Contains(deltaMatch[1], "Resolving deltas") {
			percentage, _ := strconv.Atoi(matches[1])
			g.Progress = percentage
			p.Send(progressMsg(percentage / 100)) // Send the progressMsg out
		} else {
			p.Send(finishMsg(deltaMatch[1])) // This should, in theory, send out finishMsg, but it does not.
		}
	}
	return len(raw), err
}

func finalPause() tea.Cmd {
	return tea.Tick(time.Millisecond*750, func(_ time.Time) tea.Msg {
		return nil
	})
}

type ProgressWrapper struct {
	progress progress.Model
}

func (m ProgressWrapper) Init() tea.Cmd {
	return StartClone(RepoName)
}

func (m ProgressWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit // Just quit
	case progressMsg:
		var cmds []tea.Cmd

		if msg >= 1.0 {
			cmds = append(cmds, tea.Sequentially(finalPause(), tea.Quit))
		}

		cmds = append(cmds, m.progress.SetPercent(float64(msg))) // Set the progress
		return m, tea.Batch(cmds...)
	case progress.FrameMsg:
		// Update bar
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	case finishMsg:
		return m, tea.Quit
	}
	return m, nil
}

func (m ProgressWrapper) View() string {
	return "\nDownloading...\n" +
		m.progress.View() + "\n\n" +
		"Press any key to quit"
}
