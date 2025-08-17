package flags

import (
	"fmt"
	"slices"
	"strings"
)

type Framework string

// Only Gin framework is supported
const (
	Gin Framework = "gin"
)

var AllowedProjectTypes = []string{string(Gin)}

func (f Framework) String() string {
	return string(f)
}

func (f *Framework) Type() string {
	return "Framework"
}

func (f *Framework) Set(value string) error {
	if slices.Contains(AllowedProjectTypes, value) {
		*f = Framework(value)
		return nil
	}

	return fmt.Errorf("Framework to use. Allowed values: %s", strings.Join(AllowedProjectTypes, ", "))
}
