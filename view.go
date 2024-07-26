package main

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var listBaseStyle = lipgloss.NewStyle().Margin(1, 2)

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
	return center(m, "Help is coming soon!")
}

func (m model) MainView() string {
	return center(m, listBaseStyle.Render(m.startOptionsList.View()))
}
func (m model) SearchView() string {
	s := m.inputField.View() + "\n\n"
	if len(m.packageSearchTable.Rows()) == 0 {
		s += "No results yet. Press Enter to search.\n"
	} else {
		s += baseStyle.Render(m.packageSearchTable.View()) + "\n"
	}
	return center(m, s)
}

func (m model) ProgressView() string {
	logMu.Lock()
	logger.Printf("ProgressView rendering! ")
	logMu.Unlock()

	pad := strings.Repeat(" ", padding)
	pad += "\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle("Press any key to quit")
	return center(m, pad)
}

func (m model) ListInstalledPackagesView() string {
	return center(m, "list of installed packages lol")
}
