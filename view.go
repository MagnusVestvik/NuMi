package main

import (
	"fmt"
)

const (
	MainView = iota
	SearchView
	ListInstalledPackagesView
	HelpView
)

func (m model) View() string {
	switch m.ViewState {
	case MainView:
		return m.MainView()
	case SearchView:
		return m.SearchView()
	case ListInstalledPackagesView:
		return m.ListInstalledPackagesView()
	case HelpView:
		return m.HelpView()
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
	for i, choice := range m.Choices {
		cursor := " " // no cursor
		if m.Cursor == i {
			cursor = ">" // cursor!
		}
		checked := " " // not selected
		if _, ok := m.Selected[i]; ok {
			checked = "x" // selected!
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	s += "\nPress q to quit.\n"
	return s
}
func (m model) SearchView() string {
	s := m.InputField.View() + "\n\n"
	if len(m.Choices) == 0 {
		s += "No results yet. Press Enter to search.\n"
	} else {
		for i, choice := range m.Choices {
			cursor := " " // no cursor
			if m.Cursor == i {
				cursor = ">" // cursor!
			}
			checked := " " // not selected
			if _, ok := m.Selected[i]; ok {
				checked = "x" // selected!
			}
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		}
	}
	return s
}

func (m model) ListInstalledPackagesView() string {
	return "list of installed packages lol"
}
