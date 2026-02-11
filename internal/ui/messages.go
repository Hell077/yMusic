package ui

import (
	"ymusic/internal/api"
	"ymusic/internal/player"
)

// Auth messages
type AuthCompleteMsg struct{ Token string }
type AuthErrorMsg struct{ Err error }

// API data messages
type AccountMsg struct{ Status *api.AccountStatus }
type FeedMsg struct{ Feed *api.FeedResponse }
type SearchResultMsg struct{ Result *api.SearchResult }
type PlaylistMsg struct{ Playlist *api.Playlist }
type AlbumMsg struct{ Album *api.Album }
type ArtistInfoMsg struct{ Info *api.ArtistBriefInfo }
type LikedTracksMsg struct{ Tracks []api.Track }
type UserPlaylistsMsg struct{ Playlists []api.Playlist }
type LikedAlbumsMsg struct{ Albums []api.Album }
type TracksMsg struct{ Tracks []api.Track }
type RadioTracksMsg struct {
	Tracks  []api.StationTrack
	BatchID string
}

// Player messages
type PlayerEventMsg struct{ Event player.Event }
type PlayerTickMsg struct{}
type TrackURLMsg struct {
	TrackID string
	URL     string
}

// Navigation
type NavigateMsg struct{ Page Page }
type NavigateBackMsg struct{}

// Actions
type PlayTrackMsg struct {
	Track api.Track
	Queue []api.Track
	Index int
}
type LikeToggleMsg struct{ TrackID string }
type LikeResultMsg struct {
	TrackID string
	Liked   bool
}

// UI
type ErrorMsg struct{ Err error }
type WindowSizeMsg struct {
	Width  int
	Height int
}
type ThemeChangedMsg struct{}
type OverlayToggleMsg struct{}
type FocusSidebarMsg struct{}
type FocusContentMsg struct{}
