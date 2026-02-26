package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

type navigation struct {
	Up          key.Binding
	Down        key.Binding
	Copy        key.Binding
	Quit        key.Binding
	SearchPane  key.Binding
	PinPane     key.Binding
	PreviewPane key.Binding
	Help        key.Binding
}

func (k navigation) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k navigation) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Copy, k.SearchPane},
		{k.PreviewPane, k.PinPane},
		{k.Help, k.Quit},
	}
}

var Navigation = navigation{
	Up:          key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:        key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Copy:        key.NewBinding(key.WithKeys("enter", "y"), key.WithHelp("enter/y", "copy")),
	Quit:        key.NewBinding(key.WithKeys("esc", "q"), key.WithHelp("esc/q", "quit")),
	SearchPane:  key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "SearchPane")),
	PinPane:     key.NewBinding(key.WithKeys("2"), key.WithHelp("2", "PinPane")),
	PreviewPane: key.NewBinding(key.WithKeys("1"), key.WithHelp("1", "PreviewPane")),
	Help:        key.NewBinding(key.WithKeys("h", "?"), key.WithHelp("h/?", "help")),
}
