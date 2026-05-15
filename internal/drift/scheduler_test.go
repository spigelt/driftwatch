package drift

import (
	"context"
	"sync"
	"testing"
	"time"
)

// stubRunner is a Runner replacement for tests.
type stubRunner struct {
	mu      sync.Mutex
	calls   int
	reports []Report
	err     error
}

func (s *stubRunner) Run(_ context.Context) ([]Report, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls++
	return s.reports, s.err
}

func TestScheduler_CallsOnDriftWhenDriftPresent(t *testing.T) {
	stub := &stubRunner{
		reports: []Report{
			makeReport("web", []Difference{{Field: "image", Compose: "nginx:1.24", Running: "nginx:1.25"}}),
		},
	}

	var mu sync.Mutex
	var received []Report

	onDrift := func(reports []Report) {
		mu.Lock()
		defer mu.Unlock()
		received = append(received, reports...)
	}

	sched := &Scheduler{
		runner:   &Runner{run: stub.Run},
		interval: 20 * time.Millisecond,
		onDrift:  onDrift,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	sched.Start(ctx)

	mu.Lock()
	defer mu.Unlock()
	if len(received) == 0 {
		t.Fatal("expected onDrift to be called at least once")
	}
	if received[0].Service != "web" {
		t.Errorf("expected service 'web', got %q", received[0].Service)
	}
}

func TestScheduler_NoDriftNoCallback(t *testing.T) {
	stub := &stubRunner{
		reports: []Report{
			makeReport("db", nil),
		},
	}

	called := false
	onDrift := func(_ []Report) { called = true }

	sched := &Scheduler{
		runner:   &Runner{run: stub.Run},
		interval: 20 * time.Millisecond,
		onDrift:  onDrift,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 55*time.Millisecond)
	defer cancel()
	sched.Start(ctx)

	if called {
		t.Error("onDrift should not be called when there is no drift")
	}
}
