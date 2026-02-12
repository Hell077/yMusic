package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"ymusic/internal/api"
	"ymusic/internal/config"
	"ymusic/internal/player"
	"ymusic/internal/theme"
)

const (
	sidebarWidth = 20
	playerHeight = 1
	minWidth     = 80
	minHeight    = 24
)

type FocusArea int

const (
	FocusSidebar FocusArea = iota
	FocusContent
)

type RootModel struct {
	cfg        *config.Config
	client     *api.Client
	player     *player.Controller
	queue      *Queue
	sidebar    SidebarModel
	content    ContentModel
	playerBar  PlayerBarModel
	overlay    OverlayModel
	auth       AuthModel
	nav        *NavStack
	focus      FocusArea
	width      int
	height     int
	authed     bool
	uid        int
	tooSmall   bool
	err        error
}

func NewRoot(cfg *config.Config, client *api.Client, ctrl *player.Controller) RootModel {
	q := NewQueue()
	return RootModel{
		cfg:       cfg,
		client:    client,
		player:    ctrl,
		queue:     q,
		sidebar:   NewSidebar(),
		content:   NewContent(q),
		playerBar: NewPlayerBar(q),
		overlay:   NewOverlay(),
		auth:      NewAuth(),
		nav:       NewNavStack(),
		focus:     FocusSidebar,
		authed:    client != nil,
	}
}

