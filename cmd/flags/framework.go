package flags

import (
	"fmt"
	"slices"
	"strings"
)

type Framework string

// The current frameworks supported by Go Symphony
const (
	Chi             Framework = "chi"
	StandardLibrary Framework = "standard-library"
)

var AllowedProjectTypes = []string{string(Chi), string(StandardLibrary)}

func (f Framework) String() string {
	return string(f)
}

func (f *Framework) Type() string {
	return "Framework"
}

func (f *Framework) Set(value string) error {
	if !slices.Contains(AllowedProjectTypes, value) {
		return fmt.Errorf("Framework to use. Allowed values: %s", strings.Join(AllowedProjectTypes, ", "))
	}
	// Explicit type assertion to convert string to Framework
	*f = Framework(value)
	return nil
}
