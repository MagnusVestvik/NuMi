package main

import (
	"github.com/charmbracelet/bubbles/help"
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
	navigationKeys []key.Binding
	actionKeys     []key.Binding
	helperKeys     []key.Binding
}
type keyMap struct {
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
	keys   keyMap
	help   help.Model
	style  lipgloss.Style
}

type SearchViewModel struct {
	BaseModel
	packageSearchTable table.Model
	selectSearchTable  bool
	selectedPackages   SelectedPackages
	isSearching        bool
	cursor             int
	inputField         textinput.Model
	progressBar        progress.Model
}

// TODO: legg til view for å laste ned flere pakker, enter = legge til en pakke til ny box hvor man tilslut kan laste ned alle selected pakker
type SelectedPackages struct {
	packages      list.Model
	progressBars  []progress.Model
	isDownloading bool
}

type MainViewModel struct {
	BaseModel
	viewList list.Model
}

type ListPackageViewModel struct {
	BaseModel
}

func initSearchViewModel(baseModel BaseModel) SearchViewModel {
	ti := textinput.New()
	ti.Placeholder = "Search for packages"
	ti.Focus()
	ti.Width = 20
	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))
	baseModel.style = baseStyle
	searchViewModel := SearchViewModel{
		BaseModel:         baseModel,
		inputField:        ti,
		progressBar:       progress.New(progress.WithDefaultGradient()),
		selectSearchTable: false,
		selectedPackages:  initSelectPackages(),
	}

	return searchViewModel
}

func initMainViewModel(baseModel BaseModel) MainViewModel {
	choices := getMainViewChoices()
	startViewList := list.New(choices, list.NewDefaultDelegate(), 0, 0)
	startViewList.Title = "Main Menu"
	listBaseStyle := lipgloss.NewStyle().Margin(1, 2)
	baseModel.style = listBaseStyle
	mainViewModel := MainViewModel{
		BaseModel: baseModel,
		viewList:  startViewList,
	}

	return mainViewModel
}

func initStart() tea.Model {
	return initMainViewModel(BaseModel{
		width:  80,
		height: 24,
		help:   help.New(),
		keys:   globalKeys},
	)
}

func initSelectPackages() SelectedPackages {
	packagesList := list.New(nil, list.NewDefaultDelegate(), 0, 0) // må sette width
	packagesList.Title = "Selected Packages"
	return SelectedPackages{
		packages:      packagesList,
		isDownloading: false,
	}
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

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Help, k.Quit},                // second column
	}
}

func AddPackageToSelectedPackages(packageName string, sp *SelectedPackages) {
	newItem := item{title: packageName}
	sp.packages.InsertItem(len(sp.packages.Items()), newItem)
}
