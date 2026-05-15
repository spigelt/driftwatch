package drift

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// Report holds the results of a drift check across all watched containers.
type Report struct {
	Timestamp time.Time
	Results   []Result
}

// Result holds the drift comparison outcome for a single service.
type Result struct {
	Service string
	Drifts  []Drift
}

// HasDrift returns true if any service in the report has drifted.
func (r *Report) HasDrift() bool {
	for _, res := range r.Results {
		if len(res.Drifts) > 0 {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary of the report.
func (r *Report) Summary() string {
	if !r.HasDrift() {
		return "No drift detected."
	}
	var sb strings.Builder
	for _, res := range r.Results {
		if len(res.Drifts) == 0 {
			continue
		}
		fmt.Fprintf(&sb, "Service %q has %d drift(s):\n", res.Service, len(res.Drifts))
		for _, d := range res.Drifts {
			fmt.Fprintf(&sb, "  field=%s expected=%q actual=%q\n", d.Field, d.Expected, d.Actual)
		}
	}
	return sb.String()
}

// WriteTo writes the report in a structured text format to the given writer.
func (r *Report) WriteTo(w io.Writer) (int64, error) {
	s := fmt.Sprintf("DriftWatch Report — %s\n%s",
		r.Timestamp.Format(time.RFC3339),
		r.Summary(),
	)
	n, err := fmt.Fprint(w, s)
	return int64(n), err
}
