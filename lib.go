package main

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"regexp"
	"strconv"
)

func getNamesFromSearchResult(searchResult string) string {
	re := regexp.MustCompile(`> (.*?) \|`) // TODO: finner ingen navn som matcher
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
func arrangeSearchResultTable(searchResult SearchResult) table.Model {
	columns := []table.Column{
		{Title: "Number", Width: 10},
		{Title: "Name", Width: 30},
		{Title: "Version", Width: 10},
		{Title: "Downloads", Width: 10},
		{Title: "Description", Width: 120},
	}

	searchResult.Result = searchResult.Result[1:] // Remove the first element, which is a header
	rows := make([]table.Row, len(searchResult.Result))
	for i, pkg := range searchResult.Result {
		if len(pkg) <= 2 { // fixes issue where last line is empty
			continue
		}
		logMu.Lock()
		logger.Printf("Pkg in searchResult: %v", pkg)
		logMu.Unlock()
		rows[i] = table.Row{
			strconv.Itoa(i + 1),
			getNamesFromSearchResult(pkg),
			getVersionsFromSearchResult(pkg),
			getNumDownloadsFromSearchResult(pkg),
			getDescriptionsFromSearchResult(pkg),
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(40),
		table.WithWidth(190),
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

func center(m model, s string) string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, s)

}
