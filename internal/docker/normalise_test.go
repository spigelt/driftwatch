package docker

import (
	"testing"
)

func TestNormaliseInfo_SortsEnvAndPorts(t *testing.T) {
	ci := &ContainerInfo{
		Image: "redis",
		Env:   []string{"Z=last", "A=first"},
		Ports: []string{"9000", "8000"},
	}

	NormaliseInfo(ci)

	if ci.Image != "redis:latest" {
		t.Errorf("Image: got %q, want %q", ci.Image, "redis:latest")
	}
	if ci.Env[0] != "A=first" {
		t.Errorf("Env[0]: got %q, want %q", ci.Env[0], "A=first")
	}
	if ci.Ports[0] != "8000" {
		t.Errorf("Ports[0]: got %q, want %q", ci.Ports[0], "8000")
	}
}

func TestNormaliseInfo_ImageAlreadyTagged(t *testing.T) {
	ci := &ContainerInfo{Image: "nginx:1.25"}
	NormaliseInfo(ci)
	if ci.Image != "nginx:1.25" {
		t.Errorf("Image should not change: got %q", ci.Image)
	}
}

func TestNormaliseInfo_NilSafe(t *testing.T) {
	// Should not panic.
	NormaliseInfo(nil)
}

func TestNormaliseInfo_TrimsWhitespace(t *testing.T) {
	ci := &ContainerInfo{Image: "  postgres:15  "}
	NormaliseInfo(ci)
	if ci.Image != "postgres:15" {
		t.Errorf("Image: got %q, want %q", ci.Image, "postgres:15")
	}
}
