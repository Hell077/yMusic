package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"ymusic/internal/theme"
)

type OverlayView int

const (
	OverlayMain OverlayView = iota
	OverlayThemes
	OverlayHelp
)

type OverlayItem struct {
	Label  string
	Action string
}

var mainMenuItems = []OverlayItem{
	{Label: "Themes", Action: "themes"},
	{Label: "Help", Action: "help"},
	{Label: "Quit", Action: "quit"},
}

type OverlayModel struct {
	visible    bool
	view       OverlayView
	cursor     int
	width      int
	height     int
}

func NewOverlay() OverlayModel {
	return OverlayModel{}
}

func (m *OverlayModel) Toggle() {
	m.visible = !m.visible
	if m.visible {
		m.view = OverlayMain
		m.cursor = 0
	}
}

func (m *OverlayModel) Close() {
	m.visible = false
	m.view = OverlayMain
	m.cursor = 0
}

func (m OverlayModel) Visible() bool {
	return m.visible
}

func (m OverlayModel) Update(msg tea.Msg) (OverlayModel, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "esc":
			if m.view != OverlayMain {
				m.view = OverlayMain
				m.cursor = 0
			} else {
				m.Close()
			}
			return m, nil
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			max := m.currentLen() - 1
			if max >= 0 && m.cursor < max {
				m.cursor++
			}
		case "enter":
			return m, m.handleEnter()
		case "q":
			m.Close()
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *OverlayModel) handleEnter() tea.Cmd {
	switch m.view {
	case OverlayMain:
		if m.cursor < len(mainMenuItems) {
			action := mainMenuItems[m.cursor].Action
			switch action {
			case "themes":
				m.view = OverlayThemes
				m.cursor = 0
			case "help":
				m.view = OverlayHelp
				m.cursor = 0
			case "quit":
				m.Close()
				return tea.Quit
			}
		}
	case OverlayThemes:
		if m.cursor < len(theme.Themes) {
			theme.SetThemeByIndex(m.cursor)
			m.Close()
			return func() tea.Msg { return ThemeChangedMsg{} }
		}
	}
	return nil
}

func (m OverlayModel) currentLen() int {
	switch m.view {
	case OverlayMain:
		return len(mainMenuItems)
	case OverlayThemes:
		return len(theme.Themes)
	case OverlayHelp:
		return 0
	}
	return 0
}

func (m OverlayModel) View() string {
	if !m.visible {
		return ""
	}

	var b strings.Builder

	switch m.view {
	case OverlayMain:
		b.WriteString(theme.S.Title.Render("Menu") + "\n\n")
		for i, item := range mainMenuItems {
			if i == m.cursor {
				b.WriteString(theme.S.OverlayActive.Render("▸ "+item.Label) + "\n")
			} else {
				b.WriteString(theme.S.OverlayItem.Render("  "+item.Label) + "\n")
			}
		}
	case OverlayThemes:
		b.WriteString(theme.S.Title.Render("Themes") + "\n\n")
		for i, t := range theme.Themes {
			name := t.Name
			if t.Name == theme.Current.Name {
				name += " ●"
			}
			if i == m.cursor {
				b.WriteString(theme.S.OverlayActive.Render("▸ "+name) + "\n")
			} else {
				b.WriteString(theme.S.OverlayItem.Render("  "+name) + "\n")
			}
		}
	case OverlayHelp:
		b.WriteString(theme.S.Title.Render("Keyboard Shortcuts") + "\n\n")
		b.WriteString(renderHelp("space", "Play/Pause"))
		b.WriteString(renderHelp("n/p", "Next/Previous track"))
		b.WriteString(renderHelp("+/-", "Volume up/down"))
		b.WriteString(renderHelp(">/< ", "Seek forward/back 10s"))
		b.WriteString(renderHelp("↑↓/jk", "Navigate"))
		b.WriteString(renderHelp("←→/hl", "Switch tabs"))
		b.WriteString(renderHelp("enter", "Select/Play"))
		b.WriteString(renderHelp("tab", "Switch focus"))
		b.WriteString(renderHelp("/", "Search"))
		b.WriteString(renderHelp("L", "Like track"))
		b.WriteString(renderHelp("s", "Toggle shuffle"))
		b.WriteString(renderHelp("r", "Cycle repeat"))
		b.WriteString(renderHelp("esc", "Menu / Back"))
		b.WriteString(renderHelp("q", "Quit"))
	}

	content := b.String()

	// Center overlay
	boxWidth := 40
	boxHeight := strings.Count(content, "\n") + 2
	style := theme.S.Overlay.
		Width(boxWidth).
		Height(boxHeight)

	box := style.Render(content)

	// Position in center
	x := (m.width - boxWidth - 4) / 2
	y := (m.height - boxHeight - 4) / 2
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

func renderHelp(key, desc string) string {
	return "  " + theme.S.HelpKey.Render(key) + " " +
		theme.S.HelpDesc.Render(desc) + "\n"
}

func (m *OverlayModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}
