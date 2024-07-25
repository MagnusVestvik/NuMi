package main

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	Table        table.Model
	InputField   textinput.Model
	Choices      []string
	ViewState    int
	Cursor       int
	Selected     map[int]string
	SearchTerm   string
	isInstalling bool
	ShowLogs     bool
	WindowSize   tea.WindowSizeMsg
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	logMu.Lock()
	logger.Printf("Received message: %T", msg)
	logMu.Unlock()
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case SearchResult:
		logMu.Lock()
		logger.Printf("Received SearchResult: %v", msg.PackageName)
		logMu.Unlock()
		m.Table = arrangeSearchResultTable(msg) // TODO: denne oppdaterer ikke table
		logMu.Lock()
		logger.Printf("updated table: %v", msg)
		logMu.Unlock()
		m.Cursor = 0
		return m, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+l":
			m.ShowLogs = !m.ShowLogs
			return m, nil
		}
	}
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

				switch m.Choices[m.Cursor] {
				case "List Packages In Project":
					m.ViewState = ListInstalledPackagesView
				case "Search Packages":
					m.ViewState = SearchView
					m.Choices = []string{}
				case "Help":
					m.ViewState = HelpView
				case "Quit":
					return m, tea.Quit
				}
				return m, nil
			}
		}
	case SearchView:
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
			case "enter":
				return m, SearchPackagesCmd(m.InputField.Value())
			}
			var cmd tea.Cmd
			m.InputField, cmd = m.InputField.Update(msg)
			return m, cmd
		}

	case ListInstalledPackagesView:
	}
	return m, nil
}
