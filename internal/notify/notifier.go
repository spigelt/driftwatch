package notify

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/driftwatch/internal/drift"
)

// Level controls the verbosity of notifications.
type Level int

const (
	LevelInfo  Level = iota
	LevelWarn
	LevelError
)

// Notifier sends drift reports to one or more sinks.
type Notifier struct {
	out   io.Writer
	level Level
}

// New returns a Notifier that writes to w.
func New(w io.Writer, level Level) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w, level: level}
}

// Notify formats and writes a drift report to the configured sink.
// It is a no-op when report.HasDrift() is false and level < LevelInfo.
func (n *Notifier) Notify(report drift.Report) error {
	if !report.HasDrift() && n.level > LevelInfo {
		return nil
	}

	timestamp := time.Now().UTC().Format(time.RFC3339)
	_, err := fmt.Fprintf(
		n.out,
		"[%s] %s\n",
		timestamp,
		report.Summary(),
	)
	return err
}

// NotifyAll calls Notify for each report in the slice.
func (n *Notifier) NotifyAll(reports []drift.Report) error {
	for _, r := range reports {
		if err := n.Notify(r); err != nil {
			return err
		}
	}
	return nil
}
