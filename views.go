package main

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func (mvm MainViewModel) View() string {
	s := mvm.viewList.View()
	help := mvm.help.View(mvm.keys)
	s += "\n" + help
	return center(mvm.BaseModel, mvm.style.Render(s))
}

func (svm SearchViewModel) renderProgressBar() string {
	logMu.Lock()
	logger.Printf("ProgressView rendering! ")
	logMu.Unlock()

	pad := strings.Repeat(" ", 2)
	pad += "\n" +
		pad + svm.progressBar.View() + "\n\n"

	return center(svm.BaseModel, pad)
}

func (svm SearchViewModel) renderInstalledPackages() string {
	selectedPackagesStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1).
		Width(svm.width / 2)
	return selectedPackagesStyle.Render(svm.ViewSelectedPackages())
}

func (svm SearchViewModel) renderSearchInputField() string {
	inputFieldStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1).
		Width(svm.width / 2)

	return inputFieldStyle.Render(svm.inputField.View())

}

func (svm SearchViewModel) View() string {
	if svm.isSearching {
		return svm.renderProgressBar()
	}

	container := lipgloss.NewStyle().
		Width(svm.width).
		Height(svm.height).
		Align(lipgloss.Center, lipgloss.Center)

	installedPackages := svm.renderInstalledPackages()
	inputField := svm.renderSearchInputField()

	// Indicates that no packages have been searched yet
	if len(svm.searchedPackages.Rows()) == 0 {
		return container.Render(inputField)
	}

	inputField = lipgloss.JoinVertical(lipgloss.Center, inputField, "\n", svm.style.Render(svm.searchedPackages.View())+"\n") // TROR dette bare er tull...

	if len(svm.installedPackages.packages.Items()) == 0 {
		logMu.Lock()
		logger.Printf("There are no installed packages yet")
		logMu.Unlock()
		return container.Render(inputField)
	}
	logMu.Lock()
	logger.Printf("Installed packages is not empty, it contains: ", svm.installedPackages.packages.Items())
	logMu.Unlock()

	toRender := lipgloss.JoinVertical(lipgloss.Center, installedPackages, "\n", inputField)
	return container.Render(toRender)
}

func (svm SearchViewModel) ViewSelectedPackages() string {
	s := ""
	// Renders progressbars for the selected packages
	if svm.installedPackages.isDownloading {
		for i := 0; i < len(svm.installedPackages.progressBars); i++ {
			s += svm.installedPackages.progressBars[i].View() + "\n"
		}
		return s
	}

	packages := svm.installedPackages.packages.Items()
	if len(packages) == 0 {
		logMu.Lock()
		logger.Printf("No packages selected in : " + svm.installedPackages.packages.View())
		logMu.Unlock()
		return svm.installedPackages.packages.Title
	}
	logMu.Lock()
	logger.Printf("Rendering following packages: " + svm.installedPackages.packages.View())
	logMu.Unlock()
	return svm.installedPackages.packages.View()
}

func (lvm ListPackageViewModel) View() string {
	return "" // TODO: Implement
}
