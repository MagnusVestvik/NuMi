package main

import (
	"github.com/charmbracelet/bubbles/key"
)

const (
	MainViewState = iota
	SearchViewState
	ProgressViewState
	ListInstalledPackagesViewState
	HelpViewState
)

var (
	globalKeys = keyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("?/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("?/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("?/h", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("?/l", "move right"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctr+c"),
			key.WithHelp("q/ctrl+c", "quit"),
		),
	}
)