func (m RootModel) Init() tea.Cmd {
	if !m.authed {
		return m.auth.Init()
	}
	return tea.Batch(
		m.fetchAccount(),
		m.fetchFeed(),
		m.startPlayer(),
		m.playerTick(),
	)
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Always handle window size
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = ws.Width
		m.height = ws.Height
		m.tooSmall = ws.Width < minWidth || ws.Height < minHeight
		m.updateLayout()
		return m, nil
	}

	// Route ALL messages to auth when not authenticated
	if !m.authed {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "q" || msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
		case AuthCompleteMsg:
			m.cfg.Token = msg.Token
			m.cfg.Save()
			m.client = api.NewClient(msg.Token)
			m.authed = true
			return m, tea.Batch(
				m.fetchAccount(),
				m.fetchFeed(),
				m.startPlayer(),
				m.playerTick(),
			)
		}
		var cmd tea.Cmd
		m.auth, cmd = m.auth.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.MouseMsg:
		if m.overlay.Visible() {
			var cmd tea.Cmd
			m.overlay, cmd = m.overlay.Update(msg)
			return m, cmd
		}
		contentHeight := m.height - playerHeight
		sidebarW := sidebarWidth + 1 // +1 for border

		if msg.Y >= contentHeight {
			// Player bar area
			cmd := m.playerBar.HandleMouse(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		} else if msg.X < sidebarW {
			// Sidebar area
			m.focus = FocusSidebar
			m.updateFocus()
			if handled, cmd := m.sidebar.HandleMouse(msg); handled {
				cmds = append(cmds, cmd)
			}
		} else {
			// Content area — adjust X for sidebar+border offset
			m.focus = FocusContent
			m.updateFocus()
			adjusted := tea.MouseMsg(tea.MouseEvent{
				X:      msg.X - sidebarW,
				Y:      msg.Y,
				Button: msg.Button,
				Action: msg.Action,
				Shift:  msg.Shift,
				Alt:    msg.Alt,
				Ctrl:   msg.Ctrl,
			})
			var cmd tea.Cmd
			m.content, cmd = m.content.Update(adjusted)
			cmds = append(cmds, cmd)
		}

	case tea.KeyMsg:
		// Overlay captures all input when visible
		if m.overlay.Visible() {
			var cmd tea.Cmd
			m.overlay, cmd = m.overlay.Update(msg)
			return m, cmd
		}

		// When text input is active, pass keys directly to content
		if m.content.IsTextInputActive() {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			}
			var cmd tea.Cmd
			m.content, cmd = m.content.Update(msg)
			return m, cmd
		}

		// Global keys
		switch {
		case key.Matches(msg, Keys.Escape):
			if m.nav.CanGoBack() {
				ps := m.nav.Pop()
				m.content.SetPage(ps.Page)
				m.sidebar.SetCursorForPage(ps.Page)
				return m, nil
			}
			m.overlay.Toggle()
			return m, nil
		case key.Matches(msg, Keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, Keys.Space):
			if m.player != nil {
				m.player.TogglePause()
			}
			return m, nil
		case key.Matches(msg, Keys.Next):
			return m, m.playNext()
		case key.Matches(msg, Keys.Prev):
			return m, m.playPrev()
		case key.Matches(msg, Keys.VolumeUp):
			if m.player != nil {
				st := m.player.GetState()
				vol := st.Volume + 5
				if vol > 100 {
					vol = 100
				}
				m.player.SetVolume(vol)
				m.cfg.Volume = int(vol)
				m.cfg.Save()
			}
			return m, nil
		case key.Matches(msg, Keys.VolumeDn):
			if m.player != nil {
				st := m.player.GetState()
				vol := st.Volume - 5
				if vol < 0 {
					vol = 0
				}
				m.player.SetVolume(vol)
				m.cfg.Volume = int(vol)
				m.cfg.Save()
			}
			return m, nil
		case key.Matches(msg, Keys.SeekFwd):
			if m.player != nil {
				m.player.Seek(10)
			}
			return m, nil
		case key.Matches(msg, Keys.SeekBack):
			if m.player != nil {
				m.player.Seek(-10)
			}
			return m, nil
		case key.Matches(msg, Keys.Like):
			return m, m.toggleLike()
		case key.Matches(msg, Keys.Shuffle):
			m.queue.ToggleShuffle()
			m.content.RefreshQueue()
			return m, nil
		case key.Matches(msg, Keys.Repeat):
			m.queue.CycleRepeat()
			return m, nil
		case key.Matches(msg, Keys.Search):
			if m.focus != FocusContent || m.nav.Current().Page != PageSearch {
				m.navigateTo(PageState{Page: PageSearch})
				m.focus = FocusContent
				m.updateFocus()
			}
			m.content.SearchModel().Focus()
			return m, nil

		// Focus switching with Tab
		case key.Matches(msg, Keys.Tab):
			if m.focus == FocusSidebar {
				m.focus = FocusContent
			} else {
				m.focus = FocusSidebar
			}
			m.updateFocus()
			return m, nil
		}

		// Route to focused area
		if m.focus == FocusSidebar {
			var cmd tea.Cmd
			m.sidebar, cmd = m.sidebar.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			var cmd tea.Cmd
			m.content, cmd = m.content.Update(msg)
			cmds = append(cmds, cmd)
		}

	case FeedMsg:
		var cmd tea.Cmd
		m.content, cmd = m.content.Update(msg)
		cmds = append(cmds, cmd)

	case AccountMsg:
		if msg.Status != nil && msg.Status.Account.UID > 0 {
			m.uid = msg.Status.Account.UID
			cmds = append(cmds, m.fetchLikedTracks(), m.fetchPlaylists(), m.fetchLikedAlbums())
		}

	case NavigateMsg:
		m.navigateTo(PageState{Page: msg.Page})
		m.focus = FocusContent
		m.updateFocus()
		cmds = append(cmds, m.loadPage(msg.Page, ""))

	case NavigateBackMsg:
		ps := m.nav.Pop()
		m.content.SetPage(ps.Page)
		m.sidebar.SetCursorForPage(ps.Page)

	case navigateAlbumMsg:
		ps := PageState{Page: PageAlbum, ID: fmt.Sprintf("%d", msg.id)}
		m.nav.Push(ps)
		m.content.ResetAlbum()
		m.content.SetPage(PageAlbum)
		cmds = append(cmds, m.fetchAlbum(msg.id))

	case navigateArtistMsg:
		ps := PageState{Page: PageArtist, ID: fmt.Sprintf("%d", msg.id)}
		m.nav.Push(ps)
		m.content.ResetArtist()
		m.content.SetPage(PageArtist)
		cmds = append(cmds, m.fetchArtist(msg.id))

	case navigatePlaylistMsg:
		ps := PageState{Page: PagePlaylist, ID: fmt.Sprintf("%d:%d", msg.uid, msg.kind)}
		m.nav.Push(ps)
		m.content.ResetPlaylist()
		m.content.SetPage(PagePlaylist)
		cmds = append(cmds, m.fetchPlaylist(msg.uid, msg.kind))

	case doSearchMsg:
		cmds = append(cmds, m.doSearch(msg.query))

	case PlayTrackMsg:
		m.queue.Set(msg.Queue, msg.Index)
		cmds = append(cmds, m.playCurrentTrack())
		m.content.RefreshQueue()

	case PlayPrevMsg:
		cmds = append(cmds, m.playPrev())

	case PlayNextMsg:
		cmds = append(cmds, m.playNext())

	case TogglePauseMsg:
		if m.player != nil {
			m.player.TogglePause()
		}

	case SeekToMsg:
		if m.player != nil {
			pos := msg.Position * m.player.GetState().Duration
			m.player.SeekAbsolute(pos)
		}

	case SeekRelativeMsg:
		if m.player != nil {
			m.player.Seek(msg.Seconds)
		}

	case VolumeChangeMsg:
		if m.player != nil {
			st := m.player.GetState()
			vol := st.Volume + msg.Delta
			if vol < 0 {
				vol = 0
			}
			if vol > 100 {
				vol = 100
			}
			m.player.SetVolume(vol)
			m.cfg.Volume = int(vol)
			m.cfg.Save()
		}

	case ToggleShuffleMsg:
		m.queue.ToggleShuffle()
		m.content.RefreshQueue()

	case CycleRepeatMsg:
		m.queue.CycleRepeat()

	case TrackURLMsg:
		if m.player != nil {
			m.player.LoadURL(msg.URL)
		}

	case PlayerEventMsg:
		if msg.Event.Type == "end-file" && msg.Event.Name == "eof" {
			cmds = append(cmds, m.playNext())
		}
		cmds = append(cmds, m.listenPlayerEvents())

	case PlayerTickMsg:
		if m.player != nil {
			m.playerBar.SetState(m.player.GetState())
		}
		cmds = append(cmds, m.playerTick())

	case LikeResultMsg:
		// Update liked state in track lists
		// Could propagate to collection if needed

	case listenPlayerMsg:
		cmds = append(cmds, m.listenPlayerEvents())

	case ThemeChangedMsg:
		m.cfg.Theme = theme.Current.Name
		m.cfg.Save()

	case ErrorMsg:
		m.err = msg.Err
		var cmd tea.Cmd
		m.content, cmd = m.content.Update(msg)
		cmds = append(cmds, cmd)

	default:
		// Forward data messages to content
		var cmd tea.Cmd
		m.content, cmd = m.content.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m RootModel) View() string {
	if m.tooSmall {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			theme.S.Error.Render("Terminal too small\nMinimum: 80x24"))
	}

	if !m.authed {
		return m.auth.View()
	}

	// Layout: sidebar | content
	//         player bar
	contentHeight := m.height - playerHeight
	contentWidth := m.width - sidebarWidth

	sidebarView := lipgloss.NewStyle().
		Width(sidebarWidth).
		Height(contentHeight).
		BorderRight(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(theme.Current.Border).
		Render(m.sidebar.View())

	contentView := lipgloss.NewStyle().
		Width(contentWidth).
		Height(contentHeight).
		Render(m.content.View())

	main := lipgloss.JoinHorizontal(lipgloss.Top, sidebarView, contentView)
	playerView := m.playerBar.View()

	view := lipgloss.JoinVertical(lipgloss.Left, main, playerView)

	// Overlay compositing
	if m.overlay.Visible() {
		overlayView := m.overlay.View()
		return overlayView
	}

	return view
}

