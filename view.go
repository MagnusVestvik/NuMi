package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

const (
	MainView = iota
	SearchView
	ListInstalledPackagesView
	HelpView
	ProgressView
	padding  = 2
	maxWidth = 80
)

func (m model) View() string {
	switch m.viewState {
	case MainView:
		return m.MainView()
	case SearchView:
		return m.SearchView()
	case ListInstalledPackagesView:
		return m.ListInstalledPackagesView()
	case HelpView:
		return m.HelpView()
	case ProgressView:
		return m.ProgressView()
	default:
		return ""
	}
}

func (m model) HelpView() string {
	return "Help is coming soon!"
}

func (m model) MainView() string {
	s := ""
	s += "What should we buy at the market?\n\n"
	for i, choice := range m.choices {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	s += "\nPress q to quit.\n"
	return s
}
func (m model) SearchView() string {
	s := m.inputField.View() + "\n\n"
	if len(m.table.Rows()) == 0 {
		s += "No results yet. Press Enter to search.\n"
	} else {
		s += baseStyle.Render(m.table.View()) + "\n"
	}
	return s
}

func (m model) ProgressView() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle("Press any key to quit")
}

func (m model) ListInstalledPackagesView() string {
	return "list of installed packages lol"
}
