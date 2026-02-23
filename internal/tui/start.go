package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/devamaz/clipshistory/internal/app"
	"github.com/devamaz/clipshistory/internal/store"
)

func Start(store *store.Store) error {

	// app.MonitorClipboard(store)

	m := NewModel(store)
	p := tea.NewProgram(m, tea.WithAltScreen())

	_, err := p.Run()
	return err
}
