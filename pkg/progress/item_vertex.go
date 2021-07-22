package progress

import (
	"context"
	"fmt"
	"strings"

	api "github.com/moby/buildkit/api/services/control"
	"github.com/morikuni/aec"
)

type vertexItem struct {
	t      timer
	vertex *api.Vertex
	hidden bool
}

func newVertexItem(v *api.Vertex) *vertexItem {
	return &vertexItem{
		t:      newTimer(v.Started, v.Completed),
		vertex: v,
		hidden: false,
	}
}

func (v *vertexItem) hide() bool {
	// Allow hiding completed items
	if done, _ := v.t.Completed(); done {
		v.hidden = true

		return true
	}

	return false
}

func (v *vertexItem) height() int {
	return 1
}

func (v *vertexItem) lines(width int) []string {
	if v.hidden {
		return []string{}
	}

	if outStr, ok := formatLineWithTimer(width, v.t, v.format); ok {
		return []string{outStr}
	}

	return []string{}
}

func (v *vertexItem) format() (string, aec.ANSI) {
	color := aec.BlueF
	status := ""

	if v.vertex.Error != "" {
		if strings.HasSuffix(v.vertex.Error, context.Canceled.Error()) {
			status = statusCanceled
			color = aec.YellowF
		} else {
			status = statusError
			color = aec.RedF
		}
	} else if v.vertex.Cached {
		status = statusCached
	}

	outStr := strings.ReplaceAll(v.vertex.Name, "\t", " ")

	if len(status) > 0 {
		outStr = fmt.Sprintf("%s %s", status, outStr)
	}

	return prefix + outStr, color
}
