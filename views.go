package main

import (
	"strings"
)

func (mvm MainViewModel) View() string {
	return center(mvm, mvm.style.Render(mvm.viewList.View()))
}

func (svm SearchViewModel) View() string {
	s := svm.inputField.View() + "\n\n"

	if svm.isSearching {
		logMu.Lock()
		logger.Printf("ProgressView rendering! ")
		logMu.Unlock()

		pad := strings.Repeat(" ", 2) // 2 is padding
		pad += "\n" +
			pad + svm.progressBar.View() + "\n\n"

		return center(svm, pad)
	}

	if len(svm.packageSearchTable.Rows()) == 0 {
		s += "No results yet. Press Enter to search.\n"
	} else {
		logMu.Lock()
		logger.Printf("PackageSearch table is: %v", svm.packageSearchTable.View())
		logMu.Unlock()
		s += svm.style.Render(svm.packageSearchTable.View()) + "\n"
	}
	return center(svm, s)
}

func (lvm ListPackageViewModel) View() string {
	return "" // TODO: Implement
}
