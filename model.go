package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	Choices      []string
	ViewState    int
	Cursor       int
	Selected     map[int]string
	SearchTerm   string
	isInstalling bool
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.ViewState {
	case MainView:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.Cursor > 0 {
					m.Cursor--
				}
			case "down", "j":
				if m.Cursor < len(m.Choices)-1 {
					m.Cursor++
				}
			case "enter", " ":
				_, ok := m.Selected[m.Cursor]
				if ok {
					delete(m.Selected, m.Cursor)
				} else {
					m.Selected[m.Cursor] = m.Choices[m.Cursor]
				}

				var cmd tea.Cmd
				switch m.Choices[m.Cursor] {
				case "List Packages In Project":
					m.ViewState = ListInstalledPackagesView
				case "Search Packages":
					m.ViewState = SearchView
				case "Help":
					m.ViewState = HelpView
				case "Quit":
					return m, tea.Quit
				}
				return m, cmd
			}
		}
	case SearchView:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "i", "enter":
				m.IsInstalling = true
			}

		}

	case ListInstalledPackagesView:
	}
	return m, nil
}
