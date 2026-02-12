package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"ymusic/internal/api"
	"ymusic/internal/theme"
)

type TrackListModel struct {
	tracks       []api.Track
	cursor       int
	offset       int
	height       int
	width        int
	playingID    string
	focused      bool
	likedSet     map[string]bool
}

func NewTrackList() TrackListModel {
	return TrackListModel{
		likedSet: make(map[string]bool),
	}
}

func (m *TrackListModel) SetTracks(tracks []api.Track) {
	m.tracks = tracks
	m.cursor = 0
	m.offset = 0
}

func (m *TrackListModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m *TrackListModel) SetPlaying(id string) {
	m.playingID = id
}

func (m *TrackListModel) SetFocused(f bool) {
	m.focused = f
}

func (m *TrackListModel) SetLiked(liked map[string]bool) {
	m.likedSet = liked
}

func (m *TrackListModel) ToggleLike(id string) {
	if m.likedSet[id] {
		delete(m.likedSet, id)
	} else {
		m.likedSet[id] = true
	}
}

func (m *TrackListModel) MoveUp() {
	if m.cursor > 0 {
		m.cursor--
		if m.cursor < m.offset {
			m.offset = m.cursor
		}
	}
}

func (m *TrackListModel) MoveDown() {
	if m.cursor < len(m.tracks)-1 {
		m.cursor++
		visible := m.height - 1 // header
		if m.cursor >= m.offset+visible {
			m.offset = m.cursor - visible + 1
		}
	}
}

func (m TrackListModel) Selected() *api.Track {
	if len(m.tracks) == 0 || m.cursor < 0 || m.cursor >= len(m.tracks) {
		return nil
	}
	return &m.tracks[m.cursor]
}

func (m TrackListModel) Cursor() int { return m.cursor }
func (m TrackListModel) Tracks() []api.Track { return m.tracks }

// GoToAlbumCmd returns a command to navigate to the selected track's album.
func (m TrackListModel) GoToAlbumCmd() tea.Cmd {
	t := m.Selected()
	if t == nil || len(t.Albums) == 0 {
		return nil
	}
	albumID := t.Albums[0].ID
	return func() tea.Msg {
		return navigateAlbumMsg{id: albumID}
	}
}

func (m TrackListModel) View() string {
	if len(m.tracks) == 0 {
		return theme.S.Muted.Render("  No tracks")
	}

	var b strings.Builder

	// Header
	header := fmt.Sprintf("  %-4s %-*s %-20s %6s",
		"#", m.width-38, "Title", "Artist", "Time")
	b.WriteString(theme.S.Muted.Render(header))
	b.WriteString("\n")

	visible := m.height - 1
	if visible < 1 {
		visible = 1
	}
	end := m.offset + visible
	if end > len(m.tracks) {
		end = len(m.tracks)
	}

	for i := m.offset; i < end; i++ {
		t := m.tracks[i]
		num := fmt.Sprintf("%d", i+1)

		likeIcon := " "
		if m.likedSet[t.ID] {
			likeIcon = "♥"
		}

		playIcon := " "
		if t.ID == m.playingID {
			playIcon = "▶"
		}

		dur := formatTime(t.DurationSec())
		titleWidth := m.width - 38
		if titleWidth < 10 {
			titleWidth = 10
		}

		title := truncate(t.Title, titleWidth)
		artist := truncate(t.ArtistName(), 20)

		line := fmt.Sprintf("%s%s%-4s %-*s %-20s %6s",
			playIcon, likeIcon, num, titleWidth, title, artist, dur)

		if i == m.cursor && m.focused {
			b.WriteString(theme.S.ListActive.Render(line))
		} else if t.ID == m.playingID {
			b.WriteString(theme.S.ListPlaying.Render(line))
		} else {
			b.WriteString(theme.S.ListItem.Render(line))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (m *TrackListModel) ScrollUp(n int) {
	for i := 0; i < n; i++ {
		m.MoveUp()
	}
}

func (m *TrackListModel) ScrollDown(n int) {
	for i := 0; i < n; i++ {
		m.MoveDown()
	}
}

// HandleMouse processes mouse events. offsetY is the Y position of the tracklist
// within the parent view. Returns whether the event was handled and an optional command.
func (m *TrackListModel) HandleMouse(msg tea.MouseMsg, offsetY int) (bool, tea.Cmd) {
	if len(m.tracks) == 0 {
		return false, nil
	}

	switch msg.Button {
	case tea.MouseButtonWheelUp:
		m.ScrollUp(3)
		return true, nil
	case tea.MouseButtonWheelDown:
		m.ScrollDown(3)
		return true, nil
	case tea.MouseButtonLeft:
		if msg.Action != tea.MouseActionPress {
			return false, nil
		}
		row := msg.Y - offsetY
		if row < 1 {
			return false, nil
		}
		trackIdx := m.offset + row - 1
		if trackIdx < 0 || trackIdx >= len(m.tracks) {
			return false, nil
		}
		m.cursor = trackIdx
		t := m.tracks[trackIdx]
		tracks := m.tracks
		idx := trackIdx
		return true, func() tea.Msg {
			return PlayTrackMsg{Track: t, Queue: tracks, Index: idx}
		}
	case tea.MouseButtonRight:
		if msg.Action != tea.MouseActionPress {
			return false, nil
		}
		row := msg.Y - offsetY
		if row < 1 {
			return false, nil
		}
		trackIdx := m.offset + row - 1
		if trackIdx < 0 || trackIdx >= len(m.tracks) {
			return false, nil
		}
		m.cursor = trackIdx
		t := m.tracks[trackIdx]
		if len(t.Albums) > 0 {
			albumID := t.Albums[0].ID
			return true, func() tea.Msg {
				return navigateAlbumMsg{id: albumID}
			}
		}
		return true, nil
	}
	return false, nil
}
