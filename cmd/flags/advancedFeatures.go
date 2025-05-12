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
	Docker            string = "docker"
)

var AllowedAdvancedFeatures = []string{string(GoProjectWorkflow), string(Websocket), string(Docker)}

func (f AdvancedFeatures) String() string {
	return strings.Join(f, ",")
}

func (f *AdvancedFeatures) Type() string {
	return "AdvancedFeatures"
}

func (f *AdvancedFeatures) Set(value string) error {
	if slices.Contains(AllowedAdvancedFeatures, value) {
		*f = append(*f, value)
		return nil
	}
	return fmt.Errorf("advanced Feature to use. Allowed values: %s", strings.Join(AllowedAdvancedFeatures, ", "))
}
