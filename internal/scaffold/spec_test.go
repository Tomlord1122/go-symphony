package scaffold

import (
	"testing"

	"github.com/Tomlord1122/go-symphony/cmd/flags"
)

func TestValidateSpecNoInteractiveRequiresFields(t *testing.T) {
	spec := BuildSpec("", "", nil, "", "", "", "", "", "", ExecutionOptions{NoInteractive: true})
	if err := ValidateSpec(spec); err == nil {
		t.Fatal("expected validation error for missing fields")
	}
}

func TestValidateSpecRejectsSupabaseModeWithoutSupabase(t *testing.T) {
	spec := BuildSpec("example/app", flags.Postgres, nil, flags.Skip, flags.LocalDB, flags.NoneFrontend, "", "", "", ExecutionOptions{NoInteractive: true})
	if err := ValidateSpec(spec); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidateSpecAcceptsSensibleNonInteractiveSpec(t *testing.T) {
	spec := BuildSpec("example/app", flags.Postgres, []string{flags.Sqlc}, flags.Skip, "", flags.NoneFrontend, "", "", "", ExecutionOptions{NoInteractive: true})
	if err := ValidateSpec(spec); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestBuildPlanIncludesInstallStepsUnlessSkipped(t *testing.T) {
	spec := BuildSpec("example/app", flags.Postgres, []string{flags.Sqlc}, flags.Skip, "", flags.NoneFrontend, "", "", "", ExecutionOptions{})
	plan := BuildPlan(spec, "/tmp")
	if len(plan.Steps) == 0 {
		t.Fatal("expected non-empty plan")
	}

	spec.Execution.SkipInstall = true
	plan = BuildPlan(spec, "/tmp")
	for _, step := range plan.Steps {
		if step.Name == "Run go mod tidy" {
			t.Fatal("did not expect install step when skip-install is enabled")
		}
	}
}
