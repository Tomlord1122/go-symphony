package flags

import (
	"fmt"
	"slices"
	"strings"
)

type AdvancedFeatures []string

const (
	GoProjectWorkflow string = "githubaction"
	Websocket         string = "websocket"
	Tailwind          string = "tailwind"
	Docker            string = "docker"
)

var AllowedAdvancedFeatures = []string{string(GoProjectWorkflow), string(Websocket), string(Tailwind), string(Docker)}

func (f AdvancedFeatures) String() string {
	return strings.Join(f, ",")
}

func (f *AdvancedFeatures) Type() string {
	return "AdvancedFeatures"
}

func (f *AdvancedFeatures) Set(value string) error {
	if !slices.Contains(AllowedAdvancedFeatures, value) {
		return fmt.Errorf("advanced Feature to use. Allowed values: %s", strings.Join(AllowedAdvancedFeatures, ", "))
	}

	// Explicit type assertion to convert string to AdvancedFeatures
	*f = append(*f, value)
	return nil
}
