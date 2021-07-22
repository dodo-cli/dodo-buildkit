package progress

import (
	"time"

	digest "github.com/opencontainers/go-digest"
)

const (
	statusError    = "ERROR"
	statusCanceled = "CANCELED"
	statusCached   = "CACHED"
	statusDone     = "FINISHED"

	headerPrefix = "[+] "
	prefix       = " |- "
	errorDelim   = "------"
	termHeight   = 6
)

type item interface {
	lines(int) []string
	height() int
	hide() bool
}

type itemCache map[digest.Digest][]item

func newItemCache() itemCache {
	return make(map[digest.Digest][]item)
}

func (c itemCache) clear(k digest.Digest) {
	c[k] = []item{}
}

func (c itemCache) get(k digest.Digest, gen func() []item) []item {
	if items, ok := c[k]; ok && len(items) > 0 {
		return items
	}

	newItems := gen()
	c[k] = newItems

	return newItems
}

type timer struct {
	start *time.Time
	end   *time.Time
}

func newTimer(start *time.Time, end *time.Time) timer {
	return timer{start: start, end: end}
}

func (t timer) Started() (bool, time.Time) {
	if t.start == nil {
		return false, time.Time{}
	}

	start := t.start

	return true, *start
}

func (t timer) Completed() (bool, time.Time) {
	if t.end == nil {
		return false, time.Time{}
	}

	end := t.end

	return true, *end
}

func (t timer) Running() (bool, time.Duration) {
	started, startTime := t.Started()
	if !started {
		return false, 0 * time.Second
	}

	done, endTime := t.Completed()
	if !done {
		endTime = time.Now()
	}

	runTime := endTime.Sub(startTime)
	if runTime < 50*time.Millisecond {
		return true, 0 * time.Second
	}

	return true, runTime
}
