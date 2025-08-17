package flags

import (
	"fmt"
	"slices"
	"strings"
)

type SupabaseMode string

const (
	InitOnly SupabaseMode = "init-only"
	LocalDB  SupabaseMode = "local-db"
)

var AllowedSupabaseModes = []string{string(InitOnly), string(LocalDB)}

func (s SupabaseMode) String() string {
	return string(s)
}

func (s *SupabaseMode) Type() string {
	return "SupabaseMode"
}

func (s *SupabaseMode) Set(value string) error {
	if slices.Contains(AllowedSupabaseModes, value) {
		*s = SupabaseMode(value)
		return nil
	}

	return fmt.Errorf("supabase mode to use. Allowed values: %s", strings.Join(AllowedSupabaseModes, ", "))
}
