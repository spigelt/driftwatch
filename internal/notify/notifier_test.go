package notify_test

import (
	"strings"
	"testing"

	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/notify"
)

func makeReport(container string, drifts []drift.Drift) drift.Report {
	return drift.Report{
		Container: container,
		Drifts:    drifts,
	}
}

func TestNotifier_NoDrift_InfoLevel(t *testing.T) {
	var buf strings.Builder
	n := notify.New(&buf, notify.LevelInfo)

	r := makeReport("web", nil)
	if err := n.Notify(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "web") {
		t.Errorf("expected container name in output, got: %q", buf.String())
	}
}

func TestNotifier_NoDrift_WarnLevel_Suppressed(t *testing.T) {
	var buf strings.Builder
	n := notify.New(&buf, notify.LevelWarn)

	r := makeReport("web", nil)
	_ = n.Notify(r)

	if buf.Len() != 0 {
		t.Errorf("expected no output for no-drift at warn level, got: %q", buf.String())
	}
}

func TestNotifier_WithDrift_Written(t *testing.T) {
	var buf strings.Builder
	n := notify.New(&buf, notify.LevelWarn)

	r := makeReport("api", []drift.Drift{
		{Field: "image", Compose: "nginx:1.24", Running: "nginx:1.23"},
	})
	if err := n.Notify(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "api") {
		t.Errorf("expected container name in output, got: %q", out)
	}
}

func TestNotifyAll_MultipleReports(t *testing.T) {
	var buf strings.Builder
	n := notify.New(&buf, notify.LevelInfo)

	reports := []drift.Report{
		makeReport("web", nil),
		makeReport("db", []drift.Drift{{Field: "image", Compose: "postgres:15", Running: "postgres:14"}}),
	}

	if err := n.NotifyAll(reports); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "web") || !strings.Contains(out, "db") {
		t.Errorf("expected both container names in output, got: %q", out)
	}
}
