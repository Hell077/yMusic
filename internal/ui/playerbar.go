package ui

import (
	"fmt"
	"strings"

	"ymusic/internal/api"
	"ymusic/internal/player"
	"ymusic/internal/theme"
)

type PlayerBarModel struct {
	track    *api.Track
	state    player.State
	queue    *Queue
	width    int
}

func NewPlayerBar(queue *Queue) PlayerBarModel {
	return PlayerBarModel{queue: queue}
}

func (m *PlayerBarModel) SetTrack(t *api.Track) {
	m.track = t
}

func (m *PlayerBarModel) SetState(s player.State) {
	m.state = s
}

func (m *PlayerBarModel) SetWidth(w int) {
	m.width = w
}

func (m PlayerBarModel) View() string {
	if m.track == nil {
		return theme.S.PlayerBar.Width(m.width).Render(
			theme.S.Muted.Render("No track playing"),
		)
	}

	title := theme.S.Primary.Render(truncate(m.track.Title, 30))
	artist := theme.S.Muted.Render(truncate(m.track.ArtistName(), 25))

	var playIcon string
	if m.state.Playing {
		playIcon = "▶"
	} else {
		playIcon = "⏸"
	}

	pos := formatTime(int(m.state.Position))
	dur := formatTime(int(m.state.Duration))
	timeStr := fmt.Sprintf("%s / %s", pos, dur)

	// Progress bar
	barWidth := m.width - 60
	if barWidth < 10 {
		barWidth = 10
	}
	progress := 0.0
	if m.state.Duration > 0 {
		progress = m.state.Position / m.state.Duration
	}
	filled := int(float64(barWidth) * progress)
	if filled > barWidth {
		filled = barWidth
	}
	bar := theme.S.ProgressFull.Render(strings.Repeat("━", filled)) +
		theme.S.ProgressEmpty.Render(strings.Repeat("─", barWidth-filled))

	vol := fmt.Sprintf("♪ %.0f%%", m.state.Volume)

	// Shuffle/Repeat indicators
	var shuffleIcon string
	if m.queue.IsShuffled() {
		shuffleIcon = theme.S.Primary.Render("[S]")
	} else {
		shuffleIcon = theme.S.Muted.Render("[S]")
	}
	repeatIcon := theme.S.Muted.Render(m.queue.RepeatMode().Icon())
	if m.queue.RepeatMode() != RepeatOff {
		repeatIcon = theme.S.Primary.Render(m.queue.RepeatMode().Icon())
	}

	info := fmt.Sprintf(" %s %s - %s  %s  %s  %s %s %s",
		playIcon, title, artist, bar, timeStr, vol, shuffleIcon, repeatIcon,
	)

	return theme.S.PlayerBar.Width(m.width).Render(info)
}

func formatTime(s int) string {
	if s < 0 {
		s = 0
	}
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max-1]) + "…"
}
