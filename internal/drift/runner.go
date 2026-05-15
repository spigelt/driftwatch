package drift

import (
	"context"
	"fmt"

	"github.com/yourorg/driftwatch/internal/compose"
	"github.com/yourorg/driftwatch/internal/docker"
)

// runFn abstracts the core run logic for testability.
type runFn func(ctx context.Context) ([]Report, error)

// Runner ties together compose parsing, docker inspection, and drift comparison.
type Runner struct {
	client      *docker.Client
	composeFile string
	// run is the internal implementation; swapped in tests via stubRunner.
	run runFn
}

// NewRunner constructs a Runner for the given compose file.
func NewRunner(client *docker.Client, composeFile string) *Runner {
	r := &Runner{
		client:      client,
		composeFile: composeFile,
	}
	r.run = r.runImpl
	return r
}

// Run executes a single drift-check cycle and returns one Report per service.
func (r *Runner) Run(ctx context.Context) ([]Report, error) {
	return r.run(ctx)
}

func (r *Runner) runImpl(ctx context.Context) ([]Report, error) {
	project, err := compose.ParseFile(r.composeFile)
	if err != nil {
		return nil, fmt.Errorf("parse compose: %w", err)
	}
	compose.NormaliseAll(project)

	var reports []Report
	for name, svc := range project.Services {
		info, err := r.client.Inspect(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("inspect %s: %w", name, err)
		}
		docker.NormaliseInfo(info)
		reports = append(reports, Compare(name, svc, *info))
	}
	return reports, nil
}
