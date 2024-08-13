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

func (svm SearchViewModel) View() string {
	if svm.isSearching { // TODO: Fix sizing of progress bar
		logMu.Lock()
		logger.Printf("ProgressView rendering! ")
		logMu.Unlock()

		pad := strings.Repeat(" ", 2) // 2 is padding
		pad += "\n" +
			pad + svm.progressBar.View() + "\n\n"

		return center(svm.BaseModel, pad)
	}
	container := lipgloss.NewStyle().
		Width(svm.width).
		Height(svm.height).
		Align(lipgloss.Center, lipgloss.Center)

	selectedPackagesStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1).
		Width(svm.width / 2)

	inputFieldStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1).
		Width(svm.width / 2)

	selectPackages := selectedPackagesStyle.Render(svm.ViewSelectedPackages())
	inputField := inputFieldStyle.Render(svm.inputField.View())

	if len(svm.packageSearchTable.Rows()) == 0 {
		inputField = lipgloss.JoinVertical(lipgloss.Center, inputField, "\n", "No results yet. Press Enter to search.")
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
	s := "" // Bør vær Strings.Builder istedenfor
	if svm.selectedPackages.isDownloading {
		for i := 0; i < len(svm.selectedPackages.progressBars); i++ {
			s += svm.selectedPackages.progressBars[i].View() + "\n"
		}
		return s
	}

	packages := svm.selectedPackages.packages.Items()
	if len(packages) == 0 {
		return svm.selectedPackages.packages.Title
	}

	return svm.selectedPackages.packages.View()
}

func (lvm ListPackageViewModel) View() string {
	return "" // TODO: Implement
}