// --- Layout helpers ---

func (m *RootModel) updateLayout() {
	contentH := m.height - playerHeight
	contentW := m.width - sidebarWidth - 2 // border
	m.sidebar.SetSize(sidebarWidth, contentH)
	m.content.SetSize(contentW, contentH)
	m.playerBar.SetWidth(m.width)
	m.overlay.SetSize(m.width, m.height)
	m.auth.SetSize(m.width, m.height)
}

func (m *RootModel) updateFocus() {
	m.sidebar.SetFocused(m.focus == FocusSidebar)
	m.content.SetFocused(m.focus == FocusContent)
}

func (m *RootModel) navigateTo(ps PageState) {
	m.nav.Push(ps)
	m.content.SetPage(ps.Page)
	m.sidebar.SetCursorForPage(ps.Page)
}

// --- Commands ---

func (m *RootModel) fetchAccount() tea.Cmd {
	client := m.client
	return func() tea.Msg {
		status, err := client.GetAccountStatus()
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return AccountMsg{Status: status}
	}
}

func (m *RootModel) fetchFeed() tea.Cmd {
	client := m.client
	return func() tea.Msg {
		feed, err := client.GetFeed()
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return FeedMsg{Feed: feed}
	}
}

func (m *RootModel) fetchLikedTracks() tea.Cmd {
	client := m.client
	uid := m.uid
	return func() tea.Msg {
		result, err := client.GetLikedTracks(uid)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		// Fetch full track info
		var ids []string
		for _, lt := range result.Library.Tracks {
			ids = append(ids, lt.ID)
		}
		if len(ids) == 0 {
			return LikedTracksMsg{}
		}
		// Limit to first 100
		if len(ids) > 100 {
			ids = ids[:100]
		}
		tracks, err := client.GetTracks(ids)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return LikedTracksMsg{Tracks: tracks}
	}
}

