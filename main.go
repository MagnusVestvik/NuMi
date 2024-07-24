package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	initialChoices = []string{"List Packages In Project", "Search Packages", "Help", "Quit"}
	packages       = []string{}
)

type NugetOperations interface {
	SearchPackage() ([]string, error)
	InstallPackage() (bool, error)
	UninstallPackage() (bool, error)
	ListInstalledPackages() ([]string, error)
}

type ListNugetPackages struct {
	Packages []string
}

type NugetSearchResultMessage struct {
	Result []string
	Error  error
}

type InstallNugetPackagesMessage struct {
	Package string
	Error   error
}

type UninstallNugetPackageMessage struct {
	Package []string
	Error   error
}

func (s NugetSearchResultMessage) ListPackages() ([]string, error) {
	return s.Result, nil
}

func initialModel() model {

	return model{
		Choices: initialChoices,

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		Selected: make(map[int]string),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) SelectionView(initalMessage string, choices []string) string {
	for i, choice := range choices {
		cursor := " " // no cursor
		if m.Cursor == i {
			cursor = ">" // cursor!
		}
		chcked := " " // not selected
		if _, ok := m.Selected[i]; ok {
			chcked = "x" // selected!
		}

		initalMessage += fmt.Sprintf("%s [%s] %s\n", cursor, chcked, choice)
	}
	return initalMessage
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
