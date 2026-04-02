package scaffold

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Tomlord1122/go-symphony/v2/cmd/flags"
)

func TestPlannedResultBuildsPlannedSteps(t *testing.T) {
	plan := BuildPlan(BuildSpec("example/app", flags.Postgres, nil, flags.Skip, "", flags.NoneFrontend, "", "", "", ExecutionOptions{}), "/tmp")
	result := PlannedResult(plan)

	if result.Mode != "plan" {
		t.Fatalf("expected plan mode, got %s", result.Mode)
	}
	if len(result.Steps) != len(plan.Steps) {
		t.Fatalf("expected %d steps, got %d", len(plan.Steps), len(result.Steps))
	}
	if result.Steps[0].Status != StatusPlanned {
		t.Fatalf("expected planned status, got %s", result.Steps[0].Status)
	}
}

func TestWriteApplyResultText(t *testing.T) {
	result := ApplyResult{
		Mode:      "apply",
		Succeeded: true,
		Steps: []StepResult{{
			Name:   "Initialize Go module",
			Kind:   StepRunCommand,
			Status: StatusCompleted,
		}},
	}

	var buf bytes.Buffer
	if err := WriteApplyResult(&buf, result, OutputText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Mode: apply") || !strings.Contains(output, "[completed] Initialize Go module") {
		t.Fatalf("unexpected output: %s", output)
	}
}
