package ui

type Page int

const (
	PageHome Page = iota
	PageSearch
	PageCollection
	PagePlaylist
	PageAlbum
	PageArtist
	PageQueue
	PageMyWave
)

func (p Page) String() string {
	switch p {
	case PageHome:
		return "Home"
	case PageSearch:
		return "Search"
	case PageCollection:
		return "Collection"
	case PagePlaylist:
		return "Playlist"
	case PageAlbum:
		return "Album"
	case PageArtist:
		return "Artist"
	case PageQueue:
		return "Queue"
	case PageMyWave:
		return "My Wave"
	default:
		return "Unknown"
	}
}

type NavStack struct {
	pages []PageState
}

type PageState struct {
	Page Page
	ID   string // playlist/album/artist ID
}

func NewNavStack() *NavStack {
	return &NavStack{
		pages: []PageState{{Page: PageHome}},
	}
}

func (n *NavStack) Current() PageState {
	if len(n.pages) == 0 {
		return PageState{Page: PageHome}
	}
	return n.pages[len(n.pages)-1]
}

func (n *NavStack) Push(p PageState) {
	n.pages = append(n.pages, p)
}

func (n *NavStack) Pop() PageState {
	if len(n.pages) <= 1 {
		return n.Current()
	}
	n.pages = n.pages[:len(n.pages)-1]
	return n.Current()
}

func (n *NavStack) CanGoBack() bool {
	return len(n.pages) > 1
}

func (n *NavStack) Replace(p PageState) {
	if len(n.pages) == 0 {
		n.pages = []PageState{p}
	} else {
		n.pages[len(n.pages)-1] = p
	}
}
