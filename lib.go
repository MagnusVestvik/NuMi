package main

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"strconv"
	"strings"
)

func splitFunc(r rune) bool {
	return r == ' ' || r == '|' || r == '\n' || r == '=' || r == '-' || r == '>'
}

func getNamesFromSearchResult(searchResult string) []string {
	logMu.Lock()
	logger.Printf("Recived search result for getNames function: %v", searchResult)
	logMu.Unlock()
	names := []string{strings.FieldsFunc(searchResult, splitFunc)[0]}
	logMu.Lock()
	logger.Printf("Names from recived search result: %v", names)
	logMu.Unlock()
	return names
}

func getVersionsFromSearchResult(searchResult string) []string {
	versions := []string{strings.FieldsFunc(searchResult, splitFunc)[1]}
	return versions
}

func getNumDownloadsFromSearchResult(searchResult string) []string {
	downloads := []string{strings.FieldsFunc(searchResult, splitFunc)[2]}
	return downloads
}

func getDescriptionsFromSearchResult(searchResult string) []string {
	descriptions := []string{strings.Split(searchResult, "\n")[1]}
	return descriptions
}

func arrangeSearchResultTable(searchResult SearchResult) table.Model {
	columns := []table.Column{
		{Title: "Number", Width: 10},
		{Title: "Name", Width: 30},
		{Title: "Version", Width: 10},
		{Title: "Downloads", Width: 10},
		{Title: "Description", Width: 30},
	}

	rows := []table.Row{make([]string, len(searchResult.PackageName))}
	for i, name := range searchResult.PackageName {
		rows[i] = table.Row{
			strconv.Itoa(i + 1),
			name,
			searchResult.PackageVersion[i],
			searchResult.PackageDescription[i],
			searchResult.PackageDownloads[i],
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(40),
		table.WithWidth(100),
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
