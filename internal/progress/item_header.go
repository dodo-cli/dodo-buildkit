package progress

import (
	"fmt"

	"github.com/morikuni/aec"
)

type headerItem struct {
	t        timer
	vertexes *vertexStates
}

func newHeaderItem(v *vertexStates) *headerItem {
	h := &headerItem{vertexes: v}

	if len(v.List()) > 0 {
		h.t = newTimer(v.List()[0].vtx.Started, nil)
	} else {
		h.t = newTimer(nil, nil)
	}

	return h
}

func (h *headerItem) hide() bool {
	// Never hide the header, no matter what
	return false
}

func (h *headerItem) height() int {
	return 1
}

func (h *headerItem) lines(width int) []string {
	if outStr, ok := formatLineWithTimer(width, h.t, h.format); ok {
		return []string{outStr}
	}

	return []string{}
}

func (h *headerItem) format() (string, aec.ANSI) {
	total := len(h.vertexes.List())

	done := 0

	for _, v := range h.vertexes.List() {
		if v.vtx.Completed != nil {
			done++
		}
	}

	status := ""
	// TODO: figure out when we are actually finished
	if done > 0 && done == total {
		status = statusDone
	}

	return fmt.Sprintf("%sBuilding (%d/%d) %s", headerPrefix, done, total, status), aec.WhiteF
}
