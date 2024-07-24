package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

type ListNugetPackages interface { // TODO: det her er teit kinda...
	ListPackages() ([]string, error)
}

type NugetSearchResultMessage struct {
	Result []string
	Error  error
}

type SearchNugetPackages struct {
	SearchTerm string
	Options    string
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

func SearchNugetPackes(nugetPackage SearchNugetPackages) tea.Cmd {
	return func() tea.Msg {
		// TODO: implement search
		return NugetSearchResultMessage{
			Result: []string{"test1", "test2"},
			Error:  nil,
		}
	}
}

func InstallNugetPackage(nugetPackage string) tea.Cmd {
	return func() tea.Msg {
		// TODO: implement install
		return InstallNugetPackagesMessage{
			Package: nugetPackage,
			Error:   nil,
		}
	}
}

func UninstallNugetPackage(nugetPackages []string) tea.Cmd { // todo bruk egnet struct for nugetPackages
	return func() tea.Msg {
		// TODO: implement uninstall
		return UninstallNugetPackageMessage{ // TODO: dette er feil, må returnere en kode for success eller failure ellernoe
			Package: nugetPackages,
			Error:   nil,
		}
	}
}

type model struct {
	choices  []string         // items on the to-do list
	cursor   int              // which to-do list item our cursor is pointing at
	selected map[int]struct{} // which to-do items are selected
}

func initialModel() model {
	return model{
		// Our to-do list is a grocery list
		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) View() string {
	// The header
	s := "What should we buy at the market?\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
