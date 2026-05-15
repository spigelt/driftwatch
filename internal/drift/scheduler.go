package drift

import (
	"context"
	"log"
	"time"
)

// Scheduler periodically runs drift checks at a configured interval.
type Scheduler struct {
	runner   *Runner
	interval time.Duration
	onDrift  func([]Report)
}

// NewScheduler creates a Scheduler that triggers the runner every interval.
// onDrift is called with any reports that contain drift.
func NewScheduler(runner *Runner, interval time.Duration, onDrift func([]Report)) *Scheduler {
	return &Scheduler{
		runner:   runner,
		interval: interval,
		onDrift:  onDrift,
	}
}

// Start begins the scheduling loop, blocking until ctx is cancelled.
func (s *Scheduler) Start(ctx context.Context) {
	log.Printf("[scheduler] starting drift checks every %s", s.interval)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Run immediately before waiting for first tick.
	s.runOnce(ctx)

	for {
		select {
		case <-ticker.C:
			s.runOnce(ctx)
		case <-ctx.Done():
			log.Println("[scheduler] stopping")
			return
		}
	}
}

func (s *Scheduler) runOnce(ctx context.Context) {
	reports, err := s.runner.Run(ctx)
	if err != nil {
		log.Printf("[scheduler] run error: %v", err)
		return
	}

	var drifted []Report
	for _, r := range reports {
		if r.HasDrift() {
			drifted = append(drifted, r)
		}
	}

	if len(drifted) > 0 && s.onDrift != nil {
		s.onDrift(drifted)
	}
}
