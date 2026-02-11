package api

import "encoding/json"

type Track struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	DurationMs    int      `json:"durationMs"`
	Artists       []Artist `json:"artists"`
	Albums        []Album  `json:"albums"`
	CoverURI      string   `json:"coverUri"`
	OgImage       string   `json:"ogImage"`
	Liked         bool     `json:"-"`
	Available     bool     `json:"available"`
	LyricsAvail   bool     `json:"lyricsAvailable"`
}

func (t *Track) UnmarshalJSON(data []byte) error {
	type Alias Track
	aux := &struct {
		ID json.Number `json:"id"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	t.ID = aux.ID.String()
	return nil
}

func (t Track) ArtistName() string {
	if len(t.Artists) == 0 {
		return "Unknown"
	}
	name := t.Artists[0].Name
	for i := 1; i < len(t.Artists); i++ {
		name += ", " + t.Artists[i].Name
	}
	return name
}

func (t Track) DurationSec() int {
	return t.DurationMs / 1000
}

func (t Track) AlbumTitle() string {
	if len(t.Albums) == 0 {
		return ""
	}
	return t.Albums[0].Title
}

type Artist struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Cover   *Cover   `json:"cover"`
	Various bool     `json:"various"`
	Genres  []string `json:"genres"`
}

func (a *Artist) UnmarshalJSON(data []byte) error {
	type Alias Artist
	aux := &struct {
		ID json.Number `json:"id"`
		*Alias
	}{Alias: (*Alias)(a)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	id, _ := aux.ID.Int64()
	a.ID = int(id)
	return nil
}

type Album struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Year        int      `json:"year"`
	CoverURI    string   `json:"coverUri"`
	OgImage     string   `json:"ogImage"`
	TrackCount  int      `json:"trackCount"`
	Genre       string   `json:"genre"`
	Artists     []Artist `json:"artists"`
	Volumes     [][]Track `json:"volumes"`
	LikesCount  int      `json:"likesCount"`
}

func (a Album) ArtistName() string {
	if len(a.Artists) == 0 {
		return "Unknown"
	}
	name := a.Artists[0].Name
	for i := 1; i < len(a.Artists); i++ {
		name += ", " + a.Artists[i].Name
	}
	return name
}

type Cover struct {
	URI  string `json:"uri"`
	Type string `json:"type"`
}

type Playlist struct {
	UID         int    `json:"uid"`
	Kind        int    `json:"kind"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TrackCount  int    `json:"trackCount"`
	Cover       *Cover `json:"cover"`
	OgImage     string `json:"ogImage"`
	Owner       Owner  `json:"owner"`
	DurationMs  int    `json:"durationMs"`
	Tracks      []TrackItem `json:"tracks"`
}

type Owner struct {
	UID  int    `json:"uid"`
	Name string `json:"name"`
}

type TrackItem struct {
	ID      json.Number `json:"id"`
	AlbumID json.Number `json:"albumId"`
	Track   Track       `json:"track"`
}

type SearchResult struct {
	Text    string        `json:"text"`
	Best    *SearchBest   `json:"best"`
	Tracks  *SearchBlock  `json:"tracks"`
	Albums  *SearchBlock  `json:"albums"`
	Artists *SearchBlock  `json:"artists"`
	Playlists *SearchBlock `json:"playlists"`
}

type SearchBest struct {
	Type   string      `json:"type"`
	Result interface{} `json:"result"`
}

type SearchBlock struct {
	Total   int             `json:"total"`
	Results json.RawMessage `json:"results"`
}

type DownloadInfo struct {
	Codec      string `json:"codec"`
	Bitrate    int    `json:"bitrateInKbps"`
	Src        string `json:"downloadInfoUrl"`
	Direct     bool   `json:"direct"`
	Gain       bool   `json:"gain"`
	Preview    bool   `json:"preview"`
}

type DownloadData struct {
	Host string `json:"host"`
	Path string `json:"path"`
	TS   string `json:"ts"`
	S    string `json:"s"`
}

type AccountStatus struct {
	Account Account `json:"account"`
}

type Account struct {
	UID         int    `json:"uid"`
	Login       string `json:"login"`
	DisplayName string `json:"displayName"`
	FullName    string `json:"fullName"`
}

type FeedResponse struct {
	GeneratedPlaylists []GeneratedPlaylist `json:"generatedPlaylists"`
	Days               []FeedDay           `json:"days"`
}

type GeneratedPlaylist struct {
	Type     string   `json:"type"`
	Ready    bool     `json:"ready"`
	Notify   bool     `json:"notify"`
	Data     Playlist `json:"data"`
}

type FeedDay struct {
	Day    string      `json:"day"`
	Events []FeedEvent `json:"events"`
}

type FeedEvent struct {
	ID    string  `json:"id"`
	Type  string  `json:"type"`
	Title string  `json:"title"`
	Tracks []TrackItem `json:"tracks"`
	Albums []Album     `json:"albums"`
	Artists []Artist   `json:"artists"`
}

type StationResult struct {
	Station  Station `json:"station"`
	Settings interface{} `json:"settings"`
}

type Station struct {
	ID   StationID `json:"id"`
	Name string    `json:"name"`
}

type StationID struct {
	Type string `json:"type"`
	Tag  string `json:"tag"`
}

type StationTracksResult struct {
	ID       StationID `json:"id"`
	Sequence []StationTrack `json:"sequence"`
	BatchID  string    `json:"batchId"`
}

type StationTrack struct {
	Track Track `json:"track"`
	Liked bool  `json:"liked"`
}

type LikesResult struct {
	Library LikesLibrary `json:"library"`
}

type LikesLibrary struct {
	UID    int         `json:"uid"`
	Tracks []LikeTrack `json:"tracks"`
}

type LikeTrack struct {
	ID        string `json:"id"`
	AlbumID   string `json:"albumId"`
	Timestamp string `json:"timestamp"`
}

type ArtistBriefInfo struct {
	Artist       Artist    `json:"artist"`
	Albums       []Album   `json:"albums"`
	AlsoAlbums   []Album   `json:"alsoAlbums"`
	PopularTracks []Track  `json:"popularTracks"`
	SimilarArtists []Artist `json:"similarArtists"`
}
