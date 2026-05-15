package drift

import (
	"context"
	"fmt"
	"time"

	"github.com/yourorg/driftwatch/internal/compose"
	"github.com/yourorg/driftwatch/internal/docker"
)

// Runner orchestrates a single drift-check cycle.
type Runner struct {
	client     *docker.Client
	composePath string
}

// NewRunner creates a Runner for the given compose file path.
func NewRunner(client *docker.Client, composePath string) *Runner {
	return &Runner{client: client, composePath: composePath}
}

// Run executes a full drift check and returns a populated Report.
func (r *Runner) Run(ctx context.Context) (*Report, error) {
	svcs, err := compose.ParseFile(r.composePath)
	if err != nil {
		return nil, fmt.Errorf("parse compose: %w", err)
	}
	compose.NormaliseAll(svcs)

	report := &Report{Timestamp: time.Now()}

	for name, svc := range svcs {
		info, err := r.client.Inspect(ctx, name)
		if err != nil {
			// Container not running — treat as full drift on image field.
			report.Results = append(report.Results, Result{
				Service: name,
				Drifts: []Drift{{Field: "running", Expected: "true", Actual: "false"}},
			})
			continue
		}
		docker.NormaliseInfo(info)

		drifts := Compare(svc, info)
		report.Results = append(report.Results, Result{
			Service: name,
			Drifts:  drifts,
		})
	}

	return report, nil
}
