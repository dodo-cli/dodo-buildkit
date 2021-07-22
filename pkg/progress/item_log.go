package progress

import (
	"fmt"

	"github.com/jaguilar/vt100"
	"github.com/morikuni/aec"
)

type logItem struct {
	t      timer
	term   *vt100.VT100
	hidden bool
}

func newLogItem(v *vertexState) *logItem {
	return &logItem{
		t:      newTimer(v.vtx.Started, v.vtx.Completed),
		term:   v.term,
		hidden: false,
	}
}

func (l *logItem) hide() bool {
	// We always allow hiding log lines if we need to
	l.hidden = true

	return true
}

func (l *logItem) height() int {
	return l.term.UsedHeight()
}

func (l *logItem) lines(width int) []string {
	lines := []string{}

	if l.hidden {
		return lines
	}

	// Account for the string formatting below
	l.term.Resize(termHeight, width-len(prefix)-3)

	for _, line := range l.term.Content {
		if !empty(line) {
			lines = append(lines, aec.Apply(fmt.Sprintf("%s # %s", prefix, string(line)), aec.Faint))
		}
	}

	return lines
}

func empty(l []rune) bool {
	for _, r := range l {
		if r != ' ' {
			return false
		}
	}

	return true
}
