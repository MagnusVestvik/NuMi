package refactor

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ViewModel interface {
	GetHeight() int
	GetWidth() int
	View() string
	Update(tea.Msg) (tea.Model, tea.Cmd)
}

type SearchViewModel struct {
	packageSearchTable table.Model
	inputField         textinput.Model
	progressBar        progress.Model
	style              lipgloss.Style
	height             int
	width              int
}

type MainViewModel struct {
	viewList list.Model
	height   int
	width    int
	style    lipgloss.Style
}

type ListPackageViewModel struct {
	height int
	width  int
	style  lipgloss.Style
}

type model struct { // TODO: denne skal byttes ut med main view model, sålenge alle modeller implementerer tea.model kan man bare bytte
	cursor     int
	showLogs   bool
	windowSize tea.WindowSizeMsg
	viewState  int
	views      []ViewModel
	height     int
	width      int
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
	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	return SearchViewModel{
		inputField:  textinput.New(),
		progressBar: progress.New(),
		style:       baseStyle,
	}

}

func initMainViewModel() MainViewModel {
	choices := getStartViewChoices()
	startViewList := list.New(choices, list.NewDefaultDelegate(), 0, 0)
	startViewList.Title = "Main Menu"
	listBaseStyle := lipgloss.NewStyle().Margin(1, 2)

	return MainViewModel{
		viewList: startViewList,
		style:    listBaseStyle,
	}
}

func initListPackageViewModel() ListPackageViewModel {
	return ListPackageViewModel{}
}

func (svm SearchViewModel) GetHeight() int { return svm.height }

func (mvm MainViewModel) GetHeight() int { return mvm.height }

func (lvm ListPackageViewModel) GetHeight() int { return lvm.height }

func (m model) GetHeight() int { return m.height }

func (svm SearchViewModel) GetWidth() int { return svm.width }

func (mvm MainViewModel) GetWidth() int { return mvm.width }

func (lvm ListPackageViewModel) GetWidth() int { return lvm.width }

func (m model) GetWidth() int { return m.width }
