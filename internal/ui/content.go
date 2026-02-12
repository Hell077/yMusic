package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ContentModel struct {
	home         HomeModel
	search       SearchModel
	collection   CollectionModel
	playlistView PlaylistViewModel
	albumView    AlbumViewModel
	artistView   ArtistViewModel
	queueView    QueueViewModel
	myWave       MyWaveModel
	activePage   Page
	width        int
	height       int
	focused      bool
}

func NewContent(queue *Queue) ContentModel {
	return ContentModel{
		home:         NewHome(),
		search:       NewSearch(),
		collection:   NewCollection(),
		playlistView: NewPlaylistView(),
		albumView:    NewAlbumView(),
		artistView:   NewArtistView(),
		queueView:    NewQueueView(queue),
		myWave:       NewMyWave(),
		activePage:   PageHome,
	}
}

func (m ContentModel) Init() tea.Cmd { return nil }

func (m ContentModel) Update(msg tea.Msg) (ContentModel, tea.Cmd) {
	var cmd tea.Cmd

	// Route messages to all sub-models that need data updates
	switch msg.(type) {
	case FeedMsg:
		m.home, cmd = m.home.Update(msg)
		return m, cmd
	case SearchResultMsg:
		m.search, cmd = m.search.Update(msg)
		return m, cmd
	case PlaylistMsg:
		m.playlistView, cmd = m.playlistView.Update(msg)
		return m, cmd
	case AlbumMsg:
		m.albumView, cmd = m.albumView.Update(msg)
		return m, cmd
	case ArtistInfoMsg:
		m.artistView, cmd = m.artistView.Update(msg)
		return m, cmd
	case LikedTracksMsg:
		m.collection, cmd = m.collection.Update(msg)
		return m, cmd
	case UserPlaylistsMsg:
		m.collection, cmd = m.collection.Update(msg)
		return m, cmd
	case RadioTracksMsg:
		m.myWave, cmd = m.myWave.Update(msg)
		return m, cmd
	}

	// Route key messages to active page only
	switch m.activePage {
	case PageHome:
		m.home, cmd = m.home.Update(msg)
	case PageSearch:
		m.search, cmd = m.search.Update(msg)
	case PageCollection:
		m.collection, cmd = m.collection.Update(msg)
	case PagePlaylist:
		m.playlistView, cmd = m.playlistView.Update(msg)
	case PageAlbum:
		m.albumView, cmd = m.albumView.Update(msg)
	case PageArtist:
		m.artistView, cmd = m.artistView.Update(msg)
	case PageQueue:
		m.queueView, cmd = m.queueView.Update(msg)
	case PageMyWave:
		m.myWave, cmd = m.myWave.Update(msg)
	}

	return m, cmd
}

func (m ContentModel) View() string {
	switch m.activePage {
	case PageHome:
		return m.home.View()
	case PageSearch:
		return m.search.View()
	case PageCollection:
		return m.collection.View()
	case PagePlaylist:
		return m.playlistView.View()
	case PageAlbum:
		return m.albumView.View()
	case PageArtist:
		return m.artistView.View()
	case PageQueue:
		return m.queueView.View()
	case PageMyWave:
		return m.myWave.View()
	default:
		return ""
	}
}

func (m *ContentModel) SetPage(p Page) {
	m.activePage = p
	m.updateFocus()
}

func (m *ContentModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.home.SetSize(w, h)
	m.search.SetSize(w, h)
	m.collection.SetSize(w, h)
	m.playlistView.SetSize(w, h)
	m.albumView.SetSize(w, h)
	m.artistView.SetSize(w, h)
	m.queueView.SetSize(w, h)
	m.myWave.SetSize(w, h)
}

func (m *ContentModel) SetFocused(f bool) {
	m.focused = f
	m.updateFocus()
}

func (m *ContentModel) updateFocus() {
	m.home.SetFocused(m.focused && m.activePage == PageHome)
	m.search.SetFocused(m.focused && m.activePage == PageSearch)
	m.collection.SetFocused(m.focused && m.activePage == PageCollection)
	m.playlistView.SetFocused(m.focused && m.activePage == PagePlaylist)
	m.albumView.SetFocused(m.focused && m.activePage == PageAlbum)
	m.artistView.SetFocused(m.focused && m.activePage == PageArtist)
	m.queueView.SetFocused(m.focused && m.activePage == PageQueue)
	m.myWave.SetFocused(m.focused && m.activePage == PageMyWave)
}

func (m *ContentModel) ResetPlaylist() {
	m.playlistView = NewPlaylistView()
	m.playlistView.SetSize(m.width, m.height)
}

func (m *ContentModel) ResetAlbum() {
	m.albumView = NewAlbumView()
	m.albumView.SetSize(m.width, m.height)
}

func (m *ContentModel) ResetArtist() {
	m.artistView = NewArtistView()
	m.artistView.SetSize(m.width, m.height)
}

func (m *ContentModel) RefreshQueue() {
	m.queueView.Refresh()
}

func (m *ContentModel) SetPlayingTrack(id string) {
	m.queueView.trackList.SetPlaying(id)
	m.playlistView.trackList.SetPlaying(id)
	m.albumView.trackList.SetPlaying(id)
	m.artistView.trackList.SetPlaying(id)
	m.collection.trackList.SetPlaying(id)
	m.myWave.trackList.SetPlaying(id)
}

func (m *ContentModel) SearchModel() *SearchModel {
	return &m.search
}

func (m *ContentModel) MyWaveModel() *MyWaveModel {
	return &m.myWave
}

func (m *ContentModel) IsTextInputActive() bool {
	return m.activePage == PageSearch && m.search.inputFocused
}
