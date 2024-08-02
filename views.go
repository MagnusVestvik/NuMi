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

	selectPackages := lipgloss.PlaceVertical(svm.width, lipgloss.Top, svm.ViewSelectedPackages())
	s := svm.inputField.View() + "\n\n"

	if len(svm.packageSearchTable.Rows()) == 0 {
		s += "No results yet. Press Enter to search.\n"
	} else {
		logMu.Lock()
		logger.Printf("PackageSearch table is: %v", svm.packageSearchTable.View()) // TODO: fix sizing of downloads in packageSearchTable
		logMu.Unlock()
		s += svm.style.Render(svm.packageSearchTable.View()) + "\n"
	}
	help := svm.help.View(svm.keys)

	s += "\n" + help
	toRender := selectPackages + s
	return center(svm.BaseModel, toRender)
}

func (svm SearchViewModel) ViewSelectedPackages() string {
	s := ""
	if svm.selectedPackages.isDownloading {
		for i := 0; i < len(svm.selectedPackages.prgressBars); i++ {
			s += svm.selectedPackages.prgressBars[i].View() + "\n"
		}
		return s
	}
	return svm.selectedPackages.packages.View()

}

func (lvm ListPackageViewModel) View() string {
	return "" // TODO: Implement
}
