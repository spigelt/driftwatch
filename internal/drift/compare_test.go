package drift

import (
	"testing"
)

func TestCompare_NoDrift(t *testing.T) {
	declared := ServiceInfo{
		Image:       "nginx:1.25",
		Environment: []string{"ENV=prod", "PORT=80"},
		Ports:       []string{"80/tcp"},
	}
	actual := declared // identical

	result := Compare("web", declared, actual)

	if result.Drifted {
		t.Errorf("expected no drift, got diffs: %+v", result.Diffs)
	}
	if len(result.Diffs) != 0 {
		t.Errorf("expected 0 diffs, got %d", len(result.Diffs))
	}
}

func TestCompare_ImageDrift(t *testing.T) {
	declared := ServiceInfo{Image: "nginx:1.25"}
	actual := ServiceInfo{Image: "nginx:1.24"}

	result := Compare("web", declared, actual)

	if !result.Drifted {
		t.Fatal("expected drift but got none")
	}
	if len(result.Diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(result.Diffs))
	}
	if result.Diffs[0].Field != "image" {
		t.Errorf("expected field 'image', got %q", result.Diffs[0].Field)
	}
}

func TestCompare_EnvironmentDrift(t *testing.T) {
	declared := ServiceInfo{
		Image:       "app:latest",
		Environment: []string{"DEBUG=false", "ENV=prod"},
	}
	actual := ServiceInfo{
		Image:       "app:latest",
		Environment: []string{"DEBUG=true", "ENV=prod"},
	}

	result := Compare("app", declared, actual)

	if !result.Drifted {
		t.Fatal("expected drift but got none")
	}
	found := false
	for _, d := range result.Diffs {
		if d.Field == "environment" {
			found = true
		}
	}
	if !found {
		t.Error("expected an 'environment' diff")
	}
}

func TestCompare_MultipleDrifts(t *testing.T) {
	declared := ServiceInfo{
		Image: "redis:7",
		Ports: []string{"6379/tcp"},
	}
	actual := ServiceInfo{
		Image: "redis:6",
		Ports: []string{"6380/tcp"},
	}

	result := Compare("cache", declared, actual)

	if !result.Drifted {
		t.Fatal("expected drift")
	}
	if len(result.Diffs) != 2 {
		t.Errorf("expected 2 diffs, got %d: %+v", len(result.Diffs), result.Diffs)
	}
}
