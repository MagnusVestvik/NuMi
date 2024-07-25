package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"os/exec"
	"strings"
)

type SearchResult struct {
	PackageName string
	Result      []string
}

func runNuGetCommand(args ...string) (string, error) {
	cmd := exec.Command("nuget", args...)
	logMu.Lock()
	logger.Printf("Running command: %v", cmd)
	logMu.Unlock()
	output, err := cmd.CombinedOutput()
	logMu.Lock()
	logger.Printf("Command output: %v", string(output))
	logger.Printf("Command error: %v", err)
	logMu.Unlock()
	return string(output), err
}

func SearchPackagesCmd(args ...string) tea.Cmd {
	return func() tea.Msg {
		logMu.Lock()
		logger.Printf("Executing SearchPackagesCmd with args: %v", args)
		logMu.Unlock()
		response, err := runNuGetCommand("search", args[0]) // TODO: handle multiple args
		if err != nil {
			logMu.Lock()
			logger.Printf("Error in SearchPackagesCmd: %v", err)
			logMu.Unlock()
			return err
		}
		result := SearchResult{PackageName: args[0], Result: strings.Split(response, "--------------------")}

		logMu.Lock()
		logger.Printf("SearchPackagesCmd result: %v", result)
		logMu.Unlock()
		return result
	}
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
