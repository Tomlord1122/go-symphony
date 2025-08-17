package flags

import (
	"fmt"
	"slices"
	"strings"
)

type FrontendFramework string

const (
	SvelteKitFrontend FrontendFramework = "sveltekit"
	NextJSFrontend    FrontendFramework = "nextjs"
	NoneFrontend      FrontendFramework = "none"
)

var AllowedFrontendFrameworks = []string{string(SvelteKitFrontend), string(NextJSFrontend), string(NoneFrontend)}

func (f FrontendFramework) String() string {
	return string(f)
}

func (f *FrontendFramework) Type() string {
	return "FrontendFramework"
}

func (f *FrontendFramework) Set(value string) error {
	if slices.Contains(AllowedFrontendFrameworks, value) {
		*f = FrontendFramework(value)
		return nil
	}
	return fmt.Errorf("frontend framework to use. Allowed values: %s", strings.Join(AllowedFrontendFrameworks, ", "))
}
