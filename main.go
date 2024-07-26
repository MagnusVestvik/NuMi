package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

var (
	initialChoices = []list.Item{
		item{title: "Search Packages", desc: "Search the nuget library for packages"},
		item{title: "List Packages", desc: "List all installed packages current project"},
		item{title: "Help", desc: "Get help on how to use the application"},
	}

	logger *log.Logger
	logMu  sync.Mutex
	logBuf strings.Builder
)

func initialModel() model {
	// Create a log file
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	logger = log.New(io.MultiWriter(logFile, &logBuf), "", log.Ltime|log.Lshortfile)

	startViewList := list.New(initialChoices, list.NewDefaultDelegate(), 0, 0)
	startViewList.Title = "Main Menu"

	ti := textinput.New()
	ti.Placeholder = "Search for packages"
	ti.Focus()
	ti.Width = 20

	return model{
		inputField:       ti,
		startOptionsList: startViewList,
		viewState:        MainView,
		progress:         progress.New(progress.WithDefaultGradient()),
		tableIsSelected:  false,
		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]string),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
