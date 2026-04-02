package scaffold

type StepStatus string

const (
	StatusPlanned   StepStatus = "planned"
	StatusCompleted StepStatus = "completed"
	StatusSkipped   StepStatus = "skipped"
	StatusFailed    StepStatus = "failed"
)

type StepResult struct {
	Name    string     `json:"name"`
	Kind    StepKind   `json:"kind"`
	Status  StepStatus `json:"status"`
	Path    string     `json:"path,omitempty"`
	Command []string   `json:"command,omitempty"`
	Message string     `json:"message,omitempty"`
}

type ApplyResult struct {
	Mode      string       `json:"mode"`
	Spec      CreateSpec   `json:"spec"`
	Succeeded bool         `json:"succeeded"`
	Steps     []StepResult `json:"steps"`
	Warnings  []string     `json:"warnings,omitempty"`
	NextSteps []string     `json:"next_steps,omitempty"`
}

func PlannedResult(plan Plan) ApplyResult {
	steps := make([]StepResult, 0, len(plan.Steps))
	for _, step := range plan.Steps {
		steps = append(steps, StepResult{
			Name:    step.Name,
			Kind:    step.Kind,
			Status:  StatusPlanned,
			Path:    step.Path,
			Command: step.Command,
			Message: step.Description,
		})
	}
	return ApplyResult{
		Mode:      "plan",
		Spec:      plan.Spec,
		Succeeded: true,
		Steps:     steps,
	}
}
