package drift

import (
	"fmt"
	"strings"
)

// ServiceInfo holds the normalised configuration for a single service,
// sourced from either a compose file or a running container.
type ServiceInfo struct {
	Image       string
	Environment []string
	Ports       []string
}

// Diff describes a single field-level difference between declared and actual state.
type Diff struct {
	Field    string
	Declared string
	Actual   string
}

// Result is the drift report for one service.
type Result struct {
	Service string
	Drifted bool
	Diffs   []Diff
}

// Compare returns a Result describing any drift between the declared compose
// service config and the running container's config.
func Compare(service string, declared, actual ServiceInfo) Result {
	var diffs []Diff

	if declared.Image != actual.Image {
		diffs = append(diffs, Diff{
			Field:    "image",
			Declared: declared.Image,
			Actual:   actual.Image,
		})
	}

	if d := sliceDiff("environment", declared.Environment, actual.Environment); d != nil {
		diffs = append(diffs, *d)
	}

	if d := sliceDiff("ports", declared.Ports, actual.Ports); d != nil {
		diffs = append(diffs, *d)
	}

	return Result{
		Service: service,
		Drifted: len(diffs) > 0,
		Diffs:   diffs,
	}
}

// sliceDiff compares two sorted string slices and returns a Diff when they differ.
func sliceDiff(field string, declared, actual []string) *Diff {
	declaredStr := strings.Join(declared, ", ")
	actualStr := strings.Join(actual, ", ")
	if declaredStr == actualStr {
		return nil
	}
	return &Diff{
		Field:    field,
		Declared: fmt.Sprintf("[%s]", declaredStr),
		Actual:   fmt.Sprintf("[%s]", actualStr),
	}
}
