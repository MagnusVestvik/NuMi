package main

import (
	"errors"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	logger *log.Logger
	logMu  sync.Mutex
	logBuf strings.Builder
)

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

func ChangeViewState(viewState int) (tea.Model, error) {
	switch viewState {
	case MainViewState:
		return initListPackageViewModel(), nil
	case SearchViewState:
		return initSearchViewModel(), nil
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

func center(m ViewModel, s string) string {
	return lipgloss.Place(m.GetWidth(), m.GetHeight(), lipgloss.Center, lipgloss.Center, s)

}

func wrapText(text string, width int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}
	wrapped := words[0]
	spaceLeft := width - len(wrapped)
	for _, word := range words[1:] {
		if len(word)+1 > spaceLeft {
			wrapped += "\n" + word
			spaceLeft = width - len(word)
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + len(word)
		}
	}
	return wrapped
}

// TODO: Rewrite
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

	columns := []table.Column{
		{Title: "Number", Width: 10},
		{Title: "Name", Width: 20},
		{Title: "Version", Width: 10},
		{Title: "Downloads", Width: 10},
	}

	searchResult.Result = searchResult.Result[1:] // Remove the first element and the newline at the end
	rows := make([]table.Row, len(searchResult.Result))
	logMu.Lock()
	logger.Printf("length of rows %v", len(rows))
	logger.Printf("length of searchResult.Result %v", len(searchResult.Result))
	logMu.Unlock()

	for i, pkg := range searchResult.Result {
		if !stringContainsChars(pkg) {
			continue
		}
		logMu.Lock()
		logger.Printf("currently at i %v", i)
		logger.Printf("currently at pkg %v", pkg)
		logMu.Unlock()
		rows[i] = table.Row{
			strconv.Itoa(i + 1),
			getNamesFromSearchResult(pkg),
			getVersionsFromSearchResult(pkg),
			getNumDownloadsFromSearchResult(pkg),
		}

	}

	maxWidth := FindLargestStringInSlice(searchResult.Result)
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)+2),
		table.WithWidth(maxWidth+4),
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

func FindLargestStringInSlice(s []string) int {
	max := 0
	for _, str := range s {
		if len(str) > max {
			max = len(str)
		}
	}
	return max
}
