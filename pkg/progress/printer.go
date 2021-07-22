package progress

import (
	"context"
	"fmt"
	"io"
	"strings"

	api "github.com/moby/buildkit/api/services/control"
	"github.com/morikuni/aec"
)

type Printer struct {
	writer io.Writer

	Height int
	Width  int

	vertexes *vertexStates
	cache    itemCache

	previousLines int
}

func NewPrinter(w io.Writer, height int, width int) *Printer {
	return &Printer{
		writer:   w,
		Height:   height,
		Width:    width,
		vertexes: newVertexStates(),
		cache:    newItemCache(),
	}
}

func (p *Printer) Update(s *api.StatusResponse) {
	for _, v := range s.Vertexes {
		p.vertexes.Add(v, p.Width)
		p.cache.clear(v.Digest)
	}

	for _, s := range s.Statuses {
		if v, ok := p.vertexes.Get(s.Vertex); ok {
			v.updateStatus(s)
			p.cache.clear(v.vtx.Digest)
		}
	}

	for _, l := range s.Logs {
		if v, ok := p.vertexes.Get(l.Vertex); ok {
			v.updateLogs(l)
			p.cache.clear(v.vtx.Digest)
		}
	}
}

func (p *Printer) Print(all bool) {
	p.linesUp(p.previousLines)

	fmt.Fprint(p.writer, aec.Hide)
	defer fmt.Fprint(p.writer, aec.Show)

	lines := p.printLines(all)

	if diff := p.previousLines - lines; diff > 0 {
		for i := 0; i < diff; i++ {
			p.println(strings.Repeat(" ", p.Width))
		}
		p.linesUp(diff)
	}

	p.previousLines = lines
}

func (p *Printer) PrintErrorLogs() {
	for _, v := range p.vertexes.List() {
		if v.vtx.Error == "" {
			continue
		}

		if strings.HasSuffix(v.vtx.Error, context.Canceled.Error()) {
			continue
		}

		p.println(errorDelim)
		p.println(fmt.Sprintf(" > %s:", v.vtx.Name))

		for _, l := range v.logs {
			p.println(string(l))
		}

		p.println(errorDelim)
	}
}

func (p *Printer) println(line string) {
	// We assume the terminal is in raw mode, so we always send CR+LF
	fmt.Fprintf(p.writer, "%s\r\n", line)
}

func (p *Printer) linesUp(n int) {
	fmt.Fprint(p.writer, aec.EmptyBuilder.Up(uint(n)).Column(0).ANSI)
}

func (p *Printer) printLines(all bool) int {
	var lines []string

	items := p.getItems()

	if !all {
		// Account for the prompt line and the last newline at the end
		limit := p.Height - 2
		items = filterItems(items, limit)
	}

	for _, it := range items {
		lines = append(lines, it.lines(p.Width)...)
	}

	for _, l := range lines {
		p.println(l)
	}

	return len(lines)
}

func (p *Printer) getItems() []item {
	result := []item{newHeaderItem(p.vertexes)}

	for _, v := range p.vertexes.List() {
		items := p.cache.get(v.vtx.Digest, func() []item {
			var lines []item

			lines = append(lines, newVertexItem(v.vtx))
			lines = append(lines, newLogItem(v))

			for _, s := range v.byTime {
				lines = append(lines, newStatusItem(s))
			}

			return lines
		})

		result = append(result, items...)
	}

	return result
}

func filterItems(items []item, height int) []item {
	usedHeight := 0

	for _, item := range items {
		usedHeight += item.height()
	}

	for _, item := range items {
		if usedHeight <= height {
			return items
		}

		if ok := item.hide(); ok {
			usedHeight -= item.height()
		}
	}

        // When we reach this point, we have hidden everything we can and only
        // show the bare minimum, even if it doesn't fit on the screen
        return items
}
