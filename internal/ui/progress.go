package ui

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/catppuccin/cli/internal/utils"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-git/go-git/v5"
)

type progressMsg float64 // Progress 


type GitProgress struct {
	Progress int
}


func NewProgressBar() progressWrapper {
	return progressWrapper{
		progress: progress.New(),
	}
}

// CloneRepo clones a repo into the specified location.
func CloneRepo(stagePath string, repo string) string {
	org := utils.GetEnv("ORG_OVERRIDE", "catppuccin")
	progress := GitProgress{
		Progress: 0,
	}
	_, err := git.PlainClone(stagePath, false, &git.CloneOptions{
		URL: fmt.Sprintf("https://github.com/%s/%s.git", org, repo),
		Progress: progress,
	})
	if err != nil {
		fmt.Println(err)
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
var re *regexp.Regexp = regexp.MustCompile(`Compressing objects:\s*(\d*)%`)


// Write intercepts the content and updates the percentage
func (g GitProgress) Write (raw []byte) (n int, err error) {
	data := string(raw)
	matches := re.FindStringSubmatch(data)
	if len(matches) > 1 {
		percentage, _ := strconv.Atoi(matches[1])
		g.Progress = percentage
		p.Send(progressMsg(percentage/100)) // Send the progressMsg out
	}
	return len(raw), err
}

func finalPause() tea.Cmd {
	return tea.Tick(time.Millisecond*750, func(_ time.Time) tea.Msg {
		return nil
	})
}

type progressWrapper struct {
	progress progress.Model
}

func (m progressWrapper) Init() tea.Cmd {
	return StartClone(RepoName)
}

func (m progressWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit // Just quit
	case progressMsg:
		var cmds []tea.Cmd

		if msg >= 1.0 {
			cmds = append(cmds, finalPause())
		}

		cmds = append(cmds, m.progress.SetPercent(float64(msg))) // Set the progress
		return m, tea.Batch(cmds...)
	case progress.FrameMsg:
		// Update bar
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}
	return m, nil
}

func (m progressWrapper) View() string {
	return "\nDownloading...\n" +
		m.progress.View() + "\n\n" +
		"Press any key to quit"
}
