package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"ymusic/internal/api"
	"ymusic/internal/player"
	"ymusic/internal/theme"
)

type PlayerBarModel struct {
	track    *api.Track
	state    player.State
	queue    *Queue
	width    int

	// Click area X ranges (set during View)
	prevX      [2]int // [start, end)
	playX      [2]int
	nextX      [2]int
	barX       [2]int
	barWidth   int
	shuffleX   [2]int
	repeatX    [2]int
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

func (m *PlayerBarModel) View() string {
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
	barWidth := m.width - 68
	if barWidth < 10 {
		barWidth = 10
	}
	m.barWidth = barWidth
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
	repeatIconStr := m.queue.RepeatMode().Icon()
	var repeatIcon string
	if m.queue.RepeatMode() != RepeatOff {
		repeatIcon = theme.S.Primary.Render(repeatIconStr)
	} else {
		repeatIcon = theme.S.Muted.Render(repeatIconStr)
	}

	// Track X positions for click areas (account for 1-char padding from PlayerBar style)
	x := 1 // PlayerBar has Padding(0,1) so content starts at x=1
	m.prevX = [2]int{x, x + 2}
	x += 3 // "⏮ "
	m.playX = [2]int{x, x + 2}
	x += 3 // "▶ "
	m.nextX = [2]int{x, x + 2}

	info := fmt.Sprintf(" ⏮ %s ⏭  %s - %s  %s  %s  %s %s %s",
		playIcon, title, artist, bar, timeStr, vol, shuffleIcon, repeatIcon,
	)

	// Calculate bar X position
	prefix := fmt.Sprintf(" ⏮ %s ⏭  %s - %s  ",
		playIcon, truncate(m.track.Title, 30), truncate(m.track.ArtistName(), 25))
	prefixLen := len([]rune(prefix))
	m.barX = [2]int{1 + prefixLen, 1 + prefixLen + barWidth}

	// Calculate suffix positions from the end
	// suffix: "  %s  %s %s %s" = timeStr, vol, shuffleIcon, repeatIcon
	suffix := fmt.Sprintf("  %s  %s %s %s", timeStr, vol, "[S]", repeatIconStr)
	suffixLen := len([]rune(suffix))
	endX := m.width - 1 // padding right
	suffixStart := endX - suffixLen

	// [S] and repeat icon positions
	// suffix layout: "  time  vol [S] [R]"
	volStr := fmt.Sprintf("  %s  %s ", timeStr, vol)
	volLen := len([]rune(volStr))
	m.shuffleX = [2]int{suffixStart + volLen, suffixStart + volLen + 3}
	m.repeatX = [2]int{suffixStart + volLen + 4, suffixStart + volLen + 4 + len([]rune(repeatIconStr))}

	return theme.S.PlayerBar.Width(m.width).Render(info)
}

// HandleMouse processes mouse events on the player bar.
func (m *PlayerBarModel) HandleMouse(msg tea.MouseMsg) tea.Cmd {
	if m.track == nil {
		return nil
	}

	x := msg.X

	// Wheel: on progress bar = seek, elsewhere = volume
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		if x >= m.barX[0] && x < m.barX[1] {
			return func() tea.Msg { return SeekRelativeMsg{Seconds: 5} }
		}
		return func() tea.Msg { return VolumeChangeMsg{Delta: 5} }
	case tea.MouseButtonWheelDown:
		if x >= m.barX[0] && x < m.barX[1] {
			return func() tea.Msg { return SeekRelativeMsg{Seconds: -5} }
		}
		return func() tea.Msg { return VolumeChangeMsg{Delta: -5} }
	}

	if msg.Action != tea.MouseActionPress || msg.Button != tea.MouseButtonLeft {
		return nil
	}

	if x >= m.prevX[0] && x < m.prevX[1] {
		return func() tea.Msg { return PlayPrevMsg{} }
	}
	if x >= m.playX[0] && x < m.playX[1] {
		return func() tea.Msg { return TogglePauseMsg{} }
	}
	if x >= m.nextX[0] && x < m.nextX[1] {
		return func() tea.Msg { return PlayNextMsg{} }
	}
	if x >= m.barX[0] && x < m.barX[1] && m.barWidth > 0 {
		pos := float64(x-m.barX[0]) / float64(m.barWidth)
		if pos < 0 {
			pos = 0
		}
		if pos > 1 {
			pos = 1
		}
		return func() tea.Msg { return SeekToMsg{Position: pos} }
	}
	if x >= m.shuffleX[0] && x < m.shuffleX[1] {
		return func() tea.Msg { return ToggleShuffleMsg{} }
	}
	if x >= m.repeatX[0] && x < m.repeatX[1] {
		return func() tea.Msg { return CycleRepeatMsg{} }
	}

	return nil
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
