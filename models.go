package main

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
	selectSearchTable  bool
	cursor             int
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

func initSearchViewModel() SearchViewModel {
	ti := textinput.New()
	ti.Placeholder = "Search for packages"
	ti.Focus()
	ti.Width = 20
	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	searchViewModel := SearchViewModel{
		inputField:        ti,
		progressBar:       progress.New(progress.WithDefaultGradient()),
		style:             baseStyle,
		selectSearchTable: false,
	}
	searchViewModel.SetSize(80, 24)

	return searchViewModel
}

func initMainViewModel() MainViewModel {
	choices := getMainViewChoices()
	startViewList := list.New(choices, list.NewDefaultDelegate(), 0, 0)
	startViewList.Title = "Main Menu"
	listBaseStyle := lipgloss.NewStyle().Margin(1, 2)
	mainViewModel := MainViewModel{
		viewList: startViewList,
		style:    listBaseStyle,
	}
	mainViewModel.SetSize(80, 24)

	return mainViewModel
}

func initListPackageViewModel() ListPackageViewModel {
	return ListPackageViewModel{}
}

func (svm SearchViewModel) GetHeight() int { return svm.height }

func (mvm MainViewModel) GetHeight() int { return mvm.height }

func (lvm ListPackageViewModel) GetHeight() int { return lvm.height }

func (svm SearchViewModel) GetWidth() int { return svm.width }

func (mvm MainViewModel) GetWidth() int { return mvm.width }

func (lvm ListPackageViewModel) GetWidth() int { return lvm.width }

func (svm SearchViewModel) Init() tea.Cmd { return nil }

func (mvm MainViewModel) Init() tea.Cmd { return nil }

func (lvm ListPackageViewModel) Init() tea.Cmd { return nil }

func (mvm *MainViewModel) SetSize(width, height int) {
	mvm.width = width
	mvm.height = height
	mvm.viewList.SetSize(width-4, height-4) // Adjust for margins
}

func (svm *SearchViewModel) SetSize(width, height int) {
	svm.width = width
	svm.height = height
	svm.inputField.Width = width - 4
	svm.progressBar.Width = width - 4
}
