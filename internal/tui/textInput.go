package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ProjectNameInputModel struct {
	Input textinput.Model
}

func InitialProjectNameInputModel() ProjectNameInputModel {
	ti := textinput.New()
	ti.Placeholder = "new-project"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return ProjectNameInputModel{
		Input: ti,
	}
}

func (m ProjectNameInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ProjectNameInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			m.Input.Blur()
			return m, tea.Quit
		}
	}

	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

func (m ProjectNameInputModel) View() string {
	return fmt.Sprintf("Project name %s\n", m.Input.View())
}
