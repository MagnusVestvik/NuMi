package main

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: create more models
type model struct {
	table           table.Model
	tableIsSelected bool
	inputField      textinput.Model
	choices         []string
	viewState       int
	cursor          int
	selected        map[int]string
	searchTerm      string
	isInstalling    bool
	showLogs        bool
	windowSize      tea.WindowSizeMsg
	progress        progress.Model
	width           int
	height          int
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	logMu.Lock()
	logger.Printf("Received message: %T", msg)
	logMu.Unlock()
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case SearchResult:
		m.viewState = SearchView
		m.table = arrangeSearchResultTable(msg, m.width) // TODO: denne oppdaterer ikke table
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

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		logMu.Lock()
		logger.Printf("TickMsg received: %v", msg)
		logMu.Unlock()
		if m.progress.Percent() == 1.0 {
			m.viewState = SearchView
			return m, m.progress.SetPercent(0)
		}
		// Note that you can also use progress.Model.SetPercent to set the
		// percentage value explicitly, too.
		cmd := m.progress.IncrPercent(0.25)
		logMu.Lock()
		logger.Printf("Progress incremented: %v", cmd)
		logMu.Unlock()
		return m, tea.Batch(tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		logMu.Lock()
		logger.Printf("FrameMsg received: %v", msg)
		logMu.Unlock()
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	}
	switch m.viewState {
	case MainView:
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
			case "tab":
				logMu.Lock()
				logger.Printf("Tab key was pressed and resulted in tableIsSelected: %v", m.tableIsSelected)
				logMu.Unlock()

				m.tableIsSelected = !m.tableIsSelected
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.table.Rows())-1 {
					m.cursor++
				}
			case "enter":
				m.viewState = ProgressView
				tickIncrCmd := m.progress.IncrPercent(0.25)
				return m, tea.Batch(tickCmd(), tickIncrCmd, SearchPackagesCmd(m.inputField.Value()))
			}
			var cmd tea.Cmd
			if !m.tableIsSelected {
				m.inputField, cmd = m.inputField.Update(msg)
			} else {
				m.table.SetCursor(m.cursor)
				logMu.Lock()
				logger.Printf("Table moved cursor to position: %v", m.table.Cursor())
				logMu.Unlock()
			}
			return m, cmd
		}
	case ProgressView:
		logMu.Lock()
		logger.Printf("Recived ProgressView: %v", msg)
		logMu.Unlock()

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
