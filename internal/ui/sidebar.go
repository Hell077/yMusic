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

	if msg, ok := msg.(tea.KeyMsg); ok {
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
