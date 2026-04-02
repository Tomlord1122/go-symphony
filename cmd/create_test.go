package cmd

import (
	"testing"

	"github.com/Tomlord1122/go-symphony/cmd/flags"
	"github.com/Tomlord1122/go-symphony/cmd/program"
	"github.com/Tomlord1122/go-symphony/internal/scaffold"
	"github.com/spf13/cobra"
)

func TestBuildScaffoldSpecUsesCommandFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("frontend", "", "")
	cmd.Flags().String("supabase-mode", "", "")
	cmd.Flags().String("sveltekit-template", "", "")
	cmd.Flags().String("sveltekit-types", "", "")
	cmd.Flags().String("sveltekit-package-manager", "", "")

	if err := cmd.Flags().Set("frontend", "sveltekit"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("supabase-mode", "init-only"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("sveltekit-template", "minimal"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("sveltekit-types", "ts"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("sveltekit-package-manager", "pnpm"); err != nil {
		t.Fatal(err)
	}

	project := &program.Project{
		ProjectName:     "example/app",
		DBDriver:        flags.Supabase,
		GitOptions:      flags.Skip,
		AdvancedOptions: map[string]bool{flags.Sqlc: true},
	}

	spec := buildScaffoldSpec(cmd, project, map[string]bool{}, scaffold.ExecutionOptions{NoInteractive: true})

	if spec.Frontend.Framework != flags.SvelteKitFrontend {
		t.Fatalf("expected sveltekit frontend, got %s", spec.Frontend.Framework)
	}
	if spec.SupabaseMode != flags.InitOnly {
		t.Fatalf("expected init-only supabase mode, got %s", spec.SupabaseMode)
	}
	if !spec.HasFeature(flags.Sqlc) {
		t.Fatal("expected sqlc feature to be included")
	}
}

func TestMapKeysMergesEnabledChoices(t *testing.T) {
	got := mapKeys(
		map[string]bool{"Docker": true, "Websocket": false},
		map[string]bool{"sqlc": true, "docker": true},
	)

	seen := map[string]bool{}
	for _, key := range got {
		seen[key] = true
	}

	if !seen["docker"] {
		t.Fatal("expected docker key")
	}
	if !seen["sqlc"] {
		t.Fatal("expected sqlc key")
	}
	if seen["websocket"] {
		t.Fatal("did not expect disabled websocket key")
	}
}
