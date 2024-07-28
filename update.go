package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (mvm MainViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		mvm.SetSize(msg.Width, msg.Height)
		return mvm, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			logMu.Lock()
			logger.Printf("Enter key was pressed and resulted in index: %v", mvm.viewList.Index())
			logMu.Unlock()
			newModel, err := ChangeViewState(mvm.viewList.Index() + 1) // pluss one because idx 0 is mainview
			if err != nil {
				logMu.Lock()
				logger.Printf("Error in ChangeViewState: %v", err)
				logMu.Unlock()
				return mvm, nil
			}
			return newModel, nil
		}
		var cmd tea.Cmd
		mvm.viewList, cmd = mvm.viewList.Update(msg)
		return mvm, cmd
	}
	return mvm, nil
}

func (svm SearchViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		svm.SetSize(msg.Width, msg.Height)
		return svm, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "-":
			newModel, err := ChangeViewState(MainViewState)
			if err != nil {
				logMu.Lock()
				logger.Printf("Error in ChangeViewState: %v", err)
				logMu.Unlock()
				return svm, nil
			}
			return newModel, nil

		case "tab":
			logMu.Lock()
			logger.Printf("Tab key was pressed and resulted in selectSearchTable: %v", svm.selectSearchTable)
			logMu.Unlock()

			svm.selectSearchTable = !svm.selectSearchTable
		case "ctrl+c", "q":
			return svm, tea.Quit
		case "up", "k":
			if svm.cursor > 0 {
				svm.cursor--
			}
		case "down", "j":
			if svm.cursor < len(svm.packageSearchTable.Rows())-1 {
				svm.cursor++
			}
		case "enter":
			tickIncrCmd := svm.progressBar.IncrPercent(0.25)
			return svm, tea.Batch(tickCmd(), tickIncrCmd, SearchPackagesCmd(svm.inputField.Value()))
		}
		var cmd tea.Cmd
		if !svm.selectSearchTable {
			svm.inputField, cmd = svm.inputField.Update(msg)
		} else {
			svm.packageSearchTable.SetCursor(svm.cursor)
			logMu.Lock()
			logger.Printf("Table moved cursor to position: %v", svm.packageSearchTable.Cursor())
			logMu.Unlock()
		}
		return svm, cmd
	}
	return svm, nil
}

func (lvm ListPackageViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return lvm, nil
}
