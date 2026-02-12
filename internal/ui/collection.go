package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"ymusic/internal/api"
	"ymusic/internal/theme"
)

type CollectionTab int

const (
	CollTabLiked CollectionTab = iota
	CollTabPlaylists
	CollTabAlbums
)

var collTabNames = []string{"Liked", "Playlists", "Albums"}

type CollectionModel struct {
	tab         CollectionTab
	likedTracks []api.Track
	playlists   []api.Playlist
	albums      []api.Album
	trackList   TrackListModel
	cursor      int
	loading     bool
	err         error
	width       int
	height      int
	focused     bool
}

func NewCollection() CollectionModel {
	return CollectionModel{
		trackList: NewTrackList(),
		loading:   true,
	}
}

func (m CollectionModel) Init() tea.Cmd { return nil }

func (m CollectionModel) Update(msg tea.Msg) (CollectionModel, tea.Cmd) {
	switch msg := msg.(type) {
	case LikedTracksMsg:
		m.likedTracks = msg.Tracks
		m.trackList.SetTracks(msg.Tracks)
		m.loading = false
	case UserPlaylistsMsg:
		m.playlists = msg.Playlists
		m.loading = false
	case LikedAlbumsMsg:
		m.albums = msg.Albums
	case ErrorMsg:
		m.err = msg.Err
		m.loading = false
	case tea.MouseMsg:
		// Header: title(0) + blank(1) + tabs(2) + blank(3), content starts at row 4
		if m.tab == CollTabLiked {
			if handled, cmd := m.trackList.HandleMouse(msg, 4); handled {
				return m, cmd
			}
		} else {
			switch msg.Button {
			case tea.MouseButtonWheelUp:
				if m.cursor > 0 {
					m.cursor--
				}
				return m, nil
			case tea.MouseButtonWheelDown:
				max := m.currentListLen() - 1
				if max >= 0 && m.cursor < max {
					m.cursor++
				}
				return m, nil
			case tea.MouseButtonLeft:
				if msg.Action != tea.MouseActionPress {
					return m, nil
				}
				idx := msg.Y - 4
				if idx >= 0 && idx < m.currentListLen() {
					m.cursor = idx
					return m, m.handleEnter()
				}
			}
		}
	case tea.KeyMsg:
		if !m.focused {
			return m, nil
		}
		switch msg.String() {
		case "left", "h":
			if m.tab > 0 {
				m.tab--
				m.cursor = 0
			}
		case "right", "l":
			if m.tab < CollectionTab(len(collTabNames))-1 {
				m.tab++
				m.cursor = 0
			}
		case "up", "k":
			if m.tab == CollTabLiked {
				m.trackList.MoveUp()
			} else if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.tab == CollTabLiked {
				m.trackList.MoveDown()
			} else {
				max := m.currentListLen() - 1
				if max < 0 {
					max = 0
				}
				if m.cursor < max {
					m.cursor++
				}
			}
		case "enter":
			return m, m.handleEnter()
		case "a":
			if m.tab == CollTabLiked {
				if cmd := m.trackList.GoToAlbumCmd(); cmd != nil {
					return m, cmd
				}
			}
		}
	}
	return m, nil
}

func (m CollectionModel) handleEnter() tea.Cmd {
	switch m.tab {
	case CollTabLiked:
		t := m.trackList.Selected()
		if t != nil {
			tracks := m.likedTracks
			idx := m.trackList.Cursor()
			return func() tea.Msg {
				return PlayTrackMsg{Track: *t, Queue: tracks, Index: idx}
			}
		}
	case CollTabPlaylists:
		if m.cursor < len(m.playlists) {
			p := m.playlists[m.cursor]
			return func() tea.Msg {
				return navigatePlaylistMsg{uid: p.UID, kind: p.Kind}
			}
		}
	case CollTabAlbums:
		if m.cursor < len(m.albums) {
			a := m.albums[m.cursor]
			return func() tea.Msg {
				return navigateAlbumMsg{id: a.ID}
			}
		}
	}
	return nil
}

func (m CollectionModel) currentListLen() int {
	switch m.tab {
	case CollTabPlaylists:
		return len(m.playlists)
	case CollTabAlbums:
		return len(m.albums)
	}
	return 0
}

func (m CollectionModel) View() string {
	var b strings.Builder

	b.WriteString(theme.S.Title.Render(" Collection") + "\n\n")

	// Tabs
	var tabs []string
	for i, name := range collTabNames {
		if CollectionTab(i) == m.tab {
			tabs = append(tabs, theme.S.ActiveTab.Render(name))
		} else {
			tabs = append(tabs, theme.S.Tab.Render(name))
		}
	}
	b.WriteString("  " + strings.Join(tabs, " ") + "\n\n")

	if m.loading {
		b.WriteString(theme.S.Muted.Render("  Loading...") + "\n")
		return b.String()
	}

	switch m.tab {
	case CollTabLiked:
		b.WriteString(m.trackList.View())
	case CollTabPlaylists:
		if len(m.playlists) == 0 {
			b.WriteString(theme.S.Muted.Render("  No playlists") + "\n")
		} else {
			for i, p := range m.playlists {
				line := fmt.Sprintf("  %s (%d tracks)", p.Title, p.TrackCount)
				if i == m.cursor && m.focused {
					b.WriteString(theme.S.ListActive.Render("▸"+line) + "\n")
				} else {
					b.WriteString(theme.S.ListItem.Render(" "+line) + "\n")
				}
			}
		}
	case CollTabAlbums:
		if len(m.albums) == 0 {
			b.WriteString(theme.S.Muted.Render("  No albums") + "\n")
		} else {
			for i, a := range m.albums {
				line := fmt.Sprintf("  %s - %s (%d)", a.Title, a.ArtistName(), a.Year)
				if i == m.cursor && m.focused {
					b.WriteString(theme.S.ListActive.Render("▸"+line) + "\n")
				} else {
					b.WriteString(theme.S.ListItem.Render(" "+line) + "\n")
				}
			}
		}
	}

	return b.String()
}

func (m *CollectionModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.trackList.SetSize(w, h-6)
}

func (m *CollectionModel) SetFocused(f bool) {
	m.focused = f
	m.trackList.SetFocused(f && m.tab == CollTabLiked)
}
