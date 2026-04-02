package cmd

import (
	"bytes"
	"encoding/json"
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

func TestNextStepLinesOnlyShowsDockerRunWhenDockerFeatureEnabled(t *testing.T) {
	project := &program.Project{
		ProjectName:     "example/app",
		DBDriver:        flags.Postgres,
		AdvancedOptions: map[string]bool{},
	}
	withoutDocker := strings.Join(nextStepLines(project, flags.NoneFrontend, ""), "\n")
	if strings.Contains(withoutDocker, "make docker-run") {
		t.Fatal("did not expect docker-run without docker feature")
	}

	project.AdvancedOptions[string(flags.Docker)] = true
	withDocker := strings.Join(nextStepLines(project, flags.NoneFrontend, ""), "\n")
	if !strings.Contains(withDocker, "make docker-run") {
		t.Fatal("expected docker-run with docker feature")
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

func TestFeatureFlagsAreNormalizedIntoProjectAdvancedOptions(t *testing.T) {
	project := &program.Project{AdvancedOptions: map[string]bool{}}
	for _, key := range strings.Split("docker,sqlc", ",") {
		normalized := strings.ToLower(strings.TrimSpace(key))
		if normalized != "" {
			project.AdvancedOptions[normalized] = true
		}
	}

	if !project.AdvancedOptions[flags.Docker] {
		t.Fatal("expected docker feature to be enabled")
	}
	if !project.AdvancedOptions[flags.Sqlc] {
		t.Fatal("expected sqlc feature to be enabled")
	}
}

func TestWriteApplyResultJSONContainsSteps(t *testing.T) {
	result := scaffold.PlannedResult(scaffold.BuildPlan(
		scaffold.BuildSpec("example/app", flags.Postgres, []string{flags.Docker}, flags.Skip, "", flags.NoneFrontend, "", "", "", scaffold.ExecutionOptions{}),
		"/tmp",
	))
	result.Mode = "apply"

	var buf bytes.Buffer
	if err := scaffold.WriteApplyResult(&buf, result, scaffold.OutputJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	steps, ok := decoded["steps"].([]any)
	if !ok || len(steps) == 0 {
		t.Fatalf("expected non-empty steps, got %#v", decoded["steps"])
	}
}

func TestFeatureFlagsWorkWithoutAdvancedModePrompt(t *testing.T) {
	project := &program.Project{AdvancedOptions: map[string]bool{}}
	featureFlags := "docker,sqlc"
	if featureFlags != "" {
		for _, key := range strings.Split(featureFlags, ",") {
			normalized := strings.ToLower(strings.TrimSpace(key))
			if normalized != "" {
				project.AdvancedOptions[normalized] = true
			}
		}
	}

	if !project.AdvancedOptions[flags.Docker] || !project.AdvancedOptions[flags.Sqlc] {
		t.Fatal("expected explicit feature flags to apply without advanced prompt")
	}
}

func TestNextStepLinesForPostgresWithoutDockerDoesNotSuggestDockerRun(t *testing.T) {
	project := &program.Project{
		ProjectName:     "example/app",
		DBDriver:        flags.Postgres,
		AdvancedOptions: map[string]bool{},
	}

	joined := strings.Join(nextStepLines(project, flags.NoneFrontend, ""), "\n")
	if strings.Contains(joined, "make docker-run") {
		t.Fatalf("did not expect docker-run in next steps: %s", joined)
	}
}
