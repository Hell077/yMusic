package ui

import (
	"fmt"
	"strings"

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
			b.WriteString(theme.S.Primary.Render(line))
		} else {
			b.WriteString(theme.S.ListItem.Render(line))
		}
		b.WriteString("\n")
	}

	return b.String()
}
