package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"ymusic/internal/theme"
)

type SidebarItem struct {
	Icon  string
	Label string
	Page  Page
}

var sidebarItems = []SidebarItem{
	{Icon: "♫", Label: "Home", Page: PageHome},
	{Icon: "⌕", Label: "Search", Page: PageSearch},
	{Icon: "♥", Label: "Collection", Page: PageCollection},
	{Icon: "≈", Label: "My Wave", Page: PageMyWave},
	{Icon: "≡", Label: "Queue", Page: PageQueue},
}

type SidebarModel struct {
	cursor int
	width  int
	height int
	focused bool
}

func NewSidebar() SidebarModel {
	return SidebarModel{focused: true}
}

func (m SidebarModel) Init() tea.Cmd { return nil }

func (m SidebarModel) Update(msg tea.Msg) (SidebarModel, tea.Cmd) {
	if !m.focused {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(sidebarItems)-1 {
				m.cursor++
			}
		case "enter":
			return m, func() tea.Msg {
				return NavigateMsg{Page: sidebarItems[m.cursor].Page}
			}
		}
	}
	return m, nil
}

// HandleMouse processes mouse events. Returns whether handled and an optional command.
func (m *SidebarModel) HandleMouse(msg tea.MouseMsg) (bool, tea.Cmd) {
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		if m.cursor > 0 {
			m.cursor--
		}
		return true, nil
	case tea.MouseButtonWheelDown:
		if m.cursor < len(sidebarItems)-1 {
			m.cursor++
		}
		return true, nil
	case tea.MouseButtonLeft:
		if msg.Action != tea.MouseActionPress {
			return false, nil
		}
		// Layout: title=line 0, blank=line 1, items start at line 2
		idx := msg.Y - 2
		if idx < 0 || idx >= len(sidebarItems) {
			return false, nil
		}
		m.cursor = idx
		page := sidebarItems[idx].Page
		return true, func() tea.Msg {
			return NavigateMsg{Page: page}
		}
	}
	return false, nil
}

func (m SidebarModel) View() string {
	var b strings.Builder
	b.WriteString(theme.S.Title.Render("  ymusic"))
	b.WriteString("\n\n")

	for i, item := range sidebarItems {
		label := item.Icon + " " + item.Label
		if i == m.cursor && m.focused {
			b.WriteString(theme.S.SidebarActive.Render("▸ "+label) + "\n")
		} else {
			b.WriteString(theme.S.SidebarItem.Render("  "+label) + "\n")
		}
	}

	// Pad remaining height
	lines := len(sidebarItems) + 3 // title + blank + items
	for i := lines; i < m.height; i++ {
		b.WriteString("\n")
	}

	return b.String()
}

func (m SidebarModel) SelectedPage() Page {
	return sidebarItems[m.cursor].Page
}

func (m *SidebarModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m *SidebarModel) SetFocused(f bool) {
	m.focused = f
}

func (m *SidebarModel) SetCursorForPage(p Page) {
	for i, item := range sidebarItems {
		if item.Page == p {
			m.cursor = i
			return
		}
	}
}
