package tui

// import (
// 	"strings"

// 	"github.com/charmbracelet/lipgloss"
// )

// var (
// 	appNameStyle        = lipgloss.NewStyle().Background(lipgloss.Color("99")).Padding(0, 1)
// 	faint               = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Faint(true)
// 	listEnumeratorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)
// )

// func (m model) View() string {
// 	s := appNameStyle.Render("NOTES APP") + "\n\n"

// 	if m.state == titleView {
// 		s += "Note title:\n\n"
// 		s += m.textinput.View() + "\n\n"
// 		s += faint.Render("enter - save • esc - discard")
// 	}

// 	if m.state == bodyView {
// 		s += "Note:\n\n"
// 		s += m.textarea.View() + "\n\n"
// 		s += faint.Render("ctrl+s - save • esc - discard")
// 	}

// 	if m.state == listView {
// 		for i, n := range m.notes {
// 			prefix := " "
// 			if i == m.listIndex {
// 				prefix = ">"
// 			}
// 			shortBody := strings.ReplaceAll(n.Body, "\n", " ")
// 			if len(shortBody) > 30 {
// 				shortBody = shortBody[:30]
// 			}
// 			s += listEnumeratorStyle.Render(prefix) + n.Title + " | " + faint.Render(shortBody) + "\n\n"
// 		}
// 		s += faint.Render("n - new note • q - quit")
// 	}

// 	return s
// }


// ============ former update ===============================================

	// var (
	// 	cmds []tea.Cmd
	// 	cmd  tea.Cmd
	// )

	// m.textarea, cmd = m.textarea.Update(msg)
	// cmds = append(cmds, cmd)

	// m.textinput, cmd = m.textinput.Update(msg)
	// cmds = append(cmds, cmd)

	// switch msg := msg.(type) {
	// // handle key strokes
	// case tea.KeyMsg:
	// 	key := msg.String()
	// 	switch m.state {
	// 	// List View key bindings
	// 	case listView:
	// 		switch key {
	// 		case "q":
	// 			return m, tea.Quit
	// 		case "n":
	// 			m.textinput.SetValue("")
	// 			m.textinput.Focus()
	// 			m.currNote = Note{}
	// 			m.state = titleView
	// 		case "up", "k":
	// 			if m.listIndex > 0 {
	// 				m.listIndex--
	// 			}
	// 		case "down", "j":
	// 			if m.listIndex < len(m.notes)-1 {
	// 				m.listIndex++
	// 			}
	// 		case "enter":
	// 			m.currNote = m.notes[m.listIndex]
	// 			m.state = bodyView
	// 			m.textarea.SetValue(m.currNote.Body)
	// 			m.textarea.Focus()
	// 			m.textarea.CursorEnd()
	// 		}

	// 	// Title Input View key bindings
	// 	case titleView:
	// 		switch key {
	// 		case "enter":
	// 			title := m.textinput.Value()
	// 			if title != "" {
	// 				m.currNote.Title = title

	// 				m.state = bodyView
	// 				m.textarea.SetValue("")
	// 				m.textarea.Focus()
	// 				m.textarea.CursorEnd()
	// 			}
	// 		case "esc":
	// 			m.state = listView
	// 		}

	// 	// Body Textarea key bindings
	// 	case bodyView:
	// 		switch key {
	// 		case "ctrl+s":
	// 			m.currNote.Body = m.textarea.Value()

	// 			var err error
	// 			if err = m.store.SaveNote(m.currNote); err != nil {
	// 				// TODO: handle error instead of quitting
	// 				return m, tea.Quit
	// 			}

	// 			m.notes, err = m.store.GetNotes()
	// 			if err != nil {
	// 				// TODO: handle error instead of quitting
	// 				return m, tea.Quit
	// 			}

	// 			m.state = listView
	// 		case "esc":
	// 			m.state = listView
	// 		}
	// 	}
	// }

	// return m, tea.Batch(cmds...)