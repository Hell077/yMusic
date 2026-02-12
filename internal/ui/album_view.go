package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"ymusic/internal/api"
	"ymusic/internal/theme"
)

type AlbumViewModel struct {
	album     *api.Album
	trackList TrackListModel
	loading   bool
	err       error
	width     int
	height    int
	focused   bool
}

func NewAlbumView() AlbumViewModel {
	return AlbumViewModel{
		trackList: NewTrackList(),
		loading:   true,
	}
}

func (m AlbumViewModel) Init() tea.Cmd { return nil }

func (m AlbumViewModel) Update(msg tea.Msg) (AlbumViewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case AlbumMsg:
		m.album = msg.Album
		m.loading = false
		if msg.Album != nil {
			var tracks []api.Track
			for _, vol := range msg.Album.Volumes {
				tracks = append(tracks, vol...)
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
		}
	}
	return m, nil
}

func (m AlbumViewModel) View() string {
	var b strings.Builder

	if m.loading {
		b.WriteString(theme.S.Muted.Render("  Loading album...") + "\n")
		return b.String()
	}
	if m.err != nil {
		b.WriteString(theme.S.Error.Render("  Error: "+m.err.Error()) + "\n")
		return b.String()
	}
	if m.album == nil {
		return ""
	}

	b.WriteString(theme.S.Title.Render(" "+m.album.Title) + "\n")
	info := fmt.Sprintf("  %s · %d · %s · %d tracks",
		m.album.ArtistName(), m.album.Year, m.album.Genre, m.album.TrackCount)
	b.WriteString(theme.S.Muted.Render(info) + "\n\n")

	b.WriteString(m.trackList.View())
	return b.String()
}

func (m *AlbumViewModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.trackList.SetSize(w, h-4)
}

func (m *AlbumViewModel) SetFocused(f bool) {
	m.focused = f
	m.trackList.SetFocused(f)
}
