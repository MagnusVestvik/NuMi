package main

import (
	"errors"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

var (
	logger *log.Logger
	logMu  sync.Mutex
	logBuf strings.Builder
)

type keys struct {
	navigationKeys []key.Binding // navigation keys
	actionKeys     []key.Binding // action keys
	helpKeys       []key.Binding // help key
}

func initLogger() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	logger = log.New(io.MultiWriter(logFile, &logBuf), "", log.Ltime|log.Lshortfile)
}

type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func ChangeViewState(viewState int, base BaseModel) (tea.Model, error) {
	switch viewState {
	case MainViewState:
		return initMainViewModel(base), nil
	case SearchViewState:
		return initSearchViewModel(base), nil
	}
	return nil, errors.New("Invalid view state selected with viewState: " + string(viewState)) // TODO: fix this to not return a string of runes but rather a string of the actuall number
}

func runNuGetCommand(args ...string) (string, error) {
	cmd := exec.Command("nuget", args...)
	logMu.Lock()
	logger.Printf("Running command: %v", cmd)
	logMu.Unlock()
	output, err := cmd.CombinedOutput()
	logMu.Lock()
	logger.Printf("Command output: %v", string(output))
	if err != nil {
		logger.Printf("Command error: %v", err)
	}
	logMu.Unlock()
	return string(output), err
}

func getMainViewChoices() []list.Item {
	return []list.Item{
		item{title: "Search Packages", desc: "Search the nuget library for packages"},
		item{title: "List Packages", desc: "List all installed packages in current project"},
		item{title: "Help", desc: "Get help on how to use the application"},
	}
}
func getNamesFromSearchResult(searchResult string) string {
	re := regexp.MustCompile(`> (.*?) \|`)
	matches := re.FindAllStringSubmatch(searchResult, -1)

	logMu.Lock()
	logger.Printf("Names matched: %v", matches[0][1])
	logMu.Unlock()

	name := matches[0]
	return name[1] // Returns what's inside the capture group
}

func getVersionsFromSearchResult(searchResult string) string {
	re := regexp.MustCompile(`\| (.*) \|`)
	matches := re.FindAllStringSubmatch(searchResult, -1)

	logMu.Lock()
	logger.Printf("Versions matched: %v", matches[0][1])
	logMu.Unlock()
	version := matches[0]
	return version[1] // Returns what's inside the capture group
}

func getNumDownloadsFromSearchResult(searchResult string) string {
	re := regexp.MustCompile(`Downloads: (.*)`)
	matches := re.FindAllStringSubmatch(searchResult, -1)

	logMu.Lock()
	logger.Printf("Num Downloads matched: %v", matches[0][1])
	logMu.Unlock()
	downloads := matches[0]
	return downloads[1] // Returns what's inside the capture group
}

func getDescriptionsFromSearchResult(searchResult string) string {
	re := regexp.MustCompile(`>.*\n(.*)`)
	matches := re.FindAllStringSubmatch(searchResult, -1)

	logMu.Lock()
	logger.Printf("Description matched: %v", matches[0][1])
	logMu.Unlock()
	description := matches[0]
	return description[1] // Returns what's inside the capture group
}

func stringContainsChars(s string) bool {
	re := regexp.MustCompile(`[a-zA-Z0-9]`)
	matches := re.FindAllStringSubmatch(s, -1)
	return len(matches) > 0
}

func center(m BaseModel, s string) string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, s)

}

func arrangeSearchResultTable(searchResult SearchResult, availableWidth int) table.Model {
	if strings.Contains(searchResult.Result[1], "No results found") {
		columns := []table.Column{
			{Title: "Error", Width: availableWidth - 4}, // Subtract 4 for borders
		}
		rows := []table.Row{
			{"No results found for " + searchResult.SearchTerm + ". Please try another search term."},
		}
		return table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
		)
	}

	searchResult.Result = searchResult.Result[1:] // Remove the first element and the newline at the end
	rows := make([]table.Row, len(searchResult.Result))
	logMu.Lock()
	logger.Printf("length of rows %v", len(rows))
	logger.Printf("length of searchResult.Result %v", len(searchResult.Result))
	logMu.Unlock()

	rowContentWidth := make(map[string]int)
	for i, pkg := range searchResult.Result {
		if !stringContainsChars(pkg) {
			continue
		}
		logMu.Lock()
		logger.Printf("currently at i %v", i)
		logger.Printf("currently at pkg %v", pkg)
		logMu.Unlock()

		name := getNamesFromSearchResult(pkg)
		version := getVersionsFromSearchResult(pkg)
		downloads := getNumDownloadsFromSearchResult(pkg)
		rows[i] = table.Row{
			name,
			version,
			downloads,
		}
		updateRowContentWidth(&rowContentWidth, name, version, downloads)
	}

	columns := []table.Column{
		{Title: "Name", Width: Max(rowContentWidth["Name"], len("Name"))},
		{Title: "Version", Width: Max(rowContentWidth["Version"], len("Version"))},
		{Title: "Downloads", Width: Max(rowContentWidth["Downloads"], len("Downloads"))},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t.SetStyles(s)

	return t
}

func updateRowContentWidth(rowContentWidth *map[string]int, name, version, downloads string) {
	if lipgloss.Width(name) > (*rowContentWidth)["Name"] {
		(*rowContentWidth)["Name"] = lipgloss.Width(name)
	}
	if lipgloss.Width(version) > (*rowContentWidth)["Version"] {
		(*rowContentWidth)["Version"] = lipgloss.Width(version)
	}
	if lipgloss.Width(downloads) > (*rowContentWidth)["Downloads"] {
		logMu.Lock()
		logger.Printf("DownloadsWidth: %v", (*rowContentWidth)["Downloads"])
		logger.Printf("Downloads: %v", downloads)
		logger.Printf("DownloadsWidth: %v", lipgloss.Width(downloads))
		logMu.Unlock()
		(*rowContentWidth)["Downloads"] = lipgloss.Width(downloads)
	}
}

func getTableWidth(rowWidths map[string]int) int {
	sum := 0
	for _, value := range rowWidths {
		sum += value
	}
	return sum
}

func GetMaxStringWidth(s []string) int {
	max := 0
	for _, str := range s {
		if lipgloss.Width(str) > max {
			max = lipgloss.Width(str)
		}

	}
	return max
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
