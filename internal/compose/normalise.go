package compose

import (
	"sort"
	"strings"
)

// NormalisedService is a canonical, comparable form of a service definition.
type NormalisedService struct {
	Image       string
	Environment []string // sorted KEY=VALUE pairs
	Ports       []string // sorted
	Volumes     []string // sorted
	Command     string
}

// Normalise converts a ServiceDefinition into a NormalisedService so that
// order-independent fields can be compared deterministically.
func Normalise(svc ServiceDefinition) NormalisedService {
	env := make([]string, 0, len(svc.Environment))
	for k, v := range svc.Environment {
		env = append(env, k+"="+v)
	}
	sort.Strings(env)

	ports := append([]string(nil), svc.Ports...)
	sort.Strings(ports)

	volumes := append([]string(nil), svc.Volumes...)
	sort.Strings(volumes)

	return NormalisedService{
		Image:       strings.TrimSpace(svc.Image),
		Environment: env,
		Ports:       ports,
		Volumes:     volumes,
		Command:     strings.TrimSpace(svc.Command),
	}
}

// NormaliseAll returns a map of service name → NormalisedService for every
// service declared in the ComposeFile.
func NormaliseAll(cf *ComposeFile) map[string]NormalisedService {
	out := make(map[string]NormalisedService, len(cf.Services))
	for name, svc := range cf.Services {
		out[name] = Normalise(svc)
	}
	return out
}
