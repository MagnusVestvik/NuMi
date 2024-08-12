package main

import (
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

	selectPackages := svm.ViewSelectedPackages() + "\n\n" // Her er no items satt to ganger dette skjer mest sansynlig fordi at no items er en default verdi av en tom list.model og denne rendres to ganger av en eller annen grunn
	s := svm.inputField.View() + "\n\n"

	if len(svm.packageSearchTable.Rows()) == 0 {
		s += "No results yet. Press Enter to search.\n"
	} else {
		logMu.Lock()
		logger.Printf("PackageSearch table is: %v", svm.packageSearchTable.View()) // TODO: fix sizing of downloads in packageSearchTable
		logMu.Unlock()
		s += svm.style.Render(svm.packageSearchTable.View()) + "\n"
	}
	//help := svm.help.View(svm.keys)

	toRender := selectPackages + s
	return center(svm.BaseModel, toRender)
}

func (svm SearchViewModel) ViewSelectedPackages() string {
	s := ""
	// If no packages are currently downloading, view all selected packages
	if !svm.selectedPackages.isDownloading {
		return svm.selectedPackages.packages.View()
	}

	// Render progressBars for the packages that are currently downloading.
	for i := 0; i < len(svm.selectedPackages.progressBars); i++ {
		s += svm.selectedPackages.progressBars[i].View() + "\n"
	}
	return s

}

func (lvm ListPackageViewModel) View() string {
	return "" // TODO: Implement
}
