package scaffold

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func WritePlan(w io.Writer, plan Plan, format OutputFormat) error {
	if format == OutputJSON {
		return WriteApplyResult(w, PlannedResult(plan), format)
	}

	_, err := fmt.Fprintf(w, "Planned project: %s\n", plan.Spec.ProjectName)
	if err != nil {
		return err
	}
	for i, step := range plan.Steps {
		line := fmt.Sprintf("%d. [%s] %s", i+1, step.Kind, step.Name)
		if len(step.Command) > 0 {
			line += fmt.Sprintf(" -> %s", strings.Join(step.Command, " "))
		}
		if step.Path != "" {
			line += fmt.Sprintf(" (%s)", step.Path)
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func WriteApplyResult(w io.Writer, result ApplyResult, format OutputFormat) error {
	if format == OutputJSON {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	}

	_, err := fmt.Fprintf(w, "Mode: %s\nSucceeded: %t\n", result.Mode, result.Succeeded)
	if err != nil {
		return err
	}
	for i, step := range result.Steps {
		line := fmt.Sprintf("%d. [%s] %s", i+1, step.Status, step.Name)
		if len(step.Command) > 0 {
			line += fmt.Sprintf(" -> %s", strings.Join(step.Command, " "))
		}
		if step.Path != "" {
			line += fmt.Sprintf(" (%s)", step.Path)
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}
