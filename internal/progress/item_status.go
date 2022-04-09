package progress

import (
	"fmt"

	api "github.com/moby/buildkit/api/services/control"
	"github.com/morikuni/aec"
	"github.com/tonistiigi/units"
)

type statusItem struct {
	t      timer
	status *api.VertexStatus
	hidden bool
}

func newStatusItem(s *api.VertexStatus) *statusItem {
	return &statusItem{
		t:      newTimer(s.Started, s.Completed),
		status: s,
		hidden: false,
	}
}

func (s *statusItem) hide() bool {
	// Allow hiding completed items
	if done, _ := s.t.Completed(); done {
		s.hidden = true

		return true
	}

	return false
}

func (s *statusItem) height() int {
	return 1
}

func (s *statusItem) lines(width int) []string {
	if s.hidden {
		return []string{}
	}

	if outStr, ok := formatLineWithTimer(width, s.t, s.format); ok {
		return []string{outStr}
	}

	return []string{}
}

func (s *statusItem) format() (string, aec.ANSI) {
	status := ""

	if s.status.Total != 0 {
		status = fmt.Sprintf("(%.2f / %.2f)", units.Bytes(s.status.Current), units.Bytes(s.status.Total))
	} else if s.status.Current != 0 {
		status = fmt.Sprintf("(%.2f)", units.Bytes(s.status.Current))
	}

	outStr := s.status.ID

	if len(status) > 0 {
		outStr = fmt.Sprintf("%s %s", status, outStr)
	}

	return prefix + outStr, aec.BlueF
}
