package docker

import (
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

func TestFromInspect_BasicFields(t *testing.T) {
	data := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID:   "abc123",
			Name: "/myservice",
			HostConfig: &container.HostConfig{
				PortBindings: nat.PortMap{
					"8080/tcp": []nat.PortBinding{{HostPort: "8080"}},
				},
			},
		},
		Config: &container.Config{
			Image:  "nginx:latest",
			Env:    []string{"FOO=bar", "BAZ=qux"},
			Labels: map[string]string{"app": "web"},
		},
	}

	ci := fromInspect(data)

	if ci.ID != "abc123" {
		t.Errorf("ID: got %q, want %q", ci.ID, "abc123")
	}
	if ci.Image != "nginx:latest" {
		t.Errorf("Image: got %q, want %q", ci.Image, "nginx:latest")
	}
	if len(ci.Env) != 2 {
		t.Errorf("Env length: got %d, want 2", len(ci.Env))
	}
	if len(ci.Ports) != 1 || ci.Ports[0] != "8080" {
		t.Errorf("Ports: got %v, want [8080]", ci.Ports)
	}
	if ci.Labels["app"] != "web" {
		t.Errorf("Labels[app]: got %q, want %q", ci.Labels["app"], "web")
	}
}

func TestFromInspect_NoPorts(t *testing.T) {
	data := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID:         "def456",
			Name:       "/worker",
			HostConfig: &container.HostConfig{},
		},
		Config: &container.Config{
			Image: "alpine:3.18",
		},
	}

	ci := fromInspect(data)
	if len(ci.Ports) != 0 {
		t.Errorf("expected no ports, got %v", ci.Ports)
	}
}
