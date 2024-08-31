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

func (svm SearchViewModel) renderSelectedPackages() string {
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

	selectPackages := svm.renderSelectedPackages() // TODO: dette ser bare teit ut, man kan heller bruk det til å vis hvilke pakka man har installert kanskje
	inputField := svm.renderSearchInputField()

	// Indicates that no packages have been searched yet
	if len(svm.packageSearchTable.Rows()) == 0 {
		return container.Render(inputField)
	} else {
		logMu.Lock()
		logger.Printf("PackageSearch table is: %v", svm.packageSearchTable.View())
		logMu.Unlock()
		inputField = lipgloss.JoinVertical(lipgloss.Center, inputField, "\n", svm.style.Render(svm.packageSearchTable.View())+"\n")
	}

	toRender := lipgloss.JoinVertical(lipgloss.Center, selectPackages, "\n", inputField)
	return container.Render(toRender)
}

func (svm SearchViewModel) ViewSelectedPackages() string {
	s := ""
	// Renders progressbars for the selected packages
	if svm.selectedPackages.isDownloading {
		for i := 0; i < len(svm.selectedPackages.progressBars); i++ {
			s += svm.selectedPackages.progressBars[i].View() + "\n"
		}
		return s
	}

	packages := svm.selectedPackages.packages.Items()
	if len(packages) == 0 {
		logMu.Lock()
		logger.Printf("No packages selected in : " + svm.selectedPackages.packages.View())
		logMu.Unlock()
		return svm.selectedPackages.packages.Title
	}
	logMu.Lock()
	logger.Printf("Rendering following packages: " + svm.selectedPackages.packages.View())
	logMu.Unlock()
	return svm.selectedPackages.packages.View()
}

func (lvm ListPackageViewModel) View() string {
	return "" // TODO: Implement
}
