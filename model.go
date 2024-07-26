package main

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: create more models
type model struct {
	table        table.Model
	inputField   textinput.Model
	choices      []string
	viewState    int
	cursor       int
	selected     map[int]string
	searchTerm   string
	isInstalling bool
	showLogs     bool
	windowSize   tea.WindowSizeMsg
	progress     progress.Model
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	logMu.Lock()
	logger.Printf("Received message: %T", msg)
	logMu.Unlock()
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case SearchResult:
		m.viewState = SearchView
		logMu.Lock()
		logger.Printf("Received SearchResult: %v", msg.PackageName)
		logMu.Unlock()
		m.table = arrangeSearchResultTable(msg) // TODO: denne oppdaterer ikke table
		logMu.Lock()
		logger.Printf("updated table: %v", msg)
		logMu.Unlock()
		m.cursor = 0
		return m, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+l":
			m.showLogs = !m.showLogs
			return m, nil
		}
	}
	switch m.viewState {
	case MainView:
		switch msg := msg.(type) {

		case progress.FrameMsg:
			progressModel, cmd := m.progress.Update(msg)
			m.progress = progressModel.(progress.Model)
			return m, cmd

		case tickMsg:
			if m.progress.Percent() == 1.0 {
				m.viewState = SearchView
				return m, nil
			}
			// Note that you can also use progress.Model.SetPercent to set the
			// percentage value explicitly, too.
			cmd := m.progress.IncrPercent(0.25)
			return m, tea.Batch(tickCmd(), cmd)
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			case "enter", " ":
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = m.choices[m.cursor]
				}

				switch m.choices[m.cursor] {
				case "List Packages In Project":
					m.viewState = ListInstalledPackagesView
				case "Search Packages":
					m.viewState = SearchView
					m.choices = []string{}
				case "Help":
					m.viewState = HelpView
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
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			case "enter":
				m.viewState = ProgressView
				tickIncrCmd := m.progress.IncrPercent(0.25)
				return m, tea.Batch(tickCmd(), tickIncrCmd, SearchPackagesCmd(m.inputField.Value()))
			}
			var cmd tea.Cmd
			m.inputField, cmd = m.inputField.Update(msg)
			return m, cmd
		}
	case ProgressView:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			}
		}

	case ListInstalledPackagesView:
	}
	return m, nil
}
