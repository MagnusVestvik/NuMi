package refactor

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ViewModel interface {
	View(m model) string
}

type SearchViewModel struct {
	packageSearchTable table.Model
	inputField         textinput.Model
	progress           progress.Model
}

type MainViewModel struct {
	viewList list.Model
}

type ListPackageViewModel struct{}

type model struct {
	cursor     int
	showLogs   bool
	windowSize tea.WindowSizeMsg
	viewState  int
	views      []ViewModel
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func initSearchViewModel() SearchViewModel {
	ti := textinput.New()
	ti.Placeholder = "Search for packages"
	ti.Focus()
	ti.Width = 20

	return SearchViewModel{
		inputField: textinput.New(),
		progress:   progress.New(),
	}

}

func initMainViewModel() MainViewModel {
	choices := getStartViewChoices()
	startViewList := list.New(choices, list.NewDefaultDelegate(), 0, 0)
	startViewList.Title = "Main Menu"

	return MainViewModel{
		viewList: startViewList,
	}
}

func initListPackageViewModel() ListPackageViewModel {
	return ListPackageViewModel{}
}