func (m *RootModel) fetchPlaylists() tea.Cmd {
	client := m.client
	uid := m.uid
	return func() tea.Msg {
		playlists, err := client.GetUserPlaylists(uid)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return UserPlaylistsMsg{Playlists: playlists}
	}
}

func (m *RootModel) fetchLikedAlbums() tea.Cmd {
	client := m.client
	uid := m.uid
	return func() tea.Msg {
		liked, err := client.GetLikedAlbums(uid)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		if len(liked) == 0 {
			return LikedAlbumsMsg{}
		}
		var ids []string
		for _, la := range liked {
			ids = append(ids, fmt.Sprintf("%d", la.ID))
		}
		albums, err := client.GetAlbums(ids)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return LikedAlbumsMsg{Albums: albums}
	}
}

func (m *RootModel) fetchPlaylist(uid, kind int) tea.Cmd {
	client := m.client
	return func() tea.Msg {
		p, err := client.GetPlaylist(uid, kind)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return PlaylistMsg{Playlist: p}
	}
}

func (m *RootModel) fetchAlbum(id int) tea.Cmd {
	client := m.client
	return func() tea.Msg {
		album, err := client.GetAlbumWithTracks(id)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return AlbumMsg{Album: album}
	}
}

func (m *RootModel) fetchArtist(id int) tea.Cmd {
	client := m.client
	return func() tea.Msg {
		info, err := client.GetArtistBriefInfo(id)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return ArtistInfoMsg{Info: info}
	}
}

func (m *RootModel) doSearch(query string) tea.Cmd {
	client := m.client
	return func() tea.Msg {
		result, err := client.Search(query, 0)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return SearchResultMsg{Result: result}
	}
}

func (m *RootModel) loadPage(page Page, id string) tea.Cmd {
	switch page {
	case PageMyWave:
		return m.fetchRadio()
	}
	return nil
}

func (m *RootModel) fetchRadio() tea.Cmd {
	client := m.client
	wave := m.content.MyWaveModel()
	lastID := wave.LastTrackID()
	return func() tea.Msg {
		result, err := client.GetStationTracks("user", "onyourwave", lastID)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return RadioTracksMsg{
			Tracks:  result.Sequence,
			BatchID: result.BatchID,
		}
	}
}

func (m *RootModel) playCurrentTrack() tea.Cmd {
	t := m.queue.Current()
	if t == nil {
		return nil
	}
	m.playerBar.SetTrack(t)
	m.content.SetPlayingTrack(t.ID)

	client := m.client
	trackID := t.ID
	return func() tea.Msg {
		url, err := client.GetDirectURL(trackID)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return TrackURLMsg{TrackID: trackID, URL: url}
	}
}

func (m *RootModel) playNext() tea.Cmd {
	t := m.queue.Next()
	if t == nil {
		return nil
	}
	m.playerBar.SetTrack(t)
	m.content.RefreshQueue()
	m.content.SetPlayingTrack(t.ID)
	return m.playCurrentTrack()
}

func (m *RootModel) playPrev() tea.Cmd {
	t := m.queue.Prev()
	if t == nil {
		return nil
	}
	m.playerBar.SetTrack(t)
	m.content.RefreshQueue()
	m.content.SetPlayingTrack(t.ID)
	return m.playCurrentTrack()
}

func (m *RootModel) toggleLike() tea.Cmd {
	t := m.queue.Current()
	if t == nil || m.client == nil {
		return nil
	}
	client := m.client
	uid := m.uid
	trackID := t.ID
	return func() tea.Msg {
		// Simple toggle - just try to like
		err := client.LikeTrack(uid, trackID)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return LikeResultMsg{TrackID: trackID, Liked: true}
	}
}

func (m *RootModel) startPlayer() tea.Cmd {
	ctrl := m.player
	return func() tea.Msg {
		if err := ctrl.Start(); err != nil {
			return ErrorMsg{Err: fmt.Errorf("start player: %w", err)}
		}
		return listenPlayerMsg{}
	}
}

type listenPlayerMsg struct{}

func (m *RootModel) listenPlayerEvents() tea.Cmd {
	ctrl := m.player
	return func() tea.Msg {
		event := <-ctrl.Events
		return PlayerEventMsg{Event: event}
	}
}

func (m *RootModel) playerTick() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return PlayerTickMsg{}
	})
}
