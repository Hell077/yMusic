package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"ymusic/internal/api"
	"ymusic/internal/theme"
)

type PlaylistViewModel struct {
	playlist  *api.Playlist
	trackList TrackListModel
	loading   bool
	err       error
	width     int
	height    int
	focused   bool
}

func NewPlaylistView() PlaylistViewModel {
	return PlaylistViewModel{
		trackList: NewTrackList(),
		loading:   true,
	}
}

func (m PlaylistViewModel) Init() tea.Cmd { return nil }

func (m PlaylistViewModel) Update(msg tea.Msg) (PlaylistViewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case PlaylistMsg:
		m.playlist = msg.Playlist
		m.loading = false
		if msg.Playlist != nil {
			var tracks []api.Track
			for _, ti := range msg.Playlist.Tracks {
				tracks = append(tracks, ti.Track)
			}
			m.trackList.SetTracks(tracks)
		}
	case ErrorMsg:
		m.err = msg.Err
		m.loading = false
	case tea.MouseMsg:
		// Header: title(0) + info(1) + blank(2), tracklist starts at row 3
		if handled, cmd := m.trackList.HandleMouse(msg, 3); handled {
			return m, cmd
		}
	case tea.KeyMsg:
		if !m.focused {
			return m, nil
		}
		switch msg.String() {
		case "up", "k":
			m.trackList.MoveUp()
		case "down", "j":
			m.trackList.MoveDown()
		case "enter":
			t := m.trackList.Selected()
			if t != nil {
				tracks := m.trackList.Tracks()
				idx := m.trackList.Cursor()
				return m, func() tea.Msg {
					return PlayTrackMsg{Track: *t, Queue: tracks, Index: idx}
				}
			}
		case "a":
			if cmd := m.trackList.GoToAlbumCmd(); cmd != nil {
				return m, cmd
			}
		}
	}
	return m, nil
}

func (m PlaylistViewModel) View() string {
	var b strings.Builder

	if m.loading {
		b.WriteString(theme.S.Muted.Render("  Loading playlist...") + "\n")
		return b.String()
	}
	if m.err != nil {
		b.WriteString(theme.S.Error.Render("  Error: "+m.err.Error()) + "\n")
		return b.String()
	}
	if m.playlist == nil {
		return ""
	}

	b.WriteString(theme.S.Title.Render(" "+m.playlist.Title) + "\n")
	b.WriteString(theme.S.Muted.Render(fmt.Sprintf("  %s · %d tracks",
		m.playlist.Owner.Name, m.playlist.TrackCount)) + "\n\n")

	b.WriteString(m.trackList.View())
	return b.String()
}

func (m *PlaylistViewModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.trackList.SetSize(w, h-4)
}

func (m *PlaylistViewModel) SetFocused(f bool) {
	m.focused = f
	m.trackList.SetFocused(f)
}
