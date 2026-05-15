package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
)

// ContainerInfo holds the subset of runtime data we care about for drift
// detection.
type ContainerInfo struct {
	ID      string
	Name    string
	Image   string
	Env     []string // "KEY=VALUE"
	Ports   []string // "hostPort:containerPort/proto"
	Labels  map[string]string
}

// InspectByName returns ContainerInfo for the first container whose name
// matches (Docker names are prefixed with "/").
func (dc *Client) InspectByName(ctx context.Context, name string) (*ContainerInfo, error) {
	data, err := dc.c.ContainerInspect(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("inspect %q: %w", name, err)
	}
	return fromInspect(data), nil
}

func fromInspect(data types.ContainerJSON) *ContainerInfo {
	ci := &ContainerInfo{
		ID:     data.ID,
		Name:   data.Name,
		Image:  data.Config.Image,
		Env:    data.Config.Env,
		Labels: data.Config.Labels,
	}

	for _, pb := range data.HostConfig.PortBindings {
		for _, b := range pb {
			ci.Ports = append(ci.Ports, b.HostPort)
		}
	}
	return ci
}
