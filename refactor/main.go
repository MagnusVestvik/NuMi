package refactor

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func initMainModel() model {
	modelViews := []ViewModel{
		initMainViewModel(),
		initSearchViewModel(),
		initListPackageViewModel(),
	}

	return model{
		views:     modelViews,
		viewState: MainViewState,
	}
}

func main() {
	initLogger()
	p := tea.NewProgram(initMainModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
