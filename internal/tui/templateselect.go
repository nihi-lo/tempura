package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle      = lipgloss.NewStyle()
	quitTextStyle = lipgloss.NewStyle().Margin(0, 0, 1, 0)
)

type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type TemplateSelectModel struct {
	list     list.Model
	quitting bool
	Choice   string
}

func InitialTemplateSelectModel() TemplateSelectModel {
	items := []list.Item{
		item{title: "vite-react-tw3-ts", desc: "PrettierとTailwind CSS v3が導入済みのReactプロジェクト"},
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#f0e68c")).
		BorderForeground(lipgloss.Color("#f0e68c"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("#f0e68c")).
		BorderForeground(lipgloss.Color("#f0e68c"))

	m := TemplateSelectModel{
		list:   list.New(items, delegate, 0, 0),
		Choice: "",
	}
	m.list.Title = "Choose a template"
	m.list.SetShowPagination(false)
	m.list.Styles.Title = m.list.Styles.Title.
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.NoColor{}).
		Padding(0)

	return m
}

func (m TemplateSelectModel) Init() tea.Cmd {
	return nil
}

func (m TemplateSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.Choice = ""
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.Choice = string(i.title)
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m TemplateSelectModel) View() string {
	if m.Choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("Template > %s", m.Choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Because no template was selected, the project will not be created.")
	}
	return m.list.View()
}
