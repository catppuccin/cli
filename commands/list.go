package commands

import (
	"fmt"

	"github.com/catppuccin/cli/query"
	"github.com/catppuccin/cli/shared"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"golang.org/x/exp/slices"
)

type flags struct {
	Maintained   bool `help:"List only the maintained ports" group:"X" xor:"X"`
	Unmaintained bool `help:"List only the unmaintained ports" group:"X" xor:"X"`
	Archived     bool `help:"List only the archived ports" group:"X" xor:"X"`
}

type ListCommand struct {
	Repositories repositorySubCommand `cmd:"List both ports and userstyles" default:"1" aliases:"repos"`
	Userstyles   userstylesSubCommand `cmd:"List only the userstyles"`
	Ports        portsSubCommand      `cmd:"List only the ports"`
}

type repositorySubCommand struct {
	flags
}

type userstylesSubCommand struct {
	flags
}

type portsSubCommand struct {
	flags
}

func (r *repositorySubCommand) Run(ctx *shared.Context) error {
	log.Fatal("Not implemented")
	return nil
}

func (r *userstylesSubCommand) Run(ctx *shared.Context) error {
	log.Fatal("Not implemented")
	return nil
}

type PortWithMeta struct {
	post query.Port
}

func (p PortWithMeta) FilterValue() string {
	return p.post.FilterValue()
}

func (p PortWithMeta) Render() string {
	return fmt.Sprintf("%s\n stars: 10 | issues: 3 | installs: 350", p.post.Render())
}

func (r *portsSubCommand) Run(ctx *shared.Context) error {
	ports := ctx.Datastore.Ports()
	title := "All ports"
	if r.Maintained {
		ports = query.FilterMaintained(ports)
		title = "Maintained ports"
	}

	if r.Unmaintained {
		ports = query.FilterUnmaintained(ports)
		title = "Unmaintained ports"
	}

	slices.SortFunc(ports, query.SortPortLexicographically)

	if ctx.Interactive {
		wrappedPorts := WrapItems(ports, func(port query.Port) PortWithMeta {
			return PortWithMeta{
				port,
			}
		})
		
		component := NewListComponent(ListConfig[PortWithMeta] {
			items: wrappedPorts,
			keys: []ListBinding[PortWithMeta]{
				{
					Key: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "choose")),
					OnPress: func(model *ListComponent[PortWithMeta]) tea.Cmd {
						model.showingPopup = true
						return nil
					},
				},
				{
					Key: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "jump")),
					OnPress: func(model *ListComponent[PortWithMeta]) tea.Cmd {
						model.showingPopup = !model.showingPopup
						return nil
					},
				},
			},
			title: title,
		})

		component.Show()

		return nil
	}

	fmt.Printf("There are %d ", len(ports))
	if r.Maintained {
		fmt.Println("maintained ports")
	}

	if r.Unmaintained {
		fmt.Println("unmaintained ports")
	}

	fmt.Println()

	for _, port := range ports {
		fmt.Printf("%s has %d maintainers\n", port.Name, len(port.CurrentMaintainers))
	}

	return nil
}
