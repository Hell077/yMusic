package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"ymusic/internal/api"
	"ymusic/internal/theme"
)

type HomeModel struct {
	feed       *api.FeedResponse
	playlists  []api.GeneratedPlaylist
	cursor     int
	loading    bool
	err        error
	width      int
	height     int
	focused    bool
}

func NewHome() HomeModel {
	return HomeModel{loading: true}
}

func (m HomeModel) Init() tea.Cmd { return nil }

func (m HomeModel) Update(msg tea.Msg) (HomeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case FeedMsg:
		m.feed = msg.Feed
		m.loading = false
		if msg.Feed != nil {
			m.playlists = msg.Feed.GeneratedPlaylists
		}
	case ErrorMsg:
		m.err = msg.Err
		m.loading = false
	case tea.KeyMsg:
		if !m.focused {
			return m, nil
		}
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			max := len(m.playlists) - 1
			if max < 0 {
				max = 0
			}
			if m.cursor < max {
				m.cursor++
			}
		case "enter":
			if len(m.playlists) > 0 && m.cursor < len(m.playlists) {
				p := m.playlists[m.cursor].Data
				return m, func() tea.Msg {
					return navigatePlaylistMsg{uid: p.UID, kind: p.Kind}
				}
			}
		}
	}
	return m, nil
}

func (m HomeModel) View() string {
	var b strings.Builder

	b.WriteString(theme.S.Title.Render(" Home") + "\n\n")

	if m.loading {
		b.WriteString(theme.S.Muted.Render("  Loading...") + "\n")
		return b.String()
	}

	if m.err != nil {
		b.WriteString(theme.S.Error.Render("  Error: "+m.err.Error()) + "\n")
		return b.String()
	}

	if len(m.playlists) == 0 {
		b.WriteString(theme.S.Muted.Render("  No recommendations available") + "\n")
		return b.String()
	}

	b.WriteString(theme.S.Subtitle.Render("  Playlists for you") + "\n\n")

	for i, gp := range m.playlists {
		p := gp.Data
		label := fmt.Sprintf("  %s (%d tracks)", p.Title, p.TrackCount)
		if i == m.cursor && m.focused {
			b.WriteString(theme.S.ListActive.Render("▸ "+label) + "\n")
		} else {
			b.WriteString(theme.S.ListItem.Render("  "+label) + "\n")
		}
	}

	return b.String()
}

func (m *HomeModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m *HomeModel) SetFocused(f bool) {
	m.focused = f
}
