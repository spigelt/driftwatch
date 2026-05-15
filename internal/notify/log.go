package notify

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/driftwatch/internal/drift"
)

// LogSink writes structured drift reports to an io.Writer in a
// human-readable, line-oriented format suitable for log aggregation.
type LogSink struct {
	out io.Writer
}

// NewLogSink returns a LogSink writing to w (defaults to os.Stderr).
func NewLogSink(w io.Writer) *LogSink {
	if w == nil {
		w = os.Stderr
	}
	return &LogSink{out: w}
}

// Write emits one line per Drift entry in the report.
// If the report has no drift, a single "ok" line is emitted.
func (l *LogSink) Write(report drift.Report) error {
	ts := time.Now().UTC().Format(time.RFC3339)

	if !report.HasDrift() {
		_, err := fmt.Fprintf(l.out, "%s level=info container=%s status=ok\n", ts, report.Container)
		return err
	}

	for _, d := range report.Drifts {
		_, err := fmt.Fprintf(
			l.out,
			"%s level=warn container=%s field=%s compose=%q running=%q\n",
			ts,
			report.Container,
			d.Field,
			d.Compose,
			d.Running,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
