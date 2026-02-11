package ui

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"ymusic/internal/api"
	"ymusic/internal/theme"
)

type SearchTab int

const (
	SearchTabAll SearchTab = iota
	SearchTabTracks
	SearchTabAlbums
	SearchTabArtists
)

var searchTabNames = []string{"All", "Tracks", "Albums", "Artists"}

type SearchModel struct {
	input       textinput.Model
	result      *api.SearchResult
	tracks      []api.Track
	albums      []api.Album
	artists     []api.Artist
	tab         SearchTab
	cursor      int
	inputFocused bool
	loading     bool
	err         error
	width       int
	height      int
	focused     bool
}

func NewSearch() SearchModel {
	ti := textinput.New()
	ti.Placeholder = "Search music..."
	ti.CharLimit = 100
	ti.Width = 40
	return SearchModel{
		input:        ti,
		inputFocused: true,
	}
}

func (m SearchModel) Init() tea.Cmd { return nil }

func (m SearchModel) Update(msg tea.Msg) (SearchModel, tea.Cmd) {
	switch msg := msg.(type) {
	case SearchResultMsg:
		m.result = msg.Result
		m.loading = false
		m.cursor = 0
		if msg.Result != nil {
			m.tracks, _ = parseSearchTracks(msg.Result)
			m.albums, _ = parseSearchAlbums(msg.Result)
			m.artists, _ = parseSearchArtists(msg.Result)
		}
	case ErrorMsg:
		m.err = msg.Err
		m.loading = false
	case tea.KeyMsg:
		if !m.focused {
			return m, nil
		}
		if m.inputFocused {
			switch msg.String() {
			case "enter":
				query := m.input.Value()
				if query != "" {
					m.loading = true
					m.inputFocused = false
					return m, func() tea.Msg {
						return doSearchMsg{query: query}
					}
				}
			case "esc":
				m.inputFocused = false
				return m, nil
			}
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "/":
			m.inputFocused = true
			m.input.Focus()
			return m, nil
		case "right", "l":
			if m.tab < SearchTab(len(searchTabNames))-1 {
				m.tab++
				m.cursor = 0
			}
		case "left", "h":
			if m.tab > 0 {
				m.tab--
				m.cursor = 0
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			max := m.currentListLen() - 1
			if max < 0 {
				max = 0
			}
			if m.cursor < max {
				m.cursor++
			}
		case "enter":
			return m, m.handleEnter()
		}
	}
	return m, nil
}

type doSearchMsg struct{ query string }

func (m SearchModel) handleEnter() tea.Cmd {
	switch m.tab {
	case SearchTabAll, SearchTabTracks:
		if m.cursor < len(m.tracks) {
			tracks := m.tracks
			idx := m.cursor
			return func() tea.Msg {
				return PlayTrackMsg{Track: tracks[idx], Queue: tracks, Index: idx}
			}
		}
	case SearchTabAlbums:
		if m.cursor < len(m.albums) {
			album := m.albums[m.cursor]
			return func() tea.Msg {
				return navigateAlbumMsg{id: album.ID}
			}
		}
	case SearchTabArtists:
		if m.cursor < len(m.artists) {
			artist := m.artists[m.cursor]
			return func() tea.Msg {
				return navigateArtistMsg{id: artist.ID}
			}
		}
	}
	return nil
}

type navigateAlbumMsg struct{ id int }
type navigateArtistMsg struct{ id int }
type navigatePlaylistMsg struct {
	uid  int
	kind int
}

func (m SearchModel) currentListLen() int {
	switch m.tab {
	case SearchTabAll, SearchTabTracks:
		return len(m.tracks)
	case SearchTabAlbums:
		return len(m.albums)
	case SearchTabArtists:
		return len(m.artists)
	}
	return 0
}

func (m SearchModel) View() string {
	var b strings.Builder

	b.WriteString(theme.S.Title.Render(" Search") + "\n\n")

	if m.inputFocused {
		b.WriteString("  " + theme.S.SearchInput.Render(m.input.View()) + "\n\n")
	} else {
		query := m.input.Value()
		if query != "" {
			b.WriteString(theme.S.Muted.Render("  Search: ") + theme.S.Primary.Render(query) +
				theme.S.Muted.Render("  (press / to edit)") + "\n\n")
		} else {
			b.WriteString(theme.S.Muted.Render("  Press / to search") + "\n\n")
		}
	}

	if m.loading {
		b.WriteString(theme.S.Muted.Render("  Searching...") + "\n")
		return b.String()
	}

	if m.result == nil {
		return b.String()
	}

	// Tabs
	var tabs []string
	for i, name := range searchTabNames {
		if SearchTab(i) == m.tab {
			tabs = append(tabs, theme.S.ActiveTab.Render(name))
		} else {
			tabs = append(tabs, theme.S.Tab.Render(name))
		}
	}
	b.WriteString("  " + strings.Join(tabs, " ") + "\n\n")

	switch m.tab {
	case SearchTabAll, SearchTabTracks:
		m.renderTracks(&b)
	case SearchTabAlbums:
		m.renderAlbums(&b)
	case SearchTabArtists:
		m.renderArtists(&b)
	}

	return b.String()
}

func (m SearchModel) renderTracks(b *strings.Builder) {
	if len(m.tracks) == 0 {
		b.WriteString(theme.S.Muted.Render("  No tracks found") + "\n")
		return
	}
	for i, t := range m.tracks {
		dur := formatTime(t.DurationSec())
		line := fmt.Sprintf("  %s - %s  %s", t.Title, t.ArtistName(), dur)
		if i == m.cursor && m.focused {
			b.WriteString(theme.S.ListActive.Render("▸"+line) + "\n")
		} else {
			b.WriteString(theme.S.ListItem.Render(" "+line) + "\n")
		}
	}
}

func (m SearchModel) renderAlbums(b *strings.Builder) {
	if len(m.albums) == 0 {
		b.WriteString(theme.S.Muted.Render("  No albums found") + "\n")
		return
	}
	for i, a := range m.albums {
		line := fmt.Sprintf("  %s - %s (%d)", a.Title, a.ArtistName(), a.Year)
		if i == m.cursor && m.focused {
			b.WriteString(theme.S.ListActive.Render("▸"+line) + "\n")
		} else {
			b.WriteString(theme.S.ListItem.Render(" "+line) + "\n")
		}
	}
}

func (m SearchModel) renderArtists(b *strings.Builder) {
	if len(m.artists) == 0 {
		b.WriteString(theme.S.Muted.Render("  No artists found") + "\n")
		return
	}
	for i, a := range m.artists {
		line := fmt.Sprintf("  %s", a.Name)
		if i == m.cursor && m.focused {
			b.WriteString(theme.S.ListActive.Render("▸"+line) + "\n")
		} else {
			b.WriteString(theme.S.ListItem.Render(" "+line) + "\n")
		}
	}
}

func (m *SearchModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.input.Width = w - 8
}

func (m *SearchModel) SetFocused(f bool) {
	m.focused = f
	if f && m.result == nil {
		m.inputFocused = true
		m.input.Focus()
	}
}

func (m *SearchModel) Focus() {
	m.inputFocused = true
	m.input.Focus()
}

func parseSearchTracks(result *api.SearchResult) ([]api.Track, error) {
	if result.Tracks == nil || len(result.Tracks.Results) == 0 {
		return nil, nil
	}
	var tracks []api.Track
	return tracks, json.Unmarshal(result.Tracks.Results, &tracks)
}

func parseSearchAlbums(result *api.SearchResult) ([]api.Album, error) {
	if result.Albums == nil || len(result.Albums.Results) == 0 {
		return nil, nil
	}
	var albums []api.Album
	return albums, json.Unmarshal(result.Albums.Results, &albums)
}

func parseSearchArtists(result *api.SearchResult) ([]api.Artist, error) {
	if result.Artists == nil || len(result.Artists.Results) == 0 {
		return nil, nil
	}
	var artists []api.Artist
	return artists, json.Unmarshal(result.Artists.Results, &artists)
}
