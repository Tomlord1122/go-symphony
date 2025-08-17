package flags

import (
	"fmt"
	"slices"
	"strings"
)

type SvelteKitTemplate string
type SvelteKitTypes string
type SvelteKitPackageManager string

const (
	// SvelteKit Templates
	Minimal SvelteKitTemplate = "minimal"
	Demo    SvelteKitTemplate = "demo"
	Library SvelteKitTemplate = "library"

	// Type checking options
	TypeScript SvelteKitTypes = "ts"
	JSDoc      SvelteKitTypes = "jsdoc"
	NoTypes    SvelteKitTypes = "none"

	// Package managers
	NPM  SvelteKitPackageManager = "npm"
	Yarn SvelteKitPackageManager = "yarn"
	PNPM SvelteKitPackageManager = "pnpm"
	Bun  SvelteKitPackageManager = "bun"
	Deno SvelteKitPackageManager = "deno"
)

var AllowedSvelteKitTemplates = []string{string(Minimal), string(Demo), string(Library)}
var AllowedSvelteKitTypes = []string{string(TypeScript), string(JSDoc), string(NoTypes)}
var AllowedSvelteKitPackageManagers = []string{string(NPM), string(Yarn), string(PNPM), string(Bun), string(Deno)}

// SvelteKitTemplate methods
func (s SvelteKitTemplate) String() string {
	return string(s)
}

func (s *SvelteKitTemplate) Type() string {
	return "SvelteKitTemplate"
}

func (s *SvelteKitTemplate) Set(value string) error {
	if slices.Contains(AllowedSvelteKitTemplates, value) {
		*s = SvelteKitTemplate(value)
		return nil
	}
	return fmt.Errorf("SvelteKit template to use. Allowed values: %s", strings.Join(AllowedSvelteKitTemplates, ", "))
}

// SvelteKitTypes methods
func (s SvelteKitTypes) String() string {
	return string(s)
}

func (s *SvelteKitTypes) Type() string {
	return "SvelteKitTypes"
}

func (s *SvelteKitTypes) Set(value string) error {
	if slices.Contains(AllowedSvelteKitTypes, value) {
		*s = SvelteKitTypes(value)
		return nil
	}
	return fmt.Errorf("SvelteKit type checking to use. Allowed values: %s", strings.Join(AllowedSvelteKitTypes, ", "))
}

// SvelteKitPackageManager methods
func (s SvelteKitPackageManager) String() string {
	return string(s)
}

func (s *SvelteKitPackageManager) Type() string {
	return "SvelteKitPackageManager"
}

func (s *SvelteKitPackageManager) Set(value string) error {
	if slices.Contains(AllowedSvelteKitPackageManagers, value) {
		*s = SvelteKitPackageManager(value)
		return nil
	}
	return fmt.Errorf("SvelteKit package manager to use. Allowed values: %s", strings.Join(AllowedSvelteKitPackageManagers, ", "))
}
