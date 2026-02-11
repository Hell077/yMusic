package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"ymusic/internal/theme"
)

type QueueViewModel struct {
	queue     *Queue
	trackList TrackListModel
	width     int
	height    int
	focused   bool
}

func NewQueueView(queue *Queue) QueueViewModel {
	return QueueViewModel{
		queue:     queue,
		trackList: NewTrackList(),
	}
}

func (m QueueViewModel) Init() tea.Cmd { return nil }

func (m *QueueViewModel) Refresh() {
	m.trackList.SetTracks(m.queue.Tracks())
	cur := m.queue.Current()
	if cur != nil {
		m.trackList.SetPlaying(cur.ID)
	}
}

func (m QueueViewModel) Update(msg tea.Msg) (QueueViewModel, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok && m.focused {
		switch msg.String() {
		case "up", "k":
			m.trackList.MoveUp()
		case "down", "j":
			m.trackList.MoveDown()
		case "enter":
			t := m.trackList.Selected()
			if t != nil {
				tracks := m.queue.Tracks()
				idx := m.trackList.Cursor()
				return m, func() tea.Msg {
					return PlayTrackMsg{Track: *t, Queue: tracks, Index: idx}
				}
			}
		}
	}
	return m, nil
}

func (m QueueViewModel) View() string {
	var b strings.Builder

	b.WriteString(theme.S.Title.Render(" Queue") + "\n")

	cur := m.queue.Current()
	if cur != nil {
		b.WriteString(theme.S.Subtitle.Render("  Now playing: "))
		b.WriteString(theme.S.Primary.Render(cur.Title))
		b.WriteString(theme.S.Muted.Render(" - "+cur.ArtistName()))
		b.WriteString("\n\n")
	} else {
		b.WriteString(theme.S.Muted.Render("  Nothing playing") + "\n\n")
	}

	if m.queue.Len() == 0 {
		b.WriteString(theme.S.Muted.Render("  Queue is empty") + "\n")
		return b.String()
	}

	b.WriteString(m.trackList.View())
	return b.String()
}

func (m *QueueViewModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.trackList.SetSize(w, h-5)
}

func (m *QueueViewModel) SetFocused(f bool) {
	m.focused = f
	m.trackList.SetFocused(f)
}
