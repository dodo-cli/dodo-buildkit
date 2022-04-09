package progress

import (
	"fmt"

	"github.com/morikuni/aec"
)

func formatLineWithTimer(width int, t timer, format func() (string, aec.ANSI)) (string, bool) {
	running, runTime := t.Running()
	if !running {
		return "", false
	}

	timeStr := fmt.Sprintf(" %3.1fs", runTime.Seconds())
	outLimit := width - len(timeStr) - 1
	outStr, color := format()

	if len(outStr) > outLimit {
		outStr = outStr[:outLimit]
	}

	outStr = fmt.Sprintf("%-[2]*[1]s %[3]s", outStr, outLimit, timeStr)

	if done, _ := t.Completed(); done {
		outStr = aec.Apply(outStr, color)
	}

	return outStr, true
}
