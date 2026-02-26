package tui

import (
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/devamaz/clipshistory/internal/store"
	"golang.design/x/clipboard"
)

type Position int

const (
	SearchPane Position = iota
	PreviewPane
	ContentPane
	InfoPane
	PinPane
)

type model struct {
	store         *store.Store
	state         Position
	search        textinput.Model
	currClip      *store.Clip
	clips         []store.Clip
	width         int
	height        int
	selectedIndex int
	contentScroll int
	previewScroll int
	focused       Position
	help          help.Model
	viewport      viewport.Model
}

func NewModel(store *store.Store) model {
	clips, err := store.GetClips()
	if err != nil {
		log.Fatalf("unable to get notes: %v", err)
	}

	search := textinput.New()
	search.Prompt = "Filter: "
	return model{
		store:         store,
		search:        search,
		help:          help.New(),
		viewport:      viewport.New(0, 0),
		clips:         clips,
		width:         100,
		height:        30,
		selectedIndex: 0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Update viewport dimensions to fit in the left preview panel
		// Panel width: m.width/4 - 2, Panel height: m.height - 6
		// Subtract 2 for padding (1 on each side) on both dimensions
		m.viewport.Width = m.width/4 - 4
		m.viewport.Height = m.height - 10
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Navigation.Quit):
			return m, tea.Quit
		case key.Matches(msg, Navigation.Up):
			if m.focused == PreviewPane {
				m.viewport.LineUp(1)
			} else if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case key.Matches(msg, Navigation.Down):
			if m.focused == PreviewPane {
				m.viewport.LineDown(1)
			} else if m.selectedIndex < len(m.clips)-1 {
				m.selectedIndex++
			}
		case key.Matches(msg, Navigation.Copy):
			clip := m.selectedClip()
			if clip != nil {
				clipboard.Write(clipboard.FmtText, []byte(clip.Content))
			}
		case key.Matches(msg, Navigation.SearchPane):
			m.focusPane(SearchPane)
		case key.Matches(msg, Navigation.PreviewPane):
			m.focusPane(PreviewPane)
		case key.Matches(msg, Navigation.PinPane):
			m.focusPane(PinPane)
		case key.Matches(msg, Navigation.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}
	return m, nil
}
