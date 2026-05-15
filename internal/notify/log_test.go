package notify_test

import (
	"strings"
	"testing"

	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/notify"
)

func TestLogSink_NoDrift_EmitsOK(t *testing.T) {
	var buf strings.Builder
	s := notify.NewLogSink(&buf)

	if err := s.Write(makeReport("web", nil)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "status=ok") {
		t.Errorf("expected status=ok, got: %q", out)
	}
	if !strings.Contains(out, "container=web") {
		t.Errorf("expected container=web, got: %q", out)
	}
}

func TestLogSink_WithDrift_EmitsWarnLines(t *testing.T) {
	var buf strings.Builder
	s := notify.NewLogSink(&buf)

	r := makeReport("api", []drift.Drift{
		{Field: "image", Compose: "nginx:1.24", Running: "nginx:1.23"},
		{Field: "env", Compose: "DEBUG=true", Running: ""},
	})

	if err := s.Write(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %q", len(lines), out)
	}
	for _, line := range lines {
		if !strings.Contains(line, "level=warn") {
			t.Errorf("expected level=warn in line: %q", line)
		}
	}
}

func TestLogSink_FieldsPresent(t *testing.T) {
	var buf strings.Builder
	s := notify.NewLogSink(&buf)

	r := makeReport("db", []drift.Drift{
		{Field: "image", Compose: "postgres:15", Running: "postgres:14"},
	})
	_ = s.Write(r)

	out := buf.String()
	for _, want := range []string{"container=db", "field=image", "postgres:15", "postgres:14"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %q", want, out)
		}
	}
}
