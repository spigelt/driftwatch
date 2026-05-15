package compose_test

import (
	"testing"

	"github.com/yourusername/driftwatch/internal/compose"
)

func TestNormalise_SortsEnvironment(t *testing.T) {
	svc := compose.ServiceDefinition{
		Image: "nginx:latest",
		Environment: map[string]string{
			"Z_VAR": "z",
			"A_VAR": "a",
			"M_VAR": "m",
		},
	}
	n := compose.Normalise(svc)
	if len(n.Environment) != 3 {
		t.Fatalf("expected 3 env entries, got %d", len(n.Environment))
	}
	if n.Environment[0] != "A_VAR=a" {
		t.Errorf("expected first entry A_VAR=a, got %q", n.Environment[0])
	}
	if n.Environment[2] != "Z_VAR=z" {
		t.Errorf("expected last entry Z_VAR=z, got %q", n.Environment[2])
	}
}

func TestNormalise_SortsPorts(t *testing.T) {
	svc := compose.ServiceDefinition{
		Image: "app:1.0",
		Ports: []string{"9000:9000", "443:443", "80:80"},
	}
	n := compose.Normalise(svc)
	if n.Ports[0] != "443:443" {
		t.Errorf("expected first port 443:443, got %q", n.Ports[0])
	}
}

func TestNormalise_TrimsImage(t *testing.T) {
	svc := compose.ServiceDefinition{Image: "  redis:7  "}
	n := compose.Normalise(svc)
	if n.Image != "redis:7" {
		t.Errorf("expected trimmed image, got %q", n.Image)
	}
}

func TestNormaliseAll(t *testing.T) {
	cf := &compose.ComposeFile{
		Services: map[string]compose.ServiceDefinition{
			"web": {Image: "nginx:1.25"},
			"db":  {Image: "postgres:15"},
		},
	}
	result := compose.NormaliseAll(cf)
	if len(result) != 2 {
		t.Errorf("expected 2 normalised services, got %d", len(result))
	}
	if result["web"].Image != "nginx:1.25" {
		t.Errorf("unexpected image for web: %q", result["web"].Image)
	}
}
