package flags

import (
	"fmt"
	"slices"
	"strings"
)

type Git string

const (
	Commit = "commit"
	Skip   = "skip"
)

var AllowedGitsOptions = []string{string(Commit), string(Skip)}

func (f Git) String() string {
	return string(f)
}

func (f *Git) Type() string {
	return "Git"
}

func (f *Git) Set(value string) error {
	if slices.Contains(AllowedGitsOptions, value) {
		*f = Git(value)
		return nil
	}

	return fmt.Errorf("Git to use. Allowed values: %s", strings.Join(AllowedGitsOptions, ", "))
}
