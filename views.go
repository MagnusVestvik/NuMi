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

	inputField := svm.renderSearchInputField()

	// Indicates that no packages have been searched yet
	if len(svm.searchedPackages.Rows()) == 0 {
		return container.Render(inputField)
	}

	searchedPackagesContainer := svm.style.Render(svm.searchedPackages.View()) + "\n"

	return container.Render(lipgloss.JoinVertical(lipgloss.Center, inputField, searchedPackagesContainer))

}

func (lvm ListPackageViewModel) View() string {
	return "" // TODO: Implement
}
