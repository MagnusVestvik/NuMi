package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"time"
)

type ViewState struct{ state int }

type SearchResult struct {
	Result     []string
	SearchTerm string
}

type tickMsg time.Time

func SearchPackagesCmd(args ...string) tea.Cmd {
	return func() tea.Msg {
		logMu.Lock()
		logger.Printf("Executing SearchPackagesCmd with args: %v", args)
		logMu.Unlock()
		response, err := runNuGetCommand("search", args[0]) // TODO: handle multiple args
		response = strings.Replace(response, "Source: nuget.org", "", 1)

		response = strings.Replace(response, "Source: nuget.org", "", 1)
		if err != nil {
			logMu.Lock()
			logger.Printf("Error in SearchPackagesCmd: %v", err)
			logMu.Unlock()
			return err
		}
		searchResult := SearchResult{strings.Split(response, "--------------------"), args[0]}

		return searchResult
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
