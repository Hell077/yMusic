package ui

import (
	"math/rand"
	"ymusic/internal/api"
)

type RepeatMode int

const (
	RepeatOff RepeatMode = iota
	RepeatAll
	RepeatOne
)

func (r RepeatMode) String() string {
	switch r {
	case RepeatAll:
		return "all"
	case RepeatOne:
		return "one"
	default:
		return "off"
	}
}

func (r RepeatMode) Icon() string {
	switch r {
	case RepeatAll:
		return "[R:all]"
	case RepeatOne:
		return "[R:one]"
	default:
		return "[R:off]"
	}
}

type Queue struct {
	tracks       []api.Track
	originalOrder []api.Track
	current      int
	shuffle      bool
	repeat       RepeatMode
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Set(tracks []api.Track, startIndex int) {
	q.originalOrder = make([]api.Track, len(tracks))
	copy(q.originalOrder, tracks)
	q.tracks = make([]api.Track, len(tracks))
	copy(q.tracks, tracks)
	q.current = startIndex

	if q.shuffle {
		q.doShuffle()
	}
}

func (q *Queue) Current() *api.Track {
	if len(q.tracks) == 0 || q.current < 0 || q.current >= len(q.tracks) {
		return nil
	}
	return &q.tracks[q.current]
}

func (q *Queue) Next() *api.Track {
	if len(q.tracks) == 0 {
		return nil
	}
	if q.repeat == RepeatOne {
		return q.Current()
	}
	q.current++
	if q.current >= len(q.tracks) {
		if q.repeat == RepeatAll {
			q.current = 0
		} else {
			q.current = len(q.tracks) - 1
			return nil
		}
	}
	return q.Current()
}

func (q *Queue) Prev() *api.Track {
	if len(q.tracks) == 0 {
		return nil
	}
	q.current--
	if q.current < 0 {
		if q.repeat == RepeatAll {
			q.current = len(q.tracks) - 1
		} else {
			q.current = 0
		}
	}
	return q.Current()
}

func (q *Queue) ToggleShuffle() {
	q.shuffle = !q.shuffle
	if q.shuffle {
		q.doShuffle()
	} else {
		cur := q.Current()
		q.tracks = make([]api.Track, len(q.originalOrder))
		copy(q.tracks, q.originalOrder)
		if cur != nil {
			for i, t := range q.tracks {
				if t.ID == cur.ID {
					q.current = i
					break
				}
			}
		}
	}
}

func (q *Queue) doShuffle() {
	if len(q.tracks) <= 1 {
		return
	}
	cur := q.Current()
	rand.Shuffle(len(q.tracks), func(i, j int) {
		q.tracks[i], q.tracks[j] = q.tracks[j], q.tracks[i]
	})
	if cur != nil {
		for i, t := range q.tracks {
			if t.ID == cur.ID {
				q.tracks[0], q.tracks[i] = q.tracks[i], q.tracks[0]
				break
			}
		}
		q.current = 0
	}
}

func (q *Queue) CycleRepeat() {
	q.repeat = (q.repeat + 1) % 3
}

func (q *Queue) IsShuffled() bool   { return q.shuffle }
func (q *Queue) RepeatMode() RepeatMode { return q.repeat }
func (q *Queue) Tracks() []api.Track { return q.tracks }
func (q *Queue) Index() int          { return q.current }
func (q *Queue) Len() int            { return len(q.tracks) }

func (q *Queue) Upcoming() []api.Track {
	if q.current+1 >= len(q.tracks) {
		return nil
	}
	return q.tracks[q.current+1:]
}

func (q *Queue) Append(tracks ...api.Track) {
	q.tracks = append(q.tracks, tracks...)
	q.originalOrder = append(q.originalOrder, tracks...)
}
