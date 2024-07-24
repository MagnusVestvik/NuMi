package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"os/exec"
	"strings"
)

func runNuGetCommand(args ...string) (string, error) {
	cmd := exec.Command("nuget ", args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func SearchPackagesCmd(args ...string) tea.Cmd {
	return func() tea.Msg {
		results, err := runNuGetCommand(args...)
		if err != nil {
			return err
		}
		return strings.Split(results, "\n")
	}
}

func InstallNugetPackage(packageName string) tea.Cmd {
	return func() tea.Msg {
		result, err := runNuGetCommand("install", packageName)
		if err != nil {
			return err
		}
		return result
	}
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
