package internal

import (
	"github.com/charmbracelet/bubbles/list"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	logger *log.Logger
	logMu  sync.Mutex
	logBuf strings.Builder
)

func initLogger() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	logger = log.New(io.MultiWriter(logFile, &logBuf), "", log.Ltime|log.Lshortfile)
}

type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func getStartViewChoices() []list.Item {

	return []list.Item{
		item{title: "Search Packages", desc: "Search the nuget library for packages"},
		item{title: "List Packages", desc: "List all installed packages in current project"},
		item{title: "Help", desc: "Get help on how to use the application"},
	}
}
