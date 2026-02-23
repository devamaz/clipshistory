package tui

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/devamaz/clipshistory/internal/store"
	"golang.design/x/clipboard"
)

type Position int

const (
	SearchPane Position = iota
	PreviewListPane
	ContentPane
	InfoPane
	PinListPane
)

type model struct {
	store         *store.Store
	state         Position
	textinput     textinput.Model
	currClip      *store.Clip
	clips         []store.Clip
	width         int
	height        int
	selectedIndex int
}

func NewModel(store *store.Store) model {
	clips, err := store.GetClips()
	if err != nil {
		log.Fatalf("unable to get notes: %v", err)
	}

	return model{
		store:         store,
		textinput:     textinput.New(),
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
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "k", "up":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "j", "down":
			if m.selectedIndex < len(m.clips)-1 {
				m.selectedIndex++
			}
		case "y":
			clip := m.selectedClip()
			if clip != nil {
				clipboard.Write(clipboard.FmtText, []byte(clip.Content))
			}

		}
	}
	return m, nil
}

func (m model) View() string {
	// Define styles
	searchStyle := lipgloss.NewStyle().
		Width(m.width-4).
		Height(2).
		Align(lipgloss.Center, lipgloss.Center).
		Bold(true)

	previewStyle := lipgloss.NewStyle().
		Width(m.width/4 - 2).
		Height(m.height - 8).
		Padding(1)

	contentStyle := lipgloss.NewStyle().
		Width(m.width/2 - 2).
		Height(m.height - 22).
		Padding(1)

	infoStyle := lipgloss.NewStyle().
		Width(m.width/2 - 2).
		Height(12).
		Padding(1)

	pinStyle := lipgloss.NewStyle().
		Width(m.width/4 - 2).
		Height(m.height - 8).
		Padding(1)

	// Create content
	header := borderize(searchStyle.Render(""), true, m.SearchBorderText())

	left := borderize(previewStyle.Render(m.renderPreviewList()), false, m.PreviewBorderText())

	info := borderize(infoStyle.Render(m.renderInfo()), false, m.InfoBorderText())

	center := borderize(contentStyle.Render(m.renderContent()), false, m.ContentBorderText())

	right := borderize(pinStyle.Render("Pinned clips \n\n Coming soon"), false, m.PinBorderText())

	centerLayout := lipgloss.JoinVertical(lipgloss.Left, center, info)

	// Combine the three bottom boxes horizontally
	bottomRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		centerLayout,
		right,
	)

	// Combine header and bottom row vertically
	layout := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		bottomRow,
	)

	// Add instructions at the bottom
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("\nPress ↑/↓ to navigate • 'q' to quitt")

	return layout + instructions
}

func (m model) PreviewBorderText() map[BorderPosition]string {
	return map[BorderPosition]string{
		TopLeftBorder:  " PREVIEW ",
		TopRightBorder: " 1 ",
	}
}

func (m model) SearchBorderText() map[BorderPosition]string {
	return map[BorderPosition]string{
		TopLeftBorder:  " SEARCH ",
		TopRightBorder: " 0 ",
	}
}

func (m model) PinBorderText() map[BorderPosition]string {
	return map[BorderPosition]string{
		TopLeftBorder:  " PIN ",
		TopRightBorder: " 2 ",
	}
}

func (m model) ContentBorderText() map[BorderPosition]string {
	return map[BorderPosition]string{
		TopLeftBorder:  " CONTENT ",
		TopRightBorder: " 3 ",
	}
}

func (m model) InfoBorderText() map[BorderPosition]string {
	return map[BorderPosition]string{
		TopLeftBorder:  " INFO ",
		TopRightBorder: " 4 ",
	}
}

// Render the preview list with selection indicator
func (m model) renderPreviewList() string {
	if len(m.clips) == 0 {
		return "No clips available"
	}

	var sb strings.Builder
	for i, clip := range m.clips {
		// Add selection indicator
		if i == m.selectedIndex {
			sb.WriteString("> ") // Active selection
		} else {
			sb.WriteString("  ") // Regular item
		}

		// Show preview with character count
		sb.WriteString(clip.Preview)
		if len(clip.Preview) < 15 {
			sb.WriteString(strings.Repeat(" ", 15-len(clip.Preview)))
		}
		sb.WriteString(fmt.Sprintf(" (%d chars)\n", clip.CharCount))
	}

	return sb.String()
}

// Get the currently selected clip
func (m model) selectedClip() *store.Clip {
	if len(m.clips) == 0 || m.selectedIndex >= len(m.clips) {
		return nil
	}
	return &m.clips[m.selectedIndex]
}

// Render the content of selected clip
func (m model) renderContent() string {
	clip := m.selectedClip()
	if clip == nil {
		return "No clip selected"
	}

	// Handle multi-line content
	content := strings.ReplaceAll(clip.Content, "\n", "\n")

	return fmt.Sprintf("Full Content:\n\n%s", content)
}

func formatDate(s int64) string {
	tm := time.Unix(s, 0)
	now := time.Now()

	y1, m1, d1 := tm.Date()
	y2, m2, d2 := now.Date()

	if y1 == y2 && m1 == m2 && d1 == d2 {
		return "Today at " + tm.Format("15:04:05")
	}

	yesterday := now.AddDate(0, 0, -1)
	y3, m3, d3 := yesterday.Date()
	if y1 == y3 && m1 == m3 && d1 == d3 {
		return "Yesterday at " + tm.Format("15:04:05")
	}

	return tm.Format("2006-01-02") + " at " + tm.Format("15:04:05")
}

// Render info about selected clip
func (m model) renderInfo() string {

	clip := m.selectedClip()
	if clip == nil {
		return "Select a clip to see info"
	}

	formattedDate := formatDate(clip.CreatedAt)
	formattedLastCopied := formatDate(clip.LastCopiedAt)

	return fmt.Sprintf(`Clip Details:
─────────────────
Char Count:   								%d
Times:			   								%d
Words:												%d
Pinned:       								%t
Created At:   								%s
Last Copied:  								%s`,
		clip.CharCount,
		clip.CopyCount,
		len(strings.Fields(clip.Content)),
		clip.IsPinned,
		formattedDate,
		formattedLastCopied,
	)
}
