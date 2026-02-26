package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/devamaz/clipshistory/internal/store"
)

func (m model) View() string {
	// Define styles
	searchStyle := lipgloss.NewStyle().
		Width(m.width-4).
		Height(1).
		Align(lipgloss.Center, lipgloss.Center).
		Bold(true)

	previewStyle := lipgloss.NewStyle().
		Width(m.width/4 - 2).
		Height(m.height - 6).
		MaxHeight(m.height - 6).
		Padding(1)

	contentStyle := lipgloss.NewStyle().
		Width(m.width/2 - 2).
		Height(m.height - 18).
		MaxHeight(m.height - 18).
		Padding(1)

	infoStyle := lipgloss.NewStyle().
		Width(m.width/2 - 2).
		Height(8).
		Padding(1)

	pinStyle := lipgloss.NewStyle().
		Width(m.width/4 - 2).
		Height(m.height - 6).
		MaxHeight(m.height - 6).
		Padding(1)

	// Create content
	header := borderize(searchStyle.Render(m.search.View()), true, m.SearchBorderText())

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

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("\n   y/enter: copy • ?: help • q: quit")
	showHelp := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("\n" + m.help.View(Navigation))

	if m.help.ShowAll {
		layout += showHelp
	} else {
		layout += instructions
	}
	return layout
}

func (m *model) focusPane(position Position) {
	// if _, ok := m.panes[position]; !ok {
	// 	// There is no pane to focus at requested position
	// 	return
	// }
	m.focused = position
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
		if i == m.selectedIndex {
			sb.WriteString("> ")
		} else {
			sb.WriteString("  ")
		}

		trimmed := strings.TrimSpace(clip.Preview)

		if len(trimmed) > 15 {
			// Truncate to 15 chars and add ellipsis
			fmt.Fprintf(&sb, "%.25s... ...\n", trimmed)
		} else {
			fmt.Fprintf(&sb, "%s ...\n", trimmed)
		}
	}

	m.viewport.SetContent(sb.String())

	// Auto-scroll viewport to keep selected item visible
	// Only auto-scroll when preview pane is not focused (to allow manual scrolling)
	if m.viewport.Height > 0 && m.focused != PreviewPane {
		// Try to center the selected item in the viewport
		viewportCenter := m.viewport.Height / 2
		targetLine := m.selectedIndex - viewportCenter

		// Ensure we don't scroll past the beginning or end
		maxScroll := len(m.clips) - m.viewport.Height
		if targetLine < 0 {
			targetLine = 0
		} else if targetLine > maxScroll {
			targetLine = maxScroll
		}

		m.viewport.YOffset = targetLine
	}

	return m.viewport.View()
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

	// Split content into lines
	lines := strings.Split(clip.Content, "\n")

	// Get visible window
	visibleLines := min(len(lines), m.height-20) // Account for borders/padding
	start := min(m.contentScroll, len(lines)-visibleLines)
	if start < 0 {
		start = 0
	}
	end := start + visibleLines

	// Display only visible portion
	visibleContent := strings.Join(lines[start:end], "\n")

	return fmt.Sprintf("Full Content:\n\n%s", visibleContent)
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
Times:			   						  	%d
Words:										%d
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
