package refactor

func (mvm MainViewModel) View() string {
	return center(mvm, mvm.style.Render(mvm.viewList.View()))
}

func (svm SearchViewModel) View() string {
	s := svm.inputField.View() + "\n\n"
	if len(svm.packageSearchTable.Rows()) == 0 {
		s += "No results yet. Press Enter to search.\n"
	} else {
		s += svm.style.Render(svm.packageSearchTable.View()) + "\n"
	}
	return center(svm, s)
}

func (lvm ListPackageViewModel) View() string {
	return "" // TODO: Implement
}

func (m model) View() string {
	return "" // TODO: Implement
}
