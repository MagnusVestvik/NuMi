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
	HELPER_KEYS = keys{
		navigationKeys: []key.Binding{
			key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp("?/k", "move up"),
			),
			key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp("?/j", "move down"),
			),
			key.NewBinding(
				key.WithKeys("left", "h"),
				key.WithHelp("?/h", "move left"),
			),
			key.NewBinding(
				key.WithKeys("right", "l"),
				key.WithHelp("?/l", "move right"),
			),
		},
		actionKeys: []key.Binding{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter/i", "install"),
			),
			key.NewBinding(
				key.WithKeys("q", "ctr+c"),
				key.WithHelp("q/ctrl+c", "quit"),
			),
		},
		helpKeys: []key.Binding{
			key.NewBinding(
				key.WithKeys("?"),
				key.WithHelp("helpe", "help"),
			),
		},
	}
)
