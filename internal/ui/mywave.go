package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"ymusic/internal/api"
	"ymusic/internal/theme"
)

type MyWaveModel struct {
	tracks    []api.Track
	trackList TrackListModel
	batchID   string
	loading   bool
	err       error
	width     int
	height    int
	focused   bool
}

func NewMyWave() MyWaveModel {
	return MyWaveModel{
		trackList: NewTrackList(),
	}
}

func (m MyWaveModel) Init() tea.Cmd { return nil }

func (m MyWaveModel) Update(msg tea.Msg) (MyWaveModel, tea.Cmd) {
	switch msg := msg.(type) {
	case RadioTracksMsg:
		m.loading = false
		m.batchID = msg.BatchID
		for _, st := range msg.Tracks {
			m.tracks = append(m.tracks, st.Track)
		}
		m.trackList.SetTracks(m.tracks)
	case ErrorMsg:
		m.err = msg.Err
		m.loading = false
	case tea.MouseMsg:
		// Header: title(0) + blank(1), tracklist starts at row 2
		if handled, cmd := m.trackList.HandleMouse(msg, 2); handled {
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
				tracks := m.tracks
				idx := m.trackList.Cursor()
				return m, func() tea.Msg {
					return PlayTrackMsg{Track: *t, Queue: tracks, Index: idx}
				}
			}
		}
	}
	return m, nil
}

func (m MyWaveModel) View() string {
	var b strings.Builder

	b.WriteString(theme.S.Title.Render(" ≈ My Wave") + "\n\n")

	if m.loading {
		b.WriteString(theme.S.Muted.Render("  Loading radio...") + "\n")
		return b.String()
	}
	if m.err != nil {
		b.WriteString(theme.S.Error.Render("  Error: "+m.err.Error()) + "\n")
		return b.String()
	}

	if len(m.tracks) == 0 {
		b.WriteString(theme.S.Muted.Render("  Press Enter to start My Wave") + "\n")
		return b.String()
	}

	b.WriteString(m.trackList.View())
	return b.String()
}

func (m *MyWaveModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.trackList.SetSize(w, h-4)
}

func (m *MyWaveModel) SetFocused(f bool) {
	m.focused = f
	m.trackList.SetFocused(f)
}

func (m *MyWaveModel) SetLoading() {
	m.loading = true
}

func (m MyWaveModel) BatchID() string {
	return m.batchID
}

func (m MyWaveModel) LastTrackID() string {
	if len(m.tracks) == 0 {
		return ""
	}
	return m.tracks[len(m.tracks)-1].ID
}
