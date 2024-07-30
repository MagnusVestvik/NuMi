package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ViewModel interface {
	View() string
	Update(tea.Msg) (tea.Model, tea.Cmd)
}

type GlobalKeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	Quit  key.Binding
}

type BaseModel struct {
	height int
	width  int
	keys   GlobalKeyMap
	style  lipgloss.Style
}

type SearchViewModel struct {
	BaseModel
	packageSearchTable table.Model
	selectSearchTable  bool
	isSearching        bool
	cursor             int
	inputField         textinput.Model
	progressBar        progress.Model
}

type MainViewModel struct {
	BaseModel
	viewList list.Model
}

type ListPackageViewModel struct {
	BaseModel
}

func initSearchViewModel(width int, height int) SearchViewModel {
	ti := textinput.New()
	ti.Placeholder = "Search for packages"
	ti.Focus()
	ti.Width = 20
	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))
	baseViewModel := BaseModel{
		height: height,
		width:  width,
		style:  baseStyle,
	}

	searchViewModel := SearchViewModel{
		BaseModel:         baseViewModel,
		inputField:        ti,
		progressBar:       progress.New(progress.WithDefaultGradient()),
		selectSearchTable: false,
	}
	searchViewModel.SetSize(width, height)

	return searchViewModel
}

func initMainViewModel(width int, height int) MainViewModel {
	choices := getMainViewChoices()
	startViewList := list.New(choices, list.NewDefaultDelegate(), 0, 0)
	startViewList.Title = "Main Menu"
	listBaseStyle := lipgloss.NewStyle().Margin(1, 2)
	baseViewModel := BaseModel{
		height: height,
		width:  width,
		style:  listBaseStyle,
	}
	mainViewModel := MainViewModel{
		BaseModel: baseViewModel,
		viewList:  startViewList,
	}
	mainViewModel.SetSize(width, height)

	return mainViewModel
}

func initListPackageViewModel() ListPackageViewModel {
	return ListPackageViewModel{}
}

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
	svm.progressBar.Width = width - 4 // Adjust for margins
}
