package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"os/exec"
	"strings"
	"time"
)

type SearchResult struct {
	Result     []string
	SearchTerm string
}

type tickMsg time.Time

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

		logMu.Lock()
		logger.Printf("SearchPackagesCmd result: %v", searchResult)
		logMu.Unlock()
		return searchResult
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func InstallNugetPackage(packageName string) string {
	result, err := runNuGetCommand("install", packageName)
	if err != nil {
		return err.Error()
	}
	return result
}

func UninstallNugetPackage(packageName string) tea.Cmd {
	return func() tea.Msg {
		result, err := runNuGetCommand("uninstall", packageName)
		if err != nil {
			return err
		}
		return result
	}
}

func SearchPackages(args ...string) ([]string, error) {
	results, err := runNuGetCommand(args...)
	if err != nil {
		return nil, err
	}
	return strings.Split(results, "\n"), nil
}
