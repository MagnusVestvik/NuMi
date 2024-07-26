package main

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"regexp"
	"strconv"
	"strings"
)

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

func center(m model, s string) string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, s)

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

	allColumns := []struct {
		Title      string
		GetContent func(string, int) string
	}{
		{"Number", func(pkg string, i int) string { return strconv.Itoa(i + 1) }},
		{"Name", func(pkg string, _ int) string { return getNamesFromSearchResult(pkg) }},
		{"Version", func(pkg string, _ int) string { return getVersionsFromSearchResult(pkg) }},
		{"Downloads", func(pkg string, _ int) string { return getNumDownloadsFromSearchResult(pkg) }},
	}

	searchResult.Result = searchResult.Result[1:] // Remove the first element, which is a header

	maxWidths := make([]int, len(allColumns))
	for i, col := range allColumns {
		maxWidths[i] = len(col.Title)
	}

	rows := make([]table.Row, 0, len(searchResult.Result))
	for i, pkg := range searchResult.Result {
		if len(pkg) <= 2 { // fixes issue where last line is empty
			continue
		}
		row := make(table.Row, len(allColumns))
		for j, col := range allColumns {
			content := col.GetContent(pkg, i)
			row[j] = content
			if len(content) > maxWidths[j] {
				maxWidths[j] = len(content)
			}
		}
		rows = append(rows, row)
	}

	// Determine which columns to include based on available width
	columns := make([]table.Column, 0)
	usedWidth := 0
	for i, col := range allColumns {
		columnWidth := maxWidths[i] + 2               // Add some padding
		if usedWidth+columnWidth > availableWidth-4 { // Subtract 4 for borders
			break
		}
		columns = append(columns, table.Column{Title: col.Title, Width: columnWidth})
		usedWidth += columnWidth
	}

	// Adjust rows to match selected columns
	for i, row := range rows {
		rows[i] = row[:len(columns)]
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)+2),
		table.WithWidth(usedWidth+4),
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
