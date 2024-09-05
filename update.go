package main

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

func (mvm MainViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	logMu.Lock()
	logger.Printf("Received message: %T", msg)
	logMu.Unlock()
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		mvm.SetSize(msg.Width, msg.Height)
		return mvm, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			logMu.Lock()
			logger.Printf("Enter key was pressed and will now try to change to model: %v", mvm.viewList.Index())
			logMu.Unlock()
			newModel, err := ChangeViewState(mvm.viewList.Index()+1, mvm.BaseModel) // pluss one because idx 0 is mainview
			if err != nil {
				logMu.Lock()
				logger.Printf("Error in ChangeViewState: %v", err)
				logMu.Unlock()
				return mvm, nil
			}
			logMu.Lock()
			logger.Printf("Succesffuly change viewmodel to model number: %v", mvm.viewList.Index())
			logMu.Unlock()
			return newModel, nil
		}
		var cmd tea.Cmd
		mvm.viewList, cmd = mvm.viewList.Update(msg)
		return mvm, cmd
	}
	return mvm, nil
}

func (svm SearchViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	logMu.Lock()
	logger.Printf("Received message: %T", msg)
	logMu.Unlock()
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		svm.SetSize(msg.Width, msg.Height)
		return svm, nil

	case tickMsg:
		logMu.Lock()
		logger.Printf("TickMsg received: %v", msg)
		logMu.Unlock()
		if svm.progressBar.Percent() == 1.0 {
			svm.isSearching = false
			return svm, svm.progressBar.SetPercent(0)
		}
		cmd := svm.progressBar.IncrPercent(0.25)
		logMu.Lock()
		logger.Printf("Progress incremented: %v", svm.progressBar.Percent())
		logger.Printf("Progress is currently at value: %v", svm.progressBar.View())
		logMu.Unlock()
		return svm, tea.Batch(tickCmd(), cmd)

	case progress.FrameMsg:
		progressModel, cmd := svm.progressBar.Update(msg)
		svm.progressBar = progressModel.(progress.Model)
		return svm, cmd

	case SearchResult:
		svm.searchedPackages = arrangeSearchResultTable(msg, svm.width)
		logMu.Lock()
		logger.Printf("updated table: %v", msg)
		logMu.Unlock()
		svm.cursor = 0
		return svm, nil

	case InstallPackage:
		svm.installedPackages.packages = arrangeInstalledPackagesTable(svm, msg)
		logMu.Lock()
		logger.Printf("a package was installed and installed packages now looks like this: ", svm.installedPackages.packages)
		logMu.Unlock()
		return svm, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "-":
			newModel, err := ChangeViewState(MainViewState, svm.BaseModel)
			if err != nil {
				logMu.Lock()
				logger.Printf("Error in ChangeViewState: %v", err)
				logMu.Unlock()
				return svm, nil
			}
			return newModel, nil

		case "esc":
			svm.searchedPackagesIsSelected = true
		case "tab":
			logMu.Lock()
			logger.Printf("Tab key was pressed and resulted in selectSearchTable: %v", svm.searchedPackagesIsSelected)
			logMu.Unlock()

			svm.searchedPackagesIsSelected = !svm.searchedPackagesIsSelected
		case "ctrl+c", "q":
			return svm, tea.Quit
		case "up", "k":
			if svm.cursor > 0 {
				svm.cursor--
			}
		case "down", "j":
			if svm.cursor < len(svm.searchedPackages.Rows())-1 {
				svm.cursor++
			}
		case "i":
			if svm.searchedPackagesIsSelected {
				logMu.Lock()
				logger.Printf("i key was pressed will now try to install package: %v", svm.searchedPackages.Rows()[svm.cursor][0])
				logMu.Unlock()
				return svm, InstallPackageCmd(svm.searchedPackages.Rows()[svm.cursor][0])
			}
		case "p":
			if svm.searchedPackagesIsSelected {
				logMu.Lock()
				logger.Printf("Add package to selected packages")
				logMu.Unlock()
			}
			return svm, nil
		case "enter":
			if svm.searchedPackagesIsSelected {
				return svm, InstallPackageCmd(svm.searchedPackages.Rows()[svm.cursor][0])
			}
			logMu.Lock()
			logger.Printf("Serching for package with name of %v", svm.inputField.Value())
			logMu.Unlock()
			svm.isSearching = true
			tickIncrCmd := svm.progressBar.IncrPercent(0.25)
			return svm, tea.Batch(tickCmd(), tickIncrCmd, SearchPackagesCmd(svm.inputField.Value()))
		}
		var cmd tea.Cmd
		if !svm.searchedPackagesIsSelected {
			svm.inputField, cmd = svm.inputField.Update(msg)
		} else {
			svm.searchedPackages.SetCursor(svm.cursor)
			logMu.Lock()
			logger.Printf("Table moved cursor to position: %v", svm.searchedPackages.Cursor())
			logMu.Unlock()
		}
		return svm, cmd
	}
	return svm, nil
}

func (lvm ListPackageViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return lvm, nil
}
