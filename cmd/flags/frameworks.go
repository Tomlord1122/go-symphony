package flags

import (
	"fmt"
	"slices"
	"strings"
)

type Framework string

// These are all the current frameworks supported. If you want to add one, you
// can simply copy and paste a line here. Do not forget to also add it into the
// AllowedProjectTypes slice too!
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
	// Contains isn't available in 1.20 yet
	// if AllowedProjectTypes.Contains(value) {

	if slices.Contains(AllowedProjectTypes, value) {
		// Explicit type assertion to convert string to Framework
		*f = Framework(value)
		return nil
	}
	return fmt.Errorf("Framework to use. Allowed values: %s", strings.Join(AllowedProjectTypes, ", "))
}
