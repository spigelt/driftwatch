package docker

import (
	"sort"
	"strings"
)

// NormaliseInfo sorts and trims fields in ContainerInfo so that comparisons
// against compose-derived data are deterministic.
func NormaliseInfo(ci *ContainerInfo) {
	if ci == nil {
		return
	}

	ci.Image = strings.TrimSpace(ci.Image)
	// Ensure image has an explicit tag so "nginx" == "nginx:latest" after
	// compose normalisation.
	if ci.Image != "" && !strings.Contains(ci.Image, ":") {
		ci.Image += ":latest"
	}

	sort.Strings(ci.Env)
	sort.Strings(ci.Ports)
}
