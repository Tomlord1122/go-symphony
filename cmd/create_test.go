package cmd

import (
	"strings"
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

func TestFrontendChoiceFromSelection(t *testing.T) {
	if got := frontendChoiceFromSelection("SvelteKit"); got != flags.SvelteKitFrontend {
		t.Fatalf("expected sveltekit, got %s", got)
	}
	if got := frontendChoiceFromSelection("Next.js"); got != flags.NextJSFrontend {
		t.Fatalf("expected nextjs, got %s", got)
	}
	if got := frontendChoiceFromSelection("unknown"); got != flags.NoneFrontend {
		t.Fatalf("expected none, got %s", got)
	}
}

func TestNextStepLinesForSupabaseAndFrontend(t *testing.T) {
	project := &program.Project{
		ProjectName:     "example/app",
		DBDriver:        flags.Supabase,
		AdvancedOptions: map[string]bool{flags.Sqlc: true},
	}

	lines := nextStepLines(project, flags.SvelteKitFrontend, flags.LocalDB)
	joined := strings.Join(lines, "\n")

	checks := []string{"supabase status", "sqlc generate", "pnpm dev", "cd example/app-frontend"}
	for _, check := range checks {
		if !strings.Contains(joined, check) {
			t.Fatalf("expected output to contain %q, got %s", check, joined)
		}
	}
}

func TestBuildPlanSpecFromCommand(t *testing.T) {
	cmd := &cobra.Command{Use: "plan"}
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("driver", "", "")
	cmd.Flags().String("git", "", "")
	cmd.Flags().String("feature", "", "")
	cmd.Flags().String("supabase-mode", "", "")
	cmd.Flags().String("frontend", "", "")
	cmd.Flags().String("sveltekit-template", "", "")
	cmd.Flags().String("sveltekit-types", "", "")
	cmd.Flags().String("sveltekit-package-manager", "", "")
	cmd.Flags().Bool("skip-install", false, "")
	cmd.Flags().String("output", "json", "")

	_ = cmd.Flags().Set("name", "example/app")
	_ = cmd.Flags().Set("driver", "postgres")
	_ = cmd.Flags().Set("git", "skip")
	_ = cmd.Flags().Set("feature", "sqlc,docker")
	_ = cmd.Flags().Set("frontend", "none")

	spec, err := buildPlanSpecFromCommand(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if spec.ProjectName != "example/app" {
		t.Fatalf("expected project name to round-trip, got %s", spec.ProjectName)
	}
	if !spec.HasFeature(flags.Sqlc) || !spec.HasFeature(flags.Docker) {
		t.Fatal("expected features to be included in plan spec")
	}
	if spec.Execution.Output != scaffold.OutputFormat("json") {
		t.Fatalf("expected json output, got %s", spec.Execution.Output)
	}
}
