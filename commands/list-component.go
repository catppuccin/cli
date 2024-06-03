package commands

import (
	"fmt"
	"io"
	"strings"

	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

const listHeight = 14

type RenderableItem interface {
	list.Item
	Render() string
}

type ListComponent[TItem RenderableItem] struct {
	list   list.Model
	popup  PopupModel
	config ListConfig[TItem]

	showingPopup bool
	selectedItem TItem
}

func (m ListComponent[TItem]) Init() tea.Cmd {
	return m.popup.Init()
}

func (m ListComponent[TItem]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		for _, k := range m.config.keys {
			if key.Matches(msg, k.Key) {
				return m, k.OnPress(&m)
			}
		}
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

		if m.showingPopup {
			m.popup, cmd = m.popup.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		}

	case error:
		panic(msg.Error())
	}

	m.popup, cmd = m.popup.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ListComponent[TItem]) CurrentItem() TItem {
	if i, ok := m.list.SelectedItem().(DelegatedItem[TItem]); ok {
		return i.item
	}

	panic("No selected item")
}

func (m ListComponent[TItem]) View() string {
	s := "\n"

	if m.showingPopup {
		s += lipgloss.JoinHorizontal(lipgloss.Top, m.list.View(), m.popup.View())
	} else {
		s += m.list.View()
	}

	return "\n" + s
}

func (l ListComponent[TItem]) Show() {
	if _, err := tea.NewProgram(l).Run(); err != nil {
		log.Error("Error running program:", err)
	}
}

type ListItemDelegate[TItem RenderableItem] struct {
	config ListConfig[TItem]
}

func (d ListItemDelegate[TItem]) Height() int                             { return 1 }
func (d ListItemDelegate[TItem]) Spacing() int                            { return 0 }
func (d ListItemDelegate[TItem]) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ListItemDelegate[TItem]) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(DelegatedItem[TItem])
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.item.Render())

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render(strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func (d ListItemDelegate[TItem]) ShortHelp() []key.Binding {
	keys := make([]key.Binding, 0, len(d.config.keys))
	for _, v := range d.config.keys {
		keys = append(keys, v.Key)
	}
	return keys
}

func (d ListItemDelegate[TItem]) FullHelp() [][]key.Binding {
	return nil
}

type DelegatedItem[TItem RenderableItem] struct{ item TItem }

func (d DelegatedItem[TItem]) FilterValue() string {
	return d.item.FilterValue()
}

func WrapItems[TBase any, TItem RenderableItem](items []TBase, wrapper func(TBase) TItem) []TItem {
	mapped := make([]TItem, 0, len(items))
	for _, v := range items {
		mapped = append(mapped, wrapper(v))
	}
	return mapped
}

var (
	itemStyle = lipgloss.
			NewStyle().
			PaddingLeft(2).
			Padding(0, 0, 1, 2).
			Foreground(lipgloss.Color(catppuccin.Mocha.Text().Hex))

	selectedItemStyle = lipgloss.
				NewStyle().
				Padding(0, 1, 0, 1).
				Foreground(lipgloss.Color(catppuccin.Mocha.Mauve().Hex)).
				Border(lipgloss.RoundedBorder(), false, false, true, true).
				BorderForeground(lipgloss.Color(catppuccin.Mocha.Mauve().Hex))
)

type ListBinding[TItem RenderableItem] struct {
	Key     key.Binding
	OnPress func(model *ListComponent[TItem]) tea.Cmd
}

type ListConfig[TItem RenderableItem] struct {
	items []TItem
	keys  []ListBinding[TItem]
	title string
}

func NewListComponent[TItem RenderableItem](config ListConfig[TItem]) ListComponent[TItem] {
	const defaultWidth = 20

	delegate := ListItemDelegate[TItem]{
		config: config,
	}

	mapped := make([]list.Item, 0, len(config.items))
	for _, v := range config.items {
		mapped = append(mapped, DelegatedItem[TItem]{item: v})
	}

	lst := list.New(mapped, delegate, defaultWidth, listHeight)
	lst.Title = config.title
	lst.SetShowStatusBar(true)
	lst.SetFilteringEnabled(true)
	lst.ShowFilter()

	return ListComponent[TItem]{
		list:   lst,
		popup:  NewPopup(),
		config: config,
	}
}
