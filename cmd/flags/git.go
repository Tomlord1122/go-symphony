package flags

import (
	"fmt"
	"slices"
	"strings"
)

type Git string

const (
	Commit Git = "commit"
	Skip   Git = "skip"
)

var AllowedGitsOptions = []string{string(Commit), string(Skip)}

func (g Git) String() string {
	return string(g)
}

func (g *Git) Type() string {
	return "Git"
}

func (g *Git) Set(value string) error {
	if !slices.Contains(AllowedGitsOptions, value) {
		return fmt.Errorf("Git option to use. Allowed values: %s", strings.Join(AllowedGitsOptions, ", "))
	}

	// Explicit type assertion to convert string to Git
	*g = Git(value)
	return nil
}
