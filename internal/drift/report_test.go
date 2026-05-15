package drift

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeReport(drifts map[string][]Drift) *Report {
	r := &Report{Timestamp: time.Now()}
	for svc, dd := range drifts {
		r.Results = append(r.Results, Result{Service: svc, Drifts: dd})
	}
	return r
}

func TestReport_HasDrift_False(t *testing.T) {
	r := makeReport(map[string][]Drift{
		"web": {},
	})
	if r.HasDrift() {
		t.Error("expected no drift")
	}
}

func TestReport_HasDrift_True(t *testing.T) {
	r := makeReport(map[string][]Drift{
		"web": {{Field: "image", Expected: "nginx:1.25", Actual: "nginx:1.24"}},
	})
	if !r.HasDrift() {
		t.Error("expected drift to be detected")
	}
}

func TestReport_Summary_NoDrift(t *testing.T) {
	r := makeReport(map[string][]Drift{})
	if got := r.Summary(); got != "No drift detected." {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestReport_Summary_WithDrift(t *testing.T) {
	r := makeReport(map[string][]Drift{
		"api": {
			{Field: "image", Expected: "myapp:v2", Actual: "myapp:v1"},
		},
	})
	summary := r.Summary()
	if !strings.Contains(summary, "api") {
		t.Error("summary should mention service name")
	}
	if !strings.Contains(summary, "image") {
		t.Error("summary should mention drifted field")
	}
}

func TestReport_WriteTo(t *testing.T) {
	r := makeReport(map[string][]Drift{
		"db": {{Field: "env", Expected: "PORT=5432", Actual: "PORT=5433"}},
	})
	var buf bytes.Buffer
	n, err := r.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo error: %v", err)
	}
	if n == 0 {
		t.Error("expected non-zero bytes written")
	}
	if !strings.Contains(buf.String(), "DriftWatch Report") {
		t.Error("output should contain report header")
	}
}
