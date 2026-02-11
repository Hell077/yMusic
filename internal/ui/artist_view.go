package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"ymusic/internal/api"
	"ymusic/internal/theme"
)

type ArtistViewModel struct {
	info      *api.ArtistBriefInfo
	trackList TrackListModel
	cursor    int
	section   int // 0=tracks, 1=albums, 2=similar
	loading   bool
	err       error
	width     int
	height    int
	focused   bool
}

func NewArtistView() ArtistViewModel {
	return ArtistViewModel{
		trackList: NewTrackList(),
		loading:   true,
	}
}

func (m ArtistViewModel) Init() tea.Cmd { return nil }

func (m ArtistViewModel) Update(msg tea.Msg) (ArtistViewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case ArtistInfoMsg:
		m.info = msg.Info
		m.loading = false
		if msg.Info != nil {
			m.trackList.SetTracks(msg.Info.PopularTracks)
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
			if m.section == 0 {
				m.trackList.MoveUp()
			} else if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.section == 0 {
				m.trackList.MoveDown()
			} else {
				max := m.sectionLen() - 1
				if max >= 0 && m.cursor < max {
					m.cursor++
				}
			}
		case "right", "l":
			m.section = (m.section + 1) % 3
			m.cursor = 0
		case "left", "h":
			m.section = (m.section + 2) % 3
			m.cursor = 0
		case "enter":
			return m, m.handleEnter()
		}
	}
	return m, nil
}

func (m ArtistViewModel) handleEnter() tea.Cmd {
	if m.info == nil {
		return nil
	}
	switch m.section {
	case 0:
		t := m.trackList.Selected()
		if t != nil {
			tracks := m.info.PopularTracks
			idx := m.trackList.Cursor()
			return func() tea.Msg {
				return PlayTrackMsg{Track: *t, Queue: tracks, Index: idx}
			}
		}
	case 1:
		albums := append(m.info.Albums, m.info.AlsoAlbums...)
		if m.cursor < len(albums) {
			album := albums[m.cursor]
			return func() tea.Msg {
				return navigateAlbumMsg{id: album.ID}
			}
		}
	case 2:
		if m.cursor < len(m.info.SimilarArtists) {
			artist := m.info.SimilarArtists[m.cursor]
			return func() tea.Msg {
				return navigateArtistMsg{id: artist.ID}
			}
		}
	}
	return nil
}

func (m ArtistViewModel) sectionLen() int {
	if m.info == nil {
		return 0
	}
	switch m.section {
	case 1:
		return len(m.info.Albums) + len(m.info.AlsoAlbums)
	case 2:
		return len(m.info.SimilarArtists)
	}
	return 0
}

func (m ArtistViewModel) View() string {
	var b strings.Builder

	if m.loading {
		b.WriteString(theme.S.Muted.Render("  Loading artist...") + "\n")
		return b.String()
	}
	if m.err != nil {
		b.WriteString(theme.S.Error.Render("  Error: "+m.err.Error()) + "\n")
		return b.String()
	}
	if m.info == nil {
		return ""
	}

	b.WriteString(theme.S.Title.Render(" "+m.info.Artist.Name) + "\n\n")

	sections := []string{"Popular Tracks", "Albums", "Similar Artists"}
	var tabs []string
	for i, name := range sections {
		if i == m.section {
			tabs = append(tabs, theme.S.ActiveTab.Render(name))
		} else {
			tabs = append(tabs, theme.S.Tab.Render(name))
		}
	}
	b.WriteString("  " + strings.Join(tabs, " ") + "\n\n")

	switch m.section {
	case 0:
		b.WriteString(m.trackList.View())
	case 1:
		albums := append(m.info.Albums, m.info.AlsoAlbums...)
		if len(albums) == 0 {
			b.WriteString(theme.S.Muted.Render("  No albums") + "\n")
		} else {
			for i, a := range albums {
				line := fmt.Sprintf("  %s (%d)", a.Title, a.Year)
				if i == m.cursor && m.focused {
					b.WriteString(theme.S.ListActive.Render("▸"+line) + "\n")
				} else {
					b.WriteString(theme.S.ListItem.Render(" "+line) + "\n")
				}
			}
		}
	case 2:
		if len(m.info.SimilarArtists) == 0 {
			b.WriteString(theme.S.Muted.Render("  No similar artists") + "\n")
		} else {
			for i, a := range m.info.SimilarArtists {
				line := fmt.Sprintf("  %s", a.Name)
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

func (m *ArtistViewModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.trackList.SetSize(w, h-6)
}

func (m *ArtistViewModel) SetFocused(f bool) {
	m.focused = f
	m.trackList.SetFocused(f && m.section == 0)
}
